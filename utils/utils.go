package utils

import "sort"

func Copy(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func MCopyIntInt(src map[int]int) map[int]int {
	if src == nil {
		return nil
	}
	dst := make(map[int]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func MapToSlice(m map[int]int) []int {
	ret := make([]int, 0, len(m))
	for k, v := range m {
		for i := 0; i < v; i++ {
			ret = append(ret, k)
		}
	}
	return ret
}

// NOTE it is guaranteed that in2 <= in1
// Retval is new'd
func MSub(in1, in2 map[int]int) map[int]int {
	if in2 == nil {
		return MCopyIntInt(in1)
	}
	m := make(map[int]int, 0)
	for k, v := range in1 {
		if v > in2[k] {
			m[k] = v - in2[k]
		}
	}
	return m
}

// Retval stored in in1
func MSubInplace(in1, in2 map[int]int) {
	if in2 == nil {
		return
	}
	for k := range in1 {
		if in1[k] > in2[k] {
			in1[k] -= in2[k]
		} else {
			delete(in1, k)
		}
	}
	return
}

// Retval stored in in1
func MSubSliceInplace(in map[int]int, slice []map[int]int) {
	for _, s := range slice {
		MSubInplace(in, s)
	}
	return
}

func MapToSliceLen(m map[int]int) int {
	r := 0
	for _, v := range m {
		r += v
	}
	return r
}

func InSlice(i int, slice []int) bool {
	in := Copy(slice)
	sort.Ints(in)
	var f func(int, int) bool
	f = func(l, r int) bool {
		if r-l <= 1 {
			return false
		}
		m := (l + r) / 2
		mv := in[m]
		if mv == i {
			return true
		}
		switch {
		case mv < i:
			return f(m, r)
		case mv > i:
			return f(l, m)
		}
		return false
	}
	if in[0] == i || in[len(in)-1] == i {
		return true
	}
	if i < in[0] || i > in[len(in)-1] {
		return false
	}
	return f(0, len(in)-1)
}

func Min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
func Max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
