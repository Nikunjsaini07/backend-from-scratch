package main

import (
	"log"

	"redis-demo/internals/server"
	"redis-demo/internals/db"
)

func main() {
	database := db.NewDB()
	srv := server.NewServer(":6379" , database)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}