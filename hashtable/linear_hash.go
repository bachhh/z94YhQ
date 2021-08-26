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
	slotArray   []*lhBucket

	hasher func(string) uint64

	bucketCounter int
	splitPointer  int // position of next bucket to be split
	overflowTH    int // determine overflow threshold
}

type lhBucket struct {
	key       string
	value     interface{}
	next      *lhBucket
	tombstone bool
}

func WithOverflowTH(threshold int) LHTableOption {
	return func(l *LinearHash) {
		l.overflowTH = threshold
	}
}

type LHTableOption func(l *LinearHash)

func NewLinearHash(size int, options ...LHTableOption) (l *LinearHash) {
	l = &LinearHash{
		hasher: func(in string) uint64 {
			return xxhash.ChecksumString64S(in, uint64(time.Now().UnixNano()))
		},
		splitPointer: 0,
		overflowTH:   3,
	}
	for _, f := range options {
		f(l)
	}
	return
}

func (l *LinearHash) Put(key string, value interface{}) {
	index := l.hashFunc(key)
	if l.insertBucket(index, key, value) {
		l.split()
	}
	return
}

// Get return value assigned with key, return (nil, false) if key not found in table
func (l *LinearHash) Get(key string) (interface{}, bool) {
	for z := l.slotArray[l.hashFunc(key)]; z != nil; z = z.next {
		if z.key == key {
			return z.value, true
		}
	}
	return nil, false
}

// Delete delete a key, return whether key exist, if yes, also return value return
func (l *LinearHash) Delete(key string) (value interface{}, exist bool) {
	// Here is the Linus Good Taste Linked List
	var item **lhBucket = &l.slotArray[l.hashFunc(key)]
	for (*item) != nil && (*item).key != key {
		item = &(*item).next
	}
	if (*item) == nil {
		return nil, false
	}
	value, *item = (*item).value, (*item).next
	return value, true
}

func (l *LinearHash) Size() uint64 { return l.recordCount }

func (l *LinearHash) hashFunc(key string) (index uint64) {
	index = l.hasher(key)
	if index < uint64(l.splitPointer) {
		return index % uint64(len(l.slotArray))
	}
	return index % (2 * uint64(len(l.slotArray)))
}

func (l *LinearHash) insertBucket(index uint64, key string, value interface{}) (shouldSplit bool) {
	i, f := 0, l.slotArray[index]
	for ; f != nil; f, i = f.next, i+1 {
	} // nil slot found
	*f = lhBucket{key: key, value: value} // append new value to end
	return i > l.overflowTH
}

func (l *LinearHash) split() {
	// 1. add new slot in array
	l.slotArray = append(l.slotArray, (*lhBucket)(nil))
	// 2. redistribute keys in the split bucket
	for old := l.slotArray[l.splitPointer]; old != nil; old = old.next {
		l.insertBucket(l.hashFunc(old.key), old.key, old.value)
	}
	// 3. increment split pointer
	l.splitPointer = (l.splitPointer + 1) % l.bucketCounter
}
