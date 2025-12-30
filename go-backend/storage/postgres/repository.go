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
		INSERT INTO schedule_a_receipts (
			fec_receipt_id, committee_id, amount, receipt_date, cycle,
			contributor_id, is_individual,
			contributor_name, contributor_street1, contributor_street2,
			contributor_city, contributor_state, contributor_zip,
			contributor_employer, contributor_occupation,
			conduit_committee_id, conduit_committee_name,
			conduit_committee_city, conduit_committee_state,
			conduit_committee_street1, conduit_committee_street2, conduit_committee_zip,
			memo_code, memoed_subtotal,
			file_number, image_number, pdf_url, filing_form,
			election_type, fec_election_year,
			schedule_type, schedule_type_full, line_number, line_number_label,
			receipt_type, receipt_type_desc, report_type, report_year, load_date
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39
		)
		ON CONFLICT (fec_receipt_id) DO NOTHING
	`

	for _, receipt := range receipts {
		batch.Queue(query,
			receipt.FECReceiptID,
			receipt.CommitteeID,
			receipt.Amount,
			receipt.ReceiptDate,
			receipt.Cycle,
			receipt.ContributorID,
			receipt.IsIndividual,
			receipt.ContributorName,
			receipt.ContributorStreet1,
			receipt.ContributorStreet2,
			receipt.ContributorCity,
			receipt.ContributorState,
			receipt.ContributorZip,
			receipt.ContributorEmployer,
			receipt.ContributorOccupation,
			receipt.ConduitCommitteeID,
			receipt.ConduitCommitteeName,
			receipt.ConduitCommitteeCity,
			receipt.ConduitCommitteeState,
			receipt.ConduitCommitteeStreet1,
			receipt.ConduitCommitteeStreet2,
			receipt.ConduitCommitteeZip,
			receipt.MemoCode,
			receipt.MemoedSubtotal,
			receipt.FileNumber,
			receipt.ImageNumber,
			receipt.PDFURL,
			receipt.FilingForm,
			receipt.ElectionType,
			receipt.FecElectionYear,
			receipt.ScheduleType,
			receipt.ScheduleTypeFull,
			receipt.LineNumber,
			receipt.LineNumberLabel,
			receipt.ReceiptType,
			receipt.ReceiptTypeDesc,
			receipt.ReportType,
			receipt.ReportYear,
			receipt.LoadDate,
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
		SELECT source, last_index, last_receipt_date, sort_null_only, page
		FROM ingestion_state
		WHERE source = $1
	`

	var state types.DBIngestionState
	err := r.conn.QueryRow(ctx, query, source).Scan(
		&state.Source,
		&state.LastIndex,
		&state.LastReceiptDate,
		&state.SortNullOnly,
		&state.Page,
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
		INSERT INTO ingestion_state (source, last_index, last_receipt_date, sort_null_only, page)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (source)
		DO UPDATE SET 
			last_index = EXCLUDED.last_index,
			last_receipt_date = EXCLUDED.last_receipt_date,
			sort_null_only = EXCLUDED.sort_null_only,
			page = EXCLUDED.page
	`

	_, err := r.conn.Exec(ctx, query, state.Source, state.LastIndex, state.LastReceiptDate, state.SortNullOnly, state.Page)
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

// ============================================================================
// Aggregation Queries for MongoDB ETL
// ============================================================================

// GetAllCommitteeContributors returns all committee contributors with totals for a cycle
// Only includes committee-to-committee contributions (contributor_id is not null)
func (r *Repository) GetAllCommitteeContributors(ctx context.Context, cycle int) ([]struct {
	ContributorID string
	Name          string
	State         string
	TotalAmount   float64
	ReceiptCount  int
}, error) {
	query := `
		SELECT 
			com.committee_id as contributor_id,
			com.name,
			'' as state,
			SUM(s.amount) as total_amount,
			COUNT(*) as receipt_count
		FROM schedule_a_receipts s
		JOIN committees com ON s.contributor_id = com.committee_id
		WHERE s.cycle = $1 AND s.contributor_id IS NOT NULL
		GROUP BY com.committee_id, com.name
		ORDER BY total_amount DESC
	`

	rows, err := r.conn.Query(ctx, query, cycle)
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

// GetAllCandidatesWithTotals returns all candidates with receipt totals for a cycle
func (r *Repository) GetAllCandidatesWithTotals(ctx context.Context, cycle int) ([]struct {
	CandidateID  string
	Name         string
	Party        string
	Office       string
	TotalAmount  float64
	ReceiptCount int
}, error) {
	query := `
		SELECT 
			can.candidate_id,
			can.name,
			can.party,
			can.office,
			COALESCE(SUM(s.amount), 0) as total_amount,
			COALESCE(COUNT(s.fec_receipt_id), 0) as receipt_count
		FROM candidates can
		LEFT JOIN committees com ON com.committee_id LIKE can.candidate_id || '%'
		LEFT JOIN schedule_a_receipts s ON s.committee_id = com.committee_id AND s.cycle = $1
		GROUP BY can.candidate_id, can.name, can.party, can.office
		HAVING COUNT(s.fec_receipt_id) > 0
		ORDER BY total_amount DESC
	`

	rows, err := r.conn.Query(ctx, query, cycle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		CandidateID  string
		Name         string
		Party        string
		Office       string
		TotalAmount  float64
		ReceiptCount int
	}

	for rows.Next() {
		var result struct {
			CandidateID  string
			Name         string
			Party        string
			Office       string
			TotalAmount  float64
			ReceiptCount int
		}
		if err := rows.Scan(&result.CandidateID, &result.Name, &result.Party, &result.Office, &result.TotalAmount, &result.ReceiptCount); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// GetStateContributionTotals returns contribution totals by state for a cycle
func (r *Repository) GetStateContributionTotals(ctx context.Context, cycle int) ([]struct {
	State        string
	TotalAmount  float64
	ReceiptCount int
}, error) {
	query := `
		SELECT 
			COALESCE(contributor_state, 'Unknown') as state,
			SUM(amount) as total_amount,
			COUNT(*) as receipt_count
		FROM schedule_a_receipts
		WHERE cycle = $1 AND is_individual = true
		GROUP BY contributor_state
		ORDER BY total_amount DESC
	`

	rows, err := r.conn.Query(ctx, query, cycle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		State        string
		TotalAmount  float64
		ReceiptCount int
	}

	for rows.Next() {
		var result struct {
			State        string
			TotalAmount  float64
			ReceiptCount int
		}
		if err := rows.Scan(&result.State, &result.TotalAmount, &result.ReceiptCount); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// GetTopCandidatesByState returns top candidates receiving contributions from a state
func (r *Repository) GetTopCandidatesByState(ctx context.Context, state string, cycle int, limit int) ([]types.CandidateContribution, error) {
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
		WHERE s.contributor_state = $1 AND s.cycle = $2 AND s.is_individual = true
		GROUP BY can.candidate_id, can.name, can.party
		ORDER BY total_amount DESC
		LIMIT $3
	`

	rows, err := r.conn.Query(ctx, query, state, cycle, limit)
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
