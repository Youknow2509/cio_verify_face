import { useParams, useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    Grid,
    Chip,
    Avatar,
    Divider,
    Tabs,
    Tab,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    LinearProgress,
} from '@mui/material';
import {
    ArrowBack,
    Business,
    People,
    Devices,
    Email,
    Phone,
    CalendarMonth,
    Storage,
    CreditCard,
} from '@mui/icons-material';
import { useState } from 'react';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';
import { EditLimitsDialog } from './EditLimitsDialog';
import { Edit } from '@mui/icons-material';

// Mock data
const companyDetail = {
    id: '1',
    name: 'Vinamilk',
    email: 'admin@vinamilk.com',
    phone: '028-1234-5678',
    address: '10 Tan Trao, Tan Phu, District 7, HCMC',
    plan: 'enterprise',
    status: 'active',
    employeeCount: 1250,
    employeeLimit: 2000,
    deviceCount: 12,
    deviceLimit: 50,
    storageUsed: 45.2,
    storageLimit: 100,
    registeredAt: '2023-01-15',
    expiresAt: '2024-12-15',
    admins: [
        { name: 'Nguyen Van A', email: 'a.nguyen@vinamilk.com', role: 'Super Admin' },
        { name: 'Tran Thi B', email: 'b.tran@vinamilk.com', role: 'HR Admin' },
    ],
};

const usageData = [
    { date: 'T1', checkIns: 1200, employees: 1150 },
    { date: 'T2', checkIns: 1180, employees: 1180 },
    { date: 'T3', checkIns: 1220, employees: 1200 },
    { date: 'T4', checkIns: 1250, employees: 1220 },
    { date: 'T5', checkIns: 1230, employees: 1240 },
    { date: 'T6', checkIns: 1280, employees: 1250 },
];

const auditLog = [
    { time: '2024-12-07 10:30', action: 'Thêm nhân viên mới', user: 'a.nguyen@vinamilk.com', details: 'Added emp #1251' },
    { time: '2024-12-07 09:15', action: 'Cập nhật ca làm việc', user: 'b.tran@vinamilk.com', details: 'Modified shift "Ca sáng"' },
    { time: '2024-12-06 16:45', action: 'Kích hoạt thiết bị', user: 'a.nguyen@vinamilk.com', details: 'Device #DVB-012 activated' },
];

