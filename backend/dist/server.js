"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const mongoose_1 = __importDefault(require("mongoose"));
const body_parser_1 = __importDefault(require("body-parser"));
const axios_1 = __importDefault(require("axios"));
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
const app = (0, express_1.default)();
const PORT = 8080;
app.use(body_parser_1.default.json());
const loans = () => __awaiter(void 0, void 0, void 0, function* () {
    let rtn = [];
    let page = 1;
    let response = yield fetch("https://api.open.fec.gov/v1/schedules/schedule_c/?min_amount=10000&page=1&per_page=100&sort=-incurred_date&sort_hide_null=true&sort_null_only=false&sort_nulls_last=true&api_key=mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC");
    let data = yield response.json();
    while (data.results && data.results.length > 0) {
        rtn.push(data.results);
        page += 1;
        response = yield fetch(`https://api.open.fec.gov/v1/schedules/schedule_c/?min_amount=10000&page=${page}&per_page=100&sort=-incurred_date&sort_hide_null=true&sort_null_only=false&sort_nulls_last=true&api_key=mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC`);
        data = yield response.json();
    }
    return rtn;
});
app.get("/loans", (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const data = yield loans();
    res.json(data);
}));
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
app.get("/presidential", (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    res.send(yield fetchElectionData());
}));
const uri = "mongodb+srv://zeeshawnahasnain:fgehXrHvN1Jzncqx@cluster0.u3t6l.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
mongoose_1.default.connect(uri);
mongoose_1.default.connection.on("connected", () => {
    console.log("Connected to MongoDB");
});
mongoose_1.default.connection.on("error", (err) => {
    console.error(`MongoDB connection error: ${err}`);
});
const fetchElectionData = () => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const response = yield axios_1.default.get(`${base_url}${presidential_endpoint}/`, {
            params: {
                api_key: API_KEY,
                page: 1,
                per_page: 100,
                election_year: 2024,
                contributor_state: "NC",
                sort: "-net_receipts",
                sort_hide_null: false,
                sort_null_only: false,
                sort_nulls_last: false,
            },
        });
        const { results } = response.data;
        return results;
    }
    catch (error) {
        console.error("Error fetching or updating data: ", error);
    }
});
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
