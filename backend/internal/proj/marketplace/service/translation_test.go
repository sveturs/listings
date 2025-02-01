// /data/proj/hostel-booking-system/backend/internal/proj/marketplace/service/translation_test.go
package service

import (
    "context"
    "testing"
)

func TestTranslation(t *testing.T) {
    service, err := NewTranslationService()
    if err != nil {
        t.Fatalf("Failed to create translation service: %v", err)
    }

    tests := []struct {
        name           string
        text          string
        sourceLanguage string
        targetLanguage string
        wantError     bool
    }{
        {
            name:           "Serbian to English",
            text:          "Здраво свете",
            sourceLanguage: "sr",
            targetLanguage: "en",
            wantError:     false,
        },
        {
            name:           "English to Russian",
            text:          "Hello world",
            sourceLanguage: "en",
            targetLanguage: "ru",
            wantError:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            translated, err := service.Translate(context.Background(), tt.text, tt.sourceLanguage, tt.targetLanguage)
            if (err != nil) != tt.wantError {
                t.Errorf("Translate() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if !tt.wantError && translated == "" {
                t.Error("Translate() returned empty string")
            }
        })
    }
}