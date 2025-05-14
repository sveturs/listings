import React from 'react';
import {
  Box,
  Typography,
  Chip,
  Grid,
  Paper,
  Divider,
  styled
} from '@mui/material';
import { useTranslation } from 'react-i18next';

// Названия типов атрибутов для более читаемого отображения
const attributeTypeNames: Record<string, string> = {
  text: 'Текст',
  number: 'Число',
  select: 'Выбор',
  multiselect: 'Мульти-выбор',
  boolean: 'Да/Нет',
  range: 'Диапазон',
  date: 'Дата'
};

// Стилизованная таблица для отображения атрибутов
const AttributeTable = styled(Box)(({ theme }) => ({
  width: '100%',
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(2),
  '& .attribute-row': {
    display: 'flex',
    borderBottom: `1px solid ${theme.palette.divider}`,
    '&:last-child': {
      borderBottom: 'none'
    },
    '& .attribute-name': {
      width: '40%',
      padding: theme.spacing(1.5),
      backgroundColor: theme.palette.action.hover,
      fontWeight: 500
    },
    '& .attribute-value': {
      width: '60%',
      padding: theme.spacing(1.5),
      '& .chip-list': {
        display: 'flex',
        flexWrap: 'wrap',
        gap: theme.spacing(0.5)
      }
    }
  }
}));

interface AttributeDisplayProps {
  attributes: Array<{
    attribute_id: number;
    attribute_name: string;
    display_name: string;
    attribute_type: string;
    text_value?: string | null;
    numeric_value?: number | null;
    boolean_value?: boolean | null;
    json_value?: any;
    display_value: string;
    unit?: string;
    custom_component?: string;
  }>;
  showTable?: boolean;
  showHeader?: boolean;
  groupByType?: boolean;
}

/**
 * Компонент для отображения атрибутов объявления
 */
const AttributeDisplay: React.FC<AttributeDisplayProps> = ({ 
  attributes, 
  showTable = true,
  showHeader = true,
  groupByType = false
}) => {
  const { t } = useTranslation();

  // Если нет атрибутов, ничего не отображаем
  if (!attributes || attributes.length === 0) {
    return null;
  }

  // Если нужно группировать по типу, подготавливаем группы
  const getGroupedAttributes = () => {
    const groups: Record<string, typeof attributes> = {};
    
    attributes.forEach(attr => {
      const groupName = attributeTypeNames[attr.attribute_type] || 'Прочее';
      if (!groups[groupName]) {
        groups[groupName] = [];
      }
      groups[groupName].push(attr);
    });
    
    return groups;
  };

  // Функция для отображения значения атрибута в зависимости от типа
  const renderAttributeValue = (attribute: AttributeDisplayProps['attributes'][0]) => {
    // Для булевых значений
    if (attribute.attribute_type === 'boolean') {
      return (
        <Chip 
          label={attribute.boolean_value ? t('common.yes') : t('common.no')}
          color={attribute.boolean_value ? 'success' : 'default'}
          size="small"
        />
      );
    }
    
    // Для мульти-выбора
    if (attribute.attribute_type === 'multiselect' && attribute.json_value) {
      const values = Array.isArray(attribute.json_value) 
        ? attribute.json_value 
        : typeof attribute.json_value === 'string'
          ? JSON.parse(attribute.json_value)
          : [];
          
      return (
        <Box className="chip-list">
          {values.map((value: string, index: number) => (
            <Chip key={index} label={value} size="small" />
          ))}
        </Box>
      );
    }
    
    // Для выбора из списка
    if (attribute.attribute_type === 'select') {
      return <Chip label={attribute.text_value || attribute.display_value} size="small" />;
    }
    
    // Для числовых значений с единицами измерения
    if (attribute.attribute_type === 'number' && attribute.unit) {
      return (
        <Typography variant="body2">
          {attribute.numeric_value} {attribute.unit}
        </Typography>
      );
    }
    
    // Для всех остальных случаев используем display_value
    return (
      <Typography variant="body2">
        {attribute.display_value}
      </Typography>
    );
  };

  // Отображаем в виде таблицы
  if (showTable) {
    if (groupByType) {
      const groups = getGroupedAttributes();
      
      return (
        <Box>
          {showHeader && (
            <Typography variant="subtitle1" gutterBottom>
              {t('marketplace.attributes.title')}
            </Typography>
          )}
          
          {Object.entries(groups).map(([groupName, groupAttributes]) => (
            <Box key={groupName} sx={{ mb: 3 }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>
                {groupName}
              </Typography>
              
              <AttributeTable>
                {groupAttributes.map(attribute => (
                  <Box key={attribute.attribute_id} className="attribute-row">
                    <Box className="attribute-name">
                      {attribute.display_name}
                    </Box>
                    <Box className="attribute-value">
                      {renderAttributeValue(attribute)}
                    </Box>
                  </Box>
                ))}
              </AttributeTable>
            </Box>
          ))}
        </Box>
      );
    }
    
    return (
      <Box>
        {showHeader && (
          <Typography variant="subtitle1" gutterBottom>
            {t('marketplace.attributes.title')}
          </Typography>
        )}
        
        <AttributeTable>
          {attributes.map(attribute => (
            <Box key={attribute.attribute_id} className="attribute-row">
              <Box className="attribute-name">
                {attribute.display_name}
              </Box>
              <Box className="attribute-value">
                {renderAttributeValue(attribute)}
              </Box>
            </Box>
          ))}
        </AttributeTable>
      </Box>
    );
  }
  
  // Отображаем в виде сетки чипов
  return (
    <Box>
      {showHeader && (
        <Typography variant="subtitle1" gutterBottom>
          {t('marketplace.attributes.title')}
        </Typography>
      )}
      
      <Grid container spacing={1}>
        {attributes.map(attribute => (
          <Grid item key={attribute.attribute_id}>
            <Chip
              label={`${attribute.display_name}: ${attribute.display_value}`}
              size="small"
              variant="outlined"
            />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default AttributeDisplay;