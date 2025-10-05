// src/layouts/Layout/Layout.tsx

import { Header } from '../../components/Header/Header';
import { Sidebar } from '../../components/Sidebar/Sidebar';
import { useSidebar } from '../../hooks/useSidebar';
import styles from './Layout.module.scss';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const { isCollapsed } = useSidebar();

  return (
    <div className={`${styles.layout} ${isCollapsed ? styles.collapsed : ''}`}>
      <Header className={styles.header} />
      <Sidebar className={styles.sidebar} />
      <main className={styles.content}>
        <div className={styles.contentInner}>
          {children}
        </div>
      </main>
    </div>
  );
}