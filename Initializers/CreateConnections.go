package Initializers

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Connect struct {
	router *mux.Router
	db     *sql.DB
}

func (c *Connect) CreatePostgresConnect() {

	db, err := sql.Open("postgresql", "postgres://$PGUSER:$PGPASS@$PGHOST/railway")
	if err != nil {
		log.Fatal(err)
	}
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
	c.initializeRoutes()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