export const CompanyDetailPage: React.FC = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [tab, setTab] = useState(0);
    const [limitsDialogOpen, setLimitsDialogOpen] = useState(false);

    const handleSaveLimits = (newLimits: any) => {
        console.log('Saving limits:', newLimits);
        // Todo: API call to update limits
        setLimitsDialogOpen(false);
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            {/* Header */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
                <Button startIcon={<ArrowBack />} onClick={() => navigate('/companies')}>
                    Quay lại
                </Button>
                <Box sx={{ flex: 1 }}>
                    <Typography variant="h4" fontWeight="700">
                        {companyDetail.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        ID: {id}
                    </Typography>
                </Box>
                <Chip label={companyDetail.plan.toUpperCase()} color="secondary" />
                <Chip label={companyDetail.status} color="success" />
            </Box>

            {/* Stats Cards */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
                <Grid item xs={6} md={3}>
                    <Card sx={{ position: 'relative' }}>
                        <CardContent sx={{ textAlign: 'center' }}>
                            <People color="primary" sx={{ fontSize: 32, mb: 1 }} />
                            <Typography variant="h5" fontWeight="700">
                                {companyDetail.employeeCount}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                / {companyDetail.employeeLimit} nhân viên
                            </Typography>
                            <Button
                                size="small"
                                sx={{ position: 'absolute', top: 8, right: 8, minWidth: 32, p: 0.5 }}
                                onClick={() => setLimitsDialogOpen(true)}
                            >
                                <Edit fontSize="small" />
                            </Button>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Card>
                        <CardContent sx={{ textAlign: 'center' }}>
                            <Devices color="warning" sx={{ fontSize: 32, mb: 1 }} />
                            <Typography variant="h5" fontWeight="700">
                                {companyDetail.deviceCount}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                thiết bị
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Card>
                        <CardContent sx={{ textAlign: 'center' }}>
                            <Storage color="info" sx={{ fontSize: 32, mb: 1 }} />
                            <Typography variant="h5" fontWeight="700">
                                {companyDetail.storageUsed} GB
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                / {companyDetail.storageLimit} GB
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Card>
                        <CardContent sx={{ textAlign: 'center' }}>
                            <CalendarMonth color="success" sx={{ fontSize: 32, mb: 1 }} />
                            <Typography variant="h6" fontWeight="600">
                                {new Date(companyDetail.expiresAt).toLocaleDateString('vi-VN')}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                ngày hết hạn
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>

            {/* Tabs */}
            <Card>
                <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ px: 2, pt: 1 }}>
                    <Tab label="Thông tin" />
                    <Tab label="Thống kê" />
                    <Tab label="Audit Log" />
                </Tabs>
                <Divider />
                <CardContent>
                    {tab === 0 && (
                        <Grid container spacing={4}>
                            <Grid item xs={12} md={6}>
                                <Typography variant="subtitle1" fontWeight="600" mb={2}>
                                    Thông tin liên hệ
                                </Typography>
                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 4 }}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                        <Email color="action" />
                                        <Typography>{companyDetail.email}</Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                        <Phone color="action" />
                                        <Typography>{companyDetail.phone}</Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                        <Business color="action" />
                                        <Typography>{companyDetail.address}</Typography>
                                    </Box>
                                </Box>

                                <Divider sx={{ my: 3 }} />

                                <Typography variant="subtitle1" fontWeight="600" mb={2}>
                                    Thông tin thanh toán
                                </Typography>
                                <Card variant="outlined" sx={{ mb: 2 }}>
                                    <CardContent sx={{ py: 2 }}>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
                                            <CreditCard color="primary" />
                                            <Typography fontWeight="600">Visa ending in 4242</Typography>
                                        </Box>
                                        <Typography variant="caption" color="text.secondary">
                                            Expires 12/2025 • Default Method
                                        </Typography>
                                    </CardContent>
                                </Card>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                    <Chip size="small" label="Auto-renew On" color="success" variant="outlined" />
                                    <Chip size="small" label="Next Invoice: 15/01/2025" variant="outlined" />
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Typography variant="subtitle1" fontWeight="600" mb={2}>
                                    Quản trị viên
                                </Typography>
                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 4 }}>
                                    {companyDetail.admins.map((admin) => (
                                        <Box key={admin.email} sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                            <Avatar sx={{ width: 36, height: 36, bgcolor: 'primary.main' }}>
                                                {admin.name.charAt(0)}
                                            </Avatar>
                                            <Box>
                                                <Typography variant="body2" fontWeight="600">
                                                    {admin.name}
                                                </Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    {admin.email} • {admin.role}
                                                </Typography>
                                            </Box>
                                        </Box>
                                    ))}
                                </Box>

                                <Divider sx={{ my: 3 }} />

                                <Typography variant="subtitle1" fontWeight="600" mb={2}>
                                    Lịch sử Gói dịch vụ
                                </Typography>
                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                                    {[
                                        { date: '2024-01-15', action: 'Upgraded to Enterprise', user: 'System' },
                                        { date: '2023-01-15', action: 'Subscribed to Professional', user: 'Sales Team' },
                                    ].map((history, idx) => (
                                        <Box key={idx} sx={{ display: 'flex', gap: 2 }}>
                                            <Typography variant="caption" color="text.secondary" sx={{ minWidth: 80 }}>
                                                {history.date}
                                            </Typography>
                                            <Box>
                                                <Typography variant="body2" fontWeight="500">
                                                    {history.action}
                                                </Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    by {history.user}
                                                </Typography>
                                            </Box>
                                        </Box>
                                    ))}
                                </Box>
                            </Grid>
                        </Grid>
                    )}

                    {tab === 1 && (
                        <Box>
                            <Typography variant="subtitle1" fontWeight="600" mb={2}>
                                Thống kê chấm công 6 tháng gần nhất
                            </Typography>
                            <Box sx={{ height: 300 }}>
                                <ResponsiveContainer width="100%" height="100%">
                                    <LineChart data={usageData}>
                                        <CartesianGrid strokeDasharray="3 3" />
                                        <XAxis dataKey="date" />
                                        <YAxis />
                                        <Tooltip />
                                        <Line type="monotone" dataKey="checkIns" stroke="#6366f1" name="Check-ins" strokeWidth={2} />
                                        <Line type="monotone" dataKey="employees" stroke="#10b981" name="Active Employees" strokeWidth={2} />
                                    </LineChart>
                                </ResponsiveContainer>
                            </Box>
                        </Box>
                    )}

                    {tab === 2 && (
                        <TableContainer>
                            <Table size="small">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Thời gian</TableCell>
                                        <TableCell>Hành động</TableCell>
                                        <TableCell>Người thực hiện</TableCell>
                                        <TableCell>Chi tiết</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {auditLog.map((log, index) => (
                                        <TableRow key={index}>
                                            <TableCell>{log.time}</TableCell>
                                            <TableCell>{log.action}</TableCell>
                                            <TableCell>{log.user}</TableCell>
                                            <TableCell>{log.details}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    )}
                </CardContent>
            </Card>

            <EditLimitsDialog
                open={limitsDialogOpen}
                onClose={() => setLimitsDialogOpen(false)}
                currentLimits={{
                    employees: companyDetail.employeeLimit,
                    devices: companyDetail.deviceLimit,
                    storage: companyDetail.storageLimit,
                }}
                onSave={handleSaveLimits}
            />
        </Box>
    );
};
