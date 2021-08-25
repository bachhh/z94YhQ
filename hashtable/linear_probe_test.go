package hashtable

import (
	"fmt"
	"testing"

	"github.com/zeebo/assert"
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
		for i := 0; i < 100; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		for i := 0; i < 100; i++ {
			value, _ := h.Get(fmt.Sprintf("%d", i))
			assert.Equal(t, value, fmt.Sprintf("value %d", i))
		}
	})

	t.Run("TestBadHashInsertAndGet", func(t *testing.T) {
		// TestBadHashInsertAndGet use a bad hash function with high collision rate
		h := NewLinearProbe(3,
			WithHasher(func(string) uint64 {
				return 0
			}),
		)
		for i := 0; i < 100; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		for i := 0; i < 100; i++ {
			value, _ := h.Get(fmt.Sprintf("%d", i))
			assert.Equal(t, value, fmt.Sprintf("value %d", i))
		}
	})

	// TODO(bach) test delete and get
	// TODO(bach) test insert, delete,  and get interaction
}
