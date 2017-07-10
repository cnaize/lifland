package model

import (
	"fmt"
	"sync"
	"time"
)

type Fund map[*Player]float64

type Tournament struct {
	id        int
	deposit   float64
	startTime time.Time

	mu   sync.Mutex
	open bool
	// backers (including player) and their income by the player id
	funds map[string]Fund
}

func NewTournament(id int, deposit float64) *Tournament {
	fmt.Printf("creating tournament %d, deposit: %f\n", id, deposit)
	return &Tournament{
		id:        id,
		deposit:   deposit,
		startTime: time.Now(),
		open:      true,
		funds:     make(map[string]Fund),
	}
}

func (t *Tournament) Id() int {
	return t.id
}

func (t *Tournament) Deposit() float64 {
	return t.deposit
}

func (t *Tournament) StartTime() time.Time {
	return t.startTime
}

func (t *Tournament) AddPlayer(id string, fund Fund) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.open {
		return fmt.Errorf("Tournament %d already closed", t.Id())
	}
	if fund == nil {
		return fmt.Errorf("Player %s trying to join tournament %d without fund", id, t.Id())
	}
	if _, ok := t.funds[id]; ok {
		return fmt.Errorf("Player %s already joined tournament %d", id, t.Id())
	}

	fmt.Printf("player %s joined tournament %d\n", id, t.Id())
	t.funds[id] = fund
	return nil
}

func (t *Tournament) HasPlayer(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.funds[id]
	return ok
}

func (t *Tournament) Close() (map[string]Fund, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.open {
		return nil, fmt.Errorf("Tournament %d already closed", t.Id())
	}
	return t.funds, nil
}

func (t *Tournament) IsOpen() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.open
}
