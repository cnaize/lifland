package model

import (
	"reflect"
	"testing"
)

func TestTournament(t *testing.T) {
	player100 := NewPlayer("player100")
	player100.IncrBalance(100)
	ifund := Fund{"invalid": 100}
	fund := Fund{player100.GetId(): 100}
	wfunds := map[string]Fund{player100.GetId(): fund}

	tourn := NewTournament(1, 100)
	if tourn.GetId() != 1 {
		t.Errorf("invalid tournament id: want %d, got %d", 1, tourn.GetId())
	}
	if tourn.GetDeposit() != 100.0 {
		t.Errorf("invalid tourn deposit: want %f, got %f", 100.0, tourn.GetDeposit())
	}
	if err := tourn.AddPlayer(player100.GetId(), fund); err != nil {
		t.Errorf("player %s can't join to the tournament: %+v", player100.GetId(), err)
	}
	if err := tourn.AddPlayer(player100.GetId(), fund); err == nil {
		t.Errorf("double joined player %s", player100.GetId())
	}
	if !tourn.IsOpen() {
		t.Errorf("tournament closed")
	}
	if !tourn.HasPlayer(player100.GetId()) {
		t.Errorf("player %s not in the tournament", player100.GetId())
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
