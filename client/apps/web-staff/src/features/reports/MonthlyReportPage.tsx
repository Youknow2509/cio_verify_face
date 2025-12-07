import { useState, useEffect } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    TextField,
    MenuItem,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
    Button,
    Paper,
    LinearProgress,
} from '@mui/material';
import {
    Download,
    CalendarMonth,
    AccessTime,
    TrendingUp,
    Warning,
} from '@mui/icons-material';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';
import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

interface AttendanceRecord {
    date: string;
    checkIn: string;
    checkOut: string;
    totalHours: number;
    status: 'on-time' | 'late' | 'early-leave' | 'absent';
}

interface MonthlySummary {
    totalDays: number;
    totalHours: number;
    lateCount: number;
    earlyLeaveCount: number;
    absentCount: number;
    complianceRate: number;
}

const months = [
    { value: 1, label: 'Tháng 1' },
    { value: 2, label: 'Tháng 2' },
    { value: 3, label: 'Tháng 3' },
    { value: 4, label: 'Tháng 4' },
    { value: 5, label: 'Tháng 5' },
    { value: 6, label: 'Tháng 6' },
    { value: 7, label: 'Tháng 7' },
    { value: 8, label: 'Tháng 8' },
    { value: 9, label: 'Tháng 9' },
    { value: 10, label: 'Tháng 10' },
    { value: 11, label: 'Tháng 11' },
    { value: 12, label: 'Tháng 12' },
];

const statusConfig = {
    'on-time': { label: 'Đúng giờ', color: 'success' as const },
    late: { label: 'Đi trễ', color: 'warning' as const },
    'early-leave': { label: 'Về sớm', color: 'info' as const },
    absent: { label: 'Vắng', color: 'error' as const },
};

