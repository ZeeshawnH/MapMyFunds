package db

import (
	"backend/openfec"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PopulateContributionsCollection(client *mongo.Client, year int) error {
	data, err := openfec.GetContributions(year)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("contributions")

	for _, c := range data {
		dbContribution := Contribution{
			CandidateID:      c.CandidateID,
			CandidateName:    c.CandidateLastName,
			CandidateParty:   c.CandidatePartyAffiliation,
			ContributorState: c.ContributorState,
			ElectionYear:     c.ElectionYear,
			NetReceipts:      c.NetReceipts,
		}

		filter := bson.M{
			"candidate_id":      dbContribution.CandidateID,
			"contributor_state": dbContribution.ContributorState,
			"election_year":     dbContribution.ElectionYear,
		}

		update := bson.M{
			"$set": dbContribution,
		}

		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Printf("failed to upsert %v: %v", dbContribution, err)
		}
	}

	log.Println("Successfully updated database!")
	return nil
}

func PopulateCandidatesCollection(client *mongo.Client, year int) error {
	data, err := openfec.GetCandidateData(year)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("candidates")

	var operations []mongo.WriteModel

	for _, c := range data {
		filter := bson.M{"candidate_id": c.CandidateID}
		update := bson.M{"$set": c}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	if len(operations) > 0 {
		result, err := collection.BulkWrite(ctx, operations)
		if err != nil {
			return fmt.Errorf("failed to execute bulk write: %w", err)
		}

		log.Printf("Successfully updated candidates: %d upserted, %d modified",
			result.UpsertedCount, result.ModifiedCount)
	} else {
		log.Println("No candidate data to update")
	}

	log.Println("Successfully updated Candidate info in database!")
	return nil
}
