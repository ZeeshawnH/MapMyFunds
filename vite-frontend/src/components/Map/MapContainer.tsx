import { useEffect, useState } from "react";
import Map from "./Map";
import type { StateContributions } from "../../types";
import { fetchContributionsByState } from "../../api/fetchContributions";
import { sortContributions } from "../../utils/process/contributions";

const MapContainer = () => {
  const [mapData, setMapData] = useState<StateContributions>(
    {} as StateContributions
  );
  const [isLoaded, setLoaded] = useState(false);

  useEffect(() => {
    fetchContributionsByState()
      .then((data) => sortContributions(data))
      .then((data) => setMapData(data))
      .then(() => setLoaded(true))
      .then(() => console.log("Map data loaded"));
  }, []);

  return isLoaded ? (
    <Map size={850} geojsonPath="/us-states.json" contributionData={mapData} />
  ) : (
    <div>Loading...</div>
  );
};

export default MapContainer;
