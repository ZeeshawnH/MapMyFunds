import { partyColor } from "../../utils/color/partyColors";

export interface CandidateListingProps {
  id?: string;
  name: string;
  party?: string;
  total: number;
  /**
   * sidebar: top-level sidebar rows (room for avatar)
   * sidepanel: nested dropdown rows (no avatar, more compact)
   */
  variant?: "sidebar" | "sidepanel";
  /** Optional URL for a candidate image when rendered as a top-level item. */
  imageUrl?: string;
}

const formatCurrency = (value: number) =>
  `$${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`;

export const CandidateListing = ({
  name,
  party,
  total,
  variant = "sidepanel",
  imageUrl,
}: CandidateListingProps) => {
  const isSidebar = variant === "sidebar";
  const rootClass = isSidebar
    ? "candidate-listing candidate-listing--sidebar"
    : "candidate-listing candidate-listing--sidepanel";

  const color = party ? partyColor(party) : undefined;

  return (
    <div
      className={rootClass}
      style={color ? { borderLeftColor: color } : undefined}
    >
      {isSidebar && imageUrl && (
        <div className="candidate-listing-avatar">
          <img src={imageUrl} alt={name} />
        </div>
      )}
      <div className="candidate-listing-main">
        <div className="candidate-listing-text">
          <span className="candidate-listing-name">{name}</span>
          {party && <span className="candidate-listing-party">{party}</span>}
        </div>
        <span className="candidate-listing-total">{formatCurrency(total)}</span>
      </div>
    </div>
  );
};

export default CandidateListing;
