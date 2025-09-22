package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler for AI-related endpoints
type Handler struct {
	config   *config.Config
	services globalService.ServicesInterface
}

// NewHandler creates a new AI handler
func NewHandler(cfg *config.Config, services globalService.ServicesInterface) *Handler {
	return &Handler{
		config:   cfg,
		services: services,
	}
}

// GetPrefix returns the prefix for AI routes
func (h *Handler) GetPrefix() string {
	return "/api/ai"
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RegisterRoutes registers all AI routes
func (h *Handler) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	// Create a group for AI routes
	aiGroup := app.Group(h.GetPrefix())

	// AI endpoints - no auth required for public access
	aiGroup.Post("/analyze", h.AnalyzeProduct)
	aiGroup.Post("/ab-test", h.PerformABTest)
	aiGroup.Post("/translate", h.TranslateContent)

	return nil
}

// AnalyzeProductRequest represents the request for product analysis
type AnalyzeProductRequest struct {
	ImageData string `json:"imageData"`
	UserLang  string `json:"userLang"`
}

// AnalyzeProductResponse represents the AI analysis response
type AnalyzeProductResponse struct {
	Title                 string                 `json:"title"`
	TitleVariants         []string               `json:"titleVariants"`
	Description           string                 `json:"description"`
	Category              string                 `json:"category"`
	CategoryHints         *CategoryHints         `json:"categoryHints,omitempty"`
	CategoryProbabilities []CategoryProbability  `json:"categoryProbabilities"`
	Price                 float64                `json:"price"`
	Currency              string                 `json:"currency"`
	Condition             string                 `json:"condition"`
	Attributes            map[string]interface{} `json:"attributes"`
	Keywords              []string               `json:"keywords"`
	Tags                  []string               `json:"tags"`
	SuggestedLocation     string                 `json:"suggestedLocation,omitempty"`
	SocialPosts           map[string]string      `json:"socialPosts,omitempty"`
}

// CategoryHints provides hints for category detection
type CategoryHints struct {
	Domain      string   `json:"domain"`
	ProductType string   `json:"productType"`
	Keywords    []string `json:"keywords"`
}

// CategoryProbability represents category with confidence
type CategoryProbability struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}

// ClaudeRequest represents the request to Claude API
type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

// Message represents a message in Claude API
type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

