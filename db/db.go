package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/cnaize/lifland/model"
)

const (
	dumpFileName string = "dump.db"
)

// NOTE:
// fields open only for marshaling, don't use it directly
type DB struct {
	debug       bool
	pmu         sync.Mutex
	Players     map[string]*model.Player `json:"players,omitempty"`
	tmu         sync.Mutex
	Tournaments map[int]*model.Tournament `json:"tournaments,omitempty"`
	fmu         sync.Mutex
	Funds       []model.Fund `jons:"funds,omitempty"`
}

var _ Interface = NewDB()

func NewDB() *DB {
	return &DB{
		Players:     make(map[string]*model.Player),
		Tournaments: make(map[int]*model.Tournament),
	}
}

func (db *DB) GetPlayer(id string) *model.Player {
	db.pmu.Lock()
	defer db.pmu.Unlock()

	if player, ok := db.Players[id]; ok {
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

	if _, ok := db.Players[player.GetId()]; ok {
		return fmt.Errorf("AddPlayer: player %s already exists", player.GetId())
	}
	db.Players[player.GetId()] = player
	return nil
}

func (db *DB) DelPlayer(player *model.Player) {
	if player == nil {
		return
	}

	db.pmu.Lock()
	defer db.pmu.Unlock()

	delete(db.Players, player.GetId())
}

func (db *DB) GetTournament(id int) *model.Tournament {
	db.tmu.Lock()
	defer db.tmu.Unlock()

	if tournament, ok := db.Tournaments[id]; ok {
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

	if _, ok := db.Tournaments[tournament.Id]; ok {
		return fmt.Errorf("AddTournament: tournament %d already exists", tournament.Id)
	}
	db.Tournaments[tournament.Id] = tournament
	return nil
}

func (db *DB) GetOldestTournament() *model.Tournament {
	db.tmu.Lock()
	defer db.tmu.Unlock()

	var tournament *model.Tournament
	for _, t := range db.Tournaments {
		if !t.IsOpen() {
			continue
		}
		if tournament == nil || t.GetStartTime().Before(tournament.GetStartTime()) {
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

	delete(db.Tournaments, tournament.Id)
}

func (db *DB) AddFund(fund model.Fund) error {
	if fund == nil {
		return fmt.Errorf("AddFund: fund is nil")
	}

	db.fmu.Lock()
	defer db.fmu.Unlock()

	db.Funds = append(db.Funds, fund)
	return nil
}

func (db *DB) SyncFunds() {
	synced := false
	db.fmu.Lock()
	defer func() {
		if synced {
			db.lockAll(&db.fmu)
			defer db.unlockAll(&db.fmu)
			db.dump()
		}
		db.fmu.Unlock()
	}()

	var funds []model.Fund
	for _, fund := range db.Funds {
		fmt.Println("syncing funds")
		for playerId, points := range fund {
			player := db.GetPlayer(playerId)
			if player == nil {
				fmt.Printf("ERROR: can't sync funds: player %s not found\n", playerId)
				continue
			}
			if err := player.IncrBalance(points); err == nil {
				fmt.Printf("funds %f for player %s synced\n", points, player.GetId())
				synced = true
				delete(fund, playerId)
			}
		}
		if len(fund) > 0 {
			funds = append(funds, fund)
		}
	}
	db.Funds = funds
}

func (db *DB) Reset() {
	db.lockAll(nil)
	defer db.unlockAll(nil)

	os.Remove(dumpFileName)
	db.Players = make(map[string]*model.Player)
	db.Tournaments = make(map[int]*model.Tournament)
	db.Funds = []model.Fund{}
	fmt.Println("db reseted")
}

func (db *DB) SetDebug(debug bool) {
	db.debug = !db.debug
}

func (db *DB) Dump() {
	db.lockAll(nil)
	defer db.unlockAll(nil)

	db.dump()
}

func (db *DB) Restore() {
	if b, err := ioutil.ReadFile(dumpFileName); err == nil {
		fmt.Println("restoring db")
		if err := json.Unmarshal(b, db); err != nil {
			fmt.Printf("db restore failed: %+v\n", err)
			return
		}
		fmt.Println("db restore: success")
		return
	}
}

// NOTE: not thread safe
func (db *DB) dump() {
	if db.debug {
		return
	}

	fmt.Println("dumping db")
	dump, err := json.Marshal(db)
	if err != nil {
		fmt.Printf("ERROR: db dump failed: can't marshal data: %+v\n", err)
		return
	}
	if err := ioutil.WriteFile(dumpFileName, dump, 0644); err != nil {
		fmt.Printf("ERROR: db dump failed: can't write to file: %+v\n", err)
		return
	}
	fmt.Println("db dump: success")
}

func (db *DB) lockAll(except *sync.Mutex) {
	for _, m := range []*sync.Mutex{&db.pmu, &db.tmu, &db.fmu} {
		if m != except {
			m.Lock()
		}
	}
}

func (db *DB) unlockAll(except *sync.Mutex) {
	for _, m := range []*sync.Mutex{&db.pmu, &db.tmu, &db.fmu} {
		if m != except {
			m.Unlock()
		}
	}
}
