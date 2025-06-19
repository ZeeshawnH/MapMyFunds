export interface GeoJsonFeature {
  type: string;
  properties: {
    NAME: string;
    STATE: string;
  };
  geometry: {
    type: string;
    coordinates: any;
  };
}

export interface GeoJsonData {
  type: string;
  features: GeoJsonFeature[];
}

export interface GeoJsonFeatureWithArea extends GeoJsonFeature {
  area: number;
}

export interface StateCodeMap {
  [key: string]: string;
}
