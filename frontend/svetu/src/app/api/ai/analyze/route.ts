import { NextRequest, NextResponse } from 'next/server';
import { getAnalysisPrompt } from './prompts';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { imageBase64, userLanguage = 'ru' } = body;

    if (!imageBase64 || imageBase64.length === 0) {
      console.error('No image data provided to AI analyze API');
      return NextResponse.json(
        { error: 'ai.noImageData', success: false },
        { status: 400 }
      );
    }

    // Validate base64 format
    if (imageBase64.length < 100) {
      console.error('Image data too short, likely invalid:', imageBase64.length);
      return NextResponse.json(
        { error: 'ai.invalidImageData', success: false },
        { status: 400 }
      );
    }

    if (!CLAUDE_API_KEY) {
      return NextResponse.json(
        { error: 'Claude API key not configured' },
        { status: 500 }
      );
    }

    console.log('Calling Claude API...');

    const response = await fetch(CLAUDE_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': CLAUDE_API_KEY,
        'anthropic-version': '2023-06-01',
      },
      body: JSON.stringify({
        model: 'claude-3-5-sonnet-20241022',
        max_tokens: 4096,
        messages: [
          {
            role: 'user',
            content: [
              {
                type: 'image',
                source: {
                  type: 'base64',
                  media_type: 'image/jpeg',
                  data: imageBase64,
                },
              },
              {
                type: 'text',
                text: getAnalysisPrompt(userLanguage),
              },
            ],
          },
        ],
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Claude API error:', response.status, errorText);

      // Возвращаем ошибку API
      console.error('Claude API failed with status:', response.status);

      return NextResponse.json(
        { error: `Claude API error: ${response.status}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    const content = data.content?.[0]?.text;

    if (!content) {
      return NextResponse.json(
        { error: 'No content in Claude response' },
        { status: 500 }
      );
    }

    // Simple JSON parsing
    try {
      // Try to extract JSON from the content
      let jsonStr = content;

      // If wrapped in code block, extract it
      const codeBlockMatch = content.match(/```(?:json)?\s*([\s\S]*?)\s*```/);
      if (codeBlockMatch) {
        jsonStr = codeBlockMatch[1];
      }

      // Remove any text before first { and after last }
      const firstBrace = jsonStr.indexOf('{');
      const lastBrace = jsonStr.lastIndexOf('}');
      if (firstBrace !== -1 && lastBrace !== -1) {
        jsonStr = jsonStr.substring(firstBrace, lastBrace + 1);
      }

      const analysis = JSON.parse(jsonStr);
      return NextResponse.json(analysis);
    } catch (parseError) {
      console.error('Failed to parse Claude response:', parseError);
      console.error('Content:', content.substring(0, 500));

      return NextResponse.json(
        { error: 'Failed to parse AI response' },
        { status: 500 }
      );
    }
  } catch (error) {
    console.error('API route error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
