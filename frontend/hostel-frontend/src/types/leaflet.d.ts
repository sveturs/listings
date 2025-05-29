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
  
  // Экспортируем дополнительные типы
  export class LatLngBounds {
    extend(latlng: LatLngExpression | LatLngBoundsExpression): this;
    getCenter(): LatLng;
    getSouthWest(): LatLng;
    getNorthEast(): LatLng;
  }
  
  export interface Handler {
    enable(): this;
    disable(): this;
    enabled(): boolean;
  }
  
  export const control: {
    zoom(options?: Control.ZoomOptions): Control.Zoom;
  };

  export interface MapOptions {
    keepBuffer?: number;
    updateWhenZooming?: boolean;
    updateWhenIdle?: boolean;
    wheelDebounceTime?: number;
    zoomAnimation?: boolean;
    fadeAnimation?: boolean;
    markerZoomAnimation?: boolean;
    zoom?: number;
    center?: LatLngExpression;
    zoomSnap?: number;
    zoomDelta?: number;
    wheelPxPerZoomLevel?: number;
    preferCanvas?: boolean;
    renderer?: Renderer;
    scrollWheelZoom?: boolean;
    dragging?: boolean;
    touchZoom?: boolean;
    doubleClickZoom?: boolean;
    boxZoom?: boolean;
    keyboard?: boolean;
    zoomControl?: boolean;
    attributionControl?: boolean;
    tap?: boolean;
    tapTolerance?: number;
    bounceAtZoomLimits?: boolean;
    transform3DLimit?: number;
    inertia?: boolean;
    inertiaDeceleration?: number;
    inertiaMaxSpeed?: number;
    worldCopyJump?: boolean;
    maxZoom?: number;
    minZoom?: number;
    updateInterval?: number;
  }

  // Экспортируем классы для instanceof проверок
  export class Map {
    addLayer(layer: Layer): this;
    removeLayer(layer: Layer): this;
    hasLayer(layer: Layer): boolean;
    setView(center: LatLngExpression, zoom?: number, options?: ZoomPanOptions): this;
    fitBounds(bounds: LatLngBoundsExpression, options?: FitBoundsOptions): this;
    getBounds(): LatLngBounds;
    getCenter(): LatLng;
    getZoom(): number;
    on(type: string, fn: Function, context?: any): this;
    off(type: string, fn?: Function, context?: any): this;
    invalidateSize(options?: boolean | InvalidateSizeOptions): this;
    eachLayer(fn: (layer: Layer) => void, context?: any): this;
    remove(): this;
    whenReady(fn: () => void, context?: any): this;
    dragging: Handler;
  }
  
  export class Marker {
    addTo(map: Map | LayerGroup): this;
    bindPopup(content: string | HTMLElement | Function | Popup, options?: PopupOptions): this;
    bindTooltip(content: string | HTMLElement | Function | Tooltip, options?: TooltipOptions): this;
    openPopup(): this;
    on(type: string, fn: Function, context?: any): this;
    getLatLng(): LatLng;
    setLatLng(latlng: LatLngExpression): this;
  }
  
  export class LayerGroup {
    addLayer(layer: Layer): this;
    removeLayer(layer: Layer): this;
    clearLayers(): this;
    addTo(map: Map | LayerGroup): this;
  }
  
  export class FeatureGroup extends LayerGroup {
    getBounds(): LatLngBounds;
    addLayer(layer: Layer): this;
    getLayers(): Layer[];
  }
  
  export class TileLayer {
    addTo(map: Map | LayerGroup): this;
  }
  
  export class DivIcon {}
  export class Icon {
    static Default: typeof Icon & {
      prototype: Icon;
      mergeOptions(options: IconOptions): void;
    };
  }
  export class Circle {
    addTo(map: Map | LayerGroup): this;
  }
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