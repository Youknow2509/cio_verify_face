// src/hooks/useSidebar.ts

import { useState, useCallback, useEffect } from 'react';

const SIDEBAR_STORAGE_KEY = 'sidebar-collapsed';

export function useSidebar() {
  const [isCollapsed, setIsCollapsed] = useState(() => {
    try {
      const stored = localStorage.getItem(SIDEBAR_STORAGE_KEY);
      return stored ? JSON.parse(stored) : false;
    } catch {
      return false;
    }
  });

  const toggle = useCallback(() => {
    setIsCollapsed((prev: boolean) => {
      const newValue = !prev;
      try {
        localStorage.setItem(SIDEBAR_STORAGE_KEY, JSON.stringify(newValue));
      } catch {
        // Ignore localStorage errors
      }
      return newValue;
    });
  }, []);

  const collapse = useCallback(() => {
    setIsCollapsed(true);
    try {
      localStorage.setItem(SIDEBAR_STORAGE_KEY, JSON.stringify(true));
    } catch {
      // Ignore localStorage errors
    }
  }, []);

  const expand = useCallback(() => {
    setIsCollapsed(false);
    try {
      localStorage.setItem(SIDEBAR_STORAGE_KEY, JSON.stringify(false));
    } catch {
      // Ignore localStorage errors
    }
  }, []);

  // Handle responsive behavior
  useEffect(() => {
    const handleResize = () => {
      // Auto-collapse on mobile
      if (window.innerWidth < 900) {
        setIsCollapsed(true);
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return {
    isCollapsed,
    toggle,
    collapse,
    expand
  };
}