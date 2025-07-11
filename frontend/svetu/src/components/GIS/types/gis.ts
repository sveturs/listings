export interface GeoLocation {
  latitude: number;
  longitude: number;
  accuracy?: number;
  altitude?: number;
  altitudeAccuracy?: number;
  heading?: number;
  speed?: number;
}

export interface MapViewState {
  longitude: number;
  latitude: number;
  zoom: number;
  pitch?: number;
  bearing?: number;
}

export interface MapMarkerData {
  id: string;
  position: [number, number]; // [longitude, latitude]
  title: string;
  description?: string;
  type: 'listing' | 'user' | 'poi';
  data?: any;
}

export interface MapPopupData {
  id: string;
  position: [number, number];
  title: string;
  description?: string;
  content?: React.ReactNode;
  onClose?: () => void;
}

export interface GeoSearchResult {
  id: string;
  display_name: string;
  lat: string;
  lon: string;
  boundingbox: [string, string, string, string];
  type: string;
  class: string;
  importance: number;
}

export interface GeoSearchParams {
  query: string;
  limit?: number;
  countrycodes?: string;
  bounded?: boolean;
  viewbox?: string;
  language?: string;
}

export interface DistanceCalculation {
  distance: number;
  unit: 'meters' | 'kilometers' | 'miles';
  bearing: number;
}

export interface MapBounds {
  north: number;
  south: number;
  east: number;
  west: number;
}

export interface MapControlsConfig {
  showZoom?: boolean;
  showCompass?: boolean;
  showFullscreen?: boolean;
  showGeolocate?: boolean;
  showNavigation?: boolean;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
}

export interface MapStyle {
  name: string;
  url: string;
  thumbnail?: string;
}

export interface MapTheme {
  light: MapStyle;
  dark: MapStyle;
  satellite: MapStyle;
}

// Типы для работы с Mapbox GL
export interface MapboxEvent {
  type: string;
  target: any;
  originalEvent: any;
  point: { x: number; y: number };
  lngLat: { lng: number; lat: number };
}

export interface MapboxMarkerOptions {
  color?: string;
  scale?: number;
  draggable?: boolean;
  pitchAlignment?: 'map' | 'viewport' | 'auto';
  rotationAlignment?: 'map' | 'viewport' | 'auto';
}

// Типы для интеграции с API
export interface GISApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface NearbySearchParams {
  latitude: number;
  longitude: number;
  radius: number;
  category?: string;
  limit?: number;
}

export interface RouteCalculationParams {
  origin: [number, number];
  destination: [number, number];
  mode?: 'driving' | 'walking' | 'cycling' | 'transit';
  alternatives?: boolean;
}

export interface RouteResponse {
  routes: Array<{
    geometry: any;
    distance: number;
    duration: number;
    steps: Array<{
      instruction: string;
      distance: number;
      duration: number;
      geometry: any;
    }>;
  }>;
  waypoints: Array<{
    location: [number, number];
    name: string;
  }>;
}

// Типы для кластеризации
export interface ClusterFeature {
  type: 'Feature';
  properties: {
    cluster: boolean;
    cluster_id?: number;
    point_count?: number;
    point_count_abbreviated?: string;
  };
  geometry: {
    type: 'Point';
    coordinates: [number, number];
  };
}

export interface MapClusterOptions {
  radius: number;
  maxZoom: number;
  minZoom: number;
  extent: number;
  nodeSize: number;
}

// Типы для API кластеров
export interface ClusterData {
  center: {
    lat: number;
    lng: number;
  };
  count: number;
  bounds: {
    north: number;
    south: number;
    east: number;
    west: number;
  };
  zoom_expand: number;
}

export interface ClusterResponse {
  clusters: ClusterData[];
  listings: MapMarkerData[];
  total_count: number;
}

export interface ClusterRequestParams {
  bounds: string; // "north,south,east,west"
  zoom_level: number;
  categories?: string;
  min_price?: number;
  max_price?: number;
  currency?: string;
  grid_size?: number;
}

// Типы для геозон
export interface Geofence {
  id: string;
  name: string;
  type: 'circle' | 'polygon';
  coordinates: number[][];
  radius?: number;
  active: boolean;
  events: Array<'enter' | 'exit' | 'dwell'>;
}

export interface GeofenceEvent {
  id: string;
  geofenceId: string;
  type: 'enter' | 'exit' | 'dwell';
  timestamp: Date;
  location: GeoLocation;
  userId?: string;
}

// Типы для компонентов карты
export interface MapClusterProps {
  count: number;
  onClick?: () => void;
}
