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
	ReportYear int    `json:"report_year"`
	Cycle      int    `json:"two_year_transaction_period"`
	LoadDate   string `json:"load_date"`

	// Contributor summary
	ContributorID           string  `json:"contributor_id"`
	ContributorName         string  `json:"contributor_name"`
	ContributorCity         string  `json:"contributor_city"`
	ContributorState        string  `json:"contributor_state"`
	ContributorZip          string  `json:"contributor_zip"`
	ContributorStreet1      string  `json:"contributor_street_1"`
	ContributorStreet2      string  `json:"contributor_street_2"`
	ContributorEmployer     string  `json:"contributor_employer"`
	ContributorOccupation   string  `json:"contributor_occupation"`
	ContributorAggregateYTD float64 `json:"contributor_aggregate_ytd"`
	EntityType              string  `json:"entity_type"`
	EntityTypeDesc          string  `json:"entity_type_desc"`
	IsIndividual            bool    `json:"is_individual"`

	// Nested committee receiving the money
	Committee   Committee `json:"committee"`
	CommitteeID string    `json:"committee_id"`

	// Nested contributor committee (PAC → committee)
	ContributorCommittee Committee `json:"contributor"`

	// Nested conduit committee (optional)
	ConduitCommittee *Committee `json:"conduit_committee"`

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
	CommitteeID string `json:"committee_id"`
	Name        string `json:"name"`

	AffiliatedName *string `json:"affiliated_committee_name"`

	CommitteeType     *string `json:"committee_type"`
	CommitteeTypeFull *string `json:"committee_type_full"`

	Designation     *string `json:"designation"`
	DesignationFull *string `json:"designation_full"`

	OrganizationType     *string `json:"organization_type"`
	OrganizationTypeFull *string `json:"organization_type_full"`

	Party     *string `json:"party"`
	PartyFull *string `json:"party_full"`

	City    *string `json:"city"`
	State   *string `json:"state"`
	Zip     *string `json:"zip"`
	Street1 *string `json:"street_1"`
	Street2 *string `json:"street_2"`

	TreasurerName *string `json:"treasurer_name"`

	IsActive *bool `json:"is_active"`

	CandidateIDs        []string `json:"candidate_ids"`
	SponsorCandidateIDs []string `json:"sponsor_candidate_ids"`

	Cycle  *int  `json:"cycle"`
	Cycles []int `json:"cycles"`
}

// Postgres Types

// DBCandidate represents a row in the candidates table
type DBCandidate struct {
	CandidateID string  `db:"candidate_id"`
	Name        string  `db:"name"`
	Office      string  `db:"office"`
	Party       *string `db:"party"`

	State        *string `db:"state"`
	District     *string `db:"district"`
	ElectionYear *int    `db:"election_year"`
}

// DBCommittee represents a row in the committees table
type DBCommittee struct {
	CommitteeID       string  `db:"committee_id"`
	Name              string  `db:"name"`
	CommitteeType     *string `db:"committee_type"`
	CommitteeTypeFull *string `db:"committee_type_full"`

	Designation     *string `db:"designation"`
	DesignationFull *string `db:"designation_full"`

	Party *string `db:"party"`
	State *string `db:"state"`

	CandidateIDs []string `db:"candidate_ids"`
}

// DBCandidateCommittee represents a row in the candidate_committees join table
type DBCandidateCommittee struct {
	CandidateID      string `db:"candidate_id"`
	CommitteeID      string `db:"committee_id"`
	RelationshipType string `db:"relationship_type"`
}

// DBContributor represents a row in the contributors table
type DBContributor struct {
	ContributorID string
	Name          string
	State         string
}

type DBScheduleAReceipt struct {
	// Core receipt info
	FECReceiptID string    `db:"fec_receipt_id"`
	CommitteeID  string    `db:"committee_id"`
	Amount       float64   `db:"amount"`
	ReceiptDate  time.Time `db:"receipt_date"`
	Cycle        int       `db:"cycle"`

	// Optional contributor (only for committees/PACs)
	ContributorID *string `db:"contributor_id"`

	// Individual contributor fields (always populated for context)
	IsIndividual          bool   `db:"is_individual"`
	ContributorName       string `db:"contributor_name"`
	ContributorFirstName  string `db:"contributor_first_name"`
	ContributorMiddleName string `db:"contributor_middle_name"`
	ContributorLastName   string `db:"contributor_last_name"`
	ContributorPrefix     string `db:"contributor_prefix"`
	ContributorSuffix     string `db:"contributor_suffix"`
	ContributorStreet1    string `db:"contributor_street1"`
	ContributorStreet2    string `db:"contributor_street2"`
	ContributorCity       string `db:"contributor_city"`
	ContributorState      string `db:"contributor_state"`
	ContributorZip        string `db:"contributor_zip"`
	ContributorEmployer   string `db:"contributor_employer"`
	ContributorOccupation string `db:"contributor_occupation"`

	// Conduit committee fields
	ConduitCommitteeID      *string `db:"conduit_committee_id"`
	ConduitCommitteeName    string  `db:"conduit_committee_name"`
	ConduitCommitteeCity    string  `db:"conduit_committee_city"`
	ConduitCommitteeState   string  `db:"conduit_committee_state"`
	ConduitCommitteeStreet1 string  `db:"conduit_committee_street1"`
	ConduitCommitteeStreet2 string  `db:"conduit_committee_street2"`
	ConduitCommitteeZip     string  `db:"conduit_committee_zip"`

	// Memo / special fields
	MemoCode       string `db:"memo_code"`
	MemoText       string `db:"memo_text"`
	MemoedSubtotal bool   `db:"memoed_subtotal"`

	// Filing / source metadata
	FileNumber       int       `db:"file_number"`
	ImageNumber      string    `db:"image_number"`
	PDFURL           string    `db:"pdf_url"`
	FilingForm       string    `db:"filing_form"`
	ElectionType     string    `db:"election_type"`
	FecElectionYear  int       `db:"fec_election_year"`
	ScheduleType     string    `db:"schedule_type"`
	ScheduleTypeFull string    `db:"schedule_type_full"`
	LineNumber       string    `db:"line_number"`
	LineNumberLabel  string    `db:"line_number_label"`
	ReceiptType      string    `db:"receipt_type"`
	ReceiptTypeDesc  string    `db:"receipt_type_desc"`
	ReportType       string    `db:"report_type"`
	ReportYear       int       `db:"report_year"`
	LoadDate         time.Time `db:"load_date"`

	// Optional donor committee (PAC → committee)
	DonorCommitteeName string `db:"donor_committee_name"`
}

// DBIngestionState represents a row in the ingestion_state table
type DBIngestionState struct {
	Source          string
	LastIndex       *string
	LastReceiptDate *time.Time
	SortNullOnly    bool
	Page            *int // Optional: for pagination systems that use page numbers
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
	Cycle            int                     `bson:"cycle" json:"cycle"`
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
	Cycle           int                       `bson:"cycle" json:"cycle"`
	LastUpdated     time.Time                 `bson:"last_updated" json:"last_updated"`
}

// StateStats represents aggregated state-level data in MongoDB
type StateStats struct {
	State         string                  `bson:"state" json:"state"`
	TotalAmount   float64                 `bson:"total_amount" json:"total_amount"`
	ReceiptCount  int                     `bson:"receipt_count" json:"receipt_count"`
	TopCandidates []CandidateContribution `bson:"top_candidates" json:"top_candidates"`
	Cycle         int                     `bson:"cycle" json:"cycle"`
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
	RecordsProcessed  int
	RecordsInserted   int
	RecordsFailed     int
	LastProcessedDate time.Time
	StartTime         time.Time
	EndTime           time.Time
	PagesProcessed    int
}
