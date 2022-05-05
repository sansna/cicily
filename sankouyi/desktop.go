package sankouyi

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/sansna/golang.go/card/models"
	"github.com/sansna/golang.go/card/utils"
)

const (
	CARD_COUNT   = 54
	CARD_STACK   = 2
	PLAYER_COUNT = 4
	REMAIN_CARD  = 4
)

type SKYDesktop struct {
	putHistory  []*SKYRecord
	lastPutIdx  int
	players     []*SKYPlayer
	startPlayer int
}

func (d *SKYDesktop) GetLastPut() *models.Record {
	if d.lastPutIdx == 0 {
		//return new(models.Record)
		return nil
	}
	rec := models.Record(d.putHistory[d.lastPutIdx-1])
	ret := &rec
	return ret
}

func (d *SKYDesktop) DistCard() []*models.Player {
	players := make([][]int, PLAYER_COUNT)
	for i := range players {
		players[i] = []int{}
	}
	total_cards := CARD_STACK * CARD_COUNT
	// 发牌
	rand.Seed(time.Now().UnixNano())
	permcards := rand.Perm(total_cards)
	startplayer := rand.Intn(PLAYER_COUNT)
	for i := 0; i < total_cards-REMAIN_CARD; i++ {
		pn := (i + startplayer) % PLAYER_COUNT
		players[pn] = append(players[pn], permcards[i])
	}
	// TODO 决策谁当地主
	d.startPlayer = rand.Intn(4)
	fmt.Println("本局地主", d.startPlayer)
	players[d.startPlayer] = append(players[d.startPlayer], permcards[total_cards-REMAIN_CARD:]...)
	// 理牌
	for i := 0; i < PLAYER_COUNT; i++ {
		sc := &utils.SKYSortCards{
			Cards: &players[i],
			L:     len(players[i]),
		}
		sort.Sort(sc)
		fmt.Println("第", i, "选手手牌", utils.GetAllCardsName(players[i], 54))
	}
	retPlayers := make([]*SKYPlayer, 0, PLAYER_COUNT)
	for i := 0; i < PLAYER_COUNT; i++ {
		retPlayers = append(retPlayers, &SKYPlayer{
			cards:    players[i],
			cn:       len(players[i]),
			playerId: i,
			desktop:  d,
		})
	}
	d.players = retPlayers
	ret := make([]*models.Player, 0)
	for _, p := range retPlayers {
		pp := models.Player(p)
		ret = append(ret, &pp)
	}
	return ret
}

func (d *SKYDesktop) GetFirstPut() int {
	return d.startPlayer
}

func (d *SKYDesktop) GetScores(winner int) map[int]int {
	m := make(map[int]int)
	startPoint := 15
	otherPoint := 5
	switch winner == d.startPlayer {
	case true:
		for i := 0; i < PLAYER_COUNT; i++ {
			m[i] = -1 * otherPoint
		}
		m[winner] = startPoint
	case false:
		for i := 0; i < PLAYER_COUNT; i++ {
			m[i] = 1 * otherPoint
		}
		m[winner] = -1 * startPoint
	}
	return m
}

// 用户打牌检查，通过时增加牌桌历史记录
func (d *SKYDesktop) AddPutHistory(r *models.Record) bool {
	in := &SKYRecord{
		PlayerId: (*r).GetPlayerId(),
		Cards:    (*r).GetCards(),
	}
	// 检查是否轮到该用户出牌
	if len(d.putHistory) > 0 {
		lastPlayerid := d.putHistory[len(d.putHistory)-1].PlayerId
		if (lastPlayerid+1)%PLAYER_COUNT != in.PlayerId {
			return false
		}
	}
	if d.putHistory == nil {
		d.putHistory = make([]*SKYRecord, 0, 10)
	}
	// 该用户选择pass该轮
	if in.Cards == nil || len(in.Cards) == 0 {
		d.putHistory = append(d.putHistory, in)
		return true
	}
	// 检查是否牌重复以及用户手牌中是否有该牌
	if !d.players[in.PlayerId].CheckPut(r) {
		return false
	}
	in.Type = in.GetCardsType()
	if in.Type == 0 {
		// 牌型检查不通过，打回重打
		return false
	}
	// 本次出牌与上家出牌对比大小，或者我们第一个出牌，或者上家就是我们自己，就可以出牌
	lastPustHistory := d.GetLastPut()
	if lastPustHistory == nil || (*lastPustHistory).GetPlayerId() == (*r).GetPlayerId() || in.Compare(lastPustHistory) == 1 {
		// 可以出牌， 加入记录
		d.putHistory = append(d.putHistory, in)
		d.lastPutIdx = len(d.putHistory)
		// XXX
		fmt.Println((*r).GetPlayerId(), utils.GetAllCardsName((*r).GetCards(), 54))
		return true
	}
	// 其他情况不许出牌

	return false
}
