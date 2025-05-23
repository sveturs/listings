import 'leaflet';
import 'leaflet.markercluster';

declare module 'leaflet' {
  // Расширяем типы Leaflet для поддержки всех используемых функций
  export function map(element: HTMLElement, options?: MapOptions): Map;
  export function tileLayer(urlTemplate: string, options?: TileLayerOptions): TileLayer;
  export function marker(latlng: LatLngExpression, options?: MarkerOptions): Marker;
  export function layerGroup(layers?: Layer[]): LayerGroup;
  export function featureGroup(layers?: Layer[]): FeatureGroup;
  export function divIcon(options: DivIconOptions): DivIcon;
  export function icon(options: IconOptions): Icon;
  export function circle(latlng: LatLngExpression, options?: CircleOptions): Circle;
  export function canvas(options?: RendererOptions): Canvas;
  export function point(x: number, y: number, round?: boolean): Point;
  
  export const control: {
    zoom(options?: Control.ZoomOptions): Control.Zoom;
  };

  export interface MapOptions {
    keepBuffer?: number;
    updateWhenZooming?: boolean;
    updateWhenIdle?: boolean;
    wheelDebounceTime?: number;
  }

  // Экспортируем классы для instanceof проверок
  export class Map {}
  export class Marker {}
  export class LayerGroup {}
  export class FeatureGroup {}
  export class TileLayer {}
  export class DivIcon {}
  export class Icon {
    static Default: typeof Icon & {
      prototype: Icon;
      mergeOptions(options: IconOptions): void;
    };
  }
  export class Circle {}
  export class Canvas {}
  export class Point {}
  
  // Типы для MarkerClusterGroup
  export interface MarkerClusterGroupOptions {
    maxClusterRadius?: number;
    disableClusteringAtZoom?: number;
    spiderfyOnMaxZoom?: boolean;
    showCoverageOnHover?: boolean;
    zoomToBoundsOnClick?: boolean;
    removeOutsideVisibleBounds?: boolean;
    animate?: boolean;
    animateAddingMarkers?: boolean;
    chunkedLoading?: boolean;
    iconCreateFunction?: (cluster: MarkerCluster) => Icon | DivIcon;
  }

  export class MarkerClusterGroup extends FeatureGroup {
    constructor(options?: MarkerClusterGroupOptions);
    addLayer(layer: Layer): this;
    addLayers(layers: Layer[]): this;
    clearLayers(): this;
    addTo(map: Map): this;
  }

  export function markerClusterGroup(options?: MarkerClusterGroupOptions): MarkerClusterGroup;

  export interface MarkerCluster {
    getChildCount(): number;
  }
}

// Декларация для window объекта
declare global {
  interface Window {
    mapUpdateTimeout?: NodeJS.Timeout;
    fetchListingsTimeout?: NodeJS.Timeout;
    zoomUpdateTimeout?: NodeJS.Timeout;
    dragUpdateTimeout?: NodeJS.Timeout;
  }
}