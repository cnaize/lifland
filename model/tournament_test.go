package model

import (
	"reflect"
	"testing"
)

func TestTournament(t *testing.T) {
	player100 := NewPlayer("player100")
	player100.IncrBalance(100)
	ifund := Fund{NewPlayer("invalid"): 100}
	fund := Fund{player100: 100}
	wfunds := map[string]Fund{player100.Id(): fund}

	tourn := NewTournament(1, 100)
	if tourn.Id() != 1 {
		t.Errorf("invalid tournament id: want %d, got %d", 1, tourn.Id())
	}
	if tourn.Deposit() != 100.0 {
		t.Errorf("invalid tourn deposit: want %f, got %f", 100.0, tourn.Deposit())
	}
	if err := tourn.AddPlayer(player100.Id(), fund); err != nil {
		t.Errorf("player %s can't join to the tournament: %+v", player100.Id(), err)
	}
	if err := tourn.AddPlayer(player100.Id(), fund); err == nil {
		t.Errorf("double joined player %s", player100.Id())
	}
	if !tourn.IsOpen() {
		t.Errorf("tournament closed")
	}
	if !tourn.HasPlayer(player100.Id()) {
		t.Errorf("player %s not in the tournament", player100.Id())
	}
	if err := tourn.AddPlayer("invalid", nil); err == nil {
		t.Errorf("invlid player joined to the tournament")
	}
	funds, _ := tourn.Close()
	if !reflect.DeepEqual(funds, wfunds) {
		t.Errorf("invalid funds: want %#v, got %#v", wfunds, funds)
	}
	if _, err := tourn.Close(); err == nil {
		t.Errorf("double closing tournament")
	}
	if err := tourn.AddPlayer("invalid", ifund); err == nil {
		t.Errorf("player joined closed tournament")
	}
}
