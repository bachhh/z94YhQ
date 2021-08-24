package hashtable

import (
	"time"

	"github.com/OneOfOne/xxhash"
)

var seed uint64

func init() {
	seed = uint64(time.Now().UnixNano())
}

type LinearProbe struct {
	recordCount uint64
	table       []*record
	hasher      func(string) uint64
	// when exceeding size*growthFactor, perform resize(), where table will grow by 1/growthFactor.
	growthFactor float32
}

type record struct {
	key       string
	value     interface{}
	tombstone bool
}

func NewLinearProbe(size int) (l *LinearProbe) {
	l = &LinearProbe{
		table: make([]*record, size),
		hasher: func(in string) uint64 {
			return xxhash.ChecksumString64S(in, seed)
		},
		// doubling the size keep the table at 1/2 filled,
		// this keep collision rate low
		growthFactor: 0.5,
	}
	return
}

func (l *LinearProbe) Put(key string, value interface{}) {
	if !l.shoudlResize() { // optimistic branching
		index := int(l.hashFunc(key))
		if l.table[index] == nil {
			l.table[index] = &record{key: key, value: value}
			l.recordCount++
			return
		}
		// probe at most len(table) positions
		for i := 1; i < len(l.table); i++ {
			k := (index + i) % len(l.table)
			if l.table[k] == nil {
				l.table[k] = &record{key: key, value: value}
				l.recordCount++
				return
			}
		} // still can't find a slot, resize table
	}
	l.resize()
	l.Put(key, value)
}

func (l *LinearProbe) Get(key string) interface{} {
	index := int(l.hashFunc(key))
	if l.table[index] == nil {
		return nil
	}
	if l.table[index].key == key {
		if !l.table[index].tombstone {
			return l.table[index].value
		}
		return nil // key found but marked as delete
	}

	// key occupied but record not matched, linear search until found empty slot
	for i := 0; i < len(l.table); i++ {
		index := int(l.hashFunc(key))
		if l.table[index] == nil {
			return nil
		}
		if l.table[index].key == key {
			if !l.table[index].tombstone {
				return l.table[index].value
			}
			return nil // key found but marked as delete
		}
	}
	return nil
}

func (l *LinearProbe) getRecord(key string) *record {
	index := int(l.hashFunc(key))
	if l.table[index] == nil {
		return nil
	}
	if l.table[index].key == key {
		return l.table[index]
	}

	// key occupied but record not matched, linear search until found empty slot
	for i := 0; i < len(l.table); i++ {
		index := int(l.hashFunc(key))
		if l.table[index] == nil {
			return nil
		}
		if l.table[index].key == key {
			return l.table[index]
		}
	}
	return nil
}

// Delete delete a key, return whether key exist, if yes, also return value return
func (l *LinearProbe) Delete(key string) (exist bool, value interface{}) {
	return
}

func (l *LinearProbe) Size() uint64 { return l.recordCount }

// shoudlResize return true if table needs resizing
func (l *LinearProbe) shoudlResize() bool {
	return l.recordCount > uint64(float32(len(l.table))*l.growthFactor)
}

func (l *LinearProbe) resize() {
	oldTable := l.table
	l.recordCount, l.table = 0, make([]*record, int(float32(len(l.table))*(1/l.growthFactor)))
	for _, record := range oldTable {
		if record == nil {
			continue
		}
		l.Put(record.key, record.value)
	}
}

func (l *LinearProbe) hashFunc(key string) (index uint64) {
	return l.hasher(key) % uint64(len(l.table))
}
