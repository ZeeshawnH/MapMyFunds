import axios from "axios";

// Function to fetch contributions data
export const fetchContributions = async () => {
  try {
    const response = await axios.get("http://localhost:8080/api/candidates");
    return response.data;
  } catch (error) {
    console.error("Error fetching data:", error);
    throw error;
  }
};
