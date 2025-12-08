import {
    Drawer,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Toolbar,
    Box,
    Typography,
} from '@mui/material';
import {
    People,
    Devices,
    Schedule,
    Assessment,
    Settings,
    Face,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

interface SidebarProps {
    open: boolean;
    onClose: () => void;
    drawerWidth: number;
}

const menuItems = [
    { text: 'Báo cáo', icon: <Assessment />, path: '/reports/daily' },
    { text: 'Nhân viên', icon: <People />, path: '/employees' },
    { text: 'Thiết bị', icon: <Devices />, path: '/devices' },
    { text: 'Ca làm việc', icon: <Schedule />, path: '/shifts' },
    {
        text: 'Yêu cầu cập nhật khuôn mặt',
        icon: <Face />,
        path: '/profile-update-requests',
    },
    { text: 'Cài đặt', icon: <Settings />, path: '/settings' },
];

export const Sidebar: React.FC<SidebarProps> = ({
    open,
    onClose,
    drawerWidth,
}) => {
    const navigate = useNavigate();
    const location = useLocation();

    const drawer = (
        <Box>
            <Toolbar
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <Typography variant="h6" color="primary" fontWeight="bold">
                    Face Attendance
                </Typography>
            </Toolbar>
            <List>
                {menuItems.map((item) => (
                    <ListItem key={item.text} disablePadding>
                        <ListItemButton
                            selected={location.pathname.startsWith(item.path)}
                            onClick={() => {
                                navigate(item.path);
                                onClose();
                            }}
                        >
                            <ListItemIcon>{item.icon}</ListItemIcon>
                            <ListItemText primary={item.text} />
                        </ListItemButton>
                    </ListItem>
                ))}
            </List>
        </Box>
    );

    return (
        <>
            <Drawer
                variant="temporary"
                open={open}
                onClose={onClose}
                ModalProps={{ keepMounted: true }}
                sx={{
                    display: { xs: 'block', sm: 'none' },
                    '& .MuiDrawer-paper': {
                        boxSizing: 'border-box',
                        width: drawerWidth,
                    },
                }}
            >
                {drawer}
            </Drawer>
            <Drawer
                variant="persistent"
                open={open}
                sx={{
                    display: { xs: 'none', sm: 'block' },
                    '& .MuiDrawer-paper': {
                        boxSizing: 'border-box',
                        width: drawerWidth,
                    },
                }}
            >
                {drawer}
            </Drawer>
        </>
    );
};
