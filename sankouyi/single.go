package sankouyi

import "github.com/sansna/cicily/v2/models"

type SKYSingle struct {
	winner   int
	finished bool
	desktop  *SKYDesktop
}

func (s *SKYSingle) GetWinner() int {
	if s.finished {
		return s.winner
	}
	return -1
}

func (s *SKYSingle) SetWinner(winner int) {
	s.winner = winner
	s.finished = true
	return
}

func (s *SKYSingle) GetDesktop() *models.Desktop {
	if s.desktop == nil {
		s.desktop = new(SKYDesktop)
	}
	r := models.Desktop(s.desktop)
	return &r
}
