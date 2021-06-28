package randomizer

import (
	"football-metrics/football"
	"github.com/Pallinder/go-randomdata"
	"math"
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
		BadLuck: math.Abs(rand.NormFloat64()) * 0.02,
		Offence: math.Abs(rand.NormFloat64()) * 0.2,
	}
}
