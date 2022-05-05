package models

import (
	"sort"
)

const (
	RM_INVALID = iota
	RM_SINGLE
	RM_DUIZI // 一般对子Count=1，连对Count>=3
	RM_SHUNZI
	RM_SAN // 飞机情形Count>1
	RM_SI  // 飞机情形Count>1
	RM_WANGZHA
	RM_WU
	RM_LIU
	RM_SHUANGWANGZHA
	RM_QI
	RM_BA
)

const (
	RS_NONE = iota
	RS_SINGLE
	RS_DUIZI
)

type RecordType struct {
	Count int // 用于连对之类的计数，仅当RM=RM_SHUNZI,RM_DUIZI,RM_SAN,RM_SI可以非0
	RM    int
	RS    int   // 仅当RM = RM_SAN, RM_SI时可以非0, 1: 带单， 2：带对
	Mains []int // 仅当Count > 0 时非空
	Subs  []int // 仅当Count > 0 时非空
}

// compare by int
func (rt *RecordType) GetInt() int {
	return rt.RM<<2 + rt.RS + rt.Count<<8
}

func Hasdup(in []int) bool {
	sort.Ints(in)
	last := in[0]
	for i := 1; i < len(in); i++ {
		if last == in[i] {
			return true
		}
		last = in[i]
	}
	return false
}

func Trans(in []int) map[int]int {
	ret := make(map[int]int)
	for _, i := range in {
		i %= 54
		switch i {
		case 53, 52:
			ret[i] += 1
		default:
			v := i % 13
			switch v {
			case 0:
				// A
				v = 13
			case 1:
				// 2
				v = 15 // not continuous
			default:
				// 2-12: 3-K
			}
			ret[v] += 1
		}
	}
	return ret
}

func Convert(input []int) []int {
	out := []int{}
	for _, in := range input {
		in %= 54
		switch in {
		case 52, 53:
			out = append(out, in)
		default:
			v := in % 13
			switch v {
			case 0:
				v = 13
			case 1:
				v = 15
			default:
			}
			out = append(out, v)
		}
	}
	return out
}

func ConvertAndSort(input []int) []int {
	out := Convert(input)
	sort.Ints(out)
	return out
}

// retval: 0 not shunzi, > 0 is shunzi, and shunzilen = retval
func lenShunzi(in []int) int {
	out := 1
	sort.Ints(in)
	last := in[0]
	for i := 1; i < len(in); i++ {
		if last+1 == in[i] {
			out += 1
		} else {
			return 0
		}
		last = in[i]
	}
	return out
}

func GetType4(in []int, m map[int]int) RecordType {
	inv := RecordType{
		RM: RM_INVALID,
	}
	jk, vjk := m[52], m[53]
	if jk == 2 && vjk == 2 {
		return RecordType{
			RM:    RM_SHUANGWANGZHA,
			Mains: []int{53},
		}
	}
	switch len(m) {
	case 1:
		return RecordType{
			RM:    RM_SI,
			Mains: ConvertAndSort([]int{in[0]}),
		}
	case 2:
		var rt RecordType
		for k, v := range m {
			if v == 2 {
				return inv
			}
			if v == 1 {
				rt.Subs = []int{k}
			}
			if v == 3 {
				rt.Mains = []int{k}
			}
		}
		rt.Count = 1
		rt.RM = RM_SAN
		rt.RS = RS_SINGLE
		return rt
	default:
	}
	return inv
}

func checkLianDui(m map[int]int) bool {
	is := []int{}
	for k, v := range m {
		if v != 2 {
			return false
		}
		is = append(is, k)
	}
	return lenShunzi(is) >= 3
}

// return longest possible sorted shunzi
// n = 3,4
func GetLongestShunzi(m map[int]int, n int) []int {
	is := []int{}
	for k, v := range m {
		if v >= n {
			is = append(is, k)
		}
	}
	if len(is) == 0 {
		return []int{}
	}
	sort.Ints(is)
	l, r := 0, 0
	last := is[0]
	maxl, maxr := len(is)-1, len(is)
	for r < len(is)-1 {
		r++
		rv := is[r]
		if last+1 == rv {
			last = rv
			continue
		} else {
			last = rv
			if r-l >= maxr-maxl {
				maxl, maxr = l, r
			}
			l = r
		}
	}
	if r-l+1 >= maxr-maxl {
		maxl, maxr = l, r+1
	}
	return is[maxl:maxr]
}

// retval nil if not satisfy condition
func getRemainSubs(m map[int]int, shunzi []int, size int, isSingle bool) []int {
	nm := map[int]int{}
	for _, v := range shunzi {
		nm[v] = size
	}
	out := []int{}
	step := 1
	if !isSingle {
		step = 2
	}
	for k, v := range m {
		v -= nm[k]
		if step == 2 && v%2 == 1 {
			return []int{}
		}
		for i := 0; i < v; i += step {
			out = append(out, k)
		}
	}
	sort.Ints(out)
	return out
}

