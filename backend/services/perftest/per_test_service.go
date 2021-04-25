package perftest

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	jsoniter "github.com/json-iterator/go"

	"github.com/franklee0817/t3k/backend/client/benchmark"
	"github.com/franklee0817/t3k/backend/client/tarsweb"
	"github.com/franklee0817/t3k/backend/constants"
	"github.com/franklee0817/t3k/backend/constants/errors"
	"github.com/franklee0817/t3k/backend/models/mysql"
	"github.com/franklee0817/t3k/backend/services/stats"
	"github.com/franklee0817/t3k/backend/tars-protocol/apitars"
	"github.com/franklee0817/t3k/backend/tars-protocol/bm"
	"github.com/franklee0817/t3k/backend/tools"
	"github.com/franklee0817/t3k/backend/tools/communicator"
)

// DoFuncTest 执行性能测试
func DoPerfTest(ctx context.Context, req *apitars.PerfTestReq) (apitars.PerfTestResp, error) {
	now := time.Now()
	ret := apitars.PerfTestResp{Code: 0, Msg: "succ"}
	lang := strings.ToLower(req.Lang)
	serv, ok := constants.LangMap[lang]
	if !ok {
		return ret, tars.Errorf(errors.ErrCodeParam, "unsupported performance test for language:%s", lang)
	}
	// 更新服务线程数
	err := prepareServerBeforeBM(ctx, serv, lang, int(req.ThreadCnt))
	if err != nil {
		tars.GetLogger("").Errorf(err.Error())
		return ret, err
	}

	// 开始压测
	err = StartBM(ctx, serv, req.PackageLen, int32(req.ConnCnt), int32(req.ReqFreq), int32(req.KeepAlive))
	if err != nil {
		ret.Code = uint32(tars.GetErrorCode(err))
		ret.Msg = err.Error()
		return ret, err
	}

	// 记录压测信息
	row := newPerfTest(req, now)
	_, err = row.Insert()
	if err != nil {
		tars.GetLogger("").Errorf("insert perf tests failed:%v", err)
		ret.Code = errors.ErrCodeInternalErr
		ret.Msg = err.Error()
		return ret, err
	}

	ret.TestID = uint32(row.ID)
	// 采集目标服务资源使用情况,多采集15秒
	endTime := now.Add(time.Duration(req.KeepAlive+15) * time.Second)
	stats.WatchStats(ctx, ret.TestID, endTime)
	WatchBM(ctx, ret.TestID, serv)

	return ret, err
}

func prepareServerBeforeBM(ctx context.Context, serv string, lang string, threadCnt int) error {
	adapterName := serv + "Adapter"
	err := mysql.SetAdapterThread(adapterName, threadCnt)
	if err != nil {
		tars.GetLogger("").Errorf(err.Error())
		return err
	}
	taskNo, err := restartService(ctx, constants.AppNameTestUnits, lang)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Millisecond * 500)
	for i := 0; i < 10; i++ {
		finished, err := tarsweb.IsTaskFinished(ctx, taskNo)
		if tars.GetErrorCode(err) == errors.ErrCodeInternalErr {
			return err
		}
		if finished {
			return nil
		}
		<-ticker.C
	}

	return tars.Errorf(errors.ErrCodeInternalErr, "restart service timeout")
}

func restartService(ctx context.Context, app, serv string) (string, error) {
	node, err := tarsweb.FindAppNode(ctx, app)
	if err != nil {
		return "", err
	}
	info, err := tarsweb.FindServerInfo(ctx, node.ID, app, serv)
	if err != nil {
		return "", err
	}

	return tarsweb.Restart(ctx, info.ID)
}

// WatchBM 持续收集压测信息
func WatchBM(ctx context.Context, tid uint32, serv string) {
	go doWatchBM(ctx, tid, serv)
}

func doWatchBM(ctx context.Context, tid uint32, serv string) {
	doneCh := make(chan interface{})
	t := time.NewTicker(5 * time.Second)
	for {
		<-t.C
		select {
		case <-doneCh:
			return
		default:
			FetchBM(ctx, tid, serv, "ping", doneCh)
		}
	}
}

