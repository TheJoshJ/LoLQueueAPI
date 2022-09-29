package main

import (
	"log"
	"main/initializers"
	"os"
	"os/signal"
)

func init() {
	initializers.CreatePostgresConnect()
}

func main() {
	initializers.CreateGinConnection()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
