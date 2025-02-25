import React from 'react';

const SveTuLogo = ({ width = 40, height = 40 }) => {
  return (
    <svg 
      width={width} 
      height={height} 
      viewBox="0 0 256 256" 
      xmlns="http://www.w3.org/2000/svg"
      shapeRendering="geometricPrecision"
      textRendering="geometricPrecision"
      style={{ 
        display: 'block', 
        imageRendering: 'optimizeQuality',
        transform: 'translateZ(0)'  // включает аппаратное ускорение
      }}
    >
      <defs>
        <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="4" stdDeviation="4" floodColor="#000" floodOpacity="0.25"/>
        </filter>
      </defs>
      
      {/* Убрали градиентный фон, теперь он прозрачный */}
      
      {/* Сетка из 9 разноцветных квадратиков с меньшими отступами от краев */}
      <g transform="translate(24,24)" filter="url(#shadow)">
        {/* Верхний ряд */}
        <rect x="0" y="0" width="64" height="64" rx="8" fill="#ffcc00"/>
        <text x="32" y="40" fontSize="30" fill="#fff" textAnchor="middle">🛒</text>
        <rect x="74" y="0" width="64" height="64" rx="8" fill="#ff6b6b"/>
        <text x="106" y="40" fontSize="30" fill="#fff" textAnchor="middle">🏪</text>
        <rect x="148" y="0" width="64" height="64" rx="8" fill="#4ecdc4"/>
        {/* Средний ряд */}
        <rect x="0" y="74" width="64" height="64" rx="8" fill="#1a535c"/>
        <text x="32" y="114" fontSize="30" fill="#fff" textAnchor="middle">📦</text>
        <rect x="74" y="74" width="64" height="64" rx="8" fill="#ffe66d"/>
        <rect x="148" y="74" width="64" height="64" rx="8" fill="#f7fff7"/>
        <text x="180" y="114" fontSize="30" fill="#000" textAnchor="middle">🏷️</text>
        {/* Нижний ряд */}
        <rect x="0" y="148" width="64" height="64" rx="8" fill="#ff6b6b"/>
        <rect x="74" y="148" width="64" height="64" rx="8" fill="#4ecdc4"/>
        <text x="106" y="188" fontSize="30" fill="#fff" textAnchor="middle">📍</text>
        <rect x="148" y="148" width="64" height="64" rx="8" fill="#1a535c"/>
        <text x="180" y="188" fontSize="30" fill="#fff" textAnchor="middle">💰</text>
      </g>
    </svg>
  );
};

export default SveTuLogo;