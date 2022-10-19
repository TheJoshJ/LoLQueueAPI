package riot_api

import (
	"encoding/json"
	"log"
	"main/Models"
	"net/http"
	"os"
	"strings"
)

func GetRankedInfo(Id string, server string) Models.LeagueRanked {

	var rankedArray []Models.LeagueRanked
	var rankedinfo Models.LeagueRanked
	var notRanked bool = true

	serverString := translateServerName(server)

	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("https://" + serverString + ".api.riotgames.com/lol/league/v4/entries/by-summoner/" + Id), nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	if response.StatusCode == 403 {
		log.Println("Expired riot token")
		return rankedinfo
	}

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&rankedArray)
		if err != nil {
			log.Fatalf("error decoding response into rankedArray \n%v", err)
		}

		for _, v := range rankedArray {
			if v.QueueType == "RANKED_SOLO_5x5" {
				rankedinfo = v
				notRanked = false
			}
		}
	}

	if notRanked == true {
		rankedinfo = Models.LeagueRanked{
			Tier: "",
			Rank: "",
		}
	}

	if rankedinfo.Tier == "" {
		rankedinfo.Tier = "UNRANKED"
		rankedinfo.Rank = ""
	}

	//change the all caps tier to be a title
	rankedinfo.Tier = strings.ToTitle(rankedinfo.Tier)

	return rankedinfo
}
