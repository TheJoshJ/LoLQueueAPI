package riot_api

import (
	"encoding/json"
	"log"
	"main/Models"
	"net/http"
	"os"
)

func MatchInfo(matchid string, srv string) Models.MatchData {
	var match Models.MatchData
	var region string
	if srv == "NA" || srv == "BR" || srv == "LAN" || srv == "LAS" {
		region = "americas"
	}
	if srv == "KR" || srv == "JP" {
		region = "asia"
	}
	if srv == "EUNE" || srv == "EUW" || srv == "TR" || srv == "RU" {
		region = "europe"
	}
	if srv == "OCE" {
		region = "sea"
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", ("https://" + region + ".api.riotgames.com/lol/match/v5/matches/" + matchid), nil)
	if err != nil {
		log.Println("error setting up request for match information")
	}
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, err := client.Do(request)
	if err != nil {
		log.Println("error attempting to get response for match information")
	}

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&match)
		if err != nil {
			log.Fatalf("error decoding response into []match] \n%v", err)
		}
	}
	return match
}
