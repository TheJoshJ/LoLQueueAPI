package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/Models"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//find the ID of the user that we are currently creating
	params := mux.Vars(r)

	//create a new instance of a struct for us to process
	var newUser Models.UserPost
	var bySummonerName Models.RiotBySummonerName
	var rankedinfo Models.LeagueRanked

	//process the information sent from the bot into the struct
	newUser.Username = r.FormValue("username")
	newUser.Server = r.FormValue("server")

	log.Println(newUser)
	//get the puuid, encrypted id, and ranked tier from Riot's API
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + newUser.Username + "?api_key=RGAPI-a2f8c4e7-301f-4f91-99c5-cfd894d49233")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &bySummonerName)
	log.Println(bySummonerName)

	resp2, err := http.Get("https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + bySummonerName.Id + "?api_key=RGAPI-a2f8c4e7-301f-4f91-99c5-cfd894d49233")
	if err != nil {
		log.Fatalln(err)
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body2, &rankedinfo)
	log.Println(rankedinfo)
	//save the user struct with all the udpated information to PostgresDB
	newUserDB := Models.UserDB{
		params["uid"],
		newUser.Username,
		newUser.Server,
		bySummonerName.Puuid,
		bySummonerName.Id,
		rankedinfo.Tier,
	}

	//respond accordingly
	reply, err := json.Marshal(newUserDB)
	if err != nil {
		json.NewEncoder(w).Encode(reply)
	}
	log.Println(string(reply))
}
