package hashtable

import "testing"

var a map[int]interface{}

func BenchmarkGoHash(b *testing.B) {
	a = map[int]interface{}{}
	for i := 0; i < b.N; i++ {
		a[i] = "test"
	}
}
