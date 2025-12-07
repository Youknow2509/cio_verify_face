import { useMemo } from 'react';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { RouterProvider } from 'react-router-dom';
import { getTheme } from '@/theme';
import { useThemeStore } from '@/stores/themeStore';
import { router } from '@/routes';

function App() {
    const { mode } = useThemeStore();
    const theme = useMemo(() => getTheme(mode), [mode]);

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <RouterProvider router={router} />
        </ThemeProvider>
    );
}

export default App;
