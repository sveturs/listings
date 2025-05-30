export interface ImageHost {
  protocol: 'http' | 'https';
  hostname: string;
  port?: string;
  pathname: string;
}

export interface Config {
  // API Configuration
  api: {
    url: string;
  };

  // MinIO/Storage Configuration
  storage: {
    minioUrl: string;
    imageHosts: ImageHost[];
    imagePathPattern: string;
  };

  // Environment
  env: {
    isProduction: boolean;
    isDevelopment: boolean;
  };
}

// Env variables schema for validation
export interface EnvVariables {
  NEXT_PUBLIC_API_URL?: string;
  NEXT_PUBLIC_MINIO_URL?: string;
  NEXT_PUBLIC_IMAGE_HOSTS?: string;
  NEXT_PUBLIC_IMAGE_PATH_PATTERN?: string;
  NODE_ENV?: string;
}
