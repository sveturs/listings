// Global type definitions for the project

// Allow importing .svg files as React components
declare module '*.svg' {
  import React from 'react';
  export const ReactComponent: React.FC<React.SVGProps<SVGSVGElement>>;
  const src: string;
  export default src;
}

// Image file imports
declare module '*.png';
declare module '*.jpg';
declare module '*.jpeg';
declare module '*.webp';
declare module '*.gif';

// Allow importing .css files
declare module '*.css' {
  const content: { [className: string]: string };
  export default content;
}

// Allow importing JSON files
declare module '*.json' {
  const value: any;
  export default value;
}

// Custom window properties
interface Window {
  // Add any custom window properties here
  googleMapsLoaded?: boolean;
  ENV?: {
    REACT_APP_BACKEND_URL?: string;
    REACT_APP_AUTH_URL?: string;
    REACT_APP_API_URL?: string;
    [key: string]: string | undefined;
  };
}