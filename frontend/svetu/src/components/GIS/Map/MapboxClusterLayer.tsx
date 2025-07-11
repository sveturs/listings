'use client';

import React, { useCallback, useEffect, useMemo } from 'react';
import { Source, Layer, useMap } from 'react-map-gl';
import type { LayerProps } from 'react-map-gl';
import { MapMarkerData } from '../types/gis';
import { markersToGeoJSON, getMarkerColor } from '../utils/geoJsonHelpers';

interface MapboxClusterLayerProps {
  /** Массив маркеров для кластеризации */
  markers: MapMarkerData[];
  /** Радиус кластеризации в пикселях */
  clusterRadius?: number;
  /** Максимальный уровень зума для кластеризации */
  clusterMaxZoom?: number;
  /** Минимальный размер кластера */
  clusterMinPoints?: number;
  /** Обработчик клика по кластеру */
  onClusterClick?: (clusterId: number, coordinates: [number, number]) => void;
  /** Обработчик клика по индивидуальному маркеру */
  onMarkerClick?: (marker: MapMarkerData) => void;
  /** Показывать ли цены на маркерах объявлений */
  showPrices?: boolean;
  /** Настройки стилей кластеров */
  clusterStyles?: {
    small?: {
      color?: string;
      size?: number;
      textColor?: string;
    };
    medium?: {
      color?: string;
      size?: number;
      textColor?: string;
    };
    large?: {
      color?: string;
      size?: number;
      textColor?: string;
    };
  };
}

