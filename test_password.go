package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "EmailEmail"
	existingHash := "$2a$10$C3mkAl/nqP9ZfJ5XijW3JuedWab/FAvaWwaP9L/dAmj01jsQsxdJa"
	
	// Проверяем существующий хеш
	err := bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(password))
	if err == nil {
		fmt.Println("Existing hash is VALID for password:", password)
	} else {
		fmt.Println("Existing hash is INVALID:", err)
	}
	
	// Генерируем новый хеш с default cost
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nNew hash: %s\n", string(newHash))
	
	// Проверяем новый хеш
	err = bcrypt.CompareHashAndPassword(newHash, []byte(password))
	if err == nil {
		fmt.Println("New hash verification: PASSED")
	} else {
		fmt.Println("New hash verification: FAILED")
	}
}