package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var (
	billsCollection       *mongo.Collection
	membersCollection     *mongo.Collection
	cellsCollection       *mongo.Collection
	policyAreasCollection *mongo.Collection
	subjectsCollection    *mongo.Collection
)

// Connect establishes the database connection
func Connect(uri string) (err error) {
	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err.Error())
	}
	billsCollection = client.Database("cosign").Collection("bills")
	membersCollection = client.Database("cosign").Collection("members")
	cellsCollection = client.Database("cosign").Collection("cells")
	policyAreasCollection = client.Database("cosign").Collection("policyAreas")
	subjectsCollection = client.Database("cosign").Collection("subjects")

	fmt.Println("Connected to Mongo...")

	return
}

// Disconnect tears down the database connection
func Disconnect() {
	client.Disconnect(context.Background())
	fmt.Println("Disconnected from Mongo...")
}
