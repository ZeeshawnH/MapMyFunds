import express from "express";
import { Candidate } from "../models/Candidate";
import { fetchContributionsData } from "../services/fetchData";

const router = express.Router();

router.get("/", async (req, res) => {
  try {
    const candidates = await Candidate.find();
    // const candidates = await fetchContributionsData();
    res.status(200).json(candidates);
  } catch (error) {
    res.status(500).send({ error: "Error retrieving data" });
  }
});

export default router;
