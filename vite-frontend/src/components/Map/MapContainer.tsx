import { useEffect, useState } from "react";
import Map from "./Map";
import type { StateContributions } from "../../types";
import { fetchContributionsWithCandidates } from "../../api/fetchContributions";
import { sortContributions } from "../../utils/process/contributions";

const MapContainer = () => {
  const [contributionData, setContributionData] = useState<StateContributions>(
    {} as StateContributions
  );
  const [isLoaded, setLoaded] = useState(false);

  useEffect(() => {
    fetchContributionsWithCandidates()
      .then((data) => {
        setContributionData(sortContributions(data.contributions));
        console.log(data);
      })
      .then(() => setLoaded(true))
      .then(() => console.log("Map data loaded"));
  }, []);

  return isLoaded ? (
    <Map size={850} geojsonPath="/us-states.json" contributionData={contributionData} />
  ) : (
    <div>Loading...</div>
  );
};

export default MapContainer;
