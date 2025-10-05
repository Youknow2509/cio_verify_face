// src/components/Header/Header.tsx

import { useState, useRef, useEffect } from 'react';
import { useSidebar } from '../../hooks/useSidebar';
import styles from './Header.module.scss';

interface HeaderProps {
  className?: string;
}

export function Header({ className = '' }: HeaderProps) {
  const { isCollapsed, toggle } = useSidebar();
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement>(null);

  // Close menu when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const toggleUserMenu = () => {
    setIsUserMenuOpen(!isUserMenuOpen);
  };

  const handleProfileClick = () => {
    setIsUserMenuOpen(false);
    // TODO: Navigate to profile page
    console.log('Navigating to profile page...');
  };

  const handleSettingsClick = () => {
    setIsUserMenuOpen(false);
    // TODO: Navigate to settings page
    console.log('Navigating to settings page...');
  };

  const handleLogoutClick = () => {
    setIsUserMenuOpen(false);
    // TODO: Implement logout logic
    if (window.confirm('Bạn có chắc chắn muốn đăng xuất?')) {
      console.log('Logging out...');
      // Add logout logic here
    }
  };

  const MenuIcon = () => (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
      <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
    </svg>
  );

  const SearchIcon = () => (
    <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
      <path fillRule="evenodd" d="M12.9 14.32a8 8 0 111.41-1.41l5.35 5.33-1.42 1.42-5.33-5.34zM8 14A6 6 0 108 2a6 6 0 000 12z" clipRule="evenodd" />
    </svg>
  );

  const NotificationIcon = () => (
    <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
      <path d="M10 2C7.79 2 6 3.79 6 6c0 2.12-1.19 3.84-1.67 4.96-.08.18-.33.04-.33-.04V10c0-.55-.45-1-1-1s-1 .45-1 1v1c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2v-1c0-.55-.45-1-1-1s-1 .45-1 1v.92c0 .08-.25.22-.33.04C15.19 9.84 14 8.12 14 6c0-2.21-1.79-4-4-4z" />
      <path d="M8 17c0 1.1.9 2 2 2s2-.9 2-2" />
    </svg>
  );

  const UserIcon = () => (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M8 8a3 3 0 100-6 3 3 0 000 6zM8 9a5 5 0 00-5 5 1 1 0 102 0 3 3 0 016 0 1 1 0 102 0 5 5 0 00-5-5z" />
    </svg>
  );

  const SettingsIcon = () => (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M8 4.754a3.246 3.246 0 100 6.492 3.246 3.246 0 000-6.492zM5.754 8a2.246 2.246 0 114.492 0 2.246 2.246 0 01-4.492 0z" />
      <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 01-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 01-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 01.52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 011.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 011.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 01.52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 01-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 01-1.255-.52l-.094-.319zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 002.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 001.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 00-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 00-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 00-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 001.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 003.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 002.692-1.115l.094-.319z" />
    </svg>
  );

  const LogoutIcon = () => (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M10 12.5a.5.5 0 01-.5.5h-8a.5.5 0 01-.5-.5v-9a.5.5 0 01.5-.5h8a.5.5 0 01.5.5v2a.5.5 0 001 0v-2A1.5 1.5 0 009.5 2h-8A1.5 1.5 0 000 3.5v9A1.5 1.5 0 001.5 14h8a1.5 1.5 0 001.5-1.5v-2a.5.5 0 00-1 0v2z" />
      <path d="M15.854 8.354a.5.5 0 000-.708l-3-3a.5.5 0 00-.708.708L14.293 7.5H5.5a.5.5 0 000 1h8.793l-2.147 2.146a.5.5 0 00.708.708l3-3z" />
    </svg>
  );

  const ChevronDownIcon = () => (
    <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
      <path d="M2.22 4.47a.75.75 0 011.06 0L6 7.19l2.72-2.72a.75.75 0 111.06 1.06L6.53 8.78a.75.75 0 01-1.06 0L2.22 5.53a.75.75 0 010-1.06z" />
    </svg>
  );

  return (
    <header className={`${styles.header} ${className}`}>
      <div className={styles.left}>
        <button
          className={styles.menuButton}
          onClick={toggle}
          aria-pressed={!isCollapsed}
          aria-label={isCollapsed ? 'Mở sidebar' : 'Đóng sidebar'}
        >
          <MenuIcon />
        </button>
        
        <div className={styles.logo}>
          <h1 className={styles.logoText}>Face Attendance</h1>
        </div>
      </div>

      <div className={styles.center}>
        <div className={styles.search}>
          <SearchIcon />
          <input
            type="text"
            placeholder="Tìm kiếm..."
            className={styles.searchInput}
            aria-label="Tìm kiếm"
          />
        </div>
      </div>

      <div className={styles.right}>
        <button
          className={styles.iconButton}
          aria-label="Thông báo"
        >
          <NotificationIcon />
          <span className={styles.badge}>3</span>
        </button>

        <div className={styles.userMenu} ref={userMenuRef}>
          <button
            className={styles.avatar}
            onClick={toggleUserMenu}
            aria-label="Menu người dùng"
            aria-expanded={isUserMenuOpen}
          >
            <span className={styles.avatarText}>NA</span>
            <ChevronDownIcon />
          </button>
          
          {isUserMenuOpen && (
            <div className={styles.dropdownMenu}>
              <div className={styles.userInfo}>
                <div className={styles.userAvatar}>
                  <span className={styles.avatarText}>NA</span>
                </div>
                <div className={styles.userDetails}>
                  <span className={styles.userName}>Nguyễn Admin</span>
                  <span className={styles.userEmail}>admin@company.com</span>
                </div>
              </div>
              
              <div className={styles.menuDivider}></div>
              
              <ul className={styles.menuItems}>
                <li>
                  <button className={styles.menuItem} onClick={handleProfileClick}>
                    <UserIcon />
                    <span>Hồ sơ cá nhân</span>
                  </button>
                </li>
                <li>
                  <button className={styles.menuItem} onClick={handleSettingsClick}>
                    <SettingsIcon />
                    <span>Cài đặt</span>
                  </button>
                </li>
                <li>
                  <button className={styles.menuItem} data-danger="true" onClick={handleLogoutClick}>
                    <LogoutIcon />
                    <span>Đăng xuất</span>
                  </button>
                </li>
              </ul>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}