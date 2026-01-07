import { useEffect, useState } from "react";
import MapContainer from "./components/Map/MapContainer";
import SidebarTotals from "./components/Sidebar/SidebarTotals";
import TopDonorsSidebar from "./components/Sidebar/TopDonorsSidebar";
import IntroModal from "./components/IntroModal";
import "./App.css";

function App() {
  const [year, setYear] = useState<number>(2024);
  const [showIntro, setShowIntro] = useState(false);

  useEffect(() => {
    try {
      const seen = window.localStorage.getItem("elections_intro_seen");
      if (!seen) {
        setShowIntro(true);
      }
    } catch {
      setShowIntro(true);
    }
  }, []);

  const handleCloseIntro = () => {
    setShowIntro(false);
    try {
      window.localStorage.setItem("elections_intro_seen", "1");
    } catch {
      // ignore
    }
  };

  return (
    <div id="root">
      <div className="layout-row">
        <div className="left-sidebar">
          <TopDonorsSidebar year={year} />
        </div>
        <main>
          <h1>How Each State Voted With Their Money</h1>
          <p className="app-description">
            This map summarizes the top recipients of presidential campaign
            contributions in each state in 2024. Hover over a state to see which
            candidates received the most contributions.{" "}
            <button
              type="button"
              className="more-info-link"
              onClick={() => setShowIntro(true)}
            >
              More info
            </button>
          </p>
          <MapContainer year={year} onYearChange={setYear} />
        </main>
        <div className="sidebar">
          <SidebarTotals year={year} />
        </div>
      </div>
      {showIntro && <IntroModal onClose={handleCloseIntro} />}
    </div>
  );
}

export default App;
