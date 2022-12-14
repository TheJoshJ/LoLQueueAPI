package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"main/models"
	"net/http"
	"os"
	"strconv"
)

func Idtoname(id int) string {
	var Champions models.Champions

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

func GetChampionMastery(id string, server string) []models.ChampionMasteryResp {
	var champions []models.ChampionMasteryResp
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

func GetBySummonerName(user string, server string) (models.RiotBySummonerName, error) {

	var bySummonerName models.RiotBySummonerName

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

		return bySummonerName, nil
	}
	err := errors.New("user not found")
	return bySummonerName, err
}

func MatchInfo(matchid string, srv string) models.MatchData {
	var match models.MatchData
	var region string
	if srv == "NA" || srv == "BR" || srv == "LAN" || srv == "LAS" {
		region = "AMERICAS"
	}
	if srv == "KR" || srv == "JP" {
		region = "ASIA"
	}
	if srv == "EUNE" || srv == "EUW" || srv == "TR" || srv == "RU" {
		region = "EUROPE"
	}
	if srv == "OCE" {
		region = "SEA"
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

func MatchList(id string, srv string, limit int) []string {
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

	count := strconv.Itoa(limit)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", ("https://" + region + ".api.riotgames.com/lol/match/v5/matches/by-puuid/" + id + "/ids?count=" + count), nil)
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

//func MatchListByTime(id string, srv string, start time.Time, end time.Time) []string {
//}

//func convertRegion(srv string) string {
//	var region string
//	if srv == "NA" || srv == "BR" || srv == "LAN" || srv == "LAS" {
//		region = "americas"
//	}
//	if srv == "KR" || srv == "JP" {
//		region = "asia"
//	}
//	if srv == "EUNE" || srv == "EUW" || srv == "TR" || srv == "RU" {
//		region = "europe"
//	}
//	if srv == "OCE" {
//		region = "sea"
//	}
//	return region
//}

//func convertTime()

func GetRankedInfo(Id string, server string) models.LeagueRanked {

	var rankedArray []models.LeagueRanked
	var rankedinfo models.LeagueRanked
	var notRanked = true

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
		rankedinfo = models.LeagueRanked{
			Tier: "",
			Rank: "",
		}
	}

	if rankedinfo.Tier == "" {
		rankedinfo.Tier = "UNRANKED"
		rankedinfo.Rank = ""
	}
	return rankedinfo
}

//func getLolVersion() string {
//	var versions []string
//
//	client := &http.Client{}
//	request, _ := http.NewRequest("GET", "https://ddragon.leagueoflegends.com/api/versions.json", nil)
//	request.Header.Set("X-Riot-Token", os.Getenv("RIOTKEY"))
//	response, _ := client.Do(request)
//
//	if response.StatusCode != 404 {
//		err := json.NewDecoder(response.Body).Decode(&versions)
//		if err != nil {
//			log.Fatalf("error decoding response into versions \n%v", err)
//		}
//	}
//	return versions[0]
//}

func translateServerName(server string) string {

	var serverString string

	switch server {
	case "BR":
		serverString = "BR1"
	case "EUNE":
		serverString = "EUN1"
	case "EUW":
		serverString = "EUW1"
	case "JP":
		serverString = "JP1"
	case "KR":
		serverString = "KR1"
	case "LAN":
		serverString = "LA1"
	case "LAS":
		serverString = "LA2"
	case "NA":
		serverString = "NA1"
	case "OCE":
		serverString = "OCE1"
	case "RU":
		serverString = "RU1"
	case "TR":
		serverString = "TR1"
	default:
		serverString = "NA1"
	}

	return serverString
}
