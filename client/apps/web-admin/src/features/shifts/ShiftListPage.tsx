import { useState } from 'react';
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
} from '@mui/icons-material';
import type { Shift } from '@face-attendance/types';

const DAYS_OF_WEEK_SHORT = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];

export const ShiftListPage: React.FC = () => {
    const navigate = useNavigate();
    const [shifts, setShifts] = useState<Shift[]>([
        {
            shift_id: '1',
            company_id: '1',
            name: 'Ca hành chính',
            description: 'Ca làm việc hành chính tiêu chuẩn',
            start_time: '08:00',
            end_time: '17:00',
            break_duration_minutes: 60,
            grace_period_minutes: 15,
            early_departure_minutes: 15,
            work_days: [1, 2, 3, 4, 5],
            is_flexible: false,
            overtime_after_minutes: 480,
            is_active: true,
            employee_count: 50,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        },
        {
            shift_id: '2',
            company_id: '1',
            name: 'Ca sáng',
            description: 'Ca làm việc sáng',
            start_time: '06:00',
            end_time: '14:00',
            break_duration_minutes: 30,
            grace_period_minutes: 10,
            early_departure_minutes: 10,
            work_days: [1, 2, 3, 4, 5, 6],
            is_flexible: true,
            overtime_after_minutes: 480,
            is_active: true,
            employee_count: 25,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        },
        {
            shift_id: '3',
            company_id: '1',
            name: 'Ca chiều',
            description: 'Ca làm việc chiều',
            start_time: '14:00',
            end_time: '22:00',
            break_duration_minutes: 45,
            grace_period_minutes: 15,
            early_departure_minutes: 15,
            work_days: [1, 2, 3, 4, 5],
            is_flexible: false,
            overtime_after_minutes: 480,
            is_active: false,
            employee_count: 15,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        },
    ]);

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

    const handleDelete = (shiftId: string) => {
        // TODO: Show confirmation dialog and call API
        console.log('Delete shift:', shiftId);
        setShifts(shifts.filter((s) => s.shift_id !== shiftId));
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

            <Card>
                <TableContainer>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>Tên ca</TableCell>
                                <TableCell>Giờ làm việc</TableCell>
                                <TableCell>Ngày làm việc</TableCell>
                                <TableCell align="center">Chính sách</TableCell>
                                <TableCell align="center">Nhân viên</TableCell>
                                <TableCell align="center">Trạng thái</TableCell>
                                <TableCell align="right">Thao tác</TableCell>
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
                                                {shift.break_duration_minutes}p
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
                                                    {shift.grace_period_minutes}
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
                                                    shift.employee_count || 0
                                                }
                                                size="small"
                                                color="primary"
                                                variant="outlined"
                                            />
                                        </Tooltip>
                                    </TableCell>

                                    <TableCell align="center">
                                        {shift.is_active ? (
                                            <Chip
                                                icon={<CheckCircle />}
                                                label="Hoạt động"
                                                size="small"
                                                color="success"
                                            />
                                        ) : (
                                            <Chip
                                                icon={<Cancel />}
                                                label="Tạm dừng"
                                                size="small"
                                                color="default"
                                            />
                                        )}
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
                                                    handleDelete(shift.shift_id)
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
            </Card>
        </Box>
    );
};
