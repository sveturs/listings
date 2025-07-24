'use client';

import React from 'react';

interface SveTuLogoStaticProps {
  variant?: 'gradient' | 'minimal' | 'retro' | 'neon' | 'glassmorphic';
  width?: number;
  height?: number;
}

export const SveTuLogoStatic: React.FC<SveTuLogoStaticProps> = ({
  variant = 'gradient',
  width = 200,
  height = 200,
}) => {
  const scale = Math.min(width, height) / 200;
  const tileSize = 50 * scale;
  const gap = 5 * scale;
  const totalSize = tileSize * 3 + gap * 2;
  const offset = (Math.min(width, height) - totalSize) / 2;

  const positions = [
    { x: 0, y: 0 },
    { x: 1, y: 0 },
    { x: 2, y: 0 },
    { x: 0, y: 1 },
    { x: 1, y: 1 },
    { x: 2, y: 1 },
    { x: 0, y: 2 },
    { x: 1, y: 2 },
    { x: 2, y: 2 },
  ];

  const colors = [
    '#2196F3',
    '#4CAF50',
    '#F44336',
    '#FF9800',
    '#673AB7',
    '#00BCD4',
    '#FFEB3B',
    '#607D8B',
    '#9C27B0',
  ];
  const icons = ['🛒', '🏪', '🛍️', '📦', '🏠', '🤝', '📱', '💳', '💰'];

  const renderTile = (index: number, pos: { x: number; y: number }) => {
    const x = offset + pos.x * (tileSize + gap);
    const y = offset + pos.y * (tileSize + gap);

    switch (variant) {
      case 'gradient':
        return (
          <g key={index}>
            <defs>
              <linearGradient
                id={`gradient-${index}`}
                x1="0%"
                y1="0%"
                x2="100%"
                y2="100%"
              >
                <stop offset="0%" stopColor={colors[index]} stopOpacity="1" />
                <stop
                  offset="100%"
                  stopColor={colors[index]}
                  stopOpacity="0.7"
                />
              </linearGradient>
              <filter id={`shadow-${index}`}>
                <feDropShadow
                  dx="2"
                  dy="2"
                  stdDeviation="3"
                  floodOpacity="0.2"
                />
              </filter>
            </defs>
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.15}
              fill={`url(#gradient-${index})`}
              filter={`url(#shadow-${index})`}
            />
            <text
              x={x + tileSize / 2}
              y={y + tileSize / 2}
              textAnchor="middle"
              dominantBaseline="central"
              fontSize={tileSize * 0.7}
              fill="white"
            >
              {icons[index]}
            </text>
          </g>
        );

      case 'minimal':
        return (
          <g key={index}>
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.1}
              fill={index === 4 ? colors[index] : '#E5E7EB'}
              opacity={index === 4 ? 1 : 0.5}
            />
            <text
              x={x + tileSize / 2}
              y={y + tileSize / 2}
              textAnchor="middle"
              dominantBaseline="central"
              fontSize={tileSize * 0.7}
              fill={index === 4 ? 'white' : '#6B7280'}
            >
              {icons[index]}
            </text>
          </g>
        );

      case 'retro':
        return (
          <g key={index}>
            <rect
              x={x + 2}
              y={y + 2}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.05}
              fill="black"
              opacity="0.3"
            />
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.05}
              fill={colors[index]}
              stroke="black"
              strokeWidth="2"
            />
            <rect
              x={x + tileSize * 0.1}
              y={y + tileSize * 0.1}
              width={tileSize * 0.3}
              height={tileSize * 0.3}
              fill="rgba(255,255,255,0.3)"
              rx={tileSize * 0.02}
            />
            <text
              x={x + tileSize / 2}
              y={y + tileSize / 2}
              textAnchor="middle"
              dominantBaseline="central"
              fontSize={tileSize * 0.7}
              fill="white"
              style={{ fontFamily: 'monospace', fontWeight: 'bold' }}
            >
              {icons[index]}
            </text>
          </g>
        );

      case 'neon':
        return (
          <g key={index}>
            <defs>
              <filter id={`neon-${index}`}>
                <feGaussianBlur stdDeviation="3" result="coloredBlur" />
                <feMerge>
                  <feMergeNode in="coloredBlur" />
                  <feMergeNode in="SourceGraphic" />
                </feMerge>
              </filter>
            </defs>
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.2}
              fill="none"
              stroke={colors[index]}
              strokeWidth="2"
              filter={`url(#neon-${index})`}
            />
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.2}
              fill={colors[index]}
              opacity="0.1"
            />
            <text
              x={x + tileSize / 2}
              y={y + tileSize / 2}
              textAnchor="middle"
              dominantBaseline="central"
              fontSize={tileSize * 0.7}
              fill={colors[index]}
              filter={`url(#neon-${index})`}
            >
              {icons[index]}
            </text>
          </g>
        );

      case 'glassmorphic':
        return (
          <g key={index}>
            <defs>
              <filter id={`blur-${index}`}>
                <feGaussianBlur in="SourceGraphic" stdDeviation="10" />
              </filter>
              <linearGradient
                id={`glass-${index}`}
                x1="0%"
                y1="0%"
                x2="100%"
                y2="100%"
              >
                <stop offset="0%" stopColor="rgba(255,255,255,0.4)" />
                <stop offset="100%" stopColor="rgba(255,255,255,0.1)" />
              </linearGradient>
            </defs>
            <rect
              x={x - 5}
              y={y - 5}
              width={tileSize + 10}
              height={tileSize + 10}
              rx={tileSize * 0.25}
              fill={colors[index]}
              opacity="0.2"
              filter={`url(#blur-${index})`}
            />
            <rect
              x={x}
              y={y}
              width={tileSize}
              height={tileSize}
              rx={tileSize * 0.2}
              fill={`url(#glass-${index})`}
              stroke="rgba(255,255,255,0.5)"
              strokeWidth="1"
            />
            <text
              x={x + tileSize / 2}
              y={y + tileSize / 2}
              textAnchor="middle"
              dominantBaseline="central"
              fontSize={tileSize * 0.7}
              fill={colors[index]}
            >
              {icons[index]}
            </text>
          </g>
        );

      default:
        return null;
    }
  };

  return (
    <svg
      width={width}
      height={height}
      viewBox={`0 0 ${width} ${height}`}
      style={{
        backgroundColor: variant === 'neon' ? '#0a0a0a' : 'transparent',
      }}
    >
      {positions.map((pos, index) => renderTile(index, pos))}
    </svg>
  );
};
