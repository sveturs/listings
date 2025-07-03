# 🔔 Паспорт процесса: Система уведомлений

## 📋 Метаданные
- **Код**: BP-015
- **Название**: Multi-channel Notification System
- **Версия**: 1.0.0
- **Статус**: Active
- **Критичность**: High
- **Владелец**: Platform Team

## 🎯 Краткое описание

Мультиканальная система доставки уведомлений пользователям платформы Sve Tu через in-app, email и Telegram каналы. Система обеспечивает своевременное информирование пользователей о важных событиях, поддерживает персонализированные настройки и гарантирует доставку критических уведомлений.

## 📊 Диаграмма процесса

```mermaid
flowchart TD
    Start([Событие в системе]) --> DetermineType{Определение типа<br/>уведомления}
    
    DetermineType --> CreateNotification[Создание уведомления<br/>в сервисе]
    
    CreateNotification --> CheckUser{Проверка<br/>пользователя}
    CheckUser -->|Не найден| LogError[Логирование<br/>ошибки]
    CheckUser -->|OK| LoadSettings[Загрузка настроек<br/>пользователя]
    
    LoadSettings --> CheckSettings{Проверка настроек<br/>для типа}
    CheckSettings -->|Нет настроек| UseDefaults[Использование<br/>дефолтных настроек]
    CheckSettings -->|OK| PrepareData[Подготовка данных<br/>уведомления]
    UseDefaults --> PrepareData
    
    PrepareData --> SaveInApp[Сохранение в БД<br/>для in-app]
    SaveInApp --> DeliveryChannels{Каналы<br/>доставки}
    
    %% Email канал
    DeliveryChannels -->|Email включен| PrepareEmail[Подготовка<br/>email шаблона]
    PrepareEmail --> SendEmail[Отправка через<br/>SMTP сервер]
    SendEmail --> EmailResult{Результат<br/>отправки}
    EmailResult -->|Успех| LogEmailSuccess[Лог успеха]
    EmailResult -->|Ошибка| LogEmailError[Лог ошибки<br/>+ retry queue]
    
    %% Telegram канал
    DeliveryChannels -->|Telegram включен| CheckTelegram{Подключен<br/>Telegram?}
    CheckTelegram -->|Нет| SkipTelegram[Пропустить<br/>Telegram]
    CheckTelegram -->|Да| FormatTelegram[Форматирование<br/>для Telegram]
    FormatTelegram --> SendTelegram[Отправка через<br/>Bot API]
    SendTelegram --> TelegramResult{Результат<br/>отправки}
    TelegramResult -->|Успех| LogTelegramSuccess[Лог успеха]
    TelegramResult -->|Ошибка| LogTelegramError[Лог ошибки]
    
    %% Завершение
    LogEmailSuccess --> UpdateDeliveryStatus[Обновление статуса<br/>доставки]
    LogEmailError --> UpdateDeliveryStatus
    LogTelegramSuccess --> UpdateDeliveryStatus
    LogTelegramError --> UpdateDeliveryStatus
    SkipTelegram --> UpdateDeliveryStatus
    
    UpdateDeliveryStatus --> Complete([Завершение])
    
    %% Чтение уведомлений
    UserOpens([Пользователь открывает<br/>уведомления]) --> FetchUnread[GET /api/v1/notifications]
    FetchUnread --> DisplayList[Отображение списка<br/>(Frontend TODO)]
    DisplayList --> UserReads{Пользователь<br/>читает}
    UserReads -->|Да| MarkAsRead[PUT /api/v1/notifications/:id/read]
    MarkAsRead --> UpdateReadStatus[(Обновление<br/>is_read = true)]
    
    %% Управление настройками
    UserSettings([Настройки<br/>уведомлений]) --> LoadCurrentSettings[GET /api/v1/notifications/settings]
    LoadCurrentSettings --> ShowSettings[Отображение настроек<br/>(Frontend TODO)]
    ShowSettings --> UserChanges{Изменение<br/>настроек}
    UserChanges -->|Да| SaveSettings[PUT /api/v1/notifications/settings]
    SaveSettings --> ValidateSettings{Валидация}
    ValidateSettings -->|OK| UpdateDB[(Обновление<br/>notification_settings)]
    ValidateSettings -->|Ошибка| ShowValidationError[Показ ошибки]
    
    %% Подключение Telegram
    ConnectTelegram([Подключить<br/>Telegram]) --> RequestToken[GET /api/v1/notifications/telegram/token]
    RequestToken --> GenerateToken[Генерация токена<br/>с HMAC подписью]
    GenerateToken --> ShowBotLink[Показ ссылки на бота<br/>t.me/svetubot?start=TOKEN]
    
    UserInTelegram([Пользователь в<br/>Telegram]) --> StartCommand[Команда /start TOKEN]
    StartCommand --> WebhookReceive[POST /api/v1/notifications/telegram/webhook]
    WebhookReceive --> ValidateToken{Валидация<br/>токена}
    ValidateToken -->|Invalid| SendErrorMessage[Отправка ошибки<br/>в Telegram]
    ValidateToken -->|Valid| LinkAccount[Связывание аккаунтов]
    LinkAccount --> SaveConnection[(Сохранение в<br/>user_telegram_connections)]
    SaveConnection --> SendWelcome[Отправка приветствия<br/>в Telegram]
```

