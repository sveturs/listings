import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  Tabs,
  Tab,
  Tooltip,
  Alert,
  CircularProgress
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  Code as CodeIcon,
  Preview as PreviewIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
// import CodeEditor from '@uiw/react-textarea-code-editor';
import axios from '../../api/axios';

// Простой компонент для редактора кода
const CodeEditor: React.FC<{
  value: string;
  language: string;
  placeholder?: string;
  onChange: (event: { target: { value: string } }) => void;
  padding?: number;
  style?: React.CSSProperties;
}> = ({ value, onChange, placeholder, style }) => {
  return (
    <TextField
      multiline
      fullWidth
      value={value}
      onChange={(e) => onChange({ target: { value: e.target.value } })}
      placeholder={placeholder}
      sx={{
        '& .MuiInputBase-input': {
          fontFamily: 'monospace',
          fontSize: '14px',
          backgroundColor: '#f5f5f5',
          ...style,
        }
      }}
      minRows={10}
    />
  );
};

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(3),
  marginBottom: theme.spacing(3),
}));

const ComponentTypeBadge = styled(Chip)<{ componentType: string }>(({ theme, componentType }) => ({
  backgroundColor: 
    componentType === 'category' ? theme.palette.primary.main :
    componentType === 'attribute' ? theme.palette.secondary.main :
    theme.palette.info.main,
  color: '#ffffff',
}));

