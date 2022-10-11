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
