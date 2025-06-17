package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Contribution struct {
	CandidateID      string  `bson:"candidate_id"`
	CandidateName    string  `bson:"candidate_last_name,omitempty"`
	CandidateParty   string  `bson:"candidate_party_affiliation,omitempty"`
	ContributorState string  `bson:"contributor_state"`
	ElectionYear     int     `bson:"election_year"`
	NetReceipts      float64 `bson:"net_receipts"`
}

func ConnectDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MongoDB!")

	return client, nil
}

func CreateIndexes(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "candidate_id", Value: 1},
			{Key: "contributor_state", Value: 1},
			{Key: "election_year", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("candidate_state_year_unique"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Check if error is due to index already existing
		if mongo.IsDuplicateKeyError(err) {
			log.Println("Index already exists on candidate_id, contributor_state, election_year")
			return nil
		}
		return err
	}

	log.Println("Index created on candidate_id, contributor_state, election_year")
	return nil
}
