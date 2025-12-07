import { Box, Typography } from '@mui/material';
import { Face } from '@mui/icons-material';

interface AuthLayoutProps {
    children: React.ReactNode;
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({ children }) => {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: (theme) =>
                    theme.palette.mode === 'dark'
                        ? 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)'
                        : 'linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 50%, #ddd6fe 100%)',
                position: 'relative',
                overflow: 'hidden',
                '&::before': {
                    content: '""',
                    position: 'absolute',
                    width: '400px',
                    height: '400px',
                    background: 'linear-gradient(135deg, rgba(37, 99, 235, 0.3), rgba(124, 58, 237, 0.3))',
                    borderRadius: '50%',
                    top: '-100px',
                    right: '-100px',
                    filter: 'blur(80px)',
                },
                '&::after': {
                    content: '""',
                    position: 'absolute',
                    width: '300px',
                    height: '300px',
                    background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.3), rgba(6, 182, 212, 0.3))',
                    borderRadius: '50%',
                    bottom: '-50px',
                    left: '-50px',
                    filter: 'blur(80px)',
                },
            }}
        >
            {/* Logo watermark */}
            <Box
                sx={{
                    position: 'absolute',
                    top: 24,
                    left: 24,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1,
                    zIndex: 1,
                }}
            >
                <Box
                    sx={{
                        width: 40,
                        height: 40,
                        borderRadius: 2,
                        background: 'linear-gradient(135deg, #2563eb, #7c3aed)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    <Face sx={{ color: 'white', fontSize: 24 }} />
                </Box>
                <Typography variant="h6" fontWeight="700" color="text.primary">
                    CIO Staff
                </Typography>
            </Box>

            {/* Content */}
            <Box sx={{ position: 'relative', zIndex: 1 }}>
                {children}
            </Box>
        </Box>
    );
};
