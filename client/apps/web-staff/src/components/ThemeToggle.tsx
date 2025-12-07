import { IconButton, Tooltip } from '@mui/material';
import { DarkMode, LightMode } from '@mui/icons-material';
import { useThemeStore } from '@/stores/themeStore';

export const ThemeToggle: React.FC = () => {
    const { mode, toggleTheme } = useThemeStore();

    return (
        <Tooltip title={mode === 'light' ? 'Chế độ tối' : 'Chế độ sáng'}>
            <IconButton onClick={toggleTheme} color="inherit" size="small">
                {mode === 'light' ? <DarkMode /> : <LightMode />}
            </IconButton>
        </Tooltip>
    );
};
