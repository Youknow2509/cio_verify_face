import { useState, useEffect, useCallback } from 'react';
import {
    Box,
    Card,
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
    Pagination,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Snackbar,
} from '@mui/material';
import { Download } from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';

// Types matching the Go structs
interface DailyReportDetailEmployee {
    company_id: string;
    summary_month: string; // YYYY-MM format
    work_date: string; // ISO date string
    employee_id: string;
    shift_id: string;
    actual_check_in: string | null; // ISO datetime string or null
    actual_check_out: string | null; // ISO datetime string or null
    attendance_status: number;
    late_minutes: number;
    early_leave_minutes: number;
    total_work_minutes: number;
    notes: string;
    updated_at: string; // ISO datetime string
    overtime_minutes: number;
    attendance_percentage: number;
}

interface DailyReportDetailsResponse {
    total: number;
    items: DailyReportDetailEmployee[];
    next_page?: string;
}

// Extended type with employee and shift info (will be fetched separately if needed)
interface DailyReportRecord extends DailyReportDetailEmployee {
    employee_name?: string;
    employee_code?: string;
    shift_name?: string;
}

export const DailyReportPage: React.FC = () => {
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
    const [records, setRecords] = useState<DailyReportRecord[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [nextPage, setNextPage] = useState<string | undefined>(undefined);
    const [employeeMap, setEmployeeMap] = useState<
        Map<string, { name: string; code: string }>
    >(new Map());
    const [shiftMap, setShiftMap] = useState<Map<string, string>>(new Map());

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

    // Fetch daily report data
    const fetchDailyReport = useCallback(
        async (selectedDate: string, pageToken?: string) => {
            setLoading(true);
            setError(null);
            try {
                const requestBody: any = {
                    work_date: selectedDate,
                };
                if (pageToken) {
                    requestBody.next_page = pageToken;
                }

                const response = await apiClient.post<{
                    data: DailyReportDetailsResponse;
                }>('/api/v1/daily-summaries/details', requestBody);

                if (response.data?.data) {
                    const data = response.data.data;
                    setTotal(data.total);
                    setNextPage(data.next_page);

                    // Fetch employee and shift info for all records
                    const enrichedRecords = await Promise.all(
                        data.items.map(async (item) => {
                            // Check cache first
                            let employeeInfo = employeeMap.get(
                                item.employee_id
                            );
                            let shiftName = shiftMap.get(item.shift_id);

                            // Fetch employee info if not in cache
                            if (!employeeInfo) {
                                try {
                                    const empResponse = await apiClient.get(
                                        `/api/v1/users/${item.employee_id}`
                                    );
                                    if (
                                        empResponse.status === 200 &&
                                        empResponse.data?.data
                                    ) {
                                        const emp = empResponse.data.data;
                                        employeeInfo = {
                                            name: emp.full_name || '',
                                            code: emp.employee_code || '',
                                        };
                                        setEmployeeMap((prev) => {
                                            const newMap = new Map(prev);
                                            newMap.set(
                                                item.employee_id,
                                                employeeInfo!
                                            );
                                            return newMap;
                                        });
                                    }
                                } catch (err) {
                                    console.error(
                                        'Failed to fetch employee:',
                                        err
                                    );
                                    employeeInfo = { name: '', code: '' };
                                }
                            }

                            // Fetch shift info if not in cache
                            if (!shiftName) {
                                try {
                                    const shiftResponse = await apiClient.get(
                                        `/api/v1/shift/${item.shift_id}`
                                    );
                                    if (
                                        shiftResponse.status === 200 &&
                                        shiftResponse.data?.data
                                    ) {
                                        shiftName =
                                            shiftResponse.data.data.name || '';
                                        setShiftMap((prev) => {
                                            const newMap = new Map(prev);
                                            newMap.set(
                                                item.shift_id,
                                                shiftName!
                                            );
                                            return newMap;
                                        });
                                    }
                                } catch (err) {
                                    console.error(
                                        'Failed to fetch shift:',
                                        err
                                    );
                                    shiftName = '';
                                }
                            }

                            return {
                                ...item,
                                employee_name: employeeInfo?.name || '',
                                employee_code: employeeInfo?.code || '',
                                shift_name: shiftName || '',
                            };
                        })
                    );

                    setRecords(enrichedRecords);
                }
            } catch (err: any) {
                console.error('Failed to fetch daily report:', err);
                setError(
                    err.response?.data?.message ||
                        err.response?.data?.error ||
                        'Không thể tải báo cáo. Vui lòng thử lại.'
                );
                setRecords([]);
            } finally {
                setLoading(false);
            }
        },
        [employeeMap, shiftMap]
    );

    useEffect(() => {
        fetchDailyReport(date);
    }, [date, fetchDailyReport]);

    const handleDateChange = (newDate: string) => {
        setDate(newDate);
        setPage(1);
        setNextPage(undefined);
    };

    const handlePageChange = (
        _event: React.ChangeEvent<unknown>,
        value: number
    ) => {
        setPage(value);
        if (value > page && nextPage) {
            fetchDailyReport(date, nextPage);
        } else if (value < page) {
            // For previous page, we might need to refetch from the beginning
            // This is a simplified implementation - adjust based on your API's pagination strategy
            fetchDailyReport(date);
        }
    };

    const getStatusColor = (status: number) => {
        // attendance_status: 0 = absent, 1 = present, 2 = late, 3 = early leave, etc.
        switch (status) {
            case 1:
                return 'success';
            case 2:
                return 'error';
            case 3:
                return 'warning';
            default:
                return 'default';
        }
    };

    const getStatusText = (status: number) => {
        switch (status) {
            case 0:
                return 'Vắng mặt';
            case 1:
                return 'Có mặt';
            case 2:
                return 'Trễ';
            case 3:
                return 'Về sớm';
            default:
                return '-';
        }
    };

    const formatMinutesToHours = (minutes: number) => {
        if (!minutes) return '-';
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        if (hours > 0 && mins > 0) {
            return `${hours}h${mins}m`;
        } else if (hours > 0) {
            return `${hours}h`;
        } else {
            return `${mins}m`;
        }
    };

    // Get company ID from JWT token
    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
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
            const requestBody: {
                company_id: string;
                date: string;
                format: string;
                email?: string;
            } = {
                company_id: companyId,
                date: date, // Using YYYY-MM-DD format as shown in example
                format: exportFormat,
            };

            const hasEmail = exportEmail.trim() !== '';
            if (hasEmail) {
                requestBody.email = exportEmail.trim();
            }

            // If email is provided, expect JSON response; otherwise expect blob
            const response = await apiClient.post(
                '/api/v1/reports/daily/export',
                requestBody,
                hasEmail
                    ? {} // JSON response when email is provided
                    : {
                          responseType: 'blob', // Blob response for file download
                      }
            );

            if (hasEmail) {
                // Server sent email, expect JSON response
                const message =
                    response.data?.message ||
                    `Báo cáo đã được gửi đến email ${exportEmail}`;
                setSnackbar({
                    open: true,
                    message,
                    severity: 'success',
                });
                setExportDialogOpen(false);
                setExportEmail('');
            } else {
                // Download file directly
                // Check if response is blob or needs conversion
                let blob: Blob;
                if (response.data instanceof Blob) {
                    blob = response.data;
                } else {
                    blob = new Blob([response.data], {
                        type:
                            response.headers['content-type'] ||
                            'application/octet-stream',
                    });
                }

                const url = window.URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.href = url;

                // Determine file extension based on format
                const extension =
                    exportFormat === 'excel' ? 'xlsx' : exportFormat;
                const fileName = `daily-report-${date}.${extension}`;
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
                setExportDialogOpen(false);
            }
        } catch (err: any) {
            console.error('Failed to export report:', err);

            // Try to parse error message from blob response if needed
            let errorMessage =
                err.response?.data?.message ||
                err.response?.data?.error ||
                'Không thể xuất báo cáo. Vui lòng thử lại.';

            // If response is blob but error occurred, try to parse it
            if (err.response?.data instanceof Blob) {
                try {
                    const text = await err.response.data.text();
                    const json = JSON.parse(text);
                    errorMessage = json.message || json.error || errorMessage;
                } catch {
                    // Ignore parsing errors
                }
            }

            setSnackbar({
                open: true,
                message: errorMessage,
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
                Báo cáo Chấm công Hàng ngày
            </Typography>
            <Card sx={{ mb: 3, p: 2 }}>
                <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} md={4}>
                        <TextField
                            fullWidth
                            label="Ngày"
                            type="date"
                            value={date}
                            onChange={(e) => handleDateChange(e.target.value)}
                            InputLabelProps={{ shrink: true }}
                        />
                    </Grid>
                    <Grid item xs={12} md={8}>
                        <Button
                            variant="contained"
                            startIcon={<Download />}
                            onClick={handleOpenExportDialog}
                            disabled={loading}
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

            <Card>
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
                        <TableContainer>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Nhân viên</TableCell>
                                        <TableCell>Ca làm việc</TableCell>
                                        <TableCell>Giờ vào</TableCell>
                                        <TableCell>Giờ ra</TableCell>
                                        <TableCell>Tổng giờ làm</TableCell>
                                        <TableCell>Trễ (phút)</TableCell>
                                        <TableCell>Về sớm (phút)</TableCell>
                                        <TableCell>Tăng ca (phút)</TableCell>
                                        <TableCell>Trạng thái</TableCell>
                                        <TableCell>Tỷ lệ chính xác</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {records.length === 0 ? (
                                        <TableRow>
                                            <TableCell
                                                colSpan={10}
                                                align="center"
                                            >
                                                <Typography
                                                    color="textSecondary"
                                                    py={2}
                                                >
                                                    Không có dữ liệu cho ngày đã
                                                    chọn
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        records.map((record, index) => (
                                            <TableRow
                                                key={`${record.employee_id}-${record.shift_id}-${record.work_date}-${index}`}
                                            >
                                                <TableCell>
                                                    <Box>
                                                        <Typography fontWeight="bold">
                                                            {record.employee_name ||
                                                                record.employee_id}
                                                        </Typography>
                                                        {record.employee_code && (
                                                            <Typography
                                                                variant="caption"
                                                                color="textSecondary"
                                                            >
                                                                {
                                                                    record.employee_code
                                                                }
                                                            </Typography>
                                                        )}
                                                    </Box>
                                                </TableCell>
                                                <TableCell>
                                                    {record.shift_name ||
                                                        record.shift_id}
                                                </TableCell>
                                                <TableCell>
                                                    {record.actual_check_in
                                                        ? new Date(
                                                              record.actual_check_in
                                                          ).toLocaleTimeString(
                                                              'vi-VN',
                                                              {
                                                                  hour: '2-digit',
                                                                  minute: '2-digit',
                                                              }
                                                          )
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    {record.actual_check_out
                                                        ? new Date(
                                                              record.actual_check_out
                                                          ).toLocaleTimeString(
                                                              'vi-VN',
                                                              {
                                                                  hour: '2-digit',
                                                                  minute: '2-digit',
                                                              }
                                                          )
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    {formatMinutesToHours(
                                                        record.total_work_minutes
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {record.late_minutes > 0
                                                        ? `${record.late_minutes}`
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    {record.early_leave_minutes >
                                                    0
                                                        ? `${record.early_leave_minutes}`
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    {record.overtime_minutes > 0
                                                        ? formatMinutesToHours(
                                                              record.overtime_minutes
                                                          )
                                                        : '-'}
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={getStatusText(
                                                            record.attendance_status
                                                        )}
                                                        color={getStatusColor(
                                                            record.attendance_status
                                                        )}
                                                        size="small"
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    {record.attendance_percentage >
                                                    0
                                                        ? `${(
                                                              record.attendance_percentage *
                                                              100
                                                          ).toFixed(1)}%`
                                                        : '-'}
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                        {total > 0 && (
                            <Box
                                display="flex"
                                justifyContent="space-between"
                                alignItems="center"
                                p={2}
                            >
                                <Typography
                                    variant="body2"
                                    color="textSecondary"
                                >
                                    Tổng: {total} bản ghi
                                </Typography>
                                {nextPage && (
                                    <Pagination
                                        count={Math.ceil(total / 10)} // Assuming 10 items per page, adjust as needed
                                        page={page}
                                        onChange={handlePageChange}
                                        color="primary"
                                    />
                                )}
                            </Box>
                        )}
                    </>
                )}
            </Card>

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
                            Ngày: {date}
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
