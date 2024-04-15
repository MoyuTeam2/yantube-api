//go:build linux

package memdump

import (
	"fmt"
	"os"
)

func getSystemTotalMemory() (uint64, error) {
	// read from /proc/meminfo
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalMemory uint64
	_, err = fmt.Fscanf(file, "MemTotal:%d kB", &totalMemory)
	if err != nil {
		return 0, err
	}

	return totalMemory * 1024, nil
}
