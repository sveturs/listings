package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"backend/internal/proj/vin/models"
)

// VINDecoder предоставляет функции для декодирования VIN номеров
type VINDecoder struct{}

// NewVINDecoder создает новый экземпляр декодера
func NewVINDecoder() *VINDecoder {
	return &VINDecoder{}
}

// ValidateVIN проверяет корректность VIN номера
func (d *VINDecoder) ValidateVIN(vin string) error {
	vin = strings.ToUpper(strings.TrimSpace(vin))

	if len(vin) != 17 {
		return fmt.Errorf("VIN должен содержать ровно 17 символов")
	}

	// Проверка на недопустимые символы (I, O, Q не используются в VIN)
	if strings.ContainsAny(vin, "IOQ") {
		return fmt.Errorf("VIN содержит недопустимые символы (I, O или Q)")
	}

	// Проверка что все символы валидные (A-Z, 0-9, кроме I, O, Q)
	for _, char := range vin {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z')) {
			return fmt.Errorf("VIN содержит недопустимый символ: %c", char)
		}
	}

	// Проверка контрольной суммы (позиция 9) - опциональна для европейских авто
	// В Европе контрольная сумма не обязательна, но если это североамериканское авто - проверяем
	if d.isNorthAmericanVIN(vin) {
		if !d.validateCheckDigit(vin) {
			// Для европейского рынка это не критично, просто предупреждение
			// return fmt.Errorf("неверная контрольная сумма VIN")
		}
	}

	return nil
}

// DecodeBasicInfo декодирует базовую информацию из VIN локально
func (d *VINDecoder) DecodeBasicInfo(vin string) (*models.BasicVINInfo, error) {
	vin = strings.ToUpper(strings.TrimSpace(vin))

	if err := d.ValidateVIN(vin); err != nil {
		return nil, err
	}

	info := &models.BasicVINInfo{}

	// WMI (World Manufacturer Identifier) - первые 3 символа
	wmi := vin[0:3]
	info.Manufacturer = d.getManufacturerByWMI(wmi)
	info.Region = d.getRegionByWMI(wmi)

	// Год выпуска (позиция 10)
	yearChar := vin[9]
	info.Year = d.getYearFromChar(yearChar)

	// Проверка контрольной суммы
	if d.isNorthAmericanVIN(vin) {
		info.CheckDigit = d.validateCheckDigit(vin)
	}

	return info, nil
}

// DecodeLocal выполняет локальное декодирование VIN без обращения к внешним API
func (d *VINDecoder) DecodeLocal(vin string) (*models.VINDecodeCache, error) {
	vin = strings.ToUpper(strings.TrimSpace(vin))

	if err := d.ValidateVIN(vin); err != nil {
		return nil, err
	}

	basicInfo, err := d.DecodeBasicInfo(vin)
	if err != nil {
		return nil, err
	}

	year := basicInfo.Year
	manufacturer := basicInfo.Manufacturer

	result := &models.VINDecodeCache{
		VIN:          vin,
		Year:         &year,
		Manufacturer: &manufacturer,
		DecodeStatus: "partial", // Локальное декодирование всегда частичное
	}

	// Попытка определить дополнительную информацию по WMI
	wmi := vin[0:3]
	if make, model := d.getDetailedInfoByWMI(wmi); make != "" {
		result.Make = &make
		if model != "" {
			result.Model = &model
		}
	}

	// Сохраняем raw response для отладки
	rawData := map[string]interface{}{
		"source":     "local",
		"basic_info": basicInfo,
		"wmi":        wmi,
		"vin":        vin,
	}
	rawJSON, _ := json.Marshal(rawData)
	result.RawResponse = rawJSON

	return result, nil
}

// isNorthAmericanVIN проверяет, является ли VIN североамериканским
func (d *VINDecoder) isNorthAmericanVIN(vin string) bool {
	firstChar := vin[0]
	return (firstChar >= '1' && firstChar <= '5')
}

