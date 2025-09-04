import { useCallback, useEffect, useRef } from 'react';

interface HapticPattern {
  type:
    | 'success'
    | 'warning'
    | 'error'
    | 'light'
    | 'medium'
    | 'heavy'
    | 'selection'
    | 'custom';
  duration?: number;
  intensity?: number;
  pattern?: number[];
}

interface UseHapticFeedbackReturn {
  vibrate: (pattern?: HapticPattern) => void;
  success: () => void;
  warning: () => void;
  error: () => void;
  impact: (style?: 'light' | 'medium' | 'heavy') => void;
  selection: () => void;
  notification: (type?: 'success' | 'warning' | 'error') => void;
  customPattern: (pattern: number[]) => void;
  isSupported: boolean;
  hasHapticEngine: boolean;
}

export function useHapticFeedback(): UseHapticFeedbackReturn {
  const isSupported =
    typeof navigator !== 'undefined' && 'vibrate' in navigator;
  const hasHapticEngine = useRef(false);
  const lastVibrationTime = useRef(0);

  useEffect(() => {
    // Check for iOS Haptic Engine support
    const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
    if (isIOS && 'Haptics' in window) {
      hasHapticEngine.current = true;
    }

    // Check for Android Vibration API enhancements
    if ('vibrate' in navigator && navigator.vibrate) {
      // Test if complex patterns are supported
      try {
        navigator.vibrate([1, 1, 1]);
        hasHapticEngine.current = true;
      } catch {
        // Simple vibration only
      }
    }
  }, []);

  /**
   * Prevent vibration spamming
   */
  const throttleVibration = (minInterval: number = 50): boolean => {
    const now = Date.now();
    if (now - lastVibrationTime.current < minInterval) {
      return false;
    }
    lastVibrationTime.current = now;
    return true;
  };

  /**
   * Main vibrate function
   */
  const vibrate = useCallback(
    (pattern?: HapticPattern) => {
      if (!isSupported || !throttleVibration()) return;

      try {
        if (hasHapticEngine.current && 'Haptics' in window) {
          // Use iOS Haptic Engine if available
          const haptics = (window as any).Haptics;

          switch (pattern?.type) {
            case 'success':
              haptics.notificationOccurred('success');
              break;
            case 'warning':
              haptics.notificationOccurred('warning');
              break;
            case 'error':
              haptics.notificationOccurred('error');
              break;
            case 'light':
              haptics.impactOccurred('light');
              break;
            case 'medium':
              haptics.impactOccurred('medium');
              break;
            case 'heavy':
              haptics.impactOccurred('heavy');
              break;
            case 'selection':
              haptics.selectionChanged();
              break;
            default:
              haptics.impactOccurred('light');
          }
        } else {
          // Use standard Vibration API
          let vibrationPattern: number | number[] = 10;

          if (pattern?.pattern) {
            vibrationPattern = pattern.pattern;
          } else {
            switch (pattern?.type) {
              case 'success':
                vibrationPattern = [10, 50, 10]; // Quick double tap
                break;
              case 'warning':
                vibrationPattern = [20, 100, 20]; // Medium double tap
                break;
              case 'error':
                vibrationPattern = [50, 100, 50, 100, 50]; // Triple tap
                break;
              case 'light':
                vibrationPattern = 5;
                break;
              case 'medium':
                vibrationPattern = 10;
                break;
              case 'heavy':
                vibrationPattern = 20;
                break;
              case 'selection':
                vibrationPattern = 3;
                break;
              case 'custom':
                vibrationPattern = pattern?.duration || 10;
                break;
              default:
                vibrationPattern = 10;
            }
          }

          navigator.vibrate(vibrationPattern);
        }
      } catch (error) {
        console.debug('Haptic feedback error:', error);
      }
    },
    [isSupported]
  );

  /**
   * Success haptic (e.g., task completed, form submitted)
   */
  const success = useCallback(() => {
    vibrate({ type: 'success' });
  }, [vibrate]);

  /**
   * Warning haptic (e.g., approaching limit, caution needed)
   */
  const warning = useCallback(() => {
    vibrate({ type: 'warning' });
  }, [vibrate]);

  /**
   * Error haptic (e.g., invalid input, operation failed)
   */
  const error = useCallback(() => {
    vibrate({ type: 'error' });
  }, [vibrate]);

  /**
   * Impact haptic with different intensities
   */
  const impact = useCallback(
    (style: 'light' | 'medium' | 'heavy' = 'light') => {
      vibrate({ type: style });
    },
    [vibrate]
  );

  /**
   * Selection haptic (e.g., picker wheel, toggle switch)
   */
  const selection = useCallback(() => {
    vibrate({ type: 'selection' });
  }, [vibrate]);

  /**
   * Notification haptic
   */
  const notification = useCallback(
    (type: 'success' | 'warning' | 'error' = 'success') => {
      vibrate({ type });
    },
    [vibrate]
  );

  /**
   * Custom vibration pattern
   */
  const customPattern = useCallback(
    (pattern: number[]) => {
      if (!isSupported || !throttleVibration()) return;

      try {
        navigator.vibrate(pattern);
      } catch (error) {
        console.debug('Custom pattern error:', error);
      }
    },
    [isSupported]
  );

  return {
    vibrate,
    success,
    warning,
    error,
    impact,
    selection,
    notification,
    customPattern,
    isSupported,
    hasHapticEngine: hasHapticEngine.current,
  };
}

