package Initializers

import (
	"log"
	"main/Handlers"
)

func (c *Connect) initializeRoutes() {
	c.router.HandleFunc("/ping", Handlers.Ping).Methods("GET")
	c.router.HandleFunc("/lookup/{srv}/{usr}", Handlers.ProfileLookup).Methods("GET")
	c.router.HandleFunc("/match/{srv}/{usr}", Handlers.MatchGet).Methods("GET")

	//c.router.HandleFunc("/user/{id}", Handlers.ViewUser).Methods("GET")
	c.router.HandleFunc("/user", Handlers.CreateUser).Methods("POST")

	log.Println("Loaded Routes")
}
