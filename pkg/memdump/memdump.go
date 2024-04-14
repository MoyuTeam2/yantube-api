package memdump

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	memUsedBytesLimit uint64 = 1024 * 1024 * 1024 // 1GB
	memUsedRatioLimit        = 0.8                // 80%
	dumpToDir                = "."
	checkInterval            = 5 * time.Second
	dumpInterval             = 1 * time.Minute

	configLock sync.RWMutex
)

func init() {
	go lookMemStats()
}

func SetMemUsedBytesLimit(limit uint64) {
	configLock.Lock()
	memUsedBytesLimit = limit
	configLock.Unlock()
}

func SetMemUsedRatioLimit(limit float64) {
	configLock.Lock()
	memUsedRatioLimit = limit
	configLock.Unlock()
}

func SetDumpToDir(dir string) {
	configLock.Lock()
	dumpToDir = dir
	configLock.Unlock()
}

func SetCheckInterval(interval time.Duration) {
	configLock.Lock()
	checkInterval = interval
	configLock.Unlock()
}

func SetDumpInterval(interval time.Duration) {
	configLock.Lock()
	dumpInterval = interval
	configLock.Unlock()
}

func getMemStats() (runtime.MemStats, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats, nil
}

func needMemDump() bool {
	memStats, err := getMemStats()
	if err != nil {
		return false
	}

	configLock.RLock()
	defer configLock.RUnlock()

	if memStats.Alloc >= memUsedBytesLimit {
		return true
	}

	if float64(memStats.Alloc) >= float64(memStats.Sys)*memUsedRatioLimit {
		return true
	}

	return false
}

func DumpMemStats() error {
	configLock.RLock()
	filename := fmt.Sprintf("%s/memdump_%v.prof", dumpToDir, time.Now().Format("20060102_150405"))
	configLock.RUnlock()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = pprof.Lookup("heap").WriteTo(f, 0)

	return err
}

func lookMemStats() {
	lastDumpTime := time.Now()
	for {
		if !needMemDump() {
			continue
		}

		configLock.RLock()
		checkItvl, dumpItvl := checkInterval, dumpInterval
		configLock.RUnlock()

		if time.Since(lastDumpTime) < dumpItvl {
			continue
		}

		DumpMemStats()

		lastDumpTime = time.Now()
		time.Sleep(checkItvl)
	}
}
