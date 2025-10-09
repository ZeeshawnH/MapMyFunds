package db

import (
	"backend/types"
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
	CandidateName    string  `bson:"candidate_last_name"`
	CandidateParty   string  `bson:"candidate_party_affiliation"`
	ContributorState string  `bson:"contributor_state"`
	ElectionYear     int     `bson:"election_year"`
	NetReceipts      float64 `bson:"net_receipts"`
}

type ContributionsWithCandidatesResponse struct {
	Contributions map[string][]Contribution `json:"contributions"`
	Candidates    []types.Candidate         `json:"candidates"`
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

// CleanDuplicateFields removes old/duplicate fields from all documents in the 'contributions' collection.
func CleanDuplicateFields(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("contributions")

	// List of incorrect/duplicate fields to remove
	unsetFields := bson.M{
		"candidateid":               "",
		"candidatelastname":         "",
		"candidatepartyaffiliation": "",
		"contributorstate":          "",
		"electionyear":              "",
	}

	update := bson.M{"$unset": unsetFields}

	result, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		return err
	}
	log.Printf("Removed duplicate fields from %d documents", result.ModifiedCount)
	return nil
}
