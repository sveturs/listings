import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
  Snackbar,
  Alert,
  IconButton,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import { useTranslation } from 'react-i18next';

interface AdminUser {
  id: number;
  email: string;
  created_at: string;
  created_by?: number;
  notes?: string;
}

const AdminManagerPage: React.FC = () => {
  const { t } = useTranslation();
  const [admins, setAdmins] = useState<AdminUser[]>([]);
  const [newAdminEmail, setNewAdminEmail] = useState('');
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [deleteDialog, setDeleteDialog] = useState(false);
  const [adminToDelete, setAdminToDelete] = useState<string | null>(null);
  const { user } = useAuth();

  // Загрузка данных админов
  const fetchAdmins = async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/v1/admin/admins');
      // Убедимся, что response.data - это массив, иначе используем пустой массив
      setAdmins(Array.isArray(response.data) ? response.data : []);
      setError(null);
    } catch (err) {
      console.error('Ошибка при загрузке списка администраторов:', err);
      setError('Не удалось загрузить список администраторов');
      setAdmins([]); // В случае ошибки установим пустой массив
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAdmins();
  }, []);

  // Добавление нового админа
  const handleAddAdmin = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newAdminEmail) {
      setError('Email администратора не может быть пустым');
      return;
    }

    setLoading(true);
    try {
      await axios.post('/api/v1/admin/admins', {
        email: newAdminEmail,
        notes: notes || undefined
      });
      
      setSuccess(`Администратор ${newAdminEmail} успешно добавлен`);
      setNewAdminEmail('');
      setNotes('');
      fetchAdmins();
    } catch (err) {
      console.error('Ошибка при добавлении администратора:', err);
      setError('Не удалось добавить администратора');
    } finally {
      setLoading(false);
    }
  };

  // Удаление админа
  const handleDeleteAdmin = async () => {
    if (!adminToDelete) return;
    
    setLoading(true);
    try {
      await axios.delete(`/api/v1/admin/admins/${adminToDelete}`);
      setSuccess(`Администратор ${adminToDelete} успешно удален`);
      fetchAdmins();
    } catch (err) {
      console.error('Ошибка при удалении администратора:', err);
      setError('Не удалось удалить администратора');
    } finally {
      setLoading(false);
      setDeleteDialog(false);
      setAdminToDelete(null);
    }
  };

  // Открытие диалога подтверждения удаления
  const openDeleteDialog = (email: string) => {
    setAdminToDelete(email);
    setDeleteDialog(true);
  };

  // Форматирование даты
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  // Проверка, является ли пользователь текущим
  const isCurrentUser = (email: string) => {
    return user?.email === email;
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Управление администраторами
      </Typography>

      {/* Форма добавления админа */}
      <Paper sx={{ p: 3, mb: 4 }}>
        <Typography variant="h6" gutterBottom>
          Добавить нового администратора
        </Typography>
        <Box component="form" onSubmit={handleAddAdmin} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label="Email администратора"
            value={newAdminEmail}
            onChange={(e) => setNewAdminEmail(e.target.value)}
            required
            fullWidth
            type="email"
          />
          <TextField
            label="Примечание (необязательно)"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            fullWidth
            multiline
            rows={2}
          />
          <Button 
            type="submit" 
            variant="contained" 
            color="primary" 
            disabled={loading}
            sx={{ alignSelf: 'flex-start' }}
          >
            Добавить администратора
          </Button>
        </Box>
      </Paper>

      {/* Таблица администраторов */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Список администраторов
        </Typography>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Дата добавления</TableCell>
                <TableCell>Примечание</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {admins.map((admin) => (
                <TableRow key={admin.id}>
                  <TableCell>{admin.id}</TableCell>
                  <TableCell>{admin.email}</TableCell>
                  <TableCell>{formatDate(admin.created_at)}</TableCell>
                  <TableCell>{admin.notes || '-'}</TableCell>
                  <TableCell>
                    <IconButton 
                      color="error" 
                      onClick={() => openDeleteDialog(admin.email)}
                      disabled={isCurrentUser(admin.email)} // Запрещаем удалять себя
                      title={isCurrentUser(admin.email) ? "Нельзя удалить себя" : "Удалить администратора"}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
              {admins.length === 0 && (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    {loading ? 'Загрузка...' : 'Нет администраторов'}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Диалог подтверждения удаления */}
      <Dialog
        open={deleteDialog}
        onClose={() => setDeleteDialog(false)}
      >
        <DialogTitle>Подтвердите действие</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Вы действительно хотите удалить администратора с email: {adminToDelete}?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialog(false)} color="primary">
            Отмена
          </Button>
          <Button onClick={handleDeleteAdmin} color="error" autoFocus>
            Удалить
          </Button>
        </DialogActions>
      </Dialog>

      {/* Уведомления */}
      <Snackbar open={!!error} autoHideDuration={6000} onClose={() => setError(null)}>
        <Alert onClose={() => setError(null)} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>

      <Snackbar open={!!success} autoHideDuration={6000} onClose={() => setSuccess(null)}>
        <Alert onClose={() => setSuccess(null)} severity="success" sx={{ width: '100%' }}>
          {success}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default AdminManagerPage;