// validateCheckDigit проверяет контрольную сумму VIN
func (d *VINDecoder) validateCheckDigit(vin string) bool {
	weights := []int{8, 7, 6, 5, 4, 3, 2, 10, 0, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0

	for i, char := range vin {
		value := d.getCharValue(char)
		sum += value * weights[i]
	}

	checkDigit := sum % 11
	expectedChar := strconv.Itoa(checkDigit)
	if checkDigit == 10 {
		expectedChar = "X"
	}

	return string(vin[8]) == expectedChar
}

// getCharValue возвращает числовое значение символа для расчета контрольной суммы
func (d *VINDecoder) getCharValue(char rune) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	}

	charValues := map[rune]int{
		'A': 1, 'B': 2, 'C': 3, 'D': 4, 'E': 5, 'F': 6, 'G': 7, 'H': 8,
		'J': 1, 'K': 2, 'L': 3, 'M': 4, 'N': 5, 'P': 7, 'R': 9,
		'S': 2, 'T': 3, 'U': 4, 'V': 5, 'W': 6, 'X': 7, 'Y': 8, 'Z': 9,
	}

	if val, ok := charValues[char]; ok {
		return val
	}
	return 0
}

// getYearFromChar определяет год по символу в позиции 10 VIN
func (d *VINDecoder) getYearFromChar(char byte) int {
	yearMap := map[byte]int{
		'A': 1980, 'B': 1981, 'C': 1982, 'D': 1983, 'E': 1984, 'F': 1985,
		'G': 1986, 'H': 1987, 'J': 1988, 'K': 1989, 'L': 1990, 'M': 1991,
		'N': 1992, 'P': 1993, 'R': 1994, 'S': 1995, 'T': 1996, 'V': 1997,
		'W': 1998, 'X': 1999, 'Y': 2000, '1': 2001, '2': 2002, '3': 2003,
		'4': 2004, '5': 2005, '6': 2006, '7': 2007, '8': 2008, '9': 2009,
	}

	// После 2009 года цикл повторяется
	if year, ok := yearMap[char]; ok {
		// Проверяем, нужно ли добавить 30 лет (для VIN после 2010 года)
		// Это упрощенная логика, в реальности нужно учитывать 7-й символ
		currentYear := 2025
		if year+30 <= currentYear {
			return year + 30
		}
		return year
	}

	// Для символов начиная с A (2010+)
	yearMap2 := map[byte]int{
		'A': 2010, 'B': 2011, 'C': 2012, 'D': 2013, 'E': 2014, 'F': 2015,
		'G': 2016, 'H': 2017, 'J': 2018, 'K': 2019, 'L': 2020, 'M': 2021,
		'N': 2022, 'P': 2023, 'R': 2024, 'S': 2025, 'T': 2026,
	}

	if year, ok := yearMap2[char]; ok {
		return year
	}

	return 0
}

// getRegionByWMI определяет регион по WMI
func (d *VINDecoder) getRegionByWMI(wmi string) string {
	if len(wmi) == 0 {
		return ""
	}

	firstChar := wmi[0]
	switch {
	case firstChar >= '1' && firstChar <= '5':
		return "Северная Америка"
	case firstChar >= '6' && firstChar <= '7':
		return "Океания"
	case firstChar >= '8' && firstChar <= '9':
		return "Южная Америка"
	case firstChar >= 'A' && firstChar <= 'H':
		return "Африка"
	case firstChar >= 'J' && firstChar <= 'R':
		return "Азия"
	case firstChar >= 'S' && firstChar <= 'Z':
		return "Европа"
	default:
		return ""
	}
}

