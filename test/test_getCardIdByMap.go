package main

import (
	"fmt"

	"github.com/sansna/cicily/v2/sankouyi"
	"github.com/sansna/cicily/v2/utils"
)

func main() {
	fmt.Println("vim-go")
	in := []int{2, 2}
	hc := []int{0, 2, 10, 14, 24, 32, 34, 41, 55, 57, 61, 63, 64, 67, 71, 79, 80, 81, 89, 97, 99, 101, 103, 104}
	fmt.Println(utils.GetAllCardsName(sankouyi.GetCardIdByValue(in, hc), 54), utils.GetAllCardsName(hc, 54))
}
