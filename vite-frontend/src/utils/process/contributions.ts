import type {
  Contribution,
  StateContributions,
} from "../../types/contributions";
import { candidateInfo } from "../constants/candidateConstants";

export const processContributions = (
  data: Contribution[]
): StateContributions => {
  const stateData: StateContributions = {};

  // Group by state
  data.forEach((contribution) => {
    const {
      CandidateID,
      CandidateName,
      CandidateParty,
      ContributorState,
      ElectionYear,
      NetReceipts,
    } = contribution;
    const candidate = candidateInfo[CandidateID];

    if (!candidate) return; // Skip if we don't have candidate info

    if (!stateData[ContributorState]) {
      stateData[ContributorState] = [];
    }

    // Find if candidate already exists in state
    const existingCandidate = stateData[ContributorState].find(
      (c) => c.CandidateID === CandidateID
    );

    if (existingCandidate) {
      existingCandidate.NetReceipts += NetReceipts;
    } else {
      stateData[ContributorState].push({
        CandidateID,
        CandidateName,
        CandidateParty,
        ContributorState,
        ElectionYear,
        NetReceipts,
      });
    }
  });

  // Sort candidates by total amount in each state
  Object.keys(stateData).forEach((state) => {
    stateData[state].sort((a, b) => b.NetReceipts - a.NetReceipts);
  });

  return stateData;
};
