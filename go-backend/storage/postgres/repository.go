package postgres

import (
	"context"
	"fmt"

	"backend/types"

	"github.com/jackc/pgx/v5"
)

// Repository handles all database operations for Postgres/Supabase
type Repository struct {
	conn *pgx.Conn
}

// NewRepository creates a new repository instance
func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{conn: conn}
}

// ============================================================================
// Candidate Operations
// ============================================================================

// UpsertCandidate inserts or updates a candidate record
func (r *Repository) UpsertCandidate(ctx context.Context, candidate types.DBCandidate) error {
	query := `
		INSERT INTO candidates (candidate_id, name, office, party)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (candidate_id) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			office = EXCLUDED.office,
			party = EXCLUDED.party
	`
	_, err := r.conn.Exec(ctx, query, candidate.CandidateID, candidate.Name, candidate.Office, candidate.Party)
	return err
}

// Insert or update committee record
func (r *Repository) UpsertCommittee(ctx context.Context, committee types.DBCommittee) error {
	query := `
		INSERT INTO committees (committee_id, name, committee_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (committee_id)
		DO UPDATE SET 
			name = EXCLUDED.name,
			committee_type = EXCLUDED.committee_type
	`
	_, err := r.conn.Exec(ctx, query, committee.CommitteeID, committee.Name, committee.CommitteeType)
	return err
}

// Insert or update contributor
func (r *Repository) UpsertContributor(ctx context.Context, contributor types.DBContributor) error {
	query := `
		INSERT INTO contributors (contributor_id, name, state)
		VALUES ($1, $2, $3)
		ON CONFLICT (contributor_id)
		DO UPDATE SET 
			name = EXCLUDED.name,
			state = EXCLUDED.state
	`
	_, err := r.conn.Exec(ctx, query, contributor.ContributorID, contributor.Name, contributor.State)
	return err
}

// Scheudule A Receipt
func (r *Repository) UpsertScheduleAReceipt(ctx context.Context, receipt types.DBScheduleAReceipt) error {
	query := `
		INSERT INTO schedule_a_receipts (fec_receipt_id, contributor_id, committee_id, amount, receipt_date, cycle)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (fec_receipt_id)
		DO UPDATE SET 
			contributor_id = EXCLUDED.contributor_id,
			committee_id = EXCLUDED.committee_id,
			amount = EXCLUDED.amount,
			receipt_date = EXCLUDED.receipt_date,
			cycle = EXCLUDED.cycle
	`
	_, err := r.conn.Exec(ctx, query, 
		receipt.FECReceiptID, 
		receipt.ContributorID, 
		receipt.CommitteeID, 
		receipt.Amount, 
		receipt.ReceiptDate, 
		receipt.Cycle,
	)
	return err
}

// Insert mutliple receipts at a time
func (r *Repository) BatchInsertScheduleAReceipts(ctx context.Context, receipts []types.DBScheduleAReceipt) error {
	if len(receipts) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO schedule_a_receipts (fec_receipt_id, contributor_id, committee_id, amount, receipt_date, cycle)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (fec_receipt_id) DO NOTHING
	`

	for _, receipt := range receipts {
		batch.Queue(query, 
			receipt.FECReceiptID, 
			receipt.ContributorID, 
			receipt.CommitteeID, 
			receipt.Amount, 
			receipt.ReceiptDate, 
			receipt.Cycle,
		)
	}

	results := r.conn.SendBatch(ctx, batch)
	defer results.Close()

	// Process all batch results
	for i := 0; i < len(receipts); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("batch insert failed at index %d: %w", i, err)
		}
	}

	return nil
}

// Get ingestion state progress
func (r *Repository) GetIngestionState(ctx context.Context, source string) (*types.DBIngestionState, error) {
	query := `
		SELECT source, last_index, last_receipt_date, sort_null_only
		FROM ingestion_state
		WHERE source = $1
	`

	var state types.DBIngestionState
	err := r.conn.QueryRow(ctx, query, source).Scan(
		&state.Source,
		&state.LastIndex,
		&state.LastReceiptDate,
		&state.SortNullOnly,
	)

	if err == pgx.ErrNoRows {
		// Return empty state for new source
		return &types.DBIngestionState{
			Source:       source,
			SortNullOnly: false,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &state, nil
}

// UpdateIngestionState saves the current ingestion state
func (r *Repository) UpdateIngestionState(ctx context.Context, state types.DBIngestionState) error {
	query := `
		INSERT INTO ingestion_state (source, last_index, last_receipt_date, sort_null_only)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (source)
		DO UPDATE SET 
			last_index = EXCLUDED.last_index,
			last_receipt_date = EXCLUDED.last_receipt_date,
			sort_null_only = EXCLUDED.sort_null_only
	`

	_, err := r.conn.Exec(ctx, query, state.Source, state.LastIndex, state.LastReceiptDate, state.SortNullOnly)
	return err
}

// Returns aggregated totals for a contributor
func (r *Repository) GetTotalContributionsByContributor(ctx context.Context, contributorID string) (float64, int, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) as total, COUNT(*) as count
		FROM schedule_a_receipts
		WHERE contributor_id = $1
	`

	var total float64
	var count int
	err := r.conn.QueryRow(ctx, query, contributorID).Scan(&total, &count)
	return total, count, err
}

