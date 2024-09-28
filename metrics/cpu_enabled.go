//go:build !ios && !js
// +build !ios,!js

package metrics

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/shirou/gopsutil/cpu"
)

// ReadCPUStats retrieves the current CPU stats.
func ReadCPUStats(stats *CPUStats) {
	// passing false to request all cpu times
	timeStats, err := cpu.Times(false)
	if err != nil {
		log.Error("Could not read cpu stats", "err", err)
		return
	}
	if len(timeStats) == 0 {
		log.Error("Empty cpu stats")
		return
	}
	// requesting all cpu times will always return an array with only one time stats entry
	timeStat := timeStats[0]
	stats.GlobalTime = timeStat.User + timeStat.Nice + timeStat.System
	stats.GlobalWait = timeStat.Iowait
	stats.LocalTime = getProcessCPUTime()
}
