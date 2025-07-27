'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';

interface MorphTile {
  id: number;
  x: number;
  y: number;
  originalX: number;
  originalY: number;
  color: string;
  icon: string;
  morphProgress: number;
  waveOffset: number;
  scale: number;
  opacity: number;
  blur: number;
  hue: number;
}

interface SveTuLogoMorphingProps {
  width?: number;
  height?: number;
}

export const SveTuLogoMorphing: React.FC<SveTuLogoMorphingProps> = ({
  width = 200,
  height = 200,
}) => {
  const [tiles, setTiles] = useState<MorphTile[]>([]);
  const [time, setTime] = useState(0);
  const [morphMode, setMorphMode] = useState<
    'wave' | 'spiral' | 'explosion' | 'liquid'
  >('wave');
  const animationRef = useRef<number | null>(null);

  // Инициализация плиток
  useEffect(() => {
    const scale = width < 200 ? width / 200 : 1;

    const basePositions = [
      { x: -60, y: -60 },
      { x: 0, y: -60 },
      { x: 60, y: -60 },
      { x: -60, y: 0 },
      { x: 0, y: 0 },
      { x: 60, y: 0 },
      { x: -60, y: 60 },
      { x: 0, y: 60 },
      { x: 60, y: 60 },
    ];

    const positions = basePositions.map((pos) => ({
      x: pos.x * scale,
      y: pos.y * scale,
    }));

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

    const initialTiles: MorphTile[] = positions.map((pos, index) => ({
      id: index,
      x: pos.x,
      y: pos.y,
      originalX: pos.x,
      originalY: pos.y,
      color: colors[index],
      icon: icons[index],
      morphProgress: 0,
      waveOffset: index * 0.3,
      scale: 1,
      opacity: 1,
      blur: 0,
      hue: 0,
    }));

    setTiles(initialTiles);
  }, [width]);

  // Функции морфинга
  const calculateWavePosition = useCallback(
    (tile: MorphTile, time: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const waveX = Math.sin(time * 0.8 + tile.waveOffset) * 30 * scale;
      const waveY = Math.cos(time * 1.2 + tile.waveOffset * 1.5) * 20 * scale;
      return {
        x: tile.originalX + waveX,
        y: tile.originalY + waveY,
        scale: 1 + Math.sin(time * 2 + tile.waveOffset) * 0.3,
        hue: Math.sin(time * 0.5 + tile.waveOffset) * 60,
      };
    },
    [width]
  );

  const calculateSpiralPosition = useCallback(
    (tile: MorphTile, time: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const angle = time * 0.5 + tile.id * ((Math.PI * 2) / 9);
      const radius = (40 + Math.sin(time * 0.3) * 20) * scale;
      return {
        x: Math.cos(angle) * radius,
        y: Math.sin(angle) * radius,
        scale: 1 + Math.sin(time * 3 + tile.id) * 0.2,
        hue: (time * 30 + tile.id * 40) % 360,
      };
    },
    [width]
  );

  const calculateExplosionPosition = useCallback(
    (tile: MorphTile, time: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const explosionProgress = Math.sin(time * 0.4) * 0.5 + 0.5;
      const distance = explosionProgress * 80 * scale;
      const angle = tile.id * ((Math.PI * 2) / 9);
      return {
        x: tile.originalX + Math.cos(angle) * distance,
        y: tile.originalY + Math.sin(angle) * distance,
        scale: 1 + explosionProgress * 0.5,
        hue: explosionProgress * 120,
      };
    },
    [width]
  );

  const calculateLiquidPosition = useCallback(
    (tile: MorphTile, time: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const liquidX = Math.sin(time * 0.6 + tile.id * 0.5) * 25 * scale;
      const liquidY = Math.cos(time * 0.4 + tile.id * 0.7) * 15 * scale;
      const distortion = Math.sin(time * 2 + tile.id) * 0.2;
      return {
        x: tile.originalX + liquidX,
        y: tile.originalY + liquidY,
        scale: 1 + distortion,
        hue: Math.sin(time * 0.7 + tile.id * 0.3) * 90,
      };
    },
    [width]
  );

  // Главный цикл анимации
  useEffect(() => {
    const animate = (timestamp: number) => {
      const newTime = timestamp * 0.001;
      setTime(newTime);

      setTiles((prevTiles) =>
        prevTiles.map((tile) => {
          let newPos;

          switch (morphMode) {
            case 'wave':
              newPos = calculateWavePosition(tile, newTime);
              break;
            case 'spiral':
              newPos = calculateSpiralPosition(tile, newTime);
              break;
            case 'explosion':
              newPos = calculateExplosionPosition(tile, newTime);
              break;
            case 'liquid':
              newPos = calculateLiquidPosition(tile, newTime);
              break;
            default:
              newPos = {
                x: tile.originalX,
                y: tile.originalY,
                scale: 1,
                hue: 0,
              };
          }

          // Дополнительные эффекты
          const opacity = 0.7 + Math.sin(newTime * 1.5 + tile.waveOffset) * 0.3;
          const blur = Math.abs(Math.sin(newTime * 0.8 + tile.id)) * 2;

          return {
            ...tile,
            x: newPos.x,
            y: newPos.y,
            scale: newPos.scale,
            opacity: Math.max(0.3, opacity),
            blur: blur,
            hue: newPos.hue,
          };
        })
      );

      animationRef.current = requestAnimationFrame(animate);
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [
    morphMode,
    width,
    calculateWavePosition,
    calculateSpiralPosition,
    calculateExplosionPosition,
    calculateLiquidPosition,
  ]);

  // Автоматическая смена режимов
  useEffect(() => {
    const modes: (typeof morphMode)[] = [
      'wave',
      'spiral',
      'explosion',
      'liquid',
    ];
    let currentIndex = 0;

    const switchMode = () => {
      currentIndex = (currentIndex + 1) % modes.length;
      setMorphMode(modes[currentIndex]);
    };

    const interval = setInterval(switchMode, 5000);
    return () => clearInterval(interval);
  }, []);

  // Функция для ручной смены режима
  const handleClick = () => {
    const modes: (typeof morphMode)[] = [
      'wave',
      'spiral',
      'explosion',
      'liquid',
    ];
    const currentIndex = modes.indexOf(morphMode);
    const nextIndex = (currentIndex + 1) % modes.length;
    setMorphMode(modes[nextIndex]);
  };

  // Утилита для преобразования цвета
  const adjustColor = (color: string, hue: number) => {
    // Простое изменение оттенка для демонстрации
    const colors = [
      `hsl(${210 + hue}, 70%, 55%)`, // синий
      `hsl(${120 + hue}, 70%, 55%)`, // зеленый
      `hsl(${0 + hue}, 70%, 55%)`, // красный
      `hsl(${30 + hue}, 70%, 55%)`, // оранжевый
      `hsl(${270 + hue}, 70%, 55%)`, // фиолетовый
      `hsl(${180 + hue}, 70%, 55%)`, // голубой
      `hsl(${60 + hue}, 70%, 55%)`, // желтый
      `hsl(${200 + hue}, 30%, 45%)`, // серо-синий
      `hsl(${300 + hue}, 70%, 55%)`, // пурпурный
    ];
    return colors[Math.abs(Math.round(hue / 40)) % colors.length] || color;
  };

  return (
    <div
      className="relative cursor-pointer select-none"
      style={{ width, height }}
      onClick={handleClick}
    >
      {/* Фоновые волны */}
      <div className="absolute inset-0 overflow-hidden opacity-20">
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            className="absolute inset-0 rounded-full border border-purple-300"
            style={{
              transform: `scale(${1 + Math.sin(time * 0.5 + i) * 0.3}) rotate(${time * 10 + i * 120}deg)`,
              opacity: 0.3 - i * 0.1,
            }}
          />
        ))}
      </div>

      {/* Плитки */}
      {tiles.map((tile) => {
        const scale = width < 200 ? width / 200 : 1;
        const tileSize = 50 * scale;

        return (
          <div
            key={tile.id}
            className="absolute transition-none"
            style={{
              left: '50%',
              top: '50%',
              width: `${tileSize}px`,
              height: `${tileSize}px`,
              transform: `
                translate(-50%, -50%)
                translate(${tile.x}px, ${tile.y}px)
                scale(${tile.scale})
              `,
              opacity: tile.opacity,
              filter: `blur(${tile.blur}px) hue-rotate(${tile.hue}deg)`,
            }}
          >
            {/* Энергетическое свечение */}
            <div
              className="absolute inset-0 rounded-lg"
              style={{
                background: `radial-gradient(circle, ${adjustColor(tile.color, tile.hue)}40 0%, transparent 70%)`,
                transform: `scale(${1.5 + Math.sin(time * 2 + tile.id) * 0.5})`,
              }}
            />

            {/* Основная плитка */}
            <div
              className="w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-lg relative overflow-hidden shadow-lg"
              style={{
                backgroundColor: adjustColor(tile.color, tile.hue),
                boxShadow: `0 5px 20px ${adjustColor(tile.color, tile.hue)}50`,
              }}
            >
              {/* Liquid эффект */}
              {morphMode === 'liquid' && (
                <div
                  className="absolute inset-0 opacity-30"
                  style={{
                    background: `radial-gradient(
                    ellipse ${50 + Math.sin(time * 3 + tile.id) * 20}% ${50 + Math.cos(time * 4 + tile.id) * 20}% at 
                    ${50 + Math.sin(time * 2 + tile.id) * 30}% ${50 + Math.cos(time * 2.5 + tile.id) * 30}%, 
                    rgba(255,255,255,0.8) 0%, 
                    transparent 70%
                  )`,
                  }}
                />
              )}

              {/* Взрывной эффект */}
              {morphMode === 'explosion' && (
                <div
                  className="absolute inset-0"
                  style={{
                    background: `conic-gradient(
                    from ${time * 100 + tile.id * 40}deg,
                    transparent,
                    rgba(255,255,255,0.3),
                    transparent
                  )`,
                  }}
                />
              )}

              {/* Спиральный эффект */}
              {morphMode === 'spiral' && (
                <div
                  className="absolute inset-0 opacity-20"
                  style={{
                    background: `linear-gradient(
                    ${time * 50 + tile.id * 60}deg,
                    rgba(255,255,255,0.8) 0%,
                    transparent 30%,
                    rgba(255,255,255,0.4) 70%,
                    transparent 100%
                  )`,
                  }}
                />
              )}

              {/* Иконка с эффектами */}
              <span
                className="relative z-10 transition-all duration-200"
                style={{
                  fontSize: `${1.4 + Math.sin(time * 2 + tile.id) * 0.42}rem`,
                  transform: `rotate(${Math.sin(time + tile.id) * 15}deg)`,
                  textShadow: `0 0 10px ${adjustColor(tile.color, tile.hue)}80`,
                }}
              >
                {tile.icon}
              </span>
            </div>

            {/* Следы и частицы */}
            {[...Array(2)].map((_, i) => (
              <div
                key={i}
                className="absolute w-2 h-2 rounded-full pointer-events-none"
                style={{
                  backgroundColor: adjustColor(tile.color, tile.hue),
                  top: `${20 + i * 20}%`,
                  left: `${30 + Math.sin(time * 3 + tile.id + i) * 40}%`,
                  opacity: Math.max(0, 0.6 - i * 0.3),
                  transform: `scale(${Math.sin(time * 4 + tile.id + i) * 0.5 + 0.5})`,
                  boxShadow: `0 0 10px ${adjustColor(tile.color, tile.hue)}`,
                }}
              />
            ))}
          </div>
        );
      })}

      {/* Индикатор режима */}
      <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full mt-4">
        <div className="flex items-center space-x-2">
          <div className="text-xs text-gray-500 text-center">
            Режим: <span className="font-semibold capitalize">{morphMode}</span>
          </div>
          <div className="flex space-x-1">
            {['wave', 'spiral', 'explosion', 'liquid'].map((mode) => (
              <div
                key={mode}
                className={`w-2 h-2 rounded-full transition-all duration-300 ${
                  morphMode === mode ? 'bg-purple-500' : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Интерактивные волны */}
      <div className="absolute inset-0 pointer-events-none">
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="absolute border border-purple-200 rounded-full opacity-10"
            style={{
              left: '50%',
              top: '50%',
              width: `${100 + i * 40}px`,
              height: `${100 + i * 40}px`,
              transform: `
                translate(-50%, -50%) 
                scale(${1 + Math.sin(time * 0.3 + i * 0.5) * 0.2})
                rotate(${time * 5 + i * 30}deg)
              `,
            }}
          />
        ))}
      </div>
    </div>
  );
};
