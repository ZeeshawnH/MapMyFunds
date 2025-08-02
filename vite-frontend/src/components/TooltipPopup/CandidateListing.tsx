import type { Candidate } from "../../types/contributions";
import styles from "./TooltipPopup.module.css";

interface CandidateListingProps {
  candidate: Candidate;
}

const CandidateListing = ({ candidate }: CandidateListingProps) => {
  return (
    <li key={candidate.CandidateID} className={styles.recipientItem}>
        <div className={styles.candidateInfo}>
        <span className={styles.candidateName}>{candidate.CandidateName}</span>
        <span className={styles.candidateParty}>{candidate.CandidateParty}</span>
        </div>
        <span className={styles.candidateAmount}>${candidate.NetReceipts.toLocaleString()}</span>
    </li>
  );
};

export default CandidateListing;