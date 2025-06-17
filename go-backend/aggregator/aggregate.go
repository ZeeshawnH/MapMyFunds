package aggregator

import (
	"backend/db"
	"fmt"
)

type CandidateStats struct {
	CandidateID   string
	CandidateName string
	TotalAmount   float64
	ElectionYear  int
}

func FetchContributionAmountByStateAndCandidate() (map[string][]db.Contribution, error) {
	var stateMap = map[string][]db.Contribution{
		"AL": {},
		"AK": {},
		"AZ": {},
		"AR": {},
		"CA": {},
		"CO": {},
		"CT": {},
		"DE": {},
		"FL": {},
		"GA": {},
		"HI": {},
		"ID": {},
		"IL": {},
		"IN": {},
		"IA": {},
		"KS": {},
		"KY": {},
		"LA": {},
		"ME": {},
		"MD": {},
		"MA": {},
		"MI": {},
		"MN": {},
		"MS": {},
		"MO": {},
		"MT": {},
		"NE": {},
		"NV": {},
		"NH": {},
		"NJ": {},
		"NM": {},
		"NY": {},
		"NC": {},
		"ND": {},
		"OH": {},
		"OK": {},
		"OR": {},
		"PA": {},
		"RI": {},
		"SC": {},
		"SD": {},
		"TN": {},
		"TX": {},
		"UT": {},
		"VT": {},
		"VA": {},
		"WA": {},
		"WV": {},
		"WI": {},
		"WY": {},
		"DC": {},
	}

	client, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}

	for state := range stateMap {
		contributions, err := db.GetContributionsByStateAndYear(client, 2024, state)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch contributions from db: %s", err)
		}

		stateMap[state] = contributions
	}

	return stateMap, nil
}
