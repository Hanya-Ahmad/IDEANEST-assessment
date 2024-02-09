package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.1.4")
	fmt.Println("ClientOptopm TYPE:", reflect.TypeOf(clientOptions))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		os.Exit(1)
	}
	col := client.Database("TestDB").Collection("First_Collection")
	fmt.Println("Collection Type: ", reflect.TypeOf(col))

}