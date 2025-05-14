import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  Button, 
  Divider,
  CircularProgress,
  TextField,
  InputAdornment,
  IconButton,
  Snackbar,
  Alert
} from '@mui/material';
import { 
  Add as AddIcon, 
  Search as SearchIcon,
  Clear as ClearIcon
} from '@mui/icons-material';
import axios from '../../api/axios';
import { useTranslation } from 'react-i18next';
import AttributeForm from '../../components/admin/AttributeForm';
import AttributesTable from '../../components/admin/AttributesTable';

export interface Attribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  options?: any;
  validation_rules?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_required: boolean;
  sort_order: number;
  created_at: string;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
  custom_component?: string;
}

const AttributeManagementPage: React.FC = () => {
  const { t } = useTranslation();
  const [attributes, setAttributes] = useState<Attribute[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [openForm, setOpenForm] = useState<boolean>(false);
  const [selectedAttribute, setSelectedAttribute] = useState<Attribute | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [searchText, setSearchText] = useState<string>('');
  const [filteredAttributes, setFilteredAttributes] = useState<Attribute[]>([]);

  useEffect(() => {
    fetchAttributes();
  }, []);

  useEffect(() => {
    // Фильтрация атрибутов по поисковому запросу
    if (searchText.trim() === '') {
      setFilteredAttributes(attributes);
    } else {
      const searchLower = searchText.toLowerCase();
      const filtered = attributes.filter(attr => 
        attr.name.toLowerCase().includes(searchLower) || 
        attr.display_name.toLowerCase().includes(searchLower) ||
        attr.attribute_type.toLowerCase().includes(searchLower)
      );
      setFilteredAttributes(filtered);
    }
  }, [searchText, attributes]);

  const fetchAttributes = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/admin/attributes');
      setAttributes(response.data.data);
      setFilteredAttributes(response.data.data);
      setError(null);
    } catch (err) {
      console.error('Error fetching attributes:', err);
      setError(t('admin.attributes.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleCreateAttribute = () => {
    setSelectedAttribute(null);
    setOpenForm(true);
  };

  const handleEditAttribute = (attribute: Attribute) => {
    setSelectedAttribute(attribute);
    setOpenForm(true);
  };

  const handleDeleteAttribute = async (attributeId: number) => {
    if (window.confirm(t('admin.attributes.confirmDelete'))) {
      try {
        await axios.delete(`/api/admin/attributes/${attributeId}`);
        fetchAttributes();
        setSuccess(t('admin.attributes.deleteSuccess'));
      } catch (err: any) {
        console.error('Error deleting attribute:', err);
        setError(err.response?.data?.message || t('admin.attributes.deleteError'));
      }
    }
  };

  const handleAttributeFormSubmit = async (formData: Partial<Attribute>) => {
    try {
      if (selectedAttribute) {
        // Update existing attribute
        await axios.put(`/api/admin/attributes/${selectedAttribute.id}`, formData);
        setSuccess(t('admin.attributes.updateSuccess'));
      } else {
        // Create new attribute
        await axios.post('/api/admin/attributes', formData);
        setSuccess(t('admin.attributes.createSuccess'));
      }
      setOpenForm(false);
      fetchAttributes();
    } catch (err: any) {
      console.error('Error saving attribute:', err);
      setError(err.response?.data?.message || t('admin.attributes.saveError'));
    }
  };

  const handleFormClose = () => {
    setOpenForm(false);
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value);
  };

  const handleClearSearch = () => {
    setSearchText('');
  };

  return (
    <Box sx={{ p: 3 }}>
      <Paper sx={{ p: 2 }}>
        <Grid container spacing={2} alignItems="center" sx={{ mb: 2 }}>
          <Grid item xs>
            <Typography variant="h5" component="h1">
              {t('admin.attributes.title')}
            </Typography>
          </Grid>
          <Grid item md={4} xs={12}>
            <TextField
              fullWidth
              placeholder={t('admin.attributes.search')}
              value={searchText}
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
                endAdornment: searchText && (
                  <InputAdornment position="end">
                    <IconButton onClick={handleClearSearch} size="small">
                      <ClearIcon />
                    </IconButton>
                  </InputAdornment>
                )
              }}
              size="small"
            />
          </Grid>
          <Grid item>
            <Button 
              variant="contained" 
              color="primary" 
              startIcon={<AddIcon />}
              onClick={handleCreateAttribute}
            >
              {t('admin.attributes.addButton')}
            </Button>
          </Grid>
        </Grid>

        <Divider sx={{ mb: 2 }} />

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <AttributesTable 
            attributes={filteredAttributes} 
            onEdit={handleEditAttribute}
            onDelete={handleDeleteAttribute}
          />
        )}
      </Paper>

      {openForm && (
        <AttributeForm 
          open={openForm}
          attribute={selectedAttribute}
          onSubmit={handleAttributeFormSubmit}
          onClose={handleFormClose}
        />
      )}

      <Snackbar 
        open={!!success} 
        autoHideDuration={6000} 
        onClose={() => setSuccess(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={() => setSuccess(null)} severity="success" sx={{ width: '100%' }}>
          {success}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default AttributeManagementPage;