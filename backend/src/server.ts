import express, { Request, Response } from "express";
import mongoose from "mongoose";
import bodyParser from "body-parser";
import axios from "axios";
import dotenv from "dotenv";
// import cron from "node-cron";
import { Candidate, ICandidate } from "./models/Candidate";
import candidateRoutes from "./routes/candidate";
import { fetchContributionsData } from "./services/fetchData";

dotenv.config();

const app = express();

const PORT = 8080;

app.use(bodyParser.json());
app.use(express.json());

const loans = async () => {
  let rtn = [];
  let page = 1;
  let response = await fetch(
    "https://api.open.fec.gov/v1/schedules/schedule_c/?min_amount=10000&page=1&per_page=100&sort=-incurred_date&sort_hide_null=true&sort_null_only=false&sort_nulls_last=true&api_key=mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC"
  );
  let data = await response.json();
  while (data.results && data.results.length > 0) {
    rtn.push(data.results);
    page += 1;
    response = await fetch(
      `https://api.open.fec.gov/v1/schedules/schedule_c/?min_amount=10000&page=${page}&per_page=100&sort=-incurred_date&sort_hide_null=true&sort_null_only=false&sort_nulls_last=true&api_key=mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC`
    );
    data = await response.json();
  }
  return rtn;
};

app.get("/loans", async (req, res) => {
  const data = await loans();
  res.json(data);
});

// Helper function to build query parameters
// const buildQueryParams = (params) => {
//   return Object.keys(params)
//     .map(
//       (key) => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`
//     )
//     .join("&");
// };

const API_KEY = process.env.API_KEY;
const base_url = process.env.base_url;
const presidential_endpoint = process.env.presidential;

// const presidential = async () => {
//   const queryParams = buildQueryParams({
//     api_key: API_KEY,
//     page: 1,
//     per_page: 100,
//     election_year: 2024,
//     contributor_state: "NC",
//     sort: "-net_receipts",
//     sort_hide_null: false,
//     sort_null_only: false,
//     sort_nulls_last: false,
//   });

//   let response = await fetch(
//     `${base_url}${presidential_endpoint}/?${queryParams}`
//   );
//   let data = await response.json();
//   return data;
// };

app.get("/presidential", async (req, res) => {
  res.send(await fetchElectionData());
});

const URI = process.env.MONGO_URI;

// const fetchElectionData = async () => {
//   try {
//     const response = await axios.get(`${base_url}${presidential_endpoint}/`, {
//       params: {
//         api_key: API_KEY,
//         page: 1,
//         per_page: 100,
//         election_year: 2024,
//         contributor_state: "NC",
//         sort: "-net_receipts",
//         sort_hide_null: false,
//         sort_null_only: false,
//         sort_nulls_last: false,
//       },
//     });
//     const { results } = response.data;

//     results.forEach(async (candidateData) => {
//       await Candidate.findOneAndUpdate(
//         { candidate_id: candidateData.candidate_id },
//         { $set: candidateData },
//         { upset: true, new: true }
//       );
//       console.log(`Updated candidate: ${candidateData.candidate_id}`);
//     });

//     return results;
//   } catch (error) {
//     console.error("Error fetching or updating data: ", error);
//   }
// };

// Route middleware
app.use("/api/candidates", candidateRoutes);

setInterval(fetchContributionsData, 86400000);

fetchContributionsData();

mongoose
  .connect(URI as string)
  .then(() => {
    console.log("Connected to MongoDB");
    app.listen(PORT, () => {
      console.log(`Server is running on port ${PORT}`);
    });
  })
  .catch((error) => {
    console.error("MongoDB connection error:", error);
  });
