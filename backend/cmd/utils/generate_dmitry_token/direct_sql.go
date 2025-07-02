package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type SessionData struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GoogleID   string `json:"google_id"`
	PictureURL string `json:"picture_url"`
	Provider   string `json:"provider"`
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	// Пробуем разные строки подключения
	connectionStrings := []string{
		os.Getenv("DATABASE_URL"),
		"postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable",
		"postgres://postgres:postgres@localhost:5432/hostel_db?sslmode=disable",
		"postgres://postgres:1321321321321@localhost:5432/hostel_db?sslmode=disable",
	}

	var db *sql.DB
	var err error

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

	// Пытаемся найти пользователя
	row := db.QueryRow(`
		SELECT id, name, email, google_id, COALESCE(picture_url, '') 
		FROM users 
		WHERE email = $1
	`, targetEmail)

	err = row.Scan(&userID, &name, &email, &googleID, &pictureURL)

	if err == sql.ErrNoRows {
		// Пользователь не найден - создаем нового
		fmt.Printf("🔍 Пользователь %s не найден. Создаем нового...\n", targetEmail)

		// Генерируем фиктивный Google ID
		fakeGoogleID := "google_" + generateSessionToken()[:20]

		// Создаем пользователя
		err = db.QueryRow(`
			INSERT INTO users (name, email, google_id, picture_url, provider, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "Dmitry Voroshilov", targetEmail, fakeGoogleID,
			"https://lh3.googleusercontent.com/a/default-user=s96-c",
			"google", time.Now()).Scan(&userID)
		if err != nil {
			log.Fatal("Не удалось создать пользователя:", err)
		}

		// Загружаем созданного пользователя
		row = db.QueryRow(`
			SELECT id, name, email, google_id, COALESCE(picture_url, '') 
			FROM users 
			WHERE id = $1
		`, userID)

		err = row.Scan(&userID, &name, &email, &googleID, &pictureURL)
		if err != nil {
			log.Fatal("Не удалось загрузить созданного пользователя:", err)
		}

		fmt.Printf("✅ Создан новый пользователь: %s (ID: %d)\n", email, userID)

	} else if err != nil {
		log.Fatal("Ошибка при поиске пользователя:", err)
	} else {
		fmt.Printf("✅ Найден существующий пользователь: %s (ID: %d)\n", email, userID)
	}

	// Генерируем сессионный токен
	sessionToken := generateSessionToken()

	// Создаем данные сессии
	sessionData := &SessionData{
		UserID:     userID,
		Name:       name,
		Email:      email,
		GoogleID:   googleID,
		PictureURL: pictureURL,
		Provider:   "google",
	}

	// Сохраняем сессию в базу данных
	sessionJSON, _ := json.Marshal(sessionData)
	expiry := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(`
		INSERT INTO user_sessions (id, data, expiry) 
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET data = $2, expiry = $3
	`, sessionToken, string(sessionJSON), expiry)
	if err != nil {
		log.Fatal("Не удалось сохранить сессию:", err)
	}

	// Выводим результат
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("🎉 ТОКЕН АВТОРИЗАЦИИ УСПЕШНО СОЗДАН!\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n👤 Пользователь: %s (ID: %d)\n", email, userID)
	fmt.Printf("\n🔑 Токен:\n%s\n", sessionToken)
	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n1️⃣  В браузере (DevTools Console):")
	fmt.Printf("   document.cookie = 'session_token=%s; path=/'\n", sessionToken)

	fmt.Println("\n2️⃣  В LocalStorage (для SPA):")
	fmt.Printf("   localStorage.setItem('user_session_token', '%s')\n", sessionToken)

	fmt.Println("\n3️⃣  Через URL параметр:")
	fmt.Printf("   http://localhost:3001/?session_token=%s\n", sessionToken)

	fmt.Println("\n4️⃣  cURL запрос для проверки:")
	fmt.Printf("   curl -H \"Cookie: session_token=%s\" \\\n", sessionToken)
	fmt.Println("        http://localhost:3000/api/v1/auth/session")

	fmt.Println("\n5️⃣  Для тестирования API:")
	fmt.Printf("   curl -H \"Cookie: session_token=%s\" \\\n", sessionToken)
	fmt.Println("        http://localhost:3000/api/v1/user/profile")

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⏰ Токен действителен: 24 часа")
	fmt.Println("🔒 Тип авторизации: Google OAuth (эмуляция)")
	fmt.Println(strings.Repeat("=", 80))
}

// Для форматирования вывода
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
