package service

import "context"

type TranslationServiceInterface interface {
     Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error)
}