package ingest

import (
	"backend/openfec"
	"backend/storage/postgres"
	"backend/types"
	"context"
	"fmt"
)

func IngestCandidateInfo(ctx context.Context, repo *postgres.Repository, year int) error {
	candidates, err := openfec.GetCandidateData(year)
	if err != nil {
		return fmt.Errorf("Failed to get candidate data from OpenFEC: %s", err)
	}

	fmt.Println("Candidate data retrieved from OpenFEC successfully. Now storing in database.")

	for _, candidate := range candidates {
		dbCandidate := types.DBCandidate{
			CandidateID: candidate.CandidateID,
			Name:        candidate.CandidateName,
			Office:      candidate.CandidateOffice,
			Party:       &candidate.CandidatePartyAbbr,
		}
		if err := repo.UpsertCandidate(ctx, dbCandidate); err != nil {
			message := fmt.Errorf("Error upserting candidate with ID %s: %s", candidate.CandidateID, err)
			fmt.Println(message)
			continue
		}
	}

	return nil
}