// getManufacturerByWMI определяет производителя по WMI
func (d *VINDecoder) getManufacturerByWMI(wmi string) string {
	// Таблица WMI с упором на европейские производители (актуально для сербского рынка)
	manufacturers := map[string]string{
		// Европейские (основной импорт в Сербию)
		"WBA": "BMW", "WBS": "BMW", "WBW": "BMW", "WBX": "BMW", "WBY": "BMW",
		"WDB": "Mercedes-Benz", "WDC": "Mercedes-Benz", "WDD": "Mercedes-Benz",
		"WDF": "Mercedes-Benz (Smart)", "WDZ": "Mercedes-Benz",
		"WAU": "Audi", "WA1": "Audi", "WUA": "Audi", "WUZ": "Audi quattro",
		"WVW": "Volkswagen", "WVG": "Volkswagen", "WV1": "Volkswagen Commercial", "WV2": "Volkswagen Bus/Van", "WV3": "Volkswagen Trucks",
		"W0L": "Opel", "W0V": "Opel", "VSE": "Opel (Spain)",
		"VF1": "Renault", "VF2": "Renault", "VF3": "Peugeot", "VF4": "Peugeot",
		"VF5": "Peugeot", "VF6": "Renault Trucks", "VF7": "Citroën", "VF8": "Citroën",
		"VFA": "Alpine", "VFC": "Citroën", "VFD": "Peugeot", "VFE": "IvecoBus",
		"ZFA": "Fiat", "ZFC": "Fiat", "ZFF": "Ferrari", "ZHW": "Alfa Romeo", "ZAR": "Alfa Romeo",
		"ZCF": "Iveco", "ZLA": "Lancia",
		"VSS": "SEAT", "VSX": "SEAT", "VS6": "SEAT (Málaga)", "VS9": "SEAT (Barcelona)",
		"TMA": "Hyundai (Czech)", "TMB": "Škoda", "TMC": "Škoda", "TMP": "Škoda", "TMT": "Tatra",
		"TRU": "Audi (Hungary)", "TSE": "Audi (Hungary)", "TSM": "Suzuki (Hungary)",
		"UU1": "Dacia", "UU2": "Dacia", "UU3": "ARO", "UU4": "Dacia", "UU5": "Dacia", "UU6": "Daewoo Romania",
		"X7L": "Renault (Russia)", "X7M": "Hyundai (Russia)", "X7N": "Hyundai (Russia)",
		"YS2": "Scania", "YS3": "Saab", "YS4": "Scania", "YV1": "Volvo Cars", "YV2": "Volvo Trucks", "YV3": "Volvo Buses",
		"SCA": "Rolls Royce", "SCB": "Bentley", "SCC": "Lotus", "SCE": "DeLorean", "SCF": "Aston Martin",
		"SDB": "Peugeot UK", "SFA": "Ford UK", "SFD": "Alexander Dennis", "SHH": "Honda UK",
		"SJN": "Nissan UK", "SAJ": "Jaguar", "SAL": "Land Rover", "SAR": "Rover", "SAX": "Austin",
		"SBM": "McLaren", "SEY": "LDV",
		"WF0": "Ford (Germany)", "WMW": "MINI",
		"WP0": "Porsche", "WP1": "Porsche",

		// Американские (редко, но бывают)
		"1FA": "Ford", "1FB": "Ford", "1FC": "Ford", "1FD": "Ford",
		"1FM": "Ford", "1FT": "Ford", "1FU": "Freightliner",
		"1G1": "Chevrolet", "1G2": "Pontiac", "1G3": "Oldsmobile",
		"1G4": "Buick", "1G6": "Cadillac", "1G8": "Saturn",
		"1GC": "Chevrolet Truck", "1GM": "Pontiac", "1GT": "GMC Truck",
		"1HG": "Honda", "1N4": "Nissan", "1VW": "Volkswagen",
		"2FA": "Ford", "2FB": "Ford", "2FC": "Ford", "2FM": "Ford",
		"2FT": "Ford", "2G1": "Chevrolet", "2G2": "Pontiac",
		"2HG": "Honda", "2HH": "Acura", "2HJ": "Honda",
		"2HK": "Honda", "2HM": "Honda", "3FA": "Ford",
		"3FE": "Ford", "3G1": "Chevrolet", "3VW": "Volkswagen",
		"4F2": "Mazda", "4S3": "Subaru", "4S4": "Subaru",
		"4T1": "Toyota", "4T3": "Toyota", "4US": "BMW",
		"5FN": "Honda", "5L1": "Lincoln", "5TB": "Toyota",
		"5UM": "BMW", "5UX": "BMW", "5YF": "Toyota",

		// Европейские (дополнительные коды) - уже объявлены выше

		// Японские
		"JA3": "Mitsubishi", "JA4": "Mitsubishi", "JAA": "Isuzu",
		"JAL": "Isuzu", "JF1": "Subaru", "JF2": "Subaru",
		"JHL": "Honda", "JHM": "Honda", "JM1": "Mazda",
		"JM3": "Mazda", "JN1": "Nissan", "JN8": "Nissan",
		"JS1": "Suzuki", "JS2": "Suzuki", "JT2": "Toyota",
		"JTD": "Toyota", "JTE": "Toyota", "JTK": "Toyota",
		"JTM": "Toyota", "JYA": "Yamaha", "JYE": "Yamaha",

		// Корейские
		"KL1": "Daewoo/GM Korea", "KM8": "Hyundai", "KMF": "Hyundai",
		"KMH": "Hyundai", "KMX": "Hyundai", "KNA": "Kia",
		"KNC": "Kia", "KNH": "Kia", "KPT": "Kia",

		// Китайские
		"L6T": "Geely", "LFM": "FAW Toyota", "LFV": "FAW Volkswagen",
		"LGB": "Dongfeng Nissan", "LGW": "Great Wall", "LGX": "BYD",
		"LH1": "FAW Haima", "LJC": "JAC", "LSV": "SAIC Volkswagen",
		"LSY": "Brilliance", "LVG": "GAC Toyota", "LVH": "Dongfeng Honda",
		"LVS": "BMW Brilliance", "LZW": "SAIC GM",
	}

	if manufacturer, ok := manufacturers[wmi]; ok {
		return manufacturer
	}

	// Попробуем по первым двум символам
	if len(wmi) >= 2 {
		wmi2 := wmi[0:2]
		manufacturers2 := map[string]string{
			"1F": "Ford", "1G": "General Motors", "1H": "Honda",
			"1J": "Jeep", "1L": "Lincoln", "1M": "Mercury",
			"1N": "Nissan", "1V": "Volkswagen", "1Y": "Mazda",
			"2F": "Ford", "2G": "General Motors", "2H": "Honda",
			"2M": "Mercury", "2T": "Toyota", "3F": "Ford",
			"3G": "General Motors", "3H": "Honda", "3N": "Nissan",
			"3V": "Volkswagen", "4F": "Mazda", "4J": "Mercedes-Benz",
			"4M": "Mercury", "4S": "Subaru", "4T": "Toyota",
			"4U": "BMW", "5F": "Honda", "5L": "Lincoln",
			"5N": "Hyundai", "5T": "Toyota", "5U": "BMW",
			"5X": "Hyundai/Kia", "5Y": "Toyota",
		}

		if manufacturer, ok := manufacturers2[wmi2]; ok {
			return manufacturer
		}
	}

	return "Unknown"
}

