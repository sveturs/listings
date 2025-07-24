'use client';

import React, { useState, useEffect, useRef } from 'react';

interface SpringTile {
  id: number;
  x: number;
  y: number;
  targetX: number;
  targetY: number;
  velocityX: number;
  velocityY: number;
  color: string;
  icon: string;
  scale: number;
  targetScale: number;
  scaleVelocity: number;
  rotation: number;
  targetRotation: number;
  rotationVelocity: number;
  bouncing: boolean;
}

interface SveTuLogoSpringProps {
  width?: number;
  height?: number;
}

export const SveTuLogoSpring: React.FC<SveTuLogoSpringProps> = ({
  width = 200,
  height = 200,
}) => {
  const [tiles, setTiles] = useState<SpringTile[]>([]);
  const [isAnimating, setIsAnimating] = useState(false);
  const animationRef = useRef<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Параметры физики пружины (успокоенные)
  const SPRING_STIFFNESS = 0.08; // Уменьшена жесткость пружины
  const DAMPING = 0.92; // Увеличено демпфирование
  const MIN_VELOCITY = 0.01;

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

    const initialTiles: SpringTile[] = positions.map((pos, index) => ({
      id: index,
      x: pos.x,
      y: pos.y,
      targetX: pos.x,
      targetY: pos.y,
      velocityX: 0,
      velocityY: 0,
      color: colors[index],
      icon: icons[index],
      scale: 1,
      targetScale: 1,
      scaleVelocity: 0,
      rotation: 0,
      targetRotation: 0,
      rotationVelocity: 0,
      bouncing: false,
    }));

    setTiles(initialTiles);
  }, [width]);

  // Функция обновления физики пружины
  const updateSpringPhysics = (tiles: SpringTile[]): SpringTile[] => {
    return tiles.map((tile) => {
      // Позиция
      const forceX = (tile.targetX - tile.x) * SPRING_STIFFNESS;
      const forceY = (tile.targetY - tile.y) * SPRING_STIFFNESS;

      tile.velocityX = (tile.velocityX + forceX) * DAMPING;
      tile.velocityY = (tile.velocityY + forceY) * DAMPING;

      tile.x += tile.velocityX;
      tile.y += tile.velocityY;

      // Масштаб
      const scaleForce = (tile.targetScale - tile.scale) * SPRING_STIFFNESS;
      tile.scaleVelocity = (tile.scaleVelocity + scaleForce) * DAMPING;
      tile.scale += tile.scaleVelocity;

      // Поворот
      const rotationForce =
        (tile.targetRotation - tile.rotation) * SPRING_STIFFNESS;
      tile.rotationVelocity = (tile.rotationVelocity + rotationForce) * DAMPING;
      tile.rotation += tile.rotationVelocity;

      // Проверка на завершение анимации
      const isMoving =
        Math.abs(tile.velocityX) > MIN_VELOCITY ||
        Math.abs(tile.velocityY) > MIN_VELOCITY ||
        Math.abs(tile.scaleVelocity) > MIN_VELOCITY ||
        Math.abs(tile.rotationVelocity) > MIN_VELOCITY;

      tile.bouncing = isMoving;

      return { ...tile };
    });
  };

  // Основной цикл анимации
  useEffect(() => {
    if (!isAnimating) return;

    const animate = () => {
      setTiles((prevTiles) => {
        const updatedTiles = updateSpringPhysics(prevTiles);
        const stillAnimating = updatedTiles.some((tile) => tile.bouncing);

        if (!stillAnimating) {
          setIsAnimating(false);
        }

        return updatedTiles;
      });

      if (isAnimating) {
        animationRef.current = requestAnimationFrame(animate);
      }
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [isAnimating]);

  // Функция запуска анимации
  const triggerSpringAnimation = () => {
    if (isAnimating) return;

    setIsAnimating(true);

    const scale = width < 200 ? width / 200 : 1;

    setTiles((prevTiles) => {
      return prevTiles.map((tile, index) => {
        // Случайное смещение (уменьшено для более спокойного эффекта)
        const randomX = (Math.random() - 0.5) * 120 * scale; // Было 200
        const randomY = (Math.random() - 0.5) * 120 * scale; // Было 200
        const randomScale = 0.5 + Math.random() * 0.8; // Было 0.3 + 1.4
        const randomRotation = (Math.random() - 0.5) * 360; // Было 720

        // Возвращаемся к исходным позициям через время
        setTimeout(
          () => {
            setTiles((currentTiles) =>
              currentTiles.map((t) => {
                if (t.id === tile.id) {
                  const baseOriginalPos = [
                    { x: -60, y: -60 },
                    { x: 0, y: -60 },
                    { x: 60, y: -60 },
                    { x: -60, y: 0 },
                    { x: 0, y: 0 },
                    { x: 60, y: 0 },
                    { x: -60, y: 60 },
                    { x: 0, y: 60 },
                    { x: 60, y: 60 },
                  ][tile.id];

                  const originalPos = {
                    x: baseOriginalPos.x * scale,
                    y: baseOriginalPos.y * scale,
                  };

                  return {
                    ...t,
                    targetX: originalPos.x,
                    targetY: originalPos.y,
                    targetScale: 1,
                    targetRotation: 0,
                  };
                }
                return t;
              })
            );
          },
          100 + index * 50
        );

        return {
          ...tile,
          targetX: tile.x + randomX,
          targetY: tile.y + randomY,
          targetScale: randomScale,
          targetRotation: randomRotation,
          velocityX: 0,
          velocityY: 0,
          scaleVelocity: 0,
          rotationVelocity: 0,
          bouncing: true,
        };
      });
    });
  };

  // Реакция на hover
  const handleMouseEnter = () => {
    setTiles((prevTiles) =>
      prevTiles.map((tile) => ({
        ...tile,
        targetScale: 1.05, // Уменьшено с 1.1
        targetRotation: tile.targetRotation + 8, // Уменьшено с 15
      }))
    );
    if (!isAnimating) {
      setIsAnimating(true);
    }
  };

  const handleMouseLeave = () => {
    setTiles((prevTiles) =>
      prevTiles.map((tile) => ({
        ...tile,
        targetScale: 1,
        targetRotation: 0,
      }))
    );
    if (!isAnimating) {
      setIsAnimating(true);
    }
  };

  return (
    <div
      ref={containerRef}
      className="relative cursor-pointer select-none"
      style={{ width, height }}
      onClick={triggerSpringAnimation}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {/* Центральная точка для визуального ориентира */}
      <div
        className="absolute w-2 h-2 bg-gray-400 rounded-full opacity-30"
        style={{
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)',
        }}
      />

      {/* Плитки */}
      {tiles.map((tile) => {
        const scale = width < 200 ? width / 200 : 1;
        const tileSize = 50 * scale;

        return (
          <div
            key={tile.id}
            className="absolute transition-shadow duration-200"
            style={{
              left: '50%',
              top: '50%',
              width: `${tileSize}px`,
              height: `${tileSize}px`,
              transform: `
                translate(-50%, -50%)
                translate(${tile.x}px, ${tile.y}px)
                scale(${tile.scale})
                rotate(${tile.rotation}deg)
              `,
              zIndex: tile.bouncing ? 10 : 1,
            }}
          >
            {/* Тень с пружинным эффектом */}
            <div
              className="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 rounded-full bg-black transition-all duration-100"
              style={{
                width: `${30 + tile.scale * 10}px`,
                height: `${6 + tile.scale * 2}px`,
                opacity: Math.max(0.1, 0.3 - Math.abs(tile.velocityY) * 0.01),
                filter: `blur(${Math.abs(tile.velocityY) * 0.1}px)`,
              }}
            />

            {/* Основная плитка */}
            <div
              className={`
              w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-lg
              transition-all duration-200 relative overflow-hidden
              ${tile.bouncing ? 'shadow-2xl' : 'shadow-lg'}
            `}
              style={{
                backgroundColor: tile.color,
                boxShadow: tile.bouncing
                  ? `0 ${5 + Math.abs(tile.velocityY) * 0.5}px ${15 + Math.abs(tile.velocityY)}px rgba(0,0,0,0.3)`
                  : '0 5px 15px rgba(0,0,0,0.2)',
              }}
            >
              {/* Световые блики при движении */}
              {tile.bouncing && (
                <div
                  className="absolute inset-0 opacity-30"
                  style={{
                    background: `linear-gradient(
                    ${tile.velocityX > 0 ? '45deg' : '135deg'}, 
                    rgba(255,255,255,0.8) 0%, 
                    transparent 50%
                  )`,
                  }}
                />
              )}

              {/* Индикатор скорости */}
              {tile.bouncing &&
                Math.abs(tile.velocityX) + Math.abs(tile.velocityY) > 2 && (
                  <div
                    className="absolute inset-0 animate-pulse"
                    style={{
                      background: `radial-gradient(circle, rgba(255,255,255,0.3) 0%, transparent 70%)`,
                      animationDuration: `${Math.max(0.2, 1 - (Math.abs(tile.velocityX) + Math.abs(tile.velocityY)) * 0.1)}s`,
                    }}
                  />
                )}

              {/* Иконка с пружинным эффектом */}
              <span
                className="relative z-10 transition-all duration-100"
                style={{
                  fontSize: `${1.4 + tile.scale * 0.28}rem`,
                  transform: `rotate(${-tile.rotation * 0.3}deg)`, // Противоположное вращение иконки
                  filter: tile.bouncing ? 'brightness(1.2)' : 'none',
                }}
              >
                {tile.icon}
              </span>
            </div>

            {/* Следы движения */}
            {tile.bouncing &&
              (Math.abs(tile.velocityX) > 1 ||
                Math.abs(tile.velocityY) > 1) && (
                <>
                  {[...Array(3)].map((_, i) => (
                    <div
                      key={i}
                      className="absolute w-full h-full rounded-lg pointer-events-none"
                      style={{
                        backgroundColor: tile.color,
                        opacity: Math.max(0, 0.3 - i * 0.1),
                        transform: `
                      translate(${-tile.velocityX * (i + 1) * 0.3}px, ${-tile.velocityY * (i + 1) * 0.3}px)
                      scale(${Math.max(0.8, tile.scale - i * 0.1)})
                    `,
                        zIndex: -i - 1,
                      }}
                    />
                  ))}
                </>
              )}
          </div>
        );
      })}

      {/* Инструкция */}
      <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full mt-4 text-xs text-gray-500 text-center whitespace-nowrap">
        {isAnimating ? 'Анимация...' : 'Кликните для запуска'}
      </div>

      {/* Волны от центра при клике */}
      {isAnimating && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className="absolute border border-purple-300 rounded-full animate-ping"
              style={{
                width: `${50 + i * 30}px`,
                height: `${50 + i * 30}px`,
                animationDelay: `${i * 100}ms`,
                animationDuration: '1s',
              }}
            />
          ))}
        </div>
      )}
    </div>
  );
};
