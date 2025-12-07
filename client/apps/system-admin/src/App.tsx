import { ThemeProvider, CssBaseline } from '@mui/material';
import { RouterProvider } from 'react-router-dom';
import { useThemeStore } from '@/stores/themeStore';
import { lightTheme, darkTheme } from '@/theme';
import { router } from '@/routes';
import { ErrorBoundary } from '@/components/ErrorBoundary';
import { ToastContainer } from '@/components/Toast';

export const App: React.FC = () => {
  const { mode } = useThemeStore();
  const theme = mode === 'light' ? lightTheme : darkTheme;

  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <RouterProvider router={router} />
        <ToastContainer />
      </ThemeProvider>
    </ErrorBoundary>
  );
};
