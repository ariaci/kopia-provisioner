package utils

import "iter"

func EqIterate[Slice ~[]E, E any](s Slice, eq func(a, b E) bool) iter.Seq[Slice] {
	return func(yield func(Slice) bool) {
		if len(s) == 0 {
			return
		}

		start := 0
		for i := 1; i < len(s); i++ {
			if !eq(s[i], s[start]) {
				if !yield(s[start:i:i]) {
					return
				}

				start = i
			}
		}

		yield(s[start:len(s):len(s)])
	}
}
