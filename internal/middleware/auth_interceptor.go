package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vondi-global/auth/pkg/entity"
	authservice "github.com/vondi-global/auth/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor provides JWT authentication for gRPC
type AuthInterceptor struct {
	authService *authservice.AuthService
	logger      zerolog.Logger
}

// NewAuthInterceptor creates new auth interceptor
func NewAuthInterceptor(authService *authservice.AuthService, logger zerolog.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		authService: authService,
		logger:      logger,
	}
}

// Unary returns unary server interceptor for JWT validation
func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract JWT from metadata
		token, err := a.extractToken(ctx)
		if err != nil {
			a.logger.Debug().
				Err(err).
				Str("method", info.FullMethod).
				Msg("No JWT token in request")

			// For public methods (optional auth) - allow without token
			if a.isPublicMethod(info.FullMethod) {
				return handler(ctx, req)
			}

			return nil, status.Error(codes.Unauthenticated, "missing authentication token")
		}

		// Validate JWT using auth service
		claims, err := a.authService.ValidateToken(ctx, token)
		if err != nil {
			a.logger.Warn().
				Err(err).
				Str("method", info.FullMethod).
				Msg("Invalid JWT token")
			return nil, status.Error(codes.Unauthenticated, "invalid authentication token")
		}

		// Add user claims to context
		ctx = a.enrichContext(ctx, claims)

		a.logger.Debug().
			Int("user_id", claims.UserID).
			Str("email", claims.Email).
			Str("method", info.FullMethod).
			Msg("JWT authentication successful")

		// Call handler with enriched context
		return handler(ctx, req)
	}
}

// extractToken extracts JWT token from gRPC metadata
func (a *AuthInterceptor) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	// Check Authorization header
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", fmt.Errorf("no authorization header")
	}

	// Extract token from "Bearer <token>" format
	authHeader := authHeaders[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

// enrichContext adds user claims to context
func (a *AuthInterceptor) enrichContext(ctx context.Context, validationResp *entity.TokenValidationResponse) context.Context {
	// Add claims to context using custom keys
	ctx = context.WithValue(ctx, UserIDKey{}, int64(validationResp.UserID))
	ctx = context.WithValue(ctx, EmailKey{}, validationResp.Email)
	ctx = context.WithValue(ctx, RolesKey{}, validationResp.Roles)
	return ctx
}

// isPublicMethod checks if method requires authentication
func (a *AuthInterceptor) isPublicMethod(method string) bool {
	// Public methods that don't require authentication
	publicMethods := []string{
		"/listingssvc.v1.ListingsService/GetRootCategories",
		"/listingssvc.v1.ListingsService/GetAllCategories",
		"/listingssvc.v1.ListingsService/GetCategory",
		"/listingssvc.v1.ListingsService/GetCategoryBySlug",
		"/listingssvc.v1.ListingsService/SearchListings",
		"/listingssvc.v1.ListingsService/ListListings",
		"/listingssvc.v1.ListingsService/GetListing",
		"/listingssvc.v1.ListingsService/GetSimilarListings",
		"/listingssvc.v1.ListingsService/GetProduct",
		"/listingssvc.v1.ListingsService/GetProductBySKU",
		"/listingssvc.v1.ListingsService/ListProducts",
		// Storefront public methods
		"/listingssvc.v1.ListingsService/GetStorefront",
		"/listingssvc.v1.ListingsService/GetStorefrontBySlug",
		"/listingssvc.v1.ListingsService/ListStorefronts",
		"/listingssvc.v1.ListingsService/GetMyStorefronts",
		"/listingssvc.v1.ListingsService/GetWorkingHours",
		"/listingssvc.v1.ListingsService/IsOpenNow",
		"/listingssvc.v1.ListingsService/GetStaff",
		"/listingssvc.v1.ListingsService/GetPaymentMethods",
		"/listingssvc.v1.ListingsService/GetDeliveryOptions",
		"/listingssvc.v1.ListingsService/GetMapData",
		// Attributes public methods (no auth required for viewing category attributes)
		"/listingssvc.v1.AttributeService/GetCategoryAttributes",
		"/listingssvc.v1.AttributeService/GetCategoryVariantAttributes",
		// Category service public methods
		"/categoriessvc.v1.CategoryService/GetCategories",
		"/categoriessvc.v1.CategoryService/GetCategoryBySlug",
		"/categoriessvc.v1.CategoryService/GetRootCategories",
		"/categoriessvc.v1.CategoryService/GetCategoryChildren",
		"/categoriessvc.v1.CategoryService/GetCategoryPath",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

// Context keys for user claims
type UserIDKey struct{}
type EmailKey struct{}
type RolesKey struct{}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(int64)
	return userID, ok
}

// GetEmail extracts email from context
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey{}).(string)
	return email, ok
}

// GetRoles extracts roles from context
func GetRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey{}).([]string)
	return roles, ok
}

// IsAuthenticated checks if request is authenticated
func IsAuthenticated(ctx context.Context) bool {
	_, ok := GetUserID(ctx)
	return ok
}

// HasRole checks if user has specific role
func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetRoles(ctx)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
