// src/components/FilterBar/FilterBar.tsx

import { ReactNode } from 'react';
import { Search, Filter, X } from 'lucide-react';
import styles from './FilterBar.module.scss';

interface FilterBarProps {
  children: ReactNode;
  onClear?: () => void;
  hasActiveFilters?: boolean;
  className?: string;
}

interface SearchBoxProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

interface FilterGroupProps {
  children: ReactNode;
  label?: string;
  className?: string;
}

interface FilterSelectProps {
  value: string;
  onChange: (value: string) => void;
  options: Array<{ value: string; label: string }>;
  placeholder?: string;
  className?: string;
}

export function FilterBar({ 
  children, 
  onClear, 
  hasActiveFilters = false,
  className 
}: FilterBarProps) {
  return (
    <div className={`${styles.filterBar} ${className || ''}`}>
      <div className={styles.filterContent}>
        {children}
      </div>
      {hasActiveFilters && onClear && (
        <button
          type="button"
          onClick={onClear}
          className={styles.clearButton}
          title="Xóa bộ lọc"
        >
          <X size={16} />
          Xóa bộ lọc
        </button>
      )}
    </div>
  );
}

export function SearchBox({ 
  value, 
  onChange, 
  placeholder = 'Tìm kiếm...', 
  className 
}: SearchBoxProps) {
  return (
    <div className={`${styles.searchBox} ${className || ''}`}>
      <Search className={styles.searchIcon} size={18} />
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className={styles.searchInput}
      />
      {value && (
        <button
          type="button"
          onClick={() => onChange('')}
          className={styles.clearSearchButton}
          title="Xóa tìm kiếm"
        >
          <X size={14} />
        </button>
      )}
    </div>
  );
}

export function FilterGroup({ children, label, className }: FilterGroupProps) {
  return (
    <div className={`${styles.filterGroup} ${className || ''}`}>
      {label && (
        <label className={styles.filterLabel}>
          <Filter size={14} />
          {label}
        </label>
      )}
      <div className={styles.filterGroupContent}>
        {children}
      </div>
    </div>
  );
}

export function FilterSelect({ 
  value, 
  onChange, 
  options, 
  placeholder = 'Chọn...', 
  className 
}: FilterSelectProps) {
  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={`${styles.filterSelect} ${className || ''}`}
    >
      <option value="">{placeholder}</option>
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  );
}

export default FilterBar;