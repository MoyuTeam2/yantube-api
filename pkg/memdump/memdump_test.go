package memdump

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemdump(t *testing.T) {
	SetMemUsedBytesLimit(1 * 1024 * 1024) // 1MB
	SetDumpToDir(".")
	t.Log("PWD:", os.Getenv("PWD"))

	stat, err := getMemStats()
	require.Nil(t, err)
	t.Logf("Alloc: %d, Total: %d", stat.Alloc, sysTotalMemory)
	require.False(t, needMemDump())

	data := make([]byte, 2*1024*1024) // 2MB

	stat, err = getMemStats()
	require.Nil(t, err)
	t.Logf("Alloc: %d, Total: %d", stat.Alloc, sysTotalMemory)
	require.True(t, needMemDump())

	err = DumpMemStats()
	require.Nil(t, err)
	for i := 0; i < 100; i++ {
		data[i] = 0
	}
}
