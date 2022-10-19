package riot_api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func getLolVersion() string {
	var versions []string

	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://ddragon.leagueoflegends.com/api/versions.json", nil)
	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
	response, _ := client.Do(request)

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&versions)
		if err != nil {
			log.Fatalf("error decoding response into versions \n%v", err)
		}
	}
	return versions[0]
}
