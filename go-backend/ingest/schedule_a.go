package ingest

import (
	"backend/openfec"
	"backend/storage/postgres"
	"backend/types"
	"context"
	"fmt"
	"log"
	"time"
)

const maxScheduleARetries = 3

func ingestScheduleAReceipts(ctx context.Context, receipts []types.ContributorReceipt, repo *postgres.Repository) error {
	dbReceipts := make([]types.DBScheduleAReceipt, 0, len(receipts))

receiptLoop:
	for _, receipt := range receipts {
		// Skip receipts with no contributor information
		if receipt.ContributorName == "" && receipt.ContributorCommittee.CommitteeID == "" {
			continue
		}

		// Parse receipt date
		receiptDate, err := parseDate(receipt.ContributionDate)
		if err != nil {
			// Use zero time on parse error
			receiptDate = time.Time{}
		}

		var contributorID *string
		var conduitCommitteeID *string

		// Handle committee contributors (NOT individuals)
		if !receipt.IsIndividual && receipt.ContributorCommittee.CommitteeID != "" {
			// Upsert the contributor committee with retries
			for attempt := 1; attempt <= maxScheduleARetries; attempt++ {
				err := repo.UpsertCommittee(ctx, types.DBCommittee{
					CommitteeID:   receipt.ContributorCommittee.CommitteeID,
					Name:          receipt.ContributorCommittee.Name,
					CommitteeType: receipt.ContributorCommittee.CommitteeType,
				})
				if err == nil {
					break
				}
				if attempt == maxScheduleARetries {
					log.Printf("failed to upsert contributor committee %s after %d attempts: %v",
						receipt.ContributorCommittee.CommitteeID, maxScheduleARetries, err)
					// Skip this receipt but continue processing others
					continue receiptLoop
				}
				time.Sleep(200 * time.Millisecond)
			}

			contributorID = &receipt.ContributorCommittee.CommitteeID
		}
		// For individuals: contributorID remains nil (no entity created)

		// Upsert the receiving committee
		if receipt.CommitteeID != "" {
			// Upsert the receiving committee with retries
			for attempt := 1; attempt <= maxScheduleARetries; attempt++ {
				err := repo.UpsertCommittee(ctx, types.DBCommittee{
					CommitteeID:   receipt.CommitteeID,
					Name:          receipt.Committee.Name,
					CommitteeType: receipt.Committee.CommitteeType,
				})
				if err == nil {
					break
				}
				if attempt == maxScheduleARetries {
					log.Printf("failed to upsert receiving committee %s after %d attempts: %v",
						receipt.CommitteeID, maxScheduleARetries, err)
					// Skip this receipt but continue processing others
					continue receiptLoop
				}
				time.Sleep(200 * time.Millisecond)
			}
		}

		// Handle conduit committee if present
		if receipt.ConduitCommittee != nil && receipt.ConduitCommittee.CommitteeID != "" {
			// Upsert the conduit committee with retries
			for attempt := 1; attempt <= maxScheduleARetries; attempt++ {
				err := repo.UpsertCommittee(ctx, types.DBCommittee{
					CommitteeID:   receipt.ConduitCommittee.CommitteeID,
					Name:          receipt.ConduitCommittee.Name,
					CommitteeType: receipt.ConduitCommittee.CommitteeType,
				})
				if err == nil {
					break
				}
				if attempt == maxScheduleARetries {
					log.Printf("failed to upsert conduit committee %s after %d attempts: %v",
						receipt.ConduitCommittee.CommitteeID, maxScheduleARetries, err)
					// Skip this receipt but continue processing others
					continue receiptLoop
				}
				time.Sleep(200 * time.Millisecond)
			}

			conduitCommitteeID = &receipt.ConduitCommittee.CommitteeID
		}

		// Build the DB receipt with conduit committee details
		var conduitName, conduitCity, conduitState, conduitStreet1, conduitStreet2, conduitZip string
		if receipt.ConduitCommittee != nil {
			conduitName = receipt.ConduitCommittee.Name
			if receipt.ConduitCommittee.City != nil {
				conduitCity = *receipt.ConduitCommittee.City
			}
			if receipt.ConduitCommittee.State != nil {
				conduitState = *receipt.ConduitCommittee.State
			}
			if receipt.ConduitCommittee.Street1 != nil {
				conduitStreet1 = *receipt.ConduitCommittee.Street1
			}
			if receipt.ConduitCommittee.Street2 != nil {
				conduitStreet2 = *receipt.ConduitCommittee.Street2
			}
			if receipt.ConduitCommittee.Zip != nil {
				conduitZip = *receipt.ConduitCommittee.Zip
			}
		}

		// Parse load date
		loadDate := time.Now()
		if receipt.LoadDate != "" {
			if parsed, err := parseDate(receipt.LoadDate); err == nil {
				loadDate = parsed
			}
		}

		dbReceipt := types.DBScheduleAReceipt{
			// Core fields
			FECReceiptID: receipt.SubID,
			CommitteeID:  receipt.CommitteeID,
			Amount:       receipt.ContributionAmount,
			ReceiptDate:  receiptDate,
			Cycle:        receipt.Cycle,

			// Contributor reference
			ContributorID: contributorID,

			// Individual contributor fields (inline storage)
			IsIndividual:          receipt.IsIndividual,
			ContributorName:       receipt.ContributorName,
			ContributorStreet1:    receipt.ContributorStreet1,
			ContributorStreet2:    receipt.ContributorStreet2,
			ContributorCity:       receipt.ContributorCity,
			ContributorState:      receipt.ContributorState,
			ContributorZip:        receipt.ContributorZip,
			ContributorEmployer:   receipt.ContributorEmployer,
			ContributorOccupation: receipt.ContributorOccupation,

			// Conduit committee
			ConduitCommitteeID:      conduitCommitteeID,
			ConduitCommitteeName:    conduitName,
			ConduitCommitteeCity:    conduitCity,
			ConduitCommitteeState:   conduitState,
			ConduitCommitteeStreet1: conduitStreet1,
			ConduitCommitteeStreet2: conduitStreet2,
			ConduitCommitteeZip:     conduitZip,

			// Memo fields
			MemoCode:       receipt.MemoCode,
			MemoedSubtotal: receipt.MemoedSubtotal,

			// Filing metadata
			FileNumber:       receipt.FileNumber,
			ImageNumber:      receipt.ImageNumber,
			PDFURL:           receipt.PDFURL,
			FilingForm:       receipt.FilingForm,
			ElectionType:     receipt.ElectionType,
			FecElectionYear:  0,
			ScheduleType:     receipt.ScheduleType,
			ScheduleTypeFull: receipt.ScheduleTypeFull,
			LineNumber:       receipt.LineNumber,
			LineNumberLabel:  receipt.LineNumberLabel,
			ReceiptType:      receipt.ReceiptType,
			ReceiptTypeDesc:  receipt.ReceiptTypeDesc,
			ReportType:       receipt.ReportType,
			ReportYear:       receipt.ReportYear,
			LoadDate:         loadDate,
		}

		dbReceipts = append(dbReceipts, dbReceipt)
	}

	// Batch insert all receipts. If the batch fails, fall back to per-row upserts
	// with retries so that a few bad rows don't abort the entire page.
	if len(dbReceipts) > 0 {
		if err := repo.BatchInsertScheduleAReceipts(ctx, dbReceipts); err != nil {
			log.Printf("batch insert for %d receipts failed, falling back to per-row upserts: %v", len(dbReceipts), err)

			for _, r := range dbReceipts {
				for attempt := 1; attempt <= maxScheduleARetries; attempt++ {
					if err := repo.UpsertScheduleAReceipt(ctx, r); err != nil {
						if attempt == maxScheduleARetries {
							log.Printf("failed to upsert receipt %s after %d attempts: %v", r.FECReceiptID, maxScheduleARetries, err)
						} else {
							time.Sleep(200 * time.Millisecond)
						}
					} else {
						break
					}
				}
			}
		}
	}

	return nil

}

