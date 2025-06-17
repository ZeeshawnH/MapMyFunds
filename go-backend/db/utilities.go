package db

import (
	"context"
	"log"
	"os"
	"time"

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

	// Ping to verify
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MongoDB!")

	return client, nil
}
