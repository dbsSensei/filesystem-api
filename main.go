package main

import (
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/database"
	"github.com/dbsSensei/filesystem-api/server"
)

func main() {
	// Load config
	c, err := config.LoadConfig(".")
	if err != nil {
		panic("Failed to load config!")
	}

	// Initialize database
	db, err := database.Init(c)
	if err != nil {
		panic(err)
	}

	//	Initialize server
	err = server.Init(c, db)
	if err != nil {
		panic(err)
	}
}
