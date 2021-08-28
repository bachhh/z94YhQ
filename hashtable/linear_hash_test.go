package hashtable

import (
	"testing"

	_ "github.com/davecgh/go-spew/spew" // useful, bookmark
	"github.com/stretchr/testify/require"
)

func TestLinearHash(t *testing.T) {
	t.Run("TestInsert", func(t *testing.T) {
		h := NewLinearHash(6)
		count := 100
		for i := 0; i < count; i++ {
			h.Put(kv(i))
		}
	})

	t.Run("TestInsertAndGet", func(t *testing.T) {
		h := NewLinearHash(6)
		count := 100
		keys := []int{}
		for i := 0; i < count; i++ {
			key, value := kv(i)
			h.Put(key, value)
			keys = append(keys, h.hashFunc(key))
		}
		// spew.Dump(keys)
		// spew.Dump(len(h.slotArray))
		for i := 0; i < count; i++ {
			k, v := kv(i)
			value, ok := h.Get(k)
			require.Truef(t, ok, "key %s not found", k)
			require.Equalf(t, value, v, "value not match at key %s", k)
		}
	})

	//t.Run("TestInsertDeleteGet", func(t *testing.T) {
	//	h := NewLinearHash(6)
	//	count := 100
	//	for i := 0; i < count; i++ {
	//		h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
	//	}
	//	deleted := []int{}
	//	deletedCount := min(50, count/2)

	//	// delete non-exist key
	//	for i := 0; i < deletedCount; i++ {
	//		key := rand.Intn(count) + count // outside range [0, count]
	//		_, exist := h.Delete(fmt.Sprintf("%d", key))
	//		require.False(t, exist, "key %d should not exist, but found", key)
	//	}

	//	// delete exist key
	//	for i := 0; i < deletedCount; i++ {
	//		key := rand.Intn(min(50, count/2))
	//		deleted = append(deleted, key)
	//		val, exist := h.Delete(fmt.Sprintf("%d", key))
	//		require.True(t, exist, "key %d not found", key)
	//		require.Equalf(t, val, fmt.Sprintf("value %d", key), "value not match for key %d", key)
	//	}

	//	isDeleted := func(f int) bool {
	//		for i := range deleted {
	//			if deleted[i] == f {
	//				return true
	//			}
	//		}
	//		return false
	//	}
	//	for i := 0; i < count; i++ {
	//		val, ok := h.Get(fmt.Sprintf("%d", i))
	//		if isDeleted(i) {
	//			require.Falsef(t, ok, "%d should be deleted, but found", i)
	//			require.Equal(t, val, nil)
	//		} else {
	//			require.Truef(t, ok, "%d should be retained, but deleted", i)
	//			require.Equal(t, val, fmt.Sprintf("value %d", i))
	//		}
	//	}
	//})

	//t.Run("TestReclaim", func(t *testing.T) {
	//	h := NewLinearHash(6)
	//	// count := 100
	//	// delete aggressively,
	//	//	- test if table reclaim spaces correctly
	//	//	- test if get operation still perform correctly
	//	count := 100
	//	for i := 0; i < count; i++ {
	//		h.Put(fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i))
	//	}
	//	deleted := []int{}
	//	deletedCount := (count / 4) * 3 // delete 3/4 of table

	//	// delete non-exist key
	//	for i := 0; i < deletedCount; i++ {
	//		key := rand.Intn(count) + count // outside range [0, count]
	//		_, exist := h.Delete(fmt.Sprintf("%d", key))
	//		require.False(t, exist, "key %d should not exist, but found", key)
	//	}

	//	// delete exist key
	//	for i := 0; i < deletedCount; i++ {
	//		key := rand.Intn(min(50, count/2))
	//		deleted = append(deleted, key)
	//		val, exist := h.Delete(fmt.Sprintf("%d", key))
	//		require.True(t, exist, "key %d not found", key)
	//		require.Equalf(t, val, fmt.Sprintf("value %d", key), "value not match for key %d", key)
	//	}

	//	isDeleted := func(f int) bool {
	//		for i := range deleted {
	//			if deleted[i] == f {
	//				return true
	//			}
	//		}
	//		return false
	//	}
	//	for i := 0; i < count; i++ {
	//		val, ok := h.Get(fmt.Sprintf("%d", i))
	//		if isDeleted(i) {
	//			require.Falsef(t, ok, "%d should be deleted, but found", i)
	//			require.Equal(t, val, nil)
	//		} else {
	//			require.Truef(t, ok, "%d should be retained, but deleted", i)
	//			require.Equal(t, val, fmt.Sprintf("value %d", i))
	//		}
	//	}
	//	t.Logf("table size %d", len(h.slotArray))
	//	return
	//})

}
