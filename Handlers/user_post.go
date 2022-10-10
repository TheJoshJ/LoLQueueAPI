package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/Models"
	"net/http"
	"os"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//find the ID of the user that we are currently creating
	params := mux.Vars(r)

	//create a new instance of a struct for us to process
	var newUser Models.UserPost

	//process the information sent from the bot into the struct
	newUser.Username = r.FormValue("username")
	newUser.Server = r.FormValue("server")

	//get the puuid, encrypted id from Riot's API
	bySummonerName := getBySummonerName(newUser.Username)

	//save the user struct with all the udpated information to PostgresDB
	newUserDB := Models.UserDB{
		params["uid"],
		newUser.Username,
		newUser.Server,
		bySummonerName.Puuid,
		bySummonerName.Id,
	}

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
	w.WriteHeader(http.StatusOK)
	log.Println(string(reply))
	return

}

func getBySummonerName(Username string) Models.RiotBySummonerName {

	var bySummonerName Models.RiotBySummonerName

	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + Username), nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &bySummonerName)

	return bySummonerName
}
