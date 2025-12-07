import { useLocation, useNavigate } from 'react-router-dom';
import {
    Box,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Typography,
    Divider,
    Avatar,
    Toolbar,
    Chip,
} from '@mui/material';
import {
    Dashboard,
    Business,
    MonitorHeart,
    Settings,
    Logout,
    AdminPanelSettings,
    History,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';

interface SidebarProps {
    onNavigate?: () => void;
}

const menuItems = [
    { text: 'Dashboard', icon: <Dashboard />, path: '/dashboard' },
    { text: 'Quản lý Công ty', icon: <Business />, path: '/companies' },
    { text: 'Giám sát Hệ thống', icon: <MonitorHeart />, path: '/monitoring' },
    { text: 'Audit Log', icon: <History />, path: '/audit-log' },
    { text: 'Cài đặt', icon: <Settings />, path: '/settings' },
];

export const Sidebar: React.FC<SidebarProps> = ({ onNavigate }) => {
    const location = useLocation();
    const navigate = useNavigate();
    const { user, clearAuth } = useAuthStore();

    const handleNavigation = (path: string) => {
        navigate(path);
        onNavigate?.();
    };

    const handleLogout = () => {
        clearAuth();
        navigate('/login');
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Toolbar />

            {/* Logo */}
            <Box sx={{ p: 2.5, display: 'flex', alignItems: 'center', gap: 2 }}>
                <Box
                    sx={{
                        width: 48,
                        height: 48,
                        borderRadius: 2,
                        background: 'linear-gradient(135deg, #6366f1, #ec4899)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    <AdminPanelSettings sx={{ color: 'white', fontSize: 28 }} />
                </Box>
                <Box>
                    <Typography variant="h6" fontWeight="700" color="primary">
                        System Admin
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        Platform Control
                    </Typography>
                </Box>
            </Box>

            <Divider />

            {/* Admin Info */}
            <Box sx={{ p: 2 }}>
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        p: 1.5,
                        borderRadius: 2,
                        bgcolor: 'action.hover',
                    }}
                >
                    <Avatar
                        src={user?.avatar}
                        sx={{ width: 40, height: 40, bgcolor: 'secondary.main' }}
                    >
                        {user?.full_name?.charAt(0) || 'A'}
                    </Avatar>
                    <Box sx={{ overflow: 'hidden', flex: 1 }}>
                        <Typography variant="body2" fontWeight="600" noWrap>
                            {user?.full_name || 'System Admin'}
                        </Typography>
                        <Chip
                            label={user?.role === 'super_admin' ? 'Super Admin' : 'Admin'}
                            size="small"
                            color="secondary"
                            sx={{ height: 20, fontSize: '0.65rem' }}
                        />
                    </Box>
                </Box>
            </Box>

            <Divider />

            {/* Menu Items */}
            <List sx={{ flex: 1, py: 2 }}>
                {menuItems.map((item) => (
                    <ListItem key={item.path} disablePadding>
                        <ListItemButton
                            selected={location.pathname.startsWith(item.path)}
                            onClick={() => handleNavigation(item.path)}
                            sx={{ mx: 1, borderRadius: 2 }}
                        >
                            <ListItemIcon
                                sx={{
                                    color: location.pathname.startsWith(item.path) ? 'primary.main' : 'text.secondary',
                                    minWidth: 40,
                                }}
                            >
                                {item.icon}
                            </ListItemIcon>
                            <ListItemText
                                primary={item.text}
                                primaryTypographyProps={{
                                    fontWeight: location.pathname.startsWith(item.path) ? 600 : 400,
                                    color: location.pathname.startsWith(item.path) ? 'primary.main' : 'text.primary',
                                }}
                            />
                        </ListItemButton>
                    </ListItem>
                ))}
            </List>

            <Divider />

            {/* Logout */}
            <List sx={{ pb: 2 }}>
                <ListItem disablePadding>
                    <ListItemButton onClick={handleLogout} sx={{ mx: 1, borderRadius: 2 }}>
                        <ListItemIcon sx={{ color: 'error.main', minWidth: 40 }}>
                            <Logout />
                        </ListItemIcon>
                        <ListItemText
                            primary="Đăng xuất"
                            primaryTypographyProps={{ color: 'error.main' }}
                        />
                    </ListItemButton>
                </ListItem>
            </List>
        </Box>
    );
};
