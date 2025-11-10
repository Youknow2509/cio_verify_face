import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    Autocomplete,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
    Chip,
    Alert,
} from '@mui/material';
import { Save, ArrowBack, Delete, Add } from '@mui/icons-material';
import type { Shift, EmployeeShift } from '@face-attendance/types';

interface Employee {
    id: string;
    name: string;
    employee_code: string;
    current_shift?: string;
}

export const ShiftAssignmentPage: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams(); // shift_id nếu đến từ /shifts/:id/assign

    // Mock data - should be fetched from API
    const [shifts] = useState<Shift[]>([
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
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        },
    ]);

    const [employees] = useState<Employee[]>([
        {
            id: '1',
            name: 'Nguyễn Văn A',
            employee_code: 'NV001',
            current_shift: 'Ca hành chính',
        },
        { id: '2', name: 'Trần Thị B', employee_code: 'NV002' },
        { id: '3', name: 'Lê Văn C', employee_code: 'NV003' },
        { id: '4', name: 'Phạm Thị D', employee_code: 'NV004' },
        {
            id: '5',
            name: 'Hoàng Văn E',
            employee_code: 'NV005',
            current_shift: 'Ca sáng',
        },
        { id: '6', name: 'Đặng Thị F', employee_code: 'NV006' },
    ]);

    const [assignments, setAssignments] = useState<EmployeeShift[]>([]);
    const [formData, setFormData] = useState({
        employee_ids: [] as string[], // ✅ Thay đổi: array để chọn nhiều
        shift_id: id || '', // Pre-fill nếu có shift_id trong URL
        effective_from: new Date().toISOString().split('T')[0],
        effective_to: '',
    });

    // Tìm shift được chọn (nếu có)
    const selectedShift = id ? shifts.find((s) => s.shift_id === id) : null;

    const handleAddAssignment = () => {
        // ✅ Validate: phải chọn ít nhất 1 nhân viên
        if (formData.employee_ids.length === 0 || !formData.shift_id) {
            alert('Vui lòng chọn ít nhất một nhân viên và ca làm việc');
            return;
        }

        const shift = shifts.find((s) => s.shift_id === formData.shift_id);
        if (!shift) return;

        // ✅ Tạo assignment cho TẤT CẢ nhân viên được chọn
        const newAssignments: EmployeeShift[] = formData.employee_ids.map(
            (employeeId) => {
                const employee = employees.find((e) => e.id === employeeId);

                return {
                    employee_shift_id: `temp-${Date.now()}-${employeeId}`,
                    employee_id: employeeId,
                    shift_id: formData.shift_id,
                    effective_from: formData.effective_from,
                    effective_to: formData.effective_to || undefined,
                    is_active: true,
                    created_at: new Date().toISOString(),
                    employee_name: employee?.name || '',
                    shift_name: shift.name,
                };
            }
        );

        // ✅ Thêm tất cả assignments vào danh sách
        setAssignments([...assignments, ...newAssignments]);

        // Reset form
        setFormData({
            ...formData,
            employee_ids: [], // ✅ Reset về array rỗng
            effective_to: '',
        });
    };

    const handleRemoveAssignment = (assignmentId: string) => {
        setAssignments(
            assignments.filter((a) => a.employee_shift_id !== assignmentId)
        );
    };

    const handleSubmit = () => {
        // TODO: Call API to save all assignments
        console.log('Saving assignments:', assignments);
        navigate('/shifts');
    };

    return (
        <Box>
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/shifts')}
                sx={{ mb: 2 }}
            >
                Quay lại
            </Button>

            <Grid container spacing={3}>
                {/* Assignment Form */}
                <Grid item xs={12} md={6}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="bold" mb={2}>
                                Phân công ca làm việc
                            </Typography>

                            <Grid container spacing={2}>
                                <Grid item xs={12}>
                                    <Box
                                        display="flex"
                                        justifyContent="space-between"
                                        alignItems="center"
                                        mb={1}
                                    >
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                        >
                                            Chọn nhân viên để phân công
                                        </Typography>
                                        <Box display="flex" gap={1}>
                                            <Button
                                                size="small"
                                                onClick={() =>
                                                    setFormData({
                                                        ...formData,
                                                        employee_ids:
                                                            employees.map(
                                                                (e) => e.id
                                                            ),
                                                    })
                                                }
                                            >
                                                Chọn tất cả
                                            </Button>
                                            <Button
                                                size="small"
                                                onClick={() =>
                                                    setFormData({
                                                        ...formData,
                                                        employee_ids: [],
                                                    })
                                                }
                                                disabled={
                                                    formData.employee_ids
                                                        .length === 0
                                                }
                                            >
                                                Bỏ chọn
                                            </Button>
                                        </Box>
                                    </Box>
                                    <Autocomplete
                                        multiple // ✅ Enable multi-select
                                        options={employees}
                                        disableCloseOnSelect // ✅ Không đóng khi chọn
                                        getOptionLabel={(option) =>
                                            `${option.employee_code} - ${option.name}`
                                        }
                                        value={employees.filter((e) =>
                                            formData.employee_ids.includes(e.id)
                                        )} // ✅ Lấy tất cả nhân viên được chọn
                                        onChange={(_, newValue) => {
                                            setFormData({
                                                ...formData,
                                                employee_ids: newValue.map(
                                                    (emp) => emp.id
                                                ), // ✅ Lưu array ids
                                            });
                                        }}
                                        renderInput={(params) => (
                                            <TextField
                                                {...params}
                                                label="Chọn nhân viên (có thể chọn nhiều)"
                                                required
                                                placeholder={
                                                    formData.employee_ids
                                                        .length === 0
                                                        ? 'Tìm và chọn nhân viên...'
                                                        : ''
                                                }
                                            />
                                        )}
                                        renderOption={(props, option) => (
                                            <li {...props}>
                                                <Box>
                                                    <Typography variant="body2">
                                                        {option.employee_code} -{' '}
                                                        {option.name}
                                                    </Typography>
                                                    {option.current_shift && (
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            Ca hiện tại:{' '}
                                                            {
                                                                option.current_shift
                                                            }
                                                        </Typography>
                                                    )}
                                                </Box>
                                            </li>
                                        )}
                                    />
                                </Grid>

                                {/* ✅ Hiển thị số lượng nhân viên đã chọn */}
                                {formData.employee_ids.length > 0 && (
                                    <Grid item xs={12}>
                                        <Alert severity="success" icon={false}>
                                            <Typography variant="body2">
                                                <strong>
                                                    Đã chọn{' '}
                                                    {
                                                        formData.employee_ids
                                                            .length
                                                    }{' '}
                                                    nhân viên:
                                                </strong>{' '}
                                                {employees
                                                    .filter((e) =>
                                                        formData.employee_ids.includes(
                                                            e.id
                                                        )
                                                    )
                                                    .map((e) => e.name)
                                                    .join(', ')}
                                            </Typography>
                                        </Alert>
                                    </Grid>
                                )}

                                <Grid item xs={12}>
                                    {selectedShift && (
                                        <Alert severity="info" sx={{ mb: 2 }}>
                                            <strong>Phân công cho ca:</strong>{' '}
                                            {selectedShift.name}
                                            <br />
                                            <Typography variant="caption">
                                                {selectedShift.start_time} -{' '}
                                                {selectedShift.end_time}
                                            </Typography>
                                        </Alert>
                                    )}
                                    <FormControl fullWidth required>
                                        <InputLabel>Ca làm việc</InputLabel>
                                        <Select
                                            value={formData.shift_id}
                                            onChange={(e) =>
                                                setFormData({
                                                    ...formData,
                                                    shift_id: e.target.value,
                                                })
                                            }
                                            label="Ca làm việc"
                                            disabled={!!id} // Disable nếu đã chọn shift từ URL
                                        >
                                            {shifts.map((shift) => (
                                                <MenuItem
                                                    key={shift.shift_id}
                                                    value={shift.shift_id}
                                                >
                                                    <Box>
                                                        <Typography variant="body2">
                                                            {shift.name}
                                                        </Typography>
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            {shift.start_time} -{' '}
                                                            {shift.end_time}
                                                        </Typography>
                                                    </Box>
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                </Grid>

                                <Grid item xs={12} md={6}>
                                    <TextField
                                        fullWidth
                                        label="Có hiệu lực từ"
                                        type="date"
                                        required
                                        InputLabelProps={{ shrink: true }}
                                        value={formData.effective_from}
                                        onChange={(e) =>
                                            setFormData({
                                                ...formData,
                                                effective_from: e.target.value,
                                            })
                                        }
                                    />
                                </Grid>

                                <Grid item xs={12} md={6}>
                                    <TextField
                                        fullWidth
                                        label="Có hiệu lực đến"
                                        type="date"
                                        InputLabelProps={{ shrink: true }}
                                        value={formData.effective_to}
                                        onChange={(e) =>
                                            setFormData({
                                                ...formData,
                                                effective_to: e.target.value,
                                            })
                                        }
                                        helperText="Để trống nếu không có ngày kết thúc"
                                    />
                                </Grid>

                                <Grid item xs={12}>
                                    <Button
                                        fullWidth
                                        variant="contained"
                                        startIcon={<Add />}
                                        onClick={handleAddAssignment}
                                        disabled={
                                            formData.employee_ids.length ===
                                                0 || !formData.shift_id
                                        }
                                    >
                                        {formData.employee_ids.length > 0
                                            ? `Thêm phân công cho ${formData.employee_ids.length} nhân viên`
                                            : 'Thêm phân công'}
                                    </Button>
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Assignments List */}
                <Grid item xs={12} md={6}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="bold" mb={3}>
                                Danh sách phân công ({assignments.length})
                            </Typography>

                            {assignments.length === 0 ? (
                                <Alert severity="info">
                                    Chưa có phân công nào. Vui lòng thêm phân
                                    công ở bên trái.
                                </Alert>
                            ) : (
                                <TableContainer>
                                    <Table size="small">
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Nhân viên</TableCell>
                                                <TableCell>Ca</TableCell>
                                                <TableCell>Thời gian</TableCell>
                                                <TableCell align="center">
                                                    Thao tác
                                                </TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {assignments.map((assignment) => (
                                                <TableRow
                                                    key={
                                                        assignment.employee_shift_id
                                                    }
                                                >
                                                    <TableCell>
                                                        <Typography
                                                            variant="body2"
                                                            fontWeight="medium"
                                                        >
                                                            {
                                                                assignment.employee_name
                                                            }
                                                        </Typography>
                                                    </TableCell>
                                                    <TableCell>
                                                        <Chip
                                                            label={
                                                                assignment.shift_name
                                                            }
                                                            size="small"
                                                        />
                                                    </TableCell>
                                                    <TableCell>
                                                        <Typography
                                                            variant="caption"
                                                            display="block"
                                                        >
                                                            Từ:{' '}
                                                            {
                                                                assignment.effective_from
                                                            }
                                                        </Typography>
                                                        {assignment.effective_to && (
                                                            <Typography
                                                                variant="caption"
                                                                display="block"
                                                            >
                                                                Đến:{' '}
                                                                {
                                                                    assignment.effective_to
                                                                }
                                                            </Typography>
                                                        )}
                                                    </TableCell>
                                                    <TableCell align="center">
                                                        <IconButton
                                                            size="small"
                                                            color="error"
                                                            onClick={() =>
                                                                handleRemoveAssignment(
                                                                    assignment.employee_shift_id
                                                                )
                                                            }
                                                        >
                                                            <Delete />
                                                        </IconButton>
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            )}

                            {assignments.length > 0 && (
                                <Box mt={3} display="flex" gap={2}>
                                    <Button
                                        fullWidth
                                        variant="contained"
                                        color="primary"
                                        startIcon={<Save />}
                                        onClick={handleSubmit}
                                    >
                                        Lưu tất cả phân công
                                    </Button>
                                </Box>
                            )}
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
