package model

import (
	"fmt"
	"sync"
	"time"
)

// NOTE:
// fields opened only for marshaling, don't use it directly
type Tournament struct {
	Id        int       `json:"id"`
	Deposit   float64   `json:"deposit"`
	StartTime time.Time `json:"startTime"`

	mu   sync.Mutex
	Open bool `json:"open"`
	// backers (including player) and their income by the player id
	Funds map[string]Fund `json:"funds"`
}

func NewTournament(id int, deposit float64) *Tournament {
	fmt.Printf("creating tournament %d, deposit: %f\n", id, deposit)
	return &Tournament{
		Id:        id,
		Deposit:   deposit,
		StartTime: time.Now(),
		Open:      true,
		Funds:     make(map[string]Fund),
	}
}

func (t *Tournament) GetId() int {
	return t.Id
}

func (t *Tournament) GetDeposit() float64 {
	return t.Deposit
}

func (t *Tournament) GetStartTime() time.Time {
	return t.StartTime
}

func (t *Tournament) AddPlayer(id string, fund Fund) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Open {
		return fmt.Errorf("Tournament %d already closed", t.Id)
	}
	if fund == nil {
		return fmt.Errorf("Player %s trying to join tournament %d without fund", id, t.Id)
	}
	if _, ok := t.Funds[id]; ok {
		return fmt.Errorf("Player %s already joined tournament %d", id, t.Id)
	}

	fmt.Printf("player %s joined tournament %d\n", id, t.Id)
	t.Funds[id] = fund
	return nil
}

func (t *Tournament) HasPlayer(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.Funds[id]
	return ok
}

func (t *Tournament) Close() (map[string]Fund, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Open {
		return nil, fmt.Errorf("Tournament %d already closed", t.Id)
	}
	t.Open = false
	return t.Funds, nil
}

func (t *Tournament) IsOpen() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.Open
}
