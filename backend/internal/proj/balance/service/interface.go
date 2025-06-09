// backend/internal/proj/balance/service/interface.go

package balance

import (
	"backend/internal/domain/models"
	"context"
)

type BalanceServiceInterface interface {
	GetBalance(ctx context.Context, userID int) (*models.UserBalance, error)
	CreateDeposit(ctx context.Context, userID int, amount float64, method string) (*models.BalanceTransaction, error)
	ProcessDeposit(ctx context.Context, transactionID int) error
	GetTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error)
	GetPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error)
}
