import { useEffect, useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    CircularProgress,
} from '@mui/material';
import {
    CheckCircle as CheckCircleIcon,
    AccessTime as AccessTimeIcon,
    Schedule as ScheduleIcon,
    EventAvailable as EventAvailableIcon,
} from '@mui/icons-material';
import { format } from 'date-fns';
import { vi } from 'date-fns/locale';
import { attendanceApi, shiftApi } from '@/services/api';

export const DashboardPage: React.FC = () => {
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState({
        todayAttendance: null as any,
        thisMonthStats: { total: 0, present: 0, absent: 0 },
        shifts: [],
    });

    useEffect(() => {
        loadDashboardData();
    }, []);

    const loadDashboardData = async () => {
        try {
            setLoading(true);
            const currentMonth = format(new Date(), 'yyyy-MM');
            
            // Load attendance records for current month
            const attendanceRes: any = await attendanceApi.getMyAttendanceRecords({
                year_month: currentMonth,
            });

            // Load daily summaries
            const summaryRes: any = await attendanceApi.getMySummaries({
                month: currentMonth,
            });

            // Load shifts
            const shiftsRes: any = await shiftApi.getEmployeeShifts({
                page: 1,
                size: 10,
            });

            // Calculate stats
            const todayRecords = attendanceRes.data.filter((r: any) => {
                const recordDate = format(new Date(r.RecordTime), 'yyyy-MM-dd');
                const today = format(new Date(), 'yyyy-MM-dd');
                return recordDate === today;
            });

            setStats({
                todayAttendance: todayRecords.length > 0 ? todayRecords : null,
                thisMonthStats: {
                    total: summaryRes.data.length,
                    present: summaryRes.data.filter((s: any) => s.AttendanceStatus === 1).length,
                    absent: summaryRes.data.filter((s: any) => s.AttendanceStatus === 0).length,
                },
                shifts: shiftsRes.data?.shifts || [],
            });
        } catch (error) {
            console.error('Failed to load dashboard data:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    return (
        <Box>
            <Typography variant="h4" fontWeight="700" mb={3}>
                Bảng điều khiển
            </Typography>

            <Grid container spacing={3}>
                {/* Today's Attendance */}
                <Grid item xs={12} sm={6} md={3}>
                    <Card
                        sx={{
                            background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.2) 0%, rgba(5, 150, 105, 0.2) 100%)',
                            border: '1px solid rgba(16, 185, 129, 0.3)',
                        }}
                    >
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                <Box
                                    sx={{
                                        width: 60,
                                        height: 60,
                                        borderRadius: 2,
                                        background: 'rgba(16, 185, 129, 0.2)',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                    }}
                                >
                                    <CheckCircleIcon sx={{ fontSize: 32, color: '#10b981' }} />
                                </Box>
                                <Box>
                                    <Typography variant="body2" color="text.secondary">
                                        Hôm nay
                                    </Typography>
                                    <Typography variant="h5" fontWeight="700">
                                        {stats.todayAttendance ? 'Đã chấm công' : 'Chưa chấm công'}
                                    </Typography>
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Total Days This Month */}
                <Grid item xs={12} sm={6} md={3}>
                    <Card
                        sx={{
                            background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(37, 99, 235, 0.2) 100%)',
                            border: '1px solid rgba(59, 130, 246, 0.3)',
                        }}
                    >
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                <Box
                                    sx={{
                                        width: 60,
                                        height: 60,
                                        borderRadius: 2,
                                        background: 'rgba(59, 130, 246, 0.2)',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                    }}
                                >
                                    <EventAvailableIcon sx={{ fontSize: 32, color: '#3b82f6' }} />
                                </Box>
                                <Box>
                                    <Typography variant="body2" color="text.secondary">
                                        Tháng này
                                    </Typography>
                                    <Typography variant="h5" fontWeight="700">
                                        {stats.thisMonthStats.total} ngày
                                    </Typography>
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Present Days */}
                <Grid item xs={12} sm={6} md={3}>
                    <Card
                        sx={{
                            background: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(124, 58, 237, 0.2) 100%)',
                            border: '1px solid rgba(139, 92, 246, 0.3)',
                        }}
                    >
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                <Box
                                    sx={{
                                        width: 60,
                                        height: 60,
                                        borderRadius: 2,
                                        background: 'rgba(139, 92, 246, 0.2)',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                    }}
                                >
                                    <AccessTimeIcon sx={{ fontSize: 32, color: '#8b5cf6' }} />
                                </Box>
                                <Box>
                                    <Typography variant="body2" color="text.secondary">
                                        Đi làm
                                    </Typography>
                                    <Typography variant="h5" fontWeight="700">
                                        {stats.thisMonthStats.present} ngày
                                    </Typography>
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Active Shifts */}
                <Grid item xs={12} sm={6} md={3}>
                    <Card
                        sx={{
                            background: 'linear-gradient(135deg, rgba(245, 158, 11, 0.2) 0%, rgba(217, 119, 6, 0.2) 100%)',
                            border: '1px solid rgba(245, 158, 11, 0.3)',
                        }}
                    >
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                <Box
                                    sx={{
                                        width: 60,
                                        height: 60,
                                        borderRadius: 2,
                                        background: 'rgba(245, 158, 11, 0.2)',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                    }}
                                >
                                    <ScheduleIcon sx={{ fontSize: 32, color: '#f59e0b' }} />
                                </Box>
                                <Box>
                                    <Typography variant="body2" color="text.secondary">
                                        Ca làm việc
                                    </Typography>
                                    <Typography variant="h5" fontWeight="700">
                                        {stats.shifts.length} ca
                                    </Typography>
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Active Shifts List */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={2}>
                                Ca làm việc hiện tại
                            </Typography>
                            {stats.shifts.length > 0 ? (
                                <Grid container spacing={2}>
                                    {stats.shifts.slice(0, 4).map((shift: any, index: number) => (
                                        <Grid item xs={12} sm={6} md={3} key={index}>
                                            <Box
                                                sx={{
                                                    p: 2,
                                                    borderRadius: 2,
                                                    background: 'rgba(59, 130, 246, 0.1)',
                                                    border: '1px solid rgba(59, 130, 246, 0.2)',
                                                }}
                                            >
                                                <Typography variant="subtitle2" fontWeight="600" mb={1}>
                                                    {shift.shift_name}
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    {shift.shift_start} - {shift.shift_end}
                                                </Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    {shift.is_active ? '🟢 Đang hoạt động' : '🔴 Không hoạt động'}
                                                </Typography>
                                            </Box>
                                        </Grid>
                                    ))}
                                </Grid>
                            ) : (
                                <Typography variant="body2" color="text.secondary">
                                    Chưa có ca làm việc nào được gán
                                </Typography>
                            )}
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
