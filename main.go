package main

import (
	"log"
	"time"

	"redis-demo/internals/db"
	"redis-demo/server"
)

func main() {
	database := db.NewDB()
	database.StartCleanup(1 * time.Second)

	srv := server.NewServer(":6379", database)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
