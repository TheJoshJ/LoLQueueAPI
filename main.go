package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"main/api/handler"
	_ "main/docs"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title LoLQueue API
// @version 1.0
// @description This is the documentation for the LoLQueue Api Service
// @termsOfService There are no terms of service. We accept no responsibility for your ignorance.

// @host api.LoLQueue.com

var DataSource string = "postgresql://" + os.Getenv("PGUSER") + os.Getenv("PGPASS") + "@" + os.Getenv("PGHOST") + ":" + os.Getenv("PGPORT") + "/railway"

type Connect struct {
	router *mux.Router
	db     *sql.DB
}

func main() {
	c := Connect{}
	c.CreatePostgresConnect()
	c.MuxInit()
	godotenv.Load(".env.local")
}
func (c *Connect) CreatePostgresConnect() {

	db, err := sql.Open("postgres", DataSource)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to Postgres DB")
	}

	c.db = db
}
func (c *Connect) MuxInit() {

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
	log.Println("Press Ctrl+C to exit")
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
func (c *Connect) AddRoutes() {
	c.router.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL(os.Getenv("API_URL")), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	c.router.HandleFunc("/ping", api.Ping).Methods("GET")
	c.router.HandleFunc("/lookup/{srv}/{usr}", api.ProfileLookup).Methods("GET")
	c.router.HandleFunc("/match/{srv}/{usr}", api.MatchGet).Methods("GET")

	//c.router.HandleFunc("/user/{id}", api.ViewUser).Methods("GET")
	c.router.HandleFunc("/user", api.CreateUser).Methods("POST")

	log.Println("Loaded Routes")
}
