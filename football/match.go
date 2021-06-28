package football

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"math/rand"
	"time"
)

type MatchStatus string

const (
	FirstHalf  MatchStatus = "first_half"
	SecondHalf MatchStatus = "second_half"
	Finished   MatchStatus = "finished"
)

const (
	probabilityBase = 100000
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	goals = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "football_match_goals",
		Help: "goals scored by Players",
	}, []string{"matchStartTime", "team", "player"})

	shoots = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "football_match_shoots",
		Help: "shoots by Players",
	}, []string{"matchStartTime", "team", "player"})

	goalsByPlayer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "football_player_goal_minutes",
		Help:    "player goals by match minute",
		Buckets: []float64{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85},
	}, []string{"team", "player"})

	shootsSummaryByTeam = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "football_team_shoots_seconds",
		Help: "minute in which team scored a goal",
	}, []string{"team"})

	secondsFromStart = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "football_match_seconds",
		Help: "seconds since the first whistle",
	}, []string{"matchStartTime", "hosts", "guests"})

	currentHalf = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "football_match_current_half",
		Help: "current half or 0 if match is finished",
	}, []string{"matchStartTime", "hosts", "guests"})

	currentStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "football_match_current_status",
		Help: "current status of the match",
	}, []string{"matchStartTime", "hosts", "guests", "status"})
)

type match struct {
	Hosts     Team
	Guests    Team
	StartTime time.Time
}

type Match interface {
	Play()
	CurrentStatus() MatchStatus
}

func CreateMatch(hosts Team, guests Team) Match {
	return match{hosts, guests, time.Now()}
}

func (m match) Play() {
	m.initMatch()
	go m.keepStatusUpdated()
	for m.CurrentStatus() != Finished {
		m.exposeStatus()
		m.randomizeEvents()
		time.Sleep(1 * time.Minute)
	}

	m.exposeStatus()
	time.Sleep(5 * time.Minute)
	m.cleanup()
}

func (m match) keepStatusUpdated() {
	for m.CurrentStatus() != Finished {
		m.exposeStatus()
		time.Sleep(1 * time.Second)
	}
}

func (m match) initMatch() {
	m.initTeamGoals(&m.Hosts)
	m.initTeamGoals(&m.Guests)
	m.exposeStatus()
}

func (m match) exposeStatus() {
	matchDuration := time.Now().Sub(m.StartTime).Seconds()
	secondsFromStart.WithLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name).Set(matchDuration)
	currentHalf.WithLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name).Set(float64(m.currentHalf()))
	m.exposeCurrentStatus(FirstHalf)
	m.exposeCurrentStatus(SecondHalf)
	m.exposeCurrentStatus(Finished)
}

func (m match) randomizeEvents() {
	m.randomizeTeamEvents(&m.Hosts, &m.Guests)
	m.randomizeTeamEvents(&m.Guests, &m.Hosts)
}

func (m match) randomizeTeamEvents(team *Team, opponent *Team) {
	for _, player := range team.Players {
		if randomBool(player.Offence) {
			shoots.WithLabelValues(m.startTimeString(), team.Name, player.Name).Inc()
			shootsSummaryByTeam.WithLabelValues(team.Name).Observe(time.Now().Sub(m.StartTime).Seconds())
			if !randomBool(opponent.Defence) {
				goals.WithLabelValues(m.startTimeString(), team.Name, player.Name).Inc()
				goalsByPlayer.WithLabelValues(team.Name, player.Name).Observe(time.Now().Sub(m.StartTime).Minutes())
			}
		}

		if randomBool(player.BadLuck) {
			goals.WithLabelValues(m.startTimeString(), opponent.Name, player.Name).Inc()
		}
	}
}

func randomBool(probability float64) bool {
	return float64(r.Intn(probabilityBase)) > float64(probabilityBase)-float64(probabilityBase)*probability
}

func (m match) exposeCurrentStatus(status MatchStatus) {
	isInCurrentState := m.CurrentStatus() == status
	currentStatus.WithLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name, string(status)).Set(toNumericStatus(isInCurrentState))
}

func (m match) deleteCurrentStatus(status MatchStatus) {
	currentStatus.DeleteLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name, string(status))
}

func toNumericStatus(status bool) float64 {
	if status {
		return 1
	}

	return 0
}

func (m match) currentHalf() int {
	switch m.CurrentStatus() {
	case FirstHalf:
		return 1
	case SecondHalf:
		return 2
	default:
		return 0
	}
}

func (m match) cleanup() {
	m.cleanupTeam(&m.Hosts)
	m.cleanupTeam(&m.Guests)
	secondsFromStart.DeleteLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name)
	currentHalf.DeleteLabelValues(m.startTimeString(), m.Hosts.Name, m.Guests.Name)
	m.deleteCurrentStatus(FirstHalf)
	m.deleteCurrentStatus(SecondHalf)
	m.deleteCurrentStatus(Finished)
}

func (m match) startTimeString() string {
	return m.StartTime.Format(time.RFC3339Nano)
}

func (m match) CurrentStatus() MatchStatus {
	firstHalfDeadline := m.StartTime.Add(45 * time.Minute)
	secondHalfDeadline := m.StartTime.Add(90 * time.Minute)

	if time.Now().Before(firstHalfDeadline) {
		return FirstHalf
	}

	if time.Now().Before(secondHalfDeadline) {
		return SecondHalf
	}

	return Finished
}

func (m match) initTeamGoals(t *Team) {
	for _, player := range t.Players {
		goals.WithLabelValues(m.startTimeString(), t.Name, player.Name).Add(0)
		shoots.WithLabelValues(m.startTimeString(), t.Name, player.Name).Add(0)
	}
}

func (m match) cleanupTeam(t *Team) {
	for _, player := range t.Players {
		goals.DeleteLabelValues(m.startTimeString(), m.Hosts.Name, player.Name)
		goals.DeleteLabelValues(m.startTimeString(), m.Guests.Name, player.Name)
		shoots.DeleteLabelValues(m.startTimeString(), t.Name, player.Name)
	}
}
