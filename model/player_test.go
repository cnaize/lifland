package model

import "testing"

func TestPlayer(t *testing.T) {
	player := NewPlayer("new")
	player.IncrBalance(5)
	if player.Balance() != 5.0 {
		t.Errorf("invalid balance: want %f, got %f", 5.0, player.Balance())
	}
	if err := player.IncrBalance(-10); err == nil {
		t.Errorf("invalid incr no errors")
	}
	player.IncrBalance(5)
	if player.Balance() != 10.0 {
		t.Errorf("invalid balance: want %f, got %f", 10.0, player.Balance())
	}
	player.IncrBalance(0)
	if player.Balance() != 10.0 {
		t.Errorf("invalid balance: want %f, got %f", 10.0, player.Balance())
	}
	player.IncrBalance(-10.0)
	if player.Balance() != 0.0 {
		t.Errorf("invalid balance: want %f, got %f", 0.0, player.Balance())
	}
}
