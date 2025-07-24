import MapContainer from "./components/Map/MapContainer";

function App() {
  return (
    <div id="root">
      <h1>Electoral Contributions Map</h1>
      <p className="app-description">
        This map shows the top 10 recipients of campaign contributions in each
        state.
      </p>
      <MapContainer />
    </div>
  );
}

export default App;
