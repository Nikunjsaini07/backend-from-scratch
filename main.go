package main

import (
	"log"

	"redis-demo/internals/server"
)

func main() {
	srv := server.New(":6379")

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}