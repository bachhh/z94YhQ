package hashtable

import (
	"time"

	"github.com/OneOfOne/xxhash"
)

// Linear hash differs from Linear Probe that this is a incremental, dynamic sized
// hash table every overflow operation will split and resize at most X keys from
// the table.

type LinearHash struct {
	recordCount uint64
	table       []*record

	hasher       func(string) uint64
	splitPointer int
	overflowTH   int // determine overflow threshold
}

type lhRecord struct {
	key       string
	value     interface{}
	tombstone bool
	next      *lhRecord
}

type LHTableOption func(l *LinearHash)

func LHWithHasher(hasher func(string) uint64) LHTableOption {
	return func(l *LinearHash) {
		l.hasher = hasher
	}
}

func NewLinearHash(size int, options ...LHTableOption) (l *LinearHash) {
	l = &LinearHash{
		hasher: func(in string) uint64 {
			return xxhash.ChecksumString64S(in, uint64(time.Now().UnixNano()))
		},
	}
	for _, f := range options {
		f(l)
	}
	return
}

func (l *LinearHash) Put(key string, value interface{}) {
}

// Get return value assigned with key, return (nil, false) if key not found in table
func (l *LinearHash) Get(key string) (interface{}, bool) {
	return nil, false
}

// Delete delete a key, return whether key exist, if yes, also return value return
func (l *LinearHash) Delete(key string) (exist bool, value interface{}) {
	return false, nil
}

func (l *LinearHash) Size() uint64 { return l.recordCount }