// FetchBM 查询并保存压测结果
func FetchBM(ctx context.Context, tid uint32, serv string, fn string, doneCh chan interface{}) {
	unit := &bm.BenchmarkUnit{
		Servant: serv,
		Rpcfunc: fn,
		Proto:   "json",
	}
	result, err := benchmark.Query(ctx, unit)
	errCode := tars.GetErrorCode(err)
	finished := errCode == errors.BmAdminErrNotfind || errCode == errors.BmAdminErrRunning
	if err != nil && !finished {
		return
	}

	if result != nil {
		detail := buildPerfDetail(tid, result)
		tars.GetLogger("").Debugf("perf detail:%+v", detail)
		trySavePerfDetail(detail, finished)
	}
	if finished {
		close(doneCh)
	}
}

func trySavePerfDetail(perfDetail *mysql.PerfTestDetail, finished bool) {
	if err := mysql.InsertTestDetail(perfDetail, finished); err != nil {
		time.AfterFunc(1*time.Second, func() {
			trySavePerfDetail(perfDetail, finished)
		})
	}
}

func buildPerfDetail(tid uint32, result *bm.ResultStat) *mysql.PerfTestDetail {
	avgCost := float32(math.Round(result.Total_time/float64(result.Total_request)*100) / 100)
	succRate := float64(0)
	if result.Total_request > 0 {
		succRate = math.Round(float64(result.Succ_request) / float64(result.Total_request) * 10000)
	}
	retMap, _ := jsoniter.MarshalToString(result.Ret_map)
	costMap, _ := jsoniter.MarshalToString(result.Cost_map)
	perfDetail := &mysql.PerfTestDetail{
		TestID:     int(tid),
		QPS:        int(result.Avg_speed),
		Total:      int(result.Total_request),
		Succ:       int(result.Succ_request),
		Failed:     int(result.Fail_request),
		SuccRate:   int(succRate),
		CostMax:    float32(result.Max_time),
		CostMin:    float32(result.Min_time),
		CostAvg:    avgCost,
		P90:        float32(result.P90_time),
		P99:        float32(result.P99_time),
		P999:       float32(result.P999_time),
		Send:       result.Send_bytes,
		Recv:       result.Recv_bytes,
		CostMap:    costMap,
		RetMap:     retMap,
		CreateTime: int(result.Time_stamp),
	}

	return perfDetail
}

// StartBM 开始压测
func StartBM(ctx context.Context, serv string, pkgLen uint32, links, speed, dur int32) error {

	pkg := tools.GetStr(pkgLen)
	bmReq := &bm.BenchmarkUnit{
		Servant:    serv,
		Rpcfunc:    "ping",
		Para_input: `{"1_req_string":""}`,
		Para_value: fmt.Sprintf(`{"req": "%s"}`, pkg),
		Paralist:   []string{},
		Endpoints:  communicator.GetAllEndpoints(serv),
		Links:      links,
		Speed:      speed,
		Duration:   dur,
		Proto:      "json",
	}
	err := benchmark.Startup(ctx, bmReq)
	return err
}

func newPerfTest(req *apitars.PerfTestReq, now time.Time) *mysql.PerfTests {
	row := &mysql.PerfTests{
		ServType:  req.ServType,
		Lang:      strings.ToLower(req.Lang),
		ServName:  constants.LangMap[strings.ToLower(req.Lang)],
		FnName:    "ping",
		Cores:     int(req.Cores),
		Threads:   int(req.ThreadCnt),
		ConnCnt:   int(req.ConnCnt),
		Frequency: int(req.ReqFreq),
		KeepAlive: int(req.KeepAlive),
		PkgLen:    int(req.PackageLen),
		StartTime: int(now.Unix()),
		EndTime:   int(uint32(now.Unix()) + req.KeepAlive),
	}

	return row
}

// QueryHistories 查询压测历史信息
func QueryHistories(ctx context.Context, page, pageSize uint32) (int64, []apitars.TestHistory, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 10 || pageSize > 50 {
		pageSize = 10
	}
	total, rows, err := mysql.PaginatePerfTests(page, pageSize)
	if err != nil {
		tars.GetLogger("").Errorf("failed to query perf_tests:%v", err)
		return 0, nil, err
	}

	histories := make([]apitars.TestHistory, 0)
	for _, row := range rows {
		histories = append(histories, apitars.TestHistory{
			TestID:    uint32(row.ID),
			StartTime: uint32(row.StartTime),
			EndTime:   uint32(row.EndTime),
			ServType:  row.ServType,
			Lang:      row.Lang,
			Cores:     uint32(row.Cores),
			Threads:   uint32(row.Threads),
			ConnCnt:   uint32(row.ConnCnt),
			KeepAlive: uint32(row.KeepAlive),
			PkgLen:    uint32(row.PkgLen),
		})
	}

	return total, histories, nil
}

