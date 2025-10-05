// src/components/Badge/Badge.tsx

import React from 'react';
import styles from './Badge.module.scss';

interface BadgeProps {
  children: React.ReactNode;
  variant?: 'success' | 'error' | 'warning' | 'info' | 'neutral';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  icon?: React.ReactNode;
}

export function Badge({ 
  children, 
  variant = 'neutral', 
  size = 'md',
  className = '',
  icon 
}: BadgeProps) {
  const badgeClasses = [
    styles.badge,
    styles[variant],
    styles[size],
    className
  ].filter(Boolean).join(' ');

  return (
    <span className={badgeClasses}>
      {icon && <span className={styles.icon}>{icon}</span>}
      <span className={styles.text}>{children}</span>
    </span>
  );
}

// Status-specific badge components for common use cases
export function StatusBadge({ 
  status, 
  className = '' 
}: { 
  status: 'online' | 'offline' | 'active' | 'inactive' | 'late' | 'on-time'; 
  className?: string;
}) {
  const getVariant = (status: string) => {
    switch (status) {
      case 'online':
      case 'active':
      case 'on-time':
        return 'success';
      case 'offline':
      case 'inactive':
        return 'error';
      case 'late':
        return 'warning';
      default:
        return 'neutral';
    }
  };

  const getLabel = (status: string) => {
    switch (status) {
      case 'online':
        return 'Trực tuyến';
      case 'offline':
        return 'Ngoại tuyến';
      case 'active':
        return 'Hoạt động';
      case 'inactive':
        return 'Không hoạt động';
      case 'late':
        return 'Trễ giờ';
      case 'on-time':
        return 'Đúng giờ';
      default:
        return status;
    }
  };

  const StatusIcon = ({ status }: { status: string }) => {
    const iconSize = 8;
    
    switch (status) {
      case 'online':
      case 'active':
      case 'on-time':
        return (
          <svg width={iconSize} height={iconSize} viewBox="0 0 8 8" fill="currentColor">
            <circle cx="4" cy="4" r="4" />
          </svg>
        );
      case 'offline':
      case 'inactive':
        return (
          <svg width={iconSize} height={iconSize} viewBox="0 0 8 8" fill="currentColor">
            <circle cx="4" cy="4" r="4" />
          </svg>
        );
      case 'late':
        return (
          <svg width={iconSize} height={iconSize} viewBox="0 0 8 8" fill="currentColor">
            <circle cx="4" cy="4" r="4" />
          </svg>
        );
      default:
        return null;
    }
  };

  return (
    <Badge 
      variant={getVariant(status)} 
      className={className}
      icon={<StatusIcon status={status} />}
    >
      {getLabel(status)}
    </Badge>
  );
}