import { ListingFormData } from '@/types/marketplace';

export interface DraftMetadata {
  id: string;
  userId: string | number;
  createdAt: string;
  updatedAt: string;
  currentStep: number;
  category?: {
    id: number;
    name: string;
    slug: string;
  };
  title?: string;
  isComplete: boolean;
  expiresAt: string; // Черновики истекают через 30 дней
}

export interface ListingDraft {
  metadata: DraftMetadata;
  formData: Partial<ListingFormData>;
  attributes: Record<string, any>;
  images: Array<{
    id: string;
    url: string;
    file?: File;
  }>;
}

class DraftService {
  private readonly DRAFT_KEY_PREFIX = 'listing_draft';
  private readonly DRAFT_LIST_KEY = 'listing_drafts';
  private readonly DRAFT_EXPIRY_DAYS = 30;

  /**
   * Генерирует уникальный ID для черновика
   */
  private generateDraftId(): string {
    return `draft_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  /**
   * Генерирует дату истечения черновика
   */
  private getExpiryDate(): string {
    const date = new Date();
    date.setDate(date.getDate() + this.DRAFT_EXPIRY_DAYS);
    return date.toISOString();
  }

  /**
   * Создает новый черновик
   */
  createDraft(userId: string | number): ListingDraft {
    const now = new Date().toISOString();
    const draft: ListingDraft = {
      metadata: {
        id: this.generateDraftId(),
        userId,
        createdAt: now,
        updatedAt: now,
        currentStep: 1,
        isComplete: false,
        expiresAt: this.getExpiryDate(),
      },
      formData: {},
      attributes: {},
      images: [],
    };
    return draft;
  }

  /**
   * Сохраняет черновик в localStorage
   */
  saveDraft(draft: ListingDraft, userId: string | number): void {
    if (!userId || draft.metadata.userId !== userId) {
      throw new Error('Unauthorized: Cannot save draft for another user');
    }

    const storageKey = `${this.DRAFT_KEY_PREFIX}_${draft.metadata.id}`;
    const userStorageKey = `svetu_user_${userId}_${storageKey}`;

    // Обновляем метаданные
    draft.metadata.updatedAt = new Date().toISOString();

    // Сохраняем черновик
    localStorage.setItem(userStorageKey, JSON.stringify(draft));

    // Обновляем список черновиков пользователя
    this.updateDraftsList(draft.metadata, userId);
  }

  /**
   * Обновляет список черновиков пользователя
   */
  private updateDraftsList(
    metadata: DraftMetadata,
    userId: string | number
  ): void {
    const listKey = `svetu_user_${userId}_${this.DRAFT_LIST_KEY}`;
    const existingList = this.getDraftsList(userId);

    // Удаляем старую версию если есть
    const filteredList = existingList.filter((m) => m.id !== metadata.id);

    // Добавляем обновленную версию
    filteredList.push(metadata);

    // Сортируем по дате обновления (новые первые)
    filteredList.sort(
      (a, b) =>
        new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
    );

    localStorage.setItem(listKey, JSON.stringify(filteredList));
  }

  /**
   * Получает список черновиков пользователя
   */
  getDraftsList(userId: string | number): DraftMetadata[] {
    const listKey = `svetu_user_${userId}_${this.DRAFT_LIST_KEY}`;

    try {
      const data = localStorage.getItem(listKey);
      if (!data) return [];

      const list: DraftMetadata[] = JSON.parse(data);

      // Фильтруем истекшие черновики
      const now = new Date();
      return list.filter((metadata) => new Date(metadata.expiresAt) > now);
    } catch {
      return [];
    }
  }

  /**
   * Загружает черновик по ID
   */
  getDraft(draftId: string, userId: string | number): ListingDraft | null {
    const storageKey = `${this.DRAFT_KEY_PREFIX}_${draftId}`;
    const userStorageKey = `svetu_user_${userId}_${storageKey}`;

    try {
      const data = localStorage.getItem(userStorageKey);
      if (!data) return null;

      const draft: ListingDraft = JSON.parse(data);

      // Проверяем, что черновик принадлежит пользователю
      if (draft.metadata.userId !== userId) {
        return null;
      }

      // Проверяем, что черновик не истек
      if (new Date(draft.metadata.expiresAt) < new Date()) {
        this.deleteDraft(draftId, userId);
        return null;
      }

      return draft;
    } catch {
      return null;
    }
  }

  /**
   * Удаляет черновик
   */
  deleteDraft(draftId: string, userId: string | number): void {
    const storageKey = `${this.DRAFT_KEY_PREFIX}_${draftId}`;
    const userStorageKey = `svetu_user_${userId}_${storageKey}`;

    // Удаляем черновик
    localStorage.removeItem(userStorageKey);

    // Обновляем список
    const listKey = `svetu_user_${userId}_${this.DRAFT_LIST_KEY}`;
    const existingList = this.getDraftsList(userId);
    const filteredList = existingList.filter((m) => m.id !== draftId);
    localStorage.setItem(listKey, JSON.stringify(filteredList));
  }

  /**
   * Очищает все истекшие черновики для пользователя
   */
  cleanupExpiredDrafts(userId: string | number): number {
    const list = this.getDraftsList(userId);
    const now = new Date();
    let cleaned = 0;

    list.forEach((metadata) => {
      if (new Date(metadata.expiresAt) < now) {
        this.deleteDraft(metadata.id, userId);
        cleaned++;
      }
    });

    return cleaned;
  }

  /**
   * Конвертирует черновик в данные для отправки на сервер
   */
  draftToFormData(draft: ListingDraft): FormData {
    const formData = new FormData();

    // Добавляем основные поля
    Object.entries(draft.formData).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        formData.append(key, String(value));
      }
    });

    // Добавляем атрибуты
    if (Object.keys(draft.attributes).length > 0) {
      formData.append('attributes', JSON.stringify(draft.attributes));
    }

    // Добавляем изображения (только File объекты)
    draft.images.forEach((image) => {
      if (image.file) {
        formData.append(`images`, image.file);
      }
    });

    return formData;
  }

  /**
   * Проверяет, есть ли несохраненные изменения
   */
  hasUnsavedChanges(draft: ListingDraft, lastSaved: ListingDraft): boolean {
    return JSON.stringify(draft) !== JSON.stringify(lastSaved);
  }

  /**
   * Экспортирует черновик в JSON для резервного копирования
   */
  exportDraft(draft: ListingDraft): string {
    const exportData = {
      ...draft,
      exportedAt: new Date().toISOString(),
      version: '1.0',
    };
    return JSON.stringify(exportData, null, 2);
  }

  /**
   * Импортирует черновик из JSON
   */
  importDraft(jsonData: string, userId: string | number): ListingDraft | null {
    try {
      const data = JSON.parse(jsonData);

      // Создаем новый черновик с текущим userId
      const draft: ListingDraft = {
        ...data,
        metadata: {
          ...data.metadata,
          id: this.generateDraftId(),
          userId,
          importedAt: new Date().toISOString(),
          expiresAt: this.getExpiryDate(),
        },
      };

      return draft;
    } catch {
      return null;
    }
  }
}

export const draftService = new DraftService();
