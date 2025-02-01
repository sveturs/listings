// backend/internal/proj/marketplace/service/translation.go
package service

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "log"
    "sync"
)

const GOOGLE_TRANSLATE_API_URL = "https://translation.googleapis.com/language/translate/v2"

type TranslationService struct {
    apiKey string
    cache  *sync.Map
}

type googleTranslateRequest struct {
    Q        string `json:"q"`
    Source   string `json:"source"`
    Target   string `json:"target"`
    Format   string `json:"format"`
}

type googleTranslateResponse struct {
    Data struct {
        Translations []struct {
            TranslatedText string `json:"translatedText"`
        } `json:"translations"`
    } `json:"data"`
}

func NewTranslationService(apiKey string) (*TranslationService, error) {
    if apiKey == "" {
        return nil, fmt.Errorf("Google Translate API key is required")
    }
    return &TranslationService{
        apiKey: apiKey,
        cache:  &sync.Map{},
    }, nil
}

func (s *TranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("%s:%s:%s:%s", text, sourceLanguage, targetLanguage, text)
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.(string), nil
    }

    reqBody := googleTranslateRequest{
        Q:      text,
        Source: sourceLanguage,
        Target: targetLanguage,
        Format: "text",
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %w", err)
    }

    url := fmt.Sprintf("%s?key=%s", GOOGLE_TRANSLATE_API_URL, s.apiKey)
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("translation API returned status %d", resp.StatusCode)
    }

    var result googleTranslateResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("error decoding response: %w", err)
    }

    if len(result.Data.Translations) == 0 {
        return "", fmt.Errorf("no translation returned")
    }

    translatedText := result.Data.Translations[0].TranslatedText

    // Cache the result
    s.cache.Store(cacheKey, translatedText)

    return translatedText, nil
}

// TranslateEntityFields translates specific fields of an entity
func (s *TranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
    results := make(map[string]map[string]string)
    
    for _, targetLang := range targetLanguages {
        if targetLang == sourceLanguage {
            continue
        }
        
        results[targetLang] = make(map[string]string)
        
        for fieldName, text := range fields {
            // Skip empty fields
            if text == "" {
                continue
            }
            
            translated, err := s.Translate(ctx, text, sourceLanguage, targetLang)
            if err != nil {
                log.Printf("Error translating field %s to %s: %v", fieldName, targetLang, err)
                continue
            }
            
            results[targetLang][fieldName] = translated
        }
    }
    
    return results, nil
}