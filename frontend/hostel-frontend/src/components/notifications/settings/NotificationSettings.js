// src/components/notifications/settings/NotificationSettings.js
import { Box, Typography, Switch, FormControl, FormControlLabel } from '@mui/material';
import { useNotifications } from '../../../hooks/useNotifications';

export default function NotificationSettings() {
  const { settings, updateSettings } = useNotifications();

  const notificationTypes = [
    { key: 'new_message', label: 'Новые сообщения' },
    { key: 'new_review', label: 'Новые отзывы' },
    // ...другие типы
  ];

  return (
    <Box p={3}>
      <Typography variant="h5" mb={3}>Настройки уведомлений</Typography>
      
      {notificationTypes.map(type => (
        <FormControl key={type.key} fullWidth margin="normal">
          <Typography variant="h6">{type.label}</Typography>
          <FormControlLabel
            control={
              <Switch 
                checked={settings[type.key]?.telegram || false}
                onChange={(e) => updateSettings(type.key, 'telegram', e.target.checked)}
              />
            }
            label="Telegram"
          />
          <FormControlLabel
            control={
              <Switch 
                checked={settings[type.key]?.push || false}
                onChange={(e) => updateSettings(type.key, 'push', e.target.checked)}
              />
            }
            label="Push-уведомления"
          />
        </FormControl>
      ))}
    </Box>
  );
}