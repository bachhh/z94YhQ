package hashtable

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// TODO(bach) test delete and get
	// TODO(bach) test insert, delete,  and get interaction
}
