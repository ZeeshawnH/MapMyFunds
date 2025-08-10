import type { Candidate } from "../../types/contributions";
import styles from "./TooltipPopup.module.css";
import { partyColor } from "../../utils/color/partyColors";

interface CandidateListingProps {
  candidate: Candidate;
  showCents?: boolean;
}

const CandidateListing = ({
  candidate,
  showCents = true,
}: CandidateListingProps) => {
  const formatAmount = (amount: number) => {
    if (showCents) {
      return amount.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    } else {
      return Math.round(amount).toLocaleString();
    }
  };

  return (
    <li key={candidate.CandidateID} className={styles.recipientItem}>
      <div className={styles.candidateInfo}>
        <span className={styles.candidateName}>{candidate.CandidateName}</span>
        <div className={styles.partyContainer}>
          <span
            className={styles.partyDot}
            style={{ backgroundColor: partyColor(candidate.CandidateParty) }}
          />
          <span className={styles.candidateParty}>
            {candidate.CandidateParty}
          </span>
        </div>
      </div>
      <span className={styles.candidateAmount}>
        ${formatAmount(candidate.NetReceipts)}
      </span>
    </li>
  );
};

export default CandidateListing;
