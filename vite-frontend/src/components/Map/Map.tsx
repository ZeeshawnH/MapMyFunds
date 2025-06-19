import React, { useEffect, useRef, useState, useLayoutEffect } from "react";
import * as d3 from "d3";
import type { Feature, FeatureCollection, Geometry } from "geojson";
import type { GeoPermissibleObjects } from "d3";
import type { TooltipData, GeoJsonFeature } from "../../types";
import {
  handleMouseHover,
  handleMouseMove,
  handleMouseOut,
} from "../../utils/d3";
import styles from "./Map.module.css";
import { mockContributions } from "../../mocks/mockContributions";
import { partyColor, numericToStateCode } from "../../utils/constants";

const MAP_SIZE = 850; // Change this value to control the map's size
const geojsonPath: string = "/us-states.json";
const contributionData = Object.fromEntries(
  Object.entries(mockContributions).map(([state, contributions]) => [
    state,
    contributions.sort((a, b) => b.NetReceipts - a.NetReceipts),
  ])
);

const Map: React.FC = () => {
  // svg for map
  const svgRef = useRef<SVGSVGElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  // Mouse tooltip
  const [tooltip, setTooltip] = useState<TooltipData>(null);
  const [dimensions, setDimensions] = useState({
    width: MAP_SIZE,
    height: MAP_SIZE * 0.67,
  });

  // Get container size (now just set to MAP_SIZE)
  useLayoutEffect(() => {
    setDimensions({ width: MAP_SIZE, height: MAP_SIZE * 0.67 });
  }, []);

  // Render map
  useEffect(() => {
    if (!svgRef.current) return;
    const svg = d3.select(svgRef.current);
    const width = dimensions.width;
    const height = dimensions.height;

    svg.selectAll("path").remove();
    svg.attr("width", width).attr("height", height);

    d3.json(geojsonPath).then((data) => {
      const geoData = data as FeatureCollection<
        Geometry,
        Record<string, unknown>
      >;
      if (!geoData || !geoData.features) return;

      const projection = d3.geoAlbersUsa().fitSize([width, height], geoData);
      const path = d3.geoPath().projection(projection);

      svg
        .selectAll("path")
        .data(geoData.features)
        .enter()
        .append("path")
        .attr("d", (d) => path(d as GeoPermissibleObjects) || "")
        .attr("class", styles.statePath)
        .attr("style", (d) => {
          const feature = d as GeoJsonFeature;
          const numericStateCode = feature.properties.STATE;
          const alphabeticStateCode = numericToStateCode[numericStateCode];
          const stateContributions = contributionData[alphabeticStateCode];

          if (stateContributions && stateContributions.length > 0) {
            const topRecipient = stateContributions[0];
            const color = partyColor(topRecipient.CandidateParty);
            return `--state-fill: ${color}`;
          }
          const defaultColor = partyColor("");
          return `--state-fill: ${defaultColor}`;
        })
        .on(
          "mouseover",
          function (event, d: Feature<Geometry, Record<string, unknown>>) {
            handleMouseHover(this, event, d, svg, setTooltip);
          }
        )
        .on(
          "mousemove",
          function (event, d: Feature<Geometry, Record<string, unknown>>) {
            handleMouseMove(event, d, svg, setTooltip);
          }
        )
        .on("mouseout", function () {
          handleMouseOut(this, setTooltip);
        });
    });
  }, [dimensions]);

  return (
    <div
      className={styles.mapContainer}
      ref={containerRef}
      style={{ width: MAP_SIZE, height: MAP_SIZE * 0.67 }}
    >
      <svg ref={svgRef} className={styles.mapSvg}></svg>
      {tooltip && (
        <div
          className={styles.mapTooltip}
          style={{ left: tooltip.x + 10, top: tooltip.y + 10 }}
        >
          {tooltip.name}
        </div>
      )}
    </div>
  );
};

export default Map;
