import React, { useState, useEffect, useRef } from 'react';

interface TilePosition {
  id: number;
  x: number;
  y: number;
  color: string;
  icon: string;
  scale: number;
}

interface IntermediatePoint {
  id: number;
  t: number;
  x: number;
  y: number;
  scale: number;
}

interface SveTuLogoProps {
  width?: number;
  height?: number;
}

const SveTuLogo: React.FC<SveTuLogoProps> = ({ width = 40, height = 40 }) => {
  // Состояния для хранения случайного квадрата и его размера
  const [hovering, setHovering] = useState<boolean>(false);
  const [randomTile, setRandomTile] = useState<number | null>(null);
  const [positions, setPositions] = useState<TilePosition[]>([]);
  const [animatingPositions, setAnimatingPositions] = useState<TilePosition[]>([]);
  const [targetPositions, setTargetPositions] = useState<TilePosition[] | null>(null);
  const [animationProgress, setAnimationProgress] = useState<number>(0);
  const [animationCompleted, setAnimationCompleted] = useState<boolean>(true);
  
  const animationFrameRef = useRef<number | null>(null);
  const animationStartTimeRef = useRef<number | null>(null);
  const touchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const intermediatePositionsRef = useRef<IntermediatePoint[] | null>(null);
  
  // Инициализируем начальные позиции
  useEffect(() => {
    const initialPositions: TilePosition[] = [
        { id: 0, x: 0, y: 0, color: "#2196F3", icon: "🛒", scale: 1 },     // Синий
        { id: 1, x: 74, y: 0, color: "#4CAF50", icon: "🏪", scale: 1 },    // Зеленый
        { id: 2, x: 148, y: 0, color: "#F44336", icon: "🛍️", scale: 1 },   // Красный
        { id: 3, x: 0, y: 74, color: "#FF9800", icon: "📦", scale: 1 },    // Оранжевый
        { id: 4, x: 74, y: 74, color: "#673AB7", icon: "🏠", scale: 1 },   // Фиолетовый
        { id: 5, x: 148, y: 74, color: "#00BCD4", icon: "🤝", scale: 1 },  // Голубой
        { id: 6, x: 0, y: 148, color: "#FFEB3B", icon: "📱", scale: 1 },   // Желтый
        { id: 7, x: 74, y: 148, color: "#607D8B", icon: "💳", scale: 1 },  // Серо-синий
        { id: 8, x: 148, y: 148, color: "#9C27B0", icon: "💰", scale: 1 }  // Пурпурный
      ];
    setPositions(initialPositions);
    setAnimatingPositions(initialPositions);
  }, []);

  // Генерирует промежуточные точки для плавной анимации
  const createIntermediatePositions = (startPos: TilePosition[], endPos: TilePosition[], numPoints = 5): IntermediatePoint[] => {
  if (!startPos || !endPos) return [];
  
  const points: IntermediatePoint[] = [];
  
  const selectedStartTile = startPos.find(tile => tile.id === randomTile);
  const selectedEndTile = endPos.find(tile => tile.id === randomTile);
  
  if (selectedStartTile && selectedEndTile) {
    const dx = selectedEndTile.x - selectedStartTile.x;
    const dy = selectedEndTile.y - selectedStartTile.y;
    
    const perpX = -dy;
    const perpY = dx;
    
    const length = Math.sqrt(perpX * perpX + perpY * perpY) || 1;
    const normalizedPerpX = perpX / length;
    const normalizedPerpY = perpY / length;
    
    for (let i = 1; i < numPoints; i++) {
      const t = i / numPoints;
      // Using a more natural easing function for position
      const smoothT = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
      
      const baseX = selectedStartTile.x + dx * smoothT;
      const baseY = selectedStartTile.y + dy * smoothT;
      
      // Modified deviation curve for more natural arc
      const deviationFactor = 60 * Math.sin(t * Math.PI); 
      
      // Start scaling earlier and use a custom easing for scale
      // This creates a more pleasant "anticipation" and "follow-through" effect
      const scaleEasing = (progress: number): number => {
        // Bell curve-like function for scaling - rises gradually, peaks, then falls gradually
        if (progress < 0.5) {
          // Rise phase: slow start, accelerate
          return selectedStartTile.scale + (progress * 2) * (progress * 2) * (selectedEndTile.scale - selectedStartTile.scale);
        } else {
          // Peak and fall phase: decelerate to final scale
          const fallProgress = (progress - 0.5) * 2;
          return selectedEndTile.scale - (1 - fallProgress) * (1 - fallProgress) * (selectedEndTile.scale - selectedStartTile.scale) * 0.2;
        }
      };
      
      points.push({
        id: selectedStartTile.id,
        t: t,
        x: baseX + normalizedPerpX * deviationFactor,
        y: baseY + normalizedPerpY * deviationFactor,
        // Use the custom easing function for scale
        scale: scaleEasing(t)
      });
    }
  }
  
  return points;
};

// 2. Improved calculateCurrentPositions function for smoother interpolation
const calculateCurrentPositions = (startPos: TilePosition[], endPos: TilePosition[], progress: number): TilePosition[] => {
  // Using cubic bezier easing function for more natural movement
  const easeInOutCubic = (t: number): number => t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
  const easedProgress = easeInOutCubic(progress);
  
  const intermediatePos = intermediatePositionsRef.current || [];
  
  // Better calculation of closest intermediate points
  let prevPointIndex = 0;
  let nextPointIndex = intermediatePos.length - 1;
  
  for (let i = 0; i < intermediatePos.length; i++) {
    if (intermediatePos[i].t > progress) {
      nextPointIndex = i;
      prevPointIndex = Math.max(0, i - 1);
      break;
    }
  }
  
  const prevPoint = intermediatePos[prevPointIndex];
  const nextPoint = intermediatePos[nextPointIndex];
  
  return startPos.map((start, idx) => {
    const end = endPos[idx];
    
    // For non-selected tiles, use improved easing
    if (start.id !== randomTile) {
      return {
        ...start,
        x: start.x + (end.x - start.x) * easedProgress,
        y: start.y + (end.y - start.y) * easedProgress,
        scale: start.scale + (end.scale - start.scale) * easedProgress
      };
    }
    
    // For selected tile, use intermediate points for complex path and scaling
    if (prevPoint && nextPoint && intermediatePos.length > 0) {
      // Improved point interpolation
      const pointProgress = prevPoint.t === nextPoint.t ? 
        0 : 
        easeInOutCubic((progress - prevPoint.t) / (nextPoint.t - prevPoint.t));
      
      const interpolatedX = prevPoint.x + (nextPoint.x - prevPoint.x) * pointProgress;
      const interpolatedY = prevPoint.y + (nextPoint.y - prevPoint.y) * pointProgress;
      
      // Interpolate scale directly between points with smoothing
      const interpolatedScale = prevPoint.scale + (nextPoint.scale - prevPoint.scale) * pointProgress;
      
      return {
        ...start,
        x: interpolatedX,
        y: interpolatedY,
        scale: interpolatedScale
      };
    }
    
    // Fallback with improved easing
    return {
      ...start,
      x: start.x + (end.x - start.x) * easedProgress,
      y: start.y + (end.y - start.y) * easedProgress,
      scale: start.scale + (end.scale - start.scale) * easedProgress
    };
  });
};

  // Очистка анимации при размонтировании компонента
  useEffect(() => {
    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
      if (touchTimeoutRef.current) {
        clearTimeout(touchTimeoutRef.current);
      }
    };
  }, []);

  // Функция для запуска анимации
  const runAnimation = (startPositions: TilePosition[], endPositions: TilePosition[], duration = 1500, afterComplete: (() => void) | null = null) => {
    // Clean up previous animation
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
    }
    
    // Generate more intermediate points for smoother animation
    intermediatePositionsRef.current = createIntermediatePositions(startPositions, endPositions, 15);
    
    const startTime = performance.now();
    setAnimationCompleted(false);
    console.log("dsfsdfsdfsdfdfsdfsdfsdfsdfsdfsdfdsdfsdfsdffsdsdfsdf999")
    const animate = (timestamp: number) => {
      const elapsedTime = timestamp - startTime;
      const progress = Math.min(elapsedTime / duration, 1);
      
      const currentPositions = calculateCurrentPositions(startPositions, endPositions, progress);
      setAnimatingPositions(currentPositions);
      setAnimationProgress(progress);
      
      if (progress < 1) {
        animationFrameRef.current = requestAnimationFrame(animate);
      } else {
        setPositions(endPositions);
        setAnimatingPositions(endPositions);
        setAnimationCompleted(true);
        
        if (afterComplete) {
          afterComplete();
        }
      }
    };
    
    animationFrameRef.current = requestAnimationFrame(animate);
  };

