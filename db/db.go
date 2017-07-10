package db

import (
	"fmt"
	"sync"

	"github.com/cnaize/lifland/model"
)

type DB struct {
	pmu         sync.Mutex
	players     map[string]*model.Player
	tmu         sync.Mutex
	tournaments map[int]*model.Tournament
	fmu         sync.Mutex
	funds       []model.Fund
}

var _ Interface = NewDB()

func NewDB() *DB {
	fmt.Println("creating db")
	return &DB{
		players:     make(map[string]*model.Player),
		tournaments: make(map[int]*model.Tournament),
	}
}

func (db *DB) GetPlayer(id string) *model.Player {
	db.pmu.Lock()
	defer db.pmu.Unlock()

	if player, ok := db.players[id]; ok {
		return player
	}
	return nil
}

func (db *DB) AddPlayer(player *model.Player) error {
	if player == nil {
		return fmt.Errorf("AddPlayer: player is nil")
	}

	db.pmu.Lock()
	defer db.pmu.Unlock()

	if _, ok := db.players[player.Id()]; ok {
		return fmt.Errorf("AddPlayer: player %s already exists", player.Id())
	}
	db.players[player.Id()] = player
	return nil
}

func (db *DB) DelPlayer(player *model.Player) {
	if player == nil {
		return
	}

	db.pmu.Lock()
	defer db.pmu.Unlock()

	delete(db.players, player.Id())
}

func (db *DB) GetTournament(id int) *model.Tournament {
	db.tmu.Lock()
	defer db.tmu.Unlock()

	if tournament, ok := db.tournaments[id]; ok {
		return tournament
	}
	return nil
}

func (db *DB) AddTournament(tournament *model.Tournament) error {
	if tournament == nil {
		return fmt.Errorf("AddTournament: tournament is nil")
	}

	db.tmu.Lock()
	defer db.tmu.Unlock()

	if _, ok := db.tournaments[tournament.Id()]; ok {
		return fmt.Errorf("AddTournament: tournament %d already exists", tournament.Id())
	}
	db.tournaments[tournament.Id()] = tournament
	return nil
}

func (db *DB) GetOldestTournament() *model.Tournament {
	db.tmu.Lock()
	defer db.tmu.Unlock()

	var tournament *model.Tournament
	for _, t := range db.tournaments {
		if tournament == nil || t.StartTime().Before(tournament.StartTime()) {
			tournament = t
		}
	}
	return tournament
}

func (db *DB) DelTournament(tournament *model.Tournament) {
	if tournament == nil {
		return
	}

	db.tmu.Lock()
	defer db.tmu.Unlock()

	delete(db.tournaments, tournament.Id())
}

func (db *DB) GetFunds() []model.Fund {
	db.fmu.Lock()
	defer db.fmu.Unlock()

	return db.funds
}

func (db *DB) SetFunds(funds []model.Fund) {
	db.fmu.Lock()
	defer db.fmu.Unlock()

	db.funds = funds
}

func (db *DB) AddFund(fund model.Fund) error {
	if fund == nil {
		return fmt.Errorf("AddFund: fund is nil")
	}

	db.fmu.Lock()
	defer db.fmu.Unlock()

	db.funds = append(db.funds, fund)
	return nil
}

func (db *DB) Reset() {
	db.pmu.Lock()
	db.tmu.Lock()
	db.fmu.Lock()
	defer db.pmu.Unlock()
	defer db.tmu.Unlock()
	defer db.fmu.Unlock()

	db.players = make(map[string]*model.Player)
	db.tournaments = make(map[int]*model.Tournament)
	db.funds = []model.Fund{}
	fmt.Println("db reseted")
}
