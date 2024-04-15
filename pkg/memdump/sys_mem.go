//go:build !linux

package memdump

func getSystemTotalMemory() (uint64, error) {
	return 0, nil
}
