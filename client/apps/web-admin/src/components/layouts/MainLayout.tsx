import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Box, Toolbar } from '@mui/material';
import { AppBar } from './AppBar';
import { Sidebar } from './Sidebar';

const DRAWER_WIDTH = 260;

export const MainLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState(false);

    const handleDrawerToggle = () => {
        setSidebarOpen(!sidebarOpen);
    };

    return (
        <Box sx={{ display: 'flex' }}>
            <AppBar
                onMenuClick={handleDrawerToggle}
                drawerWidth={DRAWER_WIDTH}
                open={sidebarOpen}
            />
            <Sidebar
                open={sidebarOpen}
                onClose={handleDrawerToggle}
                drawerWidth={DRAWER_WIDTH}
            />
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    p: 3,
                    width: {
                        sm: sidebarOpen
                            ? `calc(100% - ${DRAWER_WIDTH}px)`
                            : '100%',
                    },
                    minHeight: '100vh',
                    backgroundColor: 'background.default',
                }}
            >
                <Toolbar />
                <Outlet />
            </Box>
        </Box>
    );
};
