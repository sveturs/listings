import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  IconButton,
  TextField,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  CircularProgress,
  Alert,
  Chip,
  Tooltip,
  TablePagination,
  SelectChangeEvent
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Block as BlockIcon,
  CheckCircle as CheckCircleIcon,
  ArrowBack as ArrowBackIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { useNavigate } from 'react-router-dom';
import { formatDate } from '../../utils/dateUtils';
import UserDetailsDialog from '../../components/admin/UserDetailsDialog';

interface User {
  id: number;
  name: string;
  email: string;
  phone?: string;
  account_status?: string;
  created_at: string;
  [key: string]: any;
}

interface EditFormData {
  name: string;
  email: string;
  phone: string;
  status: string;
}

interface ResponseData {
  data?: User[];
  total?: number;
  message?: string;
  [key: string]: any;
}

const UsersManagementPage: React.FC = () => {
  const { t } = useTranslation(['common', 'admin']);
  const navigate = useNavigate();

  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState<number>(0);
  const [rowsPerPage, setRowsPerPage] = useState<number>(10);
  const [totalUsers, setTotalUsers] = useState<number>(0);

  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [openEditDialog, setOpenEditDialog] = useState<boolean>(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState<boolean>(false);
  const [openDetailsDialog, setOpenDetailsDialog] = useState<boolean>(false);

  const [searchTerm, setSearchTerm] = useState<string>('');
  const [searchResults, setSearchResults] = useState<User[]>([]);

  const [editFormData, setEditFormData] = useState<EditFormData>({
    name: '',
    email: '',
    phone: '',
    status: ''
  });

  const [updateLoading, setUpdateLoading] = useState<boolean>(false);
  const [updateError, setUpdateError] = useState<string | null>(null);
  const [updateSuccess, setUpdateSuccess] = useState<boolean>(false);

  // Загрузка пользователей
  const fetchUsers = async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const response = await axios.get<ResponseData>(`/api/v1/admin/users?page=${page + 1}&limit=${rowsPerPage}`);
      setUsers(response.data.data || []);
      setTotalUsers(response.data.total || 0);
      setSearchResults(response.data.data || []);
    } catch (err: any) {
      console.error('Ошибка при загрузке пользователей:', err);
      setError(err.response?.data?.message || err.message || 'Произошла ошибка при загрузке пользователей');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, [page, rowsPerPage]);

  // Поиск пользователей
  useEffect(() => {
    if (searchTerm.trim() === '') {
      setSearchResults(users);
    } else {
      const filteredResults = users.filter(user =>
        user.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.phone?.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setSearchResults(filteredResults);
    }
  }, [searchTerm, users]);

  // Обработчики пагинации
  const handleChangePage = (_event: React.MouseEvent<HTMLButtonElement> | null, newPage: number): void => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>): void => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  // Обработчики диалогов
  const handleOpenEditDialog = (user: User): void => {
    setSelectedUser(user);
    setEditFormData({
      name: user.name || '',
      email: user.email || '',
      phone: user.phone || '',
      status: user.account_status || 'active'
    });
    setOpenEditDialog(true);
  };

  const handleCloseEditDialog = (): void => {
    setOpenEditDialog(false);
    setUpdateError(null);
    setUpdateSuccess(false);
  };

  const handleOpenDeleteDialog = (user: User): void => {
    setSelectedUser(user);
    setOpenDeleteDialog(true);
  };

  const handleCloseDeleteDialog = (): void => {
    setOpenDeleteDialog(false);
  };

  const handleOpenDetailsDialog = (user: User): void => {
    setSelectedUser(user);
    setOpenDetailsDialog(true);
  };

  const handleCloseDetailsDialog = (): void => {
    setOpenDetailsDialog(false);
  };

  // Обработчики формы редактирования
  const handleEditFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>): void => {
    const { name, value } = e.target;
    setEditFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>): void => {
    setEditFormData(prev => ({
      ...prev,
      status: e.target.value
    }));
  };

  // Обновление пользователя
  const handleUpdateUser = async (): Promise<void> => {
    if (!selectedUser) return;

    setUpdateLoading(true);
    setUpdateError(null);
    setUpdateSuccess(false);

    try {
      await axios.put(`/api/v1/admin/users/${selectedUser.id}`, editFormData);
      setUpdateSuccess(true);
      fetchUsers(); // Обновляем список пользователей

      // Закрываем диалог через 1 секунду после успешного обновления
      setTimeout(() => {
        handleCloseEditDialog();
      }, 1000);
    } catch (err: any) {
      console.error('Ошибка при обновлении пользователя:', err);
      setUpdateError(err.response?.data?.message || err.message || 'Произошла ошибка при обновлении пользователя');
    } finally {
      setUpdateLoading(false);
    }
  };

  // Удаление пользователя
  const handleDeleteUser = async (): Promise<void> => {
    if (!selectedUser) return;

    setUpdateLoading(true);
    setUpdateError(null);

    try {
      await axios.delete(`/api/v1/admin/users/${selectedUser.id}`);
      fetchUsers(); // Обновляем список пользователей
      handleCloseDeleteDialog();
    } catch (err: any) {
      console.error('Ошибка при удалении пользователя:', err);
      setUpdateError(err.response?.data?.message || err.message || 'Произошла ошибка при удалении пользователя');
    } finally {
      setUpdateLoading(false);
    }
  };

  // Блокировка/разблокировка пользователя
  const handleToggleUserStatus = async (user: User): Promise<void> => {
    const newStatus = user.account_status === 'active' ? 'blocked' : 'active';

    try {
      await axios.put(`/api/v1/admin/users/${user.id}/status`, { status: newStatus });
      fetchUsers(); // Обновляем список пользователей
    } catch (err: any) {
      console.error('Ошибка при изменении статуса пользователя:', err);
      setError(err.response?.data?.message || err.message || 'Произошла ошибка при изменении статуса пользователя');
    }
  };

  // Получение цвета статуса
  const getStatusColor = (status?: string): "success" | "error" | "warning" | "default" => {
    switch (status) {
      case 'active':
        return 'success';
      case 'blocked':
        return 'error';
      case 'pending':
        return 'warning';
      default:
        return 'default';
    }
  };

  return (
    <Box sx={{ py: 4 }}>
      <Paper sx={{ p: 3, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" gutterBottom>
            Управление пользователями
          </Typography>
          <Button
            variant="outlined"
            startIcon={<ArrowBackIcon />}
            onClick={() => navigate('/admin')}
          >
            Назад к панели администратора
          </Button>
        </Box>

        <Box sx={{ mb: 3 }}>
          <TextField
            label="Поиск пользователей"
            variant="outlined"
            fullWidth
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Поиск по имени, email или телефону"
            sx={{ mb: 2 }}
          />
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Имя</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Телефон</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Дата регистрации</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : searchResults.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    Пользователи не найдены
                  </TableCell>
                </TableRow>
              ) : (
                searchResults.map((user) => (
                  <TableRow
                    key={user.id}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { backgroundColor: 'rgba(0, 0, 0, 0.04)' }
                    }}
                    onClick={() => handleOpenDetailsDialog(user)}
                  >
                    <TableCell>{user.id}</TableCell>
                    <TableCell>{user.name || '-'}</TableCell>
                    <TableCell>{user.email || '-'}</TableCell>
                    <TableCell>{user.phone || '-'}</TableCell>
                    <TableCell>
                      <Chip
                        label={user.account_status || 'active'}
                        color={getStatusColor(user.account_status)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>{formatDate(user.created_at)}</TableCell>
                    <TableCell onClick={(e) => e.stopPropagation()}>
                      <Tooltip title="Подробнее">
                        <IconButton onClick={() => handleOpenDetailsDialog(user)} size="small" color="primary">
                          <InfoIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Редактировать">
                        <IconButton onClick={() => handleOpenEditDialog(user)} size="small">
                          <EditIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title={user.account_status === 'active' ? 'Заблокировать' : 'Разблокировать'}>
                        <IconButton
                          onClick={() => handleToggleUserStatus(user)}
                          size="small"
                          color={user.account_status === 'active' ? 'default' : 'error'}
                        >
                          {user.account_status === 'active' ? <BlockIcon fontSize="small" /> : <CheckCircleIcon fontSize="small" />}
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Удалить">
                        <IconButton onClick={() => handleOpenDeleteDialog(user)} size="small" color="error">
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        <TablePagination
          component="div"
          count={totalUsers}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[5, 10, 25, 50]}
          labelRowsPerPage="Строк на странице:"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} из ${count}`}
        />
      </Paper>

      {/* Диалог редактирования пользователя */}
      <Dialog open={openEditDialog} onClose={handleCloseEditDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Редактирование пользователя</DialogTitle>
        <DialogContent>
          {updateError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {updateError}
            </Alert>
          )}
          {updateSuccess && (
            <Alert severity="success" sx={{ mb: 2 }}>
              Пользователь успешно обновлен!
            </Alert>
          )}
          <TextField
            margin="dense"
            name="name"
            label="Имя"
            type="text"
            fullWidth
            variant="outlined"
            value={editFormData.name}
            onChange={handleEditFormChange}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            name="email"
            label="Email"
            type="email"
            fullWidth
            variant="outlined"
            value={editFormData.email}
            onChange={handleEditFormChange}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            name="phone"
            label="Телефон"
            type="text"
            fullWidth
            variant="outlined"
            value={editFormData.phone || ''}
            onChange={handleEditFormChange}
            sx={{ mb: 2 }}
          />
          <TextField
            select
            margin="dense"
            name="status"
            label="Статус"
            fullWidth
            variant="outlined"
            value={editFormData.status}
            onChange={handleStatusChange}
            SelectProps={{
              native: true,
            }}
          >
            <option value="active">Активен</option>
            <option value="blocked">Заблокирован</option>
            <option value="pending">Ожидает подтверждения</option>
          </TextField>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseEditDialog}>Отмена</Button>
          <Button
            onClick={handleUpdateUser}
            variant="contained"
            color="primary"
            disabled={updateLoading}
            startIcon={updateLoading && <CircularProgress size={20} color="inherit" />}
          >
            {updateLoading ? 'Сохранение...' : 'Сохранить'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Диалог удаления пользователя */}
      <Dialog open={openDeleteDialog} onClose={handleCloseDeleteDialog}>
        <DialogTitle>Подтверждение удаления</DialogTitle>
        <DialogContent>
          {updateError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {updateError}
            </Alert>
          )}
          <Typography>
            Вы действительно хотите удалить пользователя {selectedUser?.name || selectedUser?.email}?
            Это действие нельзя отменить.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog}>Отмена</Button>
          <Button
            onClick={handleDeleteUser}
            variant="contained"
            color="error"
            disabled={updateLoading}
            startIcon={updateLoading && <CircularProgress size={20} color="inherit" />}
          >
            {updateLoading ? 'Удаление...' : 'Удалить'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Диалог просмотра детальной информации */}
      {selectedUser && (
        <UserDetailsDialog
          open={openDetailsDialog}
          onClose={handleCloseDetailsDialog}
          userId={selectedUser.id}
        />
      )}
    </Box>
  );
};

export default UsersManagementPage;