import { useEffect, useState, useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    CircularProgress,
    Alert,
    Avatar,
    Chip,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    IconButton,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Snackbar,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
} from '@mui/material';
import {
    ArrowBack,
    Edit,
    Save,
    Cancel,
    LockReset,
    Download,
    CalendarMonth,
} from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';

interface EmployeeData {
    user_id: string;
    role: number;
    avatar_url: string | null;
    email: string;
    phone: string;
    full_name: string;
    employee_code: string;
    department: string;
    hire_date: string;
    position: string;
    status: number;
}

interface DailySummary {
    employee_id: string;
    employee_name: string;
    date: string;
    total_hours: number;
    check_in: string;
    check_out: string;
    late_minutes: number;
    early_leave_minutes: number;
    status: string;
}

interface AttendanceRecord {
    record_id: string;
    employee_id: string;
    employee_name: string;
    check_in_time: string;
    check_out_time: string;
    device_id: string;
    location: string;
    status: string;
}

export const EmployeeDetailPage: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [employee, setEmployee] = useState<EmployeeData | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editField, setEditField] = useState<string | null>(null);
    const [editValue, setEditValue] = useState('');
    const [saving, setSaving] = useState(false);

    // Attendance data
    const [selectedMonth, setSelectedMonth] = useState(
        new Date().toISOString().slice(0, 7)
    );
    const [dailySummaries, setDailySummaries] = useState<DailySummary[]>([]);
    const [attendanceRecords, setAttendanceRecords] = useState<
        AttendanceRecord[]
    >([]);
    const [loadingAttendance, setLoadingAttendance] = useState(false);

    // Snackbar
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({ open: false, message: '', severity: 'success' });

    // Reset password dialog
    const [resetPasswordDialogOpen, setResetPasswordDialogOpen] =
        useState(false);
    const [resettingPassword, setResettingPassword] = useState(false);

    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
        }
    };

    // Fetch employee data
    const fetchEmployee = useCallback(async () => {
        if (!id) return;
        setLoading(true);
        setError(null);
        try {
            const response = await apiClient.get(`/api/v1/users/${id}`);
            if (response.data?.success && response.data?.data) {
                console.log('Fetched employee data:', response.data.data);
                setEmployee(response.data.data);
            } else {
                setError('Không thể tải thông tin nhân viên');
            }
        } catch (err: any) {
            console.error('Failed to fetch employee:', err);
            setError(
                err.response?.data?.message ||
                    'Không thể tải thông tin nhân viên'
            );
        } finally {
            setLoading(false);
        }
    }, [id]);

    // Fetch attendance data
    const fetchAttendanceData = useCallback(async () => {
        if (!id) return;
        const accessToken = localStorage.getItem('access_token');
        if (!accessToken) return;

        const companyId = getCompanyIdFromToken(accessToken);
        if (!companyId) return;

        setLoadingAttendance(true);
        try {
            // Fetch daily summaries
            const summariesResponse = await apiClient.get(
                `/api/v1/daily-summaries/user/${id}`,
                {
                    params: {
                        company_id: companyId,
                        year_month: selectedMonth,
                        limit: 100,
                    },
                }
            );

            if (summariesResponse.data?.success) {
                setDailySummaries(summariesResponse.data.data || []);
            }

            // Fetch attendance records
            const recordsResponse = await apiClient.get(
                `/api/v1/attendance-records/employee/${id}`,
                {
                    params: {
                        company_id: companyId,
                        year_month: selectedMonth,
                        limit: 100,
                    },
                }
            );

            if (recordsResponse.data?.success) {
                setAttendanceRecords(recordsResponse.data.data || []);
            }
        } catch (err: any) {
            console.error('Failed to fetch attendance data:', err);
        } finally {
            setLoadingAttendance(false);
        }
    }, [id, selectedMonth]);

    useEffect(() => {
        fetchEmployee();
    }, [fetchEmployee]);

    useEffect(() => {
        if (employee) {
            fetchAttendanceData();
        }
    }, [employee, fetchAttendanceData]);

    const handleEditField = (field: string, currentValue: string) => {
        setEditField(field);
        setEditValue(currentValue);
        setIsEditing(true);
    };

    const handleCancelEdit = () => {
        setEditField(null);
        setEditValue('');
        setIsEditing(false);
    };

    const handleSaveField = async () => {
        if (!id || !editField) return;

        setSaving(true);
        try {
            let endpoint = '';
            let payload: any = {};

            switch (editField) {
                case 'full_name':
                    endpoint = `/api/v1/users/${id}/name`;
                    payload = { full_name: editValue };
                    break;
                case 'phone':
                    endpoint = `/api/v1/users/${id}/phone`;
                    payload = { phone: editValue };
                    break;
                case 'department':
                    endpoint = `/api/v1/users/${id}/department`;
                    payload = { department: editValue };
                    break;
                case 'position':
                    endpoint = `/api/v1/users/${id}/position`;
                    payload = { position: editValue };
                    break;
                default:
                    return;
            }

            const response = await apiClient.put(endpoint, payload);
            if (response.data?.success) {
                setSnackbar({
                    open: true,
                    message: 'Cập nhật thành công',
                    severity: 'success',
                });
                await fetchEmployee();
                handleCancelEdit();
            } else {
                throw new Error('Cập nhật thất bại');
            }
        } catch (err: any) {
            console.error('Failed to update field:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message || 'Cập nhật thất bại',
                severity: 'error',
            });
        } finally {
            setSaving(false);
        }
    };

    const handleResetPassword = async () => {
        if (!id) return;
        setResettingPassword(true);
        try {
            const response = await apiClient.post(
                `/api/v1/users/${id}/reset-password`
            );
            if (response.data?.success) {
                setSnackbar({
                    open: true,
                    message: 'Đã gửi mật khẩu mới đến email nhân viên',
                    severity: 'success',
                });
                setResetPasswordDialogOpen(false);
            } else {
                throw new Error('Reset password thất bại');
            }
        } catch (err: any) {
            console.error('Failed to reset password:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message || 'Reset password thất bại',
                severity: 'error',
            });
        } finally {
            setResettingPassword(false);
        }
    };

    const handleExportData = () => {
        // Create CSV content
        const headers = [
            'Ngày',
            'Giờ vào',
            'Giờ ra',
            'Tổng giờ',
            'Phút muộn',
            'Phút về sớm',
            'Trạng thái',
        ];
        const rows = dailySummaries.map((summary) => [
            summary.date,
            summary.check_in || '',
            summary.check_out || '',
            summary.total_hours.toString(),
            summary.late_minutes.toString(),
            summary.early_leave_minutes.toString(),
            summary.status,
        ]);

        const csvContent =
            headers.join(',') +
            '\n' +
            rows.map((row) => row.join(',')).join('\n');

        const blob = new Blob(['\uFEFF' + csvContent], {
            type: 'text/csv;charset=utf-8;',
        });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `attendance_${employee?.employee_code}_${selectedMonth}.csv`;
        a.click();
        URL.revokeObjectURL(url);

        setSnackbar({
            open: true,
            message: 'Đã xuất dữ liệu thành công',
            severity: 'success',
        });
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'present':
            case 'on_time':
                return 'success';
            case 'late':
                return 'warning';
            case 'early_leave':
                return 'info';
            case 'absent':
                return 'error';
            default:
                return 'default';
        }
    };

    const getStatusText = (status: string) => {
        switch (status) {
            case 'present':
                return 'Có mặt';
            case 'on_time':
                return 'Đúng giờ';
            case 'late':
                return 'Trễ';
            case 'early_leave':
                return 'Về sớm';
            case 'absent':
                return 'Vắng mặt';
            default:
                return status;
        }
    };

    if (loading) {
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="400px"
            >
                <CircularProgress />
            </Box>
        );
    }

    if (error || !employee) {
        return (
            <Box>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/employees')}
                    sx={{ mb: 2 }}
                >
                    Quay lại
                </Button>
                <Alert severity="error">
                    {error || 'Không tìm thấy nhân viên'}
                </Alert>
            </Box>
        );
    }

    return (
        <>
            <Box>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/employees')}
                    sx={{ mb: 2 }}
                >
                    Quay lại
                </Button>

                <Typography variant="h4" fontWeight="bold" mb={3}>
                    Thông tin nhân viên
                </Typography>

                {/* Employee Info Card */}
                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Box display="flex" alignItems="center" gap={3} mb={3}>
                            <Avatar
                                src={employee.avatar_url || undefined}
                                sx={{ width: 80, height: 80 }}
                            >
                                {employee.full_name[0]}
                            </Avatar>
                            <Box>
                                <Typography variant="h5" fontWeight="bold">
                                    {employee.full_name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {employee.employee_code}
                                </Typography>
                                <Chip
                                    label={
                                        employee.status === 0
                                            ? 'Hoạt động'
                                            : 'Không hoạt động'
                                    }
                                    color={
                                        employee.status === 0
                                            ? 'success'
                                            : 'default'
                                    }
                                    size="small"
                                    sx={{ mt: 1 }}
                                />
                            </Box>
                        </Box>

                        <Grid container spacing={2}>
                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Họ và tên
                                    </Typography>
                                    {editField === 'full_name' ? (
                                        <Box display="flex" gap={1}>
                                            <TextField
                                                size="small"
                                                value={editValue}
                                                onChange={(e) =>
                                                    setEditValue(e.target.value)
                                                }
                                            />
                                            <IconButton
                                                size="small"
                                                color="primary"
                                                onClick={handleSaveField}
                                                disabled={saving}
                                            >
                                                <Save fontSize="small" />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={handleCancelEdit}
                                            >
                                                <Cancel fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    ) : (
                                        <Box display="flex" gap={1} alignItems="center">
                                            <Typography>{employee.full_name}</Typography>
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    handleEditField(
                                                        'full_name',
                                                        employee.full_name
                                                    )
                                                }
                                            >
                                                <Edit fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    )}
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Email
                                    </Typography>
                                    <Typography>{employee.email}</Typography>
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Số điện thoại
                                    </Typography>
                                    {editField === 'phone' ? (
                                        <Box display="flex" gap={1}>
                                            <TextField
                                                size="small"
                                                value={editValue}
                                                onChange={(e) =>
                                                    setEditValue(e.target.value)
                                                }
                                            />
                                            <IconButton
                                                size="small"
                                                color="primary"
                                                onClick={handleSaveField}
                                                disabled={saving}
                                            >
                                                <Save fontSize="small" />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={handleCancelEdit}
                                            >
                                                <Cancel fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    ) : (
                                        <Box display="flex" gap={1} alignItems="center">
                                            <Typography>{employee.phone || 'N/A'}</Typography>
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    handleEditField(
                                                        'phone',
                                                        employee.phone || ''
                                                    )
                                                }
                                            >
                                                <Edit fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    )}
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Phòng ban
                                    </Typography>
                                    {editField === 'department' ? (
                                        <Box display="flex" gap={1}>
                                            <TextField
                                                size="small"
                                                value={editValue}
                                                onChange={(e) =>
                                                    setEditValue(e.target.value)
                                                }
                                            />
                                            <IconButton
                                                size="small"
                                                color="primary"
                                                onClick={handleSaveField}
                                                disabled={saving}
                                            >
                                                <Save fontSize="small" />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={handleCancelEdit}
                                            >
                                                <Cancel fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    ) : (
                                        <Box display="flex" gap={1} alignItems="center">
                                            <Typography>
                                                {employee.department || 'N/A'}
                                            </Typography>
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    handleEditField(
                                                        'department',
                                                        employee.department || ''
                                                    )
                                                }
                                            >
                                                <Edit fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    )}
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Chức vụ
                                    </Typography>
                                    {editField === 'position' ? (
                                        <Box display="flex" gap={1}>
                                            <TextField
                                                size="small"
                                                value={editValue}
                                                onChange={(e) =>
                                                    setEditValue(e.target.value)
                                                }
                                            />
                                            <IconButton
                                                size="small"
                                                color="primary"
                                                onClick={handleSaveField}
                                                disabled={saving}
                                            >
                                                <Save fontSize="small" />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={handleCancelEdit}
                                            >
                                                <Cancel fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    ) : (
                                        <Box display="flex" gap={1} alignItems="center">
                                            <Typography>
                                                {employee.position || 'N/A'}
                                            </Typography>
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    handleEditField(
                                                        'position',
                                                        employee.position || ''
                                                    )
                                                }
                                            >
                                                <Edit fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    )}
                                </Box>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <Typography variant="body2" color="text.secondary">
                                        Ngày vào làm
                                    </Typography>
                                    <Typography>
                                        {employee.hire_date
                                            ? new Date(
                                                  employee.hire_date
                                              ).toLocaleDateString('vi-VN')
                                            : 'N/A'}
                                    </Typography>
                                </Box>
                            </Grid>

                            <Grid item xs={12}>
                                <Button
                                    variant="outlined"
                                    startIcon={<LockReset />}
                                    onClick={() => setResetPasswordDialogOpen(true)}
                                    sx={{ mt: 2 }}
                                >
                                    Reset mật khẩu
                                </Button>
                            </Grid>
                        </Grid>
                    </CardContent>
                </Card>

                {/* Attendance Section */}
                <Card>
                    <CardContent>
                        <Box
                            display="flex"
                            justifyContent="space-between"
                            alignItems="center"
                            mb={3}
                        >
                            <Typography variant="h6" fontWeight="bold">
                                Trạng thái chấm công
                            </Typography>
                            <Box display="flex" gap={2} alignItems="center">
                                <TextField
                                    type="month"
                                    label="Tháng"
                                    value={selectedMonth}
                                    onChange={(e) => setSelectedMonth(e.target.value)}
                                    InputLabelProps={{ shrink: true }}
                                    size="small"
                                />
                                <Button
                                    variant="outlined"
                                    startIcon={<Download />}
                                    onClick={handleExportData}
                                    disabled={dailySummaries.length === 0}
                                >
                                    Xuất dữ liệu
                                </Button>
                            </Box>
                        </Box>

                        {loadingAttendance ? (
                            <Box
                                display="flex"
                                justifyContent="center"
                                alignItems="center"
                                minHeight="200px"
                            >
                                <CircularProgress />
                            </Box>
                        ) : dailySummaries.length === 0 ? (
                            <Alert severity="info">
                                Không có dữ liệu chấm công cho tháng này
                            </Alert>
                        ) : (
                            <TableContainer component={Paper}>
                                <Table>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Ngày</TableCell>
                                            <TableCell>Giờ vào</TableCell>
                                            <TableCell>Giờ ra</TableCell>
                                            <TableCell>Tổng giờ</TableCell>
                                            <TableCell>Phút muộn</TableCell>
                                            <TableCell>Phút về sớm</TableCell>
                                            <TableCell>Trạng thái</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {dailySummaries.map((summary) => (
                                            <TableRow key={summary.date}>
                                                <TableCell>
                                                    {new Date(
                                                        summary.date
                                                    ).toLocaleDateString('vi-VN')}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.check_in || 'N/A'}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.check_out || 'N/A'}
                                                </TableCell>
                                                <TableCell>
                                                    {summary.total_hours.toFixed(2)}h
                                                </TableCell>
                                                <TableCell>
                                                    {summary.late_minutes} phút
                                                </TableCell>
                                                <TableCell>
                                                    {summary.early_leave_minutes} phút
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={getStatusText(
                                                            summary.status
                                                        )}
                                                        color={
                                                            getStatusColor(
                                                                summary.status
                                                            ) as any
                                                        }
                                                        size="small"
                                                    />
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        )}
                    </CardContent>
                </Card>
            </Box>

            {/* Reset Password Dialog */}
            <Dialog
                open={resetPasswordDialogOpen}
                onClose={() => setResetPasswordDialogOpen(false)}
            >
                <DialogTitle>Reset mật khẩu</DialogTitle>
                <DialogContent>
                    <Typography>
                        Bạn có chắc chắn muốn reset mật khẩu cho nhân viên{' '}
                        <strong>{employee.full_name}</strong>? Mật khẩu mới sẽ
                        được gửi đến email {employee.email}.
                    </Typography>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={() => setResetPasswordDialogOpen(false)}
                        disabled={resettingPassword}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleResetPassword}
                        variant="contained"
                        disabled={resettingPassword}
                    >
                        {resettingPassword ? 'Đang xử lý...' : 'Xác nhận'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar */}
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
        </>
    );
};

