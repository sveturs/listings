import { getEnv } from '../env';

// Mock next-runtime-env
jest.mock('next-runtime-env', () => ({
  env: jest.fn((key: string) => {
    const mockEnv: Record<string, string> = {
      NEXT_PUBLIC_API_URL: 'http://test-api.com',
      NEXT_PUBLIC_MINIO_URL: 'http://test-minio.com',
    };
    return mockEnv[key];
  }),
}));

describe('getEnv', () => {
  it('should return runtime env variable on client side', () => {
    // Mock browser environment
    global.window = {} as any;

    expect(getEnv('NEXT_PUBLIC_API_URL')).toBe('http://test-api.com');
  });

  it('should return default value if variable not found', () => {
    // Mock browser environment для консистентности
    global.window = {} as any;

    expect(getEnv('NON_EXISTENT_VAR', 'default')).toBe('default');
  });

  it('should return mocked value from env function', () => {
    // Mock browser environment
    global.window = {} as any;

    expect(getEnv('NEXT_PUBLIC_MINIO_URL')).toBe('http://test-minio.com');
  });
});
