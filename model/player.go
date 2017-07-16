package model

import (
	"fmt"
	"sync"

	"github.com/cnaize/lifland/util"
)

// NOTE:
// fields open only for marshaling, don't use it directly
type Player struct {
	Id string `json:"id"`

	mu      sync.Mutex
	Balance float64 `json:"balance"`
}

func NewPlayer(id string) *Player {
	fmt.Printf("creating player %s\n", id)
	return &Player{
		Id: id,
	}
}

func (p *Player) GetId() string {
	return p.Id
}

func (p *Player) IncrBalance(points float64) error {
	if points == 0.0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("player %s: increasing balance %f by %f points\n",
		p.Id, p.Balance, points)
	if util.Round(p.Balance+points) < 0 {
		return fmt.Errorf("player %s can't apply increasing balance %f by %f points",
			p.Id, p.Balance, points)
	}
	p.Balance += points
	return nil
}

func (p *Player) GetBalance() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Balance
}
