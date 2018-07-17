package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type getFileHtml func(string) string

func connect(connectionString string, dbName string) (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), connectionString, nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func removeAllFilesForVersion(db *mongo.Database, versionTag string) {
	coll := db.Collection("pages")
	result, err := coll.DeleteMany(context.Background(),
		bson.NewDocument(bson.EC.String("version", versionTag)))

	fmt.Printf("Deleted %v page(s) for version %s \n", result.DeletedCount, versionTag)

	if err != nil {
		fmt.Println("Error removing existing items for version " + versionTag + ": " + err.Error())
	}
}

func saveHtmlFileToDb(db *mongo.Database, fileName string, html string, versionTag string) {

	coll := db.Collection("pages")

	utcTimestamp := time.Now().UTC().UnixNano() / 1e6 // pass milliseconds

	_, err := coll.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("name", fileName),
			bson.EC.String("version", versionTag),
			bson.EC.String("html", html),
			bson.EC.DateTime("lastUpdated", utcTimestamp),
		))

	if err != nil {
		fmt.Println("Error inserting file " + fileName + ": " + err.Error())
	}
}
