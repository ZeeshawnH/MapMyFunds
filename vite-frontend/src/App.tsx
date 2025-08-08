import MapContainer from "./components/Map/MapContainer";
import "./App.css";

function App() {
  return (
    <div id="root">
      <main>
        <h1>How Each State Voted With Their Money</h1>
        <p className="app-description">
          This map shows the top 10 recipients of campaign contributions in each
          state.
        </p>
        <MapContainer />
      </main>
      <div className="sidebar">

      </div>
    </div>
  );
}

export default App;
