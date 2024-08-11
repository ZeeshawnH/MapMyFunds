import { Document, Schema, model } from "mongoose";

export interface ICandidate extends Document {
    candidate_id: string;
    [key: string]: any;
}

const candidateSchema = new Schema<ICandidate>({
    candidate_id: {type: String, required: true, unique: true},
}, {strict: false});

export const Candidate = model<ICandidate>('Candidate', candidateSchema);