package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetContributionsByStateAndYear(client *mongo.Client, year int, state string) ([]Contribution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("contributions")

	filter := bson.M{
		"contributor_state": state,
		"election_year":     year,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query contributions: %w", err)
	}
	defer cursor.Close(ctx)

	var rtn []Contribution
	if err := cursor.All(ctx, &rtn); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return rtn, nil
}
