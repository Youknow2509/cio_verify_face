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
} from '@mui/material';
import {
    Dashboard,
    Assessment,
    Face,
    Person,
    Logout,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';

interface SidebarProps {
    onNavigate?: () => void;
}

const menuItems = [
    { text: 'Dashboard', icon: <Dashboard />, path: '/dashboard' },
    { text: 'Báo cáo tháng', icon: <Assessment />, path: '/reports/monthly' },
    { text: 'Cập nhật khuôn mặt', icon: <Face />, path: '/face-update' },
    { text: 'Thông tin cá nhân', icon: <Person />, path: '/profile' },
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
            <Box sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 2 }}>
                <Box
                    sx={{
                        width: 44,
                        height: 44,
                        borderRadius: 2,
                        background: 'linear-gradient(135deg, #2563eb, #7c3aed)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    <Face sx={{ color: 'white', fontSize: 24 }} />
                </Box>
                <Box>
                    <Typography variant="h6" fontWeight="700" color="primary">
                        CIO Staff
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        Employee Portal
                    </Typography>
                </Box>
            </Box>

            <Divider />

            {/* User Info */}
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
                        sx={{ width: 40, height: 40, bgcolor: 'primary.main' }}
                    >
                        {user?.full_name?.charAt(0) || 'U'}
                    </Avatar>
                    <Box sx={{ overflow: 'hidden' }}>
                        <Typography variant="body2" fontWeight="600" noWrap>
                            {user?.full_name || 'Nhân viên'}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap>
                            {user?.department || 'Phòng ban'}
                        </Typography>
                    </Box>
                </Box>
            </Box>

            <Divider />

            {/* Menu Items */}
            <List sx={{ flex: 1, py: 2 }}>
                {menuItems.map((item) => (
                    <ListItem key={item.path} disablePadding>
                        <ListItemButton
                            selected={location.pathname === item.path}
                            onClick={() => handleNavigation(item.path)}
                            sx={{ mx: 1, borderRadius: 2 }}
                        >
                            <ListItemIcon
                                sx={{
                                    color: location.pathname === item.path ? 'primary.main' : 'text.secondary',
                                    minWidth: 40,
                                }}
                            >
                                {item.icon}
                            </ListItemIcon>
                            <ListItemText
                                primary={item.text}
                                primaryTypographyProps={{
                                    fontWeight: location.pathname === item.path ? 600 : 400,
                                    color: location.pathname === item.path ? 'primary.main' : 'text.primary',
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
