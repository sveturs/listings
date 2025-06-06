'use client';

import { useState, useEffect } from 'react';
import Lottie from 'lottie-react';

interface AnimatedEmojiProps {
  emoji: string;
  size?: number;
}

// Маппинг эмодзи на URL анимаций Google Noto Animated Emoji
// Формат URL: https://fonts.gstatic.com/s/e/notoemoji/latest/{unicode}/lottie.json
const emojiAnimations: Record<string, string> = {
  '😀': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f600/lottie.json',
  '😊': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f60a/lottie.json',
  '❤️': 'https://fonts.gstatic.com/s/e/notoemoji/latest/2764_fe0f/lottie.json',
  '🔥': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f525/lottie.json',
  '👍': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f44d/lottie.json',
  '😂': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f602/lottie.json',
  '🎉': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f389/lottie.json',
  '💕': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f495/lottie.json',
  '🥰': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f970/lottie.json',
  '😍': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f60d/lottie.json',
  '🤗': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f917/lottie.json',
  '😘': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f618/lottie.json',
  '🙂': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f642/lottie.json',
  '😎': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f60e/lottie.json',
  '😭': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f62d/lottie.json',
  '😢': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f622/lottie.json',
  '😅': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f605/lottie.json',
  '🤔': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f914/lottie.json',
  '😱': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f631/lottie.json',
  '🤯': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f92f/lottie.json',
  '😴': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f634/lottie.json',
  '🤩': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f929/lottie.json',
  '🥳': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f973/lottie.json',
  '🙏': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f64f/lottie.json',
  '👌': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f44c/lottie.json',
  '✌️': 'https://fonts.gstatic.com/s/e/notoemoji/latest/270c_fe0f/lottie.json',
  '🤞': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f91e/lottie.json',
  '💪': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f4aa/lottie.json',
  '👏': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f44f/lottie.json',
  '🙌': 'https://fonts.gstatic.com/s/e/notoemoji/latest/1f64c/lottie.json',
};

// Глобальный кэш для загруженных анимаций
const animationCache: Record<string, object> = {};

export default function AnimatedEmoji({
  emoji,
  size = 96,
}: AnimatedEmojiProps) {
  const [animationData, setAnimationData] = useState<object | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    const loadAnimation = async () => {
      // Проверяем кэш
      if (animationCache[emoji]) {
        setAnimationData(animationCache[emoji]);
        setLoading(false);
        return;
      }

      const animationUrl = emojiAnimations[emoji];

      if (!animationUrl) {
        // Если анимация не найдена, показываем обычный эмодзи
        setError(true);
        setLoading(false);
        return;
      }

      try {
        const response = await fetch(animationUrl);
        const data = await response.json();
        // Сохраняем в кэш
        animationCache[emoji] = data;
        setAnimationData(data);
        setLoading(false);
      } catch (err) {
        console.error('Failed to load emoji animation:', err);
        setError(true);
        setLoading(false);
      }
    };

    loadAnimation();
  }, [emoji]);

  // Если анимация не найдена или произошла ошибка, показываем обычный эмодзи
  if (error || !animationData) {
    return (
      <span
        className="text-6xl leading-none inline-block"
        style={{ fontSize: size }}
      >
        {emoji}
      </span>
    );
  }

  // Показываем загрузку
  if (loading) {
    return (
      <div
        className="inline-block animate-pulse bg-base-300 rounded-lg"
        style={{ width: size, height: size }}
      />
    );
  }

  return (
    <div className="inline-block" style={{ width: size, height: size }}>
      <Lottie
        animationData={animationData}
        loop={true}
        autoplay={true}
        style={{ width: '100%', height: '100%' }}
      />
    </div>
  );
}
