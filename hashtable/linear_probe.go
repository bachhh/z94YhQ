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

type TableOption func(l *LinearProbe)

func WithHasher(hasher func(string) uint64) TableOption {
	return func(l *LinearProbe) {
		l.hasher = hasher
	}
}

func WithGrowthFactor(factor float32) TableOption {
	return func(l *LinearProbe) {
		l.growthFactor = factor
	}
}

func NewLinearProbe(size int, options ...TableOption) (l *LinearProbe) {
	l = &LinearProbe{
		table: make([]*record, size),
		hasher: func(in string) uint64 {
			return xxhash.ChecksumString64S(in, seed)
		},
		// doubling the size keep the table at 1/2 filled,
		// this keep collision rate low
		growthFactor: 0.5,
	}
	for _, f := range options {
		f(l)
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
		} else if l.table[index] != nil && l.table[index].tombstone {
			// occupied but key tombstoned, overwrite
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

// Get return value assigned with key, return (nil, false) if key not found in table
func (l *LinearProbe) Get(key string) (interface{}, bool) {
	record := l.getRecord(key)
	if record == nil {
		return nil, false
	} else if !record.tombstone {
		return record.value, true
	}
	return nil, false
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
		k := (index + i) % len(l.table)
		if l.table[k] == nil {
			return nil
		}
		if l.table[k].key == key { // keep searching otherwise
			return l.table[k]
		}
	}
	return nil
}

// Delete delete a key, return whether key exist, if yes, also return value return
func (l *LinearProbe) Delete(key string) (exist bool, value interface{}) {
	record := l.getRecord(key)
	if record == nil {
		return false, nil
	}
	record.tombstone = true
	l.recordCount--
	return true, record.value
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