## 🔄 Детальный Flow

### 1️⃣ **Отправка уведомления из бизнес-логики**

```go
// Backend: Пример отправки уведомления о новом сообщении
func (s *ChatService) SendMessage(senderID, recipientID int64, message string) error {
    // Сохранение сообщения
    msg, err := s.storage.CreateMessage(senderID, recipientID, message)
    if err != nil {
        return err
    }
    
    // Отправка уведомления получателю
    notification := &NotificationData{
        UserID: recipientID,
        Type:   "new_message",
        Title:  "Новое сообщение",
        Message: fmt.Sprintf("У вас новое сообщение от %s", senderName),
        Data: map[string]interface{}{
            "chat_id":    msg.ChatID,
            "message_id": msg.ID,
            "sender_id":  senderID,
            "preview":    truncateMessage(message, 100),
        },
    }
    
    // Асинхронная отправка
    go s.notificationService.SendNotification(notification)
    
    return nil
}

// Service: Централизованная отправка уведомлений
func (s *NotificationService) SendNotification(data *NotificationData) error {
    // 1. Проверка пользователя
    user, err := s.userService.GetUser(data.UserID)
    if err != nil {
        log.Printf("User %d not found for notification: %v", data.UserID, err)
        return err
    }
    
    // 2. Загрузка настроек
    settings, err := s.storage.GetUserSettings(data.UserID, data.Type)
    if err != nil {
        // Используем дефолтные настройки
        settings = s.getDefaultSettings(data.Type)
    }
    
    // 3. Создание записи в БД (in-app всегда включен)
    notification := &Notification{
        UserID:    data.UserID,
        Type:      data.Type,
        Title:     data.Title,
        Message:   data.Message,
        Data:      data.Data,
        IsRead:    false,
        CreatedAt: time.Now(),
    }
    
    if err := s.storage.CreateNotification(notification); err != nil {
        log.Printf("Failed to save notification: %v", err)
        return err
    }
    
    // 4. Отправка по каналам
    deliveryStatus := make(map[string]bool)
    
    // Email
    if settings.EmailEnabled && user.Email != "" {
        if err := s.sendEmail(user, notification); err != nil {
            log.Printf("Email delivery failed for user %d: %v", user.ID, err)
            deliveryStatus["email"] = false
        } else {
            deliveryStatus["email"] = true
        }
    }
    
    // Telegram
    if settings.TelegramEnabled {
        if err := s.sendTelegram(user.ID, notification); err != nil {
            log.Printf("Telegram delivery failed for user %d: %v", user.ID, err)
            deliveryStatus["telegram"] = false
        } else {
            deliveryStatus["telegram"] = true
        }
    }
    
    // 5. Обновление статуса доставки
    notification.DeliveredTo = deliveryStatus
    s.storage.UpdateNotificationDelivery(notification.ID, deliveryStatus)
    
    return nil
}
```

### 2️⃣ **Email доставка**

