const toTitleCase = (value: string): string => {
  if (!value) return "";
  return value.charAt(0).toUpperCase() + value.slice(1).toLowerCase();
};

/**
 * Format FEC-style candidate names like "HARRIS, KAMALA D." into
 * a friendlier "Kamala Harris". Falls back to simple title casing
 * when no comma is present.
 */
export const formatCandidateName = (raw: string | undefined | null): string => {
  if (!raw) return "";

  const trimmed = raw.trim();
  if (!trimmed) return "";

  const parts = trimmed.split(",");

  if (parts.length >= 2) {
    const lastRaw = parts[0]?.trim() ?? "";
    const restRaw = parts.slice(1).join(",").trim();

    // Take the first token from the remainder as the given name
    const firstToken = restRaw.split(/\s+/)[0] ?? "";

    const first = toTitleCase(firstToken);
    const last = toTitleCase(lastRaw);

    if (first && last) return `${first} ${last}`;
    if (last) return last;
    return first;
  }

  // No comma: just title-case each token
  const tokens = trimmed.split(/\s+/).map(toTitleCase).filter(Boolean);
  return tokens.join(" ");
};
