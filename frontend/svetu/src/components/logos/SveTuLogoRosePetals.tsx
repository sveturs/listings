'use client';

import React, { useState, useEffect, useRef } from 'react';

interface RosePetal {
  id: number;
  x: number;
  y: number;
  z: number;
  baseX: number;
  baseY: number;
  baseZ: number;
  rotateX: number;
  rotateY: number;
  rotateZ: number;
  color: string;
  icon: string;
  scale: number;
  opacity: number;
  // Специальные параметры для лепестков
  swayPhase: number; // Фаза колыхания
  swayAmplitude: number; // Амплитуда колыхания
  floatSpeed: number; // Скорость полета
  spiralAngle: number; // Угол спирали
  windX: number; // Влияние ветра по X
  windY: number; // Влияние ветра по Y
  isFloating: boolean; // Летит ли лепесток
}

interface SveTuLogoRosePetalsProps {
  width?: number;
  height?: number;
}

export const SveTuLogoRosePetals: React.FC<SveTuLogoRosePetalsProps> = ({
  width = 200,
  height = 200,
}) => {
  const [petals, setPetals] = useState<RosePetal[]>([]);
  const [isScattering, setIsScattering] = useState(false);
  const [windForce, setWindForce] = useState({ x: 0, y: 0 });
  const animationRef = useRef<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const scatterTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Инициализация лепестков
  useEffect(() => {
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

    const initialPetals: RosePetal[] = positions.map((pos, index) => ({
      id: index,
      x: pos.x,
      y: pos.y,
      z: pos.z,
      baseX: pos.x,
      baseY: pos.y,
      baseZ: pos.z,
      rotateX: 0,
      rotateY: 0,
      rotateZ: 0,
      color: colors[index],
      icon: icons[index],
      scale: 1,
      opacity: 1,
      swayPhase: Math.random() * Math.PI * 2,
      swayAmplitude: 10 + Math.random() * 20,
      floatSpeed: 0.5 + Math.random() * 0.5,
      spiralAngle: Math.random() * Math.PI * 2,
      windX: 0,
      windY: 0,
      isFloating: false,
    }));

    setPetals(initialPetals);
  }, [width, height]);

  // Анимация колыхания лепестков
  useEffect(() => {
    const animate = (timestamp: number) => {
      const time = timestamp * 0.001;

      setPetals((prevPetals) =>
        prevPetals.map((petal) => {
          let newX = petal.x;
          let newY = petal.y;
          let newZ = petal.z;
          let newRotateX = 0;
          let newRotateY = 0;
          let newRotateZ = 0;

          if (petal.isFloating) {
            // Сложное движение лепестка в воздухе
            const swayX =
              Math.sin(time * petal.floatSpeed + petal.swayPhase) *
              petal.swayAmplitude;
            const swayY =
              Math.cos(time * petal.floatSpeed * 0.7 + petal.swayPhase) *
              petal.swayAmplitude *
              0.5;
            const swayZ =
              Math.sin(time * petal.floatSpeed * 0.5 + petal.swayPhase) *
              petal.swayAmplitude *
              0.3;

            // Спиральное движение
            petal.spiralAngle += petal.floatSpeed * 0.02;
            const spiralRadius = petal.swayAmplitude * 0.5;
            const spiralX = Math.cos(petal.spiralAngle) * spiralRadius;
            const spiralY = Math.sin(petal.spiralAngle) * spiralRadius;

            // Влияние ветра
            petal.windX += (windForce.x - petal.windX) * 0.1;
            petal.windY += (windForce.y - petal.windY) * 0.1;

            newX = petal.baseX + swayX + spiralX + petal.windX;
            newY = petal.baseY + swayY + spiralY + petal.windY;
            newZ = petal.baseZ + swayZ + Math.abs(Math.sin(time * 0.3)) * 20;

            // Вращение во всех направлениях
            newRotateX =
              Math.sin(time * petal.floatSpeed * 1.2 + petal.id) * 30;
            newRotateY =
              Math.cos(time * petal.floatSpeed * 0.8 + petal.id) * 45;
            newRotateZ =
              Math.sin(time * petal.floatSpeed * 0.6 + petal.id) * 60;
          } else {
            // Легкое покачивание в исходной позиции
            const gentleSway = Math.sin(time * 0.5 + petal.swayPhase) * 2;
            newX = petal.baseX + gentleSway;
            newY = petal.baseY + Math.cos(time * 0.3 + petal.swayPhase) * 1;
            newZ = petal.baseZ + Math.sin(time * 0.4 + petal.swayPhase) * 2;
            newRotateY = Math.sin(time * 0.4 + petal.id) * 5;
          }

          return {
            ...petal,
            x: newX,
            y: newY,
            z: newZ,
            rotateX: newRotateX,
            rotateY: newRotateY,
            rotateZ: newRotateZ,
            scale: petal.isFloating
              ? 0.8 + Math.sin(time * 2 + petal.id) * 0.2
              : 1,
            opacity: petal.isFloating
              ? 0.7 + Math.sin(time * 1.5 + petal.id) * 0.3
              : 1,
          };
        })
      );

      // Изменение направления ветра
      if (isScattering) {
        setWindForce({
          x: Math.sin(time * 0.2) * 30,
          y: Math.cos(time * 0.15) * 20,
        });
      }

      animationRef.current = requestAnimationFrame(animate);
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [isScattering, windForce]);

  // Функция разлета и сбора лепестков
  const handleClick = () => {
    if (isScattering) return;

    setIsScattering(true);

    // Сохраняем исходные позиции
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

    // Запускаем разлет лепестков
    setPetals((prevPetals) =>
      prevPetals.map((petal) => {
        // Случайное направление разлета
        const angle = Math.random() * Math.PI * 2;
        const distance = 100 + Math.random() * 200;

        return {
          ...petal,
          baseX: petal.x + Math.cos(angle) * distance,
          baseY: petal.y + Math.sin(angle) * distance,
          baseZ: Math.random() * 50,
          isFloating: true,
          swayAmplitude: 20 + Math.random() * 40,
          floatSpeed: 0.8 + Math.random() * 0.8,
        };
      })
    );

    // Через 3 секунды начинаем собирать обратно в рандомном порядке
    if (scatterTimeoutRef.current) {
      clearTimeout(scatterTimeoutRef.current);
    }

    scatterTimeoutRef.current = setTimeout(() => {
      // Перемешиваем позиции случайным образом
      const shuffledPositions = [...originalPositions].sort(
        () => Math.random() - 0.5
      );

      setPetals((prevPetals) =>
        prevPetals.map((petal, index) => ({
          ...petal,
          baseX: shuffledPositions[index].x,
          baseY: shuffledPositions[index].y,
          baseZ: shuffledPositions[index].z,
          isFloating: false,
          swayAmplitude: 10 + Math.random() * 20,
          floatSpeed: 0.5 + Math.random() * 0.5,
        }))
      );

      setIsScattering(false);
      setWindForce({ x: 0, y: 0 });
    }, 3000);
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
      onClick={handleClick}
    >
      {/* Эффект ветра */}
      {isScattering && (
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            background: `linear-gradient(${90 + windForce.x}deg, transparent 0%, rgba(100,100,200,0.1) 50%, transparent 100%)`,
            transform: `skewX(${windForce.x * 0.5}deg)`,
            transition: 'all 1s ease',
          }}
        />
      )}

      {/* Контейнер для 3D трансформации */}
      <div
        className="absolute inset-0"
        style={{
          transformStyle: 'preserve-3d',
        }}
      >
        {petals.map((petal) => (
          <div
            key={petal.id}
            className="absolute transition-none"
            style={{
              left: '50%',
              top: '50%',
              width: width <= 32 ? `${Math.max(width * 0.25, 2)}px` : '50px',
              height: width <= 32 ? `${Math.max(width * 0.25, 2)}px` : '50px',
              transform: `
                translate(-50%, -50%)
                translate3d(${petal.x}px, ${petal.y}px, ${petal.z}px)
                rotateX(${petal.rotateX}deg)
                rotateY(${petal.rotateY}deg)
                rotateZ(${petal.rotateZ}deg)
                scale(${petal.scale})
              `,
              transformStyle: 'preserve-3d',
              opacity: petal.opacity,
            }}
          >
            {/* Тень лепестка */}
            {!petal.isFloating && (
              <div
                className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 rounded-full bg-black/20 blur-sm transition-all duration-300"
                style={{
                  width: `${40 + petal.scale * 10}px`,
                  height: `${8 + petal.scale * 2}px`,
                  transform: `translateX(-50%) translateZ(-${petal.z + 20}px) scale(${Math.max(0.5, 1 - petal.z * 0.01)})`,
                }}
              />
            )}

            {/* Основная плитка-лепесток */}
            <div
              className={`
                w-full h-full rounded-lg flex items-center justify-center text-white font-bold
                transition-shadow duration-300 cursor-pointer relative overflow-hidden
                ${petal.isFloating ? 'animate-pulse' : ''}
              `}
              style={{
                backgroundColor: petal.color,
                boxShadow: petal.isFloating
                  ? `0 ${10 + petal.z * 0.3}px ${20 + petal.z * 0.5}px rgba(0,0,0,0.2)`
                  : `0 5px 15px rgba(0,0,0,0.2)`,
                borderRadius: petal.isFloating
                  ? '30% 70% 70% 30% / 30% 30% 70% 70%'
                  : '0.5rem',
                transform: `translateZ(5px)`,
              }}
            >
              {/* Градиент для эффекта лепестка */}
              <div
                className="absolute inset-0"
                style={{
                  background: petal.isFloating
                    ? `radial-gradient(ellipse at 30% 30%, rgba(255,255,255,0.4) 0%, transparent 60%)`
                    : `linear-gradient(135deg, rgba(255,255,255,0.2) 0%, transparent 50%, rgba(0,0,0,0.1) 100%)`,
                  opacity: petal.isFloating ? 0.8 : 0.5,
                }}
              />

              {/* Иконка */}
              <span
                className="relative z-10"
                style={{
                  fontSize: width <= 32 ? `${width * 0.84}px` : '1.68rem',
                  filter: petal.isFloating ? 'blur(0.5px)' : 'none',
                  lineHeight: '1',
                }}
              >
                {petal.icon}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* Инструкция */}
      <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full mt-4 text-xs text-gray-500 text-center whitespace-nowrap">
        {isScattering
          ? 'Лепестки разлетаются...'
          : 'Кликните для полета лепестков'}
      </div>
    </div>
  );
};
