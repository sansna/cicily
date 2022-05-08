package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sansna/cicily/v2/models"
)

var (
	replayMsg []string
)

func init() {
	replayMsg = make([]string, 0, 20)
}

func RepAddEmptyMsg() {
	m := &models.RepEmptyMsg{}
	msg := &models.RepMsg{
		Type: 1,
		Msg:  m,
	}
	GenRep(msg)
}
func RepAddGetCardMsg(m *models.RepGetCardMsg) {
	msg := &models.RepMsg{
		Type: 2,
		Msg:  m,
	}
	GenRep(msg)
}
func RepAddLostCardMsg(m *models.RepLostCardMsg) {
	msg := &models.RepMsg{
		Type: 3,
		Msg:  m,
	}
	GenRep(msg)
}
func RepAddDisplayMsg(m *models.RepDisplayMsg) {
	msg := &models.RepMsg{
		Type: 4,
		Msg:  m,
	}
	GenRep(msg)
}
func RepAddPrettyCardMsg(m *models.RepPrettyCardMsg) {
	msg := &models.RepMsg{
		Type: 5,
		Msg:  m,
	}
	GenRep(msg)
}
func GenRep(m interface{}) {
	byt, _ := json.Marshal(m)
	replayMsg = append(replayMsg, string(byt))
}
func GetReplayMsg() []string {
	return replayMsg
}
func ClearReplayMsg() {
	replayMsg = replayMsg[:0]
}

func RepAddScoreShowMsg(scores map[int]int) {
	var s strings.Builder
	for i := 0; i < 4; i++ {
		s.WriteString(fmt.Sprintf("第%d选手得分：%d\n", i, scores[i]))
	}
	m := &models.RepDisplayMsg{
		U: -1,
		M: s.String(),
	}
	msg := &models.RepMsg{
		Type: 4,
		Msg:  m,
	}
	GenRep(msg)
}