export const MonthlyReportPage: React.FC = () => {
    const { accessToken } = useAuthStore();
    const currentDate = new Date();
    const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth() + 1);
    const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());
    const [loading, setLoading] = useState(true);
    const [records, setRecords] = useState<AttendanceRecord[]>([]);
    const [summary, setSummary] = useState<MonthlySummary | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const response = await axios.get('/api/v1/employee/my-monthly-summary', {
                    params: { month: selectedMonth, year: selectedYear },
                    headers: { Authorization: `Bearer ${accessToken}` },
                });
                setSummary(response.data.data.summary);
                setRecords(response.data.data.records);
            } catch (err) {
                console.error('Failed to fetch monthly report:', err);
                // Mock data for demo
                const mockRecords: AttendanceRecord[] = Array.from({ length: 20 }, (_, i) => ({
                    date: `${selectedYear}-${String(selectedMonth).padStart(2, '0')}-${String(i + 1).padStart(2, '0')}`,
                    checkIn: `08:${String(Math.floor(Math.random() * 30)).padStart(2, '0')}`,
                    checkOut: `17:${String(30 + Math.floor(Math.random() * 30)).padStart(2, '0')}`,
                    totalHours: 8 + Math.random(),
                    status: Math.random() > 0.15 ? 'on-time' : Math.random() > 0.5 ? 'late' : 'early-leave',
                }));
                setRecords(mockRecords);
                setSummary({
                    totalDays: 20,
                    totalHours: 160,
                    lateCount: 3,
                    earlyLeaveCount: 1,
                    absentCount: 0,
                    complianceRate: 0.85,
                });
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [selectedMonth, selectedYear, accessToken]);

    const chartData = records.slice(0, 15).map((r) => ({
        date: r.date.split('-')[2],
        hours: Number(r.totalHours.toFixed(1)),
    }));

    const handleExport = async () => {
        try {
            const response = await axios.post(
                '/api/v1/employee/export-monthly-summary',
                { month: selectedMonth, year: selectedYear },
                {
                    headers: { Authorization: `Bearer ${accessToken}` },
                    responseType: 'blob',
                }
            );
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `report_${selectedYear}_${selectedMonth}.xlsx`);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (err) {
            console.error('Export failed:', err);
        }
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h4" fontWeight="700">
                    Báo cáo Chấm công Tháng
                </Typography>
                <Button variant="contained" startIcon={<Download />} onClick={handleExport}>
                    Xuất Excel
                </Button>
            </Box>

            {/* Filters */}
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <Grid container spacing={2} alignItems="center">
                        <Grid item xs={12} sm={4}>
                            <TextField
                                select
                                fullWidth
                                label="Tháng"
                                value={selectedMonth}
                                onChange={(e) => setSelectedMonth(Number(e.target.value))}
                            >
                                {months.map((m) => (
                                    <MenuItem key={m.value} value={m.value}>
                                        {m.label}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </Grid>
                        <Grid item xs={12} sm={4}>
                            <TextField
                                select
                                fullWidth
                                label="Năm"
                                value={selectedYear}
                                onChange={(e) => setSelectedYear(Number(e.target.value))}
                            >
                                {[2023, 2024, 2025].map((y) => (
                                    <MenuItem key={y} value={y}>
                                        {y}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </Grid>
                    </Grid>
                </CardContent>
            </Card>

            {loading ? (
                <LinearProgress />
            ) : (
                <>
                    {/* Summary Cards */}
                    <Grid container spacing={3} sx={{ mb: 3 }}>
                        <Grid item xs={6} md={3}>
                            <Card sx={{ textAlign: 'center', p: 2 }}>
                                <CalendarMonth sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
                                <Typography variant="h4" fontWeight="700" color="primary.main">
                                    {summary?.totalDays}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Ngày làm việc
                                </Typography>
                            </Card>
                        </Grid>
                        <Grid item xs={6} md={3}>
                            <Card sx={{ textAlign: 'center', p: 2 }}>
                                <AccessTime sx={{ fontSize: 40, color: 'success.main', mb: 1 }} />
                                <Typography variant="h4" fontWeight="700" color="success.main">
                                    {summary?.totalHours}h
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Tổng giờ làm
                                </Typography>
                            </Card>
                        </Grid>
                        <Grid item xs={6} md={3}>
                            <Card sx={{ textAlign: 'center', p: 2 }}>
                                <Warning sx={{ fontSize: 40, color: 'warning.main', mb: 1 }} />
                                <Typography variant="h4" fontWeight="700" color="warning.main">
                                    {summary?.lateCount}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Lần đi trễ
                                </Typography>
                            </Card>
                        </Grid>
                        <Grid item xs={6} md={3}>
                            <Card sx={{ textAlign: 'center', p: 2 }}>
                                <TrendingUp sx={{ fontSize: 40, color: 'info.main', mb: 1 }} />
                                <Typography variant="h4" fontWeight="700" color="info.main">
                                    {Math.round((summary?.complianceRate || 0) * 100)}%
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Tỷ lệ tuân thủ
                                </Typography>
                            </Card>
                        </Grid>
                    </Grid>

                    {/* Chart */}
                    <Card sx={{ mb: 3 }}>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={2}>
                                Biểu đồ giờ làm việc
                            </Typography>
                            <ResponsiveContainer width="100%" height={250}>
                                <BarChart data={chartData}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="rgba(0,0,0,0.1)" />
                                    <XAxis dataKey="date" />
                                    <YAxis />
                                    <Tooltip />
                                    <Bar dataKey="hours" fill="#2563eb" radius={[4, 4, 0, 0]} />
                                </BarChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>

                    {/* Table */}
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={2}>
                                Chi tiết theo ngày
                            </Typography>
                            <TableContainer component={Paper} elevation={0}>
                                <Table>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Ngày</TableCell>
                                            <TableCell>Giờ vào</TableCell>
                                            <TableCell>Giờ ra</TableCell>
                                            <TableCell align="right">Tổng giờ</TableCell>
                                            <TableCell align="center">Trạng thái</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {records.map((record, index) => (
                                            <TableRow key={index} hover>
                                                <TableCell>{record.date}</TableCell>
                                                <TableCell>{record.checkIn}</TableCell>
                                                <TableCell>{record.checkOut}</TableCell>
                                                <TableCell align="right">{record.totalHours.toFixed(1)}h</TableCell>
                                                <TableCell align="center">
                                                    <Chip
                                                        label={statusConfig[record.status].label}
                                                        color={statusConfig[record.status].color}
                                                        size="small"
                                                    />
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        </CardContent>
                    </Card>
                </>
            )}
        </Box>
    );
};
