// Package handler содержит обработчики HTTP запросов для работы с отзывами
package handler

import (
	"backend/internal/domain/models"
)

// Response structures for Swagger documentation

// ReviewResponse представляет ответ с отзывом
type ReviewResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    *models.Review `json:"data"`
}

// ReviewsListResponse представляет ответ со списком отзывов
type ReviewsListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []models.Review `json:"data"`
	Meta    ReviewsMeta     `json:"meta"`
}

// ReviewsMeta метаданные для списка отзывов
type ReviewsMeta struct {
	Total int `json:"total" example:"100"`
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
}

// ReviewMessageResponse представляет ответ с сообщением
type ReviewMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// VoteRequest запрос на голосование за отзыв
type VoteRequest struct {
	VoteType string `json:"vote_type" example:"helpful"`
}

// ResponseRequest запрос на добавление ответа на отзыв
type ResponseRequest struct {
	Response string `json:"response" example:"Спасибо за ваш отзыв!"`
}

// PhotosResponse ответ с загруженными фотографиями
type PhotosResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Photos uploaded successfully"`
	Photos  []string `json:"photos"`
}

// RatingResponse ответ с рейтингом
type RatingResponse struct {
	Success bool    `json:"success" example:"true"`
	Rating  float64 `json:"rating" example:"4.5"`
}

// StatsResponse ответ со статистикой отзывов
type StatsResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    *models.ReviewStats `json:"data"`
}

// UserRatingSummaryResponse ответ с рейтингом пользователя
type UserRatingSummaryResponse struct {
	Success bool                      `json:"success" example:"true"`
	Data    *models.UserRatingSummary `json:"data"`
}

// StorefrontRatingSummaryResponse ответ с рейтингом витрины
type StorefrontRatingSummaryResponse struct {
	Success bool                           `json:"success" example:"true"`
	Data    *models.StorefrontRatingSummary `json:"data"`
}