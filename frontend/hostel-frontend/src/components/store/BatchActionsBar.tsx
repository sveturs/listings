// frontend/hostel-frontend/src/components/store/BatchActionsBar.tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Button,
  ButtonGroup,
  Chip,
  Tooltip,
  Typography,
  CircularProgress
} from '@mui/material';
import { Trash2, Languages as LanguagesIcon, X } from 'lucide-react';

interface BatchActionsBarProps {
  /** Массив ID выбранных элементов */
  selectedItems: Array<number | string>;
  /** Функция очистки выбора */
  onClearSelection: () => void;
  /** Функция удаления элементов */
  onDelete: () => void;
  /** Функция перевода элементов */
  onTranslate: () => void;
  /** Флаг, указывающий на процесс перевода */
  isTranslating?: boolean;
}

/**
 * Компонент для групповых операций с объявлениями
 */
const BatchActionsBar: React.FC<BatchActionsBarProps> = ({ 
  selectedItems, 
  onClearSelection, 
  onDelete, 
  onTranslate,
  isTranslating = false
}) => {
  const { t } = useTranslation(['marketplace', 'common']);
  
  if (selectedItems.length === 0) {
    return null;
  }

  return (
    <Box
      sx={{
        position: 'sticky',
        bottom: 0,
        left: 0,
        right: 0,
        zIndex: 1000,
        bgcolor: 'background.paper',
        borderTop: 1,
        borderColor: 'divider',
        p: 2,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        boxShadow: 3,
      }}
    >
      <Box display="flex" alignItems="center">
        <Chip 
          label={`${selectedItems.length} ${t('common:common.selected', {defaultValue: 'selected'})}`}
          color="primary"
          onDelete={onClearSelection}
          deleteIcon={<X size={16} />}
          sx={{ mr: 2 }}
        />
        
        <Typography variant="body2" color="text.secondary">
          {t('store.listings.batchActionsHint', {defaultValue: 'Choose an action for selected listings'})}
        </Typography>
      </Box>
      
      <ButtonGroup variant="contained" size="small">
        <Tooltip title={t('marketplace:store.listings.translateSelected', {defaultValue: 'Translate selected'})}>
          <Button 
            startIcon={isTranslating ? <CircularProgress size={16} color="inherit" /> : <LanguagesIcon size={16} />}
            onClick={onTranslate}
            disabled={isTranslating}
            color="primary"
          >
            {t('marketplace:store.listings.translate', {defaultValue: 'Translate'})}
          </Button>
        </Tooltip>
        
        <Tooltip title={t('marketplace:store.listings.deleteSelected', {defaultValue: 'Delete selected'})}>
          <Button 
            startIcon={<Trash2 size={16} />}
            onClick={onDelete}
            color="error"
          >
            {t('marketplace:store.listings.delete', {defaultValue: 'Delete'})}
          </Button>
        </Tooltip>
      </ButtonGroup>
    </Box>
  );
};

export default BatchActionsBar;