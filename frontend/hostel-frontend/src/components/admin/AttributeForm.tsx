import React, { useState, useEffect, useRef } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControlLabel,
  Checkbox,
  Grid,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Typography,
  IconButton,
  Box,
  Divider,
  Tabs,
  Tab,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Chip
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  ExpandMore as ExpandMoreIcon,
  DragIndicator as DragIndicatorIcon,
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { Attribute } from '../../pages/admin/AttributeManagementPage';

interface AttributeFormProps {
  open: boolean;
  attribute: Attribute | null;
  onSubmit: (formData: Partial<Attribute>) => void;
  onClose: () => void;
}

// Список типов атрибутов
const attributeTypes = [
  { value: "text", label: "Текст" },
  { value: "number", label: "Число" },
  { value: "select", label: "Выбор" },
  { value: "multiselect", label: "Множественный выбор" },
  { value: "boolean", label: "Логический" },
  { value: "range", label: "Диапазон" },
  { value: "date", label: "Дата" },
];

// Список доступных кастомных компонентов
const availableComponents = [
  { value: "CarModelSelector", label: "Выбор модели авто" },
  { value: "ColorPicker", label: "Выбор цвета" },
  { value: "SizeRangeSelector", label: "Выбор размера" },
  { value: "MapLocationPicker", label: "Выбор местоположения на карте" },
];

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`attribute-tabpanel-${index}`}
      aria-labelledby={`attribute-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 2 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const AttributeForm: React.FC<AttributeFormProps> = ({
  open,
  attribute,
  onSubmit,
  onClose
}) => {
  const { t, i18n } = useTranslation();
  const [formData, setFormData] = useState<Partial<Attribute>>({
    name: '',
    display_name: '',
    attribute_type: 'text',
    is_searchable: true,
    is_filterable: true,
    is_required: false,
    sort_order: 0,
    options: { values: [] },
    validation_rules: {},
    translations: {},
    option_translations: {},
    custom_component: '',
  });
  const [languages] = useState<string[]>(['en', 'ru', 'sr']);
  const [tabValue, setTabValue] = useState(0);
  const [optionValues, setOptionValues] = useState<string[]>([]);

  // Инициализация формы при изменении атрибута
  useEffect(() => {
    if (attribute) {
      let options = attribute.options;
      if (typeof options === 'string') {
        try {
          options = JSON.parse(options);
        } catch (e) {
          options = { values: [] };
        }
      }

      let validationRules = attribute.validation_rules;
      if (typeof validationRules === 'string') {
        try {
          validationRules = JSON.parse(validationRules);
        } catch (e) {
          validationRules = {};
        }
      }

      setFormData({
        ...attribute,
        options,
        validation_rules: validationRules,
        translations: attribute.translations || {},
        option_translations: attribute.option_translations || {},
      });

      // Инициализация опций для типов select/multiselect
      if (options && options.values) {
        setOptionValues(options.values);
      } else {
        setOptionValues([]);
      }
    } else {
      setFormData({
        name: '',
        display_name: '',
        attribute_type: 'text',
        is_searchable: true,
        is_filterable: true,
        is_required: false,
        sort_order: 0,
        options: { values: [] },
        validation_rules: {},
        translations: {},
        option_translations: {},
        custom_component: '',
      });
      setOptionValues([]);
    }
  }, [attribute]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSelectChange = (e: any) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    // Если изменился тип атрибута, обновляем опции
    if (name === 'attribute_type') {
      if (value === 'select' || value === 'multiselect') {
        setFormData(prev => ({
          ...prev,
          options: { 
            ...prev.options,
            values: optionValues,
            multiselect: value === 'multiselect'
          }
        }));
      } else if (value === 'number' || value === 'range') {
        setFormData(prev => ({
          ...prev,
          options: { 
            min: 0, 
            max: 100,
            step: 1
          }
        }));
      } else {
        setFormData(prev => ({
          ...prev,
          options: {}
        }));
      }
    }
  };

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: checked
    }));
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleTranslationChange = (lang: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      translations: {
        ...prev.translations,
        [lang]: value
      }
    }));
  };

  const handleOptionTranslationChange = (lang: string, optionIndex: number, value: string) => {
    const optionKey = optionValues[optionIndex];
    setFormData(prev => {
      const prevOptionTranslations = prev.option_translations || {};
      const langTranslations = prevOptionTranslations[lang] || {};
      
      return {
        ...prev,
        option_translations: {
          ...prevOptionTranslations,
          [lang]: {
            ...langTranslations,
            [optionKey]: value
          }
        }
      };
    });
  };

  const handleNumberOptionChange = (field: 'min' | 'max' | 'step', value: string) => {
    const numValue = parseFloat(value);
    if (!isNaN(numValue)) {
      setFormData(prev => ({
        ...prev,
        options: {
          ...prev.options,
          [field]: numValue
        }
      }));
    }
  };

  // Обработка изменений опций для Select/MultiSelect
  const handleAddOption = () => {
    const newOptions = [...optionValues, ''];
    setOptionValues(newOptions);
    
    setFormData(prev => ({
      ...prev,
      options: {
        ...prev.options,
        values: newOptions
      }
    }));
  };

  const handleOptionChange = (index: number, value: string) => {
    const oldValue = optionValues[index];
    const newOptions = [...optionValues];
    newOptions[index] = value;
    setOptionValues(newOptions);
    
    setFormData(prev => {
      // Обновляем значения опций
      const updatedFormData = {
        ...prev,
        options: {
          ...prev.options,
          values: newOptions
        }
      };
      
      // Если изменилось значение существующей опции, обновляем переводы
      if (oldValue && oldValue !== value && prev.option_translations) {
        const updatedTranslations = { ...prev.option_translations };
        
        // Для каждого языка в переводах
        Object.keys(updatedTranslations).forEach(lang => {
          if (updatedTranslations[lang] && updatedTranslations[lang][oldValue]) {
            // Создаем новую запись с новым ключом и сохраняем старое значение перевода
            updatedTranslations[lang][value] = updatedTranslations[lang][oldValue];
            // Удаляем старую запись
            delete updatedTranslations[lang][oldValue];
          }
        });
        
        updatedFormData.option_translations = updatedTranslations;
      }
      
      return updatedFormData;
    });
  };

  const handleRemoveOption = (index: number) => {
    const removedValue = optionValues[index];
    const newOptions = optionValues.filter((_, i) => i !== index);
    setOptionValues(newOptions);
    
    setFormData(prev => {
      // Обновляем значения опций
      const updatedFormData = {
        ...prev,
        options: {
          ...prev.options,
          values: newOptions
        }
      };
      
      // Удаляем переводы для удаленной опции
      if (removedValue && prev.option_translations) {
        const updatedTranslations = { ...prev.option_translations };
        
        Object.keys(updatedTranslations).forEach(lang => {
          if (updatedTranslations[lang] && updatedTranslations[lang][removedValue]) {
            delete updatedTranslations[lang][removedValue];
          }
        });
        
        updatedFormData.option_translations = updatedTranslations;
      }
      
      return updatedFormData;
    });
  };

  const handleOptionDragEnd = (result: any) => {
    if (!result || !result.source || !result.destination) return;

    const items = Array.from(optionValues || []);
    const sourceIndex = result.source.index;
    const destIndex = result.destination.index;

    if (sourceIndex < 0 || destIndex < 0 || sourceIndex >= items.length) {
      return;
    }

    const [reorderedItem] = items.splice(sourceIndex, 1);
    items.splice(destIndex, 0, reorderedItem);

    setOptionValues(items);
    setFormData(prev => ({
      ...prev,
      options: {
        ...prev.options,
        values: items
      }
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Подготавливаем данные для отправки
    const submitData = { ...formData };
    
    // Проверяем тип атрибута и подготавливаем соответствующие опции
    if (formData.attribute_type === 'select' || formData.attribute_type === 'multiselect') {
      submitData.options = {
        values: optionValues.filter(v => v.trim() !== ''),
        multiselect: formData.attribute_type === 'multiselect'
      };
    }
    
    onSubmit(submitData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>
          {attribute ? t('admin.attributes.edit') : t('admin.attributes.create')}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
            <Tabs value={tabValue} onChange={handleTabChange}>
              <Tab label={t('admin.attributes.mainTab')} id="attribute-tab-0" />
              <Tab label={t('admin.attributes.optionsTab')} id="attribute-tab-1" />
              <Tab label={t('admin.attributes.translationsTab')} id="attribute-tab-2" />
              <Tab label={t('admin.attributes.customizationTab')} id="attribute-tab-3" />
            </Tabs>
          </Box>

          {/* Основные настройки */}
          <TabPanel value={tabValue} index={0}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  name="name"
                  label={t('admin.attributes.name')}
                  value={formData.name || ''}
                  onChange={handleChange}
                  required
                  margin="normal"
                  helperText={t('admin.attributes.nameHelp')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  name="display_name"
                  label={t('admin.attributes.displayName')}
                  value={formData.display_name || ''}
                  onChange={handleChange}
                  required
                  margin="normal"
                  helperText={t('admin.attributes.displayNameHelp')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth margin="normal">
                  <InputLabel id="attribute-type-label">
                    {t('admin.attributes.type')}
                  </InputLabel>
                  <Select
                    labelId="attribute-type-label"
                    name="attribute_type"
                    value={formData.attribute_type || 'text'}
                    onChange={handleSelectChange}
                    label={t('admin.attributes.type')}
                    required
                  >
                    {attributeTypes.map(type => (
                      <MenuItem key={type.value} value={type.value}>
                        {type.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  name="sort_order"
                  label={t('admin.attributes.sortOrder')}
                  type="number"
                  value={formData.sort_order || 0}
                  onChange={handleChange}
                  margin="normal"
                  helperText={t('admin.attributes.sortOrderHelp')}
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={formData.is_searchable || false}
                      onChange={handleCheckboxChange}
                      name="is_searchable"
                    />
                  }
                  label={t('admin.attributes.isSearchable')}
                />
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={formData.is_filterable || false}
                      onChange={handleCheckboxChange}
                      name="is_filterable"
                    />
                  }
                  label={t('admin.attributes.isFilterable')}
                />
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={formData.is_required || false}
                      onChange={handleCheckboxChange}
                      name="is_required"
                    />
                  }
                  label={t('admin.attributes.isRequired')}
                />
              </Grid>
            </Grid>
          </TabPanel>

          {/* Настройки опций */}
          <TabPanel value={tabValue} index={1}>
            {(formData.attribute_type === 'number' || formData.attribute_type === 'range') && (
              <Grid container spacing={2}>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    label={t('admin.attributes.minValue')}
                    type="number"
                    value={formData.options?.min ?? 0}
                    onChange={(e) => handleNumberOptionChange('min', e.target.value)}
                    margin="normal"
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    label={t('admin.attributes.maxValue')}
                    type="number"
                    value={formData.options?.max ?? 100}
                    onChange={(e) => handleNumberOptionChange('max', e.target.value)}
                    margin="normal"
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    label={t('admin.attributes.step')}
                    type="number"
                    value={formData.options?.step ?? 1}
                    onChange={(e) => handleNumberOptionChange('step', e.target.value)}
                    margin="normal"
                  />
                </Grid>
              </Grid>
            )}

            {(formData.attribute_type === 'select' || formData.attribute_type === 'multiselect') && (
              <>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="subtitle1">
                    {t('admin.attributes.optionValues')}
                  </Typography>
                  <Button
                    startIcon={<AddIcon />}
                    onClick={handleAddOption}
                    variant="outlined"
                    size="small"
                  >
                    {t('admin.attributes.addOption')}
                  </Button>
                </Box>

                <Box sx={{ mb: 2 }}>
                  {optionValues.map((option, index) => {
                    const itemId = `option-${index}`;
                    return (
                      <Box
                        key={itemId}
                        draggable
                        onDragStart={(e) => {
                          e.dataTransfer.setData('text/plain', String(index));
                          e.dataTransfer.effectAllowed = 'move';
                        }}
                        onDragOver={(e) => {
                          e.preventDefault();
                          e.dataTransfer.dropEffect = 'move';
                        }}
                        onDrop={(e) => {
                          e.preventDefault();
                          const sourceIndex = Number(e.dataTransfer.getData('text/plain'));
                          if (sourceIndex !== index) {
                            // Имитируем результат для handleOptionDragEnd
                            const result = {
                              source: { index: sourceIndex },
                              destination: { index },
                            };
                            handleOptionDragEnd(result);
                          }
                        }}
                        sx={{
                          '&:hover': {
                            backgroundColor: 'rgba(0, 0, 0, 0.04)',
                          },
                          borderRadius: 1,
                          mb: 1,
                        }}
                      >
                        <Grid
                          container
                          spacing={1}
                          alignItems="center"
                          sx={{ mb: 1, p: 0.5 }}
                        >
                          <Grid item sx={{ cursor: 'grab' }}>
                            <DragIndicatorIcon color="action" />
                          </Grid>
                          <Grid item xs>
                            <TextField
                              fullWidth
                              size="small"
                              value={option}
                              onChange={(e) => handleOptionChange(index, e.target.value)}
                              placeholder={t('admin.attributes.optionValue')}
                            />
                          </Grid>
                          <Grid item>
                            <IconButton
                              size="small"
                              onClick={() => handleRemoveOption(index)}
                              color="error"
                            >
                              <DeleteIcon />
                            </IconButton>
                          </Grid>
                        </Grid>
                      </Box>
                    );
                  })}
                </Box>

                {optionValues.length === 0 && (
                  <Typography color="textSecondary" sx={{ mt: 2, mb: 2, textAlign: 'center' }}>
                    {t('admin.attributes.noOptionsAdded')}
                  </Typography>
                )}
              </>
            )}

            {formData.attribute_type !== 'select' && 
            formData.attribute_type !== 'multiselect' && 
            formData.attribute_type !== 'number' && 
            formData.attribute_type !== 'range' && (
              <Typography sx={{ mt: 2, mb: 2, textAlign: 'center' }}>
                {t('admin.attributes.noOptionsForType')}
              </Typography>
            )}
          </TabPanel>

          {/* Переводы */}
          <TabPanel value={tabValue} index={2}>
            <Accordion defaultExpanded>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography variant="subtitle1">
                  {t('admin.attributes.attributeTranslations')}
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                <Grid container spacing={2}>
                  {languages
                    .filter(lang => lang !== i18n.language)
                    .map(lang => (
                      <Grid item xs={12} key={lang}>
                        <TextField
                          fullWidth
                          label={`${t('admin.attributes.displayName')} (${t(`languages.${lang}`)})`}
                          value={formData.translations?.[lang] || ''}
                          onChange={(e) => handleTranslationChange(lang, e.target.value)}
                          margin="normal"
                        />
                      </Grid>
                    ))}
                </Grid>
              </AccordionDetails>
            </Accordion>

            {(formData.attribute_type === 'select' || formData.attribute_type === 'multiselect') && 
              optionValues.length > 0 && (
              <Accordion sx={{ mt: 2 }}>
                <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                  <Typography variant="subtitle1">
                    {t('admin.attributes.optionTranslations')}
                  </Typography>
                </AccordionSummary>
                <AccordionDetails>
                  {languages
                    .filter(lang => lang !== i18n.language)
                    .map(lang => (
                      <Box key={lang} sx={{ mb: 3 }}>
                        <Typography variant="subtitle2" sx={{ mb: 1 }}>
                          {t(`languages.${lang}`)}
                        </Typography>
                        <Divider sx={{ mb: 2 }} />
                        
                        <Grid container spacing={2}>
                          {optionValues
                            .filter(option => option.trim() !== '')
                            .map((option, index) => (
                              <Grid item xs={12} key={`${lang}-${index}`}>
                                <TextField
                                  fullWidth
                                  label={option}
                                  value={formData.option_translations?.[lang]?.[option] || ''}
                                  onChange={(e) => handleOptionTranslationChange(lang, index, e.target.value)}
                                  margin="normal"
                                  size="small"
                                />
                              </Grid>
                            ))}
                        </Grid>
                      </Box>
                    ))}
                </AccordionDetails>
              </Accordion>
            )}
          </TabPanel>

          {/* Кастомизация */}
          <TabPanel value={tabValue} index={3}>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>
              {t('admin.attributes.customComponentSection')}
            </Typography>
            <FormControl fullWidth margin="normal">
              <InputLabel id="custom-component-label">
                {t('admin.attributes.customComponent')}
              </InputLabel>
              <Select
                labelId="custom-component-label"
                name="custom_component"
                value={formData.custom_component || ''}
                onChange={handleSelectChange}
                label={t('admin.attributes.customComponent')}
              >
                <MenuItem value="">
                  <em>{t('admin.attributes.noCustomComponent')}</em>
                </MenuItem>
                {availableComponents.map(component => (
                  <MenuItem key={component.value} value={component.value}>
                    {component.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            
            <Typography variant="body2" color="textSecondary" sx={{ mt: 2 }}>
              {t('admin.attributes.customComponentHelp')}
            </Typography>
            
            {formData.custom_component && (
              <Box sx={{ mt: 3, p: 2, border: '1px dashed', borderColor: 'primary.main', borderRadius: 1 }}>
                <Typography variant="subtitle2" color="primary">
                  {t('admin.attributes.selectedComponent')}:
                </Typography>
                <Chip
                  label={formData.custom_component}
                  color="primary"
                  variant="outlined"
                  sx={{ mt: 1 }}
                />
              </Box>
            )}
          </TabPanel>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} color="primary">
            {t('common.cancel')}
          </Button>
          <Button type="submit" color="primary" variant="contained">
            {t('common.save')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default AttributeForm;