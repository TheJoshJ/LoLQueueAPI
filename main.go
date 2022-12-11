package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	api "main/api/riot"
	_ "main/docs"
	"main/models"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"
)

// @title LoLQueue API
// @version 1.0
// @description This is the documentation for the LoLQueue Api Service
// @termsOfService There are no terms of service. We accept no responsibility for your ignorance.
// @host api.LoLQueue.com

type ProfileHandler struct {
	router *mux.Router
	db     *gorm.DB
}

func main() {
	godotenv.Load(".env")
	c := ProfileHandler{}
	c.CreatePostgresConnect()
	c.MuxInit()
}
func (c *ProfileHandler) CreatePostgresConnect() {
	var dsn = "postgresql://" + os.Getenv("PGUSER") + os.Getenv("PGPASS") + "@" + os.Getenv("PGHOST") + ":" + os.Getenv("PGPORT") + "/railway"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sucessfully created the PostgreSQL server!")
	c.db = db
}
func (c *ProfileHandler) MuxInit() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	//create the router
	c.router = mux.NewRouter()
	log.Println("Router Created")

	//define the server
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      c.router,
	}

	//run the server as a go routine, so we don't block any other processes
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	//load the Handlers
	c.AddRoutes()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	log.Println("Service is running!")
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("shutting down")
	os.Exit(0)
}
func (c *ProfileHandler) AddRoutes() {
	c.router.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL(os.Getenv("API_URL")+"/docs/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	//utility
	c.router.HandleFunc("/v1/ping", Ping).Methods("GET")

	//riot
	c.router.HandleFunc("/v1/lookup/{srv}/{usr}", c.ProfileLookup).Methods("GET")
	c.router.HandleFunc("/v1/match/{srv}/{usr}", c.GetRecentMatches).Methods("GET")

	//profiles
	c.router.HandleFunc("/v1/user/{id}", c.LookupUser).Methods("GET")
	c.router.HandleFunc("/v1/user", c.CreateUser).Methods("POST")
	c.router.HandleFunc("/v1/user", c.UpdateUser).Methods("PUT")
	c.router.HandleFunc("/v1/user/{id}", c.DeleteUser).Methods("Delete")

	//services
	c.router.HandleFunc("/v1/leaderboard", c.GetLeaderboard).Methods("GET")

	log.Println("Loaded Routes")
}

// Ping godoc
// @Summary      Pings the API service to ensure that it is active
// @Description  Ping the API service
// @Tags         utility
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      404
// @Router      /v1/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"data": "Pong"})
	if err != nil {
		log.Println("unable to encode response from Ping handler")
	}
}

