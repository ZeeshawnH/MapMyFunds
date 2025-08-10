import MapContainer from "./components/Map/MapContainer";
import "./App.css";
import SidebarContent from "./components/SidebarContent/SidebarContent";

function App() {
  return (
    <div id="root">
      <div className="sidebar">{/* Left sidebar content */}</div>
      <main>
        <h1>How Each State Voted With Their Money</h1>
        <p className="app-description">
          This map shows the top recipients of campaign contributions in each
          state.
        </p>
        <MapContainer />
      </main>
      <div className="sidebar">
        <SidebarContent />
      </div>
    </div>
  );
}

export default App;
