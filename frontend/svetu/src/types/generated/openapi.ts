/**
 * This file is auto-generated from OpenAPI specification
 * Do not edit manually
 */

// User related types
export interface User {
  id: number;
  name: string;
  email: string;
  picture_url?: string;
  google_id?: string;
  provider: 'google' | 'local';
  verified: boolean;
  last_seen?: string;
  created_at: string;
  updated_at: string;
}

export interface UserSession {
  user_id: number;
  email: string;
  provider: string;
  created_at: string;
  expires_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

// Chat related types
export interface ChatMessage {
  id: number;
  chat_id: number;
  listing_id?: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  is_read: boolean;
  has_attachments: boolean;
  attachments_count: number;
  created_at: string;
  updated_at: string;
  sender?: User;
  receiver?: User;
  attachments?: ChatAttachment[];
}

export interface ChatAttachment {
  id: number;
  message_id: number;
  file_type: 'image' | 'video' | 'document';
  file_name: string;
  file_size: number;
  content_type: string;
  public_url: string;
  thumbnail_url?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
}

export interface Chat {
  id: number;
  listing_id?: number;
  buyer_id: number;
  seller_id: number;
  last_message_at: string;
  is_archived: boolean;
  unread_count: number;
  created_at: string;
  updated_at: string;
  buyer?: User;
  seller?: User;
  other_user?: User;
  listing?: C2CListing;
  last_message?: ChatMessage;
}

// Marketplace related types
export interface C2CListing {
  id: number;
  user_id: number;
  category_id: number;
  title: string;
  description: string;
  price: number;
  currency: string;
  location: string;
  status: 'active' | 'sold' | 'archived';
  views_count: number;
  is_promoted: boolean;
  created_at: string;
  updated_at: string;
  user?: User;
  category?: Category;
  images?: MarketplaceImage[];
  attributes?: ListingAttribute[];
}

export interface MarketplaceImage {
  id: number;
  listing_id: number;
  file_path: string;
  file_name: string;
  file_size: number;
  content_type: string;
  is_main: boolean;
  public_url: string;
  created_at: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number;
  level: number;
  path: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ListingAttribute {
  attribute_id: number;
  attribute_name: string;
  value: string | number | boolean;
}

// Contact related types
export interface UserContact {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: 'pending' | 'accepted' | 'blocked';
  notes?: string;
  created_at: string;
  updated_at: string;
  contact_user?: User;
}

export interface PrivacySettings {
  user_id: number;
  allow_contact_requests: boolean;
  allow_messages_from_contacts_only: boolean;
  created_at: string;
  updated_at: string;
}

// API Request/Response types
export interface SendMessageRequest {
  chat_id?: number;
  listing_id?: number;
  receiver_id: number;
  content: string;
}

export interface GetMessagesRequest {
  chat_id?: number;
  listing_id?: number;
  receiver_id?: number;
  page?: number;
  limit?: number;
  signal?: AbortSignal;
}

export interface GetMessagesResponse {
  messages: ChatMessage[];
  total: number;
  page: number;
  limit: number;
}

export interface GetChatsResponse {
  chats: Chat[];
  total: number;
  page: number;
  limit: number;
}

export interface MarkMessagesReadRequest {
  chat_id: number;
  message_ids: number[];
}

export interface UpdateContactRequest {
  status: 'accepted' | 'blocked';
  notes?: string;
}

export interface UpdatePrivacySettingsRequest {
  allow_contact_requests?: boolean;
  allow_messages_from_contacts_only?: boolean;
}

// WebSocket message types
export interface WebSocketMessage {
  type: 'message' | 'typing' | 'read' | 'online' | 'offline';
  payload: unknown;
}
