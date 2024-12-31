package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap"
	"runtime"
	"time"
)

// updateCpuUsage 更新CPU使用率
//
// 参数:
//   - instance: 实例名称
//   - value: CPU使用率
func (m *metrics) updateCpuUsage(instance string, value float64) {
	m.cpuUsageGauge.WithLabelValues(instance).Set(value)
}

// updateMemoryUsage 更新内存使用率
//
// 参数:
//   - instance: 实例名称
//   - value: 内存使用率
func (m *metrics) updateMemoryUsage(instance string, value float64) {
	m.memoryUsageGauge.WithLabelValues(instance).Set(value)
}

// updateGoroutineNum 更新goroutine数量
//
// 参数:
//   - instance: 实例名称
//   - value: goroutine数量
func (m *metrics) updateGoroutineNum(instance string, value float64) {
	m.goroutineGauge.WithLabelValues(instance).Set(value)
}

// updateProcessNum 更新进程数量
//
// 参数:
//   - instance: 实例名称
//   - value: 进程数量
func (m *metrics) updateProcessNum(instance string, value float64) {
	m.processNumGauge.WithLabelValues(instance).Set(value)
}

func (m *metrics) AutoUpdateMetrics() {
	go func(*metrics) {
		for range time.Tick(time.Second) {
			// 此处需要保证 go 版本至少为 1.23, 以保证 time.Tick() 不会泄漏
			m.update()
		}
	}(m)
}

// update 更新指标
func (m *metrics) update() {
	// 更新CPU使用率
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		zap.L().Error("get cpu percent failed", zap.Error(err))
		m.updateCpuUsage("local", 0)
	} else {
		m.updateCpuUsage("local", cpuPercent[0])
	}

	// 更新内存使用率
	memoryPercent, err := mem.VirtualMemory()
	if err != nil {
		zap.L().Error("get memory percent failed", zap.Error(err))
		m.updateMemoryUsage("local", 0)
	} else {
		m.updateMemoryUsage("local", memoryPercent.UsedPercent)
	}

	// 更新goroutine数量
	m.updateGoroutineNum("local", float64(runtime.NumGoroutine()))

	// 更新进程数量
	processes, err := process.Processes()
	if err != nil {
		zap.L().Error("get process num failed", zap.Error(err))
		m.updateProcessNum("local", 0)
	} else {
		m.updateProcessNum("local", float64(len(processes)))
	}
}
