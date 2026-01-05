package utils

import (
	"os"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseSetting
	Server   ServerSettings
	App      Application
}

type DatabaseSetting struct {
	Url        string
	DbName     string
	Collection string
}

type ServerSettings struct {
	Port string
}

type Application struct {
	Name string
}


func ReadConfig() Config{
	config := readConfigFile()

	mongoUri := preferEnv(os.Getenv("CHATTY_MONGODB_URL"), config.Database.Url)
	dbName := preferEnv(os.Getenv("CHATTY_DB_NAME"), config.Database.DbName)
	collection := preferEnv(os.Getenv("CHATTY_COLLECTION"), config.Database.Collection)
	port := preferEnv(os.Getenv("CHATTY_SERVER_PORT"), config.Server.Port)
	appName := preferEnv(os.Getenv("CHATTY_APP_NAME"), config.App.Name)

	config = Config{
		App:      Application{Name: appName},
		Database: DatabaseSetting{Url: mongoUri, DbName: dbName, Collection: collection},
		Server:   ServerSettings{Port: port},
	}

	log.Printf("Config with environment %v", config)

	return config

}

func readConfigFile() Config {
	//Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./..")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Printf("Config with variables %v", config)

	return config
}

func preferEnv(env string, cfg string) string {
	if env != "" {
		return env
	} else {
		return cfg
	}
}
