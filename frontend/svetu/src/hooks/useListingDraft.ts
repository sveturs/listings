import { useState, useEffect, useCallback, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { draftService, ListingDraft } from '@/services/draftService';
import { toast } from '@/utils/toast';
import { useTranslations } from 'next-intl';

interface UseListingDraftOptions {
  draftId?: string;
  autoSaveInterval?: number; // в миллисекундах
  onAutoSave?: () => void;
}

export function useListingDraft(options: UseListingDraftOptions = {}) {
  const {
    draftId,
    autoSaveInterval = 30000, // 30 секунд по умолчанию
    onAutoSave,
  } = options;

  const { user } = useAuth();
  const t = useTranslations('createListing');

  const [draft, setDraft] = useState<ListingDraft | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const [hasChanges, setHasChanges] = useState(false);

  const autoSaveTimerRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const lastSavedDraftRef = useRef<ListingDraft | null>(null);

  // Загружаем черновик при монтировании
  useEffect(() => {
    if (!user?.id) {
      setIsLoading(false);
      return;
    }

    if (draftId) {
      // Загружаем существующий черновик
      const existingDraft = draftService.getDraft(draftId, user.id);
      if (existingDraft) {
        setDraft(existingDraft);
        lastSavedDraftRef.current = existingDraft;
        setLastSaved(new Date(existingDraft.metadata.updatedAt));
      } else {
        // Черновик не найден или истек
        toast.error(t('draft.notFound'));
      }
    } else {
      // Создаем новый черновик
      const newDraft = draftService.createDraft(user.id);
      setDraft(newDraft);
      lastSavedDraftRef.current = newDraft;
    }

    setIsLoading(false);
  }, [user?.id, draftId, t]);

  // Функция сохранения черновика
  const saveDraft = useCallback(
    async (showToast = true) => {
      if (!draft || !user?.id) return;

      setIsSaving(true);

      try {
        draftService.saveDraft(draft, user.id);
        lastSavedDraftRef.current = draft;
        setLastSaved(new Date());
        setHasChanges(false);

        if (showToast) {
          toast.success(t('draft.saved'));
        }

        if (onAutoSave) {
          onAutoSave();
        }
      } catch (error) {
        console.error('Error saving draft:', error);
        toast.error(t('draft.saveError'));
      } finally {
        setIsSaving(false);
      }
    },
    [draft, user?.id, t, onAutoSave]
  );

  // Функция обновления черновика
  const updateDraft = useCallback((updates: Partial<ListingDraft>) => {
    setDraft((current) => {
      if (!current) return null;

      const updated = {
        ...current,
        ...updates,
        metadata: {
          ...current.metadata,
          ...updates.metadata,
        },
      };

      // Проверяем, есть ли изменения
      const hasChanges = draftService.hasUnsavedChanges(
        updated,
        lastSavedDraftRef.current || current
      );
      setHasChanges(hasChanges);

      return updated;
    });
  }, []);

  // Обновление данных формы
  const updateFormData = useCallback(
    (data: Partial<ListingDraft['formData']>) => {
      setDraft((current) => {
        if (!current) return current;
        const updated = {
          ...current,
          formData: {
            ...current.formData,
            ...data,
          },
        };

        // Проверяем изменения
        const hasChanges = draftService.hasUnsavedChanges(
          updated,
          lastSavedDraftRef.current || current
        );
        setHasChanges(hasChanges);

        return updated;
      });
    },
    []
  );

  // Обновление атрибутов
  const updateAttributes = useCallback((attributes: Record<string, any>) => {
    setDraft((current) => {
      if (!current) return current;
      const updated = {
        ...current,
        attributes: {
          ...current.attributes,
          ...attributes,
        },
      };

      // Проверяем изменения
      const hasChanges = draftService.hasUnsavedChanges(
        updated,
        lastSavedDraftRef.current || current
      );
      setHasChanges(hasChanges);

      return updated;
    });
  }, []);

  // Обновление текущего шага
  const updateStep = useCallback((step: number) => {
    setDraft((current) => {
      if (!current) return current;
      const updated = {
        ...current,
        metadata: {
          ...current.metadata,
          currentStep: step,
        },
      };

      // Проверяем изменения
      const hasChanges = draftService.hasUnsavedChanges(
        updated,
        lastSavedDraftRef.current || current
      );
      setHasChanges(hasChanges);

      return updated;
    });
  }, []);

  // Автосохранение
  useEffect(() => {
    if (!hasChanges || !autoSaveInterval || isSaving) return;

    // Очищаем предыдущий таймер
    if (autoSaveTimerRef.current) {
      clearTimeout(autoSaveTimerRef.current);
    }

    // Устанавливаем новый таймер
    autoSaveTimerRef.current = setTimeout(() => {
      saveDraft(false); // Не показываем toast при автосохранении
    }, autoSaveInterval);

    return () => {
      if (autoSaveTimerRef.current) {
        clearTimeout(autoSaveTimerRef.current);
      }
    };
  }, [hasChanges, autoSaveInterval, isSaving, saveDraft]);

  // Удаление черновика
  const deleteDraft = useCallback(() => {
    if (!draft || !user?.id) return;

    try {
      draftService.deleteDraft(draft.metadata.id, user.id);
      setDraft(null);
      toast.success(t('draft.deleted'));
    } catch (error) {
      console.error('Error deleting draft:', error);
      toast.error(t('draft.deleteError'));
    }
  }, [draft, user?.id, t]);

  // Экспорт черновика
  const exportDraft = useCallback(() => {
    if (!draft) return;

    const json = draftService.exportDraft(draft);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `draft_${draft.metadata.id}_${new Date().toISOString()}.json`;
    a.click();
    URL.revokeObjectURL(url);

    toast.success(t('draft.exported'));
  }, [draft, t]);

  // Импорт черновика
  const importDraft = useCallback(
    (file: File) => {
      if (!user?.id) return;

      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        const imported = draftService.importDraft(content, user.id);

        if (imported) {
          setDraft(imported);
          lastSavedDraftRef.current = imported;
          toast.success(t('draft.imported'));
        } else {
          toast.error(t('draft.importError'));
        }
      };
      reader.readAsText(file);
    },
    [user?.id, t]
  );

  // Очистка при размонтировании
  useEffect(() => {
    return () => {
      if (autoSaveTimerRef.current) {
        clearTimeout(autoSaveTimerRef.current);
      }

      // Сохраняем при размонтировании если есть изменения
      if (hasChanges && draft && user?.id) {
        draftService.saveDraft(draft, user.id);
      }
    };
  }, [hasChanges, draft, user?.id]);

  // Обработка событий видимости страницы
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.hidden && hasChanges) {
        saveDraft(false);
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [hasChanges, saveDraft]);

  return {
    draft,
    isLoading,
    isSaving,
    hasChanges,
    lastSaved,
    saveDraft,
    updateDraft,
    updateFormData,
    updateAttributes,
    updateStep,
    deleteDraft,
    exportDraft,
    importDraft,
  };
}

/**
 * Hook для получения списка черновиков пользователя
 */
export function useListingDrafts() {
  const { user } = useAuth();
  const [drafts, setDrafts] = useState<
    ReturnType<typeof draftService.getDraftsList>
  >([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!user?.id) {
      setIsLoading(false);
      return;
    }

    // Очищаем истекшие черновики
    draftService.cleanupExpiredDrafts(user.id);

    // Загружаем список
    const list = draftService.getDraftsList(user.id);
    setDrafts(list);
    setIsLoading(false);
  }, [user?.id]);

  const refreshDrafts = useCallback(() => {
    if (!user?.id) return;

    const list = draftService.getDraftsList(user.id);
    setDrafts(list);
  }, [user?.id]);

  return {
    drafts,
    isLoading,
    refreshDrafts,
  };
}
