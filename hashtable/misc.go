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

func debugLH(h *LinearHash) (s string) {
	for i, v := range h.slotArray {
		switch {
		case i == h.splitPointer && i == (h.n-1):
			s += fmt.Sprintf("*%d:\t", i)
		case i == (h.n - 1):
			s += fmt.Sprintf("n%d:\t", i)
		case i == h.splitPointer:
			s += fmt.Sprintf("+%d:\t", i)
		default:
			s += fmt.Sprintf(" %d:\t", i)
		}

		for ; v != nil; v = v.next {
			s += fmt.Sprintf("'%s'-", v.key)
		}
		s += "\n"
	}
	return
}
