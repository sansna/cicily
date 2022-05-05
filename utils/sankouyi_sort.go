package utils

type SKYSortCards struct {
	Cards *[]int
	L     int
}

func (c *SKYSortCards) Less(i, j int) bool {
	iv, jv := (*c.Cards)[i]%54, (*c.Cards)[j]%54
	if iv == 53 {
		return false
	}
	if jv == 53 {
		return true
	}
	if iv == 52 {
		return false
	}
	if jv == 52 {
		return true
	}
	in, jn := iv%13, jv%13
	if in < 2 {
		in += 13
	}
	if jn < 2 {
		jn += 13
	}
	switch {
	case in < jn:
		return true
	case in > jn:
		return false
	default:
		ik, jk := iv/13, jv/13
		if ik < jk {
			return true
		}
		return false
	}
}

func (c *SKYSortCards) Swap(i, j int) {
	(*c.Cards)[i], (*c.Cards)[j] = (*c.Cards)[j], (*c.Cards)[i]
}

func (c *SKYSortCards) Len() int {
	return c.L
}
