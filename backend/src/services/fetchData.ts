import axios from "axios";
import { Candidate } from "../models/Candidate";
import dotenv from "dotenv";

dotenv.config();
const API_KEY = process.env.API_KEY;
const base_url = process.env.base_url;
const presidential_endpoint = process.env.presidential;

const URL = `${base_url}${presidential_endpoint}`;

export const fetchContributionsData = async () => {
  try {
    let page = 1;
    let totalPages = 1;

    while (page <= totalPages) {
      const response = await axios.get(URL, {
        params: {
          api_key: API_KEY,
          page,
          per_page: 100,
        },
      });

      const { pagination, results } = response.data;

      totalPages = pagination.pages;

      for (const contribution of results) {
        await Candidate.updateOne(
          { candidate_id: contribution.candidate_id },
          {
            $set: {
              candidate_id: contribution.candidate_id,
              candidate_last_name: contribution.candidate_last_name,
              candidate_party_affiliation:
                contribution.candidate_party_affiliation,
              contributor_state: contribution.contributor_state,
            },
            $push: {
              contributions: {
                election_year: contribution.election_year,
                net_receipts: contribution.net_receipts,
                rounded_net_receipts: contribution.rounded_net_receipts,
              },
            },
          },
          { upsert: true }
        );
      }

      page++;
    }

    console.log("Candidates data stored/updated successfully");
  } catch (error) {
    console.error("Error fetching or storing candidates:", error);
  }
};
