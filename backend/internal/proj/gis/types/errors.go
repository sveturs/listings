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

)
