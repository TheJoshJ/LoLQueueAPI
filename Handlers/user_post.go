package Handlers

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"main/Models"
	"main/riot_api"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var newUser Models.UserPost

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
	newUserDB := Models.UserDB{
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
