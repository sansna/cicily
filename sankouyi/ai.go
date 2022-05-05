package sankouyi

import (
	"log"
	"sort"

	"github.com/sansna/golang.go/card/models"
	"github.com/sansna/golang.go/card/utils"
)

// Analyzer Supported Types
// supposed to change over days
const (
	AT_BOMB AnalyzerType = iota
	AT_SHUNZI
	AT_FEIJI_SI
	AT_FEIJI_SAN
	AT_LIANDUI
	AT_SI
	AT_SAN
	AT_DUIZI
	AT_SINGLE
)

type AnalyzerType int

var (
	DEFAULT_RECIPE = []AnalyzerType{
		AT_BOMB,
		AT_SHUNZI,
		AT_FEIJI_SI,
		AT_FEIJI_SAN,
		AT_LIANDUI,
		AT_SI,
		AT_SAN,
		AT_DUIZI,
		AT_SINGLE,
	}
)

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
//		Cards:	  []int{cards[0]},
//	}
//	r := models.Record(v)
//	return &r
//}

// 给定手牌，以及想出的牌大小，返回手牌的id
// NOTE: in: '3': 2, 'A': 13, '2': 15, 'Joker': 53
// TODO: optimize
func GetCardIdByValue(in []int, handCard []int) []int {
	if len(in) == len(handCard) {
		return handCard
	}
	ch := models.Convert(handCard)
	ret := []int{}
	for _, i := range in {
		for j := 0; j < len(handCard); j++ {
			if ch[j] == i {
				ret = append(ret, handCard[j])
				ch[j] = -1
				break
			}
		}
	}
	return ret
}

func GetCardIdByMap(m map[int]int, handCard []int) []int {
	slice := utils.MapToSlice(m)
	return GetCardIdByValue(slice, handCard)
}

// retval: all possible bombs
func GetAllBombs(m map[int]int) []map[int]int {
	ret := make([]map[int]int, 0)
	for k, v := range m {
		if v > 3 {
			for i := 4; i <= v; i++ {
				ret = append(ret, map[int]int{
					k: i,
				})
			}
		}
	}
	jk, vjk := m[52], m[53]
	mjk := utils.Min(jk, vjk)
	switch mjk {
	case 2:
		ret = append(ret, map[int]int{
			52: 2,
			53: 2,
		})
		fallthrough
	case 1:
		ret = append(ret, map[int]int{
			52: 1,
			53: 1,
		})
	default:
	}
	return ret
}

// given n of type of shunzi, and shunzi length
// check if a truly shunzi
// n >=1
func checkLenShunzi(l int, n int) bool {
	switch n {
	case 1:
		if l >= 5 {
			return true
		}
	case 2:
		if l >= 3 {
			return true
		}
	default:
		if l >= 2 {
			return true
		}
	}
	return false
}

// all possible shunzi of *n* repitation, retval: [l, r] of shunzi
// n >= 1
func GetAllShunzi(m map[int]int, n int) []map[int]int {
	// [l-r] trans to map
	f := func(ht []int, n int) map[int]int {
		m := make(map[int]int, 0)
		for i := ht[0]; i <= ht[1]; i++ {
			m[i] = n
		}
		return m
	}
	is := []int{}
	for k, v := range m {
		if v >= n {
			is = append(is, k)
		}
	}
	ret := make([]map[int]int, 0)
	if len(is) == 0 {
		return ret
	}
	sort.Ints(is)
	l, r := 0, 0
	last := is[0]
	for r < len(is)-1 {
		r++
		rv := is[r]
		if last+1 == rv {
			last = rv
			continue
		} else {
			last = rv
			if checkLenShunzi(r-l, n) {
				//ret = append(ret, []int{l, r - 1})
				ret = append(ret, f([]int{is[l], is[r-1]}, n))
			}
			l = r
		}
	}
	if checkLenShunzi(r-l+1, n) {
		//ret = append(ret, []int{l, r})
		ret = append(ret, f([]int{is[l], is[r]}, n))
	}
	return ret
}

// GetAllSan, GetAllSi, GetAllDuizi, GetAllSingle generalization
// n := 3, 4, 2, 1
func GetAllDup(m map[int]int, n int, exact bool) []map[int]int {
	ret := make([]map[int]int, 0)
	if exact {
		// exact mode
		for k, v := range m {
			if v == n {
				ret = append(ret, map[int]int{
					k: n,
				})
			}
		}
	} else {
		// tolerance mode
		for k, v := range m {
			if v >= n {
				ret = append(ret, map[int]int{
					k: n,
				})
			}
		}
	}
	return ret
}

func SpecificAnalyzer(m map[int]int, recipe []AnalyzerType) []map[int]int {
	mm := utils.MCopyIntInt(m)
	ret := make([]map[int]int, 0)
	var rs []map[int]int
	for _, at := range recipe {
		switch at {
		case AT_BOMB:
			rs = GetAllBombs(mm)
		case AT_SHUNZI:
			rs = GetAllShunzi(mm, 1)
		case AT_SI:
			rs = GetAllDup(mm, 4, false)
		case AT_SAN:
			rs = GetAllDup(mm, 3, false)
		case AT_DUIZI:
			rs = GetAllDup(mm, 2, false)
		case AT_SINGLE:
			rs = GetAllDup(mm, 1, false)
		case AT_FEIJI_SI:
			rs = GetAllShunzi(mm, 4)
		case AT_FEIJI_SAN:
			rs = GetAllShunzi(mm, 3)
		case AT_LIANDUI:
			rs = GetAllShunzi(mm, 2)
		default:
			log.Println("Unknown AnalyzerType at:", at)
		}
		ret = append(ret, rs...)
		utils.MSubSliceInplace(mm, rs)
		rs = rs[:0]
	}
	return ret
}

