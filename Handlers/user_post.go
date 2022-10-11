package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
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

	//log.Println(newUser)

	//get the puuid, encrypted id, and ranked tier from Riot's API
	bySummonerName := getBySummonerName(newUser.Username)
	rankedinfo := getRankedInfo(bySummonerName.Id)

	//account for players who do not play ranked
	if rankedinfo.Tier == "" {
		rankedinfo.Tier = "UNRANKED"
	}

	//log.Printf("This is what we have prior to the struct being added to the DB \n%#v", rankedinfo)

	//save the user struct with all the udpated information to PostgresDB
	newUserDB := Models.UserDB{
		params["id"],
		newUser.Username,
		newUser.Server,
		bySummonerName.Puuid,
		bySummonerName.Id,
		rankedinfo.Tier,
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

func getRankedInfo(Id string) Models.LeagueRanked {

	var rawData []map[string]interface{}

	client := &http.Client{}
	request, err := http.NewRequest("GET", ("https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + Id), nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	err = json.Unmarshal(body, &rawData)

	var rankedinfo Models.LeagueRanked
	err = mapstructure.Decode(rawData[0], &rankedinfo)
	if err != nil {
		log.Fatalln(err)
	}

	return rankedinfo
}
