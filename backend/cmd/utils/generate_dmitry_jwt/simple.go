package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// JWT Claims структура
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Генерация случайного токена
func generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Генерация JWT токена
func generateJWT(userID int, email string, secret string, expHours int) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func main() {
	// JWT секрет из переменной окружения или дефолтный
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "yoursecretkey" // Дефолтный секрет из .env файла
	}

	// Строки подключения к БД
	connectionStrings := []string{
		os.Getenv("DATABASE_URL"),
		"postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable",
		"postgres://postgres:postgres@localhost:5432/hostel_db?sslmode=disable",
		"postgres://postgres:1321321321321@localhost:5432/hostel_db?sslmode=disable",
	}

	var db *sql.DB
	var err error

	// Пробуем разные строки подключения
	for _, connStr := range connectionStrings {
		if connStr == "" {
			continue
		}
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Printf("✅ Успешное подключение к базе данных\n")
				break
			}
		}
	}

	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	defer db.Close()

	// Целевой email
	targetEmail := "voroshilovdo@gmail.com"

	var userID int
	var name, email, googleID, pictureURL string

	// Получаем пользователя
	row := db.QueryRow(`
		SELECT id, name, email, google_id, COALESCE(picture_url, '') 
		FROM users 
		WHERE email = $1
	`, targetEmail)

	err = row.Scan(&userID, &name, &email, &googleID, &pictureURL)
	if err != nil {
		log.Fatal("Пользователь не найден:", err)
	}

	// Генерируем JWT токен (24 часа)
	jwtToken, err := generateJWT(userID, email, jwtSecret, 24)
	if err != nil {
		log.Fatal("Не удалось сгенерировать JWT токен:", err)
	}

	// Генерируем refresh токен
	refreshToken := generateRandomToken()

	// Сохраняем refresh токен в БД
	_, err = db.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, refreshToken, time.Now().Add(30*24*time.Hour), time.Now(), "CLI Tool", "127.0.0.1")
	if err != nil {
		// Если таблица не существует, продолжаем без refresh токена
		fmt.Printf("⚠️  Не удалось сохранить refresh токен (возможно, таблица не существует): %v\n", err)
		refreshToken = "Недоступен (таблица refresh_tokens не найдена)"
	}

	// Красивый вывод результата
	fmt.Println("\n" + repeatChar("=", 80))
	fmt.Printf("🎉 ТОКЕНЫ АВТОРИЗАЦИИ УСПЕШНО СОЗДАНЫ!\n")
	fmt.Println(repeatChar("=", 80))
	fmt.Printf("\n👤 Пользователь: %s (ID: %d)\n", email, userID)
	fmt.Printf("📧 Имя: %s\n", name)
	if pictureURL != "" {
		fmt.Printf("🖼️  Фото: %s\n", pictureURL)
	}

	fmt.Println("\n🔑 ACCESS TOKEN (JWT):")
	fmt.Println(repeatChar("-", 80))
	fmt.Printf("%s\n", jwtToken)

	if refreshToken != "Недоступен (таблица refresh_tokens не найдена)" {
		fmt.Println("\n🔄 REFRESH TOKEN:")
		fmt.Println(repeatChar("-", 80))
		fmt.Printf("%s\n", refreshToken)
	}

	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ JWT ТОКЕНА:")
	fmt.Println(repeatChar("-", 80))

	fmt.Println("\n1️⃣  Authorization header (рекомендуется):")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        http://localhost:3000/api/v1/user/profile")

	fmt.Println("\n2️⃣  В JavaScript (axios):")
	fmt.Println("   const token = '" + jwtToken + "';")
	fmt.Println("   axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;")

	fmt.Println("\n3️⃣  В JavaScript (fetch):")
	fmt.Println("   const token = '" + jwtToken + "';")
	fmt.Println("   fetch('http://localhost:3000/api/v1/user/profile', {")
	fmt.Println("     headers: {")
	fmt.Println("       'Authorization': `Bearer ${token}`")
	fmt.Println("     }")
	fmt.Println("   })")

	fmt.Println("\n4️⃣  В Frontend (localStorage):")
	fmt.Println("   const token = '" + jwtToken + "';")
	fmt.Println("   localStorage.setItem('auth_token', token);")
	fmt.Println("   // Затем используйте в запросах")

	fmt.Println("\n5️⃣  Проверка токена:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        http://localhost:3000/api/v1/auth/me")

	fmt.Println("\n6️⃣  Для тестирования в Postman:")
	fmt.Println("   - Выберите тип авторизации: Bearer Token")
	fmt.Printf("   - Вставьте токен: %s\n", jwtToken)

	fmt.Println("\n" + repeatChar("=", 80))
	fmt.Println("⏰ Access токен действителен: 24 часа")
	if refreshToken != "Недоступен (таблица refresh_tokens не найдена)" {
		fmt.Println("⏰ Refresh токен действителен: 30 дней")
	}
	fmt.Println("🔒 Тип авторизации: JWT Bearer")
	fmt.Println("🔐 Алгоритм подписи: HS256")
	fmt.Println(repeatChar("=", 80))

	// Декодируем токен для проверки
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(*Claims); ok {
			fmt.Printf("\n✅ Токен валиден. Истекает: %s\n", claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))
		}
	}
}

// Вспомогательная функция для повторения символов
func repeatChar(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}