/**
 * Haptic feedback presets for common interactions
 */
export const HapticPresets = {
  // Button interactions
  buttonTap: () => [3],
  buttonLongPress: () => [10, 10, 10],

  // Gestures
  swipeSuccess: () => [5, 20, 5],
  pullToRefresh: () => [10, 30, 10],
  pinchZoom: () => [3, 3, 3],

  // Notifications
  messageReceived: () => [20, 100, 20],
  messageSent: () => [10],

  // Form interactions
  inputFocus: () => [3],
  inputError: () => [30, 50, 30],
  formSubmit: () => [10, 50, 10, 50, 10],

  // Navigation
  pageTransition: () => [5],
  tabSwitch: () => [3],
  modalOpen: () => [10],
  modalClose: () => [5],

  // Feedback
  likeAction: () => [5, 10, 5],
  deleteAction: () => [50, 100, 50],
  copyAction: () => [10, 30, 10],

  // Gaming/Fun
  achievement: () => [10, 30, 10, 30, 10, 100],
  levelUp: () => [20, 40, 20, 40, 20, 40, 100],
  coinCollect: () => [5, 5, 5],

  // Loading states
  loadingPulse: () => [200, 200, 200, 200], // Continuous pulse
  loadComplete: () => [10, 50, 10],
};

/**
 * Hook for gesture-based haptic feedback
 */
export function useGestureHaptics() {
  const haptic = useHapticFeedback();

  return {
    onTouchStart: useCallback(() => {
      haptic.impact('light');
    }, [haptic]),

    onTouchEnd: useCallback(() => {
      haptic.selection();
    }, [haptic]),

    onLongPress: useCallback(() => {
      haptic.impact('medium');
    }, [haptic]),

    onSwipe: useCallback(() => {
      haptic.customPattern(HapticPresets.swipeSuccess());
    }, [haptic]),

    onPinch: useCallback(() => {
      haptic.customPattern(HapticPresets.pinchZoom());
    }, [haptic]),

    onDragStart: useCallback(() => {
      haptic.impact('medium');
    }, [haptic]),

    onDragEnd: useCallback(() => {
      haptic.success();
    }, [haptic]),

    onDoubleTap: useCallback(() => {
      haptic.customPattern([5, 10, 5]);
    }, [haptic]),
  };
}
