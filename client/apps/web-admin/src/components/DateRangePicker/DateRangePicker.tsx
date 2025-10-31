// src/components/DateRangePicker/DateRangePicker.tsx

import { useState } from 'react';
import { DatePicker } from '@/components/DatePicker/DatePicker';
import styles from './DateRangePicker.module.scss';

interface DateRangePickerProps {
  startDate: Date | null;
  endDate: Date | null;
  onChange: (startDate: Date | null, endDate: Date | null) => void;
  disabled?: boolean;
  className?: string;
  startPlaceholder?: string;
  endPlaceholder?: string;
}

export function DateRangePicker({
  startDate,
  endDate,
  onChange,
  disabled = false,
  className,
  startPlaceholder = 'Ngày bắt đầu',
  endPlaceholder = 'Ngày kết thúc'
}: DateRangePickerProps) {
  const [focusedInput, setFocusedInput] = useState<'start' | 'end' | null>(null);

  const handleStartDateChange = (date: Date | null) => {
    onChange(date, endDate);
  };

  const handleEndDateChange = (date: Date | null) => {
    onChange(startDate, date);
  };

  return (
    <div className={`${styles.dateRangePicker} ${className || ''}`}>
      <div className={styles.datePickerGroup}>
        <DatePicker
          selected={startDate}
          onChange={handleStartDateChange}
          maxDate={endDate || undefined}
          placeholder={startPlaceholder}
          disabled={disabled}
          className={`${styles.startDatePicker} ${focusedInput === 'start' ? styles.focused : ''}`}
        />
        
        <div className={styles.separator}>
          <span>đến</span>
        </div>
        
        <DatePicker
          selected={endDate}
          onChange={handleEndDateChange}
          minDate={startDate || undefined}
          placeholder={endPlaceholder}
          disabled={disabled}
          className={`${styles.endDatePicker} ${focusedInput === 'end' ? styles.focused : ''}`}
        />
      </div>
    </div>
  );
}

export default DateRangePicker;