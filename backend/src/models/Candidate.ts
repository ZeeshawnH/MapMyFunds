import mongoose, { Schema, model, Document } from "mongoose";

export interface IContribution {
  election_year: number;
  net_receipts: number;
  rounded_net_receipts: number;
}

const candidateSchema = new Schema({
  candidate_id: { type: String, required: true, unique: true },
  candidate_last_name: String,
  candidate_party_affiliation: String,
  contributor_state: String,
  contributions: [
    {
      election_year: Number,
      net_receipts: Number,
      rounded_net_receipts: Number,
    },
  ],
});

export const Candidate = model("Candidate", candidateSchema);
