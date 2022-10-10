package main

import (
	"main/Initializers"
)

func main() {
	c := Initializers.Connect{}
	c.CreatePostgresConnect()
	c.MuxInit()
}
