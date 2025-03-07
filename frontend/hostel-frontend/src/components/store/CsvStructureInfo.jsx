// frontend/hostel-frontend/src/components/store/CsvStructureInfo.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
  Divider,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Button
} from '@mui/material';
import { Info, ChevronDown, Download } from 'lucide-react';

const CsvStructureInfo = () => {
  const { t } = useTranslation(['marketplace', 'common']);

  // Пример CSV файла
  const sampleCsvContent = `id;title;description;price;category_id;condition;status;location;latitude;longitude;address_city;address_country;show_on_map;original_language;images
1;Уютная квартира в центре;Красивая квартира с прекрасным видом. Полностью обставлена.;500;1;new;active;Белград, улица Пример 23;44.786568;20.448922;Белград;Сербия;true;ru;apartment1.jpg,apartment2.jpg,apartment3.jpg
2;Ноутбук HP Pavilion;Почти новый ноутбук в отличном состоянии;350;5;used;active;Нови-Сад;45.267136;19.833549;Нови-Сад;Сербия;true;ru;laptop1.jpg,laptop2.jpg
3;Детская коляска;Удобная коляска для малышей;120;8;used;active;;;;;;false;ru;stroller.jpg`;

  // Описание полей CSV
  const csvFields = [
    { name: 'id', required: true, description: t('store.import.fields.id', { defaultValue: 'Уникальный идентификатор объявления (во время импорта используется только для отслеживания)' }) },
    { name: 'title', required: true, description: t('store.import.fields.title', { defaultValue: 'Заголовок объявления' }) },
    { name: 'description', required: true, description: t('store.import.fields.description', { defaultValue: 'Подробное описание объявления' }) },
    { name: 'price', required: true, description: t('store.import.fields.price', { defaultValue: 'Цена (числовое значение)' }) },
    { name: 'category_id', required: true, description: t('store.import.fields.category_id', { defaultValue: 'ID категории (если категория не найдена, используется категория "Прочее")' }) },
    { name: 'condition', required: false, description: t('store.import.fields.condition', { defaultValue: 'Состояние товара: "new" - новый, "used" - б/у (по умолчанию "new")' }) },
    { name: 'status', required: false, description: t('store.import.fields.status', { defaultValue: 'Статус объявления: "active" - активное, "inactive" - неактивное (по умолчанию "active")' }) },
    { name: 'location', required: false, description: t('store.import.fields.location', { defaultValue: 'Местоположение товара (текстовое описание)' }) },
    { name: 'latitude', required: false, description: t('store.import.fields.latitude', { defaultValue: 'Широта (для отображения на карте)' }) },
    { name: 'longitude', required: false, description: t('store.import.fields.longitude', { defaultValue: 'Долгота (для отображения на карте)' }) },
    { name: 'address_city', required: false, description: t('store.import.fields.address_city', { defaultValue: 'Город' }) },
    { name: 'address_country', required: false, description: t('store.import.fields.address_country', { defaultValue: 'Страна' }) },
    { name: 'show_on_map', required: false, description: t('store.import.fields.show_on_map', { defaultValue: 'Отображать на карте: "true" или "false" (по умолчанию "true")' }) },
    { name: 'original_language', required: false, description: t('store.import.fields.original_language', { defaultValue: 'Язык объявления: "ru", "en", "sr" и т.д. (по умолчанию "sr")' }) },
    { name: 'images', required: false, description: t('store.import.fields.images', { defaultValue: 'Пути к изображениям, разделенные запятыми. Могут быть URL-адресами или именами файлов в ZIP-архиве' }) },
  ];

  const downloadSampleCsv = () => {
    const element = document.createElement('a');
    const file = new Blob([sampleCsvContent], { type: 'text/csv' });
    element.href = URL.createObjectURL(file);
    element.download = 'sample_import.csv';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };

  return (
    <Box sx={{ mt: 3 }}>
      <Accordion>
        <AccordionSummary
          expandIcon={<ChevronDown />}
          aria-controls="csv-structure-content"
          id="csv-structure-header"
        >
          <Box display="flex" alignItems="center" gap={1}>
            <Info size={18} />
            <Typography variant="subtitle1">
              {t('store.import.csvStructureTitle', { defaultValue: 'Структура CSV файла для импорта' })}
            </Typography>
          </Box>
        </AccordionSummary>
        <AccordionDetails>
          <Alert severity="info" sx={{ mb: 2 }}>
            {t('store.import.csvDescription', { defaultValue: 'CSV файл должен быть в формате с разделителем ";" (точка с запятой). Первая строка должна содержать заголовки столбцов.' })}
          </Alert>
          
          <Button 
            variant="outlined" 
            startIcon={<Download />} 
            onClick={downloadSampleCsv}
            sx={{ mb: 2 }}
          >
            {t('store.import.downloadSample', { defaultValue: 'Скачать пример CSV файла' })}
          </Button>
          
          <Typography variant="subtitle2" gutterBottom>
            {t('store.import.requiredFields', { defaultValue: 'Поля CSV файла' })}
          </Typography>
          
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>{t('common:common.name')}</TableCell>
                  <TableCell>{t('common:common.required')}</TableCell>
                  <TableCell>{t('common:common.description')}</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {csvFields.map((field) => (
                  <TableRow key={field.name}>
                    <TableCell>
                      <Typography fontFamily="monospace" fontWeight={field.required ? 'bold' : 'normal'}>
                        {field.name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {field.required ? (
                        <Typography color="error">Да</Typography>
                      ) : (
                        <Typography color="text.secondary">Нет</Typography>
                      )}
                    </TableCell>
                    <TableCell>{field.description}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          
          <Divider sx={{ my: 2 }} />
          
          <Typography variant="subtitle2" gutterBottom>
            {t('store.import.imageInstructions', { defaultValue: 'Инструкции по работе с изображениями:' })}
          </Typography>
          
          <Box component="ul" sx={{ pl: 2 }}>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.imageInstruction1', { defaultValue: 'В колонке "images" перечислите пути к изображениям через запятую, например: "image1.jpg,image2.jpg,image3.jpg"' })}
            </Typography>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.imageInstruction2', { defaultValue: 'При загрузке ZIP-архива имена файлов в архиве должны совпадать с именами, указанными в CSV' })}
            </Typography>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.imageInstruction3', { defaultValue: 'Для изображений по URL укажите полные адреса, начинающиеся с "http://" или "https://"' })}
            </Typography>
            <Typography component="li" variant="body2">
              {t('store.import.imageInstruction4', { defaultValue: 'Первое изображение в списке будет установлено как главное' })}
            </Typography>
          </Box>
          
          <Alert severity="warning" sx={{ mt: 2 }}>
            {t('store.import.csvWarning', { defaultValue: 'Максимальный размер CSV файла - 10 МБ, ZIP-архива - 100 МБ. Поддерживаемые форматы изображений: JPEG, PNG, GIF, WebP.' })}
          </Alert>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};

export default CsvStructureInfo;