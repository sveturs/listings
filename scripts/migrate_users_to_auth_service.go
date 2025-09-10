package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Основная БД
	mainDBURL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
	
	// Auth Service БД - пароль экранируем
	authDBURL = fmt.Sprintf("postgres://auth_user:%s@localhost:25432/auth_db?sslmode=disable", 
		url.QueryEscape("AuthP@ssw0rd2025!"))
	
)

const (
	// Роль админа
	adminRoleID = 1
	userRoleID  = 2
)

type MainUser struct {
	ID        int
	Email     string
	Name      string
	Password  sql.NullString
	CreatedAt time.Time
}

type AuthUser struct {
	ID             int
	Email          string
	Name           string
	PasswordHash   sql.NullString
	GoogleID       sql.NullString
	Provider       string
	EmailVerified  bool
	IsActive       bool
	CreatedAt      time.Time
}

func main() {
	log.Println("Starting user migration from main DB to Auth Service...")

	// Подключаемся к основной БД
	mainDB, err := sql.Open("postgres", mainDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to main DB: %v", err)
	}
	defer mainDB.Close()

	// Подключаемся к Auth Service БД
	authDB, err := sql.Open("postgres", authDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to Auth Service DB: %v", err)
	}
	defer authDB.Close()

	// Создаем резервную копию текущих данных Auth Service
	log.Println("Creating backup of Auth Service users...")
	backupAuthUsers(authDB)

	// Получаем админов из основной БД
	adminEmails := getAdminEmails(mainDB)
	log.Printf("Found %d admin emails", len(adminEmails))

	// Получаем всех пользователей из основной БД
	users, err := getMainUsers(mainDB)
	if err != nil {
		log.Fatalf("Failed to get users from main DB: %v", err)
	}
	log.Printf("Found %d users in main DB", len(users))

	// Очищаем таблицу пользователей в Auth Service (кроме системных)
	log.Println("Clearing Auth Service users table...")
	if err := clearAuthUsers(authDB); err != nil {
		log.Fatalf("Failed to clear Auth Service users: %v", err)
	}

	// Мигрируем пользователей
	migrated := 0
	failed := 0
	
	for _, user := range users {
		if err := migrateUser(authDB, user, adminEmails); err != nil {
			log.Printf("Failed to migrate user %d (%s): %v", user.ID, user.Email, err)
			failed++
		} else {
			log.Printf("Migrated user %d: %s", user.ID, user.Email)
			migrated++
		}
	}

	// Сбрасываем sequence для ID
	resetSequence(authDB)

	log.Printf("\nMigration completed!")
	log.Printf("Successfully migrated: %d users", migrated)
	log.Printf("Failed: %d users", failed)
	
	// Проверяем результат
	verifyMigration(authDB)
}

func getAdminEmails(db *sql.DB) map[string]bool {
	adminEmails := make(map[string]bool)
	
	// Жестко заданные админы
	hardcodedAdmins := []string{
		"bevzenko.sergey@gmail.com",
		"voroshilovdo@gmail.com",
		"admin@svetu.rs",
		"boxmail386@gmail.com", // Добавляем текущего пользователя как админа
	}
	
	for _, email := range hardcodedAdmins {
		adminEmails[email] = true
	}
	
	// Пытаемся получить из таблицы admins если она существует
	rows, err := db.Query(`
		SELECT u.email 
		FROM users u 
		INNER JOIN admins a ON u.id = a.user_id 
		WHERE a.is_active = true
	`)
	if err == nil {
		defer rows.Close()
		
		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err == nil {
				adminEmails[email] = true
			}
		}
	}
	
	return adminEmails
}

