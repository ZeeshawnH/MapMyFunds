import { useEffect, useMemo, useState } from "react";
import { fetchContributionsWithCandidates } from "../../api/fetchContributions";
import type { StateContributions } from "../../types";
import { states as stateNames } from "../../utils/constants/stateConstants";
import CandidateListing from "./CandidateListing";
import { candidateImages } from "../../utils/constants/candidateImages";
import { formatCandidateName } from "../../utils/format/candidateName";

interface TopState {
  stateCode: string;
  stateName: string;
  total: number;
}

interface CandidateSummary {
  id: string;
  name: string;
  party: string;
  total: number;
  topStates: TopState[];
}

interface SidebarTotalsProps {
  year: number;
}

const formatCurrency = (value: number) =>
  `$${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`;

export const SidebarTotals = ({ year }: SidebarTotalsProps) => {
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [data, setData] = useState<StateContributions>({});
  const [candidateNamesById, setCandidateNamesById] = useState<
    Record<string, string>
  >({});
  const [candidateImagesById, setCandidateImagesById] = useState<
    Record<string, string>
  >({});

  useEffect(() => {
    let cancelled = false;

    fetchContributionsWithCandidates(year)
      .then((response) => {
        if (!cancelled) {
          setData(response.contributions);

          const byId: Record<string, string> = {};
          const imageById: Record<string, string> = {};
          if (Array.isArray(response.candidates)) {
            for (const cand of response.candidates) {
              if (cand && cand.candidate_id && cand.name) {
                byId[cand.candidate_id] = formatCandidateName(cand.name);
                if (cand.image_url) {
                  imageById[cand.candidate_id] = cand.image_url;
                } else if (candidateImages[cand.candidate_id]) {
                  imageById[cand.candidate_id] =
                    candidateImages[cand.candidate_id];
                }
              }
            }
          }
          setCandidateNamesById(byId);
          setCandidateImagesById(imageById);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setData({});
        }
      });

    return () => {
      cancelled = true;
    };
  }, [year]);

  const candidates: CandidateSummary[] = useMemo(() => {
    const totalsByCandidate: Record<
      string,
      {
        id: string;
        name: string;
        party: string;
        total: number;
        perState: Record<string, number>;
      }
    > = {};

    const national = data["US"] || data["USA"] || [];

    // Prefer using provided national (US) totals for each candidate when available,
    // instead of manually summing across states.
    if (national.length > 0) {
      national.forEach((c) => {
        if (
          c.CandidateName === "All candidates" ||
          c.CandidateName === "Republicans" ||
          c.CandidateName === "Democrats"
        ) {
          return;
        }

        const key = c.CandidateID;
        if (!totalsByCandidate[key]) {
          totalsByCandidate[key] = {
            id: c.CandidateID,
            name: candidateNamesById[c.CandidateID] || c.CandidateName || "",
            party: c.CandidateParty,
            total: 0,
            perState: {},
          };
        }

        totalsByCandidate[key].total += c.NetReceipts;
      });
    }

    // Build per-state breakdowns from individual states only (exclude national bucket).
    Object.entries(data).forEach(([stateCode, contributions]) => {
      if (stateCode === "US" || stateCode === "USA") {
        return;
      }

      contributions.forEach((c) => {
        if (
          c.CandidateName === "All candidates" ||
          c.CandidateName === "Republicans" ||
          c.CandidateName === "Democrats"
        ) {
          return;
        }

        const key = c.CandidateID;
        if (!totalsByCandidate[key]) {
          totalsByCandidate[key] = {
            id: c.CandidateID,
            name: candidateNamesById[c.CandidateID] || c.CandidateName || "",
            party: c.CandidateParty,
            total: 0,
            perState: {},
          };
        }

        if (!totalsByCandidate[key].perState[stateCode]) {
          totalsByCandidate[key].perState[stateCode] = 0;
        }
        totalsByCandidate[key].perState[stateCode] += c.NetReceipts;
      });
    });

    // If no national bucket is present, fall back to summing per-state values
    // to get a total per candidate.
    if (national.length === 0) {
      Object.values(totalsByCandidate).forEach((entry) => {
        entry.total = Object.values(entry.perState).reduce(
          (sum, value) => sum + value,
          0
        );
      });
    }

    const allCandidates = Object.values(totalsByCandidate);

    allCandidates.sort((a, b) => b.total - a.total);

    const topCandidates = allCandidates.slice(0, 12);

    return topCandidates.map((c) => {
      const topStatesForCandidate: TopState[] = Object.entries(c.perState)
        .map(([code, total]) => ({
          stateCode: code,
          stateName: stateNames[code] || code,
          total,
        }))
        // Exclude national aggregate bucket from per-candidate state lists
        .filter(
          (entry) => entry.stateCode !== "US" && entry.stateCode !== "USA"
        )
        .sort((a, b) => b.total - a.total)
        // Surface more states per candidate; list up to 8
        .slice(0, 8);

      return {
        id: c.id,
        name: c.name,
        party: c.party,
        total: c.total,
        topStates: topStatesForCandidate,
      };
    });
  }, [data, candidateNamesById]);

  useEffect(() => {
    if (candidates.length > 0) {
      setExpandedIds(new Set([candidates[0].id]));
    } else {
      setExpandedIds(new Set());
    }
  }, [candidates]);

  if (candidates.length === 0) {
    return (
      <aside className="sidebar-totals">
        <h2 className="sidebar-heading">Top Candidates by Funds Raised</h2>
        <div className="sidebar-empty">No data available for this year.</div>
      </aside>
    );
  }

  return (
    <aside className="sidebar-totals">
      <h2 className="sidebar-heading">Top Candidates by Funds Raised</h2>
      <ul className="sidebar-list">
        {candidates.map((candidate) => {
          const isExpanded = expandedIds.has(candidate.id);

          return (
            <li key={candidate.id} className="sidebar-item">
              <button
                type="button"
                className={`sidebar-item-header ${isExpanded ? "is-expanded" : ""}`}
                onClick={() =>
                  setExpandedIds((current) => {
                    const next = new Set(current);
                    if (next.has(candidate.id)) {
                      next.delete(candidate.id);
                    } else {
                      next.add(candidate.id);
                    }
                    return next;
                  })
                }
              >
                <CandidateListing
                  id={candidate.id}
                  name={candidate.name}
                  party={candidate.party}
                  total={candidate.total}
                  variant="sidebar"
                  imageUrl={candidateImagesById[candidate.id]}
                />
              </button>

              <div
                className={`sidebar-donors ${isExpanded ? "is-expanded" : ""}`}
              >
                <div className="sidebar-donors-label">
                  Top contributing states
                </div>
                <ul className="sidebar-donors-list">
                  {candidate.topStates.map((entry) => (
                    <li key={entry.stateCode} className="sidebar-donor-row">
                      <span className="sidebar-donor-name">
                        {entry.stateName} ({entry.stateCode})
                      </span>
                      <span className="sidebar-donor-total">
                        {formatCurrency(entry.total)}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            </li>
          );
        })}
      </ul>
    </aside>
  );
};

export default SidebarTotals;
