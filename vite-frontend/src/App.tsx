import { useEffect, useState } from "react";
import Map from "./components/Map";
import type { StateContributions } from "./types";
import { fetchContributionsByState } from "./api/fetchContributions";
import { sortContributions } from "./utils/process/contributions";

function App() {
  const [contributions, setContributions] = useState({} as StateContributions)

  useEffect(() => {
    fetchContributionsByState()
      .then((data) => sortContributions(data))
      .then((data) => setContributions(data))
  }, []);

  return (
    <div id="root">
      <h1>Electoral Contributions Map</h1>
      <Map
        size={850}
        geojsonPath="/us-states.json"
        contributionData={contributions}
      />
    </div>
  );
}

export default App;
