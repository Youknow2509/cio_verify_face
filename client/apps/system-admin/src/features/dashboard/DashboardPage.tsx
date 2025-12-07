import { useState, useEffect } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    Avatar,
    Chip,
    Alert,
    LinearProgress,
} from '@mui/material';
import { DashboardSkeleton } from '@/components/Skeletons';
import {
    Business,
    People,
    Devices,
    TouchApp,
    TrendingUp,
    Warning,
    Schedule,
} from '@mui/icons-material';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    AreaChart,
    Area,
} from 'recharts';

// Mock data
const platformStats = [
    { label: 'Công ty hoạt động', value: 47, icon: <Business />, color: '#6366f1', change: '+3' },
    { label: 'Tổng nhân viên', value: '12,458', icon: <People />, color: '#10b981', change: '+124' },
    { label: 'Thiết bị online', value: 156, icon: <Devices />, color: '#f59e0b', change: '98%' },
    { label: 'Chấm công hôm nay', value: '8,934', icon: <TouchApp />, color: '#ec4899', change: '+15%' },
];

const companyGrowth = [
    { month: 'T7', companies: 38, users: 9200 },
    { month: 'T8', companies: 40, users: 9800 },
    { month: 'T9', companies: 42, users: 10500 },
    { month: 'T10', companies: 44, users: 11200 },
    { month: 'T11', companies: 45, users: 11800 },
    { month: 'T12', companies: 47, users: 12458 },
];

const systemLoad = [
    { time: '00:00', cpu: 25, ram: 45 },
    { time: '04:00', cpu: 15, ram: 40 },
    { time: '08:00', cpu: 65, ram: 70 },
    { time: '12:00', cpu: 55, ram: 65 },
    { time: '16:00', cpu: 70, ram: 75 },
    { time: '20:00', cpu: 45, ram: 55 },
    { time: 'Now', cpu: 42, ram: 58 },
];

const alerts = [
    { id: 1, type: 'warning', message: 'Tech Corp subscription expires in 3 days', company: 'Tech Corp' },
    { id: 2, type: 'warning', message: 'Device offline for 2 hours', company: 'ABC Company' },
    { id: 3, type: 'error', message: 'Payment failed - subscription suspended', company: 'XYZ Ltd' },
];

const topCompanies = [
    { name: 'Vinamilk', employees: 1250, checkIns: 1180 },
    { name: 'FPT Software', employees: 980, checkIns: 945 },
    { name: 'Tech Corp', employees: 756, checkIns: 720 },
    { name: 'Viettel', employees: 650, checkIns: 628 },
    { name: 'VNPT', employees: 520, checkIns: 498 },
];

