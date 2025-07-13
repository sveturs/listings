package types

import "errors"

var (
	// ErrInvalidBounds ошибка некорректных границ
	ErrInvalidBounds = errors.New("invalid bounds")

	// ErrInvalidLatitude ошибка некорректной широты
	ErrInvalidLatitude = errors.New("latitude must be between -90 and 90")

	// ErrInvalidLongitude ошибка некорректной долготы
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")

	// ErrInvalidRadius ошибка некорректного радиуса
	ErrInvalidRadius = errors.New("radius must be positive")

	// ErrLocationNotFound ошибка отсутствия геоданных
	ErrLocationNotFound = errors.New("location not found")

	// ========== PHASE 2: Новые ошибки ==========

	// ErrListingNotFound ошибка отсутствия объявления
	ErrListingNotFound = errors.New("listing not found")

	// ErrAccessDenied ошибка отсутствия прав доступа
	ErrAccessDenied = errors.New("access denied")

	// ErrInvalidPrivacyLevel ошибка некорректного уровня приватности
	ErrInvalidPrivacyLevel = errors.New("invalid privacy level")

	// ErrInvalidInputMethod ошибка некорректного метода ввода
	ErrInvalidInputMethod = errors.New("invalid input method")

	// ErrGeocodingFailed ошибка неудачного геокодирования
	ErrGeocodingFailed = errors.New("geocoding failed")

	// ErrLowConfidence ошибка низкого доверия к геокодированию
	ErrLowConfidence = errors.New("geocoding confidence too low")

	// ErrCacheNotFound ошибка отсутствия в кэше
	ErrCacheNotFound = errors.New("not found in cache")

	// ErrInvalidAddressComponents ошибка некорректных компонентов адреса
	ErrInvalidAddressComponents = errors.New("invalid address components")
)
