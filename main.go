package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sansna/golang.go/card/models"
	"github.com/sansna/golang.go/card/sankouyi"
)

var (
	PLAYER_COUNT = 4
)

func matchBegin() {
	var m models.Match
	m = &sankouyi.SKYMatch{}
	singles := make([]*models.Single, 5)
	m.AddSingles(singles)

	lastWinner := -1
	var s *models.Single
	var scores map[int]int
	var n int
	for s = m.GetNewSingle(lastWinner, nil); s != nil; s = m.GetNewSingle(lastWinner, scores) {
		n++
		// 单局比赛开始
		d := (*s).GetDesktop()
		players := (*d).DistCard() // 发牌
		PLAYER_COUNT = len(players)
		st := (*d).GetFirstPut()
		for i := st; ; i++ {
			i %= PLAYER_COUNT
			// playerid i to put or pass
			p := players[i]
			var r *models.Record
			for {
				r = (*p).AiGenRecord((*d).GetLastPut())
				if r == nil || (*r).GetCardsLen() == 0 {
					(*p).Pass()
					break
				} else if (*p).Put(r) {
					break
				}
				continue
			}

			if (*p).Win() {
				fmt.Println("第", n, "局胜者", i)
				lastWinner = i
				scores = (*d).GetScores(lastWinner)
				break
			}
		}
	}
	fmt.Println(m.GetWinner(), "wins")
	fmt.Println("scores:", m.GetScore())
}

func main() {
	r := gin.Default()

}
