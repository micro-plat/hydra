package middleware

import (
	"strconv"
	"time"

	"github.com/micro-plat/lib4go/metrics"
	"github.com/micro-plat/lib4go/sysinfo/cpu"
	"github.com/micro-plat/lib4go/sysinfo/disk"
	"github.com/micro-plat/lib4go/sysinfo/memory"
	"github.com/micro-plat/lib4go/sysinfo/pipes"
)

func (m *Metric) collectCPU() {
	name := metrics.MakeName("server.cpu.used.percent", metrics.GAUGE, "host", m.ip) //堵塞计数
	counter := metrics.GetOrRegisterGaugeFloat64(name, m.currentRegistry)
	cpuInfo := cpu.GetInfo(time.Millisecond * 200)
	counter.Update(cpuInfo.UsedPercent)
}
func (m *Metric) collectMem() {
	name := metrics.MakeName("server.memory.used.percent", metrics.GAUGE, "host", m.ip) //堵塞计数
	counter := metrics.GetOrRegisterGaugeFloat64(name, m.currentRegistry)
	memoryInfo := memory.GetInfo()
	counter.Update(memoryInfo.UsedPercent)
}
func (m *Metric) collectDisk() {
	name := metrics.MakeName("server.disk.used.percent", metrics.GAUGE, "host", m.ip) //堵塞计数
	counter := metrics.GetOrRegisterGaugeFloat64(name, m.currentRegistry)
	diskInfo := disk.GetInfo()
	counter.Update(diskInfo.UsedPercent)
}
func (m *Metric) collectNetConnectCNT() {
	name := metrics.MakeName("server.net.conn.counter", metrics.GAUGE, "host", m.ip) //堵塞计数
	counter := metrics.GetOrRegisterGaugeFloat64(name, m.currentRegistry)
	counter.Update(getNetConnectCount())
}
func (m *Metric) loopCollectCPU() {
	cpuChan := m.timer.Subscribe()
	for {
		select {
		case <-m.closeChan:
			return
		case <-cpuChan:
			m.collectCPU()
		}
	}
}
func (m *Metric) loopCollectMem() {
	cpuChan := m.timer.Subscribe()
	for {
		select {
		case <-m.closeChan:
			return
		case <-cpuChan:
			m.collectMem()
		}
	}
}
func (m *Metric) loopCollectDisk() {
	cpuChan := m.timer.Subscribe()
	for {
		select {
		case <-m.closeChan:
			return
		case <-cpuChan:
			m.collectDisk()
		}
	}
}
func (m *Metric) loopNetConnCount() {
	netChan := m.timer.Subscribe()
	for {
		select {
		case <-m.closeChan:
			return
		case <-netChan:
			m.collectNetConnectCNT()
		}
	}
}

//-----------------------------------基础函数---------------------------------
func getNetConnectCount() (v float64) {
	count, err := pipes.BashRun(`netstat -an|grep tcp|wc -l`)
	if err != nil {
		return 0
	}
	x, _ := strconv.Atoi(count)
	return float64(x)
}

func getMaxOpenFiles() float64 {
	count, err := pipes.BashRun("ulimit -n")
	if err != nil {
		return 0
	}
	v, _ := strconv.Atoi(count)
	return float64(v)
}
