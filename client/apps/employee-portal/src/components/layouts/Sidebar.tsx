import {
    Drawer,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Toolbar,
    Box,
    useMediaQuery,
    useTheme,
    Divider,
} from '@mui/material';
import {
    Dashboard as DashboardIcon,
    AccessTime as AccessTimeIcon,
    CalendarMonth as CalendarIcon,
    Schedule as ScheduleIcon,
    Person as PersonIcon,
    Article as ArticleIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

const drawerWidth = 280;

interface SidebarProps {
    mobileOpen: boolean;
    onMobileClose: () => void;
}

const menuItems = [
    { text: 'Bảng điều khiển', icon: <DashboardIcon />, path: '/dashboard' },
    { text: 'Chấm công', icon: <AccessTimeIcon />, path: '/attendance' },
    { text: 'Tổng hợp ngày', icon: <CalendarIcon />, path: '/daily-summary' },
    { text: 'Ca làm việc', icon: <ScheduleIcon />, path: '/shifts' },
    { text: 'Hồ sơ cá nhân', icon: <PersonIcon />, path: '/profile' },
    { text: 'Xuất báo cáo', icon: <ArticleIcon />, path: '/export' },
];

export const Sidebar: React.FC<SidebarProps> = ({ mobileOpen, onMobileClose }) => {
    const navigate = useNavigate();
    const location = useLocation();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));

    const drawer = (
        <Box>
            <Toolbar>
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        py: 1,
                    }}
                >
                    <Box
                        sx={{
                            width: 40,
                            height: 40,
                            borderRadius: 2,
                            background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            fontSize: '1.5rem',
                        }}
                    >
                        👤
                    </Box>
                    <Box>
                        <Box sx={{ fontWeight: 700, fontSize: '1.1rem' }}>CIO</Box>
                        <Box sx={{ fontSize: '0.75rem', color: 'text.secondary' }}>
                            Employee Portal
                        </Box>
                    </Box>
                </Box>
            </Toolbar>
            <Divider sx={{ borderColor: 'rgba(71, 85, 105, 0.3)' }} />
            <List sx={{ px: 2, pt: 2 }}>
                {menuItems.map((item) => {
                    const isActive = location.pathname === item.path;
                    return (
                        <ListItem key={item.text} disablePadding sx={{ mb: 1 }}>
                            <ListItemButton
                                onClick={() => {
                                    navigate(item.path);
                                    if (isMobile) onMobileClose();
                                }}
                                sx={{
                                    borderRadius: 2,
                                    py: 1.5,
                                    background: isActive
                                        ? 'linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%)'
                                        : 'transparent',
                                    border: isActive
                                        ? '1px solid rgba(59, 130, 246, 0.3)'
                                        : '1px solid transparent',
                                    '&:hover': {
                                        background: isActive
                                            ? 'linear-gradient(135deg, rgba(59, 130, 246, 0.3) 0%, rgba(139, 92, 246, 0.3) 100%)'
                                            : 'rgba(71, 85, 105, 0.2)',
                                    },
                                    transition: 'all 0.2s ease',
                                }}
                            >
                                <ListItemIcon
                                    sx={{
                                        color: isActive ? '#3b82f6' : 'text.secondary',
                                        minWidth: 40,
                                    }}
                                >
                                    {item.icon}
                                </ListItemIcon>
                                <ListItemText
                                    primary={item.text}
                                    primaryTypographyProps={{
                                        fontWeight: isActive ? 600 : 400,
                                        color: isActive ? '#fff' : 'text.primary',
                                    }}
                                />
                            </ListItemButton>
                        </ListItem>
                    );
                })}
            </List>
        </Box>
    );

    return (
        <>
            {isMobile ? (
                <Drawer
                    variant="temporary"
                    open={mobileOpen}
                    onClose={onMobileClose}
                    ModalProps={{ keepMounted: true }}
                    sx={{
                        '& .MuiDrawer-paper': {
                            width: drawerWidth,
                            boxSizing: 'border-box',
                            background: 'rgba(15, 23, 42, 0.95)',
                            backdropFilter: 'blur(20px)',
                            borderRight: '1px solid rgba(71, 85, 105, 0.3)',
                        },
                    }}
                >
                    {drawer}
                </Drawer>
            ) : (
                <Drawer
                    variant="permanent"
                    sx={{
                        width: drawerWidth,
                        flexShrink: 0,
                        '& .MuiDrawer-paper': {
                            width: drawerWidth,
                            boxSizing: 'border-box',
                            background: 'rgba(15, 23, 42, 0.95)',
                            backdropFilter: 'blur(20px)',
                            borderRight: '1px solid rgba(71, 85, 105, 0.3)',
                        },
                    }}
                >
                    {drawer}
                </Drawer>
            )}
        </>
    );
};
