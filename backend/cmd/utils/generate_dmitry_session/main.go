package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// Генерация случайного токена сессии
func generateSessionToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	// Данные пользователя Dmitry Voroshilov из ответа API
	userID := 14
	email := "voroshilovdo@gmail.com"
	name := "Dmitry Voroshilov"

	// Генерируем сессионный токен
	sessionToken := generateSessionToken()
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("🎉 СЕССИОННЫЙ ТОКЕН ДЛЯ БРАУЗЕРА СОЗДАН!\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n👤 Пользователь: %s (ID: %d)\n", email, userID)
	fmt.Printf("📧 Имя: %s\n", name)
	
	fmt.Println("\n🍪 SESSION TOKEN (для cookie):")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s\n", sessionToken)
	
	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("\n1️⃣  В браузере через DevTools Console:")
	fmt.Printf("   document.cookie = 'session_token=%s; path=/; max-age=86400';\n", sessionToken)
	fmt.Println("   // Затем обновите страницу")
	
	fmt.Println("\n2️⃣  В LocalStorage (для frontend):")
	fmt.Printf("   localStorage.setItem('user_session_token', '%s');\n", sessionToken)
	fmt.Println("   // Frontend автоматически подхватит токен")
	
	fmt.Println("\n3️⃣  Через URL параметр:")
	fmt.Printf("   http://localhost:3001/?session_token=%s\n", sessionToken)
	
	fmt.Println("\n4️⃣  cURL запрос с cookie:")
	fmt.Printf("   curl -H \"Cookie: session_token=%s\" \\\n", sessionToken)
	fmt.Println("        http://localhost:3000/api/v1/auth/session")
	
	fmt.Println("\n⚠️  ВАЖНО:")
	fmt.Println("   Этот токен сессии НЕ сохранен на сервере!")
	fmt.Println("   Для полноценной авторизации используйте JWT токен из предыдущей утилиты.")
	fmt.Println("   Или войдите через Google OAuth для создания настоящей сессии.")
	
	fmt.Println("\n💡 АЛЬТЕРНАТИВА - Эмуляция Google OAuth:")
	fmt.Printf("   1. Откройте: http://localhost:3001\n")
	fmt.Printf("   2. Нажмите 'Войти через Google'\n")
	fmt.Printf("   3. Выберите аккаунт: %s\n", email)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⏰ Сессия действительна: пока работает сервер")
	fmt.Println("🔒 Тип авторизации: Session Cookie")
	fmt.Println(strings.Repeat("=", 80))
}