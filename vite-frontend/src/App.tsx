import { useEffect, useRef } from 'react'
import * as d3 from 'd3'
import './App.css'

function App() {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    // Clear previous content
    d3.select(svgRef.current).selectAll("*").remove();

    const svg = d3.select(svgRef.current)
      .attr("width", 800)
      .attr("height", 500);

    // Create projection
    const projection = d3.geoAlbersUsa()
      .fitSize([800, 500], { type: "FeatureCollection", features: [] });

    // Create path generator
    const path = d3.geoPath().projection(projection);

    // Load US states GeoJSON from local file
    d3.json("/us-states.json")
      .then((data: any) => {
        if (!data || !data.features) {
          console.error("Failed to load states data");
          return;
        }

        console.log("Loaded states:", data.features.length);

        // Draw states
        svg.selectAll("path")
          .data(data.features)
          .enter()
          .append("path")
          .attr("d", path as any)
          .attr("fill", "#e0e0e0")
          .attr("stroke", "#999")
          .attr("stroke-width", 0.5)
          .on("mouseover", function() {
            d3.select(this)
              .attr("fill", "#b0b0b0")
              .attr("stroke-width", 1);
          })
          .on("mouseout", function() {
            d3.select(this)
              .attr("fill", "#e0e0e0")
              .attr("stroke-width", 0.5);
          })
          .append("title")
          .text((d: any) => d.properties.name);
      })
      .catch((error) => {
        console.error("Error loading map data:", error);
        // Fallback: create a simple rectangle to show the component is working
        svg.append("rect")
          .attr("width", 800)
          .attr("height", 500)
          .attr("fill", "#f0f0f0")
          .attr("stroke", "#ccc");
        
        svg.append("text")
          .attr("x", 400)
          .attr("y", 250)
          .attr("text-anchor", "middle")
          .attr("dominant-baseline", "middle")
          .text("Map loading...");
      });
  }, []);

  return (
    <div className="App">
      <h1>Elections App</h1>
      <h2>US Election Map</h2>
      <svg ref={svgRef} style={{ border: '1px solid #ccc' }}></svg>
    </div>
  )
}

export default App
