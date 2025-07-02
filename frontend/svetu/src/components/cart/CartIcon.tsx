'use client';

import React from 'react';
import { useSelector } from 'react-redux';
import { selectCartItemsCount } from '@/store/slices/cartSlice';

interface CartIconProps {
  onClick: () => void;
  className?: string;
}

export default function CartIcon({ onClick, className = '' }: CartIconProps) {
  const itemsCount = useSelector(selectCartItemsCount);

  return (
    <div className={`indicator cursor-pointer ${className}`} onClick={onClick}>
      {itemsCount > 0 && (
        <span className="indicator-item badge badge-primary badge-sm">
          {itemsCount > 99 ? '99+' : itemsCount}
        </span>
      )}

      <svg
        className="h-6 w-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M3 3h2l.4 2M7 13h10l4-8H5.4m0 0L7 13m0 0l-2.5 5M7 13l2.5 5M17 21a2 2 0 100-4 2 2 0 000 4zM9 21a2 2 0 100-4 2 2 0 000 4z"
        />
      </svg>
    </div>
  );
}