```go
// Backend: Email отправка с шаблонами
func (s *NotificationService) sendEmail(user *User, notification *Notification) error {
    // Подготовка данных для шаблона
    templateData := map[string]interface{}{
        "UserName":     user.Name,
        "Title":        notification.Title,
        "Message":      notification.Message,
        "ActionURL":    s.buildActionURL(notification),
        "UnsubscribeURL": s.buildUnsubscribeURL(user.ID, notification.Type),
        "Year":         time.Now().Year(),
    }
    
    // Выбор шаблона по типу
    var templateName string
    switch notification.Type {
    case "new_message":
        templateName = "new_message.html"
    case "listing_status":
        templateName = "listing_status.html"
    case "new_review":
        templateName = "new_review.html"
    default:
        templateName = "default.html"
    }
    
    // Рендеринг HTML шаблона
    var htmlBody bytes.Buffer
    tmpl, err := template.ParseFiles(fmt.Sprintf("templates/email/%s", templateName))
    if err != nil {
        return fmt.Errorf("template parse error: %w", err)
    }
    
    if err := tmpl.Execute(&htmlBody, templateData); err != nil {
        return fmt.Errorf("template execute error: %w", err)
    }
    
    // Создание email сообщения
    m := gomail.NewMessage()
    m.SetHeader("From", s.config.EmailFrom)
    m.SetHeader("To", user.Email)
    m.SetHeader("Subject", notification.Title)
    m.SetBody("text/html", htmlBody.String())
    
    // Добавление заголовков для отписки
    m.SetHeader("List-Unsubscribe", fmt.Sprintf("<%s>", templateData["UnsubscribeURL"]))
    m.SetHeader("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")
    
    // Отправка через SMTP
    d := gomail.NewDialer(
        s.config.SMTPHost,
        s.config.SMTPPort,
        s.config.SMTPUser,
        s.config.SMTPPassword,
    )
    
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("smtp send error: %w", err)
    }
    
    return nil
}
```

### 3️⃣ **Telegram интеграция**

```go
// Backend: Подключение Telegram аккаунта
func (h *NotificationHandler) GetTelegramToken(c *fiber.Ctx) error {
    userID := c.Locals("userID").(int64)
    
    // Генерация токена с подписью
    token := fmt.Sprintf("%d:%d", userID, time.Now().Unix())
    signature := h.generateHMAC(token)
    fullToken := fmt.Sprintf("%s:%s", token, signature)
    
    // Сохранение токена в кеше на 15 минут
    h.cache.Set(fmt.Sprintf("tg_token:%s", fullToken), userID, 15*time.Minute)
    
    return utils.SuccessResponse(c, map[string]interface{}{
        "token":    fullToken,
        "bot_link": fmt.Sprintf("https://t.me/%s?start=%s", h.config.TelegramBotUsername, fullToken),
        "expires_in": 900, // 15 минут
    })
}

// Webhook обработка команды /start
func (h *NotificationHandler) TelegramWebhook(c *fiber.Ctx) error {
    var update TelegramUpdate
    if err := c.BodyParser(&update); err != nil {
        return c.SendStatus(fiber.StatusOK) // Telegram требует 200 OK
    }
    
    // Обработка команды /start
    if update.Message != nil && strings.HasPrefix(update.Message.Text, "/start ") {
        token := strings.TrimPrefix(update.Message.Text, "/start ")
        
        // Валидация токена
        userID, ok := h.validateTelegramToken(token)
        if !ok {
            h.sendTelegramMessage(update.Message.Chat.ID, 
                "❌ Неверный или истекший токен. Пожалуйста, получите новый токен в настройках профиля.")
            return c.SendStatus(fiber.StatusOK)
        }
        
        // Связывание аккаунтов
        connection := &UserTelegramConnection{
            UserID:           userID,
            TelegramChatID:   update.Message.Chat.ID,
            TelegramUsername: update.Message.From.Username,
            ConnectedAt:      time.Now(),
        }
        
        if err := h.storage.SaveTelegramConnection(connection); err != nil {
            h.sendTelegramMessage(update.Message.Chat.ID, 
                "❌ Ошибка при подключении. Попробуйте позже.")
            return c.SendStatus(fiber.StatusOK)
        }
        
        // Приветственное сообщение
        h.sendTelegramMessage(update.Message.Chat.ID, 
            "✅ Telegram успешно подключен!\n\n" +
            "Теперь вы будете получать уведомления от Sve Tu в этот чат.\n" +
            "Управлять настройками можно в профиле на сайте.")
    }
    
    return c.SendStatus(fiber.StatusOK)
}

// Отправка уведомления в Telegram
func (s *NotificationService) sendTelegram(userID int64, notification *Notification) error {
    // Получение chat_id
    connection, err := s.storage.GetTelegramConnection(userID)
    if err != nil {
        return fmt.Errorf("telegram not connected")
    }
    
    // Форматирование сообщения
    text := s.formatTelegramMessage(notification)
    
    // Создание inline клавиатуры
    keyboard := s.buildTelegramKeyboard(notification)
    
    // Отправка через Bot API
    bot, err := tgbotapi.NewBotAPI(s.config.TelegramBotToken)
    if err != nil {
        return err
    }
    
    msg := tgbotapi.NewMessage(connection.TelegramChatID, text)
    msg.ParseMode = "HTML"
    msg.DisableWebPagePreview = true
    
    if keyboard != nil {
        msg.ReplyMarkup = keyboard
    }
    
    _, err = bot.Send(msg)
    return err
}

// Форматирование для Telegram
func (s *NotificationService) formatTelegramMessage(n *Notification) string {
    var emoji string
    switch n.Type {
    case "new_message":
        emoji = "💬"
    case "listing_status":
        emoji = "📦"
    case "new_review":
        emoji = "⭐"
    case "payment_received":
        emoji = "💰"
    default:
        emoji = "🔔"
    }
    
    // HTML форматирование
    text := fmt.Sprintf("%s <b>%s</b>\n\n%s", emoji, n.Title, n.Message)
    
    // Добавление дополнительной информации
    if data, ok := n.Data.(map[string]interface{}); ok {
        if preview, ok := data["preview"].(string); ok {
            text += fmt.Sprintf("\n\n<i>%s</i>", html.EscapeString(preview))
        }
    }
    
    return text
}
```

