package api

import (
	"context"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/database/mongodb/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App initializes the entire app
type App struct{
	DB *models.DBClient
	router *gin.Engine
}

// NewApp constructs our app object
func NewApp(connectionURI string, dbName string)(App, error){
	db, err := ConnectToDB(connectionURI,dbName)
	if err!=nil{
		return App{}, err
	}
	return App{DB:&models.DBClient{Client: db}, router:gin.Default()}, nil
}

// Run runs the server by setting the router and calling the internal registerRoutes method
func (app *App) Run() error {
	app.registerRoutes()
	return app.router.Run(":8080")
}

// ConnectToDB connects to a database
func ConnectToDB(connectionURI string, dbName string) (*mongo.Database, error){
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}
