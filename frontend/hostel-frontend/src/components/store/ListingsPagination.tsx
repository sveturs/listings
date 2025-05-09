// frontend/hostel-frontend/src/components/store/ListingsPagination.tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Pagination,
  FormControl,
  Select,
  MenuItem,
  Typography,
  useTheme,
  useMediaQuery,
  SelectChangeEvent
} from '@mui/material';

interface ListingsPaginationProps {
  /** Общее количество объявлений */
  totalItems: number;
  /** Текущая страница */
  page: number;
  /** Текущий лимит на страницу */
  limit: number;
  /** Обработчик изменения страницы */
  onPageChange: (page: number) => void;
  /** Обработчик изменения количества на страницу */
  onLimitChange: (limit: number) => void;
}

/**
 * Компонент для пагинации и выбора количества объявлений на странице
 */
const ListingsPagination: React.FC<ListingsPaginationProps> = ({ 
  totalItems, 
  page, 
  limit, 
  onPageChange, 
  onLimitChange 
}) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  // Доступные варианты лимита на страницу
  const limitOptions = [20, 50, 100, 200, 500];
  
  // Вычисляем общее количество страниц
  const totalPages = Math.ceil(totalItems / limit);
  
  // Обработчик изменения страницы
  const handlePageChange = (_event: React.ChangeEvent<unknown>, value: number): void => {
    if (onPageChange) {
      onPageChange(value);
    }
  };
  
  // Обработчик изменения лимита
  const handleLimitChange = (event: SelectChangeEvent<number>): void => {
    const newLimit = parseInt(event.target.value as string, 10);
    if (onLimitChange) {
      onLimitChange(newLimit);
    }
  };
  
  // Если объявлений нет или всего одна страница, не показываем пагинацию
  if (totalItems === 0 || totalPages <= 1) {
    return null;
  }
  
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: isMobile ? 'column' : 'row',
        justifyContent: 'space-between',
        alignItems: isMobile ? 'center' : 'flex-end',
        mt: 2,
        pt: 2,
        borderTop: 1,
        borderColor: 'divider',
        gap: 2
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography variant="body2" color="text.secondary">
          {t('store.pagination.itemsPerPage', { defaultValue: 'Items per page:' })}
        </Typography>
        <FormControl size="small" variant="outlined">
          <Select
            value={limit}
            onChange={handleLimitChange}
            sx={{ minWidth: 80 }}
          >
            {limitOptions.map(option => (
              <MenuItem key={option} value={option}>
                {option}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        
        <Typography variant="body2" color="text.secondary" sx={{ ml: 2 }}>
          {t('store.pagination.showing', { 
            start: (page - 1) * limit + 1,
            end: Math.min(page * limit, totalItems),
            total: totalItems,
            defaultValue: 'Showing {{start}}-{{end}} of {{total}} items'
          })}
        </Typography>
      </Box>
      
      <Pagination 
        count={totalPages} 
        page={page} 
        onChange={handlePageChange}
        color="primary"
        size={isMobile ? "small" : "medium"}
        showFirstButton
        showLastButton
      />
    </Box>
  );
};

export default ListingsPagination;