import { config } from '@/config';

interface ClaudeResponse {
  content: {
    text: string;
  }[];
}

interface ProductAnalysis {
  title: string;
  titleVariants: string[];
  description: string;
  category: string;
  categoryProbabilities: { name: string; probability: number }[];
  price: string;
  priceRange: { min: number; max: number };
  attributes: Record<string, string>;
  tags: string[];
  suggestedPhotos: string[];
  translations: Record<string, { title: string; description: string }>;
  socialPosts: Record<string, string>;
}

export class ClaudeAIService {
  private apiKey: string;
  private apiUrl = 'https://api.anthropic.com/v1/messages';

  constructor() {
    this.apiKey = config.claudeApiKey || '';
    console.log('Claude API key configured:', !!this.apiKey, this.apiKey ? 'Key starts with: ' + this.apiKey.substring(0, 10) + '...' : 'No key');
  }

  async analyzeProduct(imageBase64: string): Promise<ProductAnalysis> {
    try {
      console.log('Sending product analysis request via API route...');
      console.log('Image data length:', imageBase64.length);
      
      // Use our API route instead of calling Claude directly
      const response = await fetch('/api/ai/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          imageBase64,
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
      const response = await fetch('/api/ai/description', {
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
      const response = await fetch('/api/ai/ab-test', {
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
}

export const claudeAI = new ClaudeAIService();
