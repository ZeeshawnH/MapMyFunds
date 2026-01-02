package aggregator

import (
	"backend/storage/postgres"
	"backend/types"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ScheduleAAggregator handles aggregating Schedule A data from Postgres to MongoDB
type ScheduleAAggregator struct {
	pgRepo      *postgres.Repository
	mongoDB     *mongo.Database
	mongoClient *mongo.Client
}

// NewScheduleAAggregator creates a new aggregator instance
func NewScheduleAAggregator(pgRepo *postgres.Repository, mongoClient *mongo.Client, dbName string) *ScheduleAAggregator {
	return &ScheduleAAggregator{
		pgRepo:      pgRepo,
		mongoDB:     mongoClient.Database(dbName),
		mongoClient: mongoClient,
	}
}

// RunAggregation runs all aggregation processes
func (a *ScheduleAAggregator) RunAggregation(ctx context.Context, cycle int) error {
	log.Printf("Starting aggregation for cycle %d", cycle)

	// Run aggregations sequentially (pgx.Conn is not thread-safe for concurrent use)
	if err := a.AggregateCommitteeContributors(ctx, cycle); err != nil {
		return err
	}

	if err := a.AggregateCandidateReceipts(ctx, cycle); err != nil {
		return err
	}

	if err := a.AggregateStateStats(ctx, cycle); err != nil {
		return err
	}

	log.Printf("Aggregation complete for cycle %d", cycle)
	return nil
}

// AggregateCommitteeContributors aggregates contribution data by committee contributor
// Note: Only includes committee-to-committee contributions (not individuals)
func (a *ScheduleAAggregator) AggregateCommitteeContributors(ctx context.Context, cycle int) error {
	log.Println("Aggregating committee contributors...")

	// Get all committee contributors from Postgres
	contributors, err := a.pgRepo.GetAllCommitteeContributors(ctx, cycle)
	if err != nil {
		return fmt.Errorf("failed to get committee contributors: %w", err)
	}

	collection := a.mongoDB.Collection("contributor_stats")

	for _, contrib := range contributors {
		// Get top candidates for this contributor
		topCandidates, err := a.pgRepo.GetTopCandidatesByContributor(ctx, contrib.ContributorID, 10)
		if err != nil {
			log.Printf("Warning: failed to get top candidates for %s: %v", contrib.ContributorID, err)
			topCandidates = []types.CandidateContribution{}
		} else if len(topCandidates) == 0 {
			topCandidates = []types.CandidateContribution{}
		}

		stats := types.ContributorStats{
			ContributorID:    contrib.ContributorID,
			Name:             contrib.Name,
			State:            contrib.State,
			TotalContributed: contrib.TotalAmount,
			ReceiptCount:     contrib.ReceiptCount,
			TopCandidates:    topCandidates,
			Cycle:            cycle,
			LastUpdated:      time.Now(),
		}

		// Upsert to MongoDB
		filter := bson.M{"contributor_id": contrib.ContributorID, "cycle": cycle}
		update := bson.M{"$set": stats}
		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			return fmt.Errorf("failed to upsert contributor %s: %w", contrib.ContributorID, err)
		}
	}

	log.Printf("Aggregated %d committee contributors", len(contributors))
	return nil
}

// AggregateCandidateReceipts aggregates receipt data by candidate
func (a *ScheduleAAggregator) AggregateCandidateReceipts(ctx context.Context, cycle int) error {
	log.Println("Aggregating candidate receipts...")

	// Get all candidates with receipt totals
	candidates, err := a.pgRepo.GetAllCandidatesWithTotals(ctx, cycle)
	if err != nil {
		return fmt.Errorf("failed to get candidates: %w", err)
	}

	collection := a.mongoDB.Collection("candidate_stats")

	for _, cand := range candidates {
		// Get top contributors for this candidate
		topContributors, err := a.pgRepo.GetTopContributorsByCandidate(ctx, cand.CandidateID, cycle, 10)
		if err != nil {
			log.Printf("Warning: failed to get top contributors for %s: %v", cand.CandidateID, err)
			topContributors = []types.ContributorContribution{}
		} else if len(topContributors) == 0 {
			// Ensure we store an empty array rather than null in Mongo
			topContributors = []types.ContributorContribution{}
		}

		stats := types.CandidateStats{
			CandidateID:     cand.CandidateID,
			Name:            cand.Name,
			Party:           cand.Party,
			Office:          cand.Office,
			TotalReceived:   cand.TotalAmount,
			ReceiptCount:    cand.ReceiptCount,
			TopContributors: topContributors,
			Cycle:           cycle,
			LastUpdated:     time.Now(),
		}

		// Upsert to MongoDB
		filter := bson.M{"candidate_id": cand.CandidateID, "cycle": cycle}
		update := bson.M{"$set": stats}
		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			return fmt.Errorf("failed to upsert candidate %s: %w", cand.CandidateID, err)
		}
	}

	log.Printf("Aggregated %d candidates", len(candidates))
	return nil
}

// AggregateStateStats aggregates contribution data by state
func (a *ScheduleAAggregator) AggregateStateStats(ctx context.Context, cycle int) error {
	log.Println("Aggregating state statistics...")

	// Get state-level totals
	states, err := a.pgRepo.GetStateContributionTotals(ctx, cycle)
	if err != nil {
		return fmt.Errorf("failed to get state totals: %w", err)
	}

	collection := a.mongoDB.Collection("state_stats")

	for _, state := range states {
		// Get top candidates in this state
		topCandidates, err := a.pgRepo.GetTopCandidatesByState(ctx, state.State, cycle, 10)
		if err != nil {
			log.Printf("Warning: failed to get top candidates for %s: %v", state.State, err)
			topCandidates = []types.CandidateContribution{}
		}

		stats := types.StateStats{
			State:         state.State,
			TotalAmount:   state.TotalAmount,
			ReceiptCount:  state.ReceiptCount,
			TopCandidates: topCandidates,
			Cycle:         cycle,
			LastUpdated:   time.Now(),
		}

		// Upsert to MongoDB
		filter := bson.M{"state": state.State, "cycle": cycle}
		update := bson.M{"$set": stats}
		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			return fmt.Errorf("failed to upsert state %s: %w", state.State, err)
		}
	}

	log.Printf("Aggregated %d states", len(states))
	return nil
}
