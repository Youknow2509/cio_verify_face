// src/components/Table/Table.tsx

import React from 'react';
import styles from './Table.module.scss';

interface Column<T> {
  key: keyof T | string;
  header: string;
  width?: string;
  align?: 'left' | 'center' | 'right';
  sortable?: boolean;
  render?: (value: any, row: T, index: number) => React.ReactNode;
}

interface TableProps<T> {
  data: T[];
  columns: Column<T>[];
  loading?: boolean;
  empty?: React.ReactNode;
  className?: string;
  responsive?: boolean;
  stickyHeader?: boolean;
  onSort?: (key: string, direction: 'asc' | 'desc') => void;
  sortKey?: string;
  sortDirection?: 'asc' | 'desc';
}

export function Table<T extends Record<string, any>>({
  data,
  columns,
  loading = false,
  empty,
  className = '',
  responsive = true,
  stickyHeader = false,
  onSort,
  sortKey,
  sortDirection
}: TableProps<T>) {
  const handleSort = (key: string) => {
    if (!onSort) return;
    
    const newDirection = sortKey === key && sortDirection === 'asc' ? 'desc' : 'asc';
    onSort(key, newDirection);
  };

  const getSortIcon = (key: string) => {
    if (sortKey !== key) {
      return (
        <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
          <path d="M6 2l3 3H3l3-3zM6 10L3 7h6l-3 3z" opacity="0.3" />
        </svg>
      );
    }
    
    if (sortDirection === 'asc') {
      return (
        <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
          <path d="M6 2l3 3H3l3-3z" />
        </svg>
      );
    }
    
    return (
      <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
        <path d="M6 10L3 7h6l-3 3z" />
      </svg>
    );
  };

  const tableClasses = [
    styles.table,
    stickyHeader ? styles.stickyHeader : '',
    className
  ].filter(Boolean).join(' ');

  const wrapperClasses = [
    styles.wrapper,
    responsive ? styles.responsive : ''
  ].filter(Boolean).join(' ');

  if (loading) {
    return (
      <div className={wrapperClasses}>
        <div className={styles.loading}>
          <div className={styles.spinner}></div>
          <span>Đang tải...</span>
        </div>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className={wrapperClasses}>
        <div className={styles.empty}>
          {empty || (
            <>
              <div className={styles.emptyIcon}>
                <svg width="48" height="48" viewBox="0 0 48 48" fill="currentColor">
                  <path d="M6 10c-1.1 0-2 .9-2 2v24c0 1.1.9 2 2 2h36c1.1 0 2-.9 2-2V12c0-1.1-.9-2-2-2H6zm0 2h36v4H6v-4zm0 6h36v18H6V18z" opacity="0.3" />
                </svg>
              </div>
              <p>Không có dữ liệu</p>
            </>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className={wrapperClasses}>
      <table className={tableClasses}>
        <thead>
          <tr>
            {columns.map((column, index) => (
              <th
                key={String(column.key) + index}
                style={{ width: column.width, textAlign: column.align }}
                className={column.sortable ? styles.sortable : ''}
                onClick={column.sortable ? () => handleSort(String(column.key)) : undefined}
                role={column.sortable ? 'button' : undefined}
                tabIndex={column.sortable ? 0 : undefined}
                onKeyDown={column.sortable ? (e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    handleSort(String(column.key));
                  }
                } : undefined}
              >
                <div className={styles.headerContent}>
                  <span>{column.header}</span>
                  {column.sortable && (
                    <span className={styles.sortIcon}>
                      {getSortIcon(String(column.key))}
                    </span>
                  )}
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {columns.map((column, colIndex) => {
                const value = row[column.key as keyof T];
                const content = column.render ? column.render(value, row, rowIndex) : value;
                
                return (
                  <td
                    key={String(column.key) + colIndex}
                    style={{ textAlign: column.align }}
                  >
                    {content}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}