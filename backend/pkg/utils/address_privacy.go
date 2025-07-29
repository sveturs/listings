package utils

import (
	"regexp"
	"strings"
)

// Privacy level constants
const (
	PrivacyLevelExact       = "exact"
	PrivacyLevelStreet      = "street"
	PrivacyLevelApproximate = "approximate"
	PrivacyLevelDistrict    = "district"
	PrivacyLevelCity        = "city"
	PrivacyLevelCityOnly    = "city_only"
	PrivacyLevelHidden      = "hidden"
)

// FormatAddressWithPrivacy форматирует адрес в соответствии с уровнем приватности
func FormatAddressWithPrivacy(address string, privacyLevel string) string {
	if address == "" {
		return ""
	}

	// Для точного адреса возвращаем как есть
	if privacyLevel == PrivacyLevelExact {
		return address
	}

	// Разбираем адрес на части
	parts := strings.Split(address, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	switch privacyLevel {
	case PrivacyLevelStreet, PrivacyLevelApproximate:
		// Уровень улицы - убираем номер дома
		if len(parts) > 0 {
			// Регулярное выражение для удаления номера дома
			// Удаляет числа с возможными буквами в конце (например: 15, 15a, 15б)
			re := regexp.MustCompile(`\b\d+[а-яА-Яa-zA-Z]?\b`)
			parts[0] = re.ReplaceAllString(parts[0], "")
			parts[0] = strings.TrimSpace(parts[0])

			// Если после удаления номера первая часть пустая, убираем её
			if parts[0] == "" && len(parts) > 1 {
				parts = parts[1:]
			}
		}
		return strings.Join(parts, ", ")

	case PrivacyLevelDistrict:
		// Уровень района - показываем только район и город
		if len(parts) >= 2 {
			// Берем последние 2-3 части (обычно район/город, страна)
			return strings.Join(parts[len(parts)-2:], ", ")
		}
		return address

	case PrivacyLevelCity, PrivacyLevelCityOnly:
		// Только город - берем последнюю часть адреса
		if len(parts) > 0 {
			// Обычно город находится в конце адреса
			// Для сербских адресов может быть формат: "улица дом, город почтовый_индекс, округ, страна"
			if len(parts) >= 2 {
				// Берем предпоследнюю часть (обычно город с индексом) и последнюю (страна)
				return strings.Join(parts[len(parts)-2:], ", ")
			}
			return parts[len(parts)-1]
		}
		return address

	case PrivacyLevelHidden:
		// Скрытый адрес - возвращаем только страну или общее указание
		if len(parts) > 0 {
			// Если есть страна (последняя часть), возвращаем её
			lastPart := parts[len(parts)-1]
			if len(lastPart) > 0 {
				return lastPart
			}
		}
		return "Адрес скрыт"

	default:
		return address
	}
}

// GetCoordinatesPrivacy определяет точность координат на основе уровня приватности
// ВАЖНО: Эта функция временная. Для корректной работы нужно использовать
// GetCoordinatesWithGeocoding с сервисом геокодирования
func GetCoordinatesPrivacy(lat, lng float64, privacyLevel string) (float64, float64) {
	switch privacyLevel {
	case PrivacyLevelExact:
		// Точные координаты
		return lat, lng

	case PrivacyLevelStreet, PrivacyLevelApproximate:
		// Временно округляем до ~100 метров (3 знака после запятой)
		// TODO: Использовать геокодирование для получения координат улицы
		return roundToDecimalPlaces(lat, 3), roundToDecimalPlaces(lng, 3)

	case PrivacyLevelDistrict:
		// Временно округляем до ~1 км (2 знака после запятой)
		// TODO: Использовать геокодирование для получения координат района
		return roundToDecimalPlaces(lat, 2), roundToDecimalPlaces(lng, 2)

	case PrivacyLevelCity, PrivacyLevelCityOnly:
		// Временно округляем до ~10 км (1 знак после запятой)
		// TODO: Использовать геокодирование для получения координат города
		return roundToDecimalPlaces(lat, 1), roundToDecimalPlaces(lng, 1)

	case PrivacyLevelHidden:
		// Возвращаем 0,0 для скрытых координат
		return 0, 0

	default:
		return lat, lng
	}
}

// roundToDecimalPlaces округляет число до указанного количества знаков после запятой
func roundToDecimalPlaces(value float64, places int) float64 {
	shift := 1.0
	for i := 0; i < places; i++ {
		shift *= 10
	}
	return float64(int(value*shift+0.5)) / shift
}
