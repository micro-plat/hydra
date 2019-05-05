package cpu

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total       float64 `json:"total"`
	Idle        float64 `json:"idle"`
	UsedPercent float64 `json:"percent"`
}

// GetInfo 获取当前系统CPU使用的情况数据
func GetInfo(t ...time.Duration) (useage Useage) {
	cpus, _ := cpu.Times(true)
	useage = Useage{}
	for _, value := range cpus {
		useage.Total += value.Total()
		useage.Idle += value.Idle
	}
	du := time.Second
	if len(t) > 0 {
		du = t[0]
	}
	upc, _ := cpu.Percent(du, true)
	var total float64
	for _, v := range upc {
		total += v
	}
	useage.UsedPercent = total / float64(len(upc))
	return
}
