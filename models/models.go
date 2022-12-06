package models

//structs specific to LoLQueue API

type UserPost struct {
	Username        string `json:"username"`
	RiotServer      string `json:"riotServer"`
	DiscordID       string `json:"discordID"`
	DiscordServerID string `json:"discordServerID"`
}

type UserDB struct {
	Discordid       string `json:"discordID"`
	DiscordServerID string `json:"discordServerID"`
	Username        string `json:"username"`
	Server          string `json:"server"`
	Puuid           string `json:"puuid"`
	Id              string `json:"id"`
	RankedTier      string `json:"RankedTier"`
}

type UserLookup struct {
	Username string `json:"username"`
	Server   string `json:"server"`
}

type LookupResponse struct {
	Username      string                `json:"username"`
	Tier          string                `json:"tier"`
	Rank          string                `json:"rank"`
	Level         int                   `json:"level"`
	ProfileIconId int                   `json:"profileIconId"`
	Champions     []ChampionMasteryResp `json:"champions"`
	Wins          int                   `json:"wins"`
	Losses        int                   `json:"losses"`
}

type ChampionMasteryResp struct {
	ChampionName   string `json:"championName"`
	ChampionId     int    `json:"championId"`
	ChampionLevel  int    `json:"championLevel"`
	ChampionPoints int    `json:"championPoints"`
}

type MatchDataResp struct {
	GameID       string `json:"gameID"`
	GameMode     string `json:"gameMode"`
	ChampionName string `json:"championName"`
	Assists      int    `json:"assists"`
	Deaths       int    `json:"deaths"`
	Kills        int    `json:"kills"`
	Win          bool   `json:"win"`
}

type Leaderboard struct {
	DiscordServerID  string    `json:"discordServerID"`
	CurrentObjective string    `json:"currentObjective"`
	SecTillNextObj   int       `json:"secTillNextObj"`
	Leaders          []Leaders `json:"leaders"`
}
type Leaders struct {
	SummonerName string `json:"summonerName"`
	Position     int    `json:"position"`
	Score        string `json:"score"`
}

//structs specific to Riot API

type RiotBySummonerName struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type ChampionMastery struct {
	ChampionName                 string
	ChampionId                   int    `json:"championId"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	LastPlayTime                 int64  `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int    `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int    `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool   `json:"chestGranted"`
	TokensEarned                 int    `json:"tokensEarned"`
	SummonerId                   string `json:"summonerId"`
}

type LeagueRanked struct {
	LeagueId     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerId   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

type Champions struct {
	Type    string                  `json:"type"`
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]ChampionData `json:"data"`
}

type ChampionData struct {
	Version string            `json:"version"`
	ID      string            `json:"id"`
	Key     string            `json:"key"`
	Name    string            `json:"name"`
	Title   string            `json:"title"`
	Blurb   string            `json:"blurb"`
	Info    ChampionDataInfo  `json:"info"`
	Image   ImageData         `json:"image"`
	Tags    []string          `json:"tags"`
	Partype string            `json:"partype"`
	Stats   ChampionDataStats `json:"stats"`
}

type ChampionDataInfo struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type ImageData struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

type ChampionDataStats struct {
	HealthPoints                    float64 `json:"hp"`
	HealthPointsPerLevel            float64 `json:"hpperlevel"`
	ManaPoints                      float64 `json:"mp"`
	ManaPointsPerLevel              float64 `json:"mpperlevel"`
	MovementSpeed                   float64 `json:"movespeed"`
	Armor                           float64 `json:"armor"`
	ArmorPerLevel                   float64 `json:"armorperlevel"`
	SpellBlock                      float64 `json:"spellblock"`
	SpellBlockPerLevel              float64 `json:"spellblockperlevel"`
	AttackRange                     float64 `json:"attackrange"`
	HealthPointRegeneration         float64 `json:"hpregen"`
	HealthPointRegenerationPerLevel float64 `json:"hpregenperlevel"`
	ManaPointRegeneration           float64 `json:"mpregen"`
	ManaPointRegenerationPerLevel   float64 `json:"mpregenperlevel"`
	CriticalStrikeChance            float64 `json:"crit"`
	CriticalStrikeChancePerLevel    float64 `json:"critperlevel"`
	AttackDamage                    float64 `json:"attackdamage"`
	AttackDamagePerLevel            float64 `json:"attackdamageperlevel"`
	AttackSpeedOffset               float64 `json:"attackspeedoffset"`
	AttackSpeedPerLevel             float64 `json:"attackspeedperlevel"`
}