// 6. 其他牌型先看是否同数字，再看是否顺子，最后找到出现次数最多的数字x及次数n
// 如果最大次数为2，看是否是连对
// 次数n最大为4，再看3， 针对每个n，循环尝试增加一个连续数字y，看剩下的牌是否满足对子或是单牌数量需求，满足即返回牌型，如果循环结束也没找到就返回invalid
func GetTypeDefault(in []int, m map[int]int) RecordType {
	inv := RecordType{
		RM: RM_INVALID,
	}
	// 是否同数字： >=5
	if len(m) == 1 {
		n := m[ConvertAndSort(in)[0]]
		switch n {
		case 5:
			return RecordType{
				RM:    RM_WU,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		case 6:
			return RecordType{
				RM:    RM_LIU,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		case 7:
			return RecordType{
				RM:    RM_QI,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		case 8:
			return RecordType{
				RM:    RM_BA,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		}
	}
	if len(m) == len(in) {
		t := []int{}
		for k := range m {
			t = append(t, k)
		}
		if lenShunzi(t) > 0 {
			return RecordType{
				RM:    RM_SHUNZI,
				Count: len(in),
				Mains: ConvertAndSort(in),
			}
		}
	}
	if checkLianDui(m) {
		rt := RecordType{
			RM:    RM_DUIZI,
			Count: len(in) / 2,
			Mains: GetLongestShunzi(m, 2),
		}
		return rt
	}
	// size 4 test
	ln := len(in)
	shunzi := GetLongestShunzi(m, 4)
	if ln%4 == 0 {
		// 连4检查
		if len(shunzi) == ln/4 {
			return RecordType{
				RM:    RM_SI,
				Mains: shunzi,
				Count: ln / 4,
			}
		}
	}
	if ln%8 == 0 {
		// 4带对*2检查
		for i := 0; i <= len(shunzi)-ln/6; i++ {
			subs := getRemainSubs(m, shunzi[len(shunzi)-ln/6-i:len(shunzi)-i], 4, false)
			if len(subs) > 0 {
				return RecordType{
					RM:    RM_SI,
					RS:    RS_DUIZI,
					Count: ln / 6,
					Mains: shunzi[len(shunzi)-ln/6-i : len(shunzi)-i],
					Subs:  subs,
				}
			}
		}
	}
	if ln%6 == 0 {
		// 4带1*2检查
		n := ln / 5
		if n <= len(shunzi) {
			return RecordType{
				RM:    RM_SI,
				RS:    RM_SINGLE,
				Count: n,
				Mains: shunzi[len(shunzi)-n:],
				Subs:  getRemainSubs(m, shunzi, 4, true),
			}
		}
	}

	// size 3 test
	shunzi = GetLongestShunzi(m, 3)
	if ln%3 == 0 {
		// 连3检查
		if ln/3 == len(shunzi) {
			return RecordType{
				RM:    RM_SAN,
				Mains: shunzi,
				Count: ln / 3,
			}
		}
	}
	if ln%5 == 0 {
		// 3带对检查
		for i := 0; i <= len(shunzi)-ln/5; i++ {
			subs := getRemainSubs(m, shunzi[len(shunzi)-ln/5-i:len(shunzi)-i], 3, false)
			if len(subs) > 0 {
				return RecordType{
					RM:    RM_SAN,
					RS:    RS_DUIZI,
					Count: ln / 5,
					Mains: shunzi[len(shunzi)-ln/5-i : len(shunzi)-i],
					Subs:  subs,
				}
			}
		}
	}
	if ln%4 == 0 {
		// 3带单检查
		if ln/4 <= len(shunzi) {
			return RecordType{
				RM:    RM_SAN,
				RS:    RS_SINGLE,
				Count: ln / 4,
				Mains: shunzi[len(shunzi)-ln/4:],
				Subs:  getRemainSubs(m, shunzi, 3, true),
			}
		}
	}

	return inv
}

func GetType(in []int) RecordType {
	inv := RecordType{
		RM: RM_INVALID,
	}
	// 1. 检查牌数>0
	if len(in) == 0 {
		return inv
	}
	// 1.5 检查是否有重复牌
	if Hasdup(in) {
		return inv
	}
	m := Trans(in)
	switch len(in) {
	case 1:
		// 2. 牌数1直接返回单牌
		return RecordType{
			RM:    RM_SINGLE,
			Mains: ConvertAndSort([]int{in[0]}),
		}
	case 2:
		// 3. 牌数2检查是否对子
		if len(m) == 1 {
			return RecordType{
				RM:    RM_DUIZI,
				Count: 1,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		}
		jk, vjk := m[52], m[53]
		if jk == 1 && vjk == 1 {
			return RecordType{
				RM:    RM_WANGZHA,
				Mains: []int{53},
			}
		}
		return inv
	case 3:
		// 4. 牌数3检查是否同数字
		if len(m) == 1 {
			return RecordType{
				RM:    RM_SAN,
				Count: 1,
				Mains: ConvertAndSort([]int{in[0]}),
			}
		}
		return inv
	case 4:
		// 5. 牌数4检查是否3带1，是否炸弹，是否双王炸
		return GetType4(in, m)
	default:
		// 6. 其他牌型先看是否同数字，再看是否顺子，最后找到出现次数最多的数字x及次数n
		// 次数n最大为4，再看3， 针对每个n，循环尝试增加一个连续数字y，看剩下的牌是否满足对子或是单牌数量需求，满足即返回牌型，如果循环结束也没找到就返回invalid
		// 如果最大次数为2，看是否是连对
		return GetTypeDefault(in, m)
	}
}
