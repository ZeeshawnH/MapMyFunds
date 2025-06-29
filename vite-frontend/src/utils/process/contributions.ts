import type {
  StateContributions,
} from "../../types/contributions";

export const sortContributions = (
  data: StateContributions
): StateContributions => {

  Object.keys(data).forEach((stateCode) => {
    let stateContributions = data[stateCode]
    stateContributions = stateContributions.filter(candidate => 
      candidate.CandidateName != "All candidates" && 
      candidate.CandidateName != "Republicans" && 
      candidate.CandidateName != "Democrats"
    );
    stateContributions.sort((a, b) => b.NetReceipts - a.NetReceipts);
    stateContributions = stateContributions.slice(0, 4);
    data[stateCode] = stateContributions;
  });

  return data;
};
