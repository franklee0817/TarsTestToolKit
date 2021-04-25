package cpu

import (
	fetchertars "github.com/franklee0817/t3k/fetcher/tars-protocol/TarsTestToolKit"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

// Stat 获取cpu stat信息
func Stat() ([]fetchertars.CoreInfo, error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return nil, err
	}

	cpus := make([]fetchertars.CoreInfo, 0)
	for _, s := range stat.CPUStats {
		c := fetchertars.CoreInfo{
			Total: int64(s.User + s.Nice + s.System + s.Idle + s.IOWait + s.IRQ + s.SoftIRQ + s.Steal + s.Guest + s.GuestNice),
			Idle:  int64(s.Idle),
		}
		c.Used = c.Total - c.Idle
		cpus = append(cpus, c)
	}

	return cpus, nil
}