// RunScheduleAIngestion orchestrates the full ingestion process with pagination and state tracking
func RunScheduleAIngestion(ctx context.Context, repo *postgres.Repository, cycle int, maxPages int) error {
	source := fmt.Sprintf("schedule_a_%d", cycle)

	// Load ingestion state
	state, err := repo.GetIngestionState(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to get ingestion state: %w", err)
	}

	log.Printf("Starting Schedule A ingestion for cycle %d", cycle)
	if state.LastIndex != nil {
		log.Printf("Resuming from last_index: %s", *state.LastIndex)
	}

	pageCount := 0
	totalRecords := 0

	for {
		// Check if we've hit max pages limit
		if maxPages > 0 && pageCount >= maxPages {
			log.Printf("Reached max pages limit: %d", maxPages)
			break
		}

		// Build lastReceiptDate string for API
		var lastReceiptDateStr *string
		if state.LastReceiptDate != nil && !state.LastReceiptDate.IsZero() {
			dateStr := state.LastReceiptDate.Format("2006-01-02")
			lastReceiptDateStr = &dateStr
		}

		// Fetch page from OpenFEC
		log.Printf("Fetching page %d (last_index: %v, last_date: %v)",
			pageCount+1, state.LastIndex, lastReceiptDateStr)

		response, err := openfec.FetchContributorReceiptDataFromFEC(
			state.LastIndex,
			lastReceiptDateStr,
			&state.SortNullOnly,
			[]int{cycle},
		)
		if err != nil {
			return fmt.Errorf("failed to fetch page %d: %w", pageCount+1, err)
		}

		// Check if we got any results
		if len(response.Results) == 0 {
			log.Println("No more results, ingestion complete")
			break
		}

		log.Printf("Processing %d receipts from page %d", len(response.Results), pageCount+1)

		// Process and insert this batch
		if err := ingestScheduleAReceipts(ctx, response.Results, repo); err != nil {
			return fmt.Errorf("failed to ingest page %d: %w", pageCount+1, err)
		}

		totalRecords += len(response.Results)
		pageCount++

		// Update state with pagination cursors from API response
		if response.Pagination.LastIndexes.LastIndex != "" {
			state.LastIndex = &response.Pagination.LastIndexes.LastIndex
		}

		if response.Pagination.LastIndexes.LastContributionReceiptDate != "" {
			parsedDate, err := parseDate(response.Pagination.LastIndexes.LastContributionReceiptDate)
			if err == nil {
				parsedDate.Format("2006-01-02")
				state.LastReceiptDate = &parsedDate
				fmt.Println(&parsedDate)
			}
			fmt.Println(err)
		}

		// Save checkpoint after each page
		if err := repo.UpdateIngestionState(ctx, *state); err != nil {
			log.Printf("Warning: failed to update ingestion state: %v", err)
		}

		log.Printf("Completed page %d - Total records so far: %d", pageCount, totalRecords)

		// Rate limiting - be nice to the API
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("Ingestion complete: %d pages, %d total records", pageCount, totalRecords)
	return nil
}

func parseDate(dateStr string) (time.Time, error) {
	// OpenFEC format: "2024-01-15"
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}
