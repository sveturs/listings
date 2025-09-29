package credit

import (
	"math"
	"net/http"

	"backend/internal/middleware"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler handles credit calculation endpoints
type Handler struct{}

// NewHandler creates a new credit handler
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers credit routes
func (h *Handler) RegisterRoutes(app *fiber.App, _ *middleware.Middleware) error {
	credit := app.Group("/api/v1/credit")
	credit.Post("/calculate", h.CalculateCredit)
	return nil
}

// GetPrefix returns the API prefix for credit endpoints
func (h *Handler) GetPrefix() string {
	return "/api/v1/credit"
}

// CreditCalculationRequest represents the request for credit calculation
type CreditCalculationRequest struct {
	Price        float64 `json:"price" validate:"required,gt=0"`
	DownPayment  float64 `json:"down_payment" validate:"required,gte=0"`
	Term         int     `json:"term" validate:"required,gt=0"`           // months
	InterestRate float64 `json:"interest_rate" validate:"required,gte=0"` // annual percentage
	Category     string  `json:"category"`
}

// CreditCalculationResponse represents the credit calculation result
type CreditCalculationResponse struct {
	MonthlyPayment  float64               `json:"monthly_payment"`
	TotalAmount     float64               `json:"total_amount"`
	TotalInterest   float64               `json:"total_interest"`
	LoanAmount      float64               `json:"loan_amount"`
	DownPayment     float64               `json:"down_payment"`
	Term            int                   `json:"term"`
	InterestRate    float64               `json:"interest_rate"`
	PaymentSchedule []PaymentScheduleItem `json:"payment_schedule,omitempty"`
}

// PaymentScheduleItem represents a single payment in the schedule
type PaymentScheduleItem struct {
	Month            int     `json:"month"`
	Payment          float64 `json:"payment"`
	Principal        float64 `json:"principal"`
	Interest         float64 `json:"interest"`
	RemainingBalance float64 `json:"remaining_balance"`
}

// CalculateCredit calculates credit/loan details
// @Summary Calculate credit
// @Description Calculate credit/loan details including monthly payments and total amounts
// @Tags credit
// @Accept json
// @Produce json
// @Param request body CreditCalculationRequest true "Credit calculation request"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=CreditCalculationResponse} "Credit calculation result"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Router /api/v1/credit/calculate [post]
func (h *Handler) CalculateCredit(c *fiber.Ctx) error {
	var req CreditCalculationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	// Validate down payment
	if req.DownPayment >= req.Price {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.downPaymentTooHigh", nil)
	}

	// Calculate loan amount
	loanAmount := req.Price - req.DownPayment

	// Calculate monthly interest rate
	monthlyRate := req.InterestRate / 100 / 12

	// Calculate monthly payment using the annuity formula
	var monthlyPayment float64
	if monthlyRate > 0 {
		// PMT = P * (r * (1 + r)^n) / ((1 + r)^n - 1)
		// Where:
		// P = loan amount
		// r = monthly interest rate
		// n = number of payments
		temp := math.Pow(1+monthlyRate, float64(req.Term))
		monthlyPayment = loanAmount * (monthlyRate * temp) / (temp - 1)
	} else {
		// No interest, simple division
		monthlyPayment = loanAmount / float64(req.Term)
	}

	// Calculate totals
	totalAmount := monthlyPayment * float64(req.Term)
	totalInterest := totalAmount - loanAmount

	// Generate payment schedule (optional, first 12 months)
	schedule := make([]PaymentScheduleItem, 0)
	remainingBalance := loanAmount

	scheduleMonths := req.Term
	if scheduleMonths > 12 {
		scheduleMonths = 12 // Limit to first 12 months for response size
	}

	for i := 1; i <= scheduleMonths; i++ {
		interestPayment := remainingBalance * monthlyRate
		principalPayment := monthlyPayment - interestPayment
		remainingBalance -= principalPayment

		schedule = append(schedule, PaymentScheduleItem{
			Month:            i,
			Payment:          monthlyPayment,
			Principal:        principalPayment,
			Interest:         interestPayment,
			RemainingBalance: math.Max(0, remainingBalance),
		})
	}

	response := CreditCalculationResponse{
		MonthlyPayment:  math.Round(monthlyPayment*100) / 100,
		TotalAmount:     math.Round(totalAmount*100) / 100,
		TotalInterest:   math.Round(totalInterest*100) / 100,
		LoanAmount:      loanAmount,
		DownPayment:     req.DownPayment,
		Term:            req.Term,
		InterestRate:    req.InterestRate,
		PaymentSchedule: schedule,
	}

	return utils.SendSuccessResponse(c, response, "success")
}
