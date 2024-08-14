import express from "express";
import { Candidate } from "../models/Candidate";
import { fetchContributionsData } from "../services/fetchData";

const router = express.Router();

router.get("/", async (req, res) => {
  try {
    const candidates = await Candidate.find();

    candidates.sort((a, b) => {
      const stateA = a.contributor_state ?? ""; // Use empty string if undefined or null
      const stateB = b.contributor_state ?? ""; // Use empty string if undefined or null

      // Handle sorting for strings
      return stateA.localeCompare(stateB);
    });

    res.status(200).json(candidates);
  } catch (error) {
    res.status(500).send({ error: "Error retrieving data" });
  }
});

export default router;
