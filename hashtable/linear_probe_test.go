package hashtable

import (
	"fmt"
	"testing"
)

func TestLinearProbe(t *testing.T) {
	t.Run("TestInsert", func(t *testing.T) {
		h := NewLinearProbe(3)
		for i := 0; i < 3; i++ {
			h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
		}
		t.Logf("%#v\n", h)
	})
}
