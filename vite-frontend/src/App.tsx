import Map from "./components/Map";
import { mockContributions } from "./mocks/mockContributions";

function App() {
  return (
    <div id="root">
      <h1>Electoral Contributions Map</h1>
      <Map
        size={850}
        geojsonPath="/us-states.json"
        contributionData={mockContributions}
      />
    </div>
  );
}

export default App;
