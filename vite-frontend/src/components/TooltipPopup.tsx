import React from "react";

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
      style={{
        position: "absolute",
        left: x + 10,
        top: y + 10,
        width: 220,
        minHeight: 80,
        background: "white",
        border: "1px solid #ccc",
        padding: "8px 12px",
        pointerEvents: "none",
        fontSize: 14,
        borderRadius: 4,
        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
        zIndex: 10,
        display: "flex",
        flexDirection: "column",
        justifyContent: "flex-start",
      }}
    >
      <div style={{ fontWeight: 700, marginBottom: 6 }}>{stateName}</div>
      {recipients && recipients.length > 0 ? (
        <ul style={{ margin: 0, padding: 0, listStyle: "none" }}>
          {recipients.map((recipient) => (
            <li key={recipient.CandidateID} style={{ marginBottom: 4 }}>
              <span style={{ fontWeight: 500 }}>{recipient.CandidateName}</span>
              {` (${recipient.CandidateParty})`}: $
              {recipient.NetReceipts.toLocaleString()}
            </li>
          ))}
        </ul>
      ) : (
        <div>No data</div>
      )}
    </div>
  );
};

export default TooltipPopup;
