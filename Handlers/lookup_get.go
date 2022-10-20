package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"main/Models"
	"main/riot_api"
	"net/http"
)

func ProfileLookup(w http.ResponseWriter, r *http.Request) {
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
	rankedinfo := riot_api.GetRankedInfo(bySummonerName.Id, userSearch.Server)
	championMastery := riot_api.GetChampionMastery(bySummonerName.Id, userSearch.Server)

	if bySummonerName.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	lookupResponse := Models.LookupResponse{
		Username:      rankedinfo.SummonerName,
		Rank:          rankedinfo.Rank,
		Tier:          rankedinfo.Tier,
		Champions:     championMastery,
		Level:         bySummonerName.SummonerLevel,
		ProfileIconId: bySummonerName.ProfileIconId,
		Wins:          rankedinfo.Wins,
		Losses:        rankedinfo.Losses,
	}

	reply, err := json.Marshal(lookupResponse)
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
