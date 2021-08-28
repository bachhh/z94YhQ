package hashtable

import (
	"fmt"
	"math/rand"
	"testing"

	// useful, bookmark
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearProbe(t *testing.T) {
	t.Run("TestInsert", func(t *testing.T) {
		h := NewLinearProbe(6)
		for i := 0; i < 3; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
	})

	t.Run("TestInsertAndGet", func(t *testing.T) {
		h := NewLinearProbe(6)
		count := 100
		for i := 0; i < count; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		for i := 0; i < count; i++ {
			value, _ := h.Get(fmt.Sprintf("%d", i))
			assert.Equalf(t, value, fmt.Sprintf("value %d", i), "value not match at key %d", i)
		}
	})

	t.Run("TestBadHashInsertAndGet", func(t *testing.T) {
		// TestBadHashInsertAndGet use worst hash func with 100% collision rate
		h := NewLinearProbe(3,
			LBWithHasher(func(string) uint64 {
				return 0
			}),
		)
		count := 100
		for i := 0; i < count; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		for i := 0; i < count; i++ {
			value, _ := h.Get(fmt.Sprintf("%d", i))
			assert.Equal(t, value, fmt.Sprintf("value %d", i))
		}
	})

	t.Run("TestInsertDeleteGet", func(t *testing.T) {
		h := NewLinearProbe(6)
		count := 100
		for i := 0; i < count; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		deleted := []int{}
		deletedCount := min(50, count/2)

		// delete non-exist key
		for i := 0; i < deletedCount; i++ {
			key := rand.Intn(count) + count // outside range [0, count]
			_, exist := h.Delete(fmt.Sprintf("%d", key))
			require.False(t, exist, "key %d should not exist, but found", key)
		}

		// delete exist key
		for i := 0; i < deletedCount; i++ {
			key := rand.Intn(min(50, count/2))
			deleted = append(deleted, key)
			val, exist := h.Delete(fmt.Sprintf("%d", key))
			require.True(t, exist, "key %d not found", key)
			require.Equalf(t, val, fmt.Sprintf("value %d", key), "value not match for key %d", key)
		}

		isDeleted := func(f int) bool {
			for i := range deleted {
				if deleted[i] == f {
					return true
				}
			}
			return false
		}
		for i := 0; i < count; i++ {
			val, ok := h.Get(fmt.Sprintf("%d", i))
			if isDeleted(i) {
				require.Falsef(t, ok, "%d should be deleted, but found", i)
				require.Equal(t, val, nil)
			} else {
				require.Truef(t, ok, "%d should be retained, but deleted", i)
				require.Equal(t, val, fmt.Sprintf("value %d", i))
			}
		}
	})
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
