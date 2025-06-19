import type { Contribution } from "../types/contributions";

export const mockContributions: { [key: string]: Contribution[] } = {
  AK: [
    {
      CandidateID: "P00009423",
      CandidateName: "Harris",
      CandidateParty: "DEM",
      ContributorState: "AK",
      ElectionYear: 2024,
      NetReceipts: 2000104,
    },
    {
      CandidateID: "P00000001",
      CandidateName: "Trump",
      CandidateParty: "REP",
      ContributorState: "AK",
      ElectionYear: 2024,
      NetReceipts: 1650000,
    },
  ],
  CA: [
    {
      CandidateID: "P00009423",
      CandidateName: "Harris",
      CandidateParty: "DEM",
      ContributorState: "CA",
      ElectionYear: 2024,
      NetReceipts: 171901187,
    },
    {
      CandidateID: "P00000001",
      CandidateName: "Trump",
      CandidateParty: "REP",
      ContributorState: "CA",
      ElectionYear: 2024,
      NetReceipts: 64213117,
    },
  ],
  FL: [
    {
      CandidateID: "P00009423",
      CandidateName: "Harris",
      CandidateParty: "DEM",
      ContributorState: "FL",
      ElectionYear: 2024,
      NetReceipts: 76490969,
    },
    {
      CandidateID: "P00000001",
      CandidateName: "Trump",
      CandidateParty: "REP",
      ContributorState: "FL",
      ElectionYear: 2024,
      NetReceipts: 103146477,
    },
  ],
  NY: [
    {
      CandidateID: "P00009423",
      CandidateName: "Harris",
      CandidateParty: "DEM",
      ContributorState: "NY",
      ElectionYear: 2024,
      NetReceipts: 95981414,
    },
    {
      CandidateID: "P00000001",
      CandidateName: "Trump",
      CandidateParty: "REP",
      ContributorState: "NY",
      ElectionYear: 2024,
      NetReceipts: 100596130,
    },
  ],
  TX: [
    {
      CandidateID: "P00009423",
      CandidateName: "Harris",
      CandidateParty: "DEM",
      ContributorState: "TX",
      ElectionYear: 2024,
      NetReceipts: 45000000,
    },
    {
      CandidateID: "P00000001",
      CandidateName: "Trump",
      CandidateParty: "REP",
      ContributorState: "TX",
      ElectionYear: 2024,
      NetReceipts: 80000000,
    },
  ],
};
