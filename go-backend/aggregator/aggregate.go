package aggregator

import (
	"backend/openfec"
	"fmt"
)

type CandidateStats struct {
	CandidateID   string
	CandidateName string
	TotalAmount   float64
	ElectionYear  int
}

func FetchContributionAmountByStateAndCandidate() (map[string][]openfec.CandidateContribution, error) {
	var stateMap = map[string][]openfec.CandidateContribution{
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

	for state := range stateMap {
		contributions, err := openfec.GetContributions(2024)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch contributions from FEC service: %s", err)
		}

		stateMap[state] = contributions
	}

	return stateMap, nil
}
