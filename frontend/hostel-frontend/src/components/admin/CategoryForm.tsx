import React, { useState, useEffect } from 'react';
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
  Typography
} from '@mui/material';
import { useTranslation } from 'react-i18next';

interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number | null;
  icon?: string;
  translations?: Record<string, string>;
  has_custom_ui: boolean;
  custom_ui_component?: string;
}

interface CategoryFormProps {
  open: boolean;
  category: Category | null;
  categories: Category[];
  onSubmit: (formData: Partial<Category>) => void;
  onClose: () => void;
}

// Список доступных компонентов UI для категорий
const availableComponents = [
  { value: "AutoCategoryUI", label: "Автомобили" },
  { value: "RealEstateCategoryUI", label: "Недвижимость" },
  { value: "ElectronicsCategoryUI", label: "Электроника" },
];

const CategoryForm: React.FC<CategoryFormProps> = ({
  open,
  category,
  categories,
  onSubmit,
  onClose
}) => {
  const { t, i18n } = useTranslation();
  const [formData, setFormData] = useState<Partial<Category>>({
    name: '',
    slug: '',
    parent_id: null,
    icon: '',
    has_custom_ui: false,
    custom_ui_component: '',
    translations: {}
  });
  const [languages] = useState<string[]>(['en', 'ru', 'sr']);

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name,
        slug: category.slug,
        parent_id: category.parent_id || null,
        icon: category.icon || '',
        has_custom_ui: category.has_custom_ui || false,
        custom_ui_component: category.custom_ui_component || '',
        translations: category.translations || {}
      });
    } else {
      setFormData({
        name: '',
        slug: '',
        parent_id: null,
        icon: '',
        has_custom_ui: false,
        custom_ui_component: '',
        translations: {}
      });
    }
  }, [category]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    // Автоматически генерируем slug из названия на латинице
    if (name === 'name' && !formData.slug) {
      const slug = value
        .toLowerCase()
        .replace(/[^a-z0-9\s]/g, '') // Оставляем только латиницу и цифры
        .replace(/\s+/g, '-'); // Заменяем пробелы на дефисы
      
      setFormData(prev => ({
        ...prev,
        slug
      }));
    }
  };

  const handleSelectChange = (e: any) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value === '' ? null : value
    }));
  };

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: checked
    }));

    // Если убираем флаг кастомного UI, то очищаем и компонент
    if (name === 'has_custom_ui' && !checked) {
      setFormData(prev => ({
        ...prev,
        custom_ui_component: ''
      }));
    }
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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>
          {category ? t('admin.categories.edit') : t('admin.categories.create')}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                name="name"
                label={t('admin.categories.name')}
                value={formData.name}
                onChange={handleChange}
                required
                margin="normal"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                name="slug"
                label={t('admin.categories.slug')}
                value={formData.slug}
                onChange={handleChange}
                required
                margin="normal"
                helperText={t('admin.categories.slugHelp')}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth margin="normal">
                <InputLabel id="parent-category-label">
                  {t('admin.categories.parentCategory')}
                </InputLabel>
                <Select
                  labelId="parent-category-label"
                  name="parent_id"
                  value={formData.parent_id || ''}
                  onChange={handleSelectChange}
                  label={t('admin.categories.parentCategory')}
                >
                  <MenuItem value="">
                    <em>{t('admin.categories.noParent')}</em>
                  </MenuItem>
                  {categories
                    .filter(c => !category || c.id !== category.id)
                    .map(c => (
                      <MenuItem key={c.id} value={c.id}>
                        {c.name}
                      </MenuItem>
                    ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                name="icon"
                label={t('admin.categories.icon')}
                value={formData.icon || ''}
                onChange={handleChange}
                margin="normal"
                helperText={t('admin.categories.iconHelp')}
              />
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle1" sx={{ mt: 2, mb: 1 }}>
                {t('admin.categories.translations')}
              </Typography>
              
              {languages
                .filter(lang => lang !== i18n.language)
                .map(lang => (
                  <TextField
                    key={lang}
                    fullWidth
                    name={`translation_${lang}`}
                    label={t(`languages.${lang}`)}
                    value={formData.translations?.[lang] || ''}
                    onChange={(e) => handleTranslationChange(lang, e.target.value)}
                    margin="normal"
                  />
                ))}
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle1" sx={{ mt: 2, mb: 1 }}>
                {t('admin.categories.customUi')}
              </Typography>
              
              <FormControlLabel
                control={
                  <Checkbox
                    checked={formData.has_custom_ui || false}
                    onChange={handleCheckboxChange}
                    name="has_custom_ui"
                  />
                }
                label={t('admin.categories.hasCustomUi')}
              />
            </Grid>

            {formData.has_custom_ui && (
              <Grid item xs={12}>
                <FormControl fullWidth margin="normal">
                  <InputLabel id="custom-component-label">
                    {t('admin.categories.customComponent')}
                  </InputLabel>
                  <Select
                    labelId="custom-component-label"
                    name="custom_ui_component"
                    value={formData.custom_ui_component || ''}
                    onChange={handleSelectChange}
                    label={t('admin.categories.customComponent')}
                    required={formData.has_custom_ui}
                  >
                    <MenuItem value="">
                      <em>{t('admin.categories.selectComponent')}</em>
                    </MenuItem>
                    {availableComponents.map(component => (
                      <MenuItem key={component.value} value={component.value}>
                        {component.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            )}
          </Grid>
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

export default CategoryForm;