### 4️⃣ **Управление настройками**

```typescript
// Frontend: Компонент настроек (TODO - не реализован)
const NotificationSettings: React.FC = () => {
  const [settings, setSettings] = useState<NotificationSettings[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Загрузка текущих настроек
  useEffect(() => {
    fetchSettings();
  }, []);
  
  const fetchSettings = async () => {
    try {
      const response = await api.get('/api/v1/notifications/settings');
      setSettings(response.data.data);
    } catch (error) {
      toast.error(t('settings.loadError'));
    } finally {
      setLoading(false);
    }
  };
  
  // Обновление настроек
  const handleToggle = async (type: string, channel: 'email' | 'telegram') => {
    const setting = settings.find(s => s.notification_type === type);
    if (!setting) return;
    
    const updated = {
      ...setting,
      [`${channel}_enabled`]: !setting[`${channel}_enabled`],
    };
    
    try {
      await api.put('/api/v1/notifications/settings', {
        settings: [updated],
      });
      
      // Оптимистичное обновление
      setSettings(prev => 
        prev.map(s => s.notification_type === type ? updated : s)
      );
      
      toast.success(t('settings.updated'));
    } catch (error) {
      toast.error(t('settings.updateError'));
      // Откат изменений
      fetchSettings();
    }
  };
  
  return (
    <div className="card">
      <div className="card-body">
        <h2 className="card-title">{t('notifications.settings.title')}</h2>
        
        {/* Telegram подключение */}
        <TelegramConnection />
        
        {/* Настройки по типам */}
        <div className="space-y-4 mt-6">
          {settings.map(setting => (
            <div key={setting.notification_type} className="border rounded-lg p-4">
              <h3 className="font-medium mb-2">
                {t(`notifications.types.${setting.notification_type}`)}
              </h3>
              
              <div className="flex gap-4">
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={setting.email_enabled}
                    onChange={() => handleToggle(setting.notification_type, 'email')}
                    className="checkbox"
                  />
                  <span>{t('notifications.channels.email')}</span>
                </label>
                
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={setting.telegram_enabled}
                    onChange={() => handleToggle(setting.notification_type, 'telegram')}
                    className="checkbox"
                  />
                  <span>{t('notifications.channels.telegram')}</span>
                </label>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
```

