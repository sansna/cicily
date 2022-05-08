package main

import (
	"fmt"

	"github.com/sansna/cicily/v2/utils"
)

func main() {
	for i := 0; i < 54; i++ {
		n := utils.GetCardName(i, 54)
		fmt.Println(n, i)
	}
}
