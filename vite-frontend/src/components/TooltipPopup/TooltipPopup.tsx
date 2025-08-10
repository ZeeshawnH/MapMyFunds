import React from "react";
import styles from "./TooltipPopup.module.css";
import type { Candidate } from "../../types/contributions";
import CandidateListing from "./CandidateListing";

interface TooltipPopupProps {
  x: number;
  y: number;
  stateName: string;
  candidates: Candidate[];
  visible: boolean;
}

const TooltipPopup: React.FC<TooltipPopupProps> = ({
  x,
  y,
  stateName,
  candidates,
  visible,
}) => {
  if (!visible) return null;
  return (
    <div
      className={styles.tooltip}
      style={
        {
          "--x": `${x + 10}px`,
          "--y": `${y + 10}px`,
        } as React.CSSProperties
      }
    >
      <div className={styles.stateName}>{stateName}</div>
      {candidates && candidates.length > 0 ? (
        <ul className={styles.recipientsList}>
          {candidates.map((candidate) => (
            <CandidateListing
              key={candidate.CandidateID}
              candidate={candidate}
              showCents={false}
            />
          ))}
        </ul>
      ) : (
        <div className={styles.noData}>No data</div>
      )}
    </div>
  );
};

export default TooltipPopup;
