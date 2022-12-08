package models

type Db_server struct {
	id   string
	name string
}

type Db_server_user struct {
	server_id  string
	discord_id string
}

type Db_discord_user_riot_user struct {
	discord_id string
	puuid      string
}
type Db_discord_user struct {
	id       string
	username string
}

type Db_riot_user struct {
	puuid           string
	riot_account_id string
	username        string
	server          string
}

type Db_match struct {
	match_id   string
	puuid      string
	match_time string
}

type Db_score struct {
	id              int
	puuid           string
	matches_played  int
	gold            int
	assists         int
	ward_kills      int
	dragons_claimed int
	vision_score    int
	kills           int
	wins            int
}
