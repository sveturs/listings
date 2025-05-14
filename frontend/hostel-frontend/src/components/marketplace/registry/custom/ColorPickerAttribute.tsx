import React, { useState, useEffect } from 'react';
import { Box, Typography, Grid, Paper, styled, FormHelperText } from '@mui/material';
import { AttributeComponentProps } from '../ComponentRegistry';

// Список предопределенных цветов
const PRESET_COLORS = [
  { name: 'Белый', value: '#FFFFFF', hex: '#FFFFFF', border: '#CCCCCC' },
  { name: 'Черный', value: 'Черный', hex: '#000000', border: '#000000' },
  { name: 'Серый', value: 'Серый', hex: '#808080', border: '#808080' },
  { name: 'Серебристый', value: 'Серебристый', hex: '#C0C0C0', border: '#C0C0C0' },
  { name: 'Красный', value: 'Красный', hex: '#FF0000', border: '#FF0000' },
  { name: 'Синий', value: 'Синий', hex: '#0000FF', border: '#0000FF' },
  { name: 'Зеленый', value: 'Зеленый', hex: '#008000', border: '#008000' },
  { name: 'Желтый', value: 'Желтый', hex: '#FFFF00', border: '#CCCC00' },
  { name: 'Коричневый', value: 'Коричневый', hex: '#A52A2A', border: '#A52A2A' },
  { name: 'Оранжевый', value: 'Оранжевый', hex: '#FFA500', border: '#FFA500' },
  { name: 'Фиолетовый', value: 'Фиолетовый', hex: '#800080', border: '#800080' },
  { name: 'Розовый', value: 'Розовый', hex: '#FFC0CB', border: '#FFC0CB' },
  { name: 'Золотой', value: 'Золотой', hex: '#FFD700', border: '#D4AF37' },
  { name: 'Бежевый', value: 'Бежевый', hex: '#F5F5DC', border: '#D6CDB7' },
];

// Стилизованный компонент для цветного квадрата
const ColorBox = styled(Paper)<{ selected?: boolean, hex: string, borderhex: string }>(
  ({ selected, hex, borderhex }) => ({
    width: 36,
    height: 36,
    backgroundColor: hex,
    borderRadius: 4,
    cursor: 'pointer',
    border: selected ? '3px solid #2196F3' : `1px solid ${borderhex}`,
    boxShadow: selected ? '0 0 5px #2196F3' : 'none',
    '&:hover': {
      boxShadow: '0 0 5px rgba(0, 0, 0, 0.3)',
    },
  })
);

/**
 * Кастомный компонент для выбора цвета
 * Отображает цветовую палитру с названиями цветов
 */
const ColorPickerAttribute: React.FC<AttributeComponentProps> = ({ attribute, value, onChange }) => {
  // Получаем доступные цвета из опций атрибута или используем предопределенные
  const [colors, setColors] = useState(PRESET_COLORS);
  const [selectedColor, setSelectedColor] = useState(value || '');

  useEffect(() => {
    // Если у атрибута есть опции с цветами, используем их
    if (attribute.options && attribute.options.values) {
      try {
        const attrOptions = Array.isArray(attribute.options.values) 
          ? attribute.options.values 
          : typeof attribute.options.values === 'string' 
            ? JSON.parse(attribute.options.values) 
            : [];
            
        if (attrOptions.length > 0) {
          // Если опции - просто строки с названием цвета, находим соответствие в предопределенных
          if (typeof attrOptions[0] === 'string') {
            const customColors = attrOptions.map((colorName: string) => {
              const presetMatch = PRESET_COLORS.find(c => c.name.toLowerCase() === colorName.toLowerCase());
              return presetMatch || { 
                name: colorName, 
                value: colorName, 
                hex: '#CCCCCC', 
                border: '#AAAAAA' 
              };
            });
            setColors(customColors);
          } 
          // Если опции - объекты с названием и hex-кодом, используем их
          else if (typeof attrOptions[0] === 'object') {
            const customColors = attrOptions.map((colorObj: any) => ({
              name: colorObj.name || colorObj.value || 'Цвет',
              value: colorObj.value || colorObj.name || 'Цвет',
              hex: colorObj.hex || '#CCCCCC',
              border: colorObj.border || colorObj.hex || '#AAAAAA'
            }));
            setColors(customColors);
          }
        }
      } catch (error) {
        console.error('Error parsing color options:', error);
      }
    }
  }, [attribute.options]);

  // Обрабатываем выбор цвета
  const handleColorClick = (color: string) => {
    setSelectedColor(color);
    onChange(color);
  };

  return (
    <Box sx={{ mt: 2, mb: 2 }}>
      <Typography variant="subtitle2" gutterBottom>
        {attribute.display_name}
      </Typography>
      
      <Grid container spacing={1} sx={{ mt: 1 }}>
        {colors.map((color) => (
          <Grid item key={color.value}>
            <Box sx={{ 
              display: 'flex', 
              flexDirection: 'column', 
              alignItems: 'center' 
            }}>
              <ColorBox 
                hex={color.hex} 
                borderhex={color.border}
                selected={selectedColor === color.value}
                onClick={() => handleColorClick(color.value)}
              />
              <Typography variant="caption" sx={{ mt: 0.5 }}>
                {color.name}
              </Typography>
            </Box>
          </Grid>
        ))}
      </Grid>
      
      {selectedColor && (
        <FormHelperText>
          Выбран цвет: {colors.find(c => c.value === selectedColor)?.name || selectedColor}
        </FormHelperText>
      )}
    </Box>
  );
};

export default ColorPickerAttribute;