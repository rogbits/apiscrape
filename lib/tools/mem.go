package tools

import (
	"fmt"
	"runtime"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", BytesToMegabytes(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", BytesToMegabytes(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", BytesToMegabytes(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func BytesToMegabytes(bytes uint64) uint64 {
	return bytes / 1024 / 1024
}