interface CustomComponent {
  id: number;
  name: string;
  display_name: string;
  description: string;
  component_type: 'category' | 'attribute' | 'filter';
  component_code: string;
  configuration: any;
  dependencies: any;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface ComponentUsage {
  id: number;
  component_id: number;
  category_id?: number;
  usage_context: string;
  placement: string;
  priority: number;
  configuration: any;
  is_active: boolean;
}

interface ComponentTemplate {
  id: number;
  name: string;
  display_name: string;
  description: string;
  template_code: string;
  variables: any;
  is_shared: boolean;
}

const CustomComponentsPage: React.FC = () => {
  const [components, setComponents] = useState<CustomComponent[]>([]);
  const [templates, setTemplates] = useState<ComponentTemplate[]>([]);
  const [selectedTab, setSelectedTab] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Dialog states
  const [openComponentDialog, setOpenComponentDialog] = useState(false);
  const [openTemplateDialog, setOpenTemplateDialog] = useState(false);
  const [openUsageDialog, setOpenUsageDialog] = useState(false);
  const [selectedComponent, setSelectedComponent] = useState<CustomComponent | null>(null);
  const [selectedTemplate, setSelectedTemplate] = useState<ComponentTemplate | null>(null);

  // Component form state
  const [componentForm, setComponentForm] = useState({
    name: '',
    display_name: '',
    description: '',
    component_type: 'category' as 'category' | 'attribute' | 'filter',
    component_code: '',
    configuration: '{}',
    dependencies: '{}',
    is_active: true,
  });

  // Template form state
  const [templateForm, setTemplateForm] = useState({
    name: '',
    display_name: '',
    description: '',
    template_code: '',
    variables: '{}',
    is_shared: false,
  });

  useEffect(() => {
    fetchComponents();
    fetchTemplates();
  }, []);

  const fetchComponents = async () => {
    setLoading(true);
    setError(null); // Сбрасываем ошибку
    try {
      const response = await axios.get('/api/v1/admin/custom-components');
      console.log('API response:', response.data); // Для отладки
      
      // Проверяем, что данные - это массив или null/undefined
      if (response.data === null || response.data === undefined) {
        setComponents([]);
      } else if (Array.isArray(response.data)) {
        setComponents(response.data);
      } else if (response.data && Array.isArray(response.data.components)) {
        // Если API возвращает объект с массивом components
        setComponents(response.data.components);
      } else {
        // Если нет данных или неверный формат, устанавливаем пустой массив
        setComponents([]);
        console.warn('Unexpected API response format:', response.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка загрузки компонентов');
      console.error('Error fetching components:', err);
      setComponents([]); // Устанавливаем пустой массив при ошибке
    } finally {
      setLoading(false);
    }
  };

  const fetchTemplates = async () => {
    try {
      const response = await axios.get('/api/v1/admin/custom-components/templates');
      // Проверяем, что данные - это массив
      if (Array.isArray(response.data)) {
        setTemplates(response.data);
      } else {
        setTemplates([]);
        console.warn('Unexpected templates API response format:', response.data);
      }
    } catch (err) {
      console.error('Ошибка загрузки шаблонов:', err);
      setTemplates([]); // Устанавливаем пустой массив при ошибке
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  const handleCreateComponent = () => {
    setSelectedComponent(null);
    setComponentForm({
      name: '',
      display_name: '',
      description: '',
      component_type: 'category',
      component_code: '',
      configuration: '{}',
      dependencies: '{}',
      is_active: true,
    });
    setOpenComponentDialog(true);
  };

  const handleEditComponent = (component: CustomComponent) => {
    setSelectedComponent(component);
    setComponentForm({
      name: component.name,
      display_name: component.display_name,
      description: component.description,
      component_type: component.component_type,
      component_code: component.component_code,
      configuration: JSON.stringify(component.configuration, null, 2),
      dependencies: JSON.stringify(component.dependencies, null, 2),
      is_active: component.is_active,
    });
    setOpenComponentDialog(true);
  };

  const handleDeleteComponent = async (id: number) => {
    if (window.confirm('Удалить компонент?')) {
      try {
        await axios.delete(`/api/v1/admin/custom-components/${id}`);
        fetchComponents();
      } catch (err) {
        setError('Ошибка удаления компонента');
      }
    }
  };

  const handleSaveComponent = async () => {
    try {
      const data = {
        ...componentForm,
        configuration: JSON.parse(componentForm.configuration),
        dependencies: JSON.parse(componentForm.dependencies),
      };

      if (selectedComponent) {
        await axios.put(`/api/v1/admin/custom-components/${selectedComponent.id}`, data);
      } else {
        await axios.post('/api/v1/admin/custom-components', data);
      }

      setOpenComponentDialog(false);
      fetchComponents();
    } catch (err) {
      setError('Ошибка сохранения компонента');
    }
  };

  const handleCreateTemplate = () => {
    setSelectedTemplate(null);
    setTemplateForm({
      name: '',
      display_name: '',
      description: '',
      template_code: '',
      variables: '{}',
      is_shared: false,
    });
    setOpenTemplateDialog(true);
  };

  const handleSaveTemplate = async () => {
    try {
      const data = {
        ...templateForm,
        variables: JSON.parse(templateForm.variables),
      };

      if (selectedTemplate) {
        await axios.put(`/api/v1/admin/custom-components/templates/${selectedTemplate.id}`, data);
      } else {
        await axios.post('/api/v1/admin/custom-components/templates', data);
      }

      setOpenTemplateDialog(false);
      fetchTemplates();
    } catch (err) {
      setError('Ошибка сохранения шаблона');
    }
  };

  const ComponentsTab = () => {
    if (loading) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
          <CircularProgress />
        </Box>
      );
    }

    if (error) {
      return (
        <Box p={3}>
          <Alert severity="error">{error}</Alert>
        </Box>
      );
    }

    return (
      <>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h5">Кастомные компоненты</Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={handleCreateComponent}
          >
            Создать компонент
          </Button>
        </Box>

        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Отображаемое имя</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Обновлено</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {components && components.length > 0 ? (
                components.map((component) => (
              <TableRow key={component.id}>
                <TableCell>{component.name}</TableCell>
                <TableCell>{component.display_name}</TableCell>
                <TableCell>
                  <ComponentTypeBadge
                    label={component.component_type}
                    componentType={component.component_type}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Chip
                    label={component.is_active ? 'Активен' : 'Неактивен'}
                    color={component.is_active ? 'success' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  {new Date(component.updated_at).toLocaleDateString()}
                </TableCell>
                <TableCell>
                  <Tooltip title="Редактировать">
                    <IconButton onClick={() => handleEditComponent(component)}>
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Настроить использование">
                    <IconButton>
                      <SettingsIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Удалить">
                    <IconButton onClick={() => handleDeleteComponent(component.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </Tooltip>
                </TableCell>
              </TableRow>
                ))) : (
                  <TableRow>
                    <TableCell colSpan={6} align="center">
                      <Typography variant="body2" color="textSecondary">
                        Нет кастомных компонентов. Создайте первый компонент.
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </>
      );
    };

  const TemplatesTab = () => (
    <>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h5">Шаблоны компонентов</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={handleCreateTemplate}
        >
          Создать шаблон
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Название</TableCell>
              <TableCell>Отображаемое имя</TableCell>
              <TableCell>Описание</TableCell>
              <TableCell>Общий</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {templates.map((template) => (
              <TableRow key={template.id}>
                <TableCell>{template.name}</TableCell>
                <TableCell>{template.display_name}</TableCell>
                <TableCell>{template.description}</TableCell>
                <TableCell>
                  <Chip
                    label={template.is_shared ? 'Да' : 'Нет'}
                    color={template.is_shared ? 'primary' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Tooltip title="Редактировать">
                    <IconButton>
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Предпросмотр">
                    <IconButton>
                      <PreviewIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Удалить">
                    <IconButton>
                      <DeleteIcon />
                    </IconButton>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );

  return (
    <Box>
      <StyledPaper>
        <Typography variant="h4" gutterBottom>
          Управление кастомными компонентами
        </Typography>
        
        {error && (
          <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Tabs value={selectedTab} onChange={handleTabChange} sx={{ mb: 3 }}>
          <Tab label="Компоненты" />
          <Tab label="Шаблоны" />
          <Tab label="Использование" />
        </Tabs>

        {selectedTab === 0 && <ComponentsTab />}
        {selectedTab === 1 && <TemplatesTab />}
        {selectedTab === 2 && <div>Использование компонентов (в разработке)</div>}
      </StyledPaper>

      {/* Component Dialog */}
      <Dialog
        open={openComponentDialog}
        onClose={() => setOpenComponentDialog(false)}
        maxWidth="lg"
        fullWidth
      >
        <DialogTitle>
          {selectedComponent ? 'Редактировать компонент' : 'Создать компонент'}
        </DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} sx={{ mt: 2 }}>
            <TextField
              label="Имя компонента"
              value={componentForm.name}
              onChange={(e) => setComponentForm({ ...componentForm, name: e.target.value })}
              fullWidth
              required
            />
            
            <TextField
              label="Отображаемое имя"
              value={componentForm.display_name}
              onChange={(e) => setComponentForm({ ...componentForm, display_name: e.target.value })}
              fullWidth
              required
            />
            
            <TextField
              label="Описание"
              value={componentForm.description}
              onChange={(e) => setComponentForm({ ...componentForm, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
            
            <FormControl fullWidth>
              <InputLabel>Тип компонента</InputLabel>
              <Select
                value={componentForm.component_type}
                onChange={(e) => setComponentForm({ ...componentForm, component_type: e.target.value as 'category' | 'attribute' | 'filter' })}
                label="Тип компонента"
              >
                <MenuItem value="category">Категория</MenuItem>
                <MenuItem value="attribute">Атрибут</MenuItem>
                <MenuItem value="filter">Фильтр</MenuItem>
              </Select>
            </FormControl>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Код компонента
              </Typography>
              <CodeEditor
                value={componentForm.component_code}
                language="tsx"
                placeholder="Введите код React компонента"
                onChange={(evn) => setComponentForm({ ...componentForm, component_code: evn.target.value })}
                padding={15}
                style={{
                  fontSize: 14,
                  backgroundColor: '#f5f5f5',
                  fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                }}
              />
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Конфигурация (JSON)
              </Typography>
              <CodeEditor
                value={componentForm.configuration}
                language="json"
                placeholder="{}"
                onChange={(evn) => setComponentForm({ ...componentForm, configuration: evn.target.value })}
                padding={15}
                style={{
                  fontSize: 14,
                  backgroundColor: '#f5f5f5',
                  fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                }}
              />
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Зависимости (JSON)
              </Typography>
              <CodeEditor
                value={componentForm.dependencies}
                language="json"
                placeholder="{}"
                onChange={(evn) => setComponentForm({ ...componentForm, dependencies: evn.target.value })}
                padding={15}
                style={{
                  fontSize: 14,
                  backgroundColor: '#f5f5f5',
                  fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                }}
              />
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenComponentDialog(false)}>Отмена</Button>
          <Button onClick={handleSaveComponent} variant="contained" color="primary">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>

      {/* Template Dialog */}
      <Dialog
        open={openTemplateDialog}
        onClose={() => setOpenTemplateDialog(false)}
        maxWidth="lg"
        fullWidth
      >
        <DialogTitle>
          {selectedTemplate ? 'Редактировать шаблон' : 'Создать шаблон'}
        </DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} sx={{ mt: 2 }}>
            <TextField
              label="Имя шаблона"
              value={templateForm.name}
              onChange={(e) => setTemplateForm({ ...templateForm, name: e.target.value })}
              fullWidth
              required
            />
            
            <TextField
              label="Отображаемое имя"
              value={templateForm.display_name}
              onChange={(e) => setTemplateForm({ ...templateForm, display_name: e.target.value })}
              fullWidth
              required
            />
            
            <TextField
              label="Описание"
              value={templateForm.description}
              onChange={(e) => setTemplateForm({ ...templateForm, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Код шаблона
              </Typography>
              <CodeEditor
                value={templateForm.template_code}
                language="tsx"
                placeholder="Введите код шаблона"
                onChange={(evn) => setTemplateForm({ ...templateForm, template_code: evn.target.value })}
                padding={15}
                style={{
                  fontSize: 14,
                  backgroundColor: '#f5f5f5',
                  fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                }}
              />
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Переменные (JSON)
              </Typography>
              <CodeEditor
                value={templateForm.variables}
                language="json"
                placeholder="{}"
                onChange={(evn) => setTemplateForm({ ...templateForm, variables: evn.target.value })}
                padding={15}
                style={{
                  fontSize: 14,
                  backgroundColor: '#f5f5f5',
                  fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                }}
              />
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenTemplateDialog(false)}>Отмена</Button>
          <Button onClick={handleSaveTemplate} variant="contained" color="primary">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CustomComponentsPage;