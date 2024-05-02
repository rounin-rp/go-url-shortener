package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientWrapper struct {
	mongoClient *mongo.Client
}

var (
	ctx          = context.Background()
	mongoWrapper = &MongoClientWrapper{}
	dbName       = ""
)

func InitializeMongo(dburi, dbname string) (*MongoClientWrapper, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("failed to connect to mongodb | error = ", err.Error())
		return nil, err
	}
	mongoWrapper.mongoClient = client
	dbName = dbname
	log.Printf("mongodb successfully connected at - %v", dburi)
	return mongoWrapper, nil
}

func SaveUrlMapping(shortUrl, longUrl, userId string) {
	start := time.Now()
	collection := mongoWrapper.mongoClient.Database(dbName).Collection("urls")
	// check if the short url already exist in the Database
	var searchResult bson.M
	err := collection.FindOne(ctx, bson.D{{"shortUrl", shortUrl}}).Decode(&searchResult)
	if err == nil {
		return
	}
	objectId := primitive.NewObjectID().Hex()
	res, err := collection.InsertOne(ctx, bson.M{"_id": objectId, "shortUrl": shortUrl, "longUrl": longUrl, "usreId": userId})
	if err != nil {
		log.Panic("Failed to save the url in db | shortUrl = ", shortUrl)
	} else {
		log.Printf("Short Url saved successfully in db | shortUrl = %v - Inserted Id = %v\n", shortUrl, res.InsertedID)
	}
	elapsed := time.Since(start)
	log.Printf("Total time took = %v \n", elapsed)
}

func RetrieveInitialUrl(shortUrl string) (string, error) {
	collection := mongoWrapper.mongoClient.Database(dbName).Collection("urls")
	var result bson.M
	err := collection.FindOne(ctx, bson.D{{"shortUrl", shortUrl}}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result["longUrl"].(string), nil
}
