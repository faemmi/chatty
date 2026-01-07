package database

import (
	utils "chatty/utils"
	"log"
	"strings"

	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(config utils.Config) (*mongo.Client, *mongo.Collection, func()) {
// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client

	// Split connection URL by protocol and host to add user-password
	// authentication to the URL
	var url string

	splitted := strings.Split(config.Database.Url, "://")

	if len(splitted) == 2 {
		user := utils.GetEnvWithDefault("CHATTY_DB_USER", "root")
		password := utils.GetEnvWithDefault("CHATTY_DB_PASSWORD", "testing")
		url = splitted[0] + "://" + user + ":" + password + "@" + splitted[1]
	} else {
		log.Print("Database URL does cannot be split by protocol, using as is. This may cause authentication to fail.")
		url = config.Database.Url
	}

	// Create a new client and connect to the server
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	disconnect := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}

	// Send a ping to confirm a successful connection
	var ping bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&ping); err != nil {
		panic(err)
	}
	log.Printf("Pinged your deployment. You successfully connected to MongoDB!")

	coll := client.Database(config.Database.DbName).Collection(config.Database.Collection)

	return client, coll, disconnect
}
