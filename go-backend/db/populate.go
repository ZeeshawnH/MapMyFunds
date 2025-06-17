package db

import (
	"backend/openfec"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PopulateDatabase(client *mongo.Client) error {
	data, err := openfec.GetContributions(2024)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("contributions")

	for _, c := range data {
		filter := bson.M{
			"candidate_id":      c.CandidateID,
			"contributor_state": c.ContributorState,
			"election_year":     c.ElectionYear,
		}

		update := bson.M{
			"$set": c,
		}

		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Printf("failed to upsert %v: %v", c, err)
		}
	}

	log.Println("Successfully updated database!")
	return nil
}
