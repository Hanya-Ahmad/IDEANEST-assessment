package models

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

// DBConfiguration struct to hold configuration settings
type DBConfiguration struct {
	Database struct {
		ConnectionURI string `yaml:"connectionURI"`
		DBName        string `yaml:"dbName"`
	} `yaml:"database"`
}

// DBClient is used query the database
type DBClient struct {
	Client *mongo.Database
}

// ConfigureDB loads the configuration settings for the database
func ConfigureDB(file *os.File) (DBConfiguration, error){
	var config DBConfiguration
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
				return DBConfiguration{}, err
	}
	return config, nil
}
