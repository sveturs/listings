import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText,
  Divider,
  Stack,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  CircularProgress,
  Tooltip,
  IconButton,
  alpha,
  Grid
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import DownloadIcon from '@mui/icons-material/Download';
import UploadIcon from '@mui/icons-material/Upload';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import FileUploadIcon from '@mui/icons-material/FileUpload';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import CloseIcon from '@mui/icons-material/Close';

// Типы данных
interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number | null;
  icon?: string;
  translations?: Record<string, string>;
  listing_count: number;
}

interface AttributeMapping {
  attribute_id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  is_required: boolean;
  is_enabled: boolean;
  sort_order: number;
  options?: any;
  custom_component?: string;
}

interface CategoryAttributeConfig {
  category_id: number;
  category_name: string;
  attributes: AttributeMapping[];
  exported_at: string;
}

interface CategoryAttributeExporterProps {
  categoryId: number;
  onSuccess?: (message: string) => void;
  onError?: (message: string) => void;
}

const CategoryAttributeExporter: React.FC<CategoryAttributeExporterProps> = ({
  categoryId,
  onSuccess,
  onError
}) => {
  const { t } = useTranslation();
  const [categories, setCategories] = useState<Category[]>([]);
  const [targetCategoryId, setTargetCategoryId] = useState<number | ''>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [openImportDialog, setOpenImportDialog] = useState<boolean>(false);
  const [importContent, setImportContent] = useState<string>('');
  const [importError, setImportError] = useState<string | null>(null);
  const [importSuccess, setImportSuccess] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Загрузка списка категорий при монтировании компонента
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/api/admin/categories');
        const data = response.data.data || response.data;
        
        // Фильтруем текущую категорию из списка для копирования
        const filteredCategories = data.filter((cat: Category) => cat.id !== categoryId);
        setCategories(filteredCategories);
      } catch (err) {
        console.error('Error fetching categories:', err);
        const errorMessage = t('admin.categoryAttributes.fetchCategoriesError');
        setError(errorMessage);
        if (onError) onError(errorMessage);
      }
    };

    fetchCategories();
  }, [categoryId, t, onError]);

  // Обработчик экспорта атрибутов категории
  const handleExportAttributes = async () => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      
      // Запрашиваем атрибуты для текущей категории
      const response = await axios.get(`/api/admin/categories/${categoryId}/attributes/export`);
      const data = response.data.data || response.data;
      
      // Форматируем данные для экспорта
      const exportData: CategoryAttributeConfig = {
        category_id: categoryId,
        category_name: data.category_name,
        attributes: data.attributes,
        exported_at: new Date().toISOString()
      };
      
      // Создаем и скачиваем файл
      const jsonString = JSON.stringify(exportData, null, 2);
      const blob = new Blob([jsonString], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      
      const link = document.createElement('a');
      link.href = url;
      link.download = `category_${categoryId}_attributes_${new Date().toISOString().slice(0, 10)}.json`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      const successMessage = t('admin.categoryAttributes.exportSuccess');
      setSuccess(successMessage);
      if (onSuccess) onSuccess(successMessage);
    } catch (err) {
      console.error('Error exporting category attributes:', err);
      const errorMessage = t('admin.categoryAttributes.exportError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Обработчик открытия диалога импорта
  const handleOpenImportDialog = () => {
    setOpenImportDialog(true);
    setImportContent('');
    setImportError(null);
    setImportSuccess(null);
  };

  // Обработчик закрытия диалога импорта
  const handleCloseImportDialog = () => {
    setOpenImportDialog(false);
  };

  // Обработчик импорта из JSON-файла
  const handleImportFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    setImportError(null);
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const content = e.target?.result as string;
        setImportContent(content);
      } catch (err) {
        console.error('Error reading file:', err);
        setImportError(t('admin.categoryAttributes.invalidFile'));
      }
    };
    reader.readAsText(file);
  };

  // Открыть диалог выбора файла
  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  // Обработчик подтверждения импорта
  const handleConfirmImport = async () => {
    try {
      setLoading(true);
      setImportError(null);
      setImportSuccess(null);
      
      if (!importContent.trim()) {
        setImportError(t('admin.categoryAttributes.noContentToImport'));
        return;
      }
      
      let importData: CategoryAttributeConfig;
      try {
        importData = JSON.parse(importContent);
      } catch (err) {
        setImportError(t('admin.categoryAttributes.invalidJson'));
        return;
      }
      
      // Проверяем формат данных
      if (!importData.attributes || !Array.isArray(importData.attributes)) {
        setImportError(t('admin.categoryAttributes.invalidFormat'));
        return;
      }
      
      // Отправляем данные на сервер
      await axios.post(`/api/admin/categories/${categoryId}/attributes/import`, {
        attributes: importData.attributes
      });
      
      const successMessage = t('admin.categoryAttributes.importSuccess');
      setImportSuccess(successMessage);
      if (onSuccess) onSuccess(successMessage);
      
      // Закрываем диалог после успешного импорта
      setTimeout(() => {
        handleCloseImportDialog();
      }, 1500);
    } catch (err) {
      console.error('Error importing attributes:', err);
      setImportError(t('admin.categoryAttributes.importError'));
    } finally {
      setLoading(false);
    }
  };

  // Обработчик копирования атрибутов между категориями
  const handleCopyAttributes = async () => {
    if (!targetCategoryId) {
      setError(t('admin.categoryAttributes.selectTargetCategory'));
      return;
    }
    
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      
      // Вызываем API для копирования атрибутов
      await axios.post(`/api/admin/categories/${categoryId}/attributes/copy`, {
        target_category_id: targetCategoryId
      });
      
      const successMessage = t('admin.categoryAttributes.copySuccess');
      setSuccess(successMessage);
      if (onSuccess) onSuccess(successMessage);
      
      // Сбрасываем выбранную категорию
      setTargetCategoryId('');
    } catch (err) {
      console.error('Error copying attributes:', err);
      const errorMessage = t('admin.categoryAttributes.copyError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper 
      variant="outlined" 
      sx={{ 
        p: 2, 
        mt: 3, 
        borderRadius: 1,
        position: 'relative'
      }}
    >
      <Typography variant="h6" gutterBottom>
        {t('admin.categoryAttributes.exportImport')}
      </Typography>
      
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        {t('admin.categoryAttributes.exportImportDescription')}
      </Typography>
      
      {error && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }}
          onClose={() => setError(null)}
        >
          {error}
        </Alert>
      )}
      
      {success && (
        <Alert 
          severity="success" 
          sx={{ mb: 2 }}
          onClose={() => setSuccess(null)}
        >
          {success}
        </Alert>
      )}
      
      <Grid container spacing={2}>
        {/* Экспорт атрибутов */}
        <Grid item xs={12} md={4}>
          <Paper 
            variant="outlined" 
            sx={{ 
              p: 2, 
              height: '100%',
              bgcolor: alpha('#e3f2fd', 0.3),
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <Typography variant="subtitle1" gutterBottom fontWeight="medium">
              <DownloadIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
              {t('admin.categoryAttributes.exportAttributes')}
            </Typography>
            
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2, flex: 1 }}>
              {t('admin.categoryAttributes.exportDescription')}
            </Typography>
            
            <Button
              variant="outlined"
              color="primary"
              startIcon={<DownloadIcon />}
              onClick={handleExportAttributes}
              disabled={loading}
              fullWidth
            >
              {loading ? <CircularProgress size={24} /> : t('admin.categoryAttributes.export')}
            </Button>
          </Paper>
        </Grid>
        
        {/* Импорт атрибутов */}
        <Grid item xs={12} md={4}>
          <Paper 
            variant="outlined" 
            sx={{ 
              p: 2, 
              height: '100%',
              bgcolor: alpha('#fff8e1', 0.3),
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <Typography variant="subtitle1" gutterBottom fontWeight="medium">
              <UploadIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
              {t('admin.categoryAttributes.importAttributes')}
            </Typography>
            
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2, flex: 1 }}>
              {t('admin.categoryAttributes.importDescription')}
            </Typography>
            
            <Button
              variant="outlined"
              color="secondary"
              startIcon={<UploadIcon />}
              onClick={handleOpenImportDialog}
              disabled={loading}
              fullWidth
            >
              {loading ? <CircularProgress size={24} /> : t('admin.categoryAttributes.import')}
            </Button>
          </Paper>
        </Grid>
        
        {/* Копирование атрибутов */}
        <Grid item xs={12} md={4}>
          <Paper 
            variant="outlined" 
            sx={{ 
              p: 2, 
              height: '100%',
              bgcolor: alpha('#e8f5e9', 0.3),
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <Typography variant="subtitle1" gutterBottom fontWeight="medium">
              <ContentCopyIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
              {t('admin.categoryAttributes.copyAttributes')}
            </Typography>
            
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {t('admin.categoryAttributes.copyDescription')}
            </Typography>
            
            <FormControl 
              variant="outlined" 
              size="small" 
              fullWidth 
              sx={{ mb: 2 }}
            >
              <InputLabel>{t('admin.categoryAttributes.targetCategory')}</InputLabel>
              <Select
                value={targetCategoryId}
                onChange={(e) => setTargetCategoryId(e.target.value as number)}
                label={t('admin.categoryAttributes.targetCategory')}
                disabled={loading || categories.length === 0}
              >
                <MenuItem value="">{t('admin.common.select')}</MenuItem>
                {categories.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name}
                  </MenuItem>
                ))}
              </Select>
              <FormHelperText>
                {t('admin.categoryAttributes.selectCategoryToCopy')}
              </FormHelperText>
            </FormControl>
            
            <Button
              variant="outlined"
              color="success"
              startIcon={<ContentCopyIcon />}
              onClick={handleCopyAttributes}
              disabled={loading || !targetCategoryId}
              fullWidth
            >
              {loading ? <CircularProgress size={24} /> : t('admin.categoryAttributes.copy')}
            </Button>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Диалог импорта атрибутов */}
      <Dialog 
        open={openImportDialog} 
        onClose={handleCloseImportDialog}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box>
            {t('admin.categoryAttributes.importAttributes')}
            <Typography variant="body2" color="text.secondary">
              {t('admin.categoryAttributes.importDialogDescription')}
            </Typography>
          </Box>
          <IconButton onClick={handleCloseImportDialog} size="small">
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        
        <DialogContent>
          {importError && (
            <Alert 
              severity="error" 
              sx={{ mb: 2 }}
              onClose={() => setImportError(null)}
            >
              {importError}
            </Alert>
          )}
          
          {importSuccess && (
            <Alert 
              severity="success" 
              sx={{ mb: 2 }}
              onClose={() => setImportSuccess(null)}
            >
              {importSuccess}
            </Alert>
          )}
          
          <Stack direction="column" spacing={2} sx={{ mt: 1 }}>
            <Paper 
              variant="outlined" 
              sx={{ 
                p: 2, 
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100px',
                bgcolor: alpha('#f5f5f5', 0.3),
                border: '1px dashed',
                borderColor: 'divider'
              }}
            >
              <input
                type="file"
                accept=".json"
                style={{ display: 'none' }}
                ref={fileInputRef}
                onChange={handleImportFile}
              />
              
              {!importContent ? (
                <Button
                  variant="outlined"
                  color="primary"
                  startIcon={<FileUploadIcon />}
                  onClick={handleUploadClick}
                  disabled={loading}
                >
                  {t('admin.categoryAttributes.selectFile')}
                </Button>
              ) : (
                <Box sx={{ textAlign: 'center' }}>
                  <Typography variant="subtitle1" color="primary" gutterBottom>
                    {t('admin.categoryAttributes.fileUploaded')}
                  </Typography>
                  <Button
                    variant="text"
                    color="primary"
                    size="small"
                    onClick={handleUploadClick}
                  >
                    {t('admin.categoryAttributes.chooseAnotherFile')}
                  </Button>
                </Box>
              )}
            </Paper>
            
            {importContent && (
              <Alert severity="info" icon={<InfoOutlinedIcon />}>
                {t('admin.categoryAttributes.importWarning')}
              </Alert>
            )}
          </Stack>
        </DialogContent>
        
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button 
            onClick={handleCloseImportDialog} 
            variant="outlined" 
            color="inherit"
            disabled={loading}
          >
            {t('admin.common.cancel')}
          </Button>
          <Button
            onClick={handleConfirmImport}
            variant="contained"
            color="primary"
            disabled={loading || !importContent}
            endIcon={loading ? <CircularProgress size={20} /> : <ArrowForwardIcon />}
          >
            {t('admin.categoryAttributes.confirmImport')}
          </Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
};

export default CategoryAttributeExporter;