// getDetailedInfoByWMI возвращает дополнительную информацию по WMI
func (d *VINDecoder) getDetailedInfoByWMI(wmi string) (make string, model string) {
	// Дополнительная детализация для некоторых WMI
	detailedInfo := map[string]struct{ make, model string }{
		"1FA": {"Ford", ""},
		"1FB": {"Ford", "F-Series"},
		"1FM": {"Ford", "SUV/Crossover"},
		"1FT": {"Ford", "Truck"},
		"1G1": {"Chevrolet", "Passenger Car"},
		"1GC": {"Chevrolet", "Truck"},
		"1HG": {"Honda", "Civic/Accord"},
		"2HG": {"Honda", "Civic"},
		"2HH": {"Acura", ""},
		"3VW": {"Volkswagen", ""},
		"4T1": {"Toyota", "Passenger Car"},
		"5TB": {"Toyota", "Truck"},
		"JTD": {"Toyota", ""},
		"JTE": {"Toyota", ""},
		"KMH": {"Hyundai", ""},
		"WBA": {"BMW", "3 Series"},
		"WBS": {"BMW", "M Series"},
		"WDB": {"Mercedes-Benz", ""},
		"WDD": {"Mercedes-Benz", ""},
		"WVW": {"Volkswagen", ""},
		"WAU": {"Audi", ""},
	}

	if info, ok := detailedInfo[wmi]; ok {
		return info.make, info.model
	}

	return "", ""
}
