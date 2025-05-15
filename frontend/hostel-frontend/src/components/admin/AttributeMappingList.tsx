import React, { useState, useEffect } from 'react';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Switch,
  TextField,
  Typography,
  Tooltip,
  CircularProgress,
  Alert
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';
import SettingsIcon from '@mui/icons-material/Settings';
import EditAttributeInCategory from './EditAttributeInCategory';
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
}

interface CategoryAttributeMapping {
  category_id: number;
  attribute_id: number;
  is_required: boolean;
  is_enabled: boolean;
  sort_order: number;
  custom_component?: string;
  attribute?: Attribute; // Extended info about the attribute
}

interface EditState {
  attributeId: number | null;
  is_required: boolean;
  is_enabled: boolean;
  sort_order: number;
}

interface AdvancedEditState {
  open: boolean;
  attributeId: number | null;
}

interface AttributeMappingListProps {
  categoryId: number;
  onError?: (message: string) => void;
  mappings?: CategoryAttributeMapping[];
  onUpdateAttribute?: (attributeId: number, isRequired: boolean, isEnabled: boolean, sortOrder: number) => void;
  onRemoveAttribute?: (attributeId: number) => void;
}

const AttributeMappingList: React.FC<AttributeMappingListProps> = ({ 
  categoryId, 
  onError, 
  mappings: propMappings, 
  onUpdateAttribute, 
  onRemoveAttribute 
}) => {
  const { t } = useTranslation();
  const [mappings, setMappings] = useState<CategoryAttributeMapping[]>([]);
  const [loading, setLoading] = useState<boolean>(propMappings ? false : true);
  const [error, setError] = useState<string | null>(null);
  const [editState, setEditState] = useState<EditState>({
    attributeId: null,
    is_required: false,
    is_enabled: false,
    sort_order: 0
  });
  const [customComponentValue, setCustomComponentValue] = useState<string>('');
  const [advancedEditState, setAdvancedEditState] = useState<AdvancedEditState>({
    open: false,
    attributeId: null
  });

  // Load mappings when categoryId changes or propMappings updates
  useEffect(() => {
    if (propMappings) {
      // Если переданы маппинги извне, используем их
      const sortedMappings = [...propMappings].sort((a, b) => a.sort_order - b.sort_order);
      setMappings(sortedMappings);
      setLoading(false);
    } else if (categoryId) {
      // Иначе загружаем с сервера
      loadMappings();
    }
  }, [categoryId, propMappings]);

  const loadMappings = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await axios.get(`/api/admin/categories/${categoryId}/attributes`);
      const data = response.data.data || response.data;
      
      // Sort mappings by sort_order
      const sortedMappings = [...data].sort((a, b) => a.sort_order - b.sort_order);
      setMappings(sortedMappings);
    } catch (err) {
      console.error('Error loading category attribute mappings:', err);
      const errorMessage = t('admin.categoryAttributes.fetchMappingsError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleEditClick = (mapping: CategoryAttributeMapping) => {
    setEditState({
      attributeId: mapping.attribute_id,
      is_required: mapping.is_required,
      is_enabled: mapping.is_enabled,
      sort_order: mapping.sort_order
    });
    setCustomComponentValue(mapping.custom_component || '');
  };

  const handleCancelEdit = () => {
    setEditState({
      attributeId: null,
      is_required: false,
      is_enabled: false,
      sort_order: 0
    });
    setCustomComponentValue('');
  };

  const handleOpenAdvancedEdit = (attributeId: number) => {
    setAdvancedEditState({
      open: true,
      attributeId
    });
  };

  const handleCloseAdvancedEdit = () => {
    setAdvancedEditState({
      open: false,
      attributeId: null
    });
  };

  const handleAttributeUpdated = () => {
    loadMappings();
  };

  const handleSaveEdit = async () => {
    if (!editState.attributeId) return;
    
    try {
      setLoading(true);
      
      if (onUpdateAttribute) {
        // Используем переданную функцию обновления
        await onUpdateAttribute(
          editState.attributeId, 
          editState.is_required, 
          editState.is_enabled, 
          editState.sort_order
        );
      } else {
        // Или делаем запрос напрямую
        await axios.put(`/api/admin/categories/${categoryId}/attributes/${editState.attributeId}`, {
          is_required: editState.is_required,
          is_enabled: editState.is_enabled,
          sort_order: editState.sort_order,
          custom_component: customComponentValue
        });
        
        // Если нет переданной функции, то обновим данные
        await loadMappings();
      }
      
      // Reset edit state
      handleCancelEdit();
    } catch (err) {
      console.error('Error updating attribute mapping:', err);
      const errorMessage = t('admin.categoryAttributes.updateAttributeError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveAttribute = async (attributeId: number) => {
    if (!window.confirm(t('admin.categoryAttributes.confirmRemove'))) {
      return;
    }
    
    try {
      setLoading(true);
      
      if (onRemoveAttribute) {
        // Используем переданную функцию удаления
        await onRemoveAttribute(attributeId);
      } else {
        // Или делаем запрос напрямую
        await axios.delete(`/api/admin/categories/${categoryId}/attributes/${attributeId}`);
        
        // Если нет переданной функции, то обновим данные
        await loadMappings();
      }
    } catch (err) {
      console.error('Error removing attribute from category:', err);
      const errorMessage = t('admin.categoryAttributes.removeAttributeError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleRequiredChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditState({
      ...editState,
      is_required: e.target.checked
    });
  };

  const handleEnabledChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditState({
      ...editState,
      is_enabled: e.target.checked
    });
  };

  const handleSortOrderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditState({
      ...editState,
      sort_order: parseInt(e.target.value) || 0
    });
  };

  if (loading && mappings.length === 0) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error && mappings.length === 0) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (mappings.length === 0) {
    return (
      <Typography variant="body1" sx={{ p: 2, textAlign: 'center' }}>
        {t('admin.categoryAttributes.noAttributesMapped')}
      </Typography>
    );
  }

  return (
    <Box sx={{ width: '100%', mb: 2 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      <TableContainer component={Paper} variant="outlined" sx={{ overflowX: 'auto' }}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>{t('admin.attributes.name')}</TableCell>
              <TableCell align="center">{t('admin.attributes.type')}</TableCell>
              <TableCell align="center">{t('admin.attributes.required')}</TableCell>
              <TableCell align="center">{t('admin.attributes.enabled')}</TableCell>
              <TableCell align="center">{t('admin.attributes.sortOrder')}</TableCell>
              <TableCell align="center">{t('admin.attributes.customComponent')}</TableCell>
              <TableCell align="right">{t('admin.actions')}</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {mappings.map((mapping) => (
              <TableRow key={mapping.attribute_id} hover>
                <TableCell>
                  {mapping.attribute ? mapping.attribute.display_name : `Attribute #${mapping.attribute_id}`}
                </TableCell>
                <TableCell align="center">
                  {mapping.attribute ? mapping.attribute.attribute_type : ''}
                </TableCell>
                <TableCell align="center">
                  {editState.attributeId === mapping.attribute_id ? (
                    <Switch
                      checked={editState.is_required}
                      onChange={handleRequiredChange}
                      color="primary"
                      size="small"
                    />
                  ) : (
                    <Switch 
                      checked={mapping.is_required} 
                      disabled 
                      size="small" 
                    />
                  )}
                </TableCell>
                <TableCell align="center">
                  {editState.attributeId === mapping.attribute_id ? (
                    <Switch
                      checked={editState.is_enabled}
                      onChange={handleEnabledChange}
                      color="primary"
                      size="small"
                    />
                  ) : (
                    <Switch 
                      checked={mapping.is_enabled} 
                      disabled 
                      size="small" 
                    />
                  )}
                </TableCell>
                <TableCell align="center">
                  {editState.attributeId === mapping.attribute_id ? (
                    <TextField
                      type="number"
                      value={editState.sort_order}
                      onChange={handleSortOrderChange}
                      variant="outlined"
                      size="small"
                      sx={{ width: '70px' }}
                      inputProps={{ min: 0, max: 1000 }}
                    />
                  ) : (
                    mapping.sort_order
                  )}
                </TableCell>
                <TableCell align="center">
                  {editState.attributeId === mapping.attribute_id ? (
                    <TextField
                      value={customComponentValue || ''}
                      onChange={(e) => setCustomComponentValue(e.target.value)}
                      variant="outlined"
                      size="small"
                      sx={{ width: '120px' }}
                      placeholder={t('admin.attributes.customComponentPlaceholder')}
                    />
                  ) : (
                    mapping.custom_component || '-'
                  )}
                </TableCell>
                <TableCell align="right">
                  {editState.attributeId === mapping.attribute_id ? (
                    <>
                      <IconButton 
                        onClick={handleSaveEdit} 
                        size="small" 
                        color="primary"
                        disabled={loading}
                      >
                        <SaveIcon fontSize="small" />
                      </IconButton>
                      <IconButton 
                        onClick={handleCancelEdit} 
                        size="small" 
                        color="default"
                        disabled={loading}
                      >
                        <CancelIcon fontSize="small" />
                      </IconButton>
                    </>
                  ) : (
                    <>
                      <Tooltip title={t('admin.edit')}>
                        <span>
                          <IconButton 
                            onClick={() => handleEditClick(mapping)} 
                            size="small" 
                            color="primary"
                            disabled={loading || editState.attributeId !== null}
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </span>
                      </Tooltip>
                      <Tooltip title={t('admin.categoryAttributes.advancedEdit')}>
                        <span>
                          <IconButton 
                            onClick={() => handleOpenAdvancedEdit(mapping.attribute_id)} 
                            size="small" 
                            color="secondary"
                            disabled={loading || editState.attributeId !== null}
                          >
                            <SettingsIcon fontSize="small" />
                          </IconButton>
                        </span>
                      </Tooltip>
                      <Tooltip title={t('admin.remove')}>
                        <span>
                          <IconButton 
                            onClick={() => handleRemoveAttribute(mapping.attribute_id)} 
                            size="small" 
                            color="error"
                            disabled={loading || editState.attributeId !== null}
                          >
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </span>
                      </Tooltip>
                    </>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {loading && editState.attributeId === null && (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 1 }}>
          <CircularProgress size={24} />
        </Box>
      )}

      {/* Advanced Edit Dialog */}
      <EditAttributeInCategory
        open={advancedEditState.open}
        categoryId={categoryId}
        attributeId={advancedEditState.attributeId || 0}
        onClose={handleCloseAdvancedEdit}
        onUpdate={handleAttributeUpdated}
        onError={onError}
      />
    </Box>
  );
};

export default AttributeMappingList;