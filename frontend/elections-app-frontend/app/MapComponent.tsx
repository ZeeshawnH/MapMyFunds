'use client'
import { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import { MapContainer, TileLayer, GeoJSON, Tooltip } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import axios from "axios";

// Dynamically import leaflet without server-side rendering
const Map = dynamic(() => import('react-leaflet').then(mod => mod.MapContainer), { ssr: false });

// Example GeoJSON data for US states
const USStatesGeoJSON = "@/public/us-states.json"; // You should replace this with your GeoJSON data

const MapComponent = () => {
  const [stateData, setStateData] = useState<{ [key: string]: number }>({});

  useEffect(() => {
    // Fetch contribution data
    axios.get("/api/contributions").then(response => {
      // Assuming the response contains { state: amount }
      const contributions = response.data;
      setStateData(contributions);
    }).catch(error => {
      console.error("Error fetching contribution data", error);
    });
  }, []);

  return (
    <MapContainer center={[37.8, -96]} zoom={4} style={{ height: "100vh", width: "100%" }}>
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      />
      <GeoJSON
        data={USStatesGeoJSON}
        onEachFeature={(feature, layer) => {
          const stateName = feature.properties.name; // Replace with appropriate property
          const amount = stateData[stateName] || "No data";
          layer.bindTooltip(`${stateName}: $${amount.toLocaleString()}`);
        }}
      />
    </MapContainer>
  );
};

export default MapComponent;
