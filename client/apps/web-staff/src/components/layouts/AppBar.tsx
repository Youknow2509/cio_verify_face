import { useState } from 'react';
import {
    AppBar as MuiAppBar,
    Toolbar,
    IconButton,
    Typography,
    Box,
    Badge,
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
    Divider,
} from '@mui/material';
import {
    Menu as MenuIcon,
    Notifications,
    CheckCircle,
    Info,
    Warning,
} from '@mui/icons-material';
import { ThemeToggle } from '../ThemeToggle';

interface AppBarProps {
    onMenuClick: () => void;
}

interface Notification {
    id: string;
    title: string;
    message: string;
    type: 'success' | 'info' | 'warning';
    time: string;
    read: boolean;
}

const mockNotifications: Notification[] = [
    {
        id: '1',
        title: 'Yêu cầu được duyệt',
        message: 'Yêu cầu cập nhật khuôn mặt của bạn đã được phê duyệt',
        type: 'success',
        time: '2 giờ trước',
        read: false,
    },
    {
        id: '2',
        title: 'Nhắc nhở chấm công',
        message: 'Bạn chưa chấm công ra hôm nay',
        type: 'warning',
        time: '5 giờ trước',
        read: false,
    },
    {
        id: '3',
        title: 'Thông báo hệ thống',
        message: 'Hệ thống sẽ bảo trì vào 00:00 ngày mai',
        type: 'info',
        time: '1 ngày trước',
        read: true,
    },
];

const typeIcons = {
    success: <CheckCircle color="success" fontSize="small" />,
    info: <Info color="info" fontSize="small" />,
    warning: <Warning color="warning" fontSize="small" />,
};

export const AppBar: React.FC<AppBarProps> = ({ onMenuClick }) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [notifications] = useState<Notification[]>(mockNotifications);

    const unreadCount = notifications.filter((n) => !n.read).length;

    const handleOpenMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleCloseMenu = () => {
        setAnchorEl(null);
    };

    return (
        <>
            <MuiAppBar
                position="fixed"
                elevation={0}
                sx={{
                    zIndex: (theme) => theme.zIndex.drawer + 1,
                }}
            >
                <Toolbar>
                    <IconButton
                        color="inherit"
                        onClick={onMenuClick}
                        edge="start"
                        sx={{ mr: 2, color: 'text.primary' }}
                    >
                        <MenuIcon />
                    </IconButton>

                    <Typography
                        variant="h6"
                        noWrap
                        component="div"
                        sx={{ flexGrow: 1, color: 'text.primary', fontWeight: 600 }}
                    >
                        Staff Portal
                    </Typography>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <ThemeToggle />
                        <IconButton
                            size="small"
                            sx={{ color: 'text.primary' }}
                            onClick={handleOpenMenu}
                        >
                            <Badge badgeContent={unreadCount} color="error">
                                <Notifications />
                            </Badge>
                        </IconButton>
                    </Box>
                </Toolbar>
            </MuiAppBar>

            {/* Notification Menu */}
            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleCloseMenu}
                PaperProps={{
                    sx: {
                        width: 360,
                        maxHeight: 400,
                        mt: 1.5,
                        borderRadius: 2,
                        boxShadow: '0 10px 40px rgba(0,0,0,0.15)',
                    },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
                <Box sx={{ px: 2, py: 1.5 }}>
                    <Typography variant="subtitle1" fontWeight="700">
                        Thông báo
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        Bạn có {unreadCount} thông báo chưa đọc
                    </Typography>
                </Box>
                <Divider />
                {notifications.length === 0 ? (
                    <Box sx={{ py: 4, textAlign: 'center' }}>
                        <Typography color="text.secondary">
                            Không có thông báo mới
                        </Typography>
                    </Box>
                ) : (
                    notifications.map((notification) => (
                        <MenuItem
                            key={notification.id}
                            onClick={handleCloseMenu}
                            sx={{
                                py: 1.5,
                                px: 2,
                                bgcolor: notification.read ? 'transparent' : 'action.hover',
                                '&:hover': {
                                    bgcolor: 'action.selected',
                                },
                            }}
                        >
                            <ListItemIcon sx={{ minWidth: 36 }}>
                                {typeIcons[notification.type]}
                            </ListItemIcon>
                            <ListItemText
                                primary={
                                    <Typography variant="body2" fontWeight={notification.read ? 400 : 600}>
                                        {notification.title}
                                    </Typography>
                                }
                                secondary={
                                    <>
                                        <Typography
                                            variant="caption"
                                            color="text.secondary"
                                            sx={{
                                                display: 'block',
                                                overflow: 'hidden',
                                                textOverflow: 'ellipsis',
                                                whiteSpace: 'nowrap',
                                            }}
                                        >
                                            {notification.message}
                                        </Typography>
                                        <Typography variant="caption" color="text.disabled">
                                            {notification.time}
                                        </Typography>
                                    </>
                                }
                            />
                        </MenuItem>
                    ))
                )}
                <Divider />
                <MenuItem
                    onClick={handleCloseMenu}
                    sx={{ justifyContent: 'center', py: 1.5 }}
                >
                    <Typography variant="body2" color="primary" fontWeight="500">
                        Xem tất cả thông báo
                    </Typography>
                </MenuItem>
            </Menu>
        </>
    );
};
