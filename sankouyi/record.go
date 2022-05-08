package sankouyi

import (
	"github.com/sansna/cicily/v2/models"
	"github.com/sansna/cicily/v2/utils"
)

type SKYRecord struct {
	PlayerId  int
	Cards     []int
	Type      int
	innerType *models.RecordType
}

func (r *SKYRecord) GetPlayerId() int {
	return r.PlayerId
}

func (r *SKYRecord) GetCardsLen() int {
	return len(r.Cards)
}
func (r *SKYRecord) GetCardsType() int {
	if r.Type != 0 {
		return r.Type
	}
	t := models.GetType(r.Cards)
	r.innerType = &t
	r.Type = t.GetInt()
	return r.Type
}
func (r *SKYRecord) GetCards() []int {
	return r.Cards
}

func (r *SKYRecord) Compare(in *models.Record) int {
	rr := &SKYRecord{
		PlayerId: (*in).GetPlayerId(),
		Cards:    (*in).GetCards(),
	}
	rb, rrb := r.IsBomb(), rr.IsBomb()
	switch {
	case !rb && !rrb:
		// 两边都不是炸弹
		if r.Type != rr.Type {
			// 类型不同无法比较
			return 0
		}
		if r.innerType.Mains[0] > rr.innerType.Mains[0] {
			// 牌主体部分比对方大
			return 1
		} else if r.innerType.Mains[0] == rr.innerType.Mains[0] {
			// 牌主体部分相同,看次要部分
			rs, rrs := r.innerType.Subs, rr.innerType.Subs
			if rs == nil || len(rs) == 0 {
				// 没有次要部分
				return -1
			}
			sub_all_equal := true
			for i := 0; i < len(rs); i++ {
				if rs[i] < rrs[i] {
					// 次要部分有一个比对方小
					return -1
				} else if rs[i] > rrs[i] {
					// 次要部分有比对方大的做标记
					sub_all_equal = false
				}
			}
			// 走到这里，表明所有次要项至少不比对面小
			if !sub_all_equal {
				// 存在一个次要项比对面大
				return 1
			} else {
				// 次要项完全相等
				return -1
			}
		}
		// 牌主体部分比对方小
		return -1
	case !rb && rrb:
		// 对方是炸弹
		return -1
	case rb && !rrb:
		// 我方是炸弹
		return 1
	case rb && rrb:
		// 两边都是炸弹
		if r.innerType.RM == rr.innerType.RM {
			// 两边炸弹等级相同，对比具体数字
			if r.innerType.Mains[0] > rr.innerType.Mains[0] {
				return 1
			} else {
				return -1
			}
		} else {
			if r.innerType.RM > rr.innerType.RM {
				// 等级大于对方
				return 1
			} else {
				// 等级比对方小
				return -1
			}
		}
	}
	// 走不到这里
	return 0
}

var (
	BOMB_TYPES = []*models.RecordType{
		{RM: models.RM_SI},
		{RM: models.RM_WU},
		{RM: models.RM_LIU},
		{RM: models.RM_QI},
		{RM: models.RM_BA},
		{RM: models.RM_WANGZHA},
		{RM: models.RM_SHUANGWANGZHA},
	}
	BOMB_INTS = []int{}
)

func init() {
	for _, t := range BOMB_TYPES {
		BOMB_INTS = append(BOMB_INTS, t.GetInt())
	}
}

func (r *SKYRecord) IsBomb() bool {
	t := r.GetCardsType()
	return utils.InSlice(t, BOMB_INTS)
}