func getMainUsers(db *sql.DB) ([]MainUser, error) {
	query := `
		SELECT id, email, name, password, created_at
		FROM users
		ORDER BY id
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []MainUser
	for rows.Next() {
		var user MainUser
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func clearAuthUsers(db *sql.DB) error {
	// Удаляем данные из каждой таблицы отдельно (не в транзакции)
	// чтобы ошибки отсутствующих таблиц не блокировали процесс
	
	// Сначала проверяем какие таблицы существуют
	existingTables := make(map[string]bool)
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'auth'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err == nil {
			existingTables[tableName] = true
		}
	}
	
	// Удаляем в правильном порядке для соблюдения foreign key constraints
	tables := []string{
		"user_roles",
		"refresh_tokens", 
		"sessions",
		"login_attempts",
		"oauth_states",
		"users",
	}
	
	for _, table := range tables {
		if existingTables[table] {
			query := fmt.Sprintf("DELETE FROM auth.%s", table)
			_, err = db.Exec(query)
			if err != nil {
				log.Printf("Warning: failed to clear auth.%s: %v", table, err)
			} else {
				log.Printf("Cleared table auth.%s", table)
			}
		}
	}
	
	return nil
}

func migrateUser(db *sql.DB, user MainUser, adminEmails map[string]bool) error {
	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Определяем provider
	provider := "local"
	if !user.Password.Valid || user.Password.String == "" {
		// Если нет пароля, скорее всего это OAuth пользователь
		if strings.Contains(user.Email, "gmail.com") {
			provider = "google"
		}
	}
	
	// Хешируем пароль если он есть
	var passwordHash sql.NullString
	if user.Password.Valid && user.Password.String != "" {
		// Проверяем, не является ли пароль уже хешем
		if strings.HasPrefix(user.Password.String, "$2a$") || strings.HasPrefix(user.Password.String, "$2b$") {
			passwordHash = user.Password
		} else {
			// Хешируем пароль
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password.String), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			passwordHash = sql.NullString{String: string(hash), Valid: true}
		}
	}
	
	// Вставляем пользователя с конкретным ID
	query := `
		INSERT INTO auth.users (
			id, email, email_normalized, name, password_hash, 
			provider, email_verified, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err = tx.Exec(query,
		user.ID,
		user.Email,
		strings.ToLower(user.Email),
		user.Name,
		passwordHash,
		provider,
		true, // email_verified
		true, // is_active
		user.CreatedAt,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	
	// Определяем роль
	roleID := userRoleID
	if adminEmails[user.Email] {
		roleID = adminRoleID
		log.Printf("  -> Assigning admin role to %s", user.Email)
	}
	
	// Назначаем роль
	_, err = tx.Exec(`
		INSERT INTO auth.user_roles (user_id, role_id, granted_at)
		VALUES ($1, $2, $3)
	`, user.ID, roleID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to assign role: %v", err)
	}
	
	return tx.Commit()
}

func resetSequence(db *sql.DB) {
	// Получаем максимальный ID
	var maxID int
	err := db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM auth.users").Scan(&maxID)
	if err != nil {
		log.Printf("Warning: failed to get max ID: %v", err)
		return
	}
	
	// Сбрасываем sequence
	query := fmt.Sprintf("ALTER SEQUENCE auth.users_id_seq RESTART WITH %d", maxID+1)
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Warning: failed to reset sequence: %v", err)
	} else {
		log.Printf("Reset sequence to start from %d", maxID+1)
	}
}

func backupAuthUsers(db *sql.DB) {
	// Создаем резервную таблицу
	_, err := db.Exec(`
		DROP TABLE IF EXISTS auth.users_backup_before_migration;
		CREATE TABLE auth.users_backup_before_migration AS 
		SELECT * FROM auth.users;
	`)
	if err != nil {
		log.Printf("Warning: failed to create backup table: %v", err)
	} else {
		log.Println("Backup table created: auth.users_backup_before_migration")
	}
}

func verifyMigration(db *sql.DB) {
	log.Println("\nVerifying migration...")
	
	// Проверяем количество пользователей
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM auth.users").Scan(&count)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}
	log.Printf("Total users in Auth Service: %d", count)
	
	// Проверяем админов
	rows, err := db.Query(`
		SELECT u.id, u.email, r.name
		FROM auth.users u
		JOIN auth.user_roles ur ON u.id = ur.user_id
		JOIN auth.roles r ON ur.role_id = r.id
		WHERE r.name = 'admin'
		ORDER BY u.id
	`)
	if err != nil {
		log.Printf("Failed to get admins: %v", err)
		return
	}
	defer rows.Close()
	
	log.Println("\nAdmins in Auth Service:")
	for rows.Next() {
		var id int
		var email, role string
		if err := rows.Scan(&id, &email, &role); err == nil {
			log.Printf("  - ID=%d, Email=%s, Role=%s", id, email, role)
		}
	}
	
	// Проверяем конкретных пользователей
	log.Println("\nChecking specific users:")
	checkUser(db, 3, "boxmail386@gmail.com")
	checkUser(db, 6, "4hash92@gmail.com")
}

func checkUser(db *sql.DB, expectedID int, expectedEmail string) {
	var id int
	var email string
	err := db.QueryRow("SELECT id, email FROM auth.users WHERE id = $1", expectedID).Scan(&id, &email)
	if err != nil {
		log.Printf("  ✗ User ID=%d not found: %v", expectedID, err)
	} else {
		if email == expectedEmail {
			log.Printf("  ✓ User ID=%d has correct email: %s", id, email)
		} else {
			log.Printf("  ✗ User ID=%d has wrong email: %s (expected %s)", id, email, expectedEmail)
		}
	}
}