const MapboxClusterLayer: React.FC<MapboxClusterLayerProps> = ({
  markers,
  clusterRadius = 50,
  clusterMaxZoom = 14,
  clusterMinPoints = 2,
  onClusterClick,
  onMarkerClick,
  showPrices = false,
  clusterStyles = {},
}) => {
  const mapRef = useMap();

  // Преобразуем маркеры в GeoJSON
  const geoJsonData = useMemo(() => {
    return markersToGeoJSON(markers);
  }, [markers]);

  // Дефолтные стили кластеров
  const defaultClusterStyles = useMemo(
    () => ({
      small: {
        color: '#3b82f6',
        size: 40,
        textColor: '#ffffff',
        ...clusterStyles.small,
      },
      medium: {
        color: '#059669',
        size: 50,
        textColor: '#ffffff',
        ...clusterStyles.medium,
      },
      large: {
        color: '#dc2626',
        size: 60,
        textColor: '#ffffff',
        ...clusterStyles.large,
      },
    }),
    [clusterStyles]
  );

  // Обработчик клика по кластеру
  const handleClusterClick = useCallback(
    (event: any) => {
      const features = event.features;
      if (!features || features.length === 0) return;

      const feature = features[0];
      if (!feature.properties?.cluster) return;

      const clusterId = feature.properties.cluster_id;
      const coordinates = (feature.geometry as any).coordinates as [
        number,
        number,
      ];

      if (onClusterClick) {
        onClusterClick(clusterId, coordinates);
      } else {
        // Дефолтное поведение - увеличить зум к кластеру
        if (mapRef.current) {
          const map = mapRef.current.getMap();
          const source = map.getSource('markers') as mapboxgl.GeoJSONSource;

          if (source) {
            source.getClusterExpansionZoom(clusterId, (err, zoom) => {
              if (err) return;

              mapRef.current?.flyTo({
                center: coordinates,
                zoom: zoom || 15,
                duration: 1000,
              });
            });
          }
        }
      }
    },
    [onClusterClick, mapRef]
  );

  // Обработчик клика по индивидуальному маркеру
  const handleMarkerClick = useCallback(
    (event: any) => {
      const features = event.features;
      if (!features || features.length === 0) return;

      const feature = features[0];
      if (feature.properties?.cluster) return; // Игнорируем кластеры

      const marker: MapMarkerData = {
        id: feature.properties?.id || '',
        position: [
          (feature.geometry as any).coordinates[0],
          (feature.geometry as any).coordinates[1],
        ],
        title: feature.properties?.title || '',
        description: feature.properties?.description,
        type: feature.properties?.type || 'listing',
        data: feature.properties?.data,
      };

      if (onMarkerClick) {
        onMarkerClick(marker);
      }
    },
    [onMarkerClick]
  );

  // Обработчик наведения на кластер
  const handleClusterMouseEnter = useCallback(
    (event: any) => {
      const features = event.features;
      if (!features || features.length === 0) return;

      const feature = features[0];
      if (feature.properties?.cluster) {
        if (mapRef.current) {
          mapRef.current.getCanvas().style.cursor = 'pointer';
        }
      }
    },
    [mapRef]
  );

  // Обработчик ухода курсора с кластера
  const handleClusterMouseLeave = useCallback(() => {
    if (mapRef.current) {
      mapRef.current.getCanvas().style.cursor = '';
    }
  }, [mapRef]);

  // Настройки источника данных для кластеризации
  const clusterSourceOptions = useMemo(
    () => ({
      cluster: true,
      clusterMaxZoom: clusterMaxZoom,
      clusterRadius: clusterRadius,
      clusterMinPoints: clusterMinPoints,
      clusterProperties: {
        // Подсчет количества маркеров каждого типа в кластере
        listings: ['+', ['case', ['==', ['get', 'type'], 'listing'], 1, 0]],
        users: ['+', ['case', ['==', ['get', 'type'], 'user'], 1, 0]],
        pois: ['+', ['case', ['==', ['get', 'type'], 'poi'], 1, 0]],
      },
    }),
    [clusterMaxZoom, clusterRadius, clusterMinPoints]
  );

  // Слой для кластеров
  const clusterLayer: LayerProps = useMemo(
    () => ({
      id: 'clusters',
      type: 'circle',
      source: 'markers',
      filter: ['has', 'point_count'],
      paint: {
        'circle-color': [
          'step',
          ['get', 'point_count'],
          defaultClusterStyles.small.color,
          10,
          defaultClusterStyles.medium.color,
          50,
          defaultClusterStyles.large.color,
        ],
        'circle-radius': [
          'step',
          ['get', 'point_count'],
          defaultClusterStyles.small.size / 2,
          10,
          defaultClusterStyles.medium.size / 2,
          50,
          defaultClusterStyles.large.size / 2,
        ],
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.8,
        'circle-stroke-opacity': 1,
      },
    }),
    [defaultClusterStyles]
  );

  // Слой для текста кластеров
  const clusterCountLayer: LayerProps = useMemo(
    () => ({
      id: 'cluster-count',
      type: 'symbol',
      source: 'markers',
      filter: ['has', 'point_count'],
      layout: {
        'text-field': ['get', 'point_count_abbreviated'],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 16,
        'text-allow-overlap': true,
        'text-ignore-placement': true,
      },
      paint: {
        'text-color': '#ffffff',
        'text-halo-color': 'rgba(0, 0, 0, 0.5)',
        'text-halo-width': 1,
      },
    }),
    []
  );

  // Слой для индивидуальных маркеров
  const unclusteredPointLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-point',
      type: 'circle',
      source: 'markers',
      filter: ['!', ['has', 'point_count']],
      paint: {
        'circle-color': [
          'case',
          ['==', ['get', 'type'], 'listing'],
          getMarkerColor('listing'),
          ['==', ['get', 'type'], 'user'],
          getMarkerColor('user'),
          ['==', ['get', 'type'], 'poi'],
          getMarkerColor('poi'),
          '#6b7280',
        ],
        'circle-radius': [
          'case',
          ['==', ['get', 'type'], 'listing'],
          12,
          ['==', ['get', 'type'], 'user'],
          10,
          ['==', ['get', 'type'], 'poi'],
          9,
          8,
        ],
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.9,
        'circle-stroke-opacity': 1,
      },
    }),
    []
  );

  // Слой для иконок/текста индивидуальных маркеров
  const unclusteredPointTextLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-point-text',
      type: 'symbol',
      source: 'markers',
      filter: ['!', ['has', 'point_count']],
      layout: {
        'text-field': showPrices && ['==', ['get', 'type'], 'listing']
          ? ['concat', ['get', 'data.price'], '€']
          : ['get', 'icon'],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': showPrices && ['==', ['get', 'type'], 'listing'] ? 10 : 12,
        'text-allow-overlap': true,
        'text-ignore-placement': true,
      },
      paint: {
        'text-color': '#ffffff',
        'text-halo-color': 'rgba(0, 0, 0, 0.5)',
        'text-halo-width': 1,
      },
    }),
    [showPrices]
  );

  // Эффект для обработки кликов по кластерам
  useEffect(() => {
    if (!mapRef.current) return;

    const map = mapRef.current.getMap();

    // Добавляем обработчики событий
    map.on('click', 'clusters', handleClusterClick);
    map.on('click', 'unclustered-point', handleMarkerClick);
    map.on('mouseenter', 'clusters', handleClusterMouseEnter);
    map.on('mouseleave', 'clusters', handleClusterMouseLeave);

    return () => {
      // Удаляем обработчики событий
      map.off('click', 'clusters', handleClusterClick);
      map.off('click', 'unclustered-point', handleMarkerClick);
      map.off('mouseenter', 'clusters', handleClusterMouseEnter);
      map.off('mouseleave', 'clusters', handleClusterMouseLeave);
    };
  }, [
    mapRef,
    handleClusterClick,
    handleMarkerClick,
    handleClusterMouseEnter,
    handleClusterMouseLeave,
  ]);

  return (
    <Source
      id="markers"
      type="geojson"
      data={geoJsonData}
      {...clusterSourceOptions}
    >
      <Layer {...clusterLayer} />
      <Layer {...clusterCountLayer} />
      <Layer {...unclusteredPointLayer} />
      <Layer {...unclusteredPointTextLayer} />
    </Source>
  );
};

export default MapboxClusterLayer;
