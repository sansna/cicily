package utils

import (
	"fmt"
	"sort"
)

// CARD_COUNT defaults to 54
func GetCardName(i int, CARD_COUNT int) string {
	n := i % CARD_COUNT
	switch n {
	case 53:
		return "大王"
	case 52:
		return "小王"
	default:
	}
	cn := ""
	switch n / 13 {
	case 0:
		cn += "方块"
	case 1:
		cn += "黑桃"
	case 2:
		cn += "红桃"
	case 3:
		cn += "梅花"
	default:
	}
	switch n % 13 {
	case 0:
		cn += "A"
	case 9:
		cn += "10"
	case 10:
		cn += "J"
	case 11:
		cn += "Q"
	case 12:
		cn += "K"
	default:
		cn += string(rune('1' + n%13))
	}
	return cn
}

// CARD_COUNT defaults to 54
func GetAllCardsName(in []int, CARD_COUNT int) []string {
	out := []string{}
	in = Copy(in)
	sc := &SKYSortCards{
		Cards: &in,
		L:     len(in),
	}
	sort.Sort(sc)
	for _, i := range in {
		out = append(out, GetCardName(i, CARD_COUNT))
	}
	return out
}

func P(in []string) {
	for _, i := range in {
		fmt.Println(i)
	}
}

func SortCardByNameReturnId(in []int) []int {
	cards := Copy(in)
	sc := &SKYSortCards{
		Cards: &cards,
		L:     len(cards),
	}
	sort.Sort(sc)
	return cards
}
