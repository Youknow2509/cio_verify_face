import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Button,
    Card,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
    Typography,
    Chip,
    Tooltip,
    Stack,
    CircularProgress,
    Alert,
    Snackbar,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Pagination,
} from '@mui/material';
import {
    Add,
    Edit,
    Delete,
    AccessTime,
    People,
    CheckCircle,
    Cancel,
    Assignment,
    ToggleOn,
    ToggleOff,
} from '@mui/icons-material';
import type { Shift } from '@face-attendance/types';
import {
    getShifts,
    deleteShift,
    changeShiftStatus,
} from '@face-attendance/utils';

const DAYS_OF_WEEK_SHORT = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];

export const ShiftListPage: React.FC = () => {
    const navigate = useNavigate();
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    // Snackbar state
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({
        open: false,
        message: '',
        severity: 'success',
    });

    // Delete confirmation dialog state
    const [deleteDialog, setDeleteDialog] = useState<{
        open: boolean;
        shiftId: string | null;
        shiftName: string;
    }>({
        open: false,
        shiftId: null,
        shiftName: '',
    });

    // Fetch shifts from API
    const fetchShifts = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await getShifts(page);
            if (response.code === 200 && response.data) {
                setShifts(response.data);
                // If API returns pagination info, set it here
                // setTotalPages(response.pagination?.total_pages || 1);
            } else {
                setError(
                    response.message || 'Không thể tải danh sách ca làm việc'
                );
            }
        } catch (err: any) {
            console.error('Error fetching shifts:', err);
            setError(
                err.response?.data?.error || 'Đã xảy ra lỗi khi tải dữ liệu'
            );
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchShifts();
    }, [page]);

    const calculateWorkHours = (shift: Shift): number => {
        const [startHour, startMin] = shift.start_time.split(':').map(Number);
        const [endHour, endMin] = shift.end_time.split(':').map(Number);
        const totalMinutes =
            endHour * 60 +
            endMin -
            (startHour * 60 + startMin) -
            shift.break_duration_minutes;
        return Math.max(0, totalMinutes / 60);
    };

    const handleDeleteClick = (shift: Shift) => {
        setDeleteDialog({
            open: true,
            shiftId: shift.shift_id,
            shiftName: shift.name,
        });
    };

    const handleDeleteConfirm = async () => {
        if (!deleteDialog.shiftId) return;

        try {
            const response = await deleteShift(deleteDialog.shiftId);
            if (response.code === 200) {
                setSnackbar({
                    open: true,
                    message: 'Xóa ca làm việc thành công',
                    severity: 'success',
                });
                fetchShifts(); // Refresh list
            } else {
                throw new Error(response.message || 'Xóa thất bại');
            }
        } catch (err: any) {
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.error || 'Không thể xóa ca làm việc',
                severity: 'error',
            });
        } finally {
            setDeleteDialog({ open: false, shiftId: null, shiftName: '' });
        }
    };

    const handleToggleStatus = async (shift: Shift) => {
        try {
            const newStatus = shift.is_active ? 0 : 1;
            const response = await changeShiftStatus({
                company_id: shift.company_id,
                shift_id: shift.shift_id,
                status: newStatus,
            });

            if (response.code === 200) {
                setSnackbar({
                    open: true,
                    message: `${
                        newStatus ? 'Kích hoạt' : 'Tạm dừng'
                    } ca làm việc thành công`,
                    severity: 'success',
                });
                fetchShifts(); // Refresh list
            } else {
                throw new Error(
                    response.message || 'Thay đổi trạng thái thất bại'
                );
            }
        } catch (err: any) {
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.error ||
                    'Không thể thay đổi trạng thái',
                severity: 'error',
            });
        }
    };

    const handlePageChange = (
        _event: React.ChangeEvent<unknown>,
        value: number
    ) => {
        setPage(value);
    };

    const handleCloseSnackbar = () => {
        setSnackbar({ ...snackbar, open: false });
    };

    return (
        <Box>
            <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                mb={3}
            >
                <Typography variant="h4" fontWeight="bold">
                    Quản lý Ca làm việc
                </Typography>
                <Box display="flex" gap={2}>
                    <Button
                        variant="outlined"
                        startIcon={<Assignment />}
                        onClick={() => navigate('/shifts/assign')}
                    >
                        Phân công ca
                    </Button>
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        onClick={() => navigate('/shifts/add')}
                    >
                        Thêm ca làm việc
                    </Button>
                </Box>
            </Box>

            {/* Error Alert */}
            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            <Card>
                {loading ? (
                    <Box display="flex" justifyContent="center" p={5}>
                        <CircularProgress />
                    </Box>
                ) : shifts.length === 0 ? (
                    <Box p={5} textAlign="center">
                        <Typography color="text.secondary">
                            Chưa có ca làm việc nào
                        </Typography>
                        <Button
                            variant="contained"
                            startIcon={<Add />}
                            onClick={() => navigate('/shifts/add')}
                            sx={{ mt: 2 }}
                        >
                            Tạo ca làm việc đầu tiên
                        </Button>
                    </Box>
                ) : (
                    <>
                        <TableContainer>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Tên ca</TableCell>
                                        <TableCell>Giờ làm việc</TableCell>
                                        <TableCell>Ngày làm việc</TableCell>
                                        <TableCell align="center">
                                            Chính sách
                                        </TableCell>
                                        <TableCell align="center">
                                            Nhân viên
                                        </TableCell>
                                        <TableCell align="center">
                                            Trạng thái
                                        </TableCell>
                                        <TableCell align="right">
                                            Thao tác
                                        </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {shifts.map((shift) => (
                                        <TableRow key={shift.shift_id} hover>
                                            <TableCell>
                                                <Box>
                                                    <Typography fontWeight="bold">
                                                        {shift.name}
                                                    </Typography>
                                                    {shift.description && (
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            {shift.description}
                                                        </Typography>
                                                    )}
                                                    {shift.is_flexible && (
                                                        <Chip
                                                            label="Linh hoạt"
                                                            size="small"
                                                            color="info"
                                                            sx={{ mt: 0.5 }}
                                                        />
                                                    )}
                                                </Box>
                                            </TableCell>

                                            <TableCell>
                                                <Stack spacing={0.5}>
                                                    <Box
                                                        display="flex"
                                                        alignItems="center"
                                                        gap={1}
                                                    >
                                                        <AccessTime
                                                            fontSize="small"
                                                            color="action"
                                                        />
                                                        <Typography variant="body2">
                                                            {shift.start_time} -{' '}
                                                            {shift.end_time}
                                                        </Typography>
                                                    </Box>
                                                    <Typography
                                                        variant="caption"
                                                        color="text.secondary"
                                                    >
                                                        {calculateWorkHours(
                                                            shift
                                                        ).toFixed(1)}
                                                        h làm việc, nghỉ{' '}
                                                        {
                                                            shift.break_duration_minutes
                                                        }
                                                        p
                                                    </Typography>
                                                </Stack>
                                            </TableCell>

                                            <TableCell>
                                                <Box
                                                    display="flex"
                                                    gap={0.5}
                                                    flexWrap="wrap"
                                                >
                                                    {shift.work_days
                                                        .sort()
                                                        .map((day) => (
                                                            <Chip
                                                                key={day}
                                                                label={
                                                                    DAYS_OF_WEEK_SHORT[
                                                                        day - 1
                                                                    ]
                                                                }
                                                                size="small"
                                                                variant="outlined"
                                                            />
                                                        ))}
                                                </Box>
                                            </TableCell>

                                            <TableCell>
                                                <Stack spacing={0.5}>
                                                    <Tooltip title="Cho phép đi muộn">
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            Muộn: +
                                                            {
                                                                shift.grace_period_minutes
                                                            }
                                                            p
                                                        </Typography>
                                                    </Tooltip>
                                                    <Tooltip title="Cho phép về sớm">
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            Sớm: -
                                                            {
                                                                shift.early_departure_minutes
                                                            }
                                                            p
                                                        </Typography>
                                                    </Tooltip>
                                                    <Tooltip title="Tính làm thêm giờ sau">
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            OT:{' '}
                                                            {shift.overtime_after_minutes /
                                                                60}
                                                            h
                                                        </Typography>
                                                    </Tooltip>
                                                </Stack>
                                            </TableCell>

                                            <TableCell align="center">
                                                <Tooltip title="Số nhân viên trong ca">
                                                    <Chip
                                                        icon={<People />}
                                                        label={
                                                            shift.employee_count ||
                                                            0
                                                        }
                                                        size="small"
                                                        color="primary"
                                                        variant="outlined"
                                                    />
                                                </Tooltip>
                                            </TableCell>

                                            <TableCell align="center">
                                                <Tooltip
                                                    title={
                                                        shift.is_active
                                                            ? 'Click để tạm dừng'
                                                            : 'Click để kích hoạt'
                                                    }
                                                >
                                                    <Chip
                                                        icon={
                                                            shift.is_active ? (
                                                                <CheckCircle />
                                                            ) : (
                                                                <Cancel />
                                                            )
                                                        }
                                                        label={
                                                            shift.is_active
                                                                ? 'Hoạt động'
                                                                : 'Tạm dừng'
                                                        }
                                                        size="small"
                                                        color={
                                                            shift.is_active
                                                                ? 'success'
                                                                : 'default'
                                                        }
                                                        onClick={() =>
                                                            handleToggleStatus(
                                                                shift
                                                            )
                                                        }
                                                        sx={{
                                                            cursor: 'pointer',
                                                        }}
                                                    />
                                                </Tooltip>
                                            </TableCell>

                                            <TableCell align="right">
                                                <Tooltip title="Phân công nhân viên">
                                                    <IconButton
                                                        size="small"
                                                        color="info"
                                                        onClick={() =>
                                                            navigate(
                                                                `/shifts/${shift.shift_id}/assign`
                                                            )
                                                        }
                                                    >
                                                        <Assignment />
                                                    </IconButton>
                                                </Tooltip>
                                                <Tooltip title="Chỉnh sửa">
                                                    <IconButton
                                                        size="small"
                                                        color="primary"
                                                        onClick={() =>
                                                            navigate(
                                                                `/shifts/${shift.shift_id}/edit`
                                                            )
                                                        }
                                                    >
                                                        <Edit />
                                                    </IconButton>
                                                </Tooltip>
                                                <Tooltip title="Xóa">
                                                    <IconButton
                                                        size="small"
                                                        color="error"
                                                        onClick={() =>
                                                            handleDeleteClick(
                                                                shift
                                                            )
                                                        }
                                                    >
                                                        <Delete />
                                                    </IconButton>
                                                </Tooltip>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>

                        {/* Pagination */}
                        {totalPages > 1 && (
                            <Box display="flex" justifyContent="center" p={2}>
                                <Pagination
                                    count={totalPages}
                                    page={page}
                                    onChange={handlePageChange}
                                    color="primary"
                                />
                            </Box>
                        )}
                    </>
                )}
            </Card>

            {/* Delete Confirmation Dialog */}
            <Dialog
                open={deleteDialog.open}
                onClose={() =>
                    setDeleteDialog({
                        open: false,
                        shiftId: null,
                        shiftName: '',
                    })
                }
            >
                <DialogTitle>Xác nhận xóa ca làm việc</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Bạn có chắc chắn muốn xóa ca làm việc "
                        {deleteDialog.shiftName}"? Hành động này không thể hoàn
                        tác.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={() =>
                            setDeleteDialog({
                                open: false,
                                shiftId: null,
                                shiftName: '',
                            })
                        }
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleDeleteConfirm}
                        color="error"
                        variant="contained"
                    >
                        Xóa
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar for notifications */}
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
        </Box>
    );
};
