package main

import (
	"log"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/api"
	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/database/mongodb/models"
)


func main() {
	config, err := models.ConfigureDB()
	if err!=nil{
		log.Fatal(err)
	}
	app, err := api.NewApp(config.Database.ConnectionURI, config.Database.DBName)
	if err!=nil{
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	
}