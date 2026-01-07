import { useEffect, useMemo, useState } from "react";
import USMap from "./Map";
import type { StateContributions } from "../../types";
import { fetchContributionsWithCandidates } from "../../api/fetchContributions";
import { sortContributions } from "../../utils/process/contributions";
import { partyColor } from "../../utils/color/partyColors";
import { formatCandidateName } from "../../utils/format/candidateName";

interface MapContainerProps {
  year: number;
  onYearChange: (year: number) => void;
}

const MapContainer = ({ year, onYearChange }: MapContainerProps) => {
  const [contributionData, setContributionData] = useState<StateContributions>(
    {} as StateContributions
  );
  const [isLoaded, setLoaded] = useState(false);
  const [mapSize, setMapSize] = useState(850);
  const [candidateNamesById, setCandidateNamesById] = useState<
    Record<string, string>
  >({});
  const [hoveredSegment, setHoveredSegment] = useState<null | {
    label: string;
    party: string;
    total: number;
    isOther: boolean;
    centerPercent: number;
  }>(null);

  // Keep the map sized so that sidebars + map fit within typical viewport widths
  // by shrinking the map on narrower screens.
  useEffect(() => {
    const recomputeSize = () => {
      const sidebarsWidth = 560; // two 280px sidebars
      const paddingAndGaps = 120; // approximate padding + gutter space
      const available = window.innerWidth - sidebarsWidth - paddingAndGaps;
      const clamped = Math.max(520, Math.min(850, available));
      setMapSize(Number.isFinite(clamped) ? clamped : 850);
    };

    recomputeSize();
    window.addEventListener("resize", recomputeSize);
    return () => window.removeEventListener("resize", recomputeSize);
  }, []);

  useEffect(() => {
    setLoaded(false);
    fetchContributionsWithCandidates(year)
      .then((data) => {
        setContributionData(sortContributions(data.contributions));
        console.log(data);
        const byId: Record<string, string> = {};
        if (Array.isArray(data.candidates)) {
          for (const cand of data.candidates) {
            if (cand && cand.candidate_id && cand.name) {
              byId[cand.candidate_id] = formatCandidateName(cand.name);
            }
          }
        }
        setCandidateNamesById(byId);
      })
      .then(() => setLoaded(true))
      .then(() => console.log("Map data loaded"));
  }, [year]);

  const summary = useMemo(() => {
    if (!isLoaded) return null;
    // Prefer pre-aggregated national totals (state code US/USA) when present.
    const national =
      contributionData["US"] || contributionData["USA"] || undefined;

    const source =
      national && national.length > 0
        ? national
        : Object.values(contributionData).flat();

    const rows = source.filter(
      (c) =>
        c.CandidateName !== "All candidates" &&
        c.CandidateName !== "Republicans" &&
        c.CandidateName !== "Democrats"
    );

    if (rows.length === 0) return null;

    // When using national rows, each candidate should already be an
    // aggregated total, but we still defensively group by ID.
    const totalsByCandidate = new Map<
      string,
      { name: string; party: string; total: number }
    >();

    for (const c of rows) {
      const displayName =
        candidateNamesById[c.CandidateID] || c.CandidateName || "";
      const current =
        totalsByCandidate.get(c.CandidateID) ??
        ({ name: displayName, party: c.CandidateParty, total: 0 } as const);
      totalsByCandidate.set(c.CandidateID, {
        ...current,
        total: current.total + c.NetReceipts,
      });
    }

    const candidateEntries = Array.from(totalsByCandidate.entries()).sort(
      (a, b) => b[1].total - a[1].total
    );

    const topN = 6;
    const topCandidates = candidateEntries.slice(0, topN);
    const otherCandidates = candidateEntries.slice(topN);

    const otherTotal = otherCandidates.reduce(
      (sum, [, info]) => sum + info.total,
      0
    );

    const grandTotal = candidateEntries.reduce(
      (sum, [, info]) => sum + info.total,
      0
    );

    const barSegments = [
      ...topCandidates.map(([, info]) => ({
        label: info.name,
        party: info.party,
        total: info.total,
        isOther: false,
      })),
      ...(otherTotal > 0
        ? [
            {
              label: "Other candidates",
              party: "OTHER",
              total: otherTotal,
              isOther: true,
            },
          ]
        : []),
    ];

    const topCandidate = candidateEntries[0]?.[1];

    return { topCandidate, barSegments, grandTotal };
  }, [contributionData, isLoaded, candidateNamesById]);

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <div className="map-layout">
      <div className="year-selector">
        {[2016, 2020, 2024].map((y) => (
          <button
            key={y}
            type="button"
            className={"year-pill" + (year === y ? " is-active" : "")}
            onClick={() => onYearChange(y)}
          >
            {y}
          </button>
        ))}
      </div>
      {summary && (
        <div className="summary-bar">
          <div className="summary-bar-label">
            <strong>Total raised by candidate</strong>
          </div>
          {hoveredSegment && (
            <div
              className="summary-bar-tooltip"
              style={{ left: `${hoveredSegment.centerPercent}%` }}
            >
              {hoveredSegment.label}
              {!hoveredSegment.isOther && hoveredSegment.party
                ? ` (${hoveredSegment.party})`
                : ""}
              {` \u00b7 $${hoveredSegment.total.toLocaleString()}`}
            </div>
          )}
          <div className="summary-bar-track" aria-hidden="true">
            {(() => {
              let accumulatedPercent = 0;
              return summary.barSegments.map((segment) => {
                const fraction = summary.grandTotal
                  ? (segment.total / summary.grandTotal) * 100
                  : 0;
                const startPercent = accumulatedPercent;
                accumulatedPercent += fraction;
                const centerPercent = startPercent + fraction / 2;

                const color = segment.isOther
                  ? "#d1d5db"
                  : partyColor(segment.party || "Other");

                const isHovered = hoveredSegment?.label === segment.label;

                return (
                  <div
                    key={segment.label}
                    className="summary-bar-segment"
                    style={{
                      width: `${fraction}%`,
                      backgroundColor: color,
                      opacity: hoveredSegment ? (isHovered ? 1 : 0.4) : 1,
                      transition: "opacity 120ms ease-out",
                    }}
                    onMouseEnter={() =>
                      setHoveredSegment({ ...segment, centerPercent })
                    }
                    onMouseLeave={() => setHoveredSegment(null)}
                  />
                );
              });
            })()}
          </div>
        </div>
      )}

      <USMap
        size={mapSize}
        geojsonPath="/us-states.json"
        contributionData={contributionData}
        candidateNamesById={candidateNamesById}
      />

      {summary && (
        <div className="summary-bar">
          <div className="summary-bar-bottom-track" aria-hidden="true" />
          <div className="summary-bar-bottom-text">
            All data is sourced from the U.S. Federal Election Commission (FEC)
            via the OpenFEC API. Visualizations reflect reported aggregate
            campaign finance totals by candidate committee and election cycle,
            as published by the FEC. Figures represent committee-level activity
            and may not fully attribute funds directly to individual candidates
            or donors. Data completeness and granularity are limited to the
            underlying FEC reporting and API endpoints used.
          </div>
        </div>
      )}
    </div>
  );
};

export default MapContainer;
