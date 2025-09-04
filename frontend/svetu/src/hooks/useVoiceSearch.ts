// Voice Search Hook для мобильной оптимизации
// День 25: Voice Search Integration

import { useState, useEffect, useRef, useCallback } from 'react';

interface VoiceSearchOptions {
  language?: string;
  continuous?: boolean;
  interimResults?: boolean;
  maxAlternatives?: number;
  onResult?: (transcript: string) => void;
  onError?: (error: string) => void;
  onStart?: () => void;
  onEnd?: () => void;
  autoStop?: boolean;
  autoStopTimeout?: number;
}

interface VoiceSearchState {
  isListening: boolean;
  transcript: string;
  interimTranscript: string;
  error: string | null;
  isSupported: boolean;
  confidence: number;
  alternatives: string[];
}

export const useVoiceSearch = (options: VoiceSearchOptions = {}) => {
  const {
    language = 'ru-RU',
    continuous = false,
    interimResults = true,
    maxAlternatives = 3,
    onResult,
    onError,
    onStart,
    onEnd,
    autoStop = true,
    autoStopTimeout = 3000,
  } = options;

  const [state, setState] = useState<VoiceSearchState>({
    isListening: false,
    transcript: '',
    interimTranscript: '',
    error: null,
    isSupported: false,
    confidence: 0,
    alternatives: [],
  });

  const recognitionRef = useRef<any>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const isMountedRef = useRef(true);

  // Проверка поддержки Web Speech API
  useEffect(() => {
    const SpeechRecognition =
      (window as any).SpeechRecognition ||
      (window as any).webkitSpeechRecognition;

    setState((prev) => ({ ...prev, isSupported: !!SpeechRecognition }));

    if (SpeechRecognition) {
      recognitionRef.current = new SpeechRecognition();
      setupRecognition();
    }

    return () => {
      isMountedRef.current = false;
      if (recognitionRef.current) {
        recognitionRef.current.stop();
      }
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  // Настройка распознавания речи
  const setupRecognition = useCallback(() => {
    if (!recognitionRef.current) return;

    const recognition = recognitionRef.current;

    recognition.continuous = continuous;
    recognition.interimResults = interimResults;
    recognition.maxAlternatives = maxAlternatives;
    recognition.lang = language;

    // Обработка результатов
    recognition.onresult = (event: any) => {
      if (!isMountedRef.current) return;

      let finalTranscript = '';
      let interimTranscript = '';
      const alternatives: string[] = [];
      let maxConfidence = 0;

      for (let i = event.resultIndex; i < event.results.length; i++) {
        const result = event.results[i];
        const transcript = result[0].transcript;
        const confidence = result[0].confidence || 0;

        if (confidence > maxConfidence) {
          maxConfidence = confidence;
        }

        if (result.isFinal) {
          finalTranscript += transcript;

          // Собираем альтернативные варианты
          for (let j = 0; j < Math.min(result.length, maxAlternatives); j++) {
            if (j > 0) {
              alternatives.push(result[j].transcript);
            }
          }
        } else {
          interimTranscript += transcript;
        }
      }

      setState((prev) => ({
        ...prev,
        transcript: prev.transcript + finalTranscript,
        interimTranscript,
        confidence: maxConfidence,
        alternatives,
        error: null,
      }));

      if (finalTranscript) {
        onResult?.(finalTranscript);

        // Auto-stop после получения результата
        if (autoStop) {
          resetAutoStopTimer();
        }
      }
    };

    // Обработка старта
    recognition.onstart = () => {
      if (!isMountedRef.current) return;

      setState((prev) => ({ ...prev, isListening: true, error: null }));
      onStart?.();

      // Haptic feedback
      if ('vibrate' in navigator) {
        navigator.vibrate(50);
      }
    };

    // Обработка окончания
    recognition.onend = () => {
      if (!isMountedRef.current) return;

      setState((prev) => ({ ...prev, isListening: false }));
      onEnd?.();

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };

    // Обработка ошибок
    recognition.onerror = (event: any) => {
      if (!isMountedRef.current) return;

      let errorMessage = 'Speech recognition error';

      switch (event.error) {
        case 'no-speech':
          errorMessage = 'No speech detected';
          break;
        case 'audio-capture':
          errorMessage = 'No microphone found';
          break;
        case 'not-allowed':
          errorMessage = 'Microphone access denied';
          break;
        case 'network':
          errorMessage = 'Network error';
          break;
        case 'aborted':
          errorMessage = 'Recognition aborted';
          break;
      }

      setState((prev) => ({
        ...prev,
        error: errorMessage,
        isListening: false,
      }));

      onError?.(errorMessage);
    };

    // Обработка отсутствия звука
    recognition.onspeechend = () => {
      if (!isMountedRef.current) return;

      if (autoStop && !continuous) {
        recognition.stop();
      }
    };

    // Обработка изменения звука
    recognition.onaudiostart = () => {
      // Визуальная индикация начала записи
      console.log('Audio recording started');
    };

    recognition.onaudioend = () => {
      // Визуальная индикация окончания записи
      console.log('Audio recording ended');
    };
  }, [
    continuous,
    interimResults,
    maxAlternatives,
    language,
    onResult,
    onError,
    onStart,
    onEnd,
    autoStop,
  ]);

  // Сброс таймера автостопа
  const resetAutoStopTimer = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    if (autoStop && autoStopTimeout > 0) {
      timeoutRef.current = setTimeout(() => {
        if (recognitionRef.current && state.isListening) {
          recognitionRef.current.stop();
        }
      }, autoStopTimeout);
    }
  }, [autoStop, autoStopTimeout, state.isListening]);

  // Старт распознавания
  const startListening = useCallback(() => {
    if (!recognitionRef.current || state.isListening) return;

    try {
      // Сбрасываем предыдущие результаты
      setState((prev) => ({
        ...prev,
        transcript: '',
        interimTranscript: '',
        error: null,
        alternatives: [],
        confidence: 0,
      }));

      recognitionRef.current.start();

      if (autoStop) {
        resetAutoStopTimer();
      }
    } catch (error) {
      console.error('Failed to start speech recognition:', error);
      setState((prev) => ({
        ...prev,
        error: 'Failed to start speech recognition',
      }));
    }
  }, [state.isListening, autoStop, resetAutoStopTimer]);

  // Остановка распознавания
  const stopListening = useCallback(() => {
    if (!recognitionRef.current || !state.isListening) return;

    try {
      recognitionRef.current.stop();

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    } catch (error) {
      console.error('Failed to stop speech recognition:', error);
    }
  }, [state.isListening]);

  // Переключение состояния
  const toggleListening = useCallback(() => {
    if (state.isListening) {
      stopListening();
    } else {
      startListening();
    }
  }, [state.isListening, startListening, stopListening]);

  // Очистка результатов
  const clearTranscript = useCallback(() => {
    setState((prev) => ({
      ...prev,
      transcript: '',
      interimTranscript: '',
      alternatives: [],
      confidence: 0,
    }));
  }, []);

  // Обновление языка
  const setLanguage = useCallback((newLanguage: string) => {
    if (recognitionRef.current) {
      recognitionRef.current.lang = newLanguage;
    }
  }, []);

  return {
    ...state,
    startListening,
    stopListening,
    toggleListening,
    clearTranscript,
    setLanguage,
  };
};

// Хук для голосовых команд
export const useVoiceCommands = (commands: Record<string, () => void>) => {
  const [isActive, setIsActive] = useState(false);
  const [lastCommand, setLastCommand] = useState<string | null>(null);

  const handleVoiceResult = useCallback(
    (transcript: string) => {
      const normalizedTranscript = transcript.toLowerCase().trim();

      // Поиск совпадающей команды
      for (const [command, handler] of Object.entries(commands)) {
        if (normalizedTranscript.includes(command.toLowerCase())) {
          setLastCommand(command);
          handler();

          // Haptic feedback
          if ('vibrate' in navigator) {
            navigator.vibrate([50, 50, 50]);
          }

          break;
        }
      }
    },
    [commands]
  );

  const { isListening, isSupported, error, startListening, stopListening } =
    useVoiceSearch({
      onResult: handleVoiceResult,
      continuous: true,
      autoStop: false,
    });

  const toggleVoiceCommands = useCallback(() => {
    if (isActive) {
      stopListening();
      setIsActive(false);
    } else {
      startListening();
      setIsActive(true);
    }
  }, [isActive, startListening, stopListening]);

  return {
    isActive,
    isListening,
    isSupported,
    error,
    lastCommand,
    toggleVoiceCommands,
  };
};

// Хук для диктовки текста
export const useVoiceDictation = (
  onTextReceived: (text: string) => void,
  options: { punctuation?: boolean; corrections?: boolean } = {}
) => {
  const { punctuation = true, corrections = true } = options;

  const processDictatedText = useCallback(
    (transcript: string) => {
      let processedText = transcript;

      // Добавление пунктуации на основе голосовых команд
      if (punctuation) {
        processedText = processedText
          .replace(/точка/gi, '.')
          .replace(/запятая/gi, ',')
          .replace(/вопросительный знак/gi, '?')
          .replace(/восклицательный знак/gi, '!')
          .replace(/двоеточие/gi, ':')
          .replace(/точка с запятой/gi, ';')
          .replace(/новая строка|новый абзац/gi, '\n')
          .replace(/открыть кавычки/gi, '"')
          .replace(/закрыть кавычки/gi, '"');
      }

      // Автоматические исправления
      if (corrections) {
        // Капитализация начала предложения
        processedText = processedText.replace(
          /(^|[.!?]\s+)([a-zа-я])/g,
          (match, p1, p2) => p1 + p2.toUpperCase()
        );

        // Удаление лишних пробелов
        processedText = processedText.replace(/\s+/g, ' ').trim();
      }

      onTextReceived(processedText);
    },
    [punctuation, corrections, onTextReceived]
  );

  return useVoiceSearch({
    onResult: processDictatedText,
    continuous: true,
    interimResults: true,
  });
};
