"use client";

import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";
import styles from "./MapComponent.module.css"; // Import CSS module for styling
import { fetchContributions } from "@/public/api/contributions";

interface GeoJsonFeature {
  type: string;
  properties: {
    NAME: string;
  };
  geometry: {
    type: string;
    coordinates: any;
  };
}

interface GeoJsonData {
  type: string;
  features: GeoJsonFeature[];
}

const MapComponent = () => {
  const states: {} = {
    "01": "AL",
    "02": "AK",
    "04": "AZ",
    "05": "AR",
    "06": "CA",
    "08": "CO",
    "09": "CT",
    "10": "DE",
    "12": "FL",
    "13": "GA",
    "15": "HI",
    "16": "ID",
    "17": "IL",
    "18": "IN",
    "19": "IA",
    "20": "KS",
    "21": "KY",
    "22": "LA",
    "23": "ME",
    "24": "MD",
    "25": "MA",
    "26": "MI",
    "27": "MN",
    "28": "MS",
    "29": "MO",
    "30": "MT",
    "31": "NE",
    "32": "NV",
    "33": "NH",
    "34": "NJ",
    "35": "NM",
    "36": "NY",
    "37": "NC",
    "38": "ND",
    "39": "OH",
    "40": "OK",
    "41": "OR",
    "42": "PA",
    "44": "RI",
    "45": "SC",
    "46": "SD",
    "47": "TN",
    "48": "TX",
    "49": "UT",
    "50": "VT",
    "51": "VA",
    "53": "WA",
    "54": "WV",
    "55": "WI",
    "56": "WY",
  };

  const [contributions, setContributions] = useState<
    {
      candidate_id: { type: String; required: true; unique: true };
      candidate_last_name: String;
      candidate_party_affiliation: String;
      contributions: [
        {
          contributor_state: String;
          election_year: Number;
          net_receipts: Number;
          rounded_net_receipts: Number;
        }
      ];
    }[]
  >([]);

  useEffect(() => {
    // Call the fetch function and handle the results
    const fetchAndSetContributions = async () => {
      try {
        const results = await fetchContributions();
        setContributions(results);
      } catch (error) {
        console.error("Error fetching contributions:", error);
      }
    };

    fetchAndSetContributions();
  }, []);

  const svgRef = useRef<SVGSVGElement | null>(null); // Reference to the SVG element
  const [popup, setPopup] = useState<{
    x: number;
    y: number;
    content: string;
  } | null>(null); // State for popup content

  useEffect(() => {
    if (contributions.length === 0) return;

    console.log(contributions);

    const svg = d3.select(svgRef.current);
    const width = svg.node()!.clientWidth;
    const height = svg.node()!.clientHeight;

    const projection = d3.geoAlbersUsa().translate([width / 2, height / 2]);
    const path = d3.geoPath().projection(projection);

    d3.json("/us-states.json").then((data: any) => {
      svg
        .selectAll("path")
        .data(data.features)
        .enter()
        .append("path")
        .attr("d", (d) => path(d as d3.GeoPermissibleObjects))
        .attr("fill", "#cccccc")
        .attr("stroke", "#000000")
        .attr("stroke-width", 1)
        .on("mouseover", (event, d) => {
          const feature = d as GeoJsonFeature;
          const [x, y] = path.centroid(feature as d3.GeoPermissibleObjects);
          console.log(feature.properties);

          const state: number = Number(feature.properties.STATE);
          console.log(state);
          const receipts = contributions;
          for (const receipt of receipts) {
            receipt.contributions.filter(
              (contributor_state) =>
                contributor_state === states[feature.properties.STATE]
            );
          }
          console.log(receipts);

          setPopup({
            x: x,
            y: y,
            content: receipts[0].contributions[0].net_receipts,
          });
          d3.select(event.currentTarget).attr("fill", "#ff9999");
        })
        .on("mouseout", (event) => {
          setPopup(null);
          d3.select(event.currentTarget).attr("fill", "#cccccc");
        });
    });
  });

  return (
    <div className={styles.mapContainer}>
      <svg ref={svgRef} width="100%" height="600" className={styles.mapSvg}>
        {/* SVG content will be added by D3 */}
      </svg>
      {popup && (
        <div className={styles.popup} style={{ left: popup.x, top: popup.y }}>
          {popup.content}
        </div>
      )}
    </div>
  );
};

export default MapComponent;
