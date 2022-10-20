package Models

type UserPost struct {
	Username  string `json:"username"`
	Server    string `json:"server"`
	Discordid string `json:"discordid"`
}

type UserDB struct {
	Discordid  string `json:"discordid"`
	Username   string `json:"username"`
	Server     string `json:"server"`
	Puuid      string `json:"puuid"`
	Id         string `json:"id"`
	RankedTier string `json:"RankedTier"`
}

type UserLookup struct {
	Username string `json:"username"`
	Server   string `json:"server"`
}

type LookupResponse struct {
	Username      string            `json:"username"`
	Tier          string            `json:"tier"`
	Rank          string            `json:"rank"`
	Level         int               `json:"level"`
	ProfileIconId int               `json:"profileIconId"`
	Champions     []ChampionMastery `json:"champions"`
	Wins          string            `json:"wins"`
	Losses        string            `json:"losses"`
}
