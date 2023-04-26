package database

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSubData struct {
	User_id     int64          // user chat id
	Location    GeoCoordinates //geo location coordinates of chosen city
	Language    string         // user's language
	MailingTime time.Time      // time user chosen for daily weather forecast mailing
	Code        int            //subscription step
}
type GeoCoordinates struct {
	Latitude  float32 // latitude of chosen city
	Longitude float32 // longitude of chosen city
}

// Connect to the Mongo database
func ConnectDataBase() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017/")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {

		log.Fatalf("Cannot connect to the MongoDB. Error - %s", err)
		return nil, err
	}
	return client, nil
}

// Check if there is connection with the MongoDB
func CheckDatabaseConnection(client *mongo.Client) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Cannot connect to the database. Error - %s", err)
		return false, err
	}

	log.Info("Connected to MongoDB!")
	return true, nil
}

// subscribe user for daily forecast mailing
func SubscribeUser(collection *mongo.Collection, user UserSubData) (bool, error) {
	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Errorf("Cannot subscribe the user. Error - %s", err)
		return false, err
	}
	log.Info(res)
	return true, nil
}

// check if the user already subscribed
func IsUserAlreadySubscribed(collection *mongo.Collection, userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Stop user's subscription. Delete from the Mongo database
func DeleteUsersSubscription(DBcollection *mongo.Collection, userID int64) error {
	filter := bson.M{"user_id": userID}
	result, err := DBcollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Errorf("Cannot delete subscription. Error - %s", err)
		return err
	}
	if result.DeletedCount >= 1 {
		log.Infof("User %v deleted from the database. Deleted %v document", userID, result.DeletedCount)
		return nil
	} else {
		log.Infof("User %v is not deleted from the database. Deleted %v document", userID, result.DeletedCount)
		return errors.New("0 deleted documents")
	}
}

// check whether inputted by user mailing time is really time
func IsitTime(data string) (bool, time.Time) {
	userTime, err := time.Parse("15:04", data)
	if err != nil {
		log.Errorf("Cannot parse inputted time. Error - %s", err)
		return false, time.Now()
	}
	return true, userTime
}
