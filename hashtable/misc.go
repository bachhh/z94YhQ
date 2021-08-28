package hashtable

import "fmt"

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func kv(i int) (string, string) {
	return fmt.Sprintf("%d", i), fmt.Sprintf("value %d", i)
}
