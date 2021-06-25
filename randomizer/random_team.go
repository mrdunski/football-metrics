package randomizer

import (
	"football-metrics/football"
	"github.com/Pallinder/go-randomdata"
	"math/rand"
)

func RandomTeam(name string) football.Team {
	players := make([]football.Player, 11)

	for i := 0; i < 11; i++ {
		players[i] = randomPlayer()
	}

	return football.Team{
		Name:    name,
		Players: players,
		Defence: rand.Float64(),
	}
}

func randomPlayer() football.Player {
	name := randomdata.FullName(randomdata.Male)

	return football.Player{
		Name:    name,
		BadLuck: rand.Float64() * 0.01,
		Offence: rand.Float64() * 0.1,
	}
}
