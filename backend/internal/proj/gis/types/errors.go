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

	// ========== PHASE 2.5: Ошибки для городов и контекстного поиска ==========

	// ErrMissingBounds ошибка отсутствия границ
	ErrMissingBounds = errors.New("bounds are required")

	// ErrMissingCenter ошибка отсутствия центра
	ErrMissingCenter = errors.New("center point is required")

	// ErrCityNotFound ошибка отсутствия города
	ErrCityNotFound = errors.New("city not found")

	// ErrDistrictNotFound ошибка отсутствия района
	ErrDistrictNotFound = errors.New("district not found")

	// ErrMunicipalityNotFound ошибка отсутствия муниципалитета
	ErrMunicipalityNotFound = errors.New("municipality not found")

	// ErrNoCitiesInViewport ошибка отсутствия городов в области просмотра
	ErrNoCitiesInViewport = errors.New("no cities found in viewport")

	// ErrNoDistrictsAvailable ошибка отсутствия доступных районов
	ErrNoDistrictsAvailable = errors.New("no districts available for this city")
)
