package randomizer

import "football-metrics/football"

var poland = football.Team{
	Name: "Polska",
	Players: []football.Player{
		{
			Name:    "Robert Lewandowski",
			Offence: 0.25,
			BadLuck: 0,
		},
		{
			Name:    "Wojciech Szczęsny",
			Offence: 0,
			BadLuck: 0.05,
		},
		{
			Name:    "Jakub Błaszczykowski",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Arkadiusz Milik",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Grzegorz Krychowiak",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Kamil Glik",
			Offence: 0,
			BadLuck: 0,
		},
		{
			Name:    "Piotr Zieliński",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Kamil Jóźwiak",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Karol Świderski",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Jakub Moder",
			Offence: 0.01,
			BadLuck: 0,
		},
		{
			Name:    "Michał Helik",
			Offence: 0,
			BadLuck: 0,
		},
	},
	Defence: 0.5,
}

var otherTeams = []string{
	"Niemcy",
	"Włochy",
	"Hiszpania",
	"Wielka Brytania",
	"Norwegia",
	"Szwecja",
	"Portugalia",
	"Białoruś",
	"Czechy",
	"Litwa",
	"Dania",
	"Białoruś",
	"Chorwacja",
}

func RandomTeams() []football.Team {
	teams := make([]football.Team, len(otherTeams)+1)
	teams[0] = poland

	for i, team := range otherTeams {
		teams[i+1] = RandomTeam(team)
	}

	return teams
}
