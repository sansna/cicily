package sankouyi

import "github.com/sansna/golang.go/card/models"

type SKYMatch struct {
	parts     []*SKYSingle
	winner    int
	idx       int
	finished  bool
	highscore int
	scores    map[int]int
}

// 假定最高积分不会相等
func (m *SKYMatch) GetWinner() int {
	if m.finished {
		return m.winner
	}
	// 比赛尚未完成
	if m.idx != len(m.parts) {
		return -1
	}
	// 比赛已完成,计算胜者
	scores := m.GetScore()
	max := 0
	winner := -1
	for k, v := range scores {
		if v > max {
			winner = k
			max = v
		}
	}
	m.finished = true
	m.winner = winner
	m.highscore = max
	return m.winner
}
func (m *SKYMatch) GetScore() map[int]int {
	return m.scores
}

func (m *SKYMatch) GetCurrentSingle() *models.Single {
	v := m.parts[m.idx]
	r := models.Single(v)
	return &r
}

func (m *SKYMatch) GetNewSingle(winner int, add_score map[int]int) *models.Single {
	// 重复调用不增加
	if m.idx == len(m.parts) {
		return nil
	}
	// 第一局特殊处理
	if winner == -1 {
		v := m.parts[0]
		r := models.Single(v)
		return &r
	}
	m.parts[m.idx].SetWinner(winner)
	m.idx++
	// 积分
	if m.scores == nil {
		m.scores = make(map[int]int)
	}
	for playerId, score := range add_score {
		m.scores[playerId] += score
	}
	if len(m.parts) > m.idx {
		v := m.parts[m.idx]
		r := models.Single(v)
		return &r
	}
	// 比赛结束
	return nil
}

func (m *SKYMatch) AddSingles(in []*models.Single) bool {
	// TODO: Optimize for Single specifications
	for i := 0; i < len(in); i++ {
		m.parts = append(m.parts, &SKYSingle{})
	}
	return true
}
