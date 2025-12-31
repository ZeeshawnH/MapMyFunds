package ingest

import (
	"backend/openfec"
	"backend/storage/postgres"
	"backend/types"
	"context"
	"fmt"
	"log"
)

const (
	ingestionStateName = "committees"
)

func IngestCommitteeInfo(ctx context.Context, repo *postgres.Repository) error {
	state, err := repo.GetIngestionState(ctx, ingestionStateName)
	if err != nil {
		return fmt.Errorf("Could not get ingestion progress: %w", err)
	}

	// Initialize page if nil
	page := 1
	if state.Page != nil {
		page = *state.Page + 1
	}

	committeeResponse, err := openfec.FetchCommitteeDataFromFEC(page)
	if err != nil {
		return fmt.Errorf("Could not fetch OpenFEC data: %w", err)
	}

	totalPages := committeeResponse.Pagination.Pages
	for page < totalPages {
		committeeResponse, err := openfec.FetchCommitteeDataFromFEC(page)
		if err != nil {
			return fmt.Errorf("Could not fetch OpenFEC data: %w", err)
		}

		results := committeeResponse.Results

		for _, committee := range results {
			dbCommittee := types.DBCommittee{
				CommitteeID:       committee.CommitteeID,
				Name:              committee.Name,
				CommitteeType:     committee.CommitteeType,
				CommitteeTypeFull: committee.CommitteeTypeFull,
				Designation:       committee.Designation,
				DesignationFull:   committee.DesignationFull,
				Party:             committee.Party,
				State:             committee.State,
				CandidateIDs:      committee.CandidateIDs,
			}

			repo.UpsertCommittee(ctx, dbCommittee)
		}

		page += 1
		state.Page = &page
		repo.UpdateIngestionState(ctx, *state)
		log.Println("Done with page ", page)
	}

	return nil
}
