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

	log.Printf("%#v", bySummonerName)
	log.Printf("%#v", matchList)

	for i, matchid := range matchList {
		matchesData[i] = riot_api.MatchInfo(matchid, userSearch.Server)
	}

	for _, mdata := range matchesData {
		for i, participant := range mdata.Info.Participants {
			if participant.Puuid == bySummonerName.Puuid {
				matchDataReturn[i] = participant
			}
		}
	}

	for i := range matchesData {
		matchDataReturn[i].GameMode = matchesData[i].Info.GameMode
	}

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
