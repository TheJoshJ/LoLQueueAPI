package Initializers

import (
	"log"
	"main/Handlers"
)

func (c *Connect) initializeRoutes() {
	c.router.HandleFunc("/ping", Handlers.Ping).Methods("GET")

	//c.router.HandleFunc("/user/{uid}", Handlers.ViewUser).Methods("GET")
	c.router.HandleFunc("/user/{id}", Handlers.CreateUser).Methods("POST")

	log.Println("Loaded Routes")
}
