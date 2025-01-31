package service

import (
    "context"
    "cloud.google.com/go/translate"
    "golang.org/x/text/language"
)

type TranslationService struct {
    client *translate.Client
}

func NewTranslationService() (*TranslationService, error) {
    ctx := context.Background()
    client, err := translate.NewClient(ctx)
    if err != nil {
        return nil, err
    }
    return &TranslationService{client: client}, nil
}

// Изменяем имя метода с TranslateText на Translate
func (s *TranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
    if text == "" {
        return "", nil
    }

    source, err := language.Parse(sourceLanguage)
    if err != nil {
        return "", err
    }

    target, err := language.Parse(targetLanguage)
    if err != nil {
        return "", err
    }

    translations, err := s.client.Translate(ctx, []string{text}, target, &translate.Options{
        Source: source,
    })

    if err != nil {
        return "", err
    }

    if len(translations) == 0 {
        return "", nil
    }

    return translations[0].Text, nil
}