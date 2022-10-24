package riot_api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func MatchList(id string, srv string) []string {
	var matches []string
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
	request, _ := http.NewRequest("GET", ("https://" + region + ".api.riotgames.com/lol/match/v5/matches/by-puuid/" + id + "/ids?count=12"), nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&matches)
		if err != nil {
			log.Fatalf("error decoding response into []matches] \n%v", err)
		}
	}

	return matches
}
