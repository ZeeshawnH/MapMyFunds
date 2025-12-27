package types

import "time"


// OpenFEC Types 

// ContributorReceipt represents a Schedule A receipt response from OpenFEC API
type ContributorReceipt struct {
	// Core receipt info
	SubID                      string  `json:"sub_id"`
	TransactionID              string  `json:"transaction_id"`
	BackReferenceTransactionID string  `json:"back_reference_transaction_id"`
	BackReferenceScheduleName  string  `json:"back_reference_schedule_name"`
	ContributionAmount         float64 `json:"contribution_receipt_amount"`
	ContributionDate           string  `json:"contribution_receipt_date"`

	ScheduleType     string `json:"schedule_type"`
	ScheduleTypeFull string `json:"schedule_type_full"`
	LineNumber       string `json:"line_number"`
	LineNumberLabel  string `json:"line_number_label"`

	ReceiptType     string `json:"receipt_type"`
	ReceiptTypeDesc string `json:"receipt_type_desc"`

	MemoCode       string `json:"memo_code"`
	MemoedSubtotal bool   `json:"memoed_subtotal"`

	ReportType string `json:"report_type"`
	ReportYear string `json:"report_year"`
	Cycle      int    `json:"two_year_transaction_period"`
	LoadDate   string `json:"load_date"`

	// Contributor summary
	ContributorID               string  `json:"contributor_id"`
	ContributorName             string  `json:"contributor_name"`
	ContributorCity             string  `json:"contributor_city"`
	ContributorState            string  `json:"contributor_state"`
	ContributorZip              string  `json:"contributor_zip"`
	ContributorStreet1          string  `json:"contributor_street_1"`
	ContributorStreet2          string  `json:"contributor_street_2"`
	ContributorEmployer         string  `json:"contributor_employer"`
	ContributorOccupation       string  `json:"contributor_occupation"`
	ContributorAggregateYTD     float64 `json:"contributor_aggregate_ytd"`
	EntityType                  string  `json:"entity_type"`
	EntityTypeDesc              string  `json:"entity_type_desc"`
	IsIndividual                bool    `json:"is_individual"`

	// Nested committee receiving the money
	Committee   Committee `json:"committee"`
	CommitteeID string    `json:"committee_id"`

	// Nested contributor committee (PAC → committee)
	ContributorCommittee Committee `json:"contributor"`

	// Filing metadata
	FileNumber      int    `json:"file_number"`
	ImageNumber     string `json:"image_number"`
	PDFURL          string `json:"pdf_url"`
	FilingForm      string `json:"filing_form"`
	ElectionType    string `json:"election_type"`
	FecElectionYear string `json:"fec_election_year"`
}

// Committee represents committee data response from OpenFEC API
type Committee struct {
	CommitteeID              string `json:"committee_id"`
	Name                     string `json:"name"`
	AffiliatedName           string `json:"affiliated_committee_name"`
	CommitteeType            string `json:"committee_type"`
	CommitteeTypeFull        string `json:"committee_type_full"`
	Designation              string `json:"designation"`
	DesignationFull          string `json:"designation_full"`
	OrganizationType         string `json:"organization_type"`
	OrganizationTypeFull     string `json:"organization_type_full"`
	Party                    string `json:"party"`
	PartyFull                string `json:"party_full"`
	City                     string `json:"city"`
	State                    string `json:"state"`
	Zip                      string `json:"zip"`
	Street1                  string `json:"street_1"`
	Street2                  string `json:"street_2"`
	TreasurerName            string `json:"treasurer_name"`
	IsActive                 bool   `json:"is_active"`
	Cycle                    int    `json:"cycle"`
	Cycles                   []int  `json:"cycles"`
}


// Postgres Types

// DBCandidate represents a row in the candidates table
type DBCandidate struct {
	CandidateID string
	Name        string
	Office      string
	Party       string
}

// DBCommittee represents a row in the committees table
type DBCommittee struct {
	CommitteeID   string
	Name          string
	CommitteeType string
}

// DBContributor represents a row in the contributors table
type DBContributor struct {
	ContributorID string
	Name          string
	State         string
}

// DBScheduleAReceipt represents a row in the schedule_a_receipts table
type DBScheduleAReceipt struct {
	FECReceiptID  string
	ContributorID string
	CommitteeID   string
	Amount        float64
	ReceiptDate   time.Time
	Cycle         int
}

// DBIngestionState represents a row in the ingestion_state table
type DBIngestionState struct {
	Source          string
	LastIndex       *string
	LastReceiptDate *time.Time
	SortNullOnly    bool
}

// MongoDB Types

// ContributorStats represents aggregated contributor data in MongoDB
type ContributorStats struct {
	ContributorID    string                  `bson:"contributor_id" json:"contributor_id"`
	Name             string                  `bson:"name" json:"name"`
	State            string                  `bson:"state" json:"state"`
	TotalContributed float64                 `bson:"total_contributed" json:"total_contributed"`
	ReceiptCount     int                     `bson:"receipt_count" json:"receipt_count"`
	TopCandidates    []CandidateContribution `bson:"top_candidates" json:"top_candidates"`
	LastUpdated      time.Time               `bson:"last_updated" json:"last_updated"`
}

// CandidateStats represents aggregated candidate data in MongoDB
type CandidateStats struct {
	CandidateID     string                    `bson:"candidate_id" json:"candidate_id"`
	Name            string                    `bson:"name" json:"name"`
	Party           string                    `bson:"party" json:"party"`
	Office          string                    `bson:"office" json:"office"`
	TotalReceived   float64                   `bson:"total_received" json:"total_received"`
	ReceiptCount    int                       `bson:"receipt_count" json:"receipt_count"`
	TopContributors []ContributorContribution `bson:"top_contributors" json:"top_contributors"`
	LastUpdated     time.Time                 `bson:"last_updated" json:"last_updated"`
}

// StateStats represents aggregated state-level data in MongoDB
type StateStats struct {
	State         string                  `bson:"state" json:"state"`
	TotalAmount   float64                 `bson:"total_amount" json:"total_amount"`
	ReceiptCount  int                     `bson:"receipt_count" json:"receipt_count"`
	TopCandidates []CandidateContribution `bson:"top_candidates" json:"top_candidates"`
	LastUpdated   time.Time               `bson:"last_updated" json:"last_updated"`
}

// CandidateContribution represents a candidate's contribution summary
type CandidateContribution struct {
	CandidateID  string  `bson:"candidate_id" json:"candidate_id"`
	Name         string  `bson:"name" json:"name"`
	Party        string  `bson:"party" json:"party"`
	TotalAmount  float64 `bson:"total_amount" json:"total_amount"`
	ReceiptCount int     `bson:"receipt_count" json:"receipt_count"`
}

// ContributorContribution represents a contributor's contribution summary
type ContributorContribution struct {
	ContributorID string  `bson:"contributor_id" json:"contributor_id"`
	Name          string  `bson:"name" json:"name"`
	State         string  `bson:"state" json:"state"`
	TotalAmount   float64 `bson:"total_amount" json:"total_amount"`
	ReceiptCount  int     `bson:"receipt_count" json:"receipt_count"`
}

// Helper Types for Data Transformation

// ContributorIdentifier contains the components used to generate a stable contributor ID
type ContributorIdentifier struct {
	Name  string
	City  string
	State string
	Zip   string
}

// IngestionMetrics tracks statistics during the ingestion process
type IngestionMetrics struct {
	RecordsProcessed    int
	RecordsInserted     int
	RecordsFailed       int
	LastProcessedDate   time.Time
	StartTime           time.Time
	EndTime             time.Time
	PagesProcessed      int
}
