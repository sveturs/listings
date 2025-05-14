import React, { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Chip,
  Box,
  Typography,
  Tooltip,
  TablePagination,
  TableSortLabel
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Code as CodeIcon,
  Check as CheckIcon,
  Close as CloseIcon,
  List as ListIcon,
  Numbers as NumbersIcon,
  TextFields as TextFieldsIcon,
  ToggleOn as ToggleOnIcon,
  DateRange as DateRangeIcon
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { Attribute } from '../../pages/admin/AttributeManagementPage';

interface AttributesTableProps {
  attributes: Attribute[];
  onEdit: (attribute: Attribute) => void;
  onDelete: (attributeId: number) => void;
}

type Order = 'asc' | 'desc';
type SortColumn = 
  | 'id' 
  | 'name' 
  | 'display_name' 
  | 'attribute_type' 
  | 'is_searchable' 
  | 'is_filterable' 
  | 'is_required'
  | 'sort_order';

const AttributesTable: React.FC<AttributesTableProps> = ({ attributes, onEdit, onDelete }) => {
  const { t } = useTranslation();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [order, setOrder] = useState<Order>('asc');
  const [orderBy, setOrderBy] = useState<SortColumn>('sort_order');

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRequestSort = (column: SortColumn) => {
    const isAsc = orderBy === column && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(column);
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'text':
        return <TextFieldsIcon fontSize="small" />;
      case 'number':
      case 'range':
        return <NumbersIcon fontSize="small" />;
      case 'select':
      case 'multiselect':
        return <ListIcon fontSize="small" />;
      case 'boolean':
        return <ToggleOnIcon fontSize="small" />;
      case 'date':
        return <DateRangeIcon fontSize="small" />;
      default:
        return <TextFieldsIcon fontSize="small" />;
    }
  };

  // Функция сортировки с защитой от null и undefined
  const sortAttributes = (attrA: Attribute, attrB: Attribute, column: SortColumn, sortOrder: Order) => {
    // Проверка на случай, если какой-то из атрибутов undefined или null
    if (!attrA || !attrB) return 0;

    let a = column in attrA ? attrA[column] : null;
    let b = column in attrB ? attrB[column] : null;

    // Если оба значения пустые или undefined, считаем их равными
    if ((a === null || a === undefined) && (b === null || b === undefined)) {
      return 0;
    }

    // Если только одно из значений отсутствует
    if (a === null || a === undefined) return sortOrder === 'asc' ? 1 : -1;
    if (b === null || b === undefined) return sortOrder === 'asc' ? -1 : 1;

    // Для текстовых полей
    if (typeof a === 'string' && typeof b === 'string') {
      return sortOrder === 'asc'
        ? a.localeCompare(b)
        : b.localeCompare(a);
    }

    // Для булевых и числовых полей
    if (a === b) return 0;
    if (sortOrder === 'asc') {
      return a < b ? -1 : 1;
    } else {
      return a < b ? 1 : -1;
    }
  };

  // Проверяем, что атрибуты не пустые и являются массивом
  const validAttributes = Array.isArray(attributes) ? attributes : [];

  // Применяем сортировку к списку атрибутов
  const sortedAttributes = [...validAttributes].sort((a, b) =>
    sortAttributes(a, b, orderBy, order)
  );

  // Применяем пагинацию
  const paginatedAttributes = sortedAttributes.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  return (
    <Paper variant="outlined">
      <TableContainer component={Paper} elevation={0}>
        <Table sx={{ minWidth: 650 }} size="medium">
          <TableHead>
            <TableRow>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'id'}
                  direction={orderBy === 'id' ? order : 'asc'}
                  onClick={() => handleRequestSort('id')}
                >
                  ID
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'name'}
                  direction={orderBy === 'name' ? order : 'asc'}
                  onClick={() => handleRequestSort('name')}
                >
                  {t('admin.attributes.name')}
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'display_name'}
                  direction={orderBy === 'display_name' ? order : 'asc'}
                  onClick={() => handleRequestSort('display_name')}
                >
                  {t('admin.attributes.displayName')}
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'attribute_type'}
                  direction={orderBy === 'attribute_type' ? order : 'asc'}
                  onClick={() => handleRequestSort('attribute_type')}
                >
                  {t('admin.attributes.type')}
                </TableSortLabel>
              </TableCell>
              <TableCell align="center">
                <TableSortLabel
                  active={orderBy === 'is_searchable'}
                  direction={orderBy === 'is_searchable' ? order : 'asc'}
                  onClick={() => handleRequestSort('is_searchable')}
                >
                  {t('admin.attributes.isSearchable')}
                </TableSortLabel>
              </TableCell>
              <TableCell align="center">
                <TableSortLabel
                  active={orderBy === 'is_filterable'}
                  direction={orderBy === 'is_filterable' ? order : 'asc'}
                  onClick={() => handleRequestSort('is_filterable')}
                >
                  {t('admin.attributes.isFilterable')}
                </TableSortLabel>
              </TableCell>
              <TableCell align="center">
                <TableSortLabel
                  active={orderBy === 'is_required'}
                  direction={orderBy === 'is_required' ? order : 'asc'}
                  onClick={() => handleRequestSort('is_required')}
                >
                  {t('admin.attributes.isRequired')}
                </TableSortLabel>
              </TableCell>
              <TableCell align="center">
                <TableSortLabel
                  active={orderBy === 'sort_order'}
                  direction={orderBy === 'sort_order' ? order : 'asc'}
                  onClick={() => handleRequestSort('sort_order')}
                >
                  {t('admin.attributes.sortOrder')}
                </TableSortLabel>
              </TableCell>
              <TableCell align="right">{t('common.actions')}</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedAttributes.length > 0 ? (
              paginatedAttributes.map((attribute) => (
                <TableRow key={attribute.id} hover>
                  <TableCell>{attribute.id}</TableCell>
                  <TableCell>
                    {attribute.name}
                    {attribute.custom_component && (
                      <Tooltip title={t('admin.attributes.hasCustomComponent')}>
                        <Chip
                          icon={<CodeIcon />}
                          label={attribute.custom_component}
                          size="small"
                          color="primary"
                          variant="outlined"
                          sx={{ ml: 1 }}
                        />
                      </Tooltip>
                    )}
                  </TableCell>
                  <TableCell>{attribute.display_name}</TableCell>
                  <TableCell>
                    <Chip
                      icon={getTypeIcon(attribute.attribute_type)}
                      label={t(`admin.attributes.types.${attribute.attribute_type}`, attribute.attribute_type)}
                      size="small"
                      color="default"
                    />
                  </TableCell>
                  <TableCell align="center">
                    {attribute.is_searchable ? (
                      <CheckIcon color="success" fontSize="small" />
                    ) : (
                      <CloseIcon color="error" fontSize="small" />
                    )}
                  </TableCell>
                  <TableCell align="center">
                    {attribute.is_filterable ? (
                      <CheckIcon color="success" fontSize="small" />
                    ) : (
                      <CloseIcon color="error" fontSize="small" />
                    )}
                  </TableCell>
                  <TableCell align="center">
                    {attribute.is_required ? (
                      <CheckIcon color="success" fontSize="small" />
                    ) : (
                      <CloseIcon color="error" fontSize="small" />
                    )}
                  </TableCell>
                  <TableCell align="center">{attribute.sort_order}</TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => onEdit(attribute)}
                      aria-label={t('common.edit')}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => onDelete(attribute.id)}
                      aria-label={t('common.delete')}
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={9} align="center">
                  <Typography variant="body1" sx={{ py: 2 }}>
                    {t('admin.attributes.noAttributes')}
                  </Typography>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[5, 10, 25, 50]}
        component="div"
        count={validAttributes.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        labelRowsPerPage={t('common.rowsPerPage')}
      />
    </Paper>
  );
};

export default AttributesTable;