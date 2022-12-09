package models

type Server struct {
	Id   string
	Name string
}

type Server_user struct {
	Server_id  string
	Discord_id string
}

type Discord_user_riot_user struct {
	Discord_id string
	Puuid      string
}

type Discord_user struct {
	Id       string
	Username string
}

type Riot_user struct {
	Puuid           string
	Riot_account_id string
	Username        string
	Server          string
}

type Match struct {
	Match_id   string
	Puuid      string
	Match_time string
}

type Score struct {
	Id              int
	Puuid           string
	Matches_played  int
	Gold            int
	Assists         int
	Ward_kills      int
	Dragons_claimed int
	Vision_score    int
	Kills           int
	Wins            int
}
