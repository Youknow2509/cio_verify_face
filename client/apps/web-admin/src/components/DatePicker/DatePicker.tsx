// src/components/DatePicker/DatePicker.tsx

import { forwardRef } from 'react';
import ReactDatePicker from 'react-datepicker';
import { CalendarDays } from 'lucide-react';
import styles from './DatePicker.module.scss';
import 'react-datepicker/dist/react-datepicker.css';

interface DatePickerProps {
  selected?: Date | null;
  onChange: (date: Date | null) => void;
  placeholder?: string;
  disabled?: boolean;
  minDate?: Date;
  maxDate?: Date;
  dateFormat?: string;
  isClearable?: boolean;
  selectsRange?: boolean;
  startDate?: Date | null;
  endDate?: Date | null;
  className?: string;
  showMonthDropdown?: boolean;
  showYearDropdown?: boolean;
  dropdownMode?: 'scroll' | 'select';
}

// Custom input component
const CustomInput = forwardRef<HTMLInputElement, any>(
  ({ value, onClick, placeholder, disabled, className }, ref) => (
    <div className={`${styles.datePickerWrapper} ${className || ''}`}>
      <input
        ref={ref}
        value={value}
        onClick={onClick}
        placeholder={placeholder}
        disabled={disabled}
        readOnly
        className={styles.dateInput}
      />
      <CalendarDays className={styles.calendarIcon} size={18} />
    </div>
  )
);

CustomInput.displayName = 'CustomInput';

export function DatePicker({
  selected,
  onChange,
  placeholder = 'Chọn ngày',
  disabled = false,
  minDate,
  maxDate,
  dateFormat = 'dd/MM/yyyy',
  isClearable = true,
  selectsRange = false,
  startDate,
  endDate,
  className,
  showMonthDropdown = true,
  showYearDropdown = true,
  dropdownMode = 'select',
}: DatePickerProps) {
  const datePickerProps: any = {
    selected,
    onChange,
    customInput: <CustomInput className={className} />,
    placeholderText: placeholder,
    disabled,
    minDate,
    maxDate,
    dateFormat,
    isClearable,
    showMonthDropdown,
    showYearDropdown,
    dropdownMode,
    calendarClassName: styles.calendar,
    wrapperClassName: styles.wrapper,
    popperClassName: styles.popper,
    dayClassName: () => styles.day,
  };

  if (selectsRange) {
    datePickerProps.selectsRange = true;
    datePickerProps.startDate = startDate;
    datePickerProps.endDate = endDate;
  }

  return <ReactDatePicker {...datePickerProps} />;
}

export default DatePicker;
