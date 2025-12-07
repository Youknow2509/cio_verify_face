import { useEffect, useState } from 'react';
import {
    Box,
    Grid,
    Card,
    CardContent,
    Typography,
    Avatar,
    Chip,
    LinearProgress,
} from '@mui/material';
import {
    AccessTime,
    CalendarMonth,
    TrendingUp,
    EmojiEvents,
} from '@mui/icons-material';
import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

interface Stats {
    totalDays: number;
    totalHours: number;
    lateCount: number;
    onTimeRate: number;
}

export const DashboardPage: React.FC = () => {
    const { user, accessToken } = useAuthStore();
    const [stats, setStats] = useState<Stats | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const response = await axios.get('/api/v1/employee/my-stats', {
                    headers: { Authorization: `Bearer ${accessToken}` },
                });
                setStats(response.data.data);
            } catch (err) {
                console.error('Failed to fetch stats:', err);
                // Use mock data for demo
                setStats({
                    totalDays: 18,
                    totalHours: 144,
                    lateCount: 2,
                    onTimeRate: 0.89,
                });
            } finally {
                setLoading(false);
            }
        };

        fetchStats();
    }, [accessToken]);

    const currentHour = new Date().getHours();
    const greeting =
        currentHour < 12 ? 'Chào buổi sáng' : currentHour < 18 ? 'Chào buổi chiều' : 'Chào buổi tối';

    const statsCards = [
        {
            title: 'Ngày làm việc',
            value: stats?.totalDays || 0,
            unit: 'ngày',
            icon: <CalendarMonth />,
            color: '#2563eb',
            bgGradient: 'linear-gradient(135deg, rgba(37, 99, 235, 0.1), rgba(37, 99, 235, 0.05))',
        },
        {
            title: 'Tổng giờ làm',
            value: stats?.totalHours || 0,
            unit: 'giờ',
            icon: <AccessTime />,
            color: '#10b981',
            bgGradient: 'linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(16, 185, 129, 0.05))',
        },
        {
            title: 'Số lần đi trễ',
            value: stats?.lateCount || 0,
            unit: 'lần',
            icon: <TrendingUp />,
            color: '#f59e0b',
            bgGradient: 'linear-gradient(135deg, rgba(245, 158, 11, 0.1), rgba(245, 158, 11, 0.05))',
        },
        {
            title: 'Tỷ lệ đúng giờ',
            value: Math.round((stats?.onTimeRate || 0) * 100),
            unit: '%',
            icon: <EmojiEvents />,
            color: '#7c3aed',
            bgGradient: 'linear-gradient(135deg, rgba(124, 58, 237, 0.1), rgba(124, 58, 237, 0.05))',
        },
    ];

    if (loading) {
        return <LinearProgress />;
    }

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            {/* Welcome Section */}
            <Card
                sx={{
                    mb: 3,
                    background: 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                    color: 'white',
                }}
            >
                <CardContent sx={{ p: 4 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                        <Avatar
                            src={user?.avatar}
                            sx={{
                                width: 80,
                                height: 80,
                                fontSize: '2rem',
                                bgcolor: 'rgba(255,255,255,0.2)',
                                border: '3px solid rgba(255,255,255,0.5)',
                            }}
                        >
                            {user?.full_name?.charAt(0) || 'U'}
                        </Avatar>
                        <Box>
                            <Typography variant="h4" fontWeight="700" mb={0.5}>
                                {greeting}, {user?.full_name?.split(' ').pop() || 'bạn'}! 👋
                            </Typography>
                            <Typography variant="body1" sx={{ opacity: 0.9 }}>
                                Chúc bạn một ngày làm việc hiệu quả
                            </Typography>
                            <Box sx={{ display: 'flex', gap: 1, mt: 1.5 }}>
                                <Chip
                                    label={user?.department || 'Phòng ban'}
                                    size="small"
                                    sx={{ bgcolor: 'rgba(255,255,255,0.2)', color: 'white' }}
                                />
                                <Chip
                                    label={user?.position || 'Chức vụ'}
                                    size="small"
                                    sx={{ bgcolor: 'rgba(255,255,255,0.2)', color: 'white' }}
                                />
                            </Box>
                        </Box>
                    </Box>
                </CardContent>
            </Card>

            {/* Stats Cards */}
            <Typography variant="h6" fontWeight="600" mb={2}>
                Thống kê tháng này
            </Typography>
            <Grid container spacing={3}>
                {statsCards.map((card, index) => (
                    <Grid item xs={12} sm={6} md={3} key={index}>
                        <Card
                            sx={{
                                height: '100%',
                                background: card.bgGradient,
                                border: '1px solid',
                                borderColor: 'divider',
                                transition: 'transform 0.2s, box-shadow 0.2s',
                                '&:hover': {
                                    transform: 'translateY(-4px)',
                                    boxShadow: `0 12px 20px -10px ${card.color}40`,
                                },
                            }}
                        >
                            <CardContent>
                                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                                    <Box>
                                        <Typography variant="body2" color="text.secondary" mb={1}>
                                            {card.title}
                                        </Typography>
                                        <Typography variant="h4" fontWeight="700" sx={{ color: card.color }}>
                                            {card.value}
                                            <Typography component="span" variant="body1" sx={{ ml: 0.5 }}>
                                                {card.unit}
                                            </Typography>
                                        </Typography>
                                    </Box>
                                    <Avatar
                                        sx={{
                                            bgcolor: card.color,
                                            width: 48,
                                            height: 48,
                                        }}
                                    >
                                        {card.icon}
                                    </Avatar>
                                </Box>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>

            {/* Quick Actions */}
            <Typography variant="h6" fontWeight="600" mb={2} mt={4}>
                Hành động nhanh
            </Typography>
            <Grid container spacing={2}>
                <Grid item xs={12} md={4}>
                    <Card
                        sx={{
                            cursor: 'pointer',
                            '&:hover': { transform: 'translateY(-2px)' },
                            transition: 'transform 0.2s',
                        }}
                        onClick={() => (window.location.href = '/reports/monthly')}
                    >
                        <CardContent sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <Avatar sx={{ bgcolor: 'primary.main' }}>
                                <CalendarMonth />
                            </Avatar>
                            <Box>
                                <Typography fontWeight="600">Xem báo cáo tháng</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Chi tiết chấm công theo tháng
                                </Typography>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={12} md={4}>
                    <Card
                        sx={{
                            cursor: 'pointer',
                            '&:hover': { transform: 'translateY(-2px)' },
                            transition: 'transform 0.2s',
                        }}
                        onClick={() => (window.location.href = '/face-update')}
                    >
                        <CardContent sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <Avatar sx={{ bgcolor: 'secondary.main' }}>
                                <AccessTime />
                            </Avatar>
                            <Box>
                                <Typography fontWeight="600">Cập nhật khuôn mặt</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Gửi yêu cầu cập nhật ảnh
                                </Typography>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={12} md={4}>
                    <Card
                        sx={{
                            cursor: 'pointer',
                            '&:hover': { transform: 'translateY(-2px)' },
                            transition: 'transform 0.2s',
                        }}
                        onClick={() => (window.location.href = '/profile')}
                    >
                        <CardContent sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <Avatar sx={{ bgcolor: 'warning.main' }}>
                                <EmojiEvents />
                            </Avatar>
                            <Box>
                                <Typography fontWeight="600">Thông tin cá nhân</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Quản lý hồ sơ của bạn
                                </Typography>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
