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
	"sync"
	"time"
)

func play(hosts, guests football.Team, group *sync.WaitGroup) {
	time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	football.CreateMatch(hosts, guests).Play()
	group.Done()
}

func runTournament(concurrency int) {
	teams := randomizer.RandomTeams()

	for true {
		group := sync.WaitGroup{}
		rand.Shuffle(len(teams), func(a, b int) {
			teamA := teams[a]
			teamB := teams[b]

			teams[a] = teamB
			teams[b] = teamA

		})
		for i := 0; i < concurrency; i++ {
			hosts := teams[i*2]
			guests := teams[i*2+1]

			group.Add(1)
			go play(hosts, guests, &group)
		}

		group.Wait()
	}
}

func main() {
	address := getListeningAddress()
	fmt.Printf("Listening on %s\n", address)

	go runTournament(6)

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
