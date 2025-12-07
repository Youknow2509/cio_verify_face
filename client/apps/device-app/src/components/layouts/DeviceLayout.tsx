import { Box } from '@mui/material';

interface DeviceLayoutProps {
    children: React.ReactNode;
}

export const DeviceLayout: React.FC<DeviceLayoutProps> = ({ children }) => {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                width: '100%',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 50%, #cbd5e1 100%)',
                position: 'relative',
                overflow: 'hidden',
                // Decorative circles
                '&::before': {
                    content: '""',
                    position: 'absolute',
                    width: '500px',
                    height: '500px',
                    background: 'radial-gradient(circle, rgba(37, 99, 235, 0.08) 0%, transparent 70%)',
                    borderRadius: '50%',
                    top: '-150px',
                    right: '-150px',
                    pointerEvents: 'none',
                },
                '&::after': {
                    content: '""',
                    position: 'absolute',
                    width: '400px',
                    height: '400px',
                    background: 'radial-gradient(circle, rgba(16, 185, 129, 0.06) 0%, transparent 70%)',
                    borderRadius: '50%',
                    bottom: '-100px',
                    left: '-100px',
                    pointerEvents: 'none',
                },
            }}
        >
            {children}
        </Box>
    );
};
