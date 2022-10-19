package riot_api

import (
	"encoding/json"
	"log"
	"main/Models"
	"net/http"
	"strconv"
)

func Idtoname(id int) string {
	var Champions Models.Champions

	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("http://ddragon.leagueoflegends.com/cdn/12.20.1/data/en_US/champion.json"), nil)
	response, _ := client.Do(request)

	if response.StatusCode != 404 {
		err := json.NewDecoder(response.Body).Decode(&Champions)
		if err != nil {
			log.Fatalf("error decoding response into champs array champidtoname \n%v", err)
		}
		for _, champion := range Champions.Data {
			if champion.Key == strconv.Itoa(id) {
				return champion.Name
			}
		}
	}
	return "Error"
}
