package main

import (
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/database"
	"github.com/dbsSensei/filesystem-api/server"
	"github.com/dbsSensei/filesystem-api/service"
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

	// Initialize service
	s := service.Init(db)

	//	Initialize server
	err = server.Init(c, db, s)
	if err != nil {
		panic(err)
	}
}
