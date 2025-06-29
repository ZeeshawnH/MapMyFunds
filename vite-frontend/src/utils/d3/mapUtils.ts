import type { Feature, Geometry } from "geojson";
import type { Dispatch, SetStateAction } from "react";
import type { TooltipData } from "../../types/map";
import * as d3 from "d3";
import { numericToStateCode } from "../constants/stateConstants";

/**
 * Helper to get state name from GeoJSON properties
 * @param properties GeoJSON properties
 * @returns state name as string
 */
export function getStateName(properties: unknown): string {
  if (
    properties &&
    typeof properties === "object" &&
    "NAME" in properties &&
    typeof (properties as { NAME: unknown }).NAME === "string"
  ) {
    return (properties as { NAME: string }).NAME;
  }
  return "";
}

/**
 * Highlights state when hovering over
 * @param element map element
 * @param event Mouse Event
 * @param d state being hovered over
 * @param svg d3 svg
 * @param onSetTooltip Tooltip setter
 */
export const handleMouseHover = (
  element: SVGPathElement,
  event: MouseEvent,
  d: Feature<Geometry, Record<string, unknown>>,
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  onSetTooltip: Dispatch<SetStateAction<TooltipData>>
) => {
  d3.select(element);
  const [x, y] = d3.pointer(event, svg.node());
  const stateCode =
    d.properties && 
    d.properties.STATE &&
    typeof d.properties.STATE !== "undefined"
      ? numericToStateCode[d.properties.STATE as string]
      : "";
  onSetTooltip({
    x,
    y,
    name: getStateName(d.properties),
    stateCode,
  });
};

/**
 * Handles mouse moving over state
 * @param event Mouse Event
 * @param d state being hovered over
 * @param svg d3 svg
 * @param onSetTooltip Tooltip setter
 */
export const handleMouseMove = (
  event: MouseEvent,
  d: Feature<Geometry, Record<string, unknown>>,
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  onSetTooltip: Dispatch<SetStateAction<TooltipData>>
) => {
  const [x, y] = d3.pointer(event, svg.node());
  const stateCode =
    d.properties && 
    d.properties.STATE &&
    typeof d.properties.STATE !== "undefined"
      ? numericToStateCode[d.properties.STATE as string]
      : "";
  onSetTooltip({
    x,
    y,
    name: getStateName(d.properties),
    stateCode,
  });
};

/**
 * Handle mouse leaving hover state and unhighlighting state
 * @param element map element
 * @param onSetTooltip Tooltip setter
 */
export const handleMouseOut = (
  element: SVGPathElement,
  onSetTooltip: Dispatch<SetStateAction<TooltipData>>
) => {
  d3.select(element).attr("fill", "#e0e0e0").attr("stroke-width", 0.5);
  onSetTooltip(null);
};
