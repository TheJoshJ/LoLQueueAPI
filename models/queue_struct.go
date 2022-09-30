package models

type Queue struct {
	Gamemode  string `json: "gamemode"`
	Primary   string `json: "primary"`
	Secondary string `json: "secondary"`
	Fill      string `json: "fill"`
}
