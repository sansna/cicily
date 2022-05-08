package sankouyi

import (
	"github.com/sansna/cicily/v2/models"
	"github.com/sansna/cicily/v2/utils"
)

type SKYPlayer struct {
	cards    []int
	cn       int
	playerId int
	desktop  *SKYDesktop
}

// couple with single/desktop 检查是否可打该牌,可打时进行牌处理
func (p *SKYPlayer) Put(in *models.Record) bool {
	r := &SKYRecord{
		//PlayerId: (*in).GetPlayerId(),
		Cards: (*in).GetCards(),
	}
	r.PlayerId = p.playerId
	if p.desktop.AddPutHistory(in) {
		// 对应牌从手牌中清掉
		cardOrderMap := map[int]int{}
		for idx, c := range p.cards {
			cardOrderMap[c] = idx
		}
		for _, c := range r.Cards {
			//pos := sort.SearchInts(p.cards, c)
			pos := cardOrderMap[c]
			p.cards[pos] = -1
			//p.Cards = append(p.Cards[:pos], p.Cards[pos+1:]...)
		}
		left := 0
		for i := 0; i < p.cn; i++ {
			if p.cards[i] >= 0 {
				if left != i {
					p.cards[left] = p.cards[i]
				}
				left++
			}
		}
		p.cn -= len(r.Cards)
		p.cards = p.cards[:p.cn]
		return true
	}
	return false
}

// 自查逻辑
func (p *SKYPlayer) CheckPut(in *models.Record) bool {
	r := &SKYRecord{
		PlayerId: (*in).GetPlayerId(),
		Cards:    (*in).GetCards(),
	}
	// 1. check if dup
	if models.Hasdup(r.Cards) {
		return false
	}
	// 2. check if all exist in hand
	if len(r.Cards) > p.cn {
		return false
	}
	for _, c := range r.Cards {
		if !utils.InSlice(c, p.cards) {
			return false
		}
	}
	// 客户端初步检查通过
	return true
}

// XXX: See sankouyi/ai.go
//func (p *SKYPlayer) AiGenRecord(in *models.Record) *models.Record {
//	cards := (*p).GetCards()
//	cn := len(cards)
//	// 没牌出当然不出牌
//	if cn == 0 {
//		return nil
//	}
//	// 随机拿个牌
//	v := &SKYRecord{
//		PlayerId: p.playerId,
//		Cards:    []int{cards[0]},
//	}
//	r := models.Record(v)
//	return &r
//}

// 直接进行pass操作
func (p *SKYPlayer) Pass() {
	r := &SKYRecord{
		PlayerId: p.playerId,
		Cards:    nil,
	}
	rr := models.Record(r)
	p.desktop.AddPutHistory(&rr)
}

func (p *SKYPlayer) Win() bool {
	return p.cn == 0
}

func (p *SKYPlayer) GetRecord() *models.Record {
	v := &SKYRecord{
		PlayerId: p.playerId,
		Cards:    p.cards,
	}
	r := models.Record(v)
	return &r
}

func (p *SKYPlayer) GetCards() []int {
	return p.cards
}
