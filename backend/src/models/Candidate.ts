import mongoose, { Schema } from "mongoose";

// export interface ICandidate extends Document {
//   candidate_id: string;
//   [key: string]: any;
// }

const contributionSchema = new mongoose.Schema({
  candidate_id: { type: String, required: true, unique: true },
  candidate_last_name: String,
  candidate_party_affiliation: String,
  contributions: [
    {
      contributor_state: String,
      election_year: Number,
      net_receipts: Number,
      rounded_net_receipts: Number,
    },
  ],
});

export const Candidate = mongoose.model("Candidate", contributionSchema);
