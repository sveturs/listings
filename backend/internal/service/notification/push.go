// В backend/internal/service/notification/push.go

package notification
import "log"
import "github.com/SherClockHolmes/webpush-go"

// GenerateVAPIDKeys генерирует ключи для VAPID

func GenerateVAPIDKeys() (publicKey, privateKey string) {
    privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
    if err != nil {
        log.Fatal(err)
    }
    return publicKey, privateKey
}