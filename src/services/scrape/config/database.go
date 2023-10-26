package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func InitDB() {
	godotenv.Load("../../../../env/.env")
	var err error
	var cred options.Credential
	cred.Username = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cred.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	ctx := context.Background()
	dataSourceName := os.Getenv("DB_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dataSourceName).SetServerAPIOptions(serverAPI).SetAuth(cred)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		errors.New("Failed to connect to the database")
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		errors.New("Can't ping the database")
	}

	database = client.Database("Nocaine")

	fmt.Println("Successfully connected to the database")
}

func Getdb() *mongo.Database {
	return database
}