// GetTestDetail 查询测试详情
func GetTestDetail(ctx context.Context, tid, timestamp uint32) (bool, []apitars.PerfTestDetail, []apitars.PerfResDetail, error) {
	test, err := mysql.GetPerfTest(tid)
	if err != nil {
		return false, nil, nil, err
	}
	status := test.Finished == 1
	testDetail, err := mysql.QueryTestDetail(tid, timestamp)
	if err != nil {
		tars.GetLogger("").Errorf("query test_detail for %v failed %v", tid, err)
		return status, nil, nil, err
	}
	perfDetails := make([]apitars.PerfTestDetail, 0)
	for _, detail := range testDetail {
		perfDetails = append(perfDetails, buildPerfTestDetailFromDB(detail))
	}

	cpuStats, err := mysql.QueryCpuStats(tid)
	if err != nil {
		tars.GetLogger("").Errorf("query cpu_stats for %v failed %v", tid, err)
		return status, nil, nil, err
	}
	memStats, err := mysql.QueryMemStats(tid)
	if err != nil {
		tars.GetLogger("").Errorf("query mem_stats for %v failed %v", tid, err)
		return status, nil, nil, err
	}
	resDetails := buildResDetailFromDB(cpuStats, memStats)

	return status, perfDetails, resDetails, nil
}

func buildResDetailFromDB(cpuStats []mysql.CpuStats, memStats []mysql.MemStats) []apitars.PerfResDetail {
	cpuStatsMp := make(map[int64][]mysql.CpuStats)
	memStatsMp := make(map[int64]mysql.MemStats)
	for _, cpu := range cpuStats {
		if _, ok := cpuStatsMp[cpu.CreateTime]; !ok {
			cpuStatsMp[cpu.CreateTime] = make([]mysql.CpuStats, 0)
		}
		cpuStatsMp[cpu.CreateTime] = append(cpuStatsMp[cpu.CreateTime], cpu)
	}
	for _, mem := range memStats {
		memStatsMp[mem.CreateTime] = mem
	}
	ret := make([]apitars.PerfResDetail, 0)
	for ts, cpus := range cpuStatsMp {
		item := apitars.PerfResDetail{
			Timestamp: uint32(ts),
		}
		for _, cpu := range cpus {
			p := math.Round(float64(cpu.Used)/float64(cpu.Total)*100) / 100
			item.Cpu = append(item.Cpu, apitars.CoreUsage{Percent: float32(p)})
		}
		mem := memStatsMp[ts]
		item.Mem.Total = mem.Total
		item.Mem.Used = mem.Used
		ret = append(ret, item)
	}

	return ret
}

func buildPerfTestDetailFromDB(detail mysql.PerfTestDetail) apitars.PerfTestDetail {
	costMap := make(map[string]uint32)
	retMap := make(map[string]uint32)
	_ = jsoniter.UnmarshalFromString(detail.CostMap, &costMap)
	_ = jsoniter.UnmarshalFromString(detail.RetMap, &retMap)
	perfDetail := apitars.PerfTestDetail{
		Timestamp:  uint32(detail.CreateTime),
		Qps:        uint32(detail.QPS),
		TotalReq:   uint32(detail.Total),
		Succ:       uint32(detail.Succ),
		Failed:     uint32(detail.Failed),
		SuccRate:   fmt.Sprintf("%.2f%", float64(detail.SuccRate)/100),
		CostMax:    detail.CostMax,
		CostMin:    detail.CostMin,
		CostAvg:    detail.CostAvg,
		P90:        detail.P90,
		P99:        detail.P99,
		P999:       detail.P999,
		SendByte:   uint32(detail.Send),
		RecvByte:   uint32(detail.Recv),
		CostMap:    retMap,
		RetCodeMap: costMap,
	}

	return perfDetail
}
