// frontend/hostel-frontend/src/components/store/ImportModal.jsx
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Modal,
  Box,
  Typography,
  Button,
  Alert,
  CircularProgress,
  Divider,
  Paper,
  Stack,
  FormHelperText
} from '@mui/material';
import { Upload, FileArchive, FileType, Info } from 'lucide-react';
import axios from '../../api/axios';
import CsvStructureInfo from './CsvStructureInfo';

const ImportModal = ({ open, onClose, sourceId, onSuccess }) => {
  const { t } = useTranslation(['common', 'marketplace']);
  const [csvFile, setCsvFile] = useState(null);
  const [zipFile, setZipFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const handleImport = async () => {
    if (!csvFile) {
      setError(t('marketplace:store.import.noCsvSelected'));
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const formData = new FormData();
      formData.append('file', csvFile);
      
      // Добавляем ZIP файл с изображениями, если он выбран
      if (zipFile) {
        formData.append('images_zip', zipFile);
      }

      const response = await axios.post(
        `/api/v1/storefronts/import-sources/${sourceId}/run`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }
      );

      if (response.data.success) {
        const result = response.data.data;
        setSuccess(
          t('marketplace:store.import.importSuccess', {
            imported: result.items_imported,
            total: result.items_total,
            status: result.status,
            defaultValue: 'Import completed. Status: {{status}}. Imported: {{imported}}/{{total}} items.'
          })
        );
        
        // Уведомляем родительский компонент об успешном импорте
        if (onSuccess) {
          onSuccess(result);
        }
        
        // Автоматически закрываем после успешного импорта через 3 секунды
        setTimeout(() => {
          onClose();
        }, 3000);
      } else {
        setError(response.data.error || t('marketplace:store.import.unknownError'));
      }
    } catch (err) {
      console.error('Import error:', err);
      setError(err.response?.data?.error || t('marketplace:store.import.serverError'));
    } finally {
      setLoading(false);
    }
  };

  const handleCsvFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      if (file.type !== 'text/csv' && !file.name.endsWith('.csv')) {
        setError(t('marketplace:store.import.invalidCsvFormat'));
        setCsvFile(null);
      } else {
        setError(null);
        setCsvFile(file);
      }
    }
  };

  const handleZipFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      if (file.type !== 'application/zip' && !file.name.endsWith('.zip')) {
        setError(t('marketplace:store.import.invalidZipFormat'));
        setZipFile(null);
      } else {
        setZipFile(file);
      }
    }
  };

  return (
    <Modal
      open={open}
      onClose={() => !loading && onClose()}
      aria-labelledby="import-modal-title"
    >
      <Paper
        sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: { xs: '90%', sm: 600 },
          bgcolor: 'background.paper',
          boxShadow: 24,
          p: 4,
          maxHeight: '90vh',
          overflow: 'auto'
        }}
      >
        <Typography id="import-modal-title" variant="h6" component="h2" gutterBottom>
          {t('marketplace:store.import.importData')}
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            {success}
          </Alert>
        )}

        <Stack spacing={3} sx={{ mt: 2 }}>
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              {t('marketplace:store.import.selectCsvFile')}
            </Typography>
            
            <Stack direction="row" spacing={2} alignItems="center">
              <Button
                variant="outlined"
                component="label"
                startIcon={<FileType />}
                disabled={loading}
              >
                {t('marketplace:store.import.browseCsv')}
                <input
                  type="file"
                  hidden
                  accept=".csv"
                  onChange={handleCsvFileChange}
                />
              </Button>
              
              {csvFile && (
                <Typography variant="body2" color="text.secondary">
                  {csvFile.name} ({Math.round(csvFile.size / 1024)} KB)
                </Typography>
              )}
            </Stack>
            
            <FormHelperText>
              {t('marketplace:store.import.csvHelp', { defaultValue: 'CSV-файл должен содержать обязательные поля: id, title, description, price, category_id' })}
            </FormHelperText>
          </Box>
          
          <Divider />
          
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              {t('marketplace:store.import.selectZipFile')}
            </Typography>
            
            <Stack direction="row" spacing={2} alignItems="center">
              <Button
                variant="outlined"
                component="label"
                startIcon={<FileArchive />}
                disabled={loading}
              >
                {t('marketplace:store.import.browseZip')}
                <input
                  type="file"
                  hidden
                  accept=".zip"
                  onChange={handleZipFileChange}
                />
              </Button>
              
              {zipFile && (
                <Typography variant="body2" color="text.secondary">
                  {zipFile.name} ({Math.round(zipFile.size / 1024)} KB)
                </Typography>
              )}
            </Stack>
            
            <FormHelperText>
              {t('marketplace:store.import.zipHelp', {
                defaultValue: 'Опционально: добавьте ZIP-архив с изображениями. Имена файлов должны совпадать с именами в колонке "images" CSV файла.'
              })}
            </FormHelperText>
          </Box>
          
          <Alert severity="info" icon={<Info />}>
            {t('marketplace:store.import.imageNamingInfo', {
              defaultValue: 'В колонке "images" CSV файла укажите пути к изображениям через запятую. При загрузке ZIP-архива эти пути должны соответствовать файлам внутри архива.'
            })}
          </Alert>
          
          <CsvStructureInfo />
        </Stack>

        <Box sx={{ mt: 4, display: 'flex', justifyContent: 'space-between' }}>
          <Button
            variant="outlined"
            onClick={onClose}
            disabled={loading}
          >
            {t('common:buttons.cancel')}
          </Button>
          
          <Button
            variant="contained"
            color="primary"
            onClick={handleImport}
            disabled={loading || !csvFile}
            startIcon={loading ? <CircularProgress size={24} /> : <Upload />}
          >
            {loading ? t('marketplace:store.import.importing') : t('marketplace:store.import.startImport')}
          </Button>
        </Box>
      </Paper>
    </Modal>
  );
};

export default ImportModal;