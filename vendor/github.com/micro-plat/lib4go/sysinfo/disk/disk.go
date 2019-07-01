package disk

import (
	"runtime"

	"github.com/shirou/gopsutil/disk"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total       uint64  `json:"total"`
	Idle        uint64  `json:"idle"`
	UsedPercent float64 `json:"percent"`
}

// GetInfo 获取磁盘使用信息
func GetInfo() (useage Useage) {
	dir := "/"
	if runtime.GOOS == "windows" {
		dir = "c:"
	}
	sm, _ := disk.Usage(dir)

	useage.Total = sm.Total
	useage.Idle = sm.Total - sm.Used
	useage.UsedPercent = sm.UsedPercent
	return
}
