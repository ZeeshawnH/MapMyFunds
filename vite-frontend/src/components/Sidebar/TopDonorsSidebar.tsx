import { useEffect, useMemo, useState } from "react";
import CandidateListing from "./CandidateListing";
import { fetchContributionsWithCandidates } from "../../api/fetchContributions";
import type { StateContributions } from "../../types";
import { states as stateNames } from "../../utils/constants/stateConstants";
import { formatCandidateName } from "../../utils/format/candidateName";

interface StateCandidate {
  id: string;
  name: string;
  party: string;
  total: number;
}

interface StateSummary {
  code: string;
  name: string;
  total: number;
  candidates: StateCandidate[];
}

interface TopDonorsSidebarProps {
  year: number;
}

const formatCurrency = (value: number) =>
  `$${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`;

export const TopDonorsSidebar = ({ year }: TopDonorsSidebarProps) => {
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [data, setData] = useState<StateContributions>({});
  const [candidateNamesById, setCandidateNamesById] = useState<
    Record<string, string>
  >({});

  useEffect(() => {
    let cancelled = false;

    fetchContributionsWithCandidates(year)
      .then((response) => {
        if (!cancelled) {
          setData(response.contributions);

          const byId: Record<string, string> = {};
          if (Array.isArray(response.candidates)) {
            for (const cand of response.candidates) {
              if (cand && cand.candidate_id && cand.name) {
                byId[cand.candidate_id] = formatCandidateName(cand.name);
              }
            }
          }
          setCandidateNamesById(byId);
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

  const states: StateSummary[] = useMemo(() => {
    const stateTotals: Record<
      string,
      {
        code: string;
        name: string;
        total: number;
        candidates: Record<string, StateCandidate>;
      }
    > = {};

    Object.entries(data).forEach(([stateCode, contributions]) => {
      if (!stateTotals[stateCode]) {
        stateTotals[stateCode] = {
          code: stateCode,
          name: stateNames[stateCode] || stateCode,
          total: 0,
          candidates: {},
        };
      }

      let aggregateTotal: number | undefined;
      let candidateTotal = 0;

      contributions.forEach((c) => {
        if (c.CandidateName === "All candidates") {
          aggregateTotal = c.NetReceipts;
          return;
        }
        if (
          c.CandidateName === "Republicans" ||
          c.CandidateName === "Democrats"
        ) {
          // Ignore party aggregates in the per-candidate breakdown
          return;
        }

        candidateTotal += c.NetReceipts;

        const candidateId = c.CandidateID;
        if (!stateTotals[stateCode].candidates[candidateId]) {
          stateTotals[stateCode].candidates[candidateId] = {
            id: candidateId,
            name: candidateNamesById[c.CandidateID] || c.CandidateName || "",
            party: c.CandidateParty,
            total: 0,
          };
        }

        stateTotals[stateCode].candidates[candidateId].total += c.NetReceipts;
      });

      stateTotals[stateCode].total =
        aggregateTotal !== undefined ? aggregateTotal : candidateTotal;
    });

    const allStates = Object.values(stateTotals).filter(
      (s) => s.code !== "US" && s.code !== "USA"
    );

    allStates.sort((a, b) => b.total - a.total);

    // Allow more states to surface; include up to 15
    const topStates = allStates.slice(0, 15);

    return topStates.map((s) => ({
      code: s.code,
      name: s.name,
      total: s.total,
      candidates: Object.values(s.candidates)
        .sort((a, b) => b.total - a.total)
        // Show more candidates per state
        .slice(0, 12),
    }));
  }, [data, candidateNamesById]);

  useEffect(() => {
    if (states.length > 0) {
      setExpandedIds(new Set([states[0].code]));
    } else {
      setExpandedIds(new Set());
    }
  }, [states]);

  if (states.length === 0) {
    return (
      <aside className="sidepanel sidepanel-donors">
        <h2 className="sidepanel-heading">Top Contributing States</h2>
        <div className="sidepanel-empty">No data available for this year.</div>
      </aside>
    );
  }

  return (
    <aside className="sidepanel sidepanel-donors">
      <h2 className="sidepanel-heading">Top Contributing States</h2>
      <ul className="sidepanel-list">
        {states.map((state) => {
          const isExpanded = expandedIds.has(state.code);

          return (
            <li key={state.code} className="sidepanel-item">
              <button
                type="button"
                className={`sidepanel-item-header ${
                  isExpanded ? "is-expanded" : ""
                }`}
                onClick={() =>
                  setExpandedIds((current) => {
                    const next = new Set(current);
                    if (next.has(state.code)) {
                      next.delete(state.code);
                    } else {
                      next.add(state.code);
                    }
                    return next;
                  })
                }
              >
                <div className="sidepanel-item-primary">
                  <span className="sidepanel-item-name">
                    {state.code === "US" || state.code === "USA"
                      ? "United States (overall)"
                      : `${state.name} (${state.code})`}
                  </span>
                </div>
                <div className="sidepanel-item-total">
                  <span
                    className={
                      state.code === "US" || state.code === "USA"
                        ? "sidepanel-item-total-us"
                        : undefined
                    }
                  >
                    {formatCurrency(state.total)}
                  </span>
                </div>
              </button>

              <div
                className={`sidepanel-detail ${
                  isExpanded ? "is-expanded" : ""
                }`}
              >
                <div className="sidepanel-detail-label">
                  {state.code === "US" || state.code === "USA"
                    ? "Top recipient campaigns nationally"
                    : "Top recipient campaigns in this state"}
                </div>
                <ul className="sidepanel-detail-list">
                  {state.candidates.map((cand) => (
                    <li key={cand.id} className="sidepanel-detail-row">
                      <CandidateListing
                        id={cand.id}
                        name={cand.name}
                        party={cand.party}
                        total={cand.total}
                        variant="sidepanel"
                      />
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

export default TopDonorsSidebar;