export const DashboardPage: React.FC = () => {
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Simulate API fetch delay
        const timer = setTimeout(() => setLoading(false), 800);
        return () => clearTimeout(timer);
    }, []);

    if (loading) {
        return <DashboardSkeleton />;
    }

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            {/* Header */}
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" fontWeight="700" mb={1}>
                    Platform Dashboard
                </Typography>
                <Typography variant="body1" color="text.secondary">
                    Tổng quan hệ thống Face Attendance SaaS
                </Typography>
            </Box>

            {/* Stats Cards */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
                {platformStats.map((stat, index) => (
                    <Grid item xs={12} sm={6} lg={3} key={index}>
                        <Card sx={{ height: '100%' }}>
                            <CardContent>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                    <Box>
                                        <Typography variant="body2" color="text.secondary" mb={1}>
                                            {stat.label}
                                        </Typography>
                                        <Typography variant="h4" fontWeight="700">
                                            {stat.value}
                                        </Typography>
                                        <Chip
                                            label={stat.change}
                                            size="small"
                                            color="success"
                                            sx={{ mt: 1, height: 22, fontSize: '0.7rem' }}
                                        />
                                    </Box>
                                    <Avatar
                                        sx={{
                                            bgcolor: `${stat.color}20`,
                                            color: stat.color,
                                            width: 48,
                                            height: 48,
                                        }}
                                    >
                                        {stat.icon}
                                    </Avatar>
                                </Box>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>

            {/* Charts Row */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
                {/* Company Growth Chart */}
                <Grid item xs={12} lg={8}>
                    <Card sx={{ height: '100%' }}>
                        <CardContent>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                                <Box>
                                    <Typography variant="h6" fontWeight="600">
                                        Tăng trưởng Platform
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Số công ty & nhân viên 6 tháng gần nhất
                                    </Typography>
                                </Box>
                                <TrendingUp color="success" />
                            </Box>
                            <Box sx={{ height: 280 }}>
                                <ResponsiveContainer width="100%" height="100%">
                                    <BarChart data={companyGrowth}>
                                        <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                                        <XAxis dataKey="month" />
                                        <YAxis yAxisId="left" orientation="left" />
                                        <YAxis yAxisId="right" orientation="right" />
                                        <Tooltip />
                                        <Bar yAxisId="left" dataKey="companies" fill="#6366f1" radius={[4, 4, 0, 0]} name="Công ty" />
                                        <Bar yAxisId="right" dataKey="users" fill="#ec4899" radius={[4, 4, 0, 0]} name="Nhân viên" />
                                    </BarChart>
                                </ResponsiveContainer>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* System Load */}
                <Grid item xs={12} lg={4}>
                    <Card sx={{ height: '100%' }}>
                        <CardContent>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                                <Typography variant="h6" fontWeight="600">
                                    System Load
                                </Typography>
                                <Chip label="Live" size="small" color="success" />
                            </Box>
                            <Box sx={{ height: 200 }}>
                                <ResponsiveContainer width="100%" height="100%">
                                    <AreaChart data={systemLoad}>
                                        <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                                        <XAxis dataKey="time" tick={{ fontSize: 10 }} />
                                        <YAxis />
                                        <Tooltip />
                                        <Area type="monotone" dataKey="cpu" stackId="1" stroke="#6366f1" fill="#6366f1" fillOpacity={0.3} name="CPU %" />
                                        <Area type="monotone" dataKey="ram" stackId="2" stroke="#ec4899" fill="#ec4899" fillOpacity={0.3} name="RAM %" />
                                    </AreaChart>
                                </ResponsiveContainer>
                            </Box>
                            <Box sx={{ mt: 2, display: 'flex', gap: 3, justifyContent: 'center' }}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <Box sx={{ width: 12, height: 12, borderRadius: 1, bgcolor: '#6366f1' }} />
                                    <Typography variant="caption">CPU: 42%</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <Box sx={{ width: 12, height: 12, borderRadius: 1, bgcolor: '#ec4899' }} />
                                    <Typography variant="caption">RAM: 58%</Typography>
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>

            {/* Bottom Row */}
            <Grid container spacing={3}>
                {/* Alerts */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                                <Warning color="warning" />
                                <Typography variant="h6" fontWeight="600">
                                    Cảnh báo Hệ thống
                                </Typography>
                            </Box>
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                                {alerts.map((alert) => (
                                    <Alert
                                        key={alert.id}
                                        severity={alert.type as 'warning' | 'error'}
                                        sx={{ borderRadius: 2 }}
                                    >
                                        <Box>
                                            <Typography variant="body2" fontWeight="600">
                                                {alert.company}
                                            </Typography>
                                            <Typography variant="caption">
                                                {alert.message}
                                            </Typography>
                                        </Box>
                                    </Alert>
                                ))}
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Top Companies */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                                <Schedule color="primary" />
                                <Typography variant="h6" fontWeight="600">
                                    Top Công ty Chấm công Hôm nay
                                </Typography>
                            </Box>
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                                {topCompanies.map((company, index) => (
                                    <Box key={company.name} sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                        <Avatar
                                            sx={{
                                                width: 32,
                                                height: 32,
                                                bgcolor: index === 0 ? 'warning.main' : 'action.hover',
                                                fontSize: '0.8rem',
                                                fontWeight: 700,
                                            }}
                                        >
                                            {index + 1}
                                        </Avatar>
                                        <Box sx={{ flex: 1 }}>
                                            <Typography variant="body2" fontWeight="600">
                                                {company.name}
                                            </Typography>
                                            <LinearProgress
                                                variant="determinate"
                                                value={(company.checkIns / company.employees) * 100}
                                                sx={{ height: 6, borderRadius: 1, mt: 0.5 }}
                                            />
                                        </Box>
                                        <Typography variant="body2" color="text.secondary">
                                            {company.checkIns}/{company.employees}
                                        </Typography>
                                    </Box>
                                ))}
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
