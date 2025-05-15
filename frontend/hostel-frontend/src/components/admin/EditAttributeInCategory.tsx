import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
  Typography,
  TextField,
  FormControlLabel,
  Checkbox,
  Button,
  Grid,
  CircularProgress,
  Alert,
  Divider,
  Tabs,
  Tab,
  IconButton,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Paper,
  InputLabel,
  FormControl,
  Select,
  MenuItem,
  ListItem,
  List,
  ListItemText,
  Tooltip
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  ExpandMore as ExpandMoreIcon,
  DragIndicator as DragIndicatorIcon
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';

interface Attribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  options?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_required: boolean;
  sort_order: number;
  custom_component?: string;
  created_at: string;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
  validation_rules?: any;
}

interface CategoryAttributeMapping {
  category_id: number;
  attribute_id: number;
  is_required: boolean;
  is_enabled: boolean;
  sort_order: number;
  custom_component?: string;
  attribute?: Attribute;
  hint?: string; 
  description?: string;
  translations?: Record<string, any>;
  options?: any;
  unit?: string;
  unit_translations?: Record<string, string>;
}

interface EditAttributeInCategoryProps {
  open: boolean;
  categoryId: number;
  attributeId: number;
  onClose: () => void;
  onUpdate: () => void;
  onError?: (message: string) => void;
}

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