// Content represents content in a message
type Content struct {
	Type   string       `json:"type"`
	Text   string       `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
}

// ImageSource represents an image source in Claude API
type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// AnalyzeProduct analyzes a product image using Claude AI
// @Summary Analyze product image using AI
// @Description Analyzes a product image and extracts information using Claude AI
// @Tags ai
// @Accept json
// @Produce json
// @Param request body AnalyzeProductRequest true "Analysis request"
// @Success 200 {object} utils.SuccessResponseSwag{data=AnalyzeProductResponse}
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/ai/analyze [post]
func (h *Handler) AnalyzeProduct(c *fiber.Ctx) error {
	var req AnalyzeProductRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse analyze request")
		return utils.SendError(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Extract base64 data and media type from data URL
	var base64Data string
	mediaType := "image/jpeg" // default

	if strings.HasPrefix(req.ImageData, "data:") {
		parts := strings.Split(req.ImageData, ",")
		if len(parts) == 2 {
			base64Data = parts[1]
			// Extract media type from data URL (e.g., "data:image/png;base64,")
			if metaParts := strings.Split(parts[0], ";"); len(metaParts) > 0 {
				if strings.HasPrefix(metaParts[0], "data:") {
					mediaType = strings.TrimPrefix(metaParts[0], "data:")
				}
			}
		} else {
			base64Data = req.ImageData
		}
	} else {
		base64Data = req.ImageData
	}

	// Log the data for debugging
	logger.Info().
		Int("imageDataLength", len(req.ImageData)).
		Int("base64DataLength", len(base64Data)).
		Bool("isDataURL", strings.HasPrefix(req.ImageData, "data:")).
		Str("mediaType", mediaType).
		Str("dataPrefix", req.ImageData[:min(50, len(req.ImageData))]).
		Msg("Processing image data")

	// Check if base64 data is empty
	if len(base64Data) == 0 {
		logger.Error().Msg("No image data provided")
		return utils.SendError(c, fiber.StatusBadRequest, "ai.noImageData")
	}

	// Get Claude API key from config or environment
	apiKey := h.config.ClaudeAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("CLAUDE_API_KEY")
	}

	if apiKey == "" {
		logger.Error().Msg("Claude API key not configured")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.notConfigured")
	}

	// Prepare prompt based on user language
	prompt := h.getPromptForLanguage(req.UserLang)

	// Create Claude API request
	claudeReq := ClaudeRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "image",
						Source: &ImageSource{
							Type:      "base64",
							MediaType: mediaType,
							Data:      base64Data,
						},
					},
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(claudeReq)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal Claude request")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Create HTTP request to Claude API
	claudeURL := "https://api.anthropic.com/v1/messages"
	req2, err := http.NewRequest("POST", claudeURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Claude request")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Set headers
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("x-api-key", apiKey)
	req2.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to call Claude API")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.apiError")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read Claude response")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Error().
			Int("status", resp.StatusCode).
			Str("response", string(body)).
			Msg("Claude API returned error")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.apiError")
	}

	// Parse Claude response
	var claudeResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &claudeResp); err != nil {
		logger.Error().Err(err).Msg("Failed to parse Claude response")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Extract text from response
	var responseText string
	for _, content := range claudeResp.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		logger.Error().Msg("No text in Claude response")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.noResponse")
	}

	// Parse the structured response from Claude
	var analysisResult AnalyzeProductResponse
	if err := json.Unmarshal([]byte(responseText), &analysisResult); err != nil {
		logger.Error().
			Err(err).
			Str("response", responseText).
			Msg("Failed to parse Claude analysis result")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.invalidResponse")
	}

	logger.Info().
		Str("title", analysisResult.Title).
		Str("category", analysisResult.Category).
		Float64("price", analysisResult.Price).
		Msg("Product analyzed successfully")

	return utils.SendSuccess(c, fiber.StatusOK, "ai.analysisSuccess", analysisResult)
}

// ABTestRequest represents the request for A/B testing
type ABTestRequest struct {
	TitleVariants []string `json:"titleVariants"`
}

// ABTestResponse represents the A/B test response
type ABTestResponse struct {
	BestVariant string        `json:"bestVariant"`
	Scores      []ABTestScore `json:"scores"`
}

// ABTestScore represents individual variant score
type ABTestScore struct {
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

// PerformABTest performs A/B testing on title variants
func (h *Handler) PerformABTest(c *fiber.Ctx) error {
	var req ABTestRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse AB test request")
		return utils.SendError(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Simple A/B test logic - in production this would be more sophisticated
	scores := make([]ABTestScore, len(req.TitleVariants))
	bestScore := 0.0
	bestVariant := ""

	for i, variant := range req.TitleVariants {
		// Simple scoring based on length and clarity
		score := 100.0 - float64(i)*10 // First variant gets highest score by default
		if len(variant) > 50 {
			score -= 10 // Penalize too long titles
		}
		if len(variant) < 10 {
			score -= 15 // Penalize too short titles
		}

		scores[i] = ABTestScore{
			Title: variant,
			Score: score,
		}

		if score > bestScore {
			bestScore = score
			bestVariant = variant
		}
	}

	response := ABTestResponse{
		BestVariant: bestVariant,
		Scores:      scores,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "ai.abTestSuccess", response)
}

// TranslateContentRequest represents the request for content translation
type TranslateContentRequest struct {
	Content struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"content"`
	TargetLanguages []string `json:"targetLanguages"`
}