// Общая функция для анимации, используемая для hover и touch
const animateShuffle = () => {
  setHovering(true);

  // Проверяем, что positions не пустой
  if (!positions || positions.length === 0) {
    console.warn("Positions array is not initialized yet.");
    return;
  }

  const tilesWithIcons = positions
    .map((tile, index) => ({ tile, index }));

  // Выбираем случайную плитку
  const randomIconTile = tilesWithIcons[Math.floor(Math.random() * tilesWithIcons.length)];
  
  // Проверяем, что randomIconTile определён
  if (!randomIconTile) {
    console.error("No valid tile selected for animation.");
    return;
  }

  const randomIndex = randomIconTile.index;

  setRandomTile(randomIndex);

  // Создаем копию позиций и перемешиваем
  const newPositions = [...positions].map(pos => ({ ...pos }));

  // Перемешиваем плитки, кроме выбранной
  for (let i = newPositions.length - 1; i > 0; i--) {
    if (i === randomIndex) continue;

    const j = Math.floor(Math.random() * (i + 1));
    if (j === randomIndex) continue;

    const tempX = newPositions[i].x;
    const tempY = newPositions[i].y;
    newPositions[i].x = newPositions[j].x;
    newPositions[i].y = newPositions[j].y;
    newPositions[j].x = tempX;
    newPositions[j].y = tempY;
  }

  const selectedTileOriginalX = newPositions[randomIndex].x;
  const selectedTileOriginalY = newPositions[randomIndex].y;

  const centerX = 74;
  const centerY = 74;

  const currentCenterTileIndex = newPositions.findIndex((tile, idx) =>
    idx !== randomIndex && tile.x === centerX && tile.y === centerY
  );

  if (currentCenterTileIndex !== -1) {
    newPositions[currentCenterTileIndex].x = selectedTileOriginalX;
    newPositions[currentCenterTileIndex].y = selectedTileOriginalY;
  }

  newPositions[randomIndex].x = centerX;
  newPositions[randomIndex].y = centerY;
  newPositions[randomIndex].scale = 3.2;

  setTargetPositions(newPositions);

  const startPositions = [...animatingPositions].map(pos => ({ ...pos }));

  runAnimation(startPositions, newPositions);
};
  
  // Возврат к исходному состоянию
  const resetAnimation = () => {
    setHovering(false);
    
    if (randomTile !== null) {
      // Создаем целевые позиции, где все плитки вернутся к исходному размеру,
      // но останутся на своих текущих позициях
      const resetPositions = (targetPositions || animatingPositions).map(tile => ({
        ...tile, 
        scale: 1
      }));
      
      // Создаем начальные позиции для анимации возврата
      const startPositions = [...animatingPositions].map(pos => ({...pos}));
      
      // Запускаем быструю анимацию возврата к нормальному размеру
      runAnimation(startPositions, resetPositions, 400, () => {
        setRandomTile(null);
      });
    }
  };

  // Обработчики событий для hover
  const handleMouseEnter = () => {
    animateShuffle();
  };

  const handleMouseLeave = () => {
    resetAnimation();
  };
  
  // Обработчики touch-событий для мобильных устройств
  const handleTouchStart = (e: React.TouchEvent) => {
    // Предотвращаем скролл при касании логотипа
    e.preventDefault();
    animateShuffle();
  };
  
  const handleTouchEnd = () => {
    // Добавляем небольшую задержку перед сбросом для мобильных устройств
    touchTimeoutRef.current = setTimeout(() => {
      resetAnimation();
    }, 1800); // Держим анимацию чуть дольше на мобильных устройствах
  };

  return (
    <svg 
      width={width} 
      height={height} 
      viewBox="0 0 256 256" 
      xmlns="http://www.w3.org/2000/svg"
      style={{ 
        display: 'block', 
        cursor: 'pointer',
        overflow: 'visible'
      }}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      <defs>
        <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" floodColor="#000" floodOpacity="0.25"/>
        </filter>
        <filter id="enhancedShadow" x="-30%" y="-30%" width="160%" height="160%">
          <feDropShadow dx="0" dy="4" stdDeviation="4" floodColor="#000" floodOpacity="0.35"/>
        </filter>
      </defs>
      
      <g transform="translate(24,24)">
        {/* Рендерим сначала все обычные плитки */}
        {animatingPositions.map((tile) => {
          // Пропускаем увеличенную плитку, чтобы отрисовать её последней
          if (tile.id === randomTile) return null;
          
          return (
            <g 
              key={tile.id} 
              filter="url(#shadow)" 
              style={{ 
                isolation: 'isolate'
              }}
            >
              <rect 
                x={tile.x} 
                y={tile.y} 
                width={64} 
                height={64} 
                rx={8} 
                fill={tile.color}
              />
              {tile.icon && (
                <text 
                  x={tile.x + 32} 
                  y={tile.y + 38} 
                  fontSize={45} // размер иконок
                  fill={tile.color === "#f7fff7" ? "#000" : "#fff"} 
                  textAnchor="middle"
                  dominantBaseline="middle"
                >
                  {tile.icon}
                </text>
              )}
            </g>
          );
        })}
        
        {/* Рендерим увеличенную плитку последней, чтобы она всегда была сверху */}
        {randomTile !== null && animatingPositions.find(tile => tile.id === randomTile) && (
          <g 
            filter="url(#enhancedShadow)" 
            style={{ 
              isolation: 'isolate',
              zIndex: 100 
            }}
          >
            {(() => {
              const tile = animatingPositions.find(t => t.id === randomTile);
              
              if (!tile) return null;
              
              // Рассчитываем центр для трансформации
              const centerX = tile.x + 32;
              const centerY = tile.y + 32;
              
              // Определяем трансформацию для увеличения вокруг центра
              const transform = `translate(${centerX}, ${centerY}) scale(${tile.scale}) translate(-${centerX}, -${centerY})`;
              
              return (
                <g 
                  key={tile.id} 
                  transform={transform} 
                >
                  <rect 
                    x={tile.x} 
                    y={tile.y} 
                    width={64} 
                    height={64} 
                    rx={8} 
                    fill={tile.color}
                  />
                  {tile.icon && (
                    <text 
                      x={tile.x + 32} 
                      y={tile.y + 38} 
                      fontSize={56} // Увеличил размер иконок в 2 раза
                      fill={tile.color === "#f7fff7" ? "#000" : "#fff"} 
                      textAnchor="middle"
                      dominantBaseline="middle"
                    >
                      {tile.icon}
                    </text>
                  )}
                </g>
              );
            })()}
          </g>
        )}
      </g>
    </svg>
  );
};

export default SveTuLogo;