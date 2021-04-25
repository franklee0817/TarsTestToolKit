package memory

import (
	stat "github.com/mackerelio/go-osstat/memory"

	fetchertars "github.com/franklee0817/t3k/fetcher/tars-protocol/TarsTestToolKit"
)

// Stat 获取内存统计信息
func Stat() (*fetchertars.MemInfo, error) {
	s, err := stat.Get()
	if err != nil {
		return nil, err
	}

	ret := &fetchertars.MemInfo{
		Total:     int64(s.Total),
		Used:      int64(s.Used),
		Cached:    int64(s.Cached),
		Free:      int64(s.Free),
		Active:    int64(s.Active),
		Inactive:  int64(s.Inactive),
		SwapTotal: int64(s.SwapTotal),
		SwapUsed:  int64(s.SwapUsed),
		SwapFree:  int64(s.SwapFree),
	}

	return ret, nil
}
