export interface Contribution {
  CandidateID: string;
  CandidateName: string;
  CandidateParty: string;
  ContributorState: string;
  ElectionYear: number;
  NetReceipts: number;
}

export interface Candidate {
  candidate_id: string;
  candidate_last_name: string;
  candidate_party_affiliation: string;
  total_amount: number;
}

export interface StateContributions {
  [stateCode: string]: Contribution[];
}

export interface CandidateInfo {
  last_name: string;
  party: string;
}

export interface CandidateImageMap {
  [key: string]: string;
}

export interface ApiResponse<T> {
  api_version: string;
  pagination: {
    count: number;
    is_count_exact: boolean;
    page: number;
    pages: number;
    per_page: number;
  };
  results: T[];
}
