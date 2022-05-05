package main

import (
	"fmt"

	"github.com/sansna/golang.go/card/utils"
)

func main() {
	for i := 0; i < 54; i++ {
		n := utils.GetCardName(i, 54)
		fmt.Println(n, i)
	}
}
