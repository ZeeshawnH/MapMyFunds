package ingest

import (
	"backend/storage/postgres"
	"backend/types"
	"context"
	"fmt"
	"log"
)

// IngestCandidateCommitteeRelations populates the candidate_committees join table
// based on the candidate_ids array stored on committees.
//
// This is intended to be run as a separate stage after candidates and committees
// have been ingested.
func IngestCandidateCommitteeRelations(ctx context.Context, repo *postgres.Repository) error {
	committees, err := repo.GetCommitteesWithCandidateIDs(ctx)
	if err != nil {
		return fmt.Errorf("could not load committees with candidate_ids: %w", err)
	}

	for _, committee := range committees {
		for _, candidateID := range committee.CandidateIDs {
			// Ensure the candidate exists before inserting into the join table
			exists, err := repo.CandidateExists(ctx, candidateID)
			if err != nil {
				log.Printf("error checking candidate existence (candidate=%s): %v", candidateID, err)
				continue
			}
			if !exists {
				// Skip relationships for candidates we haven't ingested yet
				continue
			}

			rel := types.DBCandidateCommittee{
				CandidateID:      candidateID,
				CommitteeID:      committee.CommitteeID,
				RelationshipType: "primary", // relationship type based on OpenFEC candidate_ids
			}

			if err := repo.UpsertCandidateCommittee(ctx, rel); err != nil {
				log.Printf("failed to upsert candidate-committee relation (candidate=%s, committee=%s): %v", candidateID, committee.CommitteeID, err)
			}
		}
	}

	return nil
}
