'use client';

import React, { useState } from 'react';

interface TooltipProps {
  content: React.ReactNode;
  children: React.ReactNode;
  position?: 'top' | 'bottom' | 'left' | 'right';
  className?: string;
}

export const Tooltip: React.FC<TooltipProps> = ({
  content,
  children,
  position = 'top',
  className = '',
}) => {
  const [isVisible, setIsVisible] = useState(false);

  const positionClasses = {
    top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
    bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
    left: 'right-full top-1/2 -translate-y-1/2 mr-2',
    right: 'left-full top-1/2 -translate-y-1/2 ml-2',
  };

  const arrowClasses = {
    top: 'top-full left-1/2 -translate-x-1/2 border-t-base-content',
    bottom: 'bottom-full left-1/2 -translate-x-1/2 border-b-base-content',
    left: 'left-full top-1/2 -translate-y-1/2 border-l-base-content',
    right: 'right-full top-1/2 -translate-y-1/2 border-r-base-content',
  };

  return (
    <div 
      className="relative inline-block"
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {children}
      
      {isVisible && (
        <div
          className={`
            absolute z-50 
            bg-base-content text-base-100 
            rounded-lg shadow-lg
            min-w-max max-w-xs
            pointer-events-none
            animate-fadeIn
            ${positionClasses[position]}
            ${className}
          `}
        >
          {content}
          
          {/* Arrow */}
          <div
            className={`
              absolute w-0 h-0
              border-4 border-transparent
              ${arrowClasses[position]}
            `}
          />
        </div>
      )}
    </div>
  );
};