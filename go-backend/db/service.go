package db

import (
	"backend/types"
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

func GetAllContributionsWithCandidates(client *mongo.Client, year int) (*ContributionsWithCandidatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all contributions for the year
	collection := client.Database("election_data").Collection("contributions")
	filter := bson.M{"election_year": year}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query contributions: %w", err)
	}
	defer cursor.Close(ctx)

	var allContributions []Contribution
	if err := cursor.All(ctx, &allContributions); err != nil {
		return nil, fmt.Errorf("failed to decode contributions: %w", err)
	}

	// Group by state
	contributionsByState := make(map[string][]Contribution)
	candidateIDSet := make(map[string]bool)

	for _, contrib := range allContributions {
		contributionsByState[contrib.ContributorState] = append(
			contributionsByState[contrib.ContributorState],
			contrib,
		)
		candidateIDSet[contrib.CandidateID] = true
	}

	// Unique candidate IDs
	var candidateIDs []string
	for id := range candidateIDSet {
		candidateIDs = append(candidateIDs, id)
	}

	candidates, err := getCandidatesByIDs(client, candidateIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates: %w", err)
	}

	return &ContributionsWithCandidatesResponse{
		Contributions: contributionsByState,
		Candidates:    candidates,
	}, nil
}

func getCandidatesByIDs(client *mongo.Client, candidateIDs []string) ([]types.Candidate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("election_data").Collection("candidates")
	filter := bson.M{"candidate_id": bson.M{"$in": candidateIDs}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var candidates []types.Candidate
	if err := cursor.All(ctx, &candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}
