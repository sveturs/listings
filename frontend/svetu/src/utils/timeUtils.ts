import {
  differenceInMinutes,
  differenceInHours,
  differenceInDays,
  differenceInWeeks,
  differenceInMonths,
} from 'date-fns';

type TranslateFunction = {
  (key: string, values?: Record<string, string | number | Date>): string;
};

export function getLastSeenText(
  lastSeen: string | Date,
  t: TranslateFunction
): string {
  const now = new Date();
  const lastSeenDate =
    typeof lastSeen === 'string' ? new Date(lastSeen) : lastSeen;

  const minutes = differenceInMinutes(now, lastSeenDate);
  const hours = differenceInHours(now, lastSeenDate);
  const days = differenceInDays(now, lastSeenDate);
  const weeks = differenceInWeeks(now, lastSeenDate);
  const months = differenceInMonths(now, lastSeenDate);

  if (minutes < 1) {
    return t('lastSeenJustNow');
  } else if (minutes < 60) {
    return t('lastSeenMinutesAgo', { count: minutes });
  } else if (hours < 24) {
    return t('lastSeenHoursAgo', { count: hours });
  } else if (days < 7) {
    return t('lastSeenDaysAgo', { count: days });
  } else if (weeks < 4) {
    return t('lastSeenWeeksAgo', { count: weeks });
  } else {
    return t('lastSeenMonthsAgo', { count: months });
  }
}
