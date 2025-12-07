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
    Chip,
} from '@mui/material';
import {
    Menu as MenuIcon,
    Notifications,
    Warning,
    Error as ErrorIcon,
    CheckCircle,
} from '@mui/icons-material';
import { ThemeToggle } from '../ThemeToggle';

interface AppBarProps {
    onMenuClick: () => void;
}

interface SystemAlert {
    id: string;
    title: string;
    message: string;
    type: 'error' | 'warning' | 'success';
    time: string;
}

const mockAlerts: SystemAlert[] = [
    {
        id: '1',
        title: 'Database High Load',
        message: 'PostgreSQL connections at 85% capacity',
        type: 'warning',
        time: '5 phút trước',
    },
    {
        id: '2',
        title: 'Company Subscription Expiring',
        message: 'Tech Corp subscription expires in 3 days',
        type: 'warning',
        time: '1 giờ trước',
    },
    {
        id: '3',
        title: 'Service Recovered',
        message: 'Attendance Service is back online',
        type: 'success',
        time: '2 giờ trước',
    },
];

const typeIcons = {
    error: <ErrorIcon color="error" fontSize="small" />,
    warning: <Warning color="warning" fontSize="small" />,
    success: <CheckCircle color="success" fontSize="small" />,
};

export const AppBar: React.FC<AppBarProps> = ({ onMenuClick }) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const activeAlerts = mockAlerts.filter((a) => a.type !== 'success').length;

    return (
        <>
            <MuiAppBar
                position="fixed"
                elevation={0}
                sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
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
                        sx={{ flexGrow: 1, color: 'text.primary', fontWeight: 600 }}
                    >
                        System Administration
                    </Typography>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Chip
                            label="Production"
                            size="small"
                            color="success"
                            variant="outlined"
                            sx={{ fontWeight: 500 }}
                        />
                        <ThemeToggle />
                        <IconButton
                            size="small"
                            sx={{ color: 'text.primary' }}
                            onClick={(e) => setAnchorEl(e.currentTarget)}
                        >
                            <Badge badgeContent={activeAlerts} color="error">
                                <Notifications />
                            </Badge>
                        </IconButton>
                    </Box>
                </Toolbar>
            </MuiAppBar>

            {/* Alerts Menu */}
            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={() => setAnchorEl(null)}
                PaperProps={{
                    sx: {
                        width: 380,
                        maxHeight: 400,
                        mt: 1.5,
                        borderRadius: 2,
                    },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
                <Box sx={{ px: 2, py: 1.5 }}>
                    <Typography variant="subtitle1" fontWeight="700">
                        System Alerts
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        {activeAlerts} thông báo cần xử lý
                    </Typography>
                </Box>
                <Divider />
                {mockAlerts.map((alert) => (
                    <MenuItem
                        key={alert.id}
                        onClick={() => setAnchorEl(null)}
                        sx={{ py: 1.5, px: 2 }}
                    >
                        <ListItemIcon sx={{ minWidth: 36 }}>
                            {typeIcons[alert.type]}
                        </ListItemIcon>
                        <ListItemText
                            primary={
                                <Typography variant="body2" fontWeight="600">
                                    {alert.title}
                                </Typography>
                            }
                            secondary={
                                <>
                                    <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                                        {alert.message}
                                    </Typography>
                                    <Typography variant="caption" color="text.disabled">
                                        {alert.time}
                                    </Typography>
                                </>
                            }
                        />
                    </MenuItem>
                ))}
                <Divider />
                <MenuItem sx={{ justifyContent: 'center', py: 1.5 }}>
                    <Typography variant="body2" color="primary" fontWeight="500">
                        View All Alerts
                    </Typography>
                </MenuItem>
            </Menu>
        </>
    );
};
