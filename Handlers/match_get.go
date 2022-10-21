package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/Models"
	"main/riot_api"
	"net/http"
)

func MatchGet(w http.ResponseWriter, r *http.Request) {
	var matchList []string
	matchesData := make([]Models.MatchData, 10)
	matchDataReturn := make([]Models.Participants, 10)

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var userSearch Models.UserLookup

	//process the information sent via the PostForm request
	vars := mux.Vars(r)
	userSearch.Username = vars["usr"]
	userSearch.Server = vars["srv"]

	if userSearch.Username == "" || userSearch.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bySummonerName := riot_api.GetBySummonerName(userSearch.Username, userSearch.Server)
	matchList = riot_api.MatchList(bySummonerName.Puuid, userSearch.Server)

	for i, matchid := range matchList {
		matchesData[i] = riot_api.MatchInfo(matchid, userSearch.Server)
	}

	for mi, mdata := range matchesData {
		for pi, participant := range mdata.Info.Participants {
			if participant.Puuid == bySummonerName.Puuid {
				matchDataReturn[mi].GameMode = matchesData[pi].Info.GameMode
				log.Println(mi)
				log.Println(matchDataReturn[mi])
			}
			matchDataReturn[mi].GameID = matchesData[mi].Metadata.MatchId
			matchDataReturn[mi] = participant
		}
	}

	for i := range matchDataReturn {
		log.Println(i)
		log.Println(matchDataReturn[i])
	}

	log.Printf("%#v", matchDataReturn)

	reply, err := json.Marshal(matchDataReturn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