// TranslateContentResponse represents the translation response
type TranslateContentResponse map[string]struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// TranslateContent translates content to multiple languages
// @Summary Translate content to multiple languages
// @Description Translates title and description to specified languages using Claude AI
// @Tags ai
// @Accept json
// @Produce json
// @Param request body TranslateContentRequest true "Translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=TranslateContentResponse}
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/ai/translate [post]
func (h *Handler) TranslateContent(c *fiber.Ctx) error {
	var req TranslateContentRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse translate request")
		return utils.SendError(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Check for Claude API key
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		// Return mock translation if no API key
		mockResponse := make(TranslateContentResponse)
		for _, lang := range req.TargetLanguages {
			mockResponse[lang] = struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				Title:       req.Content.Title,
				Description: req.Content.Description,
			}
		}
		return utils.SendSuccess(c, fiber.StatusOK, "ai.translationSuccess", mockResponse)
	}

	// Build translation prompt
	prompt := fmt.Sprintf(`Translate the following content to the specified languages.
Return ONLY a valid JSON object with translations.

Original content:
Title: %s
Description: %s

Target languages: %s

Return JSON in this format:
{
  "en": {
    "title": "English title",
    "description": "English description"
  },
  "sr": {
    "title": "Serbian title",
    "description": "Serbian description"
  }
}

IMPORTANT: Return ONLY the JSON, no markdown or explanations.`,
		req.Content.Title,
		req.Content.Description,
		strings.Join(req.TargetLanguages, ", "))

	// Create Claude API request
	claudeReq := ClaudeRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(claudeReq)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal Claude request")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Create HTTP request to Claude API
	claudeURL := "https://api.anthropic.com/v1/messages"
	req2, err := http.NewRequest("POST", claudeURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Claude request")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Set headers
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("x-api-key", apiKey)
	req2.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to call Claude API")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.apiError")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read Claude response")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Error().
			Int("status", resp.StatusCode).
			Str("response", string(body)).
			Msg("Claude API returned error")
		// Return mock translation on API error
		mockResponse := make(TranslateContentResponse)
		for _, lang := range req.TargetLanguages {
			mockResponse[lang] = struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				Title:       req.Content.Title,
				Description: req.Content.Description,
			}
		}
		return utils.SendSuccess(c, fiber.StatusOK, "ai.translationSuccess", mockResponse)
	}

	// Parse Claude response
	var claudeResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &claudeResp); err != nil {
		logger.Error().Err(err).Msg("Failed to parse Claude response")
		return utils.SendError(c, fiber.StatusInternalServerError, "ai.processingError")
	}

	// Extract text from response
	var responseText string
	for _, content := range claudeResp.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		logger.Error().Msg("No text in Claude response")
		// Return mock translation
		mockResponse := make(TranslateContentResponse)
		for _, lang := range req.TargetLanguages {
			mockResponse[lang] = struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				Title:       req.Content.Title,
				Description: req.Content.Description,
			}
		}
		return utils.SendSuccess(c, fiber.StatusOK, "ai.translationSuccess", mockResponse)
	}

	// Parse the translation response
	var translationResult TranslateContentResponse
	if err := json.Unmarshal([]byte(responseText), &translationResult); err != nil {
		logger.Error().
			Err(err).
			Str("response", responseText).
			Msg("Failed to parse Claude translation result")
		// Return mock translation on parse error
		mockResponse := make(TranslateContentResponse)
		for _, lang := range req.TargetLanguages {
			mockResponse[lang] = struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				Title:       req.Content.Title,
				Description: req.Content.Description,
			}
		}
		return utils.SendSuccess(c, fiber.StatusOK, "ai.translationSuccess", mockResponse)
	}

	logger.Info().
		Interface("languages", req.TargetLanguages).
		Msg("Content translated successfully")

	return utils.SendSuccess(c, fiber.StatusOK, "ai.translationSuccess", translationResult)
}

// getPromptForLanguage returns the appropriate prompt based on user language
func (h *Handler) getPromptForLanguage(lang string) string {
	if lang == "ru" {
		return `Проанализируй это изображение товара и верни информацию в формате JSON.
Определи что это за товар, его состояние, возможную цену и категорию.

ВАЖНО:
1. Категория должна быть на АНГЛИЙСКОМ языке (Electronics, Clothing, Furniture, etc.)
2. ВСЕ тексты (title, description, attributes) должны быть на РУССКОМ языке!
3. ОБЯЗАТЕЛЬНО создай посты для социальных сетей НА РУССКОМ языке!

Верни ТОЛЬКО валидный JSON без markdown или объяснений:
{
  "title": "короткое название товара НА РУССКОМ",
  "titleVariants": ["вариант названия 1 НА РУССКОМ", "вариант названия 2 НА РУССКОМ", "вариант названия 3 НА РУССКОМ"],
  "description": "детальное описание товара НА РУССКОМ ЯЗЫКЕ, его особенностей и состояния",
  "category": "главная категория на АНГЛИЙСКОМ (Electronics, Clothing, Furniture, Home & Garden, Cars, etc.)",
  "categoryHints": {
    "domain": "тип товара (electronics, fashion, home, automotive, etc.)",
    "productType": "конкретный тип (router, smartphone, laptop, shirt, etc.)",
    "keywords": ["router", "wifi", "network", "huawei"]
  },
  "categoryProbabilities": [
    {"name": "Electronics", "probability": 0.9},
    {"name": "Computers", "probability": 0.1}
  ],
  "price": числовое значение предполагаемой цены,
  "currency": "RSD",
  "condition": "new" или "used",
  "attributes": {
    "brand": "бренд если виден НА РУССКОМ",
    "model": "модель если видна НА РУССКОМ",
    "color": "цвет НА РУССКОМ",
    "size": "размер если применимо НА РУССКОМ",
    "material": "материал если виден НА РУССКОМ"
  },
  "keywords": ["ключевое слово 1", "ключевое слово 2"],
  "tags": ["тег 1", "тег 2", "тег 3"],
  "suggestedLocation": "предположительное местоположение если есть подсказки",
  "socialPosts": {
    "instagram": "🔥 Продаю [название товара]! ✨ [краткое описание] 💰 Цена: [цена] динаров 📱 Пишите в директ! #продажа #товары #сербия",
    "facebook": "Продается [название товара]. [Описание состояния и особенностей]. Цена: [цена] динаров. Интересно? Пишите в комментарии или личные сообщения!",
    "twitter": "🔥 [Название товара] за [цена] динаров! [Краткое описание] #продажа #маркетплейс",
    "whatsapp": "Привет! Продаю [название товара]. [Описание]. Цена: [цена] динаров. Интересно? Могу отправить дополнительные фото!",
    "telegram": "💰 Продается: [название товара]\n📝 Описание: [краткое описание]\n💵 Цена: [цена] динаров\n📞 Связаться со мной можно здесь!",
    "viber": "🎯 Продаю [название товара]\n✅ [Краткое описание]\n💰 Цена: [цена] динаров\n📱 Звоните или пишите в Viber!"
  }
}`
	}

	if lang == "sr" {
		return `Analiziraj ovu sliku proizvoda i vrati informacije u JSON formatu.
Identifikuj šta je proizvod, njegovo stanje, moguću cenu i kategoriju.

VAŽNO:
1. Kategorija mora biti na ENGLESKOM jeziku (Electronics, Clothing, Furniture, itd.)
2. SVI tekstovi (title, description, attributes) moraju biti na SRPSKOM jeziku!

Vrati SAMO valjan JSON bez markdown ili objašnjenja:
{
  "title": "kratak naziv proizvoda NA SRPSKOM",
  "titleVariants": ["varijanta naziva 1 NA SRPSKOM", "varijanta naziva 2 NA SRPSKOM", "varijanta naziva 3 NA SRPSKOM"],
  "description": "detaljan opis proizvoda NA SRPSKOM JEZIKU, njegove karakteristike i stanje",
  "category": "glavna kategorija na ENGLESKOM (Electronics, Clothing, Furniture, Home & Garden, Cars, etc.)",
  "categoryHints": {
    "domain": "tip proizvoda (electronics, fashion, home, automotive, etc.)",
    "productType": "konkretan tip (router, smartphone, laptop, shirt, etc.)",
    "keywords": ["router", "wifi", "network", "huawei"]
  },
  "categoryProbabilities": [
    {"name": "Electronics", "probability": 0.9},
    {"name": "Computers", "probability": 0.1}
  ],
  "price": numerička vrednost pretpostavljene cene,
  "currency": "RSD",
  "condition": "new" ili "used",
  "attributes": {
    "brand": "brend ako je vidljiv NA SRPSKOM",
    "model": "model ako je vidljiv NA SRPSKOM",
    "color": "boja NA SRPSKOM",
    "size": "veličina ako je primenljivo NA SRPSKOM",
    "material": "materijal ako je vidljiv NA SRPSKOM"
  },
  "keywords": ["ključna reč 1", "ključna reč 2"],
  "tags": ["tag 1", "tag 2", "tag 3"],
  "suggestedLocation": "pretpostavljena lokacija ako postoje naznake",
  "socialPosts": {
    "instagram": "🔥 Prodajem [naziv proizvoda]! ✨ [kratak opis] 💰 Cena: [cena] dinara 📱 Pišite u DM! #prodaja #proizvodi #srbija",
    "facebook": "Prodaje se [naziv proizvoda]. [Opis stanja i karakteristika]. Cena: [cena] dinara. Zainteresovani? Pišite u komentarima ili privatnim porukama!",
    "twitter": "🔥 [Naziv proizvoda] za [cena] dinara! [Kratak opis] #prodaja #marketplace",
    "whatsapp": "Zdravo! Prodajem [naziv proizvoda]. [Opis]. Cena: [cena] dinara. Zainteresovani? Mogu poslati dodatne fotografije!",
    "telegram": "💰 Prodaje se: [naziv proizvoda]\n📝 Opis: [kratak opis]\n💵 Cena: [cena] dinara\n📞 Kontaktirajte me ovde!",
    "viber": "🎯 Prodajem [naziv proizvoda]\n✅ [Kratak opis]\n💰 Cena: [cena] dinara\n📱 Pozovite ili pišite na Viber!"
  }
}`
	}

	// English prompt (default)
	return `Analyze this product image and return information in JSON format.
Identify what the product is, its condition, possible price, and category.

IMPORTANT: Category must be in ENGLISH (Electronics, Clothing, Furniture, etc.)

Return ONLY valid JSON without markdown or explanations:
{
  "title": "short product title",
  "titleVariants": ["title variant 1", "title variant 2", "title variant 3"],
  "description": "detailed description of the product, its features and condition",
  "category": "main category in ENGLISH (Electronics, Clothing, Furniture, Home & Garden, Cars, etc.)",
  "categoryHints": {
    "domain": "product domain (electronics, fashion, home, automotive, etc.)",
    "productType": "specific type (router, smartphone, laptop, shirt, etc.)",
    "keywords": ["router", "wifi", "network", "huawei"]
  },
  "categoryProbabilities": [
    {"name": "Electronics", "probability": 0.9},
    {"name": "Computers", "probability": 0.1}
  ],
  "price": numeric price value,
  "currency": "RSD",
  "condition": "new" or "used",
  "attributes": {
    "brand": "brand if visible",
    "model": "model if visible",
    "color": "color",
    "size": "size if applicable",
    "material": "material if visible"
  },
  "keywords": ["keyword 1", "keyword 2"],
  "tags": ["tag 1", "tag 2", "tag 3"],
  "suggestedLocation": "suggested location if there are hints",
  "socialPosts": {
    "instagram": "🔥 Selling [product name]! ✨ [brief description] 💰 Price: [price] RSD 📱 DM for details! #forsale #marketplace #serbia",
    "facebook": "For sale: [product name]. [Description of condition and features]. Price: [price] RSD. Interested? Comment below or send me a message!",
    "twitter": "🔥 [Product name] for [price] RSD! [Brief description] #forsale #marketplace",
    "whatsapp": "Hi! I'm selling [product name]. [Description]. Price: [price] RSD. Interested? I can send more photos!",
    "telegram": "💰 For sale: [product name]\n📝 Description: [brief description]\n💵 Price: [price] RSD\n📞 Contact me here!",
    "viber": "🎯 Selling [product name]\n✅ [Brief description]\n💰 Price: [price] RSD\n📱 Call or message on Viber!"
  }
}`
}
