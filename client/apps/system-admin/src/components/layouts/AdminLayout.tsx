import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Box, Toolbar, useMediaQuery, useTheme, Drawer } from '@mui/material';
import { AppBar } from './AppBar';
import { Sidebar } from './Sidebar';

const DRAWER_WIDTH = 280;

export const AdminLayout: React.FC = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const [sidebarOpen, setSidebarOpen] = useState(!isMobile);

    const handleDrawerToggle = () => {
        setSidebarOpen(!sidebarOpen);
    };

    return (
        <Box sx={{ display: 'flex' }}>
            <AppBar onMenuClick={handleDrawerToggle} />

            <Drawer
                variant={isMobile ? 'temporary' : 'persistent'}
                open={sidebarOpen}
                onClose={handleDrawerToggle}
                ModalProps={{ keepMounted: true }}
                PaperProps={{
                    sx: {
                        width: DRAWER_WIDTH,
                        borderRight: '1px solid',
                        borderColor: 'divider',
                    },
                }}
                sx={{
                    width: isMobile ? 0 : (sidebarOpen ? DRAWER_WIDTH : 0),
                    flexShrink: 0,
                    transition: 'width 0.2s ease',
                }}
            >
                <Sidebar onNavigate={isMobile ? handleDrawerToggle : undefined} />
            </Drawer>

            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    p: 3,
                    minHeight: '100vh',
                    maxWidth: '100%',
                    overflow: 'auto',
                    backgroundColor: 'background.default',
                }}
            >
                <Toolbar />
                <Outlet />
            </Box>
        </Box>
    );
};
