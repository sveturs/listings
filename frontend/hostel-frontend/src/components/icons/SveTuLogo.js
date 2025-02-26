import React, { useState, useEffect, useRef } from 'react';

const SveTuLogo = ({ width = 40, height = 40 }) => {
  // Состояния для хранения случайного квадрата и его размера
  const [hovering, setHovering] = useState(false);
  const [randomTile, setRandomTile] = useState(null);
  const [positions, setPositions] = useState([]);
  const [animatingPositions, setAnimatingPositions] = useState([]);
  const [targetPositions, setTargetPositions] = useState(null);
  const [animationProgress, setAnimationProgress] = useState(0);
  const [animationCompleted, setAnimationCompleted] = useState(true);
  
  const animationFrameRef = useRef(null);
  const animationStartTimeRef = useRef(null);
  const touchTimeoutRef = useRef(null);
  const intermediatePositionsRef = useRef(null);
  
  // Инициализируем начальные позиции
  useEffect(() => {
    const initialPositions = [
        { id: 0, x: 0, y: 0, color: "#2196F3", icon: "🛒", scale: 1 },     // Синий
        { id: 1, x: 74, y: 0, color: "#4CAF50", icon: "🏪", scale: 1 },    // Зеленый
        { id: 2, x: 148, y: 0, color: "#F44336", icon: "🔍", scale: 1 },   // Красный
        { id: 3, x: 0, y: 74, color: "#FF9800", icon: "📦", scale: 1 },    // Оранжевый
        { id: 4, x: 74, y: 74, color: "#673AB7", icon: "🏠", scale: 1 },   // Фиолетовый
        { id: 5, x: 148, y: 74, color: "#00BCD4", icon: "🏷️", scale: 1 },  // Голубой
        { id: 6, x: 0, y: 148, color: "#FFEB3B", icon: "📱", scale: 1 },   // Желтый
        { id: 7, x: 74, y: 148, color: "#607D8B", icon: "📍", scale: 1 },  // Серо-синий
        { id: 8, x: 148, y: 148, color: "#9C27B0", icon: "💰", scale: 1 }  // Пурпурный
      ];
    setPositions(initialPositions);
    setAnimatingPositions(initialPositions);
  }, []);

  // Генерирует промежуточные точки для плавной анимации
  const createIntermediatePositions = (startPos, endPos, numPoints = 5) => {
    if (!startPos || !endPos) return [];
    
    // Создаем список промежуточных точек для выбранной плитки
    const points = [];
    
    // Только для выбранной плитки создаем сложный путь
    const selectedStartTile = startPos.find(tile => tile.id === randomTile);
    const selectedEndTile = endPos.find(tile => tile.id === randomTile);
    
    if (selectedStartTile && selectedEndTile) {
      // Средняя точка с большим отклонением для выбранной плитки
      const midX = (selectedStartTile.x + selectedEndTile.x) / 2;
      const midY = (selectedStartTile.y + selectedEndTile.y) / 2;
      
      // Вычисляем вектор пути
      const dx = selectedEndTile.x - selectedStartTile.x;
      const dy = selectedEndTile.y - selectedStartTile.y;
      
      // Создаем перпендикулярный вектор для отклонения
      const perpX = -dy;
      const perpY = dx;
      
      // Нормализуем перпендикулярный вектор
      const length = Math.sqrt(perpX * perpX + perpY * perpY) || 1;
      const normalizedPerpX = perpX / length;
      const normalizedPerpY = perpY / length;
      
      // Создаем промежуточные точки с отклонением
      for (let i = 1; i < numPoints; i++) {
        const t = i / numPoints;
        const smoothT = t * t * (3 - 2 * t); // Плавная функция
        
        // Базовая позиция на прямой
        const baseX = selectedStartTile.x + dx * smoothT;
        const baseY = selectedStartTile.y + dy * smoothT;
        
        // Отклонение, максимальное в середине пути
        const deviationFactor = 50 * Math.sin(t * Math.PI); // Максимальное отклонение в середине
        
        // Добавляем к базовой позиции отклонение в перпендикулярном направлении
        points.push({
          id: selectedStartTile.id,
          t: t,
          x: baseX + normalizedPerpX * deviationFactor,
          y: baseY + normalizedPerpY * deviationFactor,
          scale: 1 + (t > 0.6 ? (t - 0.6) / 0.4 * (selectedEndTile.scale - 1) : 0) // Увеличение во второй половине пути
        });
      }
    }
    
    return points;
  };

  // Функция для расчета текущих позиций на основе прогресса анимации
  const calculateCurrentPositions = (startPos, endPos, progress) => {
    // Кубическая функция плавности для естественного движения
    const easeInOut = t => t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    const easedProgress = easeInOut(progress);
    
    // Находим индекс ближайшей промежуточной точки для выбранной плитки
    const intermediatePos = intermediatePositionsRef.current || [];
    const closestIndex = intermediatePos.findIndex(pos => pos.t >= progress) - 1;
    const prevPoint = intermediatePos[Math.max(0, closestIndex)];
    const nextPoint = intermediatePos[Math.min(intermediatePos.length - 1, closestIndex + 1)];
    
    // Расчет позиций для всех плиток
    return startPos.map((start, idx) => {
      const end = endPos[idx];
      
      // Для невыбранных плиток используем обычную интерполяцию
      if (start.id !== randomTile) {
        return {
          ...start,
          x: start.x + (end.x - start.x) * easedProgress,
          y: start.y + (end.y - start.y) * easedProgress,
          scale: start.scale
        };
      }
      
      // Для выбранной плитки используем промежуточные точки, если они есть
      if (prevPoint && nextPoint && intermediatePos.length > 0) {
        // Интерполируем между точками пути
        const pointProgress = prevPoint.t === nextPoint.t ? 0 : (progress - prevPoint.t) / (nextPoint.t - prevPoint.t);
        const interpolatedX = prevPoint.x + (nextPoint.x - prevPoint.x) * pointProgress;
        const interpolatedY = prevPoint.y + (nextPoint.y - prevPoint.y) * pointProgress;
        
        // Функция для плавного масштабирования во второй половине пути
        const scaleProgress = progress > 0.6 ? (progress - 0.6) / 0.4 : 0;
        const superSmoothScale = scaleProgress * scaleProgress * (3 - 2 * scaleProgress);
        const targetScale = 1 + (end.scale - 1) * superSmoothScale;
        
        return {
          ...start,
          x: interpolatedX,
          y: interpolatedY,
          scale: targetScale
        };
      }
      
      // Fallback на обычную интерполяцию, если нет промежуточных точек
      return {
        ...start,
        x: start.x + (end.x - start.x) * easedProgress,
        y: start.y + (end.y - start.y) * easedProgress,
        scale: progress > 0.6 ? start.scale + (end.scale - start.scale) * ((progress - 0.6) / 0.4) : start.scale
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
  const runAnimation = (startPositions, endPositions, duration = 1500, afterComplete = null) => {
    // Очистим предыдущую анимацию
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
    }
    
    // Генерируем промежуточные точки
    intermediatePositionsRef.current = createIntermediatePositions(startPositions, endPositions);
    
    // Сохраняем время начала анимации
    const startTime = performance.now();
    setAnimationCompleted(false);
    
    const animate = (timestamp) => {
      const elapsedTime = timestamp - startTime;
      const progress = Math.min(elapsedTime / duration, 1);
      
      // Вычисляем текущие позиции на основе прогресса
      const currentPositions = calculateCurrentPositions(startPositions, endPositions, progress);
      setAnimatingPositions(currentPositions);
      setAnimationProgress(progress);
      
      // Проверяем, завершена ли анимация
      if (progress < 1) {
        animationFrameRef.current = requestAnimationFrame(animate);
      } else {
        // Анимация завершена
        setPositions(endPositions);
        setAnimatingPositions(endPositions);
        setAnimationCompleted(true);
        
        if (afterComplete) {
          afterComplete();
        }
      }
    };
    
    // Запускаем анимацию
    animationFrameRef.current = requestAnimationFrame(animate);
  };

  // Общая функция для анимации, используемая для hover и touch
  const animateShuffle = () => {
    setHovering(true);
    
    // Теперь все плитки имеют иконки, так что выбираем из всех
    const tilesWithIcons = positions
      .map((tile, index) => ({ tile, index }));
    
    // Выбираем случайную плитку
    const randomIconTile = tilesWithIcons[Math.floor(Math.random() * tilesWithIcons.length)];
    const randomIndex = randomIconTile.index;
    
    setRandomTile(randomIndex);
    
    // Создаем копию позиций и перемешиваем
    const newPositions = [...positions].map(pos => ({...pos}));
    
    // Перемешиваем плитки, кроме выбранной
    for (let i = newPositions.length - 1; i > 0; i--) {
      if (i === randomIndex) continue; // Пропускаем выбранную плитку
      
      const j = Math.floor(Math.random() * (i + 1));
      if (j === randomIndex) continue; // Пропускаем выбранную плитку
      
      // Меняем местами позиции x и y
      const tempX = newPositions[i].x;
      const tempY = newPositions[i].y;
      newPositions[i].x = newPositions[j].x;
      newPositions[i].y = newPositions[j].y;
      newPositions[j].x = tempX;
      newPositions[j].y = tempY;
    }
    
    // Запоминаем начальное положение выбранной плитки
    const selectedTileOriginalX = newPositions[randomIndex].x;
    const selectedTileOriginalY = newPositions[randomIndex].y;
    
    // Перемещаем выбранную плитку в центр в целевой позиции
    const centerX = 74;
    const centerY = 74;
    
    // Находим плитку, которая оказалась в центре после перемешивания
    const currentCenterTileIndex = newPositions.findIndex((tile, idx) => 
      idx !== randomIndex && tile.x === centerX && tile.y === centerY
    );
    
    // Если какая-то плитка уже в центре, перемещаем её на место выбранной
    if (currentCenterTileIndex !== -1) {
      newPositions[currentCenterTileIndex].x = selectedTileOriginalX;
      newPositions[currentCenterTileIndex].y = selectedTileOriginalY;
    }
    
    // Устанавливаем целевую позицию для выбранной плитки - центр
    newPositions[randomIndex].x = centerX;
    newPositions[randomIndex].y = centerY;
    
    // Увеличиваем выбранную плитку до 3.2
    newPositions[randomIndex].scale = 3.2;
    
    // Сохраняем целевые позиции
    setTargetPositions(newPositions);
    
    // Создаем начальные позиции для анимации
    const startPositions = [...animatingPositions].map(pos => ({...pos}));
    
    // Запускаем анимацию
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
  const handleTouchStart = (e) => {
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
        {animatingPositions.map((tile, index) => {
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