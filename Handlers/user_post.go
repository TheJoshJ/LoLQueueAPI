package Handlers

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	"main/Models"
	"net/http"
	"os"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var newUser Models.UserPost

	//process the information sent via the PostForm request
	newUser.Discordid = r.FormValue("discordid")
	newUser.Username = r.FormValue("username")
	newUser.Server = r.FormValue("server")

	//check to see if the user already exists

	//get the puuid, encrypted id, and ranked tier from Riot's API
	bySummonerName := getBySummonerName(newUser.Username)
	rankedinfo := getRankedInfo(bySummonerName.Id)

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
		newUser.Discordid,
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

	if response.StatusCode != 404 {

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &bySummonerName)

		log.Println(bySummonerName)

		return bySummonerName
	}
	return bySummonerName
}

func getRankedInfo(Id string) Models.LeagueRanked {

	var rankedinfo Models.LeagueRanked
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

	if response.ContentLength == 0 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(body))

		err = json.Unmarshal(body, &rawData)
		if err != nil {
			log.Printf("error unmarshalling %v", err)
		}
		err = mapstructure.Decode(rawData[0], &rankedinfo)
		if err != nil {
			log.Printf("error decoding %v", err)
		}
		return rankedinfo
	} else {
		return rankedinfo
	}
}
