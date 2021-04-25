package benchmark

import (
	"context"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/franklee0817/t3k/backend/constants"
	"github.com/franklee0817/t3k/backend/constants/errors"
	"github.com/franklee0817/t3k/backend/tars-protocol/bm"
	"github.com/franklee0817/t3k/backend/tools/communicator"
)

func getAdminApp() *bm.Admin {
	app := new(bm.Admin)
	cmm := communicator.GetCommunicator()
	cmm.StringToProxy(constants.ServiceNameBMAdmin, app)

	return app
}

// test 对指定服务进行功能测试
func test(ctx context.Context, req *bm.BenchmarkUnit) int32 {
	app := getAdminApp()

	var resp, msg string
	ret, err := app.Test(req, &resp, &msg)
	tars.GetLogger("").Debugf("request servant:%s func:%s return:%v err:%v resp:%s msg:%s",
		req.Servant, req.Rpcfunc, ret, err, resp, msg)
	if err != nil {
		tars.GetLogger("").Errorf("request servant:%s func:%s got error:%s", req.Servant, req.Rpcfunc, err.Error())
	}

	return ret
}

func ping(ctx context.Context, servant string) int32 {
	req := &bm.BenchmarkUnit{
		Servant:    servant,
		Rpcfunc:    "ping",
		Para_input: `{"1_req_string":""}`,
		Para_value: `{"req": ""}`,
		Paralist:   []string{},
		Endpoints:  communicator.GetAllEndpoints(servant),
		Links:      0,
		Speed:      0,
		Duration:   0,
		Proto:      "json",
	}
	return test(ctx, req)
}

func Startup(ctx context.Context, req *bm.BenchmarkUnit) error {
	app := getAdminApp()

	var resp, msg string
	ret, err := app.Startup(req)
	tars.GetLogger("").Debugf("start benchmark servant:%s func:%s return:%v err:%v resp:%s msg:%s",
		req.Servant, req.Rpcfunc, ret, err, resp, msg)
	if err != nil {
		tars.GetLogger("").Errorf("start benchmark servant:%s func:%s got error:%s", req.Servant, req.Rpcfunc, err.Error())
		return err
	}
	if ret != errors.BmSucc {
		return tars.Errorf(ret, "benchmark returns error code:%v", ret)
	}

	return nil
}

// PingCpp ping cpp 服务
func PingCpp(ctx context.Context) error {
	ret := ping(ctx, constants.ServiceNameCpp)
	if ret == errors.BmSucc {
		return nil
	}

	return tars.Errorf(ret, "ping cpp failed")
}

// PingJava ping java 服务
func PingJava(ctx context.Context) error {
	ret := ping(ctx, constants.ServiceNameJava)
	if ret == errors.BmSucc {
		return nil
	}

	return tars.Errorf(ret, "ping java failed")
}

// PingGo ping golang 服务
func PingGo(ctx context.Context) error {
	ret := ping(ctx, constants.ServiceNameGolang)
	if ret == errors.BmSucc {
		return nil
	}

	return tars.Errorf(ret, "ping golang failed")
}

// PingNodeJs ping nodejs 服务
func PingNodeJs(ctx context.Context) error {
	ret := ping(ctx, constants.ServiceNameNodeJs)
	if ret == errors.BmSucc {
		return nil
	}

	return tars.Errorf(ret, "ping nodejs failed")
}

// PingPhp ping php服务
func PingPhp(ctx context.Context) error {
	ret := ping(ctx, constants.ServiceNamePhp)
	if ret == errors.BmSucc {
		return nil
	}

	return tars.Errorf(ret, "ping php failed")
}

// Query 查询压测结果
func Query(ctx context.Context, req *bm.BenchmarkUnit) (*bm.ResultStat, error) {
	app := getAdminApp()

	resp := new(bm.ResultStat)

	ret, err := app.Query(req, resp)
	if err != nil {
		tars.GetLogger("").Errorf("failed to query benchmark result for test:%+v, error:%v", *req, err)
		return resp, err
	}
	if ret != errors.BmSucc {
		tars.GetLogger("").Errorf("failed to query benchmark result for test:%+v, ret:%v", *req, ret)
		return resp, tars.Errorf(ret, "failed to query benchmark result")
	}

	return resp, nil
}
