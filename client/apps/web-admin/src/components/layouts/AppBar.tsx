import {
    AppBar as MuiAppBar,
    Toolbar,
    IconButton,
    Typography,
    Avatar,
    Menu,
    MenuItem,
    Box,
    Badge,
} from '@mui/material';
import {
    Menu as MenuIcon,
    Notifications as NotificationsIcon,
    AccountCircle,
} from '@mui/icons-material';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';

interface AppBarProps {
    onMenuClick: () => void;
    drawerWidth: number;
    open: boolean;
}

export const AppBar: React.FC<AppBarProps> = ({
    onMenuClick,
    drawerWidth,
    open,
}) => {
    const navigate = useNavigate();
    const { user, clearAuth } = useAuthStore();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = () => {
        clearAuth();
        navigate('/login');
    };

    return (
        <MuiAppBar
            position="fixed"
            sx={{
                width: { sm: open ? `calc(100% - ${drawerWidth}px)` : '100%' },
                ml: { sm: open ? `${drawerWidth}px` : 0 },
            }}
        >
            <Toolbar>
                <IconButton
                    color="inherit"
                    edge="start"
                    onClick={onMenuClick}
                    sx={{ mr: 2 }}
                    aria-label="toggle sidebar"
                >
                    <MenuIcon />
                </IconButton>
                <Typography
                    variant="h6"
                    noWrap
                    component="div"
                    sx={{ flexGrow: 1 }}
                >
                    Hệ thống Chấm công Khuôn mặt
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <IconButton color="inherit">
                        <Badge badgeContent={3} color="error">
                            <NotificationsIcon />
                        </Badge>
                    </IconButton>
                    <IconButton onClick={handleMenu} color="inherit">
                        {user?.avatar_url ? (
                            <Avatar
                                src={user.avatar_url}
                                sx={{ width: 32, height: 32 }}
                            />
                        ) : (
                            <AccountCircle />
                        )}
                    </IconButton>
                    <Menu
                        anchorEl={anchorEl}
                        open={Boolean(anchorEl)}
                        onClose={handleClose}
                    >
                        <MenuItem disabled>
                            <Typography variant="body2">
                                {user?.full_name}
                            </Typography>
                        </MenuItem>
                        <MenuItem
                            onClick={() => {
                                handleClose();
                                navigate('/settings');
                            }}
                        >
                            Cài đặt
                        </MenuItem>
                        <MenuItem onClick={handleLogout}>Đăng xuất</MenuItem>
                    </Menu>
                </Box>
            </Toolbar>
        </MuiAppBar>
    );
};