// ProfileLookup godoc
// @Summary      Show an account
// @Description  Gets the users account information by their Username and Server
// @Tags         Riot
// @Accept       json
// @Produce      json
// @Param        srv   path      string  true  "Riot Server"
// @Param        usr   path      string  true  "Username"
// @Success      200   {object}  models.LookupResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /v1/lookup/{srv}/{usr} [get]
func (c *ProfileHandler) ProfileLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var userSearch models.UserLookup

	//process the information from the URL
	vars := mux.Vars(r)
	userSearch.Username = vars["usr"]
	userSearch.Server = vars["srv"]

	if userSearch.Username == "" || userSearch.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bySummonerName, err := api.GetBySummonerName(userSearch.Username, userSearch.Server)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rankedinfo := api.GetRankedInfo(bySummonerName.Id, userSearch.Server)
	championMastery := api.GetChampionMastery(bySummonerName.Id, userSearch.Server)

	lookupResponse := models.LookupResponse{
		Username:      bySummonerName.Name,
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

// GetRecentMatches godoc
// @Summary      Show recent matches
// @Description  Show the past 10 matches
// @Tags         Riot
// @Accept       json
// @Produce      json
// @Param        srv   path      string  true  "Riot Server"
// @Param        usr   path      string  true  "Username"
// @Param        usr   query      string  true  "count"
// @Success      200  {array}    models.MatchDataResp
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /v1/match/{srv}/{usr} [get]
func (c *ProfileHandler) GetRecentMatches(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("count")
	log.Println(limit)
	if limit == "" {
		// id.asc is the default sort query
		limit = "10"
	}

	limitInt, _ := strconv.Atoi(limit)
	if limitInt > 50 {
		limitInt = 20
	}

	var matchList []string
	matchesData := make([]models.MatchData, 0)
	matchDataReturn := make([]models.Participants, limitInt)
	matchDataResp := make([]models.MatchDataResp, limitInt)

	w.Header().Set("Content-Type", "application/json")

	//create a new instance of a struct for us to process
	var userSearch models.UserLookup

	//process the information sent via the PostForm request
	vars := mux.Vars(r)
	userSearch.Username = vars["usr"]
	userSearch.Server = vars["srv"]

	if userSearch.Username == "" || userSearch.Server == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bySummonerName, err := api.GetBySummonerName(userSearch.Username, userSearch.Server)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	matchList = api.MatchList(bySummonerName.Puuid, userSearch.Server, limitInt)

	ch := make(chan models.MatchData)
	for _, matchid := range matchList {
		go func(matchid string) {
			match := api.MatchInfo(matchid, userSearch.Server)
			ch <- match
		}(matchid)
	}
	for i := 0; i < limitInt; i++ {
		match := <-ch
		matchesData = append(matchesData, match)
	}
	close(ch)
	for idx, mdata := range matchesData {
		log.Println(idx)
		for _, participant := range mdata.Info.Participants {
			if participant.Puuid == bySummonerName.Puuid {
				matchDataReturn[idx] = participant
			}
			matchDataResp[idx].GameID = strings.ReplaceAll(matchesData[idx].Metadata.MatchId, "NA1_", "")
			matchDataResp[idx].GameMode = matchesData[idx].Info.GameMode
			matchDataResp[idx].Assists = participant.Assists
			matchDataResp[idx].Deaths = participant.Deaths
			matchDataResp[idx].Kills = participant.Kills
			matchDataResp[idx].Win = participant.Win
			matchDataResp[idx].ChampionName = participant.ChampionName
		}
	}

	sort.Slice(matchDataResp, func(i, j int) bool {
		return matchDataResp[i].GameID > matchDataResp[j].GameID
	})

	reply, err := json.Marshal(matchDataResp)
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

// GetLeaderboard godoc
// @Summary      Get leaderboard data
// @Description  Get the leaderboards for a specifc discord server ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.Leaderboard
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /v1/leaderboard [get]
func (c *ProfileHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := os.ReadFile("./models/MOCK_LEADERBOARD.json")
	if err != nil {
		log.Fatalf("Failed to read leaderboard - %s.", err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CreateUser godoc
// @Summary      Create an account
// @Description  Creates and stores the users data to be used when executing commands/api calls.
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Success      201
// @Failure      400
// @Failure      404
// @Failure      409
// @Failure      500
// @Router      /v1/user [post]
func (c *ProfileHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var server models.Server
	var discordUser models.Discord_user
	var serverUser models.Server_user
	var newUser models.UserPost
	var riotUser models.Riot_user
	var duru models.Discord_user_riot_user
	var q *gorm.DB

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//does the server already exist in server?
	q = c.db.Table("server").First(&server, "id = ?", newUser.DiscordServerID)
	//if it doesn't, add it to the table
	if q.RowsAffected == 0 {
		c.db.Table("server").Create(&models.Server{
			Id:   newUser.DiscordServerID,
			Name: newUser.DiscordServerName,
		})
	}

	//does the user already exist alongside this server?
	q = c.db.Table("server_user").Where(&models.Server_user{Discord_id: newUser.DiscordID, Server_id: newUser.DiscordServerID}).First(&serverUser)
	if q.RowsAffected == 0 {
		//does the user exist in discord_user
		q = c.db.Table("discord_user").First(&discordUser, "id = ?", newUser.DiscordID)
		if q.RowsAffected == 0 {
			//if not, add data to discord_user
			c.db.Table("discord_user").Create(&models.Discord_user{
				Id:       newUser.DiscordID,
				Username: newUser.DiscordUsername,
			})
		}
		//add the data to server_user
		c.db.Table("server_user").Create(&models.Server_user{
			Server_id:  newUser.DiscordServerID,
			Discord_id: newUser.DiscordID,
		})
	} else {
		//user already exists alongside this server - return
		w.WriteHeader(http.StatusConflict)
		return
	}
	//does the username already exist within the riot_user table?
	q = c.db.Table("riot_user").Where(&models.Riot_user{Username: newUser.RiotUsername}).First(&riotUser)
	if q.RowsAffected != 0 {
		//username does exist within the riot_user table. Does the puuid exist in the duru table?
		q = c.db.Table("discord_user_riot_user").Where(&models.Discord_user_riot_user{Puuid: riotUser.Puuid}).First(&duru)
		if q.RowsAffected != 0 {
			//puuid does exist within the duru table - return
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	//since the username did not exist in the riot_user table get the information from Riot.
	rr, err := api.GetBySummonerName(newUser.RiotUsername, newUser.RiotServer)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//does the puuid already exist in the duru table?
	q = c.db.Table("discord_user_riot_user").Where("puuid = ?", rr.Puuid).Or("discord_id = ?", newUser.DiscordID).First(&duru)
	log.Println(q.RowsAffected)
	if q.RowsAffected != 0 {
		//puuid or discord_id does exist within the duru table - return
		w.WriteHeader(http.StatusConflict)
		return
	}
	//puuid does not exist, is the username valid?
	if err != nil {
		//add to server_user map and exit - not found
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//since the username is valid, add data to riot_user
	c.db.Table("riot_user").Create(&models.Riot_user{
		Puuid:           rr.Puuid,
		Username:        newUser.RiotUsername,
		Server:          newUser.RiotServer,
		Riot_account_id: rr.AccountId,
	})
	//add the data to duru
	c.db.Table("discord_user_riot_user").Create(&models.Discord_user_riot_user{
		Puuid:      rr.Puuid,
		Discord_id: newUser.DiscordID,
	})
	//exit - created
	w.WriteHeader(http.StatusCreated)
}

// LookupUser godoc
// @Summary      Look up a profile
// @Description  Looks up a users profile to display which servers they are part of in addidion to their stored riot account information.
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}   models.UserLookupResponse
// @Failure      404
// @Failure      500
// @Router       /v1/user/{id} [get]
func (c *ProfileHandler) LookupUser(w http.ResponseWriter, r *http.Request) {
	var q *gorm.DB
	var servers []models.Server_user
	var srv models.Server
	var response models.UserLookupResponse

	//get the discord ID of the user from the request URL
	vars := mux.Vars(r)
	discordID := vars["id"]

	//get all server names associated with the discord ID
	q = c.db.Table("server_user").Where(&models.Server_user{Discord_id: discordID}).Find(&servers)
	if q.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for i := 0; i < len(servers); i++ {
		c.db.Table("server").Where("id = ?", servers[i].Server_id).Find(&srv)
		response.Servers = append(response.Servers, srv.Name)
		srv.Id = ""
	}

	//get the puuid associated with the Discord ID
	var puuid string
	c.db.Table("discord_user_riot_user").Where(&models.Discord_user_riot_user{Discord_id: discordID}).Find(&puuid)

	//get the Riot Username and Server associated with the puuid
	var riotInfo models.Riot_user
	c.db.Table("riot_user").Where(&models.Riot_user{Puuid: puuid}).Find(&riotInfo)
	response.RiotServer = riotInfo.Server
	response.RiotUsername = riotInfo.Username

	w.WriteHeader(http.StatusOK)
	reply, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

// UpdateUser godoc
// @Summary      Update a profile
// @Description  Updates the users profile.
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      409
// @Failure      500
// @Router       /v1/user/{id} [put]
func (c *ProfileHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var q *gorm.DB
	var uUser models.UserPut
	var discordUser models.Discord_user
	var serverUser models.Server_user
	var duru models.Discord_user_riot_user

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&uUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//does the user already exist in discord_user
	q = c.db.Table("discord_user").First(&discordUser, "id = ?", uUser.DiscordID)
	if q.RowsAffected == 0 {
		//if not, exit 404 not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//does the user have a profile with this server?
	q = c.db.Table("server_user").Where(&models.Server_user{Discord_id: uUser.DiscordID, Server_id: uUser.DiscordServerID}).First(&serverUser)
	if q.RowsAffected == 0 {
		//if not, exit 404 not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//does the user have any information in RUDU?
	q = c.db.Table("discord_user_riot_user").Where("discord_id = ?", uUser.DiscordID).First(&duru)
	if q.RowsAffected == 0 {
		//if not, exit 404 status conflict
		w.WriteHeader(http.StatusConflict)
		return
	}

	//if no, is the provided username valid?
	rr, err := api.GetBySummonerName(uUser.RiotUsername, uUser.RiotServer)
	if err != nil {
		//if no, exit 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//if the username is valid, update the riot_user table.
	q = c.db.Table("riot_user").Where("puuid = ?", duru.Puuid).Updates(models.Riot_user{Puuid: rr.Puuid, Username: uUser.RiotUsername, Server: uUser.RiotServer, Riot_account_id: rr.AccountId})
	if q.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

// DeleteUser godoc
// @Summary      Delete a profile
// @Description  Deletes the users profile by supplying the discord id in the path along with the server ID as a URL param.
// @Tags         profiles
// @Accept       json
// @Param        id   path      string  true  "id"
// @Param        discordID    query     string     false  "string serverid"
// @Produce      json
// @Success      200
// @Failure      404
// @Failure      500
// @Router       /v1/user/{id} [delete]
func (c *ProfileHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var q *gorm.DB
	var uDel models.UserDel
	var discordUser models.Discord_user
	var serverUser models.Server_user
	var riotUser models.Riot_user
	var duru models.Discord_user_riot_user

	vars := mux.Vars(r)
	uDel.DiscordID = vars["id"]
	uDel.DiscordServerID = r.URL.Query().Get("serverid")

	//process the request
	w.Header().Set("Content-Type", "application/json")

	if uDel.DiscordServerID == "" {
		//internal service error here because that means that the bot/website didn't specify a server correctly
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//does the user already exist
	q = c.db.Table("discord_user").Where("id = ?", uDel.DiscordID).First(&discordUser)
	log.Println(q.RowsAffected)
	if q.RowsAffected == 0 {
		//if not, exit 404
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//does the user have a profile with this server?
	q = c.db.Table("server_user").Where("server_id = ? AND discord_id = ?", uDel.DiscordServerID, uDel.DiscordID).First(&serverUser)
	if q.RowsAffected == 0 {
		//if not, exit 404
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//get the puuid from duru
	q = c.db.Table("discord_user_riot_user").Where("discord_id = ?", uDel.DiscordID).First(&duru)
	if q.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//delete the row where the puuid matches in riot_user
	q = c.db.Table("riot_user").Where(models.Riot_user{Puuid: duru.Puuid}).Delete(&riotUser)
	if q.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//delete the row where the discord and server id match in server_user
	q = c.db.Table("server_user").Where("server_id = ? AND discord_id = ?", uDel.DiscordServerID, uDel.DiscordID).Delete(&uDel)
	if q.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
