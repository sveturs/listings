// Детальная проверка карты и маркеров
const checkMap = () => {
  const result = {
    timestamp: new Date().toISOString(),
    url: window.location.href,
    checks: {}
  };

  // 1. Проверка window._map
  result.checks.windowMap = {
    exists: typeof window._map \!== 'undefined',
    type: typeof window._map,
    constructor: window._map ? window._map.constructor.name : null
  };

  // 2. Проверка Leaflet
  result.checks.leaflet = {
    exists: typeof L \!== 'undefined',
    version: typeof L \!== 'undefined' && L.version ? L.version : null
  };

  // 3. Поиск маркеров через _layers
  if (window._map && window._map._layers) {
    const layers = Object.values(window._map._layers);
    const markers = layers.filter(layer => {
      return layer instanceof L.Marker || 
             layer._icon || 
             layer.options?.icon ||
             layer.constructor.name.includes('Marker');
    });
    
    result.checks.markersFromLayers = {
      totalLayers: layers.length,
      markerCount: markers.length,
      markerDetails: markers.map((m, i) => ({
        index: i,
        hasIcon: \!\!m._icon,
        hasLatLng: \!\!m._latlng,
        latLng: m._latlng ? [m._latlng.lat, m._latlng.lng] : null,
        className: m._icon ? m._icon.className : null
      }))
    };
  }

  // 4. Поиск маркеров через DOM
  const domMarkers = document.querySelectorAll('.leaflet-marker-icon');
  result.checks.domMarkers = {
    count: domMarkers.length,
    details: Array.from(domMarkers).map((el, i) => ({
      index: i,
      className: el.className,
      style: el.style.cssText,
      src: el.src || null
    }))
  };

  // 5. Поиск кластеров
  const clusters = document.querySelectorAll('.leaflet-marker-cluster');
  result.checks.clusters = {
    count: clusters.length,
    details: Array.from(clusters).map((el, i) => ({
      index: i,
      text: el.textContent,
      className: el.className
    }))
  };

  // 6. Проверка других элементов карты
  result.checks.mapElements = {
    popups: document.querySelectorAll('.leaflet-popup').length,
    tooltips: document.querySelectorAll('.leaflet-tooltip').length,
    panes: document.querySelectorAll('.leaflet-pane').length,
    tiles: document.querySelectorAll('.leaflet-tile').length
  };

  // 7. Проверка глобальных переменных
  result.checks.globals = {
    hasMarkers: typeof window.markers \!== 'undefined',
    hasMarkerCluster: typeof window.markerCluster \!== 'undefined',
    hasMapData: typeof window.mapData \!== 'undefined'
  };

  return JSON.stringify(result, null, 2);
};

checkMap();
EOF < /dev/null