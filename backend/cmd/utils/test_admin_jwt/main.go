package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Используем статичные данные для тестирования
	userID := int64(9) // ID админа из базы данных
	email := "admin@test.com"
	secret := "yoursecretkey" // Из .env файла
	expirationHours := 24

	// Генерируем JWT токен
	jwtToken, err := generateJWT(userID, email, secret, expirationHours)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate JWT token: %v", err))
	}

	// Выводим результат
	fmt.Println("================================================================================")
	fmt.Println("🎉 ТЕСТОВЫЙ ADMIN JWT ТОКЕН СОЗДАН!")
	fmt.Println("================================================================================")
	fmt.Printf("\n👤 Администратор: %s (ID: %d)\n", email, userID)
	fmt.Printf("🔑 Admin права: true\n")

	fmt.Println("\n🔑 ACCESS TOKEN (JWT):")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%s\n", jwtToken)

	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Println("\n1️⃣  Тест API синонимов:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms?page=1&limit=20&language=ru'")

	fmt.Println("\n2️⃣  Добавление синонима:")
	fmt.Printf("   curl -X POST -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        -H \"Content-Type: application/json\" \\")
	fmt.Println("        -d '{\"word\": \"телефон\", \"synonyms\": [\"смартфон\", \"мобильный\"], \"language\": \"ru\"}' \\")
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms'")

	fmt.Println("\n3️⃣  Для использования в браузере (установка в localStorage):")
	fmt.Printf("   localStorage.setItem('access_token', '%s');\n", jwtToken)

	fmt.Println("\n================================================================================")
	fmt.Printf("⏰ Access токен действителен: %d часов\n", expirationHours)
	fmt.Println("🔒 Тип авторизации: JWT Bearer")
	fmt.Println("================================================================================")
}

func generateJWT(userID int64, email string, secret string, expirationHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iss":     "svetu-backend",
		"sub":     fmt.Sprintf("user:%d", userID),
		"exp":     time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
