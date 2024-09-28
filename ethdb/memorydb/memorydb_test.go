package memorydb

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/dbtest"
)

func TestMemoryDB(t *testing.T) {
	t.Run("DatabaseSuite", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func() ethdb.KeyValueStore {
			return New()
		})
	})
}

// BenchmarkBatchAllocs measures the time/allocs for storing 120 kB of data
func BenchmarkBatchAllocs(b *testing.B) {
	b.ReportAllocs()
	var key = make([]byte, 20)
	var val = make([]byte, 100)
	// 120 * 1_000 -> 120_000 == 120kB
	for i := 0; i < b.N; i++ {
		batch := New().NewBatch()
		for j := uint64(0); j < 1000; j++ {
			binary.BigEndian.PutUint64(key, j)
			binary.BigEndian.PutUint64(val, j)
			batch.Put(key, val)
		}
		batch.Write()
	}
}
