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
  FormHelperText,
  Tabs,
  Tab
} from '@mui/material';
import { Upload, FileArchive, FileType, Info, FileCode } from 'lucide-react';
import axios from '../../api/axios';
import CsvStructureInfo from './CsvStructureInfo';
import XmlStructureInfo from './XmlStructureInfo';

const ImportModal = ({ open, onClose, sourceId, onSuccess }) => {
  const { t } = useTranslation(['common', 'marketplace']);
  const [csvFile, setCsvFile] = useState(null);
  const [zipFile, setZipFile] = useState(null);
  const [xmlZipFile, setXmlZipFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [importTab, setImportTab] = useState(0); // 0 - CSV, 1 - XML

  const handleTabChange = (event, newValue) => {
    setImportTab(newValue);
    // Сбрасываем файлы при переключении вкладок
    setCsvFile(null);
    setZipFile(null);
    setXmlZipFile(null);
    setError(null);
    setSuccess(null);
  };

  const handleImport = async () => {
    // Проверяем выбранную вкладку и наличие файлов
    if (importTab === 0 && !csvFile) {
      setError(t('marketplace:store.import.noCsvSelected'));
      return;
    } else if (importTab === 1 && !xmlZipFile) {
      setError(t('marketplace:store.import.noXmlZipSelected'));
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const formData = new FormData();
      
      if (importTab === 0) {
        // CSV импорт
        formData.append('file', csvFile);
        
        // Добавляем ZIP файл с изображениями, если он выбран
        if (zipFile) {
          formData.append('images_zip', zipFile);
        }
      } else {
        // XML импорт
        formData.append('xml_zip', xmlZipFile);
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

  const handleXmlZipFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      if (file.type !== 'application/zip' && !file.name.endsWith('.zip')) {
        setError(t('marketplace:store.import.invalidZipFormat'));
        setXmlZipFile(null);
      } else {
        setError(null);
        setXmlZipFile(file);
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

        <Tabs 
          value={importTab} 
          onChange={handleTabChange} 
          aria-label="import format tabs"
          sx={{ mb: 2 }}
        >
          <Tab label={t('marketplace:store.import.csvImport', { defaultValue: 'CSV Import' })} />
          <Tab label={t('marketplace:store.import.xmlImport', { defaultValue: 'XML Import' })} />
        </Tabs>

        {importTab === 0 ? (
          // CSV импорт
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
        ) : (
          // XML импорт
          <Stack spacing={3} sx={{ mt: 2 }}>
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                {t('marketplace:store.import.selectXmlZipFile', { defaultValue: 'Выберите ZIP-архив с XML файлом' })}
              </Typography>
              
              <Stack direction="row" spacing={2} alignItems="center">
                <Button
                  variant="outlined"
                  component="label"
                  startIcon={<FileCode />}
                  disabled={loading}
                >
                  {t('marketplace:store.import.browseXmlZip', { defaultValue: 'Выбрать ZIP с XML' })}
                  <input
                    type="file"
                    hidden
                    accept=".zip"
                    onChange={handleXmlZipFileChange}
                  />
                </Button>
                
                {xmlZipFile && (
                  <Typography variant="body2" color="text.secondary">
                    {xmlZipFile.name} ({Math.round(xmlZipFile.size / 1024)} KB)
                  </Typography>
                )}
              </Stack>
              
              <FormHelperText>
                {t('marketplace:store.import.xmlZipHelp', {
                  defaultValue: 'ZIP-архив должен содержать XML файл с каталогом товаров. Система автоматически найдет и обработает XML файл внутри архива.'
                })}
              </FormHelperText>
            </Box>
            
            <Alert severity="info" icon={<Info />}>
              {t('marketplace:store.import.xmlFormatInfo', {
                defaultValue: 'XML файл должен содержать элементы <artikal> с информацией о товарах. Поддерживаются изображения, указанные в тегах <slika>.'
              })}
            </Alert>
            
            <Divider />
            
            <Typography variant="subtitle2" gutterBottom>
              {t('marketplace:store.import.alternativeUrlImport', { defaultValue: 'Альтернативно: импорт по URL' })}
            </Typography>
            
            <Alert severity="info">
              {t('marketplace:store.import.xmlUrlInfo', {
                defaultValue: 'Вы также можете добавить прямую ссылку на ZIP-архив с XML при создании источника импорта. Система автоматически распознает ZIP-архив по расширению .zip в URL и обработает его как XML источник.'
              })}
            </Alert>
            
            <XmlStructureInfo />
          </Stack>
        )}

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
            disabled={loading || (importTab === 0 && !csvFile) || (importTab === 1 && !xmlZipFile)}
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