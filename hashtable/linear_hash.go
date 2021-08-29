package hashtable

import (
	"time"

	"github.com/OneOfOne/xxhash"
)

// minimum size for the slot array, to stop reclaming spaces
const LH_MIN_TABLE_SIZE = 6

// minimum length for buckets
const LH_MIN_BUCKET_SIZE = 6

// Linear hash differs from Linear Probe that this is a incremental, dynamic sized
// hash table every overflow operation will split and resize at most X keys from
// the table.
type LinearHash struct {
	recordCount uint64
	slotArray   []*lhBucket

	hasher func(string) uint64
	// n is the hash modulo number, doubles when split counter wrap
	// around
	n int

	splitPointer  int // position of next bucket to be split
	maxBucketSize int // determine overflow threshold
}

type lhBucket struct {
	key       string
	value     interface{}
	next      *lhBucket
	tombstone bool
}

// WithmaxBucketSize set the maximum bucket size
func WithmaxBucketSize(size int) LHTableOption {
	size = max(size, LH_MIN_BUCKET_SIZE)
	return func(l *LinearHash) {
		l.maxBucketSize = size
	}
}

type LHTableOption func(l *LinearHash)

func NewLinearHash(size int, options ...LHTableOption) (l *LinearHash) {
	size = max(size, LH_MIN_BUCKET_SIZE)
	seed := uint64(time.Now().UnixNano())
	l = &LinearHash{
		hasher: func(in string) uint64 {
			return xxhash.ChecksumString64S(in, seed)
		},
		slotArray:     make([]*lhBucket, size),
		n:             size,
		splitPointer:  0,
		maxBucketSize: LH_MIN_BUCKET_SIZE,
	}

	for _, f := range options {
		f(l)
	}
	return
}

func (l *LinearHash) Put(key string, value interface{}) {
	if l.insertBucket(key, value) {
		l.split()
	}
	l.recordCount++
	return
}

func (l *LinearHash) split() {
	l.slotArray = append(l.slotArray, (*lhBucket)(nil))

	// redistribute keys in the splitted bucket
	old := l.slotArray[l.splitPointer]
	l.slotArray[l.splitPointer] = nil
	l.splitPointer++
	for ; old != nil; old = old.next {
		l.insertBucket(old.key, old.value)
	}

	if l.splitPointer >= l.n {
		l.n = len(l.slotArray)
		l.splitPointer = 0
	}
}

func (l *LinearHash) insertBucket(key string, value interface{}) (shouldSplit bool) {
	index := l.hashFunc(key)
	i, f := 0, &l.slotArray[index]
	// Here is the Linus Good Taste Linked List
	for ; *f != nil; f, i = &(*f).next, i+1 {
	} // nil slot found
	*f = &lhBucket{key: key, value: value} // append new value to end
	return i > l.maxBucketSize
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
	l.recordCount--
	value, *item = (*item).value, (*item).next
	l.unsplit()
	return value, true
}

func (l *LinearHash) Size() uint64 { return l.recordCount }

func (l *LinearHash) hashFunc(key string) int {
	index := l.hasher(key)
	if (index % uint64(l.n)) < uint64(l.splitPointer) {
		return int(index % uint64(l.n*2))
	}
	return int(index % uint64(l.n))
}

func (l *LinearHash) unsplit() {
	if l.slotArray[len(l.slotArray)-1] != nil ||
		len(l.slotArray) < LH_MIN_BUCKET_SIZE {
		return // nothing to reclaim
	}
	// shorten
	l.slotArray = l.slotArray[:len(l.slotArray)-1]
	// move split pointer up
	if l.splitPointer--; l.splitPointer < 0 {
		l.n /= 2
		l.splitPointer = l.n / 2
	}
}

func (l *LinearHash) Stats() map[string]interface{} {
	bucketSizeSlize := []int{}
	for _, v := range l.slotArray {
		c := 0
		for ; v != nil; v, c = v.next, c+1 {
		}
		bucketSizeSlize = append(bucketSizeSlize, c)
	}
	return map[string]interface{}{
		"n":                 l.n,
		"capacity":          len(l.slotArray),
		"size":              int(l.recordCount),
		"split_pointer":     l.splitPointer,
		"max_bucket_size":   l.maxBucketSize,
		"bucket_size_stats": bucketSizeSlize,
	}
}
