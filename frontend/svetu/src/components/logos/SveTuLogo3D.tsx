'use client';

import React, { useState, useEffect, useRef } from 'react';

interface Tile3D {
  id: number;
  x: number;
  y: number;
  z: number;
  baseX: number; // базовая позиция для расчета влияния мыши
  baseY: number;
  baseZ: number;
  rotateX: number;
  rotateY: number;
  rotateZ: number;
  color: string;
  icon: string;
  scale: number;
  opacity: number;
  glowing: boolean;
}

interface SveTuLogo3DProps {
  width?: number;
  height?: number;
}

export const SveTuLogo3D: React.FC<SveTuLogo3DProps> = ({
  width = 200,
  height = 200,
}) => {
  const [tiles, setTiles] = useState<Tile3D[]>([]);
  const [isHovering, setIsHovering] = useState(false);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
  const [isExploding, setIsExploding] = useState(false);
  const animationRef = useRef<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Инициализация плиток
  useEffect(() => {
    // Масштабируем позиции для всех логотипов меньше стандартного размера
    const scale = width < 200 ? width / 200 : 1;
    const basePositions = [
      { x: -60, y: -60, z: 0 },
      { x: 0, y: -60, z: 3 },
      { x: 60, y: -60, z: 0 },
      { x: -60, y: 0, z: 5 },
      { x: 0, y: 0, z: 8 },
      { x: 60, y: 0, z: 2 },
      { x: -60, y: 60, z: 0 },
      { x: 0, y: 60, z: 4 },
      { x: 60, y: 60, z: 6 },
    ];

    const positions = basePositions.map((pos) => ({
      x: pos.x * scale,
      y: pos.y * scale,
      z: pos.z * scale,
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

    const initialTiles: Tile3D[] = positions.map((pos, index) => ({
      id: index,
      x: pos.x,
      y: pos.y,
      z: pos.z,
      baseX: pos.x, // сохраняем базовые позиции
      baseY: pos.y,
      baseZ: pos.z,
      rotateX: 0,
      rotateY: 0,
      rotateZ: 0,
      color: colors[index],
      icon: icons[index],
      scale: 1,
      opacity: 1,
      glowing: false,
    }));

    setTiles(initialTiles);
  }, [width, height]);

  // Анимация плавания
  useEffect(() => {
    const animate = (timestamp: number) => {
      setTiles((prevTiles) =>
        prevTiles.map((tile, index) => {
          const time = timestamp * 0.001;
          const baseFreq = 0.4 + index * 0.08; // Восстановлена нормальная частота

          // Уменьшаем амплитуду движения для маленьких логотипов
          const ampScale = width <= 32 ? 0.25 : 1; // Уменьшаем в 4 раза для маленьких

          // Живое, но контролируемое плавание
          const baseFloatY = Math.sin(time * baseFreq) * 6 * ampScale;
          const baseFloatZ = Math.cos(time * baseFreq * 0.7) * 10 * ampScale;
          const baseRotateY = Math.sin(time * baseFreq * 0.8) * 6 * ampScale;
          const baseRotateX = Math.cos(time * baseFreq * 0.6) * 4 * ampScale;

          // Реакция на мышь (только при hover)
          let mouseInfluenceX = 0;
          let mouseInfluenceY = 0;
          let mouseInfluenceZ = 0;

          if (isHovering) {
            // Расстояние от базовой позиции плитки до мыши (в пикселях контейнера)
            const tileCenterX = tile.baseX + width / 2;
            const tileCenterY = tile.baseY + height / 2;
            const mouseDistance = Math.sqrt(
              Math.pow(mousePos.x - tileCenterX, 2) +
                Math.pow(mousePos.y - tileCenterY, 2)
            );

            // Влияние мыши уменьшается с расстоянием
            const maxInfluenceDistance = 60; // уменьшен радиус влияния
            const influence = Math.max(
              0,
              (maxInfluenceDistance - mouseDistance) / maxInfluenceDistance
            );

            if (influence > 0) {
              // Направление от плитки к мыши
              const dirX = (mousePos.x - tileCenterX) / (mouseDistance || 1);
              const dirY = (mousePos.y - tileCenterY) / (mouseDistance || 1);

              // Плитка отталкивается от мыши (или притягивается - можно поменять знак)
              const mouseScale = width <= 32 ? 0.25 : 1; // Уменьшаем для маленьких логотипов
              mouseInfluenceX = -dirX * influence * 12 * mouseScale;
              mouseInfluenceY = -dirY * influence * 12 * mouseScale;
              mouseInfluenceZ = influence * 8 * mouseScale;
            }
          }

          // Пульсирование при hover (более мягкое)
          const scalePulse = isHovering
            ? 1 + Math.sin(time * 2 + index) * 0.05
            : 1; // было 0.1

          // Эффект свечения (более мягкий)
          const shouldGlow = isHovering && Math.sin(time * 1.5 + index) > 0.3; // было > 0.5

          return {
            ...tile,
            x: tile.baseX + mouseInfluenceX, // базовая позиция + влияние мыши
            y: tile.baseY + baseFloatY + mouseInfluenceY,
            z: tile.baseZ + baseFloatZ + mouseInfluenceZ,
            rotateX: baseRotateX,
            rotateY: baseRotateY,
            scale: scalePulse,
            glowing: shouldGlow,
            opacity: isHovering
              ? 0.9 + Math.sin(time * 1.5 + index * 0.5) * 0.1
              : 1, // было 0.8 + 0.2
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
  }, [isHovering, mousePos, isExploding, width, height]);

  // Отслеживание мыши для эффекта параллакса
  const handleMouseMove = (e: React.MouseEvent) => {
    if (containerRef.current) {
      const rect = containerRef.current.getBoundingClientRect();
      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;

      setMousePos({
        x: (e.clientX - centerX) / rect.width,
        y: (e.clientY - centerY) / rect.height,
      });
    }
  };

  const handleMouseEnter = () => {
    setIsHovering(true);
  };

  const handleMouseLeave = () => {
    setIsHovering(false);
    setMousePos({ x: 0, y: 0 });
  };

  // Функция разлета и сбора плиток
  const handleClick = () => {
    if (isExploding) return; // Предотвращаем повторные клики

    setIsExploding(true);

    // Сохраняем исходные позиции с учетом масштаба
    const scale = width < 200 ? width / 200 : 1;
    const basePositions = [
      { x: -60, y: -60, z: 0 },
      { x: 0, y: -60, z: 3 },
      { x: 60, y: -60, z: 0 },
      { x: -60, y: 0, z: 5 },
      { x: 0, y: 0, z: 8 },
      { x: 60, y: 0, z: 2 },
      { x: -60, y: 60, z: 0 },
      { x: 0, y: 60, z: 4 },
      { x: 60, y: 60, z: 6 },
    ];

    const originalPositions = basePositions.map((pos) => ({
      x: pos.x * scale,
      y: pos.y * scale,
      z: pos.z * scale,
    }));

    // Создаем случайные новые позиции для разлета (масштабируем для маленьких логотипов)
    const explodeRadius = width <= 32 ? width * 5 : 400;
    const explodePositions = tiles.map(() => ({
      x: (Math.random() - 0.5) * explodeRadius,
      y: (Math.random() - 0.5) * explodeRadius,
      z:
        Math.random() * (width <= 32 ? width * 1.5 : 100) +
        (width <= 32 ? width * 0.5 : 20),
    }));

    // Анимация разлета
    setTiles((prevTiles) =>
      prevTiles.map((tile, index) => ({
        ...tile,
        baseX: explodePositions[index].x,
        baseY: explodePositions[index].y,
        baseZ: explodePositions[index].z,
        scale: 0.5 + Math.random() * 0.8, // Рандомный размер
        glowing: true,
      }))
    );

    // Через 800ms сразу начинаем собирать обратно в рандомном порядке
    setTimeout(() => {
      // Создаем рандомную перестановку позиций
      const shuffledPositions = [...originalPositions].sort(
        () => Math.random() - 0.5
      );

      setTiles((prevTiles) =>
        prevTiles.map((tile, index) => ({
          ...tile,
          baseX: shuffledPositions[index].x,
          baseY: shuffledPositions[index].y,
          baseZ: shuffledPositions[index].z,
          scale: 1,
          glowing: false,
        }))
      );

      // Заканчиваем анимацию взрыва
      setTimeout(() => {
        setIsExploding(false);
      }, 800); // Уменьшено время для более плавного перехода
    }, 800); // Уменьшено время паузы
  };

  return (
    <div
      ref={containerRef}
      className="relative cursor-pointer"
      style={{
        width: width,
        height: height,
        perspective: '1000px',
        perspectiveOrigin: 'center center',
      }}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
    >
      {/* Фоновое свечение */}
      {isHovering && (
        <div
          className="absolute inset-0 rounded-full opacity-30 blur-3xl transition-all duration-1000"
          style={{
            background:
              'radial-gradient(circle, rgba(103,58,183,0.4) 0%, rgba(33,150,243,0.2) 50%, transparent 100%)',
            transform: `scale(${1.5 + Math.abs(mousePos.x) * 0.5})`,
          }}
        />
      )}

      {/* Контейнер для 3D трансформации */}
      <div
        className="absolute inset-0 transition-transform duration-300 ease-out"
        style={{
          transformStyle: 'preserve-3d',
          transform: `rotateX(${mousePos.y * 20}deg) rotateY(${mousePos.x * 20}deg)`,
        }}
      >
        {tiles.map((tile) => (
          <div
            key={tile.id}
            className="absolute transition-all duration-300 ease-out"
            style={{
              left: '50%',
              top: '50%',
              width: width <= 32 ? `${Math.max(width * 0.25, 2)}px` : '50px',
              height: width <= 32 ? `${Math.max(width * 0.25, 2)}px` : '50px',
              transform: `
                translate(-50%, -50%)
                translate3d(${tile.x}px, ${tile.y}px, ${tile.z}px)
                rotateX(${tile.rotateX}deg)
                rotateY(${tile.rotateY}deg)
                rotateZ(${tile.rotateZ}deg)
                scale(${tile.scale})
              `,
              transformStyle: 'preserve-3d',
              opacity: tile.opacity,
            }}
          >
            {/* Тень плитки */}
            <div
              className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 rounded-full bg-black/20 blur-sm transition-all duration-300"
              style={{
                width: `${40 + tile.scale * 10}px`,
                height: `${8 + tile.scale * 2}px`,
                transform: `translateX(-50%) translateZ(-${tile.z + 20}px) scale(${Math.max(0.5, 1 - tile.z * 0.01)})`,
              }}
            />

            {/* Основная плитка */}
            <div
              className={`
                w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-xl
                transition-all duration-300 cursor-pointer relative overflow-hidden
                ${tile.glowing ? 'shadow-2xl' : 'shadow-lg'}
              `}
              style={{
                backgroundColor: tile.color,
                boxShadow: tile.glowing
                  ? `0 0 30px ${tile.color}80, 0 10px 25px rgba(0,0,0,0.3)`
                  : `0 ${5 + tile.z * 0.3}px ${15 + tile.z * 0.5}px rgba(0,0,0,0.2)`,
                transform: `translateZ(5px)`,
              }}
            >
              {/* Градиентное покрытие для 3D эффекта */}
              <div
                className="absolute inset-0 opacity-20"
                style={{
                  background: `linear-gradient(
                    135deg, 
                    rgba(255,255,255,0.8) 0%, 
                    transparent 50%, 
                    rgba(0,0,0,0.3) 100%
                  )`,
                }}
              />

              {/* Анимированные блики */}
              {tile.glowing && (
                <div
                  className="absolute inset-0 animate-pulse"
                  style={{
                    background: `radial-gradient(circle at 30% 30%, rgba(255,255,255,0.4) 0%, transparent 70%)`,
                  }}
                />
              )}

              {/* Иконка */}
              <span
                className="relative z-10 transition-transform duration-300"
                style={{
                  fontSize:
                    width <= 32
                      ? isHovering
                        ? `${width * 0.98}px`
                        : `${width * 0.84}px`
                      : isHovering
                        ? '1.96rem'
                        : '1.68rem',
                  filter: tile.glowing
                    ? 'brightness(1.3) contrast(1.2)'
                    : 'none',
                  lineHeight: '1',
                }}
              >
                {tile.icon}
              </span>
            </div>

            {/* Частицы света */}
            {tile.glowing && (
              <>
                {[...Array(3)].map((_, i) => (
                  <div
                    key={i}
                    className="absolute w-1 h-1 bg-white rounded-full animate-ping"
                    style={{
                      top: `${20 + i * 15}%`,
                      left: `${30 + i * 20}%`,
                      animationDelay: `${i * 200}ms`,
                      animationDuration: '1s',
                    }}
                  />
                ))}
              </>
            )}
          </div>
        ))}
      </div>

      {/* Отражение на поверхности */}
      {isHovering && (
        <div
          className="absolute bottom-0 left-0 right-0 h-16 opacity-20 blur-sm"
          style={{
            background:
              'linear-gradient(to top, rgba(103,58,183,0.3) 0%, transparent 100%)',
            transform: 'scaleY(-0.3) translateY(100%)',
            maskImage: 'linear-gradient(to top, black 0%, transparent 100%)',
          }}
        />
      )}
    </div>
  );
};
