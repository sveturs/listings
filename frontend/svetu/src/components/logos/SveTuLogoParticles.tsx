'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';

interface Particle {
  id: number;
  x: number;
  y: number;
  vx: number;
  vy: number;
  life: number;
  maxLife: number;
  size: number;
  color: string;
  alpha: number;
  type: 'spark' | 'glow' | 'trail' | 'magic';
}

interface MagicTile {
  id: number;
  x: number;
  y: number;
  originalX: number;
  originalY: number;
  color: string;
  icon: string;
  energy: number;
  glowing: boolean;
  pulsing: boolean;
  sparkling: boolean;
  magneticField: number;
}

interface SveTuLogoParticlesProps {
  width?: number;
  height?: number;
}

export const SveTuLogoParticles: React.FC<SveTuLogoParticlesProps> = ({
  width = 200,
  height = 200,
}) => {
  const [tiles, setTiles] = useState<MagicTile[]>([]);
  const [particles, setParticles] = useState<Particle[]>([]);
  const [isActive, setIsActive] = useState(false);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
  const animationRef = useRef<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const particleIdCounter = useRef(0);

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

    const initialTiles: MagicTile[] = positions.map((pos, index) => ({
      id: index,
      x: pos.x,
      y: pos.y,
      originalX: pos.x,
      originalY: pos.y,
      color: colors[index],
      icon: icons[index],
      energy: 0,
      glowing: false,
      pulsing: false,
      sparkling: false,
      magneticField: 0,
    }));

    setTiles(initialTiles);
  }, [width]);

  // Создание частиц
  const createParticle = (
    x: number,
    y: number,
    type: Particle['type'],
    color: string
  ): Particle => {
    const angle = Math.random() * Math.PI * 2;
    const speed = Math.random() * 3 + 1;

    return {
      id: particleIdCounter.current++,
      x,
      y,
      vx: Math.cos(angle) * speed,
      vy: Math.sin(angle) * speed,
      life: 1,
      maxLife: type === 'magic' ? 120 : type === 'glow' ? 80 : 60,
      size: type === 'magic' ? Math.random() * 4 + 2 : Math.random() * 3 + 1,
      color,
      alpha: 1,
      type,
    };
  };

  // Генерация частиц от плиток
  const generateParticles = useCallback(
    (tile: MagicTile, count: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const newParticles: Particle[] = [];

      for (let i = 0; i < count; i++) {
        const offsetX = (Math.random() - 0.5) * 50 * scale;
        const offsetY = (Math.random() - 0.5) * 50 * scale;

        if (tile.sparkling && Math.random() < 0.7) {
          newParticles.push(
            createParticle(
              tile.x + offsetX + width / 2,
              tile.y + offsetY + height / 2,
              'spark',
              tile.color
            )
          );
        }

        if (tile.glowing && Math.random() < 0.5) {
          newParticles.push(
            createParticle(
              tile.x + offsetX * 0.5 + width / 2,
              tile.y + offsetY * 0.5 + height / 2,
              'glow',
              tile.color
            )
          );
        }

        if (tile.energy > 0.8 && Math.random() < 0.3) {
          newParticles.push(
            createParticle(
              tile.x + offsetX * 0.3 + width / 2,
              tile.y + offsetY * 0.3 + height / 2,
              'magic',
              '#ffffff'
            )
          );
        }
      }

      return newParticles;
    },
    [width, height]
  );

  // Обновление частиц
  const updateParticles = (particles: Particle[]): Particle[] => {
    return particles
      .map((particle) => {
        // Гравитация и сопротивление
        if (particle.type === 'spark') {
          particle.vy += 0.05; // гравитация
          particle.vx *= 0.98; // сопротивление воздуха
          particle.vy *= 0.98;
        } else if (particle.type === 'glow') {
          particle.vx *= 0.95;
          particle.vy *= 0.95;
        } else if (particle.type === 'magic') {
          // Магические частицы двигаются по синусоиде
          particle.vx += Math.sin(particle.life * 0.1) * 0.1;
          particle.vy += Math.cos(particle.life * 0.1) * 0.1;
        }

        // Обновление позиции
        particle.x += particle.vx;
        particle.y += particle.vy;

        // Уменьшение жизни
        particle.life--;
        particle.alpha = particle.life / particle.maxLife;

        // Изменение размера
        if (particle.type === 'magic') {
          particle.size = 2 + Math.sin(particle.life * 0.2) * 1;
        }

        return particle;
      })
      .filter((particle) => particle.life > 0);
  };

  // Магнитное притяжение к мыши
  const applyMagneticEffect = useCallback(
    (tiles: MagicTile[], mouseX: number, mouseY: number) => {
      const scale = width < 200 ? width / 200 : 1;
      const magneticRadius = 100 * scale;

      return tiles.map((tile) => {
        const dx = mouseX - (tile.originalX + width / 2);
        const dy = mouseY - (tile.originalY + height / 2);
        const distance = Math.sqrt(dx * dx + dy * dy);

        if (distance < magneticRadius && isActive) {
          const force = Math.max(
            0,
            (magneticRadius - distance) / magneticRadius
          );
          const pullX = (dx / distance) * force * 15 * scale;
          const pullY = (dy / distance) * force * 15 * scale;

          return {
            ...tile,
            x: tile.originalX + pullX,
            y: tile.originalY + pullY,
            energy: Math.min(1, force * 2),
            glowing: force > 0.3,
            sparkling: force > 0.6,
            pulsing: force > 0.8,
            magneticField: force,
          };
        } else {
          // Возврат к исходной позиции
          const returnForce = 0.1;
          return {
            ...tile,
            x: tile.x + (tile.originalX - tile.x) * returnForce,
            y: tile.y + (tile.originalY - tile.y) * returnForce,
            energy: Math.max(0, tile.energy - 0.02),
            glowing: tile.energy > 0.2,
            sparkling: tile.energy > 0.4,
            pulsing: tile.energy > 0.6,
            magneticField: Math.max(0, tile.magneticField - 0.05),
          };
        }
      });
    },
    [width, height, isActive]
  );

  // Главный цикл анимации
  useEffect(() => {
    const animate = () => {
      // Обновление плиток
      setTiles((prevTiles) =>
        applyMagneticEffect(prevTiles, mousePos.x, mousePos.y)
      );

      // Обновление частиц
      setParticles((prevParticles) => {
        let updatedParticles = updateParticles(prevParticles);

        // Генерация новых частиц
        if (isActive) {
          tiles.forEach((tile) => {
            if (tile.energy > 0.1) {
              const particleCount = Math.floor(tile.energy * 3);
              const newParticles = generateParticles(tile, particleCount);
              updatedParticles = [...updatedParticles, ...newParticles];
            }
          });
        }

        // Ограничение количества частиц
        return updatedParticles.slice(-200);
      });

      animationRef.current = requestAnimationFrame(animate);
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [isActive, mousePos, tiles, applyMagneticEffect, generateParticles]);

  // Обработчики событий
  const handleMouseMove = (e: React.MouseEvent) => {
    if (containerRef.current) {
      const rect = containerRef.current.getBoundingClientRect();
      setMousePos({
        x: e.clientX - rect.left,
        y: e.clientY - rect.top,
      });
    }
  };

  const handleMouseEnter = () => {
    setIsActive(true);
  };

  const handleMouseLeave = () => {
    setIsActive(false);
    setMousePos({ x: width / 2, y: height / 2 });
  };

  const handleClick = () => {
    // Создание взрывного эффекта
    const explosionParticles: Particle[] = [];
    for (let i = 0; i < 30; i++) {
      explosionParticles.push(
        createParticle(
          mousePos.x,
          mousePos.y,
          'magic',
          `hsl(${Math.random() * 360}, 70%, 60%)`
        )
      );
    }
    setParticles((prev) => [...prev, ...explosionParticles]);
  };

  return (
    <div
      ref={containerRef}
      className="relative cursor-none select-none overflow-hidden"
      style={{ width, height }}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
    >
      {/* Фоновое магическое поле */}
      {isActive && (
        <div
          className="absolute inset-0 opacity-20 transition-opacity duration-500"
          style={{
            background: `radial-gradient(circle at ${mousePos.x}px ${mousePos.y}px, 
              rgba(147, 51, 234, 0.3) 0%, 
              rgba(59, 130, 246, 0.2) 40%, 
              transparent 70%)`,
          }}
        />
      )}

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
                scale(${1 + tile.energy * 0.3})
              `,
            }}
          >
            {/* Магнитное поле вокруг плитки */}
            {tile.magneticField > 0 && (
              <div
                className="absolute inset-0 rounded-full animate-pulse"
                style={{
                  background: `radial-gradient(circle, ${tile.color}20 0%, transparent 70%)`,
                  transform: `scale(${2 + tile.magneticField})`,
                  opacity: tile.magneticField,
                }}
              />
            )}

            {/* Энергетические кольца */}
            {tile.pulsing && (
              <>
                {[...Array(3)].map((_, i) => (
                  <div
                    key={i}
                    className="absolute inset-0 rounded-full border animate-ping"
                    style={{
                      borderColor: `${tile.color}60`,
                      transform: `scale(${1.5 + i * 0.5})`,
                      animationDelay: `${i * 200}ms`,
                      animationDuration: '1s',
                    }}
                  />
                ))}
              </>
            )}

            {/* Основная плитка */}
            <div
              className="w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-lg relative overflow-hidden"
              style={{
                backgroundColor: tile.color,
                boxShadow: tile.glowing
                  ? `0 0 30px ${tile.color}, 0 0 60px ${tile.color}40, 0 10px 25px rgba(0,0,0,0.3)`
                  : '0 5px 15px rgba(0,0,0,0.2)',
                filter: tile.sparkling
                  ? 'brightness(1.4) saturate(1.2)'
                  : 'none',
              }}
            >
              {/* Энергетический градиент */}
              <div
                className="absolute inset-0"
                style={{
                  background: `linear-gradient(
                  45deg, 
                  transparent 0%, 
                  rgba(255,255,255,${tile.energy * 0.3}) 50%, 
                  transparent 100%
                )`,
                  opacity: tile.energy,
                }}
              />

              {/* Магические искры */}
              {tile.sparkling && (
                <>
                  {[...Array(5)].map((_, i) => (
                    <div
                      key={i}
                      className="absolute w-1 h-1 bg-white rounded-full animate-pulse"
                      style={{
                        top: `${Math.random() * 80 + 10}%`,
                        left: `${Math.random() * 80 + 10}%`,
                        animationDelay: `${Math.random() * 1000}ms`,
                        animationDuration: `${500 + Math.random() * 500}ms`,
                        boxShadow: '0 0 4px white',
                      }}
                    />
                  ))}
                </>
              )}

              {/* Иконка */}
              <span
                className="relative z-10 transition-all duration-200"
                style={{
                  fontSize: `${1.68 + tile.energy * 0.56}rem`,
                  transform: `rotate(${tile.magneticField * 10}deg)`,
                  textShadow: tile.glowing ? `0 0 10px ${tile.color}` : 'none',
                  filter: `brightness(${1 + tile.energy * 0.5})`,
                }}
              >
                {tile.icon}
              </span>
            </div>
          </div>
        );
      })}

      {/* Частицы */}
      {particles.map((particle) => (
        <div
          key={particle.id}
          className="absolute pointer-events-none"
          style={{
            left: particle.x,
            top: particle.y,
            width: particle.size,
            height: particle.size,
            backgroundColor: particle.color,
            borderRadius: '50%',
            opacity: particle.alpha,
            transform: 'translate(-50%, -50%)',
            boxShadow:
              particle.type === 'magic'
                ? `0 0 ${particle.size * 2}px ${particle.color}`
                : 'none',
            filter: particle.type === 'glow' ? 'blur(1px)' : 'none',
          }}
        />
      ))}

      {/* Магический курсор */}
      {isActive && (
        <div
          className="absolute w-6 h-6 pointer-events-none transition-opacity duration-300"
          style={{
            left: mousePos.x,
            top: mousePos.y,
            transform: 'translate(-50%, -50%)',
          }}
        >
          <div className="w-full h-full rounded-full bg-white opacity-50 animate-pulse" />
          <div
            className="absolute inset-0 rounded-full border-2 border-purple-400 animate-spin"
            style={{ animationDuration: '2s' }}
          />
        </div>
      )}

      {/* Инструкции */}
      <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full mt-4 text-xs text-gray-500 text-center">
        {isActive
          ? 'Двигайте мышью • Кликните для взрыва'
          : 'Наведите для магии'}
      </div>

      {/* Энергометр */}
      {isActive && (
        <div className="absolute top-0 right-0 transform translate-x-full ml-4">
          <div className="text-xs text-gray-500 mb-2">Энергия</div>
          {tiles.map((tile) => (
            <div key={tile.id} className="flex items-center mb-1">
              <div
                className="w-1 h-3 mr-1 rounded-full"
                style={{
                  backgroundColor: tile.color,
                  opacity: tile.energy,
                }}
              />
              <div className="text-xs text-gray-400">{tile.icon}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
