package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sansna/cicily/v2/models"
	"github.com/sansna/cicily/v2/sankouyi"
	"github.com/sansna/cicily/v2/utils"
)

var (
	PLAYER_COUNT = 4
)

func matchBegin() {
	utils.ClearReplayMsg()

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
		utils.RepAddEmptyMsg()
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
					utils.RepAddLostCardMsg(&models.RepLostCardMsg{
						U:        i,
						PutCards: []int{},
					})
					break
				} else if (*p).Put(r) {
					utils.RepAddLostCardMsg(&models.RepLostCardMsg{
						U:         i,
						PutCards:  utils.SortCardByNameReturnId((*(*d).GetLastPut()).GetCards()),
						HandCards: utils.SortCardByNameReturnId((*p).GetCards()),
					})
					break
				}
				continue
			}

			if (*p).Win() {
				//fmt.Println("第", n, "局胜者", i)
				lastWinner = i
				scores = (*d).GetScores(lastWinner)
				utils.RepAddDisplayMsg(&models.RepDisplayMsg{
					U: -1,
					M: fmt.Sprintf("第%d局胜者: %d", n, i),
				})
				utils.RepAddScoreShowMsg(scores)
				break
			}
		}
	}
	//fmt.Println(m.GetWinner(), "wins")
	utils.RepAddDisplayMsg(&models.RepDisplayMsg{
		U: -1,
		M: fmt.Sprintf("比赛结束，第%d选手获胜，总分:%d", m.GetWinner(), m.GetScore()[m.GetWinner()]),
	})
	utils.RepAddScoreShowMsg(m.GetScore())
	//fmt.Println("scores:", m.GetScore())
}

func f(c *gin.Context) {
	matchBegin()
	repmsgs := utils.GetReplayMsg()
	type Ret struct {
		Msgs []string
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, &Ret{
		Msgs: repmsgs,
	})

}

func main() {
	r := gin.Default()
	r.GET("/getmatch", f)
	r.Run(":8888")
}
