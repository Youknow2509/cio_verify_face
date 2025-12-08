import { useEffect, useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    CircularProgress,
    TextField,
    Chip,
    Paper,
} from '@mui/material';
import { format } from 'date-fns';
import { attendanceApi } from '@/services/api';

export const DailySummaryPage: React.FC = () => {
    const [loading, setLoading] = useState(true);
    const [summaries, setSummaries] = useState<any[]>([]);
    const [selectedMonth, setSelectedMonth] = useState(format(new Date(), 'yyyy-MM'));

    useEffect(() => {
        loadDailySummaries();
    }, [selectedMonth]);

    const loadDailySummaries = async () => {
        try {
            setLoading(true);
            const response: any = await attendanceApi.getMySummaries({
                month: selectedMonth,
            });
            setSummaries(response.data || []);
        } catch (error) {
            console.error('Failed to load daily summaries:', error);
        } finally {
            setLoading(false);
        }
    };

    const getAttendanceStatusLabel = (status: number) => {
        switch (status) {
            case 1:
                return { label: 'Có mặt', color: 'success' as const };
            case 0:
                return { label: 'Vắng mặt', color: 'error' as const };
            default:
                return { label: 'Không rõ', color: 'default' as const };
        }
    };

    const formatMinutes = (minutes: number) => {
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        return `${hours}h ${mins}p`;
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
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h4" fontWeight="700">
                    Tổng hợp theo ngày
                </Typography>
                <TextField
                    type="month"
                    value={selectedMonth}
                    onChange={(e) => setSelectedMonth(e.target.value)}
                    size="small"
                    sx={{ width: 200 }}
                />
            </Box>

            <Card>
                <CardContent>
                    <Typography variant="h6" fontWeight="600" mb={2}>
                        Tổng hợp chấm công theo ngày - Tháng {selectedMonth}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={3}>
                        Tổng số ngày: {summaries.length} ngày
                    </Typography>

                    <TableContainer component={Paper} sx={{ background: 'transparent' }}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 600 }}>Ngày</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Giờ vào</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Giờ ra</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Trạng thái</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Đi muộn</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Về sớm</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Tổng giờ</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {summaries.length > 0 ? (
                                    summaries.map((summary, index) => {
                                        const statusInfo = getAttendanceStatusLabel(summary.AttendanceStatus);
                                        return (
                                            <TableRow key={index}>
                                                <TableCell>
                                                    {format(new Date(summary.WorkDate), 'dd/MM/yyyy')}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.ActualCheckIn
                                                        ? format(new Date(summary.ActualCheckIn), 'HH:mm:ss')
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.ActualCheckOut
                                                        ? format(new Date(summary.ActualCheckOut), 'HH:mm:ss')
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={statusInfo.label}
                                                        color={statusInfo.color}
                                                        size="small"
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    {summary.LateMinutes > 0 ? (
                                                        <Chip
                                                            label={`${summary.LateMinutes} phút`}
                                                            color="warning"
                                                            size="small"
                                                            variant="outlined"
                                                        />
                                                    ) : (
                                                        '-'
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.EarlyLeaveMinutes > 0 ? (
                                                        <Chip
                                                            label={`${summary.EarlyLeaveMinutes} phút`}
                                                            color="warning"
                                                            size="small"
                                                            variant="outlined"
                                                        />
                                                    ) : (
                                                        '-'
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.TotalWorkMinutes
                                                        ? formatMinutes(summary.TotalWorkMinutes)
                                                        : '-'}
                                                </TableCell>
                                            </TableRow>
                                        );
                                    })
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={7} align="center">
                                            <Typography variant="body2" color="text.secondary" py={4}>
                                                Không có dữ liệu tổng hợp trong tháng này
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
        </Box>
    );
};
