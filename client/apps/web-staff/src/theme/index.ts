import { createTheme, Theme } from '@mui/material/styles';

const getDesignTokens = (mode: 'light' | 'dark') => ({
    palette: {
        mode,
        primary: {
            main: '#2563eb',
            light: '#60a5fa',
            dark: '#1e40af',
        },
        secondary: {
            main: '#10b981',
            light: '#34d399',
            dark: '#059669',
        },
        error: {
            main: '#ef4444',
            light: '#f87171',
            dark: '#dc2626',
        },
        warning: {
            main: '#f59e0b',
            light: '#fbbf24',
            dark: '#d97706',
        },
        info: {
            main: '#06b6d4',
            light: '#22d3ee',
            dark: '#0891b2',
        },
        success: {
            main: '#10b981',
            light: '#34d399',
            dark: '#059669',
        },
        background: {
            default: mode === 'light' ? '#f8fafc' : '#0f172a',
            paper: mode === 'light' ? '#ffffff' : '#1e293b',
        },
        text: {
            primary: mode === 'light' ? '#1e293b' : '#f1f5f9',
            secondary: mode === 'light' ? '#64748b' : '#94a3b8',
        },
        divider: mode === 'light' ? 'rgba(0, 0, 0, 0.08)' : 'rgba(255, 255, 255, 0.08)',
    },
    typography: {
        fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif",
        h1: { fontWeight: 700, fontSize: '2.5rem' },
        h2: { fontWeight: 700, fontSize: '2rem' },
        h3: { fontWeight: 600, fontSize: '1.75rem' },
        h4: { fontWeight: 600, fontSize: '1.5rem' },
        h5: { fontWeight: 600, fontSize: '1.25rem' },
        h6: { fontWeight: 600, fontSize: '1rem' },
        button: { textTransform: 'none' as const, fontWeight: 500 },
    },
    shape: {
        borderRadius: 12,
    },
});

export const getTheme = (mode: 'light' | 'dark'): Theme => {
    const baseTheme = createTheme(getDesignTokens(mode));

    return createTheme(baseTheme, {
        components: {
            MuiButton: {
                styleOverrides: {
                    root: {
                        borderRadius: 10,
                        padding: '10px 20px',
                        fontWeight: 600,
                    },
                    contained: {
                        boxShadow: 'none',
                        '&:hover': {
                            boxShadow: '0 4px 12px rgba(37, 99, 235, 0.25)',
                        },
                    },
                },
            },
            MuiCard: {
                styleOverrides: {
                    root: {
                        borderRadius: 16,
                        boxShadow: mode === 'light'
                            ? '0 1px 3px rgba(0,0,0,0.08), 0 1px 2px rgba(0,0,0,0.12)'
                            : '0 1px 3px rgba(0,0,0,0.4), 0 1px 2px rgba(0,0,0,0.3)',
                        border: mode === 'light'
                            ? '1px solid rgba(226, 232, 240, 0.8)'
                            : '1px solid rgba(51, 65, 85, 0.5)',
                    },
                },
            },
            MuiPaper: {
                styleOverrides: {
                    root: {
                        backgroundImage: 'none',
                    },
                },
            },
            MuiTextField: {
                defaultProps: {
                    variant: 'outlined',
                    size: 'small',
                },
                styleOverrides: {
                    root: {
                        '& .MuiOutlinedInput-root': {
                            borderRadius: 10,
                        },
                    },
                },
            },
            MuiDrawer: {
                styleOverrides: {
                    paper: {
                        backgroundColor: baseTheme.palette.background.paper,
                        borderRight: mode === 'light'
                            ? '1px solid rgba(226, 232, 240, 0.8)'
                            : '1px solid rgba(51, 65, 85, 0.5)',
                    },
                },
            },
            MuiAppBar: {
                styleOverrides: {
                    root: {
                        backgroundColor: mode === 'light'
                            ? 'rgba(255, 255, 255, 0.9)'
                            : 'rgba(30, 41, 59, 0.9)',
                        backdropFilter: 'blur(10px)',
                        boxShadow: mode === 'light'
                            ? '0 1px 3px rgba(0,0,0,0.05)'
                            : '0 1px 3px rgba(0,0,0,0.3)',
                    },
                },
            },
            MuiListItemButton: {
                styleOverrides: {
                    root: {
                        borderRadius: 10,
                        margin: '4px 8px',
                        '&.Mui-selected': {
                            backgroundColor: mode === 'light'
                                ? 'rgba(37, 99, 235, 0.1)'
                                : 'rgba(37, 99, 235, 0.2)',
                            '&:hover': {
                                backgroundColor: mode === 'light'
                                    ? 'rgba(37, 99, 235, 0.15)'
                                    : 'rgba(37, 99, 235, 0.25)',
                            },
                        },
                    },
                },
            },
        },
    });
};

export const lightTheme = getTheme('light');
export const darkTheme = getTheme('dark');