### 5️⃣ **In-app уведомления и real-time**

```typescript
// Frontend: Компонент уведомлений (TODO - не реализован)
const NotificationBell: React.FC = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [isOpen, setIsOpen] = useState(false);
  
  // WebSocket подключение для real-time
  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/notifications`);
    
    ws.onmessage = (event) => {
      const notification = JSON.parse(event.data);
      
      // Добавление нового уведомления
      setNotifications(prev => [notification, ...prev]);
      setUnreadCount(prev => prev + 1);
      
      // Показ browser notification
      if (Notification.permission === 'granted') {
        new Notification(notification.title, {
          body: notification.message,
          icon: '/logo.png',
        });
      }
    };
    
    return () => ws.close();
  }, []);
  
  // Загрузка уведомлений
  const fetchNotifications = async () => {
    try {
      const response = await api.get('/api/v1/notifications', {
        params: { limit: 20, unread_only: false },
      });
      
      setNotifications(response.data.data);
      setUnreadCount(response.data.unread_count);
    } catch (error) {
      console.error('Failed to fetch notifications:', error);
    }
  };
  
  // Отметка как прочитанное
  const markAsRead = async (id: number) => {
    try {
      await api.put(`/api/v1/notifications/${id}/read`);
      
      setNotifications(prev =>
        prev.map(n => n.id === id ? { ...n, is_read: true } : n)
      );
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (error) {
      console.error('Failed to mark as read:', error);
    }
  };
  
  return (
    <div className="relative">
      <button
        className="btn btn-ghost btn-circle"
        onClick={() => setIsOpen(!isOpen)}
      >
        <Bell className="w-5 h-5" />
        {unreadCount > 0 && (
          <span className="badge badge-sm badge-error absolute -top-1 -right-1">
            {unreadCount > 99 ? '99+' : unreadCount}
          </span>
        )}
      </button>
      
      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 card bg-base-100 shadow-xl">
          <div className="card-body p-0">
            <div className="p-4 border-b">
              <h3 className="font-bold">{t('notifications.title')}</h3>
            </div>
            
            <div className="max-h-96 overflow-y-auto">
              {notifications.length === 0 ? (
                <div className="p-4 text-center text-gray-500">
                  {t('notifications.empty')}
                </div>
              ) : (
                notifications.map(notification => (
                  <NotificationItem
                    key={notification.id}
                    notification={notification}
                    onRead={markAsRead}
                  />
                ))
              )}
            </div>
            
            <div className="p-2 border-t">
              <Link href="/notifications" className="btn btn-sm btn-block">
                {t('notifications.viewAll')}
              </Link>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
```

## 🔐 Безопасность и валидация

### Безопасность каналов
- ✅ HMAC подпись для Telegram токенов
- ✅ Валидация webhook запросов от Telegram
- ✅ SMTP аутентификация для email
- ✅ Rate limiting для предотвращения спама
- ✅ Изоляция данных между пользователями

### Приватность
- ✅ One-click отписка в email
- ✅ Гранулярные настройки по типам
- ✅ Логирование только агрегированных метрик
- ✅ Шифрование sensitive данных

### Валидация
- ✅ Проверка формата email
- ✅ Валидация Telegram chat_id
- ✅ Ограничение размера сообщений
- ✅ Санитизация HTML в сообщениях

## 📊 Аналитика и метрики

### Отслеживаемые события
```typescript
// Доставка уведомлений
analytics.track('notification_sent', {
  user_id: userId,
  type: notificationType,
  channels: ['email', 'telegram'],
  success: true,
});

// Взаимодействие с уведомлениями
analytics.track('notification_read', {
  user_id: userId,
  notification_id: notificationId,
  time_to_read: timeToRead,
});

// Управление настройками
analytics.track('notification_settings_updated', {
  user_id: userId,
  changes: changedSettings,
});

// Подключение каналов
analytics.track('telegram_connected', {
  user_id: userId,
  connection_method: 'bot_command',
});
```

### KPI метрики
- **Delivery Rate**: % успешно доставленных уведомлений
- **Read Rate**: % прочитанных in-app уведомлений
- **Channel Preference**: распределение по каналам
- **Response Time**: время реакции на уведомление
- **Opt-out Rate**: % отписок по типам

## 🧪 Тестирование

### Unit тесты
```go
// Backend: notification_service_test.go
func TestNotificationDelivery(t *testing.T) {
    service := NewNotificationService(mockConfig)
    
    // Тест отправки по всем каналам
    notification := &NotificationData{
        UserID:  1,
        Type:    "new_message",
        Title:   "Test",
        Message: "Test message",
    }
    
    // Mock настройки - все каналы включены
    mockStorage.On("GetUserSettings", 1, "new_message").Return(&NotificationSettings{
        EmailEnabled:    true,
        TelegramEnabled: true,
    }, nil)
    
    err := service.SendNotification(notification)
    assert.NoError(t, err)
    
    // Проверка вызовов
    mockStorage.AssertCalled(t, "CreateNotification", mock.Anything)
    mockEmailService.AssertCalled(t, "Send", mock.Anything)
    mockTelegramService.AssertCalled(t, "Send", mock.Anything)
}

func TestTelegramTokenValidation(t *testing.T) {
    handler := NewNotificationHandler(mockConfig)
    
    // Генерация валидного токена
    token := handler.generateToken(123)
    
    // Проверка валидации
    userID, valid := handler.validateToken(token)
    assert.True(t, valid)
    assert.Equal(t, int64(123), userID)
    
    // Проверка невалидного токена
    _, valid = handler.validateToken("invalid")
    assert.False(t, valid)
}
```

### Integration тесты
```typescript
// Frontend: NotificationBell.test.tsx (когда будет реализован)
describe('NotificationBell', () => {
  it('should display unread count', async () => {
    mockAPI.get.mockResolvedValue({
      data: {
        data: mockNotifications,
        unread_count: 3,
      },
    });
    
    const { getByText } = render(<NotificationBell />);
    
    await waitFor(() => {
      expect(getByText('3')).toBeInTheDocument();
    });
  });
  
  it('should mark notification as read', async () => {
    const { getByTestId } = render(<NotificationBell />);
    
    // Клик на уведомление
    fireEvent.click(getByTestId('notification-1'));
    
    expect(mockAPI.put).toHaveBeenCalledWith('/api/v1/notifications/1/read');
  });
});
```

## ⚡ Производительность и оптимизации

### Backend оптимизации
- 🚀 Асинхронная отправка уведомлений
- 🚀 Batch отправка для массовых рассылок
- 🚀 Кеширование настроек пользователей
- 🚀 Connection pooling для SMTP
- 🚀 Rate limiting для защиты от спама

### Frontend оптимизации (планируемые)
- 🚀 WebSocket для real-time обновлений
- 🚀 Service Worker для offline уведомлений
- 🚀 Виртуализация списка уведомлений
- 🚀 Lazy loading старых уведомлений
- 🚀 IndexedDB для локального кеша

### Рекомендации по масштабированию
- 📈 Очередь сообщений (RabbitMQ/Kafka)
- 📈 Отдельный микросервис для email
- 📈 Horizontal scaling для Telegram ботов
- 📈 CDN для email шаблонов
- 📈 Распределенный кеш для настроек

## 🐛 Известные проблемы и ограничения

1. **Frontend**: Полностью отсутствуют UI компоненты
2. **Real-time**: Нет WebSocket поддержки
3. **Push**: Не реализованы push-уведомления
4. **Retry**: Отсутствует механизм повторной отправки
5. **Метрики**: Нет dashboards для мониторинга доставки

## 🔄 Связанные процессы

- **[BP-005] Коммуникация** - уведомления о новых сообщениях
- **[BP-006] Процесс покупки** - уведомления о платежах
- **[BP-014] Отзывы** - уведомления о новых отзывах
- **[BP-003] Публикация объявлений** - уведомления о статусах

## 📚 Дополнительные ресурсы

- [API документация Notifications](/docs/api/notifications)
- [Telegram Bot настройка](/docs/telegram-bot-setup)
- [Email шаблоны](/templates/email/)
- [Push notifications план](/docs/push-notifications-roadmap)