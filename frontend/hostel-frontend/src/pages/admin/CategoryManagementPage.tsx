import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  Button, 
  IconButton, 
  Divider,
  CircularProgress
} from '@mui/material';
import { 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon
} from '@mui/icons-material';
import axios from '../../api/axios';
import { useTranslation } from 'react-i18next';
import CategoryForm from '../../components/admin/CategoryForm';
import CategoryTreeView from '../../components/admin/CategoryTreeView';

interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number | null;
  icon?: string;
  created_at: string;
  translations?: Record<string, string>;
  listing_count: number;
  has_custom_ui: boolean;
  custom_ui_component?: string;
}

const CategoryManagementPage: React.FC = () => {
  const { t } = useTranslation();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [openForm, setOpenForm] = useState<boolean>(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/admin/categories');
      setCategories(response.data.data);
      setError(null);
    } catch (err) {
      console.error('Error fetching categories:', err);
      setError(t('admin.categories.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCategory = () => {
    setSelectedCategory(null);
    setOpenForm(true);
  };

  const handleEditCategory = (category: Category) => {
    setSelectedCategory(category);
    setOpenForm(true);
  };

  const handleDeleteCategory = async (categoryId: number) => {
    if (window.confirm(t('admin.categories.confirmDelete'))) {
      try {
        await axios.delete(`/api/admin/categories/${categoryId}`);
        fetchCategories();
      } catch (err) {
        console.error('Error deleting category:', err);
        setError(t('admin.categories.deleteError'));
      }
    }
  };

  const handleCategoryFormSubmit = async (formData: Partial<Category>) => {
    try {
      if (selectedCategory) {
        // Update existing category
        await axios.put(`/api/admin/categories/${selectedCategory.id}`, formData);
      } else {
        // Create new category
        await axios.post('/api/admin/categories', formData);
      }
      setOpenForm(false);
      fetchCategories();
    } catch (err) {
      console.error('Error saving category:', err);
      setError(t('admin.categories.saveError'));
    }
  };

  const handleFormClose = () => {
    setOpenForm(false);
  };

  const handleCategoryReorder = async (orderedIds: number[]) => {
    try {
      await axios.post(`/api/admin/categories/${orderedIds[0]}/reorder`, { ordered_ids: orderedIds });
      fetchCategories();
    } catch (err) {
      console.error('Error reordering categories:', err);
      setError(t('admin.categories.reorderError'));
    }
  };

  const handleCategoryMove = async (categoryId: number, newParentId: number) => {
    try {
      await axios.put(`/api/admin/categories/${categoryId}/move`, { new_parent_id: newParentId });
      fetchCategories();
    } catch (err) {
      console.error('Error moving category:', err);
      setError(t('admin.categories.moveError'));
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Paper sx={{ p: 2 }}>
        <Grid container spacing={2} alignItems="center" sx={{ mb: 2 }}>
          <Grid item xs>
            <Typography variant="h5" component="h1">
              {t('admin.categories.title')}
            </Typography>
          </Grid>
          <Grid item>
            <Button 
              variant="contained" 
              color="primary" 
              startIcon={<AddIcon />}
              onClick={handleCreateCategory}
            >
              {t('admin.categories.addButton')}
            </Button>
          </Grid>
        </Grid>

        <Divider sx={{ mb: 2 }} />

        {error && (
          <Typography color="error" sx={{ mb: 2 }}>
            {error}
          </Typography>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <CategoryTreeView 
            categories={categories} 
            onEdit={handleEditCategory}
            onDelete={handleDeleteCategory}
            onReorder={handleCategoryReorder}
            onMove={handleCategoryMove}
          />
        )}
      </Paper>

      {openForm && (
        <CategoryForm 
          open={openForm}
          category={selectedCategory}
          categories={categories}
          onSubmit={handleCategoryFormSubmit}
          onClose={handleFormClose}
        />
      )}
    </Box>
  );
};

export default CategoryManagementPage;