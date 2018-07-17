package conversion

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func connect(connectionString string, dbName string) (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), connectionString, nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func removeAllItemsInCollection(db *mongo.Database, versionTag string, collectionName string) {
	coll := db.Collection(collectionName)
	result, err := coll.DeleteMany(context.Background(),
		bson.NewDocument(bson.EC.String("version", versionTag)))

	fmt.Printf("Deleted %v page(s) for version %s \n", result.DeletedCount, versionTag)

	if err != nil {
		fmt.Println("Error removing existing items for version " + versionTag + ": " + err.Error())
	}
}

func saveHtmlFileToDb(db *mongo.Database, fileName string, html string, versionTag string) {

	coll := db.Collection("pages")

	_, err := coll.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("name", fileName),
			bson.EC.String("version", versionTag),
			bson.EC.String("html", html),
			bson.EC.DateTime("lastUpdated", getCurrentTimestampUTC()),
		))

	if err != nil {
		fmt.Println("Error inserting file " + fileName + ": " + err.Error())
	}
}

func saveHintToDb(db *mongo.Database, hintKey string, hintText string, versionTag string) {

	coll := db.Collection("hints")

	_, err := coll.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("key", hintKey),
			bson.EC.String("version", versionTag),
			bson.EC.String("hint", hintText),
			bson.EC.DateTime("lastUpdated", getCurrentTimestampUTC()),
		))

	if err != nil {
		fmt.Println("Error inserting hint " + hintKey + ": " + err.Error())
	}
}

func saveHintsToDb(db *mongo.Database, hints []Hint, versionTag string) {

	docs := make([]interface{}, 0)

	for _, hint := range hints {
		doc := bson.NewDocument(
			bson.EC.String("key", hint.Key),
			bson.EC.String("version", versionTag),
			bson.EC.String("hint", hint.Text),
			bson.EC.DateTime("lastUpdated", getCurrentTimestampUTC()),
		)
		docs = append(docs, doc)
	}

	coll := db.Collection("hints")

	_, err := coll.InsertMany(
		context.Background(), docs)

	if err != nil {
		fmt.Println("Error inserting hints: " + err.Error())
	}
}

func getCurrentTimestampUTC() int64 {
	return time.Now().UTC().UnixNano() / 1e6 // pass milliseconds
}
