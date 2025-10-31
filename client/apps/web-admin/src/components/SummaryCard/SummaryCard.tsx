// src/components/SummaryCard/SummaryCard.tsx

import { ReactNode } from 'react';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import styles from './SummaryCard.module.scss';

interface SummaryCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: ReactNode;
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  trend?: {
    value: number;
    type: 'up' | 'down' | 'neutral';
    text?: string;
  };
  className?: string;
}

export function SummaryCard({
  title,
  value,
  subtitle,
  icon,
  variant = 'neutral',
  trend,
  className
}: SummaryCardProps) {
  const renderTrendIcon = () => {
    if (!trend) return null;
    
    switch (trend.type) {
      case 'up':
        return <TrendingUp size={16} />;
      case 'down':
        return <TrendingDown size={16} />;
      default:
        return <Minus size={16} />;
    }
  };

  const getTrendClassName = () => {
    if (!trend) return '';
    
    switch (trend.type) {
      case 'up':
        return styles.trendUp;
      case 'down':
        return styles.trendDown;
      default:
        return styles.trendNeutral;
    }
  };

  return (
    <div className={`${styles.summaryCard} ${styles[variant]} ${className || ''}`}>
      <div className={styles.header}>
        {icon && (
          <div className={styles.iconWrapper}>
            {icon}
          </div>
        )}
        <div className={styles.headerContent}>
          <h3 className={styles.title}>{title}</h3>
          {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
        </div>
      </div>
      
      <div className={styles.body}>
        <div className={styles.value}>
          {typeof value === 'number' ? value.toLocaleString('vi-VN') : value}
        </div>
        
        {trend && (
          <div className={`${styles.trend} ${getTrendClassName()}`}>
            {renderTrendIcon()}
            <span className={styles.trendValue}>
              {Math.abs(trend.value)}%
            </span>
            {trend.text && (
              <span className={styles.trendText}>{trend.text}</span>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default SummaryCard;