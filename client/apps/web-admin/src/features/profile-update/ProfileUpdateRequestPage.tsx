import { useState, useEffect, useCallback } from 'react';
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
    TablePagination,
    Chip,
    Button,
    IconButton,
    CircularProgress,
    Alert,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Paper,
    Snackbar,
    Tooltip,
    Grid,
    Divider,
    Stack,
} from '@mui/material';
import {
    CheckCircle,
    Cancel,
    Refresh,
    Visibility,
    CalendarToday,
    Info,
} from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';
import { format } from 'date-fns';

// Types
interface ProfileUpdateRequest {
    request_id: string;
    user_id: string;
    company_id: string;
    status: 'pending' | 'approved' | 'rejected';
    request_month: string;
    request_count_in_month: number;
    reason: string;
    meta_data: {
        client_ip?: string;
        user_agent?: string;
    };
    created_at: string;
    updated_at: string;
}

interface PendingRequestsResponse {
    success: boolean;
    message: string;
    data: {
        requests: ProfileUpdateRequest[];
        total: number;
    };
}

interface ApiResponse {
    success: boolean;
    message: string;
    data?: any;
}

export const ProfileUpdateRequestPage: React.FC = () => {
    const [requests, setRequests] = useState<ProfileUpdateRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [total, setTotal] = useState(0);
    const [processingId, setProcessingId] = useState<string | null>(null);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error' | 'info';
    }>({ open: false, message: '', severity: 'success' });
    const [detailDialog, setDetailDialog] = useState<{
        open: boolean;
        request: ProfileUpdateRequest | null;
    }>({ open: false, request: null });
    const [actionDialog, setActionDialog] = useState<{
        open: boolean;
        type: 'approve' | 'reject' | null;
        requestId: string | null;
    }>({ open: false, type: null, requestId: null });

    // Fetch pending requests
    const fetchPendingRequests = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const offset = page * rowsPerPage;
            const response = await apiClient.get<PendingRequestsResponse>(
                '/api/v1/profile-update/requests/pending',
                {
                    params: {
                        limit: rowsPerPage,
                        offset: offset,
                    },
                }
            );

            if (response.data.success) {
                setRequests(response.data.data.requests);
                setTotal(response.data.data.total);
            } else {
                setError('Không thể tải danh sách yêu cầu');
            }
        } catch (err: any) {
            console.error('Failed to fetch pending requests:', err);
            setError(
                err.response?.data?.message ||
                    'Không thể tải danh sách yêu cầu. Vui lòng thử lại.'
            );
            setRequests([]);
        } finally {
            setLoading(false);
        }
    }, [page, rowsPerPage]);

    useEffect(() => {
        fetchPendingRequests();
    }, [fetchPendingRequests]);

    // Handle approve
    const handleApprove = async (requestId: string) => {
        setProcessingId(requestId);
        try {
            const response = await apiClient.post<ApiResponse>(
                `/api/v1/profile-update/requests/${requestId}/approve`
            );

            if (response.data.success) {
                setSnackbar({
                    open: true,
                    message: 'Phê duyệt yêu cầu thành công',
                    severity: 'success',
                });
                setActionDialog({ open: false, type: null, requestId: null });
                // Refresh the list
                await fetchPendingRequests();
            } else {
                setSnackbar({
                    open: true,
                    message: response.data.message || 'Phê duyệt thất bại',
                    severity: 'error',
                });
            }
        } catch (err: any) {
            console.error('Failed to approve request:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message ||
                    'Không thể phê duyệt yêu cầu. Vui lòng thử lại.',
                severity: 'error',
            });
        } finally {
            setProcessingId(null);
        }
    };

    // Handle reject
    const handleReject = async (requestId: string) => {
        setProcessingId(requestId);
        try {
            const response = await apiClient.post<ApiResponse>(
                `/api/v1/profile-update/requests/${requestId}/reject`
            );

            if (response.data.success) {
                setSnackbar({
                    open: true,
                    message: 'Từ chối yêu cầu thành công',
                    severity: 'success',
                });
                setActionDialog({ open: false, type: null, requestId: null });
                // Refresh the list
                await fetchPendingRequests();
            } else {
                setSnackbar({
                    open: true,
                    message: response.data.message || 'Từ chối thất bại',
                    severity: 'error',
                });
            }
        } catch (err: any) {
            console.error('Failed to reject request:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message ||
                    'Không thể từ chối yêu cầu. Vui lòng thử lại.',
                severity: 'error',
            });
        } finally {
            setProcessingId(null);
        }
    };

    // Handle page change
    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage);
    };

    // Handle rows per page change
    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    // Format date
    const formatDate = (dateString: string) => {
        try {
            return format(new Date(dateString), 'dd/MM/yyyy HH:mm');
        } catch {
            return dateString;
        }
    };

    // Open detail dialog
    const handleOpenDetail = (request: ProfileUpdateRequest) => {
        setDetailDialog({ open: true, request });
    };

    // Open action dialog
    const handleOpenAction = (
        type: 'approve' | 'reject',
        requestId: string
    ) => {
        setActionDialog({ open: true, type, requestId });
    };

    // Close dialogs
    const handleCloseDetail = () => {
        setDetailDialog({ open: false, request: null });
    };

    const handleCloseAction = () => {
        if (!processingId) {
            setActionDialog({ open: false, type: null, requestId: null });
        }
    };

    // Confirm action
    const handleConfirmAction = () => {
        if (actionDialog.requestId && actionDialog.type) {
            if (actionDialog.type === 'approve') {
                handleApprove(actionDialog.requestId);
            } else {
                handleReject(actionDialog.requestId);
            }
        }
    };

    // Close snackbar
    const handleCloseSnackbar = () => {
        setSnackbar({ ...snackbar, open: false });
    };

    if (loading && requests.length === 0) {
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

    return (
        <>
            <Box>
                <Box
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                    mb={3}
                >
                    <Typography variant="h4" fontWeight="bold">
                        Yêu cầu cập nhật khuôn mặt
                    </Typography>
                    <Button
                        variant="outlined"
                        startIcon={<Refresh />}
                        onClick={fetchPendingRequests}
                        disabled={loading}
                    >
                        Làm mới
                    </Button>
                </Box>

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
                    <CardContent>
                        <Box mb={2}>
                            <Typography variant="body2" color="text.secondary">
                                Tổng số yêu cầu đang chờ:{' '}
                                <strong>{total}</strong>
                            </Typography>
                        </Box>
                        <TableContainer component={Paper} variant="outlined">
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>ID Yêu cầu</TableCell>
                                        <TableCell>User ID</TableCell>
                                        <TableCell>Tháng yêu cầu</TableCell>
                                        <TableCell>
                                            Số lần trong tháng
                                        </TableCell>
                                        <TableCell>Lý do</TableCell>
                                        <TableCell>Ngày tạo</TableCell>
                                        <TableCell>Trạng thái</TableCell>
                                        <TableCell align="right">
                                            Thao tác
                                        </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {requests.length === 0 ? (
                                        <TableRow>
                                            <TableCell
                                                colSpan={8}
                                                align="center"
                                                sx={{ py: 4 }}
                                            >
                                                <Typography
                                                    variant="body2"
                                                    color="text.secondary"
                                                >
                                                    Không có yêu cầu nào đang
                                                    chờ
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        requests.map((request) => (
                                            <TableRow
                                                key={request.request_id}
                                                hover={
                                                    processingId !==
                                                    request.request_id
                                                }
                                                sx={{
                                                    '&:hover': {
                                                        backgroundColor:
                                                            processingId ===
                                                            request.request_id
                                                                ? 'transparent'
                                                                : 'action.hover',
                                                    },
                                                    backgroundColor:
                                                        processingId ===
                                                        request.request_id
                                                            ? 'action.selected'
                                                            : 'transparent',
                                                    opacity:
                                                        processingId ===
                                                        request.request_id
                                                            ? 0.7
                                                            : 1,
                                                }}
                                            >
                                                <TableCell>
                                                    <Typography
                                                        variant="body2"
                                                        sx={{
                                                            fontFamily:
                                                                'monospace',
                                                            fontSize: '0.75rem',
                                                        }}
                                                    >
                                                        {request.request_id.slice(
                                                            0,
                                                            8
                                                        )}
                                                        ...
                                                    </Typography>
                                                </TableCell>
                                                <TableCell>
                                                    <Typography
                                                        variant="body2"
                                                        sx={{
                                                            fontFamily:
                                                                'monospace',
                                                            fontSize: '0.75rem',
                                                        }}
                                                    >
                                                        {request.user_id.slice(
                                                            0,
                                                            8
                                                        )}
                                                        ...
                                                    </Typography>
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={
                                                            request.request_month
                                                        }
                                                        size="small"
                                                        icon={<CalendarToday />}
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={
                                                            request.request_count_in_month
                                                        }
                                                        size="small"
                                                        color="warning"
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    <Tooltip
                                                        title={request.reason}
                                                    >
                                                        <Typography
                                                            variant="body2"
                                                            sx={{
                                                                maxWidth: 200,
                                                                overflow:
                                                                    'hidden',
                                                                textOverflow:
                                                                    'ellipsis',
                                                                whiteSpace:
                                                                    'nowrap',
                                                            }}
                                                        >
                                                            {request.reason}
                                                        </Typography>
                                                    </Tooltip>
                                                </TableCell>
                                                <TableCell>
                                                    {formatDate(
                                                        request.created_at
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {processingId ===
                                                    request.request_id ? (
                                                        <Chip
                                                            label="Đang xử lý..."
                                                            color="info"
                                                            size="small"
                                                            icon={
                                                                <CircularProgress
                                                                    size={12}
                                                                    sx={{
                                                                        color: 'inherit',
                                                                    }}
                                                                />
                                                            }
                                                        />
                                                    ) : (
                                                        <Chip
                                                            label={
                                                                request.status ===
                                                                'pending'
                                                                    ? 'Chờ xử lý'
                                                                    : request.status ===
                                                                      'approved'
                                                                    ? 'Đã phê duyệt'
                                                                    : 'Đã từ chối'
                                                            }
                                                            color={
                                                                request.status ===
                                                                'pending'
                                                                    ? 'warning'
                                                                    : request.status ===
                                                                      'approved'
                                                                    ? 'success'
                                                                    : 'error'
                                                            }
                                                            size="small"
                                                        />
                                                    )}
                                                </TableCell>
                                                <TableCell align="right">
                                                    <Stack
                                                        direction="row"
                                                        spacing={1}
                                                        justifyContent="flex-end"
                                                    >
                                                        <Tooltip
                                                            title={
                                                                processingId ===
                                                                request.request_id
                                                                    ? 'Đang xử lý...'
                                                                    : 'Xem chi tiết'
                                                            }
                                                        >
                                                            <span>
                                                                <IconButton
                                                                    size="small"
                                                                    onClick={() =>
                                                                        handleOpenDetail(
                                                                            request
                                                                        )
                                                                    }
                                                                    disabled={
                                                                        processingId ===
                                                                        request.request_id
                                                                    }
                                                                    sx={{
                                                                        color:
                                                                            processingId ===
                                                                            request.request_id
                                                                                ? 'action.disabled'
                                                                                : 'inherit',
                                                                    }}
                                                                >
                                                                    <Visibility />
                                                                </IconButton>
                                                            </span>
                                                        </Tooltip>
                                                        <Tooltip
                                                            title={
                                                                processingId ===
                                                                request.request_id
                                                                    ? 'Đang xử lý...'
                                                                    : 'Phê duyệt'
                                                            }
                                                        >
                                                            <span>
                                                                <IconButton
                                                                    size="small"
                                                                    color="success"
                                                                    onClick={() =>
                                                                        handleOpenAction(
                                                                            'approve',
                                                                            request.request_id
                                                                        )
                                                                    }
                                                                    disabled={
                                                                        processingId ===
                                                                        request.request_id
                                                                    }
                                                                    sx={{
                                                                        color:
                                                                            processingId ===
                                                                            request.request_id
                                                                                ? 'action.disabled'
                                                                                : 'success.main',
                                                                    }}
                                                                >
                                                                    {processingId ===
                                                                    request.request_id ? (
                                                                        <CircularProgress
                                                                            size={
                                                                                16
                                                                            }
                                                                        />
                                                                    ) : (
                                                                        <CheckCircle />
                                                                    )}
                                                                </IconButton>
                                                            </span>
                                                        </Tooltip>
                                                        <Tooltip
                                                            title={
                                                                processingId ===
                                                                request.request_id
                                                                    ? 'Đang xử lý...'
                                                                    : 'Từ chối'
                                                            }
                                                        >
                                                            <span>
                                                                <IconButton
                                                                    size="small"
                                                                    color="error"
                                                                    onClick={() =>
                                                                        handleOpenAction(
                                                                            'reject',
                                                                            request.request_id
                                                                        )
                                                                    }
                                                                    disabled={
                                                                        processingId ===
                                                                        request.request_id
                                                                    }
                                                                    sx={{
                                                                        color:
                                                                            processingId ===
                                                                            request.request_id
                                                                                ? 'action.disabled'
                                                                                : 'error.main',
                                                                    }}
                                                                >
                                                                    {processingId ===
                                                                    request.request_id ? (
                                                                        <CircularProgress
                                                                            size={
                                                                                16
                                                                            }
                                                                        />
                                                                    ) : (
                                                                        <Cancel />
                                                                    )}
                                                                </IconButton>
                                                            </span>
                                                        </Tooltip>
                                                    </Stack>
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                        <TablePagination
                            component="div"
                            count={total}
                            page={page}
                            onPageChange={handleChangePage}
                            rowsPerPage={rowsPerPage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            rowsPerPageOptions={[5, 10, 25, 50]}
                            labelRowsPerPage="Số dòng mỗi trang:"
                            labelDisplayedRows={({ from, to, count }) =>
                                `${from}-${to} của ${
                                    count !== -1 ? count : `nhiều hơn ${to}`
                                }`
                            }
                        />
                    </CardContent>
                </Card>
            </Box>

            {/* Detail Dialog */}
            <Dialog
                open={detailDialog.open}
                onClose={handleCloseDetail}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>
                    <Box display="flex" alignItems="center" gap={1}>
                        <Info color="primary" />
                        <Typography variant="h6">Chi tiết yêu cầu</Typography>
                    </Box>
                </DialogTitle>
                <DialogContent>
                    {detailDialog.request && (
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <Divider sx={{ my: 1 }} />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Request ID
                                </Typography>
                                <Typography
                                    variant="body2"
                                    sx={{ fontFamily: 'monospace', mt: 0.5 }}
                                >
                                    {detailDialog.request.request_id}
                                </Typography>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    User ID
                                </Typography>
                                <Typography
                                    variant="body2"
                                    sx={{ fontFamily: 'monospace', mt: 0.5 }}
                                >
                                    {detailDialog.request.user_id}
                                </Typography>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Company ID
                                </Typography>
                                <Typography
                                    variant="body2"
                                    sx={{ fontFamily: 'monospace', mt: 0.5 }}
                                >
                                    {detailDialog.request.company_id}
                                </Typography>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Trạng thái
                                </Typography>
                                <Box mt={0.5}>
                                    <Chip
                                        label={
                                            detailDialog.request.status ===
                                            'pending'
                                                ? 'Chờ xử lý'
                                                : detailDialog.request
                                                      .status === 'approved'
                                                ? 'Đã phê duyệt'
                                                : 'Đã từ chối'
                                        }
                                        color={
                                            detailDialog.request.status ===
                                            'pending'
                                                ? 'warning'
                                                : detailDialog.request
                                                      .status === 'approved'
                                                ? 'success'
                                                : 'error'
                                        }
                                        size="small"
                                    />
                                </Box>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Tháng yêu cầu
                                </Typography>
                                <Typography variant="body2" mt={0.5}>
                                    {detailDialog.request.request_month}
                                </Typography>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Số lần yêu cầu trong tháng
                                </Typography>
                                <Typography variant="body2" mt={0.5}>
                                    {
                                        detailDialog.request
                                            .request_count_in_month
                                    }
                                </Typography>
                            </Grid>
                            <Grid item xs={12}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Lý do
                                </Typography>
                                <Typography variant="body2" mt={0.5}>
                                    {detailDialog.request.reason}
                                </Typography>
                            </Grid>
                            <Grid item xs={12}>
                                <Divider sx={{ my: 1 }} />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Ngày tạo
                                </Typography>
                                <Typography variant="body2" mt={0.5}>
                                    {formatDate(
                                        detailDialog.request.created_at
                                    )}
                                </Typography>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                >
                                    Ngày cập nhật
                                </Typography>
                                <Typography variant="body2" mt={0.5}>
                                    {formatDate(
                                        detailDialog.request.updated_at
                                    )}
                                </Typography>
                            </Grid>
                            {detailDialog.request.meta_data && (
                                <>
                                    <Grid item xs={12}>
                                        <Divider sx={{ my: 1 }} />
                                        <Typography
                                            variant="subtitle2"
                                            gutterBottom
                                        >
                                            Thông tin metadata
                                        </Typography>
                                    </Grid>
                                    {detailDialog.request.meta_data
                                        .client_ip && (
                                        <Grid item xs={12} sm={6}>
                                            <Typography
                                                variant="caption"
                                                color="text.secondary"
                                            >
                                                Client IP
                                            </Typography>
                                            <Typography
                                                variant="body2"
                                                sx={{
                                                    fontFamily: 'monospace',
                                                    mt: 0.5,
                                                }}
                                            >
                                                {
                                                    detailDialog.request
                                                        .meta_data.client_ip
                                                }
                                            </Typography>
                                        </Grid>
                                    )}
                                    {detailDialog.request.meta_data
                                        .user_agent && (
                                        <Grid item xs={12}>
                                            <Typography
                                                variant="caption"
                                                color="text.secondary"
                                            >
                                                User Agent
                                            </Typography>
                                            <Typography
                                                variant="body2"
                                                sx={{
                                                    fontFamily: 'monospace',
                                                    mt: 0.5,
                                                }}
                                            >
                                                {
                                                    detailDialog.request
                                                        .meta_data.user_agent
                                                }
                                            </Typography>
                                        </Grid>
                                    )}
                                </>
                            )}
                        </Grid>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleCloseDetail}>Đóng</Button>
                </DialogActions>
            </Dialog>

            {/* Action Confirmation Dialog */}
            <Dialog
                open={actionDialog.open}
                onClose={handleCloseAction}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>
                    {actionDialog.type === 'approve'
                        ? 'Xác nhận phê duyệt'
                        : 'Xác nhận từ chối'}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {actionDialog.type === 'approve'
                            ? 'Bạn có chắc chắn muốn phê duyệt yêu cầu cập nhật khuôn mặt này không?'
                            : 'Bạn có chắc chắn muốn từ chối yêu cầu cập nhật khuôn mặt này không?'}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={handleCloseAction}
                        disabled={!!processingId}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleConfirmAction}
                        variant="contained"
                        color={
                            actionDialog.type === 'approve'
                                ? 'success'
                                : 'error'
                        }
                        disabled={!!processingId}
                        startIcon={
                            processingId ? (
                                <CircularProgress size={16} />
                            ) : actionDialog.type === 'approve' ? (
                                <CheckCircle />
                            ) : (
                                <Cancel />
                            )
                        }
                    >
                        {processingId
                            ? 'Đang xử lý...'
                            : actionDialog.type === 'approve'
                            ? 'Phê duyệt'
                            : 'Từ chối'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={handleCloseSnackbar}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={handleCloseSnackbar}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </>
    );
};
