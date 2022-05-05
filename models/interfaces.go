package models

type Record interface {
	GetPlayerId() int
	GetCardsLen() int
	GetCardsType() int
	GetCards() []int
	Compare(*Record) int
}

type Desktop interface {
	GetLastPut() *Record
	DistCard() []*Player
	GetFirstPut() int
	AddPutHistory(*Record) bool
	GetScores(int) map[int]int
}
type Player interface {
	AiGenRecord(*Record) *Record
	Put(*Record) bool
	CheckPut(*Record) bool
	Pass()
	Win() bool
	GetCards() []int
	GetRecord() *Record
}
type Single interface {
	GetWinner() int
	GetDesktop() *Desktop
	SetWinner(int)
}
type Match interface {
	GetWinner() int
	GetScore() map[int]int // map playerid -> score
	GetNewSingle(int, map[int]int) *Single
	GetCurrentSingle() *Single
	AddSingles([]*Single) bool
}
