package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	App      App
	Mongo        Mongo
	Fixer        Fixer
}

type App struct {
	Name string
	Port string
}

type Mongo struct {
	User       string
	Password   string
	Urls       []string
	Port       string
	ReplicaSet string
	Database   string
}

type Fixer struct {
	Url    string
	ApiKey string
}

var Config *Configuration

func init() {
	loadConfig()

	Config = &Configuration{
		App: App{
			Name: viper.GetString("app.name"),
			Port: viper.GetString("app.port"),
		},
		Mongo: Mongo{
			User:       viper.GetString("mongodb.username"),
			Password:   viper.GetString("mongodb.password"),
			Urls:       viper.GetStringSlice("mongodb.urls"),
			Port:       viper.GetString("mongodb.port"),
			ReplicaSet: viper.GetString("mongodb.replica_set"),
			Database:   viper.GetString("mongodb.database"),
		},
		Fixer: Fixer{
			Url:    viper.GetString("webservices.fixer.url"),
			ApiKey: viper.GetString("webservices.fixer.api_key"),
		},
	}
}

func loadConfig() {
	appEnv := os.Getenv("APP_ENV")

	if IsTestingAppEnv(appEnv) {
		return
	}

	if !isValidAppEnv(appEnv) {
		fmt.Printf("Failed to start, since APP_ENV = %s is not registered\n", appEnv)
		os.Exit(1)
	}

	configFileName := appEnv + "-config"
	viper.SetConfigName(configFileName) // name of config file (without extension)
	viper.SetConfigType("yaml")         // path to look for the config file in
	viper.AddConfigPath("./config")     // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func isValidAppEnv(appEnv string) bool {
	return appEnv == "localhost" || appEnv == "development" || appEnv == "staging" || appEnv == "production"
}

func IsTestingAppEnv(appEnv string) bool {
	return appEnv == "testing"
}
