// frontend/hostel-frontend/src/components/store/XmlStructureInfo.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  Paper,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Divider,
  Alert
} from '@mui/material';
import { Info, ChevronDown } from 'lucide-react';

const XmlStructureInfo = () => {
  const { t } = useTranslation(['marketplace', 'common']);

  // Пример XML структуры
  const xmlExample = `<artikal>
  <id>38538</id>
  <sifra>3037</sifra>
  <naziv><![CDATA[Батарея для LG KU800 900 mAh.]]></naziv>
  <kategorija1><![CDATA[ОБОРУДОВАНИЕ ДЛЯ МОБИЛЬНЫХ]]></kategorija1>
  <kategorija2><![CDATA[БАТАРЕИ]]></kategorija2>
  <kategorija3><![CDATA[БАТАРЕИ OUTLET]]></kategorija3>
  <uvoznik>Digital Vision doo</uvoznik>
  <godinaUvoza>2025.</godinaUvoza>
  <zemljaPorekla>Китай</zemljaPorekla>
  <vpCena>121</vpCena>
  <mpCena>690.0000</mpCena>
  <dostupan>1</dostupan>
  <naAkciji>1</naAkciji>
  <opis><![CDATA[<p>Стандартная сменная батарея для LG.</p><p>Гарантия: 6 месяцев</p>]]></opis>
  <barKod></barKod>
  <slike>
    <slika><![CDATA[https://digitalvision.rs/p/38/38538/bat-lg-ku800.png]]></slika>
  </slike>
</artikal>`;

  // Описание полей XML
  const xmlFields = [
    { name: 'id', description: t('store.import.xml.fields.id', { defaultValue: 'Уникальный идентификатор товара' }) },
    { name: 'sifra', description: t('store.import.xml.fields.sifra', { defaultValue: 'Артикул (код) товара' }) },
    { name: 'naziv', description: t('store.import.xml.fields.naziv', { defaultValue: 'Название товара' }) },
    { name: 'kategorija1', description: t('store.import.xml.fields.kategorija1', { defaultValue: 'Основная категория товара' }) },
    { name: 'kategorija2', description: t('store.import.xml.fields.kategorija2', { defaultValue: 'Подкатегория товара' }) },
    { name: 'kategorija3', description: t('store.import.xml.fields.kategorija3', { defaultValue: 'Дополнительная подкатегория' }) },
    { name: 'vpCena', description: t('store.import.xml.fields.vpCena', { defaultValue: 'Оптовая цена (не используется)' }) },
    { name: 'mpCena', description: t('store.import.xml.fields.mpCena', { defaultValue: 'Розничная цена товара' }) },
    { name: 'dostupan', description: t('store.import.xml.fields.dostupan', { defaultValue: 'Доступность товара (1 - доступен, 0 - недоступен)' }) },
    { name: 'naAkciji', description: t('store.import.xml.fields.naAkciji', { defaultValue: 'Товар на акции (1 - да, 0 - нет)' }) },
    { name: 'opis', description: t('store.import.xml.fields.opis', { defaultValue: 'Описание товара (может содержать HTML-теги)' }) },
    { name: 'slike', description: t('store.import.xml.fields.slike', { defaultValue: 'Блок с изображениями товара' }) },
    { name: 'slika', description: t('store.import.xml.fields.slika', { defaultValue: 'Ссылка на изображение товара (внутри тега slike)' }) },
  ];

  return (
    <Box sx={{ mt: 3 }}>
      <Accordion>
        <AccordionSummary
          expandIcon={<ChevronDown />}
          aria-controls="xml-structure-content"
          id="xml-structure-header"
        >
          <Box display="flex" alignItems="center" gap={1}>
            <Info size={18} />
            <Typography variant="subtitle1">
              {t('store.import.xmlStructureTitle', { defaultValue: 'Структура XML файла для импорта' })}
            </Typography>
          </Box>
        </AccordionSummary>
        <AccordionDetails>
          <Alert severity="info" sx={{ mb: 2 }}>
            {t('store.import.xmlDescription', { defaultValue: 'XML файл должен содержать элементы <artikal> с информацией о товарах. Система автоматически импортирует все товары из XML файла.' })}
          </Alert>
          
          <Typography variant="subtitle2" gutterBottom>
            {t('store.import.xmlExample', { defaultValue: 'Пример структуры элемента товара:' })}
          </Typography>
          
          <Paper variant="outlined" sx={{ p: 2, mb: 3, bgcolor: 'grey.50', overflow: 'auto' }}>
            <pre style={{ margin: 0, fontFamily: 'monospace', fontSize: '0.875rem' }}>
              {xmlExample}
            </pre>
          </Paper>
          
          <Typography variant="subtitle2" gutterBottom>
            {t('store.import.xmlFields', { defaultValue: 'Поля XML файла' })}
          </Typography>
          
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>{t('common:common.name')}</TableCell>
                  <TableCell>{t('common:common.description')}</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {xmlFields.map((field) => (
                  <TableRow key={field.name}>
                    <TableCell>
                      <Typography fontFamily="monospace">
                        {field.name}
                      </Typography>
                    </TableCell>
                    <TableCell>{field.description}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          
          <Divider sx={{ my: 2 }} />
          
          <Typography variant="subtitle2" gutterBottom>
            {t('store.import.xmlImageInfo', { defaultValue: 'Информация об изображениях:' })}
          </Typography>
          
          <Box component="ul" sx={{ pl: 2 }}>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.xmlImageInfo1', { defaultValue: 'Изображения задаются внутри элемента <slike> в виде отдельных элементов <slika>' })}
            </Typography>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.xmlImageInfo2', { defaultValue: 'Каждый элемент <slika> должен содержать URL или путь к изображению' })}
            </Typography>
            <Typography component="li" variant="body2" sx={{ mb: 1 }}>
              {t('store.import.xmlImageInfo3', { defaultValue: 'Первое изображение в списке будет установлено как главное' })}
            </Typography>
            <Typography component="li" variant="body2">
              {t('store.import.xmlImageInfo4', { defaultValue: 'Для элементов с CDATA внутри, содержимое будет корректно обработано' })}
            </Typography>
          </Box>
          
          <Alert severity="warning" sx={{ mt: 2 }}>
            {t('store.import.xmlWarning', { defaultValue: 'Максимальный размер ZIP-архива - 100 МБ. Поддерживаемые форматы XML: UTF-8, Windows-1251.' })}
          </Alert>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};

export default XmlStructureInfo;