// src/components/Card/Card.tsx

import React from 'react';
import styles from './Card.module.scss';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  padding?: 'none' | 'sm' | 'md' | 'lg';
  elevation?: 0 | 1 | 2 | 3 | 4 | 5;
  onClick?: () => void;
}

export function Card({ 
  children, 
  className = '', 
  padding = 'md',
  elevation = 1,
  onClick 
}: CardProps) {
  const cardClasses = [
    styles.card,
    styles[`padding-${padding}`],
    styles[`elevation-${elevation}`],
    onClick ? styles.clickable : '',
    className
  ].filter(Boolean).join(' ');

  const handleClick = () => {
    if (onClick) {
      onClick();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (onClick && (e.key === 'Enter' || e.key === ' ')) {
      e.preventDefault();
      onClick();
    }
  };

  if (onClick) {
    return (
      <div 
        className={cardClasses}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        role="button"
        tabIndex={0}
        aria-label="Clickable card"
      >
        {children}
      </div>
    );
  }

  return (
    <div className={cardClasses}>
      {children}
    </div>
  );
}

interface CardHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export function CardHeader({ children, className = '' }: CardHeaderProps) {
  return (
    <div className={`${styles.header} ${className}`}>
      {children}
    </div>
  );
}

interface CardTitleProps {
  children: React.ReactNode;
  className?: string;
  level?: 1 | 2 | 3 | 4 | 5 | 6;
}

export function CardTitle({ children, className = '', level = 3 }: CardTitleProps) {
  const Tag = `h${level}` as keyof JSX.IntrinsicElements;
  
  return (
    <Tag className={`${styles.title} ${className}`}>
      {children}
    </Tag>
  );
}

interface CardContentProps {
  children: React.ReactNode;
  className?: string;
}

export function CardContent({ children, className = '' }: CardContentProps) {
  return (
    <div className={`${styles.content} ${className}`}>
      {children}
    </div>
  );
}

interface CardActionsProps {
  children: React.ReactNode;
  className?: string;
  align?: 'left' | 'right' | 'center' | 'between';
}

export function CardActions({ children, className = '', align = 'right' }: CardActionsProps) {
  return (
    <div className={`${styles.actions} ${styles[`align-${align}`]} ${className}`}>
      {children}
    </div>
  );
}