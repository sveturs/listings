export interface UserContact {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: 'pending' | 'accepted' | 'blocked';
  notes?: string;
  added_from_chat_id?: number;
  created_at: string;
  updated_at: string;
  contact_user?: {
    id: number;
    name: string;
    email: string;
  };
}

export interface ContactsResponse {
  contacts: UserContact[];
  total_count: number;
  page: number;
  limit: number;
}

export interface PrivacySettings {
  user_id: number;
  allow_contact_requests: boolean;
  allow_messages_from_contacts_only: boolean;
  created_at: string;
  updated_at: string;
}

export interface AddContactRequest {
  contact_user_id: number;
  notes?: string;
  added_from_chat_id?: number;
}

export interface UpdateContactRequest {
  status: 'accepted' | 'blocked';
  notes?: string;
}

export interface UpdatePrivacySettingsRequest {
  allow_contact_requests?: boolean;
  allow_messages_from_contacts_only?: boolean;
}

export interface ContactStatus {
  are_contacts: boolean;
  user_id: number;
  contact_id: number;
}
