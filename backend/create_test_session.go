package main

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Создаем тестовую сессию через добавление записи в сессию
	client := &http.Client{}
	
	// Вызываем test endpoint, который создаст сессию
	data := url.Values{}
	data.Set("test_user_id", "3")
	
	resp, err := client.PostForm("http://localhost:3000/api/test/create-session", data)
	if err \!= nil {
		log.Fatal("Error:", err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err \!= nil {
		log.Fatal("Error reading response:", err)
	}
	
	fmt.Printf("Response: %s\n", string(body))
}
EOF < /dev/null
