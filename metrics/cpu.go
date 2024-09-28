package metrics

// CPUStats is the system and process CPU stats.
// All values are in seconds.
type CPUStats struct {
	GlobalTime float64 // Time spent by the CPU working on all processes
	GlobalWait float64 // Time spent by waiting on disk for all processes
	LocalTime  float64 // Time spent by the CPU working on this process
}