const EditAttributeInCategory: React.FC<EditAttributeInCategoryProps> = ({
  open,
  categoryId,
  attributeId,
  onClose,
  onUpdate,
  onError
}) => {
  const { t, i18n } = useTranslation();
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [mapping, setMapping] = useState<CategoryAttributeMapping | null>(null);
  const [originalAttribute, setOriginalAttribute] = useState<Attribute | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [languages] = useState<string[]>(['en', 'ru', 'sr']);
  const [optionValues, setOptionValues] = useState<string[]>([]);
  const [attributeOptions, setAttributeOptions] = useState<string[]>([]);

  // Fetch attribute mapping details
  useEffect(() => {
    if (open && categoryId && attributeId) {
      fetchAttributeMapping();
    }
  }, [open, categoryId, attributeId]);

  const fetchAttributeMapping = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch attribute mapping details from the API
      const response = await axios.get(`/api/admin/categories/${categoryId}/attributes/${attributeId}/details`);
      const data = response.data.data || response.data;
      
      setMapping({
        ...data,
        options: data.options || {},
        translations: data.translations || {},
        unit_translations: data.unit_translations || {}
      });
      
      if (data.attribute) {
        setOriginalAttribute(data.attribute);
        
        // Initialize options for select/multiselect
        let attrOptions = [];
        if (data.attribute.attribute_type === 'select' || data.attribute.attribute_type === 'multiselect') {
          let options = data.attribute.options;
          if (typeof options === 'string') {
            try {
              options = JSON.parse(options);
            } catch (e) {
              options = { values: [] };
            }
          }
          
          attrOptions = options?.values || [];
          setAttributeOptions(attrOptions);
          
          // Initialize category-specific options if available
          if (data.options && data.options.values) {
            setOptionValues(data.options.values);
          } else {
            setOptionValues([...attrOptions]);
          }
        }
      }
    } catch (err) {
      console.error('Error fetching attribute mapping details:', err);
      const errorMessage = t('admin.categoryAttributes.fetchDetailsError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleRequiredChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        is_required: e.target.checked
      });
    }
  };

  const handleEnabledChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        is_enabled: e.target.checked
      });
    }
  };

  const handleSortOrderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        sort_order: parseInt(e.target.value) || 0
      });
    }
  };

  const handleHintChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        hint: e.target.value
      });
    }
  };

  const handleDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        description: e.target.value
      });
    }
  };

  const handleCustomComponentChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        custom_component: e.target.value
      });
    }
  };

  const handleUnitChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (mapping) {
      setMapping({
        ...mapping,
        unit: e.target.value
      });
    }
  };

  const handleTranslationChange = (field: string, lang: string, value: string) => {
    if (mapping) {
      setMapping({
        ...mapping,
        translations: {
          ...mapping.translations,
          [lang]: {
            ...(mapping.translations?.[lang] || {}),
            [field]: value
          }
        }
      });
    }
  };

  const handleUnitTranslationChange = (lang: string, value: string) => {
    if (mapping) {
      setMapping({
        ...mapping,
        unit_translations: {
          ...mapping.unit_translations,
          [lang]: value
        }
      });
    }
  };

  // Handle option translations
  const handleOptionTranslationChange = (lang: string, optionValue: string, value: string) => {
    if (mapping) {
      setMapping({
        ...mapping,
        options: {
          ...mapping.options,
          translations: {
            ...(mapping.options?.translations || {}),
            [lang]: {
              ...(mapping.options?.translations?.[lang] || {}),
              [optionValue]: value
            }
          }
        }
      });
    }
  };

  // Manage options for select/multiselect
  const handleAddOption = () => {
    const newOptions = [...optionValues, ''];
    setOptionValues(newOptions);
    
    if (mapping) {
      setMapping({
        ...mapping,
        options: {
          ...mapping.options,
          values: newOptions
        }
      });
    }
  };

  const handleOptionChange = (index: number, value: string) => {
    const oldValue = optionValues[index];
    const newOptions = [...optionValues];
    newOptions[index] = value;
    setOptionValues(newOptions);
    
    if (mapping) {
      // Update options values
      const updatedMapping = {
        ...mapping,
        options: {
          ...mapping.options,
          values: newOptions
        }
      };
      
      // If an existing option value changed, update translations
      if (oldValue && oldValue !== value && mapping.options?.translations) {
        const updatedTranslations = { ...mapping.options.translations };
        
        // Update translations for each language
        Object.keys(updatedTranslations).forEach(lang => {
          if (updatedTranslations[lang] && updatedTranslations[lang][oldValue]) {
            // Create new entry with the new key, keeping old translation value
            updatedTranslations[lang][value] = updatedTranslations[lang][oldValue];
            // Remove old entry
            delete updatedTranslations[lang][oldValue];
          }
        });
        
        updatedMapping.options.translations = updatedTranslations;
      }
      
      setMapping(updatedMapping);
    }
  };

  const handleRemoveOption = (index: number) => {
    const removedValue = optionValues[index];
    const newOptions = optionValues.filter((_, i) => i !== index);
    setOptionValues(newOptions);
    
    if (mapping) {
      // Update options values
      const updatedMapping = {
        ...mapping,
        options: {
          ...mapping.options,
          values: newOptions
        }
      };
      
      // Remove translations for the removed option
      if (removedValue && mapping.options?.translations) {
        const updatedTranslations = { ...mapping.options.translations };
        
        Object.keys(updatedTranslations).forEach(lang => {
          if (updatedTranslations[lang] && updatedTranslations[lang][removedValue]) {
            delete updatedTranslations[lang][removedValue];
          }
        });
        
        updatedMapping.options.translations = updatedTranslations;
      }
      
      setMapping(updatedMapping);
    }
  };

  const handleOptionDragEnd = (result: any) => {
    if (!result || !result.source || !result.destination) return;

    const items = Array.from(optionValues);
    const sourceIndex = result.source.index;
    const destIndex = result.destination.index;

    if (sourceIndex < 0 || destIndex < 0 || sourceIndex >= items.length) {
      return;
    }

    const [reorderedItem] = items.splice(sourceIndex, 1);
    items.splice(destIndex, 0, reorderedItem);

    setOptionValues(items);
    
    if (mapping) {
      setMapping({
        ...mapping,
        options: {
          ...mapping.options,
          values: items
        }
      });
    }
  };

  const handleSave = async () => {
    if (!mapping) return;
    
    try {
      setSaving(true);
      setError(null);
      
      // Prepare data for API
      const updateData = {
        is_required: mapping.is_required,
        is_enabled: mapping.is_enabled,
        sort_order: mapping.sort_order,
        hint: mapping.hint,
        description: mapping.description,
        unit: mapping.unit,
        custom_component: mapping.custom_component,
        translations: mapping.translations,
        unit_translations: mapping.unit_translations,
        options: mapping.options
      };
      
      // Update attribute mapping
      await axios.put(`/api/admin/categories/${categoryId}/attributes/${attributeId}`, updateData);
      
      // Notify parent component
      onUpdate();
      onClose();
    } catch (err) {
      console.error('Error updating attribute in category:', err);
      const errorMessage = t('admin.categoryAttributes.updateAttributeError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setSaving(false);
    }
  };

  const getAttributeTypeName = (type: string) => {
    const attributeTypes: Record<string, string> = {
      'text': t('admin.attributeTypes.text'),
      'number': t('admin.attributeTypes.number'),
      'select': t('admin.attributeTypes.select'),
      'multiselect': t('admin.attributeTypes.multiselect'),
      'boolean': t('admin.attributeTypes.boolean'),
      'range': t('admin.attributeTypes.range'),
      'date': t('admin.attributeTypes.date')
    };
    
    return attributeTypes[type] || type;
  };

  if (!open) return null;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {t('admin.categoryAttributes.editAttributeInCategory')}
      </DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {originalAttribute && (
              <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                <Typography variant="subtitle1" fontWeight="bold">
                  {originalAttribute.display_name}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {originalAttribute.name} ({getAttributeTypeName(originalAttribute.attribute_type)})
                </Typography>
              </Paper>
            )}
            
            <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
              <Tabs value={tabValue} onChange={handleTabChange}>
                <Tab label={t('admin.categoryAttributes.basicTab')} id="attribute-tab-0" />
                <Tab label={t('admin.categoryAttributes.descriptionsTab')} id="attribute-tab-1" />
                {(originalAttribute?.attribute_type === 'select' || originalAttribute?.attribute_type === 'multiselect') && (
                  <Tab label={t('admin.categoryAttributes.optionsTab')} id="attribute-tab-2" />
                )}
                {(originalAttribute?.attribute_type === 'number' || originalAttribute?.attribute_type === 'range') && (
                  <Tab label={t('admin.categoryAttributes.unitsTab')} id="attribute-tab-3" />
                )}
              </Tabs>
            </Box>
            
            {/* Basic Settings Tab */}
            <TabPanel value={tabValue} index={0}>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={mapping?.is_required || false}
                        onChange={handleRequiredChange}
                        color="primary"
                      />
                    }
                    label={t('admin.attributes.required')}
                  />
                </Grid>
                <Grid item xs={12} md={6}>
                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={mapping?.is_enabled || false}
                        onChange={handleEnabledChange}
                        color="primary"
                      />
                    }
                    label={t('admin.attributes.enabled')}
                  />
                </Grid>
                <Grid item xs={12} md={6}>
                  <TextField
                    fullWidth
                    label={t('admin.attributes.sortOrder')}
                    type="number"
                    value={mapping?.sort_order || 0}
                    onChange={handleSortOrderChange}
                    InputProps={{ inputProps: { min: 0 } }}
                    margin="normal"
                  />
                </Grid>
                <Grid item xs={12} md={6}>
                  <TextField
                    fullWidth
                    label={t('admin.attributes.customComponent')}
                    value={mapping?.custom_component || ''}
                    onChange={handleCustomComponentChange}
                    margin="normal"
                    helperText={t('admin.attributes.customComponentHelp')}
                  />
                </Grid>
              </Grid>
            </TabPanel>
            
            {/* Descriptions Tab */}
            <TabPanel value={tabValue} index={1}>
              <Typography variant="subtitle2" gutterBottom>
                {t('admin.categoryAttributes.defaultLanguage')} ({i18n.language.toUpperCase()})
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label={t('admin.categoryAttributes.hint')}
                    value={mapping?.hint || ''}
                    onChange={handleHintChange}
                    margin="normal"
                    helperText={t('admin.categoryAttributes.hintHelp')}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label={t('admin.categoryAttributes.description')}
                    value={mapping?.description || ''}
                    onChange={handleDescriptionChange}
                    margin="normal"
                    multiline
                    rows={3}
                    helperText={t('admin.categoryAttributes.descriptionHelp')}
                  />
                </Grid>
              </Grid>
              
              <Divider sx={{ my: 3 }} />
              
              <Accordion>
                <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                  <Typography variant="subtitle2">
                    {t('admin.categoryAttributes.translations')}
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
                          <Grid item xs={12}>
                            <TextField
                              fullWidth
                              label={t('admin.categoryAttributes.hint')}
                              value={mapping?.translations?.[lang]?.hint || ''}
                              onChange={(e) => handleTranslationChange('hint', lang, e.target.value)}
                              margin="normal"
                            />
                          </Grid>
                          <Grid item xs={12}>
                            <TextField
                              fullWidth
                              label={t('admin.categoryAttributes.description')}
                              value={mapping?.translations?.[lang]?.description || ''}
                              onChange={(e) => handleTranslationChange('description', lang, e.target.value)}
                              margin="normal"
                              multiline
                              rows={3}
                            />
                          </Grid>
                        </Grid>
                      </Box>
                    ))}
                </AccordionDetails>
              </Accordion>
            </TabPanel>
            
            {/* Options Tab for Select/Multiselect */}
            {(originalAttribute?.attribute_type === 'select' || originalAttribute?.attribute_type === 'multiselect') && (
              <TabPanel value={tabValue} index={2}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="subtitle1">
                    {t('admin.categoryAttributes.customizeOptions')}
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
                  {optionValues.length > 0 ? (
                    <List>
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
                                <Tooltip title={t('admin.remove')}>
                                  <IconButton
                                    size="small"
                                    onClick={() => handleRemoveOption(index)}
                                    color="error"
                                  >
                                    <DeleteIcon />
                                  </IconButton>
                                </Tooltip>
                              </Grid>
                            </Grid>
                          </Box>
                        );
                      })}
                    </List>
                  ) : (
                    <Typography color="textSecondary" sx={{ mt: 2, mb: 2, textAlign: 'center' }}>
                      {t('admin.categoryAttributes.noCustomOptions')}
                    </Typography>
                  )}
                </Box>
                
                {attributeOptions.length > 0 && (
                  <Box sx={{ mt: 3 }}>
                    <Typography variant="subtitle2" sx={{ mb: 1 }}>
                      {t('admin.categoryAttributes.globalOptionsAvailable')}:
                    </Typography>
                    <Paper variant="outlined" sx={{ p: 1 }}>
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                        {attributeOptions.map((option) => (
                          <Tooltip 
                            key={option} 
                            title={t('admin.categoryAttributes.clickToAdd')}
                          >
                            <Button 
                              size="small" 
                              variant="outlined"
                              onClick={() => {
                                if (!optionValues.includes(option)) {
                                  const newOptions = [...optionValues, option];
                                  setOptionValues(newOptions);
                                  
                                  if (mapping) {
                                    setMapping({
                                      ...mapping,
                                      options: {
                                        ...mapping.options,
                                        values: newOptions
                                      }
                                    });
                                  }
                                }
                              }}
                              disabled={optionValues.includes(option)}
                            >
                              {option}
                            </Button>
                          </Tooltip>
                        ))}
                      </Box>
                    </Paper>
                  </Box>
                )}
                
                {/* Options Translations */}
                {optionValues.length > 0 && (
                  <Accordion sx={{ mt: 3 }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                      <Typography variant="subtitle2">
                        {t('admin.categoryAttributes.optionTranslations')}
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
                                  <Grid item xs={12} md={6} key={`${lang}-${index}`}>
                                    <TextField
                                      fullWidth
                                      label={option}
                                      value={mapping?.options?.translations?.[lang]?.[option] || ''}
                                      onChange={(e) => handleOptionTranslationChange(lang, option, e.target.value)}
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
            )}
            
            {/* Units Tab for Number/Range */}
            {(originalAttribute?.attribute_type === 'number' || originalAttribute?.attribute_type === 'range') && (
              <TabPanel value={tabValue} index={3}>
                <Typography variant="subtitle2" gutterBottom>
                  {t('admin.categoryAttributes.defaultLanguage')} ({i18n.language.toUpperCase()})
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <TextField
                      fullWidth
                      label={t('admin.categoryAttributes.unit')}
                      value={mapping?.unit || ''}
                      onChange={handleUnitChange}
                      margin="normal"
                      helperText={t('admin.categoryAttributes.unitHelp')}
                    />
                  </Grid>
                </Grid>
                
                <Divider sx={{ my: 3 }} />
                
                <Accordion>
                  <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography variant="subtitle2">
                      {t('admin.categoryAttributes.unitTranslations')}
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
                            <Grid item xs={12}>
                              <TextField
                                fullWidth
                                label={t('admin.categoryAttributes.unit')}
                                value={mapping?.unit_translations?.[lang] || ''}
                                onChange={(e) => handleUnitTranslationChange(lang, e.target.value)}
                                margin="normal"
                              />
                            </Grid>
                          </Grid>
                        </Box>
                      ))}
                  </AccordionDetails>
                </Accordion>
              </TabPanel>
            )}
          </>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          {t('common.cancel')}
        </Button>
        <Button 
          onClick={handleSave} 
          color="primary" 
          variant="contained"
          disabled={loading || saving}
          startIcon={saving ? <CircularProgress size={20} /> : null}
        >
          {t('common.save')}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default EditAttributeInCategory;