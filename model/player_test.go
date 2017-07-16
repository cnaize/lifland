package model

import "testing"

func TestPlayer(t *testing.T) {
	player := NewPlayer("new")
	player.IncrBalance(5)
	if player.GetBalance() != 5.0 {
		t.Errorf("invalid balance: want %f, got %f", 5.0, player.GetBalance())
	}
	if err := player.IncrBalance(-10); err == nil {
		t.Errorf("invalid incr no errors")
	}
	player.IncrBalance(5)
	if player.GetBalance() != 10.0 {
		t.Errorf("invalid balance: want %f, got %f", 10.0, player.GetBalance())
	}
	player.IncrBalance(0)
	if player.GetBalance() != 10.0 {
		t.Errorf("invalid balance: want %f, got %f", 10.0, player.GetBalance())
	}
	player.IncrBalance(-10.0)
	if player.GetBalance() != 0.0 {
		t.Errorf("invalid balance: want %f, got %f", 0.0, player.GetBalance())
	}
}
