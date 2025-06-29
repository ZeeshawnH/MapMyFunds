import axios from "axios";

export const fetchContributionsByState = async () => {
    const url = import.meta.env.VITE_API_URL;
    const path = import.meta.env.VITE_CONTRIBUTIONS_ENDPOINT

    try {
        const response = await axios.get(`https://${url}${path}`);
        return response.data;
    } catch (error) {
        console.error(error);
        throw error;
    }
}