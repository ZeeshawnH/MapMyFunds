export interface Contribution {
  CandidateID: string;
  CandidateName: string;
  CandidateParty: string;
  ContributorState: string;
  ElectionYear: number;
  NetReceipts: number;
}

export interface Candidate {
  CandidateID: string;
  CandidateName: string;
  CandidateParty: string;
  NetReceipts: number;
}

export interface CandidateResponse {
  candidate_id: string,
  name: string,
  office: string,
  office_full: string,
  party: string,
  party_full: string,
  image_url: string
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
