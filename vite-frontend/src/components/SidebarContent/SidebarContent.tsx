import CandidateListing from "../TooltipPopup/CandidateListing";
import { mockCandidateTotals } from "../../mocks/mockCandidateTotals";
import styles from "./SidebarContent.module.css";

const SidebarContent = () => {
  return (
    <div className={styles.sidebarContent}>
      <h2 className={styles.title}>Top Recipients</h2>
      <ul className={styles.list}>
        {mockCandidateTotals.map((candidate) => (
          <CandidateListing
            key={candidate.CandidateID}
            candidate={candidate}
            showCents={false}
          />
        ))}
      </ul>
    </div>
  );
};

export default SidebarContent;
