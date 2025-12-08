import { useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    TextField,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    CircularProgress,
    Alert,
    Grid,
} from '@mui/material';
import { Download as DownloadIcon } from '@mui/icons-material';
import { format } from 'date-fns';
import { useAuthStore } from '@/stores/authStore';
import { attendanceApi } from '@/services/api';

export const ExportPage: React.FC = () => {
    const { user } = useAuthStore();
    const [loading, setLoading] = useState(false);
    const [selectedMonth, setSelectedMonth] = useState(format(new Date(), 'yyyy-MM'));
    const [exportFormat, setExportFormat] = useState('excel');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleExport = async () => {
        if (!user?.email) {
            setError('Không tìm thấy thông tin email');
            return;
        }

        try {
            setLoading(true);
            setError('');
            setSuccess('');

            const response: any = await attendanceApi.exportMonthlySummary({
                email: user.email,
                format: exportFormat,
                month: selectedMonth,
            });

            setSuccess(
                `Yêu cầu xuất báo cáo đã được tạo thành công! Mã công việc: ${response.data?.job_id}. Báo cáo sẽ được gửi qua email khi hoàn tất.`
            );
        } catch (err: any) {
            setError(err.message || 'Không thể xuất báo cáo. Vui lòng thử lại.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box>
            <Typography variant="h4" fontWeight="700" mb={3}>
                Xuất báo cáo chấm công
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError('')}>
                    {error}
                </Alert>
            )}

            {success && (
                <Alert severity="success" sx={{ mb: 3 }} onClose={() => setSuccess('')}>
                    {success}
                </Alert>
            )}

            <Grid container spacing={3}>
                <Grid item xs={12} md={8}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={3}>
                                Cấu hình xuất báo cáo
                            </Typography>

                            <Box sx={{ mb: 3 }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    Email nhận báo cáo
                                </Typography>
                                <TextField
                                    fullWidth
                                    value={user?.email || ''}
                                    disabled
                                    size="small"
                                />
                            </Box>

                            <Box sx={{ mb: 3 }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    Chọn tháng
                                </Typography>
                                <TextField
                                    type="month"
                                    value={selectedMonth}
                                    onChange={(e) => setSelectedMonth(e.target.value)}
                                    size="small"
                                    fullWidth
                                />
                            </Box>

                            <Box sx={{ mb: 3 }}>
                                <FormControl fullWidth size="small">
                                    <InputLabel>Định dạng file</InputLabel>
                                    <Select
                                        value={exportFormat}
                                        label="Định dạng file"
                                        onChange={(e) => setExportFormat(e.target.value)}
                                    >
                                        <MenuItem value="excel">Excel (.xlsx)</MenuItem>
                                        <MenuItem value="pdf">PDF (.pdf)</MenuItem>
                                        <MenuItem value="csv">CSV (.csv)</MenuItem>
                                    </Select>
                                </FormControl>
                            </Box>

                            <Button
                                variant="contained"
                                startIcon={loading ? <CircularProgress size={20} /> : <DownloadIcon />}
                                onClick={handleExport}
                                disabled={loading}
                                fullWidth
                                size="large"
                            >
                                {loading ? 'Đang xử lý...' : 'Xuất báo cáo'}
                            </Button>
                        </CardContent>
                    </Card>
                </Grid>

                <Grid item xs={12} md={4}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={2}>
                                Thông tin
                            </Typography>
                            <Typography variant="body2" color="text.secondary" paragraph>
                                📊 Báo cáo sẽ bao gồm:
                            </Typography>
                            <Box component="ul" sx={{ pl: 2, color: 'text.secondary' }}>
                                <Typography component="li" variant="body2" mb={1}>
                                    Tổng hợp chấm công theo ngày
                                </Typography>
                                <Typography component="li" variant="body2" mb={1}>
                                    Thời gian vào/ra
                                </Typography>
                                <Typography component="li" variant="body2" mb={1}>
                                    Tổng giờ làm việc
                                </Typography>
                                <Typography component="li" variant="body2" mb={1}>
                                    Số ngày đi muộn/về sớm
                                </Typography>
                                <Typography component="li" variant="body2">
                                    Thống kê tổng hợp
                                </Typography>
                            </Box>
                            <Typography variant="body2" color="text.secondary" mt={2}>
                                ⏱️ Thời gian xử lý: 1-5 phút
                            </Typography>
                            <Typography variant="body2" color="text.secondary" mt={1}>
                                📧 Báo cáo sẽ được gửi qua email
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
