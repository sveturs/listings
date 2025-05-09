// frontend/hostel-frontend/src/components/store/ImportSourceForm.tsx
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText,
  Alert,
  Paper,
  Stack,
  Divider,
  SelectChangeEvent
} from '@mui/material';
import { Plus, Save, Clock } from 'lucide-react';
import axios from '../../api/axios';

export interface ImportSource {
  id?: number | string;
  type: 'csv' | 'xml';
  url: string;
  schedule: string;
  storefront_id: number | string;
  [key: string]: any;
}

interface ImportSourceFormProps {
  onClose: () => void;
  onSuccess: (data: ImportSource) => void;
  storefrontId: number | string;
  initialData?: ImportSource | null;
}

const ImportSourceForm: React.FC<ImportSourceFormProps> = ({ 
  onClose, 
  onSuccess, 
  storefrontId, 
  initialData = null 
}) => {
  const { t } = useTranslation(['common', 'marketplace']);
  const isEditing = !!initialData;
  
  const [formData, setFormData] = useState<ImportSource>({
    type: initialData?.type || 'csv',
    url: initialData?.url || '',
    schedule: initialData?.schedule || '',
    storefront_id: storefrontId
  });
  
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | SelectChangeEvent): void => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (!formData.type) {
        setError(t('marketplace:store.import.selectTypeError'));
        setLoading(false);
        return;
      }

      let response;
      if (isEditing && initialData?.id) {
        response = await axios.put(`/api/v1/storefronts/import-sources/${initialData.id}`, formData);
      } else {
        response = await axios.post('/api/v1/storefronts/import-sources', formData);
      }
      
      if (onSuccess) {
        onSuccess(response.data.data);
      }
      
      onClose();
    } catch (err: any) {
      console.error('Error submitting import source:', err);
      setError(err.response?.data?.error || t('marketplace:store.import.submitError'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom>
        {isEditing
          ? t('marketplace:store.import.editSource')
          : t('marketplace:store.import.createSource')}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Stack spacing={3}>
          <FormControl fullWidth>
            {/* @ts-ignore - Ignoring TypeScript error for InputLabel due to MUI API */}
            <InputLabel id="import-type-label">
              {t('marketplace:store.import.typeLabel')}
            </InputLabel>
            <Select
              labelId="import-type-label"
              id="import-type"
              name="type"
              value={formData.type}
              onChange={handleChange}
              label={t('marketplace:store.import.typeLabel')}
              disabled={loading}
            >
              <MenuItem value="csv">{t('marketplace:store.import.typeCSV')}</MenuItem>
              <MenuItem value="xml">{t('marketplace:store.import.typeXML')}</MenuItem>
            </Select>
            <FormHelperText>
              {t('marketplace:store.import.typeHelp')}
            </FormHelperText>
          </FormControl>

          <TextField
            name="url"
            label={t('marketplace:store.import.urlLabel')}
            value={formData.url}
            onChange={handleChange}
            fullWidth
            placeholder={t('marketplace:store.import.urlPlaceholder')}
            helperText={
              formData.type === 'xml'
                ? t('marketplace:store.import.xmlUrlHelp', {
                    defaultValue: 'Укажите URL к ZIP-архиву с XML файлом. URL должен заканчиваться на .zip'
                  })
                : t('marketplace:store.import.urlHelp')
            }
            disabled={loading}
          />

          <Divider />
          
          <Box>
            <Typography variant="subtitle2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Clock size={16} />
              {t('marketplace:store.import.scheduleLabel', { defaultValue: 'Расписание обновления' })}
            </Typography>
            
            <FormControl fullWidth>
              {/* @ts-ignore - Ignoring TypeScript error for InputLabel due to MUI API */}
              <InputLabel id="schedule-label">
                {t('marketplace:store.import.scheduleFrequency', { defaultValue: 'Частота обновления' })}
              </InputLabel>
              <Select
                labelId="schedule-label"
                id="schedule"
                name="schedule"
                value={formData.schedule}
                onChange={handleChange}
                label={t('marketplace:store.import.scheduleFrequency')}
                disabled={loading || !formData.url}
              >
                <MenuItem value="">
                  {t('marketplace:store.import.scheduleManualy', { defaultValue: 'Только вручную' })}
                </MenuItem>
                <MenuItem value="hourly">
                  {t('marketplace:store.import.scheduleHourly', { defaultValue: 'Каждый час' })}
                </MenuItem>
                <MenuItem value="daily">
                  {t('marketplace:store.import.scheduleDaily', { defaultValue: 'Ежедневно' })}
                </MenuItem>
                <MenuItem value="weekly">
                  {t('marketplace:store.import.scheduleWeekly', { defaultValue: 'Еженедельно' })}
                </MenuItem>
                <MenuItem value="monthly">
                  {t('marketplace:store.import.scheduleMonthly', { defaultValue: 'Ежемесячно' })}
                </MenuItem>
              </Select>
              <FormHelperText>
                {formData.url
                  ? t('marketplace:store.import.scheduleHelp', {
                      defaultValue: 'Выберите, как часто автоматически обновлять данные из указанного URL'
                    })
                  : t('marketplace:store.import.scheduleUrlRequired', {
                      defaultValue: 'Для настройки расписания необходимо указать URL'
                    })}
              </FormHelperText>
            </FormControl>
          </Box>
          
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2, mt: 2 }}>
            <Button
              variant="outlined"
              onClick={onClose}
              disabled={loading}
            >
              {t('common:buttons.cancel')}
            </Button>
            <Button
              type="submit"
              variant="contained"
              startIcon={isEditing ? <Save /> : <Plus />}
              disabled={loading}
            >
              {loading
                ? t('common:common.processing')
                : isEditing
                ? t('common:buttons.save')
                : t('marketplace:store.import.createSource')}
            </Button>
          </Box>
        </Stack>
      </form>
    </Paper>
  );
};

export default ImportSourceForm;