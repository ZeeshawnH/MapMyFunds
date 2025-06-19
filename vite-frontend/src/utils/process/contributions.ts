import {
  Contribution,
  StateContributions,
  Candidate,
} from "../../types/contributions";
import { candidateInfo } from "../constants/candidateConstants";

export const processContributions = (
  data: Contribution[]
): StateContributions => {
  const stateData: StateContributions = {};

  // Group by state
  data.forEach((contribution) => {
    const { contribution_state, candidate_id, contribution_receipt_amount } =
      contribution;
    const candidate = candidateInfo[candidate_id];

    if (!candidate) return; // Skip if we don't have candidate info

    if (!stateData[contribution_state]) {
      stateData[contribution_state] = [];
    }

    // Find if candidate already exists in state
    const existingCandidate = stateData[contribution_state].find(
      (c) => c.candidate_id === candidate_id
    );

    if (existingCandidate) {
      existingCandidate.total_amount += contribution_receipt_amount;
    } else {
      stateData[contribution_state].push({
        candidate_id,
        candidate_last_name: candidate.last_name,
        candidate_party_affiliation: candidate.party,
        total_amount: contribution_receipt_amount,
      });
    }
  });

  // Sort candidates by total amount in each state
  Object.keys(stateData).forEach((state) => {
    stateData[state].sort((a, b) => b.total_amount - a.total_amount);
  });

  return stateData;
};
