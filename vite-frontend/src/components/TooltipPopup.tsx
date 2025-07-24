import React from "react";
import styles from "./TooltipPopup.module.css";

interface Recipient {
  CandidateID: string;
  CandidateName: string;
  CandidateParty: string;
  NetReceipts: number;
}

interface TooltipPopupProps {
  x: number;
  y: number;
  stateName: string;
  recipients: Recipient[];
  visible: boolean;
}

const TooltipPopup: React.FC<TooltipPopupProps> = ({
  x,
  y,
  stateName,
  recipients,
  visible,
}) => {
  if (!visible) return null;
  return (
    <div
      className={styles.tooltip}
      style={{
        "--x": `${x + 10}px`,
        "--y": `${y + 10}px`,
      } as React.CSSProperties}
    >
      <div className={styles.stateName}>{stateName}</div>
      {recipients && recipients.length > 0 ? (
        <ul className={styles.recipientsList}>
          {recipients.map((recipient) => (
            <li key={recipient.CandidateID} className={styles.recipientItem}>
              <span className={styles.candidateName}>{recipient.CandidateName}</span>
              {` (${recipient.CandidateParty})`}: $
              {recipient.NetReceipts.toLocaleString()}
            </li>
          ))}
        </ul>
      ) : (
        <div className={styles.noData}>No data</div>
      )}
    </div>
  );
};

export default TooltipPopup;
