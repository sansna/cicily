package models

type RepMsg struct {
	Type int
	Msg  interface{}
}
type RepEmptyMsg struct {
}
type RepGetCardMsg struct {
	U     int
	Cards []int
}
type RepLostCardMsg struct {
	U         int
	HandCards []int
	PutCards  []int
}
type RepDisplayMsg struct {
	U int
	M string
}
type RepPrettyCardMsg struct {
	U     int
	Cards []int
}
