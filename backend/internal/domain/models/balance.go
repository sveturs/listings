// backend/internal/domain/models/balance.go

package models

import (
    "time"
    "encoding/json"
)

type UserBalance struct {
    UserID        int       `json:"user_id"`
    Balance       float64   `json:"balance"`
    FrozenBalance float64   `json:"frozen_balance"`
    Currency      string    `json:"currency"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type BalanceTransaction struct {
    ID             int             `json:"id"`
    UserID         int             `json:"user_id"`
    Type           string          `json:"type"`
    Amount         float64         `json:"amount"`
    Currency       string          `json:"currency"`
    Status         string          `json:"status"`
    PaymentMethod  string          `json:"payment_method"`
    PaymentDetails json.RawMessage `json:"payment_details"`
    Description    string          `json:"description"`
    CreatedAt      time.Time       `json:"created_at"`
    CompletedAt    *time.Time      `json:"completed_at,omitempty"`
}

type PaymentMethod struct {
    ID            int             `json:"id"`
    Name          string          `json:"name"`
    Code          string          `json:"code"`
    Type          string          `json:"type"`
    IsActive      bool            `json:"is_active"`
    MinimumAmount float64         `json:"minimum_amount"`
    MaximumAmount float64         `json:"maximum_amount"`
    FeePercentage float64         `json:"fee_percentage"`
    FixedFee      float64         `json:"fixed_fee"`
    Credentials   json.RawMessage `json:"credentials,omitempty"`
    CreatedAt     time.Time       `json:"created_at"`
}