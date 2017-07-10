package model

import (
	"fmt"
	"math"
	"sync"
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
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("player %s: increasing balance %f by %f points\n", p.Id(), p.balance, points)
	if p.balance+points < 0 || p.balance+math.Abs(points) < p.balance {
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
