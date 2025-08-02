import {
  type CandidateInfo,
  type CandidateImageMap,
} from "../../types/contributions";

// Candidate mapping
// TODO update with api data
export const candidateInfo: Record<string, CandidateInfo> = {
  P00000001: { last_name: "Trump", party: "REP" },
  P00000002: { last_name: "Harris", party: "DEM" },
  P00009423: { last_name: "Kennedy", party: "IND" },
};

// Update color scale
export const partyColor = (party: string, isHovered: boolean = false) => {
  switch (party?.toUpperCase()) {
    case "DEM":
      return isHovered ? "#4d94ff" : "#1a75ff";
    case "REP":
      return isHovered ? "#ff4d4d" : "#ff1a1a";
    case "IND":
      return isHovered ? "#ffd700" : "#ffcc00";
    default:
      return isHovered ? "#e6e6e6" : "#cccccc";
  }
};

export type Party = "DEM" | "REP" | "LIB" | "";

// Add image mapping
export const candidateImages: CandidateImageMap = {
  Trump: "https://mapmyfunds-images.s3.us-east-1.amazonaws.com/donaldtrump.jpeg",
  Harris: "https://mapmyfunds-images.s3.us-east-1.amazonaws.com/kamalaharris.png",
  Kennedy: "https://mapmyfunds-images.s3.us-east-1.amazonaws.com/rfk.png",
};
