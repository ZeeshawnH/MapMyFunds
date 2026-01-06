import axios from "axios";

export const fetchContributionsByState = async () => {
  const url = import.meta.env.VITE_API_URL;
  const path = "/api/contributions";

  try {
    const response = await axios.get(`http://${url}${path}`);
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
};

export const fetchContributionsWithCandidates = async (year?: number) => {
  const url = import.meta.env.VITE_API_URL;
  const path = "/api/contributions/withCandidates";

  try {
    const query = year ? `?year=${year}` : "";
    const response = await axios.get(`http://${url}${path}${query}`);
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
};
