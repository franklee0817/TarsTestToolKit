package impl

import (
	"context"

	"github.com/franklee0817/t3k/backend/services/functest"
	"github.com/franklee0817/t3k/backend/services/perftest"
	"github.com/franklee0817/t3k/backend/tars-protocol/apitars"
)

// APIImpl servant implementation
type APIImpl struct {
}

func (imp *APIImpl) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Init servant init
func (imp *APIImpl) InitFramework(ctx context.Context) (ret apitars.SimpleResp, err error) {
	//initialize servant here:
	//...
	return
}

func (imp *APIImpl) DoPerfTest(tarsCtx context.Context, req *apitars.PerfTestReq) (ret apitars.PerfTestResp, err error) {
	return perftest.DoPerfTest(tarsCtx, req)
}

func (imp *APIImpl) DoFuncTest(tarsCtx context.Context) (ret apitars.FuncTestResp, err error) {
	return functest.DoFuncTest(tarsCtx)
}

func (imp *APIImpl) GetTestDetail(tarsCtx context.Context, testID uint32, timestamp uint32) (
	ret apitars.TestDetailResp, err error) {
	var finished = false
	finished, ret.PerfDetail, ret.ResUsage, err = perftest.GetTestDetail(tarsCtx, testID, timestamp)
	if err != nil {
		return ret, err
	}

	ret.Code = 0
	if finished == false {
		// code 1 表示running
		ret.Code = 1
	}
	ret.Msg = "succ"
	return ret, err
}

func (imp *APIImpl) GetTestHistories(tarsCtx context.Context, req *apitars.QueryTestHistoryReq) (ret apitars.QueryTestHistoryResp, err error) {
	total, perfTests, err := perftest.QueryHistories(tarsCtx, req.Page, req.PageSize)
	ret.Total = uint32(total)
	ret.Page = req.Page
	ret.Histories = perfTests

	return ret, err
}
