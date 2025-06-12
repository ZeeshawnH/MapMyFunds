package aggregator

import (
	"elections-backend/openfec"
)

type FECResult = openfec.FECResult

// AggregateContributions aggregates contributions by state and returns a map of state to total contributions.
func TopCandidatesByState(contributions []openfec.FECResult) (map[string][]openfec.FECResult, error) {
	states := make(map[string][]FECResult)

	for _, contribution := range contributions {
		if contribution.ContributionState == "" {
			continue // Skip contributions without a state
		}

		// Initialize the state if it doesn't exist
		if _, exists := states[contribution.ContributionState]; !exists {
			states[contribution.ContributionState] = []FECResult{}
		}

		states[contribution.ContributionState] = append(states[contribution.ContributionState], contribution)
	}

	return states, nil
}
