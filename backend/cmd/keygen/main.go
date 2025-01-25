// cmd/keygen/main.go
package main

import (
    "fmt"
    "github.com/SherClockHolmes/webpush-go"
)

func main() {
    privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Public Key: %s\nPrivate Key: %s\n", publicKey, privateKey)
}