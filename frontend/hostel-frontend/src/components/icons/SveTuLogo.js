import React, { useState, useEffect, useRef } from 'react';

const SveTuLogo = ({ width = 40, height = 40 }) => {
  // Состояния для хранения случайного квадрата и его размера
  const [hovering, setHovering] = useState(false);
  const [randomTile, setRandomTile] = useState(null);
  const [positions, setPositions] = useState([]);
  const [animatingPositions, setAnimatingPositions] = useState([]);
  const [targetPositions, setTargetPositions] = useState(null);
  const animationFrameRef = useRef(null);
  const animationStartTimeRef = useRef(null);
  const touchTimeoutRef = useRef(null);
  
  // Инициализируем начальные позиции
  useEffect(() => {
    const initialPositions = [
      { id: 0, x: 0, y: 0, color: "#ffcc00", icon: "🛒", scale: 1, wobble: 0, originalX: 0, originalY: 0 },
      { id: 1, x: 74, y: 0, color: "#ff6b6b", icon: "🏪", scale: 1, wobble: 0, originalX: 74, originalY: 0 },
      { id: 2, x: 148, y: 0, color: "#4ecdc4", icon: "", scale: 1, wobble: 0, originalX: 148, originalY: 0 },
      { id: 3, x: 0, y: 74, color: "#1a535c", icon: "📦", scale: 1, wobble: 0, originalX: 0, originalY: 74 },
      { id: 4, x: 74, y: 74, color: "#ffe66d", icon: "", scale: 1, wobble: 0, originalX: 74, originalY: 74 },
      { id: 5, x: 148, y: 74, color: "#f7fff7", icon: "🏷️", scale: 1, wobble: 0, originalX: 148, originalY: 74 },
      { id: 6, x: 0, y: 148, color: "#ff6b6b", icon: "", scale: 1, wobble: 0, originalX: 0, originalY: 148 },
      { id: 7, x: 74, y: 148, color: "#4ecdc4", icon: "📍", scale: 1, wobble: 0, originalX: 74, originalY: 148 },
      { id: 8, x: 148, y: 148, color: "#1a535c", icon: "💰", scale: 1, wobble: 0, originalX: 148, originalY: 148 }
    ];
    setPositions(initialPositions);
    setAnimatingPositions(initialPositions);
  }, []);

  // Функция для анимации перемещения плиток
  const animateTiles = (startTime, fromPositions, toPositions, duration = 800) => {
    const currentTime = performance.now();
    const elapsedTime = currentTime - startTime;
    const progress = Math.min(elapsedTime / duration, 1);
    
    // Кубическая функция плавности для естественного движения (ease-in-out)
    const easeInOut = t => t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    const easedProgress = easeInOut(progress);
    
    // Расчет текущих позиций на основе прогресса анимации
    const currentPositions = fromPositions.map((startPos, index) => {
      const targetPos = toPositions[index];
      
      // Промежуточный scale - плавное увеличение выбранной плитки
      const currentScale = startPos.id === randomTile 
        ? startPos.scale + (targetPos.scale - startPos.scale) * easedProgress
        : startPos.scale;
        
      return {
        ...startPos,
        x: startPos.x + (targetPos.x - startPos.x) * easedProgress,
        y: startPos.y + (targetPos.y - startPos.y) * easedProgress,
        scale: currentScale
      };
    });
    
    setAnimatingPositions(currentPositions);
    
    if (progress < 1) {
      animationFrameRef.current = requestAnimationFrame(() => {
        animateTiles(startTime, fromPositions, toPositions, duration);
      });
    } else {
      // Анимация завершена, сохраняем конечные позиции
      setPositions(toPositions);
      setAnimatingPositions(toPositions);
      setTargetPositions(null);
    }
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

  // Общая функция для анимации, используемая и для hover и для touch
  const animateShuffle = () => {
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
    }
    
    setHovering(true);
    
    // Выбираем только плитки с иконками
    const tilesWithIcons = positions
      .map((tile, index) => ({ tile, index }))
      .filter(item => item.tile.icon !== "");
    
    // Если нет плиток с иконками, ничего не делаем
    if (tilesWithIcons.length === 0) return;
    
    // Выбираем случайную плитку с иконкой
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
    
    // Увеличиваем выбранную плитку в 2.5 раза
    newPositions[randomIndex].scale = 3.2;
    
    // Сохраняем целевые позиции
    setTargetPositions(newPositions);
    
    // Создаем начальные позиции для анимации, где выбранная плитка еще в исходном положении
    const startPositions = [...animatingPositions].map(pos => ({...pos}));
    
    // Запускаем анимацию перемещения
    animationStartTimeRef.current = performance.now();
    animateTiles(animationStartTimeRef.current, startPositions, newPositions);
  };
  
  // Возврат к исходному состоянию
  const resetAnimation = () => {
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
    }
    
    setHovering(false);
    
    if (randomTile !== null) {
      // Создаем целевые позиции, где все плитки вернутся к исходному размеру
      const resetPositions = [...animatingPositions].map(tile => ({
        ...tile, 
        scale: 1
      }));
      
      // Запускаем анимацию возврата
      animationStartTimeRef.current = performance.now();
      animateTiles(animationStartTimeRef.current, animatingPositions, resetPositions, 500);
      setRandomTile(null);
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
    }, 1500); // Держим анимацию чуть дольше на мобильных устройствах
  };

  // Создаем четкий SVG с улучшенным разрешением
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
                  fontSize={28} 
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
                      fontSize={28} 
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