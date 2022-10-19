package riot_api

import (
	"encoding/json"
	"io"
	"log"
	"main/Models"
	"net/http"
	"os"
)

func GetBySummonerName(user string, server string) Models.RiotBySummonerName {

	var bySummonerName Models.RiotBySummonerName

	serverString := translateServerName(server)

	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("https://" + serverString + ".api.riotgames.com/lol/summoner/v4/summoners/by-name/" + user), nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	if response.StatusCode != 404 {

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &bySummonerName)

		//log.Println(bySummonerName)

		return bySummonerName
	}
	return bySummonerName
}
