"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Candidate = void 0;
const mongoose_1 = require("mongoose");
const candidateSchema = new mongoose_1.Schema({
    candidate_id: { type: String, required: true, unique: true },
}, { strict: false });
exports.Candidate = (0, mongoose_1.model)('Candidate', candidateSchema);
