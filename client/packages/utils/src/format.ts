import { format, parseISO } from 'date-fns';

export const formatDate = (date: string | Date, formatStr: string = 'dd/MM/yyyy'): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date;
    return format(dateObj, formatStr);
  } catch {
    return '';
  }
};

export const formatTime = (time: string | Date, formatStr: string = 'HH:mm'): string => {
  try {
    const timeObj = typeof time === 'string' ? parseISO(time) : time;
    return format(timeObj, formatStr);
  } catch {
    return '';
  }
};

export const formatDateTime = (dateTime: string | Date): string => {
  return formatDate(dateTime, 'dd/MM/yyyy HH:mm');
};

export const formatDuration = (hours: number): string => {
  const h = Math.floor(hours);
  const m = Math.round((hours - h) * 60);
  return `${h}h ${m}m`;
};

export const formatPercentage = (value: number, decimals: number = 1): string => {
  return `${(value * 100).toFixed(decimals)}%`;
};