// hand score: the less the better.
func CardsScore(m map[int]int, recipe []AnalyzerType) int {
	mm := utils.MCopyIntInt(m)
	//ret := make([]map[int]int, 0)
	score := 0
	// base score for length
	for _, v := range mm {
		score += v
	}
	var rs []map[int]int
	for _, at := range recipe {
		switch at {
		case AT_BOMB:
			rs = GetAllBombs(mm)
			// simple calc
			score -= 8 * len(rs)
		case AT_SHUNZI:
			rs = GetAllShunzi(mm, 1)
			score += 5 * len(rs)
		case AT_SI:
			rs = GetAllDup(mm, 4, false)
			score += 5 * len(rs)
		case AT_SAN:
			rs = GetAllDup(mm, 3, false)
			score += 5 * len(rs)
		case AT_DUIZI:
			rs = GetAllDup(mm, 2, false)
			score += 6 * len(rs)
		case AT_SINGLE:
			rs = GetAllDup(mm, 1, false)
			score += 7 * len(rs)
		case AT_FEIJI_SI:
			rs = GetAllShunzi(mm, 4)
			score += 3 * len(rs)
		case AT_FEIJI_SAN:
			rs = GetAllShunzi(mm, 3)
			score += 4 * len(rs)
		case AT_LIANDUI:
			rs = GetAllShunzi(mm, 2)
			score += 5 * len(rs)
		default:
			log.Println("Unknown AnalyzerType at:", at)
		}
		utils.MSubSliceInplace(mm, rs)
		rs = rs[:0]
	}
	// 如果还有未解析牌继续加分
	for _, v := range mm {
		score += v * 1
	}
	return score
}

func GenRecordFromMap(m map[int]int, handCard []int, playerId int) *models.Record {
	cards := GetCardIdByMap(m, handCard)
	v := &SKYRecord{
		Cards:    cards,
		PlayerId: playerId,
	}
	r := models.Record(v)
	return &r
}

// Retval: true: catch it
// false: pass
func DecideCatch(p *SKYPlayer, lastPut *SKYRecord) bool {
	// 抢队友出牌
	if p.playerId != p.desktop.startPlayer {
		// 所有手牌符合条件能直接扔了就赢就直接赢了
		r := models.Record(lastPut)
		if (*(*p).GetRecord()).Compare(&r) == 1 {
			return true
		}
		// TODO 如果跟牌后手牌积分下降百分比特别高的，也可以抢着出
	}
	// 本人就是地主，抓
	// 本人不是地主，上个出牌的是地主，抓
	// 本人不是地主，且上个出牌的也不是地主，不抓
	if p.playerId != p.desktop.startPlayer && lastPut.PlayerId != p.desktop.startPlayer {
		// TODO: 当地主手牌为1或者2时，队友认为牌面不够大，可以进行增强
		return false
	}
	return true
}

func (p *SKYPlayer) AiGenRecord(in *models.Record) *models.Record {
	var lpIdx int
	var lp *SKYPlayer
	if in != nil {
		lpIdx = (*in).GetPlayerId()
		lp = p.desktop.players[lpIdx]
	}
	ret := SpecificAnalyzer(models.Trans(p.GetCards()), DEFAULT_RECIPE)
	// 如果我是第一个出牌不用管上次
	// 如果上次牌就是我打的，那就不用管上次的牌，随便出个非炸弹长度最大的
	// 决定出牌时的长度
	// TODO: 非地主查看自己后面的队友的手牌数，如果够小，需要放小牌让他跑
	if in == nil || lp.playerId == p.playerId {
		//ret[len(ret)-1]
		//return GenRecordFromMap(ret[len(ret)-1], p.GetCards(), p.playerId)
		resp := GenRecordFromMap(ret[len(ret)-1], p.GetCards(), p.playerId)
		maxl := 0
		for i := 0; i < len(ret); i++ {
			tmp := ret[len(ret)-1-i]
			tmpv := &SKYRecord{
				Cards:    GetCardIdByMap(tmp, p.GetCards()),
				PlayerId: p.playerId,
			}
			if tmpv.IsBomb() {
				break
			} else if utils.MapToSliceLen(tmp) > maxl {
				maxl = utils.MapToSliceLen(tmp)
				t := models.Record(tmpv)
				resp = &t
			}
		}
		return resp
	}
	lv := &SKYRecord{
		Cards:    (*in).GetCards(),
		PlayerId: lpIdx,
	}
	lr := models.Record(lv)
	// 决策是否跟牌/抓
	if !DecideCatch(p, lv) {
		return nil
	}
	// 需要对比其他人出牌
	for i := 0; i < len(ret); i++ {
		rr := ret[len(ret)-1-i]
		v := &SKYRecord{
			Cards:    GetCardIdByMap(rr, p.GetCards()),
			PlayerId: p.playerId,
		}
		if v.Compare(&lr) == 1 {
			r := models.Record(v)
			return &r
		}
	}
	return nil
}
