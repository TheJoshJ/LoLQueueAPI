package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/api/riot_api"
	"main/models"
	"net/http"
)

// Ping godoc
// @Summary      Pings the API service to ensure that it is active
// @Description  Ping the API service
// @Tags         utility
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      404
// @Router       /ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"data": "Pong"})
	if err != nil {
		log.Println("unable to encode response from Ping handler")
	}
}

// ProfileLookup godoc
// @Summary      Show an account
// @Description  Gets the users account information by their Username and Server
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        srv   path      string  true  "Riot Server"
// @Param        usr   path      string  true  "Username"
// @Success      200   {object}  models.LookupResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /lookup/{srv}/{usr} [get]
func ProfileLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var userSearch models.UserLookup

	//process the information sent via the PostForm request
	vars := mux.Vars(r)
	userSearch.Username = vars["usr"]
	userSearch.Server = vars["srv"]

	if userSearch.Username == "" || userSearch.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bySummonerName := riot_api.GetBySummonerName(userSearch.Username, userSearch.Server)
	rankedinfo := riot_api.GetRankedInfo(bySummonerName.Id, userSearch.Server)
	championMastery := riot_api.GetChampionMastery(bySummonerName.Id, userSearch.Server)

	if bySummonerName.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	lookupResponse := models.LookupResponse{
		Username:      bySummonerName.Name,
		Rank:          rankedinfo.Rank,
		Tier:          rankedinfo.Tier,
		Champions:     championMastery,
		Level:         bySummonerName.SummonerLevel,
		ProfileIconId: bySummonerName.ProfileIconId,
		Wins:          rankedinfo.Wins,
		Losses:        rankedinfo.Losses,
	}

	reply, err := json.Marshal(lookupResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetRecentMatches godoc
// @Summary      Show recent matches
// @Description  Show the past 10 matches
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        srv   path      string  true  "Riot Server"
// @Param        usr   path      string  true  "Username"
// @Success      200  {array}    models.MatchDataResp
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /match/{srv}/{usr} [get]
func GetRecentMatches(w http.ResponseWriter, r *http.Request) {
	var matchList []string
	matchesData := make([]models.MatchData, 10)
	matchDataReturn := make([]models.Participants, 10)
	matchDataResp := make([]models.MatchDataResp, 10)

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var userSearch models.UserLookup

	//process the information sent via the PostForm request
	vars := mux.Vars(r)
	userSearch.Username = vars["usr"]
	userSearch.Server = vars["srv"]

	if userSearch.Username == "" || userSearch.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bySummonerName := riot_api.GetBySummonerName(userSearch.Username, userSearch.Server)
	matchList = riot_api.MatchList(bySummonerName.Puuid, userSearch.Server)

	for i, matchid := range matchList {
		matchesData[i] = riot_api.MatchInfo(matchid, userSearch.Server)
	}

	for i, mdata := range matchesData {
		for _, participant := range mdata.Info.Participants {
			if participant.Puuid == bySummonerName.Puuid {
				matchDataReturn[i] = participant
				log.Println(i)
				log.Println(matchDataReturn[i])
			}
			matchDataResp[i].GameID = matchesData[i].Metadata.MatchId
			matchDataResp[i].GameMode = matchesData[i].Info.GameMode
			matchDataResp[i].Assists = participant.Assists
			matchDataResp[i].Deaths = participant.Deaths
			matchDataResp[i].Kills = participant.Kills
			matchDataResp[i].Win = participant.Win
		}
	}

	reply, err := json.Marshal(matchDataResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// CreateUser godoc
// @Summary      Create an account
// @Description  Creates and stores the users data to be used when executing commands/api calls.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.UserDB
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var newUser models.UserPost
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		log.Printf("error decoding user post resposne\n %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if newUser.Discordid == "" || newUser.Username == "" || newUser.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check to see if the user already exists

	//get the puuid, encrypted id, and ranked tier from Riot's API
	bySummonerName := riot_api.GetBySummonerName(newUser.Username, newUser.Server)
	rankedinfo := riot_api.GetRankedInfo(bySummonerName.Id, newUser.Server)

	//account for players who do not play ranked or does not exist
	if bySummonerName.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if rankedinfo.Tier == "" {
		rankedinfo.Tier = "UNRANKED"
	}

	//save the user struct with all the udpated information to PostgresDB
	newUserDB := models.UserDB{
		Discordid:  newUser.Discordid,
		Username:   newUser.Username,
		Server:     newUser.Server,
		Puuid:      bySummonerName.Puuid,
		Id:         bySummonerName.Id,
		RankedTier: rankedinfo.Tier,
	}

	//add the data to the database

	//respond accordingly
	reply, err := json.Marshal(newUserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
