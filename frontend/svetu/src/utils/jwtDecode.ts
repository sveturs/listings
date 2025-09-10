import { jwtDecode } from 'jwt-decode';
import type { User } from '@/types/auth';

interface JWTPayload {
  iss: string;
  sub: string;
  aud: string[];
  exp: number;
  nbf: number;
  iat: number;
  jti: string;
  user_id: number;
  email: string;
  name: string;
  roles: string[];
  provider: string;
}

export function decodeUserFromToken(token: string): User | null {
  try {
    const decoded = jwtDecode<JWTPayload>(token);

    // Проверяем, не истёк ли токен
    const now = Date.now() / 1000;
    if (decoded.exp < now) {
      console.log('[jwtDecode] Token has expired');
      return null;
    }

    // Преобразуем payload в User объект
    return {
      id: decoded.user_id,
      email: decoded.email,
      name: decoded.name || decoded.email,
      is_admin: decoded.roles?.includes('admin') || false,
      provider: decoded.provider,
    };
  } catch (error) {
    console.error('[jwtDecode] Failed to decode token:', error);
    return null;
  }
}
