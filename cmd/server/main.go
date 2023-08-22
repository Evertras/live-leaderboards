package main

import (
	"log"

	"github.com/Evertras/live-leaderboards/pkg/server"
)

func main() {
	s := server.New(":8037")

	log.Fatal(s.ListenAndServe())
}
