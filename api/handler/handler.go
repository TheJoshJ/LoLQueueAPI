package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/api/riot_api"
	"main/models"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"data": "Pong"})
	if err != nil {
		log.Println("unable to encode response from Ping handler")
	}
}

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

func MatchGet(w http.ResponseWriter, r *http.Request) {
	var matchList []string
	matchData := make([]models.MatchData, 10)

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
	matchList = riot_api.MatchListByCount(bySummonerName.Puuid, userSearch.Server, 10)

	log.Printf("%#v", bySummonerName)
	log.Printf("%#v", matchList)

	for i, matchid := range matchList {
		matchData[i] = riot_api.MatchInfo(matchid, userSearch.Server)
	}

	reply, err := json.Marshal(matchData)
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

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var newUser models.UserPost

	//process the information sent via the PostForm request
	newUser.Discordid = r.FormValue("discordid")
	newUser.Username = r.FormValue("username")
	newUser.Server = r.FormValue("server")

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
