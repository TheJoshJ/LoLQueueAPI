package riot_api

import (
	"encoding/json"
	"log"
	"main/Models"
	"net/http"
	"os"
)

func GetChampionMastery(id string, server string) []Models.ChampionMastery {
	var champions []Models.ChampionMastery
	serverString := translateServerName(server)

	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("https://" + serverString + ".api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/" + id + "/top?count=6"), nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	if response.StatusCode == 403 {
		log.Println("Expired riot token")
		return champions
	}

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&champions)
		if err != nil {
			log.Fatalf("error decoding response into []champions \n%v", err)
		}

		for i, champ := range champions {
			champions[i].ChampionName = Idtoname(champ.ChampionId)
		}
		return champions
	}
	return champions
}
