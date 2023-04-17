package server

type Player struct {
	Name   string
	Games  []int
	Status string
}

var Plrs = map[string]Player{}
