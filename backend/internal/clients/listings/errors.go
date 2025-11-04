package listings

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Domain errors - маппинг gRPC ошибок на доменные ошибки
var (
	// ErrServiceUnavailable - микросервис временно недоступен
	ErrServiceUnavailable = errors.New("listings service temporarily unavailable")

	// ErrListingNotFound - listing не найден
	ErrListingNotFound = errors.New("listing not found")

	// ErrInvalidInput - невалидные входные данные
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized - недостаточно прав доступа
	ErrUnauthorized = errors.New("unauthorized")

	// ErrAlreadyExists - listing уже существует
	ErrAlreadyExists = errors.New("listing already exists")

	// ErrInternal - внутренняя ошибка микросервиса
	ErrInternal = errors.New("internal service error")
)

// MapGRPCError преобразует gRPC ошибку в доменную ошибку
func MapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		// Если не gRPC ошибка - возвращаем как есть
		return err
	}

	//nolint:exhaustive // We handle the most common cases; others return the original error
	switch st.Code() {
	case codes.OK:
		return nil

	case codes.NotFound:
		return ErrListingNotFound

	case codes.InvalidArgument:
		return ErrInvalidInput

	case codes.PermissionDenied, codes.Unauthenticated:
		return ErrUnauthorized

	case codes.AlreadyExists:
		return ErrAlreadyExists

	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
		return ErrServiceUnavailable

	case codes.Internal, codes.Unknown, codes.DataLoss:
		return ErrInternal

	default:
		// Для остальных кодов возвращаем оригинальную ошибку
		return err
	}
}

// IsNotFound проверяет, является ли ошибка "not found"
func IsNotFound(err error) bool {
	return errors.Is(err, ErrListingNotFound)
}

// IsInvalidInput проверяет, является ли ошибка невалидным вводом
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsUnauthorized проверяет, является ли ошибка неавторизованным доступом
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsServiceUnavailable проверяет, недоступен ли сервис
func IsServiceUnavailable(err error) bool {
	return errors.Is(err, ErrServiceUnavailable)
}
