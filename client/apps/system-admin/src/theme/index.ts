import { createTheme, alpha } from '@mui/material/styles';

const primaryColor = '#6366f1'; // Indigo - distinguishes from web-admin blue
const secondaryColor = '#ec4899'; // Pink accent

export const lightTheme = createTheme({
    palette: {
        mode: 'light',
        primary: {
            main: primaryColor,
            light: '#818cf8',
            dark: '#4f46e5',
        },
        secondary: {
            main: secondaryColor,
        },
        background: {
            default: '#f8fafc',
            paper: '#ffffff',
        },
        success: { main: '#10b981' },
        warning: { main: '#f59e0b' },
        error: { main: '#ef4444' },
        info: { main: '#06b6d4' },
    },
    typography: {
        fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, sans-serif",
        h1: { fontWeight: 700 },
        h2: { fontWeight: 700 },
        h3: { fontWeight: 600 },
        h4: { fontWeight: 600 },
        h5: { fontWeight: 600 },
        h6: { fontWeight: 600 },
    },
    shape: { borderRadius: 12 },
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    fontWeight: 600,
                    borderRadius: 8,
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    borderRadius: 16,
                    boxShadow: '0 1px 3px rgba(0,0,0,0.08)',
                    border: '1px solid rgba(0,0,0,0.05)',
                },
            },
        },
        MuiAppBar: {
            styleOverrides: {
                root: {
                    backgroundColor: '#ffffff',
                    boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
                },
            },
        },
        MuiDrawer: {
            styleOverrides: {
                paper: {
                    backgroundColor: '#ffffff',
                    borderRight: '1px solid rgba(0,0,0,0.08)',
                },
            },
        },
        MuiChip: {
            styleOverrides: {
                root: { fontWeight: 500 },
            },
        },
    },
});

export const darkTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#818cf8',
            light: '#a5b4fc',
            dark: '#6366f1',
        },
        secondary: {
            main: '#f472b6',
        },
        background: {
            default: '#0f172a',
            paper: '#1e293b',
        },
        success: { main: '#34d399' },
        warning: { main: '#fbbf24' },
        error: { main: '#f87171' },
        info: { main: '#22d3ee' },
    },
    typography: {
        fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, sans-serif",
        h1: { fontWeight: 700 },
        h2: { fontWeight: 700 },
        h3: { fontWeight: 600 },
        h4: { fontWeight: 600 },
        h5: { fontWeight: 600 },
        h6: { fontWeight: 600 },
    },
    shape: { borderRadius: 12 },
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    fontWeight: 600,
                    borderRadius: 8,
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    borderRadius: 16,
                    boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
                    border: '1px solid rgba(255,255,255,0.05)',
                    backgroundImage: 'none',
                },
            },
        },
        MuiAppBar: {
            styleOverrides: {
                root: {
                    backgroundColor: '#1e293b',
                    boxShadow: '0 1px 3px rgba(0,0,0,0.3)',
                },
            },
        },
        MuiDrawer: {
            styleOverrides: {
                paper: {
                    backgroundColor: '#1e293b',
                    borderRight: '1px solid rgba(255,255,255,0.05)',
                },
            },
        },
        MuiChip: {
            styleOverrides: {
                root: { fontWeight: 500 },
            },
        },
    },
});