// Returns aggregated totals for a candidate
// This joins through committees since contributions go to committees, not directly to candidates
func (r *Repository) GetTotalContributionsByCandidate(ctx context.Context, candidateID string) (float64, int, error) {
	query := `
		SELECT COALESCE(SUM(s.amount), 0) as total, COUNT(*) as count
		FROM schedule_a_receipts s
		JOIN committees c ON s.committee_id = c.committee_id
		WHERE c.committee_id LIKE $1 || '%'
	`

	var total float64
	var count int
	err := r.conn.QueryRow(ctx, query, candidateID).Scan(&total, &count)
	return total, count, err
}

// GetContributionsByState returns total contributions from a specific state
func (r *Repository) GetContributionsByState(ctx context.Context, state string) (float64, int, error) {
	query := `
		SELECT COALESCE(SUM(s.amount), 0) as total, COUNT(*) as count
		FROM schedule_a_receipts s
		JOIN contributors c ON s.contributor_id = c.contributor_id
		WHERE c.state = $1
	`

	var total float64
	var count int
	err := r.conn.QueryRow(ctx, query, state).Scan(&total, &count)
	return total, count, err
}

// GetTopContributorsByAmount returns the top N contributors by total contribution amount
func (r *Repository) GetTopContributorsByAmount(ctx context.Context, limit int) ([]struct {
	ContributorID string
	Name          string
	State         string
	TotalAmount   float64
	ReceiptCount  int
}, error) {
	query := `
		SELECT 
			c.contributor_id,
			c.name,
			c.state,
			SUM(s.amount) as total_amount,
			COUNT(*) as receipt_count
		FROM contributors c
		JOIN schedule_a_receipts s ON c.contributor_id = s.contributor_id
		GROUP BY c.contributor_id, c.name, c.state
		ORDER BY total_amount DESC
		LIMIT $1
	`

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		ContributorID string
		Name          string
		State         string
		TotalAmount   float64
		ReceiptCount  int
	}

	for rows.Next() {
		var result struct {
			ContributorID string
			Name          string
			State         string
			TotalAmount   float64
			ReceiptCount  int
		}
		if err := rows.Scan(&result.ContributorID, &result.Name, &result.State, &result.TotalAmount, &result.ReceiptCount); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// GetTopCandidatesByContributor returns the top N candidates that a contributor has contributed to
func (r *Repository) GetTopCandidatesByContributor(ctx context.Context, contributorID string, limit int) ([]types.CandidateContribution, error) {
	query := `
		SELECT 
			can.candidate_id,
			can.name,
			can.party,
			SUM(s.amount) as total_amount,
			COUNT(*) as receipt_count
		FROM schedule_a_receipts s
		JOIN committees com ON s.committee_id = com.committee_id
		JOIN candidates can ON com.committee_id LIKE can.candidate_id || '%'
		WHERE s.contributor_id = $1
		GROUP BY can.candidate_id, can.name, can.party
		ORDER BY total_amount DESC
		LIMIT $2
	`

	rows, err := r.conn.Query(ctx, query, contributorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.CandidateContribution
	for rows.Next() {
		var cc types.CandidateContribution
		if err := rows.Scan(&cc.CandidateID, &cc.Name, &cc.Party, &cc.TotalAmount, &cc.ReceiptCount); err != nil {
			return nil, err
		}
		results = append(results, cc)
	}

	return results, rows.Err()
}

// GetTopContributorsByCandidate returns the top N contributors to a specific candidate
func (r *Repository) GetTopContributorsByCandidate(ctx context.Context, candidateID string, limit int) ([]types.ContributorContribution, error) {
	query := `
		SELECT 
			c.contributor_id,
			c.name,
			c.state,
			SUM(s.amount) as total_amount,
			COUNT(*) as receipt_count
		FROM schedule_a_receipts s
		JOIN contributors c ON s.contributor_id = c.contributor_id
		JOIN committees com ON s.committee_id = com.committee_id
		WHERE com.committee_id LIKE $1 || '%'
		GROUP BY c.contributor_id, c.name, c.state
		ORDER BY total_amount DESC
		LIMIT $2
	`

	rows, err := r.conn.Query(ctx, query, candidateID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.ContributorContribution
	for rows.Next() {
		var cc types.ContributorContribution
		if err := rows.Scan(&cc.ContributorID, &cc.Name, &cc.State, &cc.TotalAmount, &cc.ReceiptCount); err != nil {
			return nil, err
		}
		results = append(results, cc)
	}

	return results, rows.Err()
}

