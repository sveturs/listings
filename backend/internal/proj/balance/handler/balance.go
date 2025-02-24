// backend/internal/proj/balance/handler/balance.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
    balance "backend/internal/proj/balance/service"  
	"backend/internal/domain/models" 
    "strconv"
	"log"
	"strings"  
)

type BalanceHandler struct {
    balanceService balance.BalanceServiceInterface // Используем balance вместо service
}

func NewBalanceHandler(balanceService balance.BalanceServiceInterface) *BalanceHandler { // Используем balance вместо service
    return &BalanceHandler{
        balanceService: balanceService,
    }
}
 
func (h *BalanceHandler) GetBalance(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    log.Printf("Getting balance for user %d", userID)

    balance, err := h.balanceService.GetBalance(c.Context(), userID)
    if err != nil {
        log.Printf("Error getting balance for user %d: %v", userID, err)
        // Если записи нет, возвращаем нулевой баланс
        if err.Error() == "no rows in result set" {
            return utils.SuccessResponse(c, &models.UserBalance{
                UserID:    userID,
                Balance:   0,
                Currency: "RSD",
            })
        }
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get balance")
    }

    return utils.SuccessResponse(c, balance)
}


func (h *BalanceHandler) CreateDeposit(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)

    var request struct {
        Amount        float64 `json:"amount"`
        PaymentMethod string  `json:"payment_method"`
    }

    if err := c.BodyParser(&request); err != nil {
        log.Printf("Error parsing deposit request: %v", err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
    }

    log.Printf("Processing deposit request: amount=%f, method=%s", request.Amount, request.PaymentMethod)

    transaction, err := h.balanceService.CreateDeposit(c.Context(), userID, request.Amount, request.PaymentMethod)
    if err != nil {
        log.Printf("Error creating deposit: %v", err)
        
        switch {
        case strings.Contains(err.Error(), "below minimum allowed"):
            return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
        case strings.Contains(err.Error(), "above maximum allowed"):
            return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
        default:
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обработки депозита")
        }
    }

    log.Printf("Successfully created deposit transaction: %+v", transaction)
    return utils.SuccessResponse(c, transaction)
}


func (h *BalanceHandler) GetTransactions(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    // Заменяем utils.QueryInt на собственную реализацию
    limit := 20
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    offset := 0
    if offsetStr := c.Query("offset"); offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }

    transactions, err := h.balanceService.GetTransactions(c.Context(), userID, limit, offset)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get transactions")
    }

    return utils.SuccessResponse(c, transactions)
}

func (h *BalanceHandler) GetPaymentMethods(c *fiber.Ctx) error {
    methods, err := h.balanceService.GetPaymentMethods(c.Context())
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get payment methods")
    }

    return utils.SuccessResponse(c, methods)
}