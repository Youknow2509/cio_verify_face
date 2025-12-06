import { useState, useEffect, useCallback } from 'react';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
    Chip,
    Grid,
    CircularProgress,
    Alert,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Snackbar,
    ToggleButton,
    ToggleButtonGroup,
} from '@mui/material';
import {
    Download,
    People,
    CheckCircle,
    Cancel,
    Schedule,
    ExitToApp,
    AccessTime,
    CalendarMonth,
    CalendarToday,
} from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';

// Types for Daily Attendance Status API
interface DailyEmployee {
    employee_id: string;
    name: string;
    check_in: string | null; // HH:mm:ss format
    check_out: string | null; // HH:mm:ss format
    status: 'on_time' | 'late' | 'early_leave' | 'absent' | 'overtime';
    late_minutes: number;
    total_hours: number;
}

interface DailyStatistics {
    total_employees: number;
    checked_in: number;
    not_checked_in: number;
    on_time: number;
    late: number;
    early_leave: number;
    overtime: number;
}

interface DailyAttendanceStatusResponse {
    success: boolean;
    data: {
        date: string; // YYYY-MM-DD
        statistics: DailyStatistics;
        employees: DailyEmployee[];
    };
}

// Types for Monthly Summary API
interface MonthlyEmployeeSummary {
    employee_id: string;
    name: string;
    present_days: number;
    absent_days: number;
    late_days: number;
    total_hours: number;
    average_hours_per_day: number;
    overtime_hours: number;
}

interface MonthlyStatistics {
    average_attendance_rate: number;
    total_late_instances: number;
    total_overtime_hours: number;
}

interface MonthlySummaryResponse {
    success: boolean;
    data: {
        month: string; // YYYY-MM
        total_working_days: number;
        employees_summary: MonthlyEmployeeSummary[];
        statistics: MonthlyStatistics;
    };
}

// Export response type
interface ExportResponse {
    success: boolean;
    data: {
        file_id: string;
        download_url: string;
        expires_at: string;
    };
}

type ReportMode = 'daily' | 'monthly';

export const DailyReportPage: React.FC = () => {
    const [mode, setMode] = useState<ReportMode>('daily');
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
    const [month, setMonth] = useState(
        new Date().toISOString().slice(0, 7) // YYYY-MM format
    );
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Daily report states
    const [dailyData, setDailyData] = useState<
        DailyAttendanceStatusResponse['data'] | null
    >(null);

    // Monthly report states
    const [monthlyData, setMonthlyData] = useState<
        MonthlySummaryResponse['data'] | null
    >(null);

    // Export states
    const [exportDialogOpen, setExportDialogOpen] = useState(false);
    const [exportFormat, setExportFormat] = useState<'excel' | 'pdf' | 'csv'>(
        'excel'
    );
    const [exportEmail, setExportEmail] = useState('');
    const [exportLoading, setExportLoading] = useState(false);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({ open: false, message: '', severity: 'success' });

    // Get company ID from JWT token
    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
        }
    };

    // Fetch daily attendance status
    const fetchDailyReport = useCallback(async (selectedDate: string) => {
        const accessToken = localStorage.getItem('access_token');
        if (!accessToken) {
            setError('Không tìm thấy token xác thực');
            return;
        }

        const companyId = getCompanyIdFromToken(accessToken);
        if (!companyId) {
            setError('Không tìm thấy company_id trong token');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await apiClient.get<DailyAttendanceStatusResponse>(
                '/api/v1/company/daily-attendance-status',
                {
                    params: {
                        company_id: companyId,
                        date: selectedDate,
                    },
                }
            );

            if (response.data?.success && response.data?.data) {
                setDailyData(response.data.data);
            } else {
                setError('Không thể tải dữ liệu báo cáo');
            }
        } catch (err: any) {
            console.error('Failed to fetch daily report:', err);
            setError(
                err.response?.data?.message ||
                    err.response?.data?.error ||
                    'Không thể tải báo cáo. Vui lòng thử lại.'
            );
            setDailyData(null);
        } finally {
            setLoading(false);
        }
    }, []);

    // Fetch monthly summary
    const fetchMonthlyReport = useCallback(async (selectedMonth: string) => {
        const accessToken = localStorage.getItem('access_token');
        if (!accessToken) {
            setError('Không tìm thấy token xác thực');
            return;
        }

        const companyId = getCompanyIdFromToken(accessToken);
        if (!companyId) {
            setError('Không tìm thấy company_id trong token');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await apiClient.get<MonthlySummaryResponse>(
                '/api/v1/company/monthly-summary',
                {
                    params: {
                        company_id: companyId,
                        month: selectedMonth,
                    },
                }
            );

            if (response.data?.success && response.data?.data) {
                setMonthlyData(response.data.data);
            } else {
                setError('Không thể tải dữ liệu báo cáo');
            }
        } catch (err: any) {
            console.error('Failed to fetch monthly report:', err);
            setError(
                err.response?.data?.message ||
                    err.response?.data?.error ||
                    'Không thể tải báo cáo. Vui lòng thử lại.'
            );
            setMonthlyData(null);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        if (mode === 'daily') {
            fetchDailyReport(date);
        } else {
            fetchMonthlyReport(month);
        }
    }, [mode, date, month, fetchDailyReport, fetchMonthlyReport]);

    const handleModeChange = (
        _event: React.MouseEvent<HTMLElement>,
        newMode: ReportMode | null
    ) => {
        if (newMode !== null) {
            setMode(newMode);
        }
    };

    const handleDateChange = (newDate: string) => {
        setDate(newDate);
    };

    const handleMonthChange = (newMonth: string) => {
        setMonth(newMonth);
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'on_time':
                return 'success';
            case 'late':
                return 'warning';
            case 'early_leave':
                return 'info';
            case 'absent':
                return 'error';
            case 'overtime':
                return 'secondary';
            default:
                return 'default';
        }
    };

    const getStatusText = (status: string) => {
        switch (status) {
            case 'on_time':
                return 'Đúng giờ';
            case 'late':
                return 'Trễ';
            case 'early_leave':
                return 'Về sớm';
            case 'absent':
                return 'Vắng mặt';
            case 'overtime':
                return 'Tăng ca';
            default:
                return '-';
        }
    };

    const formatHours = (hours: number) => {
        if (!hours) return '-';
        const h = Math.floor(hours);
        const m = Math.round((hours - h) * 60);
        if (h > 0 && m > 0) {
            return `${h}h${m}m`;
        } else if (h > 0) {
            return `${h}h`;
        } else {
            return `${m}m`;
        }
    };

    // Handle export report
    const handleExportReport = async () => {
        const accessToken = localStorage.getItem('access_token');
        if (!accessToken) {
            setSnackbar({
                open: true,
                message: 'Không tìm thấy token xác thực',
                severity: 'error',
            });
            return;
        }

        const companyId = getCompanyIdFromToken(accessToken);
        if (!companyId) {
            setSnackbar({
                open: true,
                message: 'Không tìm thấy company_id trong token',
                severity: 'error',
            });
            return;
        }

        setExportLoading(true);
        try {
            const endpoint =
                mode === 'daily'
                    ? '/api/v1/company/export-daily-status'
                    : '/api/v1/company/export-monthly-summary';

            const requestBody: {
                company_id: string;
                date?: string;
                month?: string;
                format: string;
                email?: string;
            } = {
                company_id: companyId,
                format: exportFormat,
            };

            if (mode === 'daily') {
                requestBody.date = date;
            } else {
                requestBody.month = month;
            }

            if (exportEmail.trim()) {
                requestBody.email = exportEmail.trim();
            }

            const response = await apiClient.post<ExportResponse>(
                endpoint,
                requestBody
            );

            if (response.data?.success && response.data?.data) {
                const { download_url } = response.data.data;

                if (exportEmail.trim()) {
                    // Email sent
                    setSnackbar({
                        open: true,
                        message: `Báo cáo đã được gửi đến email ${exportEmail}`,
                        severity: 'success',
                    });
                } else {
                    // Download file
                    try {
                        const downloadResponse = await apiClient.get(
                            download_url,
                            {
                                responseType: 'blob',
                            }
                        );

                        const blob = new Blob([downloadResponse.data], {
                            type:
                                downloadResponse.headers['content-type'] ||
                                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
                        });

                        const url = window.URL.createObjectURL(blob);
                        const link = document.createElement('a');
                        link.href = url;
                        const extension =
                            exportFormat === 'excel' ? 'xlsx' : exportFormat;
                        const fileName = `${mode}-report-${
                            mode === 'daily' ? date : month
                        }.${extension}`;
                        link.setAttribute('download', fileName);
                        document.body.appendChild(link);
                        link.click();
                        link.remove();
                        window.URL.revokeObjectURL(url);

                        setSnackbar({
                            open: true,
                            message: 'Xuất báo cáo thành công',
                            severity: 'success',
                        });
                    } catch (downloadErr: any) {
                        console.error('Failed to download file:', downloadErr);
                        setSnackbar({
                            open: true,
                            message:
                                'Xuất báo cáo thành công nhưng không thể tải xuống. Vui lòng thử lại.',
                            severity: 'error',
                        });
                    }
                }

                setExportDialogOpen(false);
                setExportEmail('');
            } else {
                setSnackbar({
                    open: true,
                    message: 'Không thể xuất báo cáo. Vui lòng thử lại.',
                    severity: 'error',
                });
            }
        } catch (err: any) {
            console.error('Failed to export report:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message ||
                    err.response?.data?.error ||
                    'Không thể xuất báo cáo. Vui lòng thử lại.',
                severity: 'error',
            });
        } finally {
            setExportLoading(false);
        }
    };

    const handleOpenExportDialog = () => {
        setExportDialogOpen(true);
        setExportFormat('excel');
        setExportEmail('');
    };

    const handleCloseExportDialog = () => {
        if (!exportLoading) {
            setExportDialogOpen(false);
            setExportEmail('');
        }
    };

    return (
        <Box>
            <Typography variant="h4" fontWeight="bold" mb={3}>
                Báo cáo Chấm công
            </Typography>

            {/* Mode Toggle */}
            <Card sx={{ mb: 3, p: 2 }}>
                <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} md={4}>
                        <FormControl fullWidth>
                            <ToggleButtonGroup
                                value={mode}
                                exclusive
                                onChange={handleModeChange}
                                aria-label="report mode"
                                fullWidth
                            >
                                <ToggleButton value="daily" aria-label="daily">
                                    <CalendarToday sx={{ mr: 1 }} />
                                    Theo ngày
                                </ToggleButton>
                                <ToggleButton
                                    value="monthly"
                                    aria-label="monthly"
                                >
                                    <CalendarMonth sx={{ mr: 1 }} />
                                    Theo tháng
                                </ToggleButton>
                            </ToggleButtonGroup>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12} md={4}>
                        {mode === 'daily' ? (
                            <TextField
                                fullWidth
                                label="Ngày"
                                type="date"
                                value={date}
                                onChange={(e) =>
                                    handleDateChange(e.target.value)
                                }
                                InputLabelProps={{ shrink: true }}
                            />
                        ) : (
                            <TextField
                                fullWidth
                                label="Tháng"
                                type="month"
                                value={month}
                                onChange={(e) =>
                                    handleMonthChange(e.target.value)
                                }
                                InputLabelProps={{ shrink: true }}
                            />
                        )}
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <Button
                            variant="contained"
                            startIcon={<Download />}
                            onClick={handleOpenExportDialog}
                            disabled={loading}
                            fullWidth
                        >
                            Xuất báo cáo
                        </Button>
                    </Grid>
                </Grid>
            </Card>

            {error && (
                <Alert
                    severity="error"
                    sx={{ mb: 2 }}
                    onClose={() => setError(null)}
                >
                    {error}
                </Alert>
            )}

            {loading ? (
                <Box
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    minHeight="200px"
                >
                    <CircularProgress />
                </Box>
            ) : (
                <>
                    {/* Daily Report View */}
                    {mode === 'daily' && dailyData && dailyData.statistics && (
                        <>
                            {/* Statistics Cards */}
                            <Grid container spacing={2} sx={{ mb: 3 }}>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Tổng nhân viên
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                    >
                                                        {dailyData.statistics
                                                            ?.total_employees ??
                                                            0}
                                                    </Typography>
                                                </Box>
                                                <People
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'primary.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đã chấm công
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="success.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.checked_in ?? 0}
                                                    </Typography>
                                                </Box>
                                                <CheckCircle
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'success.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Chưa chấm công
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="error.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.not_checked_in ??
                                                            0}
                                                    </Typography>
                                                </Box>
                                                <Cancel
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'error.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đúng giờ
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="success.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.on_time ?? 0}
                                                    </Typography>
                                                </Box>
                                                <CheckCircle
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'success.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đi trễ
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="warning.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.late ?? 0}
                                                    </Typography>
                                                </Box>
                                                <Schedule
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'warning.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Về sớm
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="info.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.early_leave ?? 0}
                                                    </Typography>
                                                </Box>
                                                <ExitToApp
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'info.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Tăng ca
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="secondary.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.overtime ?? 0}
                                                    </Typography>
                                                </Box>
                                                <AccessTime
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'secondary.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                            </Grid>

                            {/* Employees Table */}
                            <Card>
                                <TableContainer>
                                    <Table>
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Nhân viên</TableCell>
                                                <TableCell>Giờ vào</TableCell>
                                                <TableCell>Giờ ra</TableCell>
                                                <TableCell>
                                                    Tổng giờ làm
                                                </TableCell>
                                                <TableCell>
                                                    Trễ (phút)
                                                </TableCell>
                                                <TableCell>
                                                    Trạng thái
                                                </TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {!dailyData.employees ||
                                            dailyData.employees.length === 0 ? (
                                                <TableRow>
                                                    <TableCell
                                                        colSpan={6}
                                                        align="center"
                                                    >
                                                        <Typography
                                                            color="textSecondary"
                                                            py={2}
                                                        >
                                                            Không có dữ liệu cho
                                                            ngày đã chọn
                                                        </Typography>
                                                    </TableCell>
                                                </TableRow>
                                            ) : (
                                                dailyData.employees.map(
                                                    (employee) => (
                                                        <TableRow
                                                            key={
                                                                employee.employee_id
                                                            }
                                                        >
                                                            <TableCell>
                                                                <Typography fontWeight="bold">
                                                                    {
                                                                        employee.name
                                                                    }
                                                                </Typography>
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.check_in ||
                                                                    '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.check_out ||
                                                                    '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                {formatHours(
                                                                    employee.total_hours
                                                                )}
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.late_minutes >
                                                                0
                                                                    ? `${employee.late_minutes}`
                                                                    : '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                <Chip
                                                                    label={getStatusText(
                                                                        employee.status
                                                                    )}
                                                                    color={getStatusColor(
                                                                        employee.status
                                                                    )}
                                                                    size="small"
                                                                />
                                                            </TableCell>
                                                        </TableRow>
                                                    )
                                                )
                                            )}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </Card>
                        </>
                    )}

                    {/* Monthly Report View */}
                    {mode === 'monthly' &&
                        monthlyData &&
                        monthlyData.statistics && (
                            <>
                                {/* Statistics Cards */}
                                <Grid container spacing={2} sx={{ mb: 3 }}>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng ngày làm việc
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                        >
                                                            {monthlyData.total_working_days ??
                                                                0}
                                                        </Typography>
                                                    </Box>
                                                    <CalendarMonth
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'primary.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tỷ lệ chấm công TB
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="success.main"
                                                        >
                                                            {(
                                                                monthlyData
                                                                    .statistics
                                                                    ?.average_attendance_rate ??
                                                                0
                                                            ).toFixed(1)}
                                                            %
                                                        </Typography>
                                                    </Box>
                                                    <CheckCircle
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'success.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng lần đi trễ
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="warning.main"
                                                        >
                                                            {monthlyData
                                                                .statistics
                                                                ?.total_late_instances ??
                                                                0}
                                                        </Typography>
                                                    </Box>
                                                    <Schedule
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'warning.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng giờ tăng ca
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="secondary.main"
                                                        >
                                                            {formatHours(
                                                                monthlyData
                                                                    .statistics
                                                                    ?.total_overtime_hours ??
                                                                    0
                                                            )}
                                                        </Typography>
                                                    </Box>
                                                    <AccessTime
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'secondary.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                </Grid>

                                {/* Employees Summary Table */}
                                <Card>
                                    <TableContainer>
                                        <Table>
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell>
                                                        Nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày có mặt
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày vắng mặt
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày đi trễ
                                                    </TableCell>
                                                    <TableCell>
                                                        Tổng giờ làm
                                                    </TableCell>
                                                    <TableCell>
                                                        TB giờ/ngày
                                                    </TableCell>
                                                    <TableCell>
                                                        Giờ tăng ca
                                                    </TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {!monthlyData.employees_summary ||
                                                monthlyData.employees_summary
                                                    .length === 0 ? (
                                                    <TableRow>
                                                        <TableCell
                                                            colSpan={7}
                                                            align="center"
                                                        >
                                                            <Typography
                                                                color="textSecondary"
                                                                py={2}
                                                            >
                                                                Không có dữ liệu
                                                                cho tháng đã
                                                                chọn
                                                            </Typography>
                                                        </TableCell>
                                                    </TableRow>
                                                ) : (
                                                    monthlyData.employees_summary.map(
                                                        (employee) => (
                                                            <TableRow
                                                                key={
                                                                    employee.employee_id
                                                                }
                                                            >
                                                                <TableCell>
                                                                    <Typography fontWeight="bold">
                                                                        {
                                                                            employee.name
                                                                        }
                                                                    </Typography>
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.present_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.absent_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.late_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.total_hours
                                                                    )}
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.average_hours_per_day
                                                                    )}
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.overtime_hours
                                                                    )}
                                                                </TableCell>
                                                            </TableRow>
                                                        )
                                                    )
                                                )}
                                            </TableBody>
                                        </Table>
                                    </TableContainer>
                                </Card>
                            </>
                        )}

                    {/* No data message */}
                    {((mode === 'daily' && !dailyData) ||
                        (mode === 'monthly' && !monthlyData)) && (
                        <Card>
                            <Box
                                display="flex"
                                justifyContent="center"
                                alignItems="center"
                                minHeight="200px"
                            >
                                <Typography color="textSecondary">
                                    Không có dữ liệu để hiển thị
                                </Typography>
                            </Box>
                        </Card>
                    )}
                </>
            )}

            {/* Export Dialog */}
            <Dialog
                open={exportDialogOpen}
                onClose={handleCloseExportDialog}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>Xuất báo cáo chấm công</DialogTitle>
                <DialogContent>
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'column',
                            gap: 2,
                            pt: 2,
                        }}
                    >
                        <FormControl fullWidth>
                            <InputLabel>Định dạng</InputLabel>
                            <Select
                                value={exportFormat}
                                label="Định dạng"
                                onChange={(e) =>
                                    setExportFormat(
                                        e.target.value as
                                            | 'excel'
                                            | 'pdf'
                                            | 'csv'
                                    )
                                }
                            >
                                <MenuItem value="excel">Excel (.xlsx)</MenuItem>
                                <MenuItem value="pdf">PDF (.pdf)</MenuItem>
                                <MenuItem value="csv">CSV (.csv)</MenuItem>
                            </Select>
                        </FormControl>
                        <TextField
                            fullWidth
                            label="Email (tùy chọn)"
                            type="email"
                            value={exportEmail}
                            onChange={(e) => setExportEmail(e.target.value)}
                            placeholder="Nhập email để nhận báo cáo qua email"
                            helperText="Để trống nếu muốn tải xuống trực tiếp"
                        />
                        <Typography variant="body2" color="textSecondary">
                            {mode === 'daily'
                                ? `Ngày: ${date}`
                                : `Tháng: ${month}`}
                        </Typography>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={handleCloseExportDialog}
                        disabled={exportLoading}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleExportReport}
                        variant="contained"
                        disabled={exportLoading}
                        startIcon={
                            exportLoading ? (
                                <CircularProgress size={20} />
                            ) : (
                                <Download />
                            )
                        }
                    >
                        {exportLoading ? 'Đang xử lý...' : 'Xuất báo cáo'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar for notifications */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={() => setSnackbar({ ...snackbar, open: false })}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={() => setSnackbar({ ...snackbar, open: false })}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};
