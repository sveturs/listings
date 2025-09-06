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
  /** Обработчик наведения на маркер */
  onMarkerHover?: (
    marker: MapMarkerData,
    event: { x: number; y: number }
  ) => void;
  /** Обработчик ухода курсора с маркера */
  onMarkerLeave?: () => void;
  /** Обработчик наведения на кластер */
  onClusterHover?: (
    clusterId: number,
    coordinates: [number, number],
    event: { x: number; y: number }
  ) => void;
  /** Обработчик ухода курсора с кластера */
  onClusterLeave?: () => void;
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
  onMarkerHover,
  onMarkerLeave,
  onClusterHover,
  onClusterLeave,
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
            source.getClusterExpansionZoom(clusterId, (_err, zoom) => {
              if (_err) return;

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

      const coordinates = (feature.geometry as any).coordinates;
      const marker: MapMarkerData = {
        id: feature.properties?.id || '',
        position: [coordinates[0], coordinates[1]],
        longitude: coordinates[0],
        latitude: coordinates[1],
        title: feature.properties?.title || '',
        description: feature.properties?.description,
        type: feature.properties?.type || 'listing',
        data: feature.properties?.data,
        imageUrl: feature.properties?.imageUrl,
        metadata: feature.properties?.metadata,
      };

      if (onMarkerClick) {
        onMarkerClick(marker);
      }
    },
    [onMarkerClick]
  );

  // Обработчик наведения на кластер
  const handleClusterMouseMove = useCallback(
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

      if (mapRef.current) {
        mapRef.current.getCanvas().style.cursor = 'pointer';
      }

      if (onClusterHover) {
        onClusterHover(clusterId, coordinates, {
          x: event.point.x,
          y: event.point.y,
        });
      }
    },
    [mapRef, onClusterHover]
  );

  // Обработчик ухода курсора с кластера
  const handleClusterMouseLeave = useCallback(() => {
    if (mapRef.current) {
      mapRef.current.getCanvas().style.cursor = '';
    }
    if (onClusterLeave) {
      onClusterLeave();
    }
  }, [mapRef, onClusterLeave]);

  // Обработчик наведения на маркер
  const handleMarkerMouseMove = useCallback(
    (event: any) => {
      const features = event.features;
      if (!features || features.length === 0) return;

      const feature = features[0];
      if (feature.properties?.cluster) return; // Игнорируем кластеры

      const coordinates = (feature.geometry as any).coordinates;
      const marker: MapMarkerData = {
        id: feature.properties?.id || '',
        position: [coordinates[0], coordinates[1]],
        longitude: coordinates[0],
        latitude: coordinates[1],
        title: feature.properties?.title || '',
        description: feature.properties?.description,
        type: feature.properties?.type || 'listing',
        data: feature.properties?.data,
        imageUrl: feature.properties?.imageUrl,
        metadata: feature.properties?.metadata,
      };

      if (mapRef.current) {
        mapRef.current.getCanvas().style.cursor = 'pointer';
      }

      if (onMarkerHover) {
        onMarkerHover(marker, { x: event.point.x, y: event.point.y });
      }
    },
    [onMarkerHover, mapRef]
  );

  // Обработчик ухода курсора с маркера
  const handleMarkerMouseLeave = useCallback(() => {
    if (mapRef.current) {
      mapRef.current.getCanvas().style.cursor = '';
    }
    if (onMarkerLeave) {
      onMarkerLeave();
    }
  }, [onMarkerLeave, mapRef]);

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
        // Агрегация цен в кластере с корректной обработкой типов
        minPrice: [
          'min',
          [
            'case',
            [
              'all',
              ['==', ['get', 'type'], 'listing'],
              ['has', 'metadata'],
              ['has', 'price', ['get', 'metadata']],
              ['!=', ['get', 'price', ['get', 'metadata']], null],
            ],
            ['number', ['get', 'price', ['get', 'metadata']]],
            999999,
          ],
        ],
        maxPrice: [
          'max',
          [
            'case',
            [
              'all',
              ['==', ['get', 'type'], 'listing'],
              ['has', 'metadata'],
              ['has', 'price', ['get', 'metadata']],
              ['!=', ['get', 'price', ['get', 'metadata']], null],
            ],
            ['number', ['get', 'price', ['get', 'metadata']]],
            0,
          ],
        ],
        totalListings: [
          '+',
          ['case', ['==', ['get', 'type'], 'listing'], 1, 0],
        ],
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

  // Слой для отображения диапазона цен под кластерами
  const clusterPriceLayer: LayerProps = useMemo(
    () => ({
      id: 'cluster-price',
      type: 'symbol',
      source: 'markers',
      filter: [
        'all',
        ['has', 'point_count'],
        ['>', ['get', 'totalListings'], 0],
        ['<', ['get', 'minPrice'], 999999], // Показываем только если есть реальные цены
        ['>', ['get', 'maxPrice'], 0], // Показываем только если есть реальные цены
      ],
      layout: {
        'text-field': [
          'case',
          ['==', ['get', 'minPrice'], ['get', 'maxPrice']],
          [
            'concat',
            [
              'number-format',
              ['get', 'minPrice'],
              { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
            ],
            ' RSD',
          ],
          [
            'concat',
            [
              'number-format',
              ['get', 'minPrice'],
              { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
            ],
            '-',
            [
              'number-format',
              ['get', 'maxPrice'],
              { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
            ],
            ' RSD',
          ],
        ],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 11,
        'text-allow-overlap': false,
        'text-ignore-placement': false,
        'text-anchor': 'top',
        'text-offset': [0, 2.5],
      },
      paint: {
        'text-color': '#1f2937',
        'text-halo-color': '#ffffff',
        'text-halo-width': 3,
        'text-halo-blur': 1,
      },
    }),
    []
  );

  // Слой для отображения цен под индивидуальными маркерами объявлений
  const unclusteredPriceLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-price',
      type: 'symbol',
      source: 'markers',
      filter: [
        'all',
        ['!', ['has', 'point_count']],
        ['==', ['get', 'type'], 'listing'],
        ['has', 'metadata'],
        ['has', 'price', ['get', 'metadata']],
        ['!=', ['get', 'price', ['get', 'metadata']], null],
      ],
      layout: {
        'text-field': [
          'concat',
          [
            'number-format',
            ['number', ['get', 'price', ['get', 'metadata']]],
            { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
          ],
          ' RSD',
        ],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 10,
        'text-allow-overlap': false,
        'text-ignore-placement': false,
        'text-anchor': 'top',
        'text-offset': [0, 1.8],
      },
      paint: {
        'text-color': '#1f2937',
        'text-halo-color': '#ffffff',
        'text-halo-width': 3,
        'text-halo-blur': 1,
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

  // Слой для иконок индивидуальных маркеров
  const unclusteredPointTextLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-point-text',
      type: 'symbol',
      source: 'markers',
      filter: ['!', ['has', 'point_count']],
      layout: {
        'text-field': ['get', 'icon'],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 12,
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

  // Эффект для обработки кликов по кластерам
  useEffect(() => {
    if (!mapRef.current) return;

    const map = mapRef.current.getMap();

    // Добавляем обработчики событий
    map.on('click', 'clusters', handleClusterClick);
    map.on('click', 'unclustered-point', handleMarkerClick);
    map.on('mousemove', 'clusters', handleClusterMouseMove);
    map.on('mouseleave', 'clusters', handleClusterMouseLeave);
    map.on('mousemove', 'unclustered-point', handleMarkerMouseMove);
    map.on('mouseleave', 'unclustered-point', handleMarkerMouseLeave);

    return () => {
      // Удаляем обработчики событий
      map.off('click', 'clusters', handleClusterClick);
      map.off('click', 'unclustered-point', handleMarkerClick);
      map.off('mousemove', 'clusters', handleClusterMouseMove);
      map.off('mouseleave', 'clusters', handleClusterMouseLeave);
      map.off('mousemove', 'unclustered-point', handleMarkerMouseMove);
      map.off('mouseleave', 'unclustered-point', handleMarkerMouseLeave);
    };
  }, [
    mapRef,
    handleClusterClick,
    handleMarkerClick,
    handleClusterMouseMove,
    handleClusterMouseLeave,
    handleMarkerMouseMove,
    handleMarkerMouseLeave,
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
      <Layer {...clusterPriceLayer} />
      <Layer {...unclusteredPointLayer} />
      <Layer {...unclusteredPointTextLayer} />
      <Layer {...unclusteredPriceLayer} />
    </Source>
  );
};

export default MapboxClusterLayer;
