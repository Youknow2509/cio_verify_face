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

export const AttendancePage: React.FC = () => {
    const [loading, setLoading] = useState(true);
    const [records, setRecords] = useState<any[]>([]);
    const [selectedMonth, setSelectedMonth] = useState(format(new Date(), 'yyyy-MM'));

    useEffect(() => {
        loadAttendanceRecords();
    }, [selectedMonth]);

    const loadAttendanceRecords = async () => {
        try {
            setLoading(true);
            const response: any = await attendanceApi.getMyAttendanceRecords({
                year_month: selectedMonth,
            });
            setRecords(response.data || []);
        } catch (error) {
            console.error('Failed to load attendance records:', error);
        } finally {
            setLoading(false);
        }
    };

    const getRecordTypeLabel = (type: number) => {
        switch (type) {
            case 0:
                return { label: 'Vào', color: 'success' as const };
            case 1:
                return { label: 'Ra', color: 'error' as const };
            default:
                return { label: 'Khác', color: 'default' as const };
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
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h4" fontWeight="700">
                    Lịch sử chấm công
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
                        Bản ghi chấm công - Tháng {selectedMonth}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={3}>
                        Tổng số: {records.length} bản ghi
                    </Typography>

                    <TableContainer component={Paper} sx={{ background: 'transparent' }}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 600 }}>Thời gian</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Loại</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Phương thức</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Điểm số</TableCell>
                                    <TableCell sx={{ fontWeight: 600 }}>Trạng thái</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {records.length > 0 ? (
                                    records.map((record, index) => {
                                        const typeInfo = getRecordTypeLabel(record.RecordType);
                                        return (
                                            <TableRow key={index}>
                                                <TableCell>
                                                    {format(new Date(record.RecordTime), 'dd/MM/yyyy HH:mm:ss')}
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={typeInfo.label}
                                                        color={typeInfo.color}
                                                        size="small"
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    {record.VerificationMethod === 'face' ? 'Khuôn mặt' : record.VerificationMethod}
                                                </TableCell>
                                                <TableCell>
                                                    {record.VerificationScore?.Float
                                                        ? (record.VerificationScore.Float * 100).toFixed(1) + '%'
                                                        : 'N/A'}
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={record.SyncStatus === 'synced' ? 'Đã đồng bộ' : 'Chưa đồng bộ'}
                                                        color={record.SyncStatus === 'synced' ? 'success' : 'warning'}
                                                        size="small"
                                                        variant="outlined"
                                                    />
                                                </TableCell>
                                            </TableRow>
                                        );
                                    })
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={5} align="center">
                                            <Typography variant="body2" color="text.secondary" py={4}>
                                                Không có bản ghi chấm công nào trong tháng này
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
