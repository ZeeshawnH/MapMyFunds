import axios from "axios";

export const fetchContributionsByState = async () => {
  const url = "api.zeeshawnh.com";
  const path = "/api/contributions";

  try {
    const response = await axios.get(`https://${url}${path}`);
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
};

export const fetchContributionsWithCandidates = async () => {
  const url = "localhost:8080";
  const path = "/api/contributions/withCandidates";

  try {
    const response = await axios.get(`http://${url}${path}`);
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
};
