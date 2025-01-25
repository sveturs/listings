// backend/internal/proj/notifications/handler/handler.go
package handler

import (
    globalService "backend/internal/proj/global/service"
)

type Handler struct {
    Notification *NotificationHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
    return &Handler{
        Notification: NewNotificationHandler(services.Notification()),
    }
}