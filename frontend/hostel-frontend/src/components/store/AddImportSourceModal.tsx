// frontend/hostel-frontend/src/components/store/AddImportSourceModal.tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Paper,
  Dialog,
  DialogProps
} from '@mui/material';
import ImportSourceForm from './ImportSourceForm';

export interface ImportSource {
  id?: number | string;
  name: string;
  source_type: 'csv' | 'xml' | 'api';
  url?: string;
  file_path?: string;
  auth_type?: string;
  auth_token?: string;
  auth_username?: string;
  auth_password?: string;
  update_interval?: number;
  storefront_id: number | string;
  mapping_rules?: Record<string, any>;
  [key: string]: any;
}

interface AddImportSourceModalProps {
  /** Флаг открытия модального окна */
  open: boolean;
  /** Функция закрытия модального окна */
  onClose: () => void;
  /** Функция вызываемая при успешном создании/обновлении источника импорта */
  onSuccess: () => void;
  /** ID магазина, для которого создается источник импорта */
  storefrontId: number | string;
  /** Начальные данные для редактирования существующего источника */
  initialData?: ImportSource | null;
}

const AddImportSourceModal: React.FC<AddImportSourceModalProps> = ({ 
  open, 
  onClose, 
  onSuccess, 
  storefrontId, 
  initialData = null 
}) => {
  const { t } = useTranslation(['common', 'marketplace']);

  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="add-import-source-modal-title"
      maxWidth="sm"
      fullWidth
    >
      <Box p={3}>
        <ImportSourceForm
          onClose={onClose}
          onSuccess={onSuccess}
          storefrontId={storefrontId}
          initialData={initialData ? {
            id: initialData.id,
            type: initialData.source_type as 'csv' | 'xml',
            url: initialData.url || '',
            schedule: initialData.update_interval ? 'daily' : '',
            storefront_id: initialData.storefront_id
          } : null}
        />
      </Box>
    </Dialog>
  );
};

export default AddImportSourceModal;