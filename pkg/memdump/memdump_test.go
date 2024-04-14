package memdump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemdump(t *testing.T) {
	SetMemUsedBytesLimit(1 * 1024 * 1024) // 1MB
	stat, err := getMemStats()
	require.Nil(t, err)
	t.Logf("Alloc: %d, Sys: %d", stat.Alloc, stat.Sys)
	SetDumpToDir("/Users/icceey/Projects/yantube-api")
	require.False(t, needMemDump())
	data := make([]byte, 2*1024*1024) // 2MB
	t.Logf("Alloc: %d, Sys: %d", stat.Alloc, stat.Sys)
	require.True(t, needMemDump())
	err = DumpMemStats()
	require.Nil(t, err)
	for i := 0; i < 100; i++ {
		data[i] = 0
	}
}
