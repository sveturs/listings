import React from 'react';
import { TextField, Checkbox, FormControlLabel, MenuItem, FormControl, InputLabel, Select, Slider, Typography, Box } from '@mui/material';
import ComponentRegistry, { AttributeComponentProps } from './ComponentRegistry';

/**
 * Компонент для отображения атрибута типа "text"
 */
const TextAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  return (
    <TextField
      fullWidth
      label={attribute.display_name}
      value={value || ''}
      onChange={(e) => onChange(e.target.value)}
      margin="normal"
      size="small"
    />
  );
};

/**
 * Компонент для отображения атрибута типа "number"
 */
const NumberAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const numValue = e.target.value === '' ? '' : Number(e.target.value);
    onChange(numValue);
  };

  const min = attribute.options?.min !== undefined ? attribute.options.min : undefined;
  const max = attribute.options?.max !== undefined ? attribute.options.max : undefined;
  const step = attribute.options?.step !== undefined ? attribute.options.step : 1;

  return (
    <TextField
      fullWidth
      label={attribute.display_name}
      type="number"
      value={value !== undefined && value !== null ? value : ''}
      onChange={handleChange}
      margin="normal"
      size="small"
      inputProps={{
        min: min,
        max: max,
        step: step
      }}
    />
  );
};

/**
 * Компонент для отображения атрибута типа "select"
 */
const SelectAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  const options = attribute.options?.values || [];
  
  return (
    <FormControl fullWidth margin="normal" size="small">
      <InputLabel>{attribute.display_name}</InputLabel>
      <Select
        value={value || ''}
        onChange={(e) => onChange(e.target.value)}
        label={attribute.display_name}
      >
        <MenuItem value="">
          <em>Не выбрано</em>
        </MenuItem>
        {options.map((option: string) => (
          <MenuItem key={option} value={option}>
            {option}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

/**
 * Компонент для отображения атрибута типа "multiselect"
 */
const MultiSelectAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  const options = attribute.options?.values || [];
  const selectedValues = Array.isArray(value) ? value : [];
  
  const handleChange = (event: any) => {
    const {
      target: { value: newValue },
    } = event;
    onChange(typeof newValue === 'string' ? newValue.split(',') : newValue);
  };
  
  return (
    <FormControl fullWidth margin="normal" size="small">
      <InputLabel>{attribute.display_name}</InputLabel>
      <Select
        multiple
        value={selectedValues}
        onChange={handleChange}
        label={attribute.display_name}
      >
        {options.map((option: string) => (
          <MenuItem key={option} value={option}>
            {option}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

/**
 * Компонент для отображения атрибута типа "boolean"
 */
const BooleanAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  return (
    <FormControlLabel
      control={
        <Checkbox
          checked={!!value}
          onChange={(e) => onChange(e.target.checked)}
        />
      }
      label={attribute.display_name}
    />
  );
};

/**
 * Компонент для отображения атрибута типа "range"
 */
const RangeAttributeInput: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  const min = attribute.options?.min !== undefined ? attribute.options.min : 0;
  const max = attribute.options?.max !== undefined ? attribute.options.max : 100;
  const step = attribute.options?.step !== undefined ? attribute.options.step : 1;
  
  const initialValue = value !== undefined ? value : [min, max];
  
  return (
    <Box sx={{ width: '100%', padding: '10px 0', mt: 2 }}>
      <Typography variant="body2" gutterBottom>
        {attribute.display_name}
      </Typography>
      <Slider
        value={initialValue}
        onChange={(e, newValue) => onChange(newValue)}
        valueLabelDisplay="auto"
        min={min}
        max={max}
        step={step}
      />
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
        <Typography variant="caption" color="text.secondary">
          {initialValue[0]}
        </Typography>
        <Typography variant="caption" color="text.secondary">
          {initialValue[1]}
        </Typography>
      </Box>
    </Box>
  );
};

/**
 * Фабрика компонентов для атрибутов
 * Выбирает подходящий компонент в зависимости от типа атрибута
 */
const AttributeComponentFactory: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  // Проверяем наличие кастомного компонента
  if (attribute.custom_component) {
    const CustomComponent = ComponentRegistry.getAttributeComponent(attribute.custom_component);
    if (CustomComponent) {
      return <CustomComponent attribute={attribute} value={value} onChange={onChange} />;
    }
  }

  // Выбираем стандартный компонент в зависимости от типа атрибута
  switch (attribute.attribute_type) {
    case 'text':
      return <TextAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'number':
      return <NumberAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'select':
      return <SelectAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'multiselect':
      return <MultiSelectAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'boolean':
      return <BooleanAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'range':
      return <RangeAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    default:
      return <TextAttributeInput attribute={attribute} value={value} onChange={onChange} />;
  }
};

export default AttributeComponentFactory;

// Регистрируем стандартные компоненты
ComponentRegistry.registerAttributeComponent('TextAttributeInput', TextAttributeInput);
ComponentRegistry.registerAttributeComponent('NumberAttributeInput', NumberAttributeInput);
ComponentRegistry.registerAttributeComponent('SelectAttributeInput', SelectAttributeInput);
ComponentRegistry.registerAttributeComponent('MultiSelectAttributeInput', MultiSelectAttributeInput);
ComponentRegistry.registerAttributeComponent('BooleanAttributeInput', BooleanAttributeInput);
ComponentRegistry.registerAttributeComponent('RangeAttributeInput', RangeAttributeInput);