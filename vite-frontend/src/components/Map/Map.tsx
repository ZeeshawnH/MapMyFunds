import React, { useEffect, useRef, useState, useLayoutEffect } from "react";
import * as d3 from "d3";
import type { Feature, FeatureCollection, Geometry } from "geojson";
import type { GeoPermissibleObjects } from "d3";
import type { TooltipData, GeoJsonFeature } from "../../types";
import type { StateContributions } from "../../types/contributions";
import {
  handleMouseHover,
  handleMouseMove,
  handleMouseOut,
} from "../../utils/d3";
import styles from "./Map.module.css";
import { partyColor, numericToStateCode } from "../../utils/constants";
import TooltipPopup from "../TooltipPopup/TooltipPopup";

interface MapProps {
  size: number;
  geojsonPath: string;
  contributionData: StateContributions;
}

const Map: React.FC<MapProps> = ({ size, geojsonPath, contributionData }) => {
  // svg for map
  const svgRef = useRef<SVGSVGElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  // Mouse tooltip
  const [tooltip, setTooltip] = useState<TooltipData>(null);
  const [dimensions, setDimensions] = useState({
    width: size,
    height: size * 0.67,
  });

  // Get container size (now just set to MAP_SIZE)
  useLayoutEffect(() => {
    setDimensions({ width: size, height: size * 0.67 });
  }, [size]);

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
  }, [dimensions, geojsonPath, contributionData]);

  return (
    <div
      className={styles["map-container"]}
      ref={containerRef}
      style={{ width: size, height: size * 0.67 }}
    >
      <svg ref={svgRef} className={styles.mapSvg}></svg>
      <TooltipPopup
        x={tooltip?.x ?? 0}
        y={tooltip?.y ?? 0}
        stateName={tooltip?.name ?? ""}
        candidates={
          tooltip?.stateCode && contributionData[tooltip.stateCode]
            ? contributionData[tooltip.stateCode]
            : []
        }
        visible={!!tooltip}
      />
    </div>
  );
};

export default Map;
