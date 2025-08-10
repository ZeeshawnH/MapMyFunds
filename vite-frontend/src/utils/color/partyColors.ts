export const partyColor = (
  party: string | null,
  isHovered: boolean = false
): string => {
  if (!party) {
    return isHovered ? "#9ca3af" : "#6b7280"; // Neutral gray
  }

  switch (party.toUpperCase()) {
    case "DEM":
      return isHovered ? "#4d94ff" : "#1a75ff";
    case "REP":
      return isHovered ? "#ff4d4d" : "#ff1a1a";
    case "IND":
      return isHovered ? "#ffd700" : "#ffcc00"; // Gold for independent
    default:
      return isHovered ? "#e6e6e6" : "#cccccc";
  }
};
