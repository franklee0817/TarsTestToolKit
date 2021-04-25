package impl

import (
	"context"

	"github.com/franklee0817/t3k/fetcher/services/cpu"
	"github.com/franklee0817/t3k/fetcher/services/memory"
	fetchertars "github.com/franklee0817/t3k/fetcher/tars-protocol/TarsTestToolKit"
)

// FetcherImp servant implementation
type FetcherImp struct {
}

// Init servant init
func (imp *FetcherImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destory
func (imp *FetcherImp) Destroy() {
	//destroy servant here:
	//...
}

// FetchResInfo 获取系统资源
func (imp *FetcherImp) FetchResInfo(tarsCtx context.Context) (ret fetchertars.ResResp, err error) {
	memInfo, err := memory.Stat()
	if err != nil {
		return ret, err
	}

	cores, err := cpu.Stat()
	if err != nil {
		return ret, err
	}

	ret.MemInfo = *memInfo
	ret.Cores = cores

	return ret, nil
}
