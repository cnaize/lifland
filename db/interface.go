package db

import "github.com/cnaize/lifland/model"

type Interface interface {
	GetPlayer(id string) *model.Player
	AddPlayer(player *model.Player) error
	DelPlayer(player *model.Player)

	GetTournament(id int) *model.Tournament
	AddTournament(tournament *model.Tournament) error
	DelTournament(tournament *model.Tournament)
	GetOldestTournament() *model.Tournament

	AddFund(fund model.Fund) error
	SyncFunds()

	Dump()
	Restore()
	Reset()

	SetDebug(debug bool)
}
