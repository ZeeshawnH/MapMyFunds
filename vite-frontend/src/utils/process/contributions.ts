import type { StateContributions } from "../../types/contributions";

// Sort contributions for each state in-place by net receipts.
// We now keep all rows, including aggregate ones like "All candidates",
// and filter them at the visualization layer instead.
export const sortContributions = (
  data: StateContributions
): StateContributions => {
  Object.keys(data).forEach((stateCode) => {
    const stateContributions = data[stateCode];
    stateContributions.sort((a, b) => b.NetReceipts - a.NetReceipts);
    data[stateCode] = stateContributions;
  });

  return data;
};
