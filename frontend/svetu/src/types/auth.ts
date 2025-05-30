export interface User {
  id: number;
  name: string;
  email: string;
  provider: string;
  picture_url?: string;
  is_admin?: boolean;
  city?: string;
  country?: string;
  phone?: string;
}

export interface SessionResponse {
  authenticated: boolean;
  user?: User;
}

export interface UserProfile extends User {
  bio?: string;
  notification_email: boolean;
  timezone: string;
  last_seen?: string;
  account_status: string;
  settings?: Record<string, unknown>;
}

export interface RegisterUserRequest {
  name: string;
  email: string;
}

export interface UpdateProfileRequest {
  name?: string;
  phone?: string;
  bio?: string;
  city?: string;
  country?: string;
  notification_email?: boolean;
  timezone?: string;
}
