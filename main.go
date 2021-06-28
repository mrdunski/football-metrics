package main

import (
	"fmt"
	"football-metrics/football"
	"football-metrics/randomizer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var poland = football.Team{
	Name: "Polska",
	Players: []football.Player{
		{
			Name:    "Robert Lewandowski",
			Offence: 0.05,
			BadLuck: 0,
		},
		{
			Name:    "Wojciech Szczęsny",
			Offence: 0,
			BadLuck: 0.03,
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

func randomMatcher(teams []football.Team) {
	hosts := teams[rand.Intn(len(teams))]
	guests := teams[rand.Intn(len(teams))]

	for hosts.Name == guests.Name {
		guests = teams[rand.Intn(len(teams))]
	}

	football.CreateMatch(hosts, guests).Play()
}

func runTournament() {
	teams := make([]football.Team, len(otherTeams)+1)
	teams[0] = poland

	for i, team := range otherTeams {
		teams[i+1] = randomizer.RandomTeam(team)
	}

	for true {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		randomMatcher(teams)
		time.Sleep(time.Duration(rand.Intn(90)) * time.Minute)
	}
}

func main() {
	address := getListeningAddress()
	fmt.Printf("Listening on %s\n", address)

	go runTournament()
	go runTournament()
	go runTournament()
	go runTournament()

	http.HandleFunc("/", home)
	http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(address, nil))
}

func getListeningAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}

	return ":" + port
}

func home(w http.ResponseWriter, r *http.Request) {
	writeResponse("Hello!", w)
}

func health(w http.ResponseWriter, r *http.Request) {
	writeResponse("ok", w)
}

func writeResponse(r string, w http.ResponseWriter) {
	_, err := w.Write([]byte(r))
	if err != nil {
		log.Printf("Couldn't reply: %s\n", err)
	}
}
