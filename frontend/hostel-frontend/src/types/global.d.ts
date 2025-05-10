// src/types/global.d.ts
interface Window {
  ENV?: {
    REACT_APP_BACKEND_URL?: string;
    REACT_APP_AUTH_URL?: string;
    REACT_APP_MINIO_URL?: string;
    REACT_APP_MAPS_API_KEY?: string;
    [key: string]: string | undefined;
  };
}

// Fix for JSX namespace issue
declare namespace JSX {
  interface Element extends React.ReactElement<any, any> {}
  interface IntrinsicElements {
    [elemName: string]: any;
  }
}

// Additional MUI styled component types
declare namespace React {
  interface DOMAttributes<T> {
    children?: React.ReactNode;
  }
}