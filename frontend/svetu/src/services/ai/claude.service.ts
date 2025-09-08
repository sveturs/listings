import configManager from '@/config';

interface ProductAnalysis {
  title: string;
  titleVariants: string[];
  description: string;
  categoryHints?: {
    domain: string;
    productType: string;
    keywords: string[];
  };
  category: string;
  categoryProbabilities: { name: string; probability: number }[];
  price: string;
  priceRange: { min: number; max: number };
  attributes: Record<string, string>;
  tags: string[];
  suggestedPhotos: string[];
  translations: Record<string, { title: string; description: string }>;
  socialPosts: Record<string, string>;
  location?: {
    city?: string;
    region?: string;
    suggestedLocation?: string;
  };
  condition?: 'new' | 'used' | 'refurbished';
  insights?: Record<
    string,
    {
      demand: string;
      audience: string;
      recommendations: string;
    }
  >;
}

export class ClaudeAIService {
  private apiKey: string;
  private apiUrl = 'https://api.anthropic.com/v1/messages';

  constructor() {
    this.apiKey = configManager.getConfig().claudeApiKey || '';
    console.log(
      'Claude API:',
      this.apiKey
        ? 'Key configured, starts with: ' + this.apiKey.substring(0, 10) + '...'
        : 'Using mock mode (no key configured)'
    );
  }

  async analyzeProduct(
    imageBase64: string,
    userLanguage: string = 'ru'
  ): Promise<ProductAnalysis> {
    try {
      console.log('Sending product analysis request via API route...');
      console.log('Image data length:', imageBase64.length);
      console.log('User language:', userLanguage);

      // Use our API route instead of calling Claude directly
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(`${apiUrl}/api/ai/analyze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          imageBase64,
          userLanguage,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        console.error('Claude API error:', response.status, errorData);
        throw new Error(`Claude API error: ${response.status} - ${errorData}`);
      }

      const analysis = await response.json();
      console.log('Product analysis completed:', analysis.title);

      return analysis;
    } catch (error) {
      console.error('Claude AI analysis error:', error);
      throw error;
    }
  }

  async generateOptimizedDescription(
    title: string,
    category: string,
    attributes: Record<string, string>
  ): Promise<string> {
    try {
      // Use our API route instead of calling Claude directly
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(`${apiUrl}/api/ai/description`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title,
          category,
          attributes,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        console.error('Description API error:', response.status, errorData);
        throw new Error(`Description API error: ${response.status}`);
      }

      const data = await response.json();
      return data.description || '';
    } catch (error) {
      console.error('Claude description generation error:', error);
      throw error;
    }
  }

  async performABTesting(titleVariants: string[]): Promise<{
    bestVariant: string;
    scores: { title: string; score: number }[];
  }> {
    try {
      // Use our API route instead of calling Claude directly
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(`${apiUrl}/api/ai/ab-test`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          titleVariants,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        console.error('A/B testing API error:', response.status, errorData);
        throw new Error(`A/B testing API error: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Claude A/B testing error:', error);
      // Fallback: return first variant
      return {
        bestVariant: titleVariants[0],
        scores: titleVariants.map((title, i) => ({
          title,
          score: 90 - i * 10,
        })),
      };
    }
  }

  async translateContent(
    content: { title: string; description: string },
    targetLanguages: string[] = ['en', 'ru', 'sr']
  ): Promise<Record<string, { title: string; description: string }>> {
    try {
      console.log('Translating content to languages:', targetLanguages);

      // Use our API route for translation
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(`${apiUrl}/api/ai/translate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content,
          targetLanguages,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        console.error('Translation API error:', response.status, errorData);
        throw new Error(`Translation API error: ${response.status}`);
      }

      const translations = await response.json();
      console.log('Translations completed');

      return translations;
    } catch (error) {
      console.error('Claude translation error:', error);

      // Fallback: вернуть исходный контент для всех языков
      const fallbackTranslations: Record<
        string,
        { title: string; description: string }
      > = {};
      targetLanguages.forEach((lang) => {
        fallbackTranslations[lang] = {
          title: content.title,
          description: content.description,
        };
      });

      return fallbackTranslations;
    }
  }
}

export const claudeAI = new ClaudeAIService();
