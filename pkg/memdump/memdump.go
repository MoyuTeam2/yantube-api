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

	sysTotalMemory uint64
)

func init() {
	sysTotalMemory, _ = getSystemTotalMemory()
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

	if sysTotalMemory > 0 && float64(memStats.Alloc) >= float64(sysTotalMemory)*memUsedRatioLimit {
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
	keepHighMem := false
	for {
		configLock.RLock()
		checkItvl, dumpItvl := checkInterval, dumpInterval
		configLock.RUnlock()

		// wait for next check
		time.Sleep(checkItvl)

		// continue if no need to dump
		if !needMemDump() {
			keepHighMem = false
			continue
		}

		// continue if keeping high memory
		if keepHighMem {
			continue
		}
		keepHighMem = true

		// continue if dump too frequently
		if time.Since(lastDumpTime) < dumpItvl {
			continue
		}

		DumpMemStats()

		lastDumpTime = time.Now()

	}
}
