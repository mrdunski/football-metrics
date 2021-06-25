package football

type Team struct {
	Name    string
	Players []Player
	Defence float64
}

type Player struct {
	Name    string
	Offence float64
	BadLuck float64
}
