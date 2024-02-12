package main

import (
	"log"
	"os"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/api"
	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/database/mongodb/models"
)


func main() {
	file,err := os.Open("./config/database-config.yaml")
	if err!=nil{
		log.Fatal(err)
	}
	defer file.Close()
	config, err := models.ConfigureDB(file)
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