// src/components/NavItem/NavItem.tsx

import { NavLink } from 'react-router-dom';
import styles from './NavItem.module.scss';

interface NavItemProps {
  to: string;
  icon: React.ComponentType;
  label: string;
  collapsed?: boolean;
  className?: string;
}

export function NavItem({ to, icon: Icon, label, collapsed = false, className = '' }: NavItemProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) => 
        `${styles.navItem} ${isActive ? styles.active : ''} ${collapsed ? styles.collapsed : ''} ${className}`
      }
      title={collapsed ? label : undefined}
    >
      <span className={styles.icon}>
        <Icon />
      </span>
      <span className={styles.label}>{label}</span>
    </NavLink>
  );
}