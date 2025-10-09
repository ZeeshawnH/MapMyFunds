package types

type Candidate struct {
	CandidateID         string `json:"candidate_id" bson:"candidate_id"`
	CandidateName       string `json:"name" bson:"name"`
	CandidateOffice     string `json:"office" bson:"office"`
	CandidateOfficeFull string `json:"office_full" bson:"office_full"`
	CandidatePartyAbbr  string `json:"party" bson:"party"`
	CandidatePartyFull  string `json:"party_full" bson:"party_full"`
	ImageURL            string `json:"image_url" bson:"image_url"`
}
