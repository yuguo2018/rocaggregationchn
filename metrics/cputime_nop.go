//go:build windows || js
// +build windows js

package metrics

// getProcessCPUTime returns 0 on Windows as there is no system call to resolve
// the actual process' CPU time.
func getProcessCPUTime() float64 {
	return 0
}
