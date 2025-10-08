import { User } from './auth';
import { C2CItem } from './c2c';

export interface MarketplaceChat {
  id: number;
  listing_id: number;
  storefront_product_id?: number;
  buyer_id: number;
  seller_id: number;
  last_message_at: string;
  created_at: string;
  updated_at: string;
  is_archived: boolean;

  // Дополнительные поля
  buyer?: User;
  seller?: User;
  other_user?: User;
  listing?: C2CItem;
  last_message?: MarketplaceMessage;
  unread_count: number;
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
  metadata?: {
    duration?: number; // для видео в секундах
    pages?: number; // для PDF
    width?: number; // для изображений
    height?: number;
    [key: string]: string | number | boolean | undefined;
  };
  created_at: string;
}

export interface MarketplaceMessage {
  id: number;
  chat_id: number;
  listing_id: number;
  storefront_product_id?: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  is_read: boolean;
  created_at: string;
  updated_at: string;

  // Дополнительные поля
  sender?: User;
  receiver?: User;
  listing?: C2CItem;
  storefront_product?: any; // TODO: добавить тип B2CProduct

  // Мультиязычность
  original_language?: string;
  translations?: Record<string, string>; // { "en": "Hello", "sr": "Здраво" }

  // Поля для вложений
  has_attachments?: boolean;
  attachments_count?: number;
  attachments?: ChatAttachment[];
}

export interface SendMessagePayload {
  listing_id?: number;
  storefront_product_id?: number;
  chat_id?: number;
  receiver_id?: number;
  content: string;
}

export interface GetMessagesParams {
  listing_id?: number;
  chat_id?: number;
  receiver_id?: number;
  page?: number;
  limit?: number;
  signal?: AbortSignal;
}

export interface MarkMessagesReadPayload {
  chat_id: number;
  message_ids: number[];
}

export interface ChatListResponse {
  chats: MarketplaceChat[];
  total: number;
  page: number;
  limit: number;
}

export interface MessagesResponse {
  messages: MarketplaceMessage[];
  total: number;
  page: number;
  limit: number;
}

export interface UnreadCountResponse {
  unread_count: number;
}

// Типы для загрузки файлов
export interface UploadingFile {
  id: string;
  name: string;
  size: number;
  type: string;
  progress: number;
  status: 'pending' | 'uploading' | 'success' | 'error';
  error?: string;
}

export interface FileUploadResponse {
  attachments: ChatAttachment[];
}

// WebSocket события
export interface WSMessage {
  type:
    | 'new_message'
    | 'message_read'
    | 'user_typing'
    | 'user_online'
    | 'user_offline'
    | 'attachment_upload'
    | 'attachment_delete';
  payload: unknown;
}

export interface WSNewMessage {
  type: 'new_message';
  payload: MarketplaceMessage;
}

export interface WSMessageRead {
  type: 'message_read';
  payload: {
    chat_id: number;
    message_ids: number[];
    reader_id: number;
  };
}

export interface WSUserTyping {
  type: 'user_typing';
  payload: {
    chat_id: number;
    user_id: number;
    is_typing: boolean;
  };
}

export interface WSUserStatus {
  type: 'user_online' | 'user_offline';
  payload: {
    user_id: number;
    status: 'online' | 'offline';
    last_seen?: string;
  };
}

// Типы для переводов сообщений
export interface TranslationMetadata {
  translated_from: string;
  translated_to: string;
  translated_at: string;
  cache_hit: boolean;
  provider: string;
}

export interface TranslationResponse {
  message_id: number;
  original_text: string;
  translated_text: string;
  source_language: string;
  target_language: string;
  metadata: TranslationMetadata;
}

export interface GetTranslationParams {
  messageId: number;
  language: string;
}
