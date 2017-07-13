package model

import (
	"fmt"
	"sync"

	"github.com/cnaize/lifland/util"
)

type Player struct {
	id string

	mu      sync.Mutex
	balance float64
}

func NewPlayer(id string) *Player {
	fmt.Printf("creating player %s\n", id)
	return &Player{
		id: id,
	}
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) IncrBalance(points float64) error {
	if points == 0.0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("player %s: increasing balance %f by %f points\n", p.Id(), p.balance, points)
	if util.Round(p.balance+points) < 0 {
		return fmt.Errorf("player %s can't apply increasing balance %f by %f points",
			p.Id(), p.balance, points)
	}
	p.balance += points
	return nil
}

func (p *Player) Balance() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.balance
}
