import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  try {
    const { content, targetLanguages } = await request.json();

    if (!content || !content.title || !content.description) {
      return NextResponse.json(
        { error: 'Content with title and description is required' },
        { status: 400 }
      );
    }

    if (
      !targetLanguages ||
      !Array.isArray(targetLanguages) ||
      targetLanguages.length === 0
    ) {
      return NextResponse.json(
        { error: 'Target languages array is required' },
        { status: 400 }
      );
    }

    if (!CLAUDE_API_KEY) {
      return NextResponse.json(
        { error: 'Claude API key is not configured' },
        { status: 500 }
      );
    }

    console.log('Translating content to languages:', targetLanguages);

    const prompt = `You are a professional translator for an online marketplace. Translate the following product listing into the specified languages.

Original content:
Title: ${content.title}
Description: ${content.description}

Please translate this content into the following languages: ${targetLanguages.join(', ')}.

Important guidelines:
1. Maintain the marketing appeal and SEO optimization
2. Use natural, native-sounding language for each translation
3. Adapt cultural references when necessary
4. Keep the same level of detail and information
5. For Serbian (sr), use Cyrillic script
6. For Russian (ru), use Cyrillic script
7. For English (en), ensure proper grammar and natural flow

Return the translations in this exact JSON format:
{
  "en": {
    "title": "English title here",
    "description": "English description here"
  },
  "ru": {
    "title": "Русский заголовок здесь",
    "description": "Русское описание здесь"
  },
  "sr": {
    "title": "Српски наслов овде",
    "description": "Српски опис овде"
  }
}

Only include the languages that were requested. Return ONLY the JSON object, no additional text.`;

    const response = await fetch(CLAUDE_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': CLAUDE_API_KEY,
        'anthropic-version': '2023-06-01',
      },
      body: JSON.stringify({
        model: 'claude-3-5-sonnet-20241022',
        max_tokens: 2000,
        messages: [
          {
            role: 'user',
            content: prompt,
          },
        ],
        temperature: 0.3, // Lower temperature for more consistent translations
      }),
    });

    if (!response.ok) {
      const errorData = await response.text();
      console.error('Claude API error:', response.status, errorData);
      return NextResponse.json(
        { error: `Claude API error: ${response.status}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    const responseText = data.content[0]?.text || '';

    // Try to parse the JSON response
    let translations;
    try {
      // Find JSON in the response (Claude might add some text around it)
      const jsonMatch = responseText.match(/\{[\s\S]*\}/);
      if (jsonMatch) {
        translations = JSON.parse(jsonMatch[0]);
      } else {
        throw new Error('No JSON found in response');
      }
    } catch {
      console.error('Failed to parse Claude response:', responseText);
      throw new Error('Invalid response format from Claude');
    }

    // Validate that we got translations for all requested languages
    const missingLanguages = targetLanguages.filter(
      (lang) => !translations[lang]
    );
    if (missingLanguages.length > 0) {
      console.warn('Missing translations for languages:', missingLanguages);

      // Add fallback for missing languages
      missingLanguages.forEach((lang) => {
        translations[lang] = {
          title: content.title,
          description: content.description,
        };
      });
    }

    console.log('Translation completed successfully');
    return NextResponse.json(translations);
  } catch (error) {
    console.error('Translation error:', error);

    if (error instanceof Error && error.message.includes('api_key')) {
      return NextResponse.json(
        { error: 'Claude API key not configured' },
        { status: 500 }
      );
    }

    return NextResponse.json(
      { error: 'Failed to translate content' },
      { status: 500 }
    );
  }
}
