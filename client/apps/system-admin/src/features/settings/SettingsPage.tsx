import { useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    Button,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    ListItemButton,
    Divider,
} from '@mui/material';
import {
    Settings as SettingsIcon,
    Security,
    Notifications,
    Storage,
    Language,
    CreditCard,
    ChevronRight,
} from '@mui/icons-material';

const settingLinks = [
    {
        title: 'Gói Dịch vụ (Service Plans)',
        description: 'Quản lý các gói cước và giới hạn tài nguyên',
        icon: <CreditCard />,
        path: '/settings/plans',
    },
    {
        title: 'System Security',
        description: 'Cấu hình bảo mật, MFA và chính sách mật khẩu',
        icon: <Security />,
        path: '/settings/security',
    },
    {
        title: 'System Notifications',
        description: 'Cấu hình email server và mẫu thông báo',
        icon: <Notifications />,
        path: '/settings/notifications',
    },
    {
        title: 'Storage Configuration',
        description: 'Kết nối MinIO/S3 và chính sách lưu trữ',
        icon: <Storage />,
        path: '/settings/storage',
    },
];

export const SettingsPage: React.FC = () => {
    const navigate = useNavigate();

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Typography variant="h4" fontWeight="700" mb={1}>
                Cài đặt Hệ thống
            </Typography>
            <Typography variant="body1" color="text.secondary" mb={4}>
                Cấu hình toàn bộ tham số vận hành của nền tảng
            </Typography>

            <Grid container spacing={3}>
                <Grid item xs={12} md={8}>
                    <Card>
                        <List>
                            {settingLinks.map((link, index) => (
                                <Box key={link.title}>
                                    <ListItem disablePadding>
                                        <ListItemButton onClick={() => navigate(link.path)} sx={{ py: 2 }}>
                                            <ListItemIcon sx={{ minWidth: 56, color: 'primary.main' }}>
                                                {link.icon}
                                            </ListItemIcon>
                                            <ListItemText
                                                primary={<Typography variant="subtitle1" fontWeight="600">{link.title}</Typography>}
                                                secondary={link.description}
                                            />
                                            <ChevronRight color="action" />
                                        </ListItemButton>
                                    </ListItem>
                                    {index < settingLinks.length - 1 && <Divider />}
                                </Box>
                            ))}
                        </List>
                    </Card>
                </Grid>

                <Grid item xs={12} md={4}>
                    <Card sx={{ bgcolor: 'primary.main', color: 'primary.contrastText' }}>
                        <CardContent>
                            <Typography variant="h6" fontWeight="700" gutterBottom>
                                System Information
                            </Typography>
                            <Box sx={{ mt: 2, display: 'flex', flexDirection: 'column', gap: 1 }}>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Version:</Typography>
                                    <Typography variant="body2" fontWeight="600">v2.5.0-stable</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Build:</Typography>
                                    <Typography variant="body2" fontWeight="600">20241207.1</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Environment:</Typography>
                                    <Typography variant="body2" fontWeight="600">Production</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Timezone:</Typography>
                                    <Typography variant="body2" fontWeight="600">Asia/Ho_Chi_Minh</Typography>
                                </Box>
                            </Box>
                            <Button
                                variant="contained"
                                color="secondary"
                                size="small"
                                fullWidth
                                sx={{ mt: 3 }}
                            >
                                Check for Updates
                            </Button>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
