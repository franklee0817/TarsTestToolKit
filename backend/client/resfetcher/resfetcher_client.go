package resfetcher

import (
	"fmt"

	fetchertars "github.com/franklee0817/t3k/fetcher/tars-protocol/TarsTestToolKit"

	"github.com/franklee0817/t3k/backend/constants"
	"github.com/franklee0817/t3k/backend/tools/communicator"
)

const (
	ServiceName = "TarsTestToolKit.ResFetcher.fetcherObj"
)

// FetchResInfo 请求ResFetcher获取硬件信息
func FetchResInfo() ([]fetchertars.CoreInfo, *fetchertars.MemInfo, error) {
	app := new(fetchertars.Fetcher)
	communicator.StringToProxy(ServiceName, app)
	resp, err := app.FetchResInfo()
	if err != nil {
		return nil, nil, err
	}
	if resp.Code != 0 {
		return nil, nil, fmt.Errorf("request service:%s api:FetchResInfo failed, err:%s",
			constants.ServiceNameResFetcher, resp.Msg)
	}

	return resp.Cores, &resp.MemInfo, nil
}
