package models

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Played   int    `json:"played"`
	Won      int    `json:"won"`
	Drawn    int    `json:"drawn"`
	Lost     int    `json:"lost"`
	GF       int    `json:"gf"`
	GA       int    `json:"ga"`
	GD       int    `json:"gd"`
	Points   int    `json:"points"`
	Strength int    `json:"strength"`
}