type MatchData struct {
	Metadata struct {
		DataVersion  string   `json:"dataVersion"`
		MatchId      string   `json:"matchId"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		GameCreation       int64          `json:"gameCreation"`
		GameDuration       int            `json:"gameDuration"`
		GameEndTimestamp   int64          `json:"gameEndTimestamp"`
		GameId             int64          `json:"gameId"`
		GameMode           string         `json:"gameMode"`
		GameName           string         `json:"gameName"`
		GameStartTimestamp int64          `json:"gameStartTimestamp"`
		GameType           string         `json:"gameType"`
		GameVersion        string         `json:"gameVersion"`
		MapId              int            `json:"mapId"`
		Participants       []Participants `json:"participants"`
		PlatformId         string         `json:"platformId"`
		QueueId            int            `json:"queueId"`
		Teams              []struct {
			Bans       []interface{} `json:"bans"`
			Objectives struct {
				Baron struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"baron"`
				Champion struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"champion"`
				Dragon struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"dragon"`
				Inhibitor struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"inhibitor"`
				RiftHerald struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"riftHerald"`
				Tower struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"tower"`
			} `json:"objectives"`
			TeamId int  `json:"teamId"`
			Win    bool `json:"win"`
		} `json:"teams"`
		TournamentCode string `json:"tournamentCode"`
	} `json:"info"`
}

type Participants struct {
	GameID                         string `json:"gameID"`
	GameMode                       string `json:"gameMode"`
	Assists                        int    `json:"assists"`
	BaronKills                     int    `json:"baronKills"`
	BasicPings                     int    `json:"basicPings"`
	BountyLevel                    int    `json:"bountyLevel"`
	ChampExperience                int    `json:"champExperience"`
	ChampLevel                     int    `json:"champLevel"`
	ChampionId                     int    `json:"championId"`
	ChampionName                   string `json:"championName"`
	ChampionTransform              int    `json:"championTransform"`
	ConsumablesPurchased           int    `json:"consumablesPurchased"`
	DamageDealtToBuildings         int    `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int    `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int    `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int    `json:"damageSelfMitigated"`
	Deaths                         int    `json:"deaths"`
	DetectorWardsPlaced            int    `json:"detectorWardsPlaced"`
	DoubleKills                    int    `json:"doubleKills"`
	DragonKills                    int    `json:"dragonKills"`
	EligibleForProgression         bool   `json:"eligibleForProgression"`
	FirstBloodAssist               bool   `json:"firstBloodAssist"`
	FirstBloodKill                 bool   `json:"firstBloodKill"`
	FirstTowerAssist               bool   `json:"firstTowerAssist"`
	FirstTowerKill                 bool   `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool   `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool   `json:"gameEndedInSurrender"`
	GoldEarned                     int    `json:"goldEarned"`
	GoldSpent                      int    `json:"goldSpent"`
	IndividualPosition             string `json:"individualPosition"`
	InhibitorKills                 int    `json:"inhibitorKills"`
	InhibitorTakedowns             int    `json:"inhibitorTakedowns"`
	InhibitorsLost                 int    `json:"inhibitorsLost"`
	Item0                          int    `json:"item0"`
	Item1                          int    `json:"item1"`
	Item2                          int    `json:"item2"`
	Item3                          int    `json:"item3"`
	Item4                          int    `json:"item4"`
	Item5                          int    `json:"item5"`
	Item6                          int    `json:"item6"`
	ItemsPurchased                 int    `json:"itemsPurchased"`
	KillingSprees                  int    `json:"killingSprees"`
	Kills                          int    `json:"kills"`
	Lane                           string `json:"lane"`
	LargestCriticalStrike          int    `json:"largestCriticalStrike"`
	LargestKillingSpree            int    `json:"largestKillingSpree"`
	LargestMultiKill               int    `json:"largestMultiKill"`
	LongestTimeSpentLiving         int    `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int    `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int    `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int    `json:"magicDamageTaken"`
	NeutralMinionsKilled           int    `json:"neutralMinionsKilled"`
	NexusKills                     int    `json:"nexusKills"`
	NexusLost                      int    `json:"nexusLost"`
	NexusTakedowns                 int    `json:"nexusTakedowns"`
	ObjectivesStolen               int    `json:"objectivesStolen"`
	ObjectivesStolenAssists        int    `json:"objectivesStolenAssists"`
	ParticipantId                  int    `json:"participantId"`
	PentaKills                     int    `json:"pentaKills"`
	PhysicalDamageDealt            int    `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int    `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int    `json:"physicalDamageTaken"`
	ProfileIcon                    int    `json:"profileIcon"`
	Puuid                          string `json:"puuid"`
	QuadraKills                    int    `json:"quadraKills"`
	RiotIdName                     string `json:"riotIdName"`
	RiotIdTagline                  string `json:"riotIdTagline"`
	Role                           string `json:"role"`
	SightWardsBoughtInGame         int    `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int    `json:"spell1Casts"`
	Spell2Casts                    int    `json:"spell2Casts"`
	Spell3Casts                    int    `json:"spell3Casts"`
	Spell4Casts                    int    `json:"spell4Casts"`
	Summoner1Casts                 int    `json:"summoner1Casts"`
	Summoner1Id                    int    `json:"summoner1Id"`
	Summoner2Casts                 int    `json:"summoner2Casts"`
	Summoner2Id                    int    `json:"summoner2Id"`
	SummonerId                     string `json:"summonerId"`
	SummonerLevel                  int    `json:"summonerLevel"`
	SummonerName                   string `json:"summonerName"`
	TeamEarlySurrendered           bool   `json:"teamEarlySurrendered"`
	TeamId                         int    `json:"teamId"`
	TeamPosition                   string `json:"teamPosition"`
	TimeCCingOthers                int    `json:"timeCCingOthers"`
	TimePlayed                     int    `json:"timePlayed"`
	TotalDamageDealt               int    `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int    `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int    `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int    `json:"totalDamageTaken"`
	TotalHeal                      int    `json:"totalHeal"`
	TotalHealsOnTeammates          int    `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int    `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int    `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int    `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int    `json:"totalUnitsHealed"`
	TripleKills                    int    `json:"tripleKills"`
	TrueDamageDealt                int    `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int    `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int    `json:"trueDamageTaken"`
	TurretKills                    int    `json:"turretKills"`
	TurretTakedowns                int    `json:"turretTakedowns"`
	TurretsLost                    int    `json:"turretsLost"`
	UnrealKills                    int    `json:"unrealKills"`
	VisionScore                    int    `json:"visionScore"`
	VisionWardsBoughtInGame        int    `json:"visionWardsBoughtInGame"`
	WardsKilled                    int    `json:"wardsKilled"`
	WardsPlaced                    int    `json:"wardsPlaced"`
	Win                            bool   `json:"win"`
}
