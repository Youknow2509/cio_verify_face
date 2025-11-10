import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch,
    Chip,
    OutlinedInput,
    SelectChangeEvent,
    Divider,
    Alert,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';
import type { ShiftFormData } from '@face-attendance/types';

const DAYS_OF_WEEK = [
    { value: 1, label: 'Thứ 2' },
    { value: 2, label: 'Thứ 3' },
    { value: 3, label: 'Thứ 4' },
    { value: 4, label: 'Thứ 5' },
    { value: 5, label: 'Thứ 6' },
    { value: 6, label: 'Thứ 7' },
    { value: 7, label: 'Chủ nhật' },
];

export const ShiftFormPage: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams();

    const [formData, setFormData] = useState<ShiftFormData>({
        name: '',
        description: '',
        start_time: '08:00',
        end_time: '17:00',
        break_duration_minutes: 60,
        grace_period_minutes: 15,
        early_departure_minutes: 15,
        work_days: [1, 2, 3, 4, 5], // Monday to Friday by default
        is_flexible: false,
        overtime_after_minutes: 480, // 8 hours
        is_active: true,
    });

    const [errors, setErrors] = useState<
        Partial<Record<keyof ShiftFormData, string>>
    >({});

    useEffect(() => {
        // TODO: If editing, fetch shift data by id
        if (id) {
            // Fetch shift data and populate formData
            // Example: fetchShift(id).then(data => setFormData(data));
        }
    }, [id]);

    const handleChange =
        (field: keyof ShiftFormData) =>
        (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
            const value =
                e.target.type === 'checkbox'
                    ? (e.target as HTMLInputElement).checked
                    : e.target.value;

            setFormData({ ...formData, [field]: value });
            // Clear error when user starts typing
            if (errors[field]) {
                setErrors({ ...errors, [field]: undefined });
            }
        };

    const handleNumberChange =
        (field: keyof ShiftFormData) =>
        (e: React.ChangeEvent<HTMLInputElement>) => {
            const value = parseInt(e.target.value) || 0;
            setFormData({ ...formData, [field]: value });
            if (errors[field]) {
                setErrors({ ...errors, [field]: undefined });
            }
        };

    const handleSwitchChange =
        (field: keyof ShiftFormData) =>
        (e: React.ChangeEvent<HTMLInputElement>) => {
            setFormData({ ...formData, [field]: e.target.checked });
        };

    const handleWorkDaysChange = (event: SelectChangeEvent<number[]>) => {
        const value = event.target.value;
        setFormData({
            ...formData,
            work_days: typeof value === 'string' ? [] : value,
        });
    };

    const validateForm = (): boolean => {
        const newErrors: Partial<Record<keyof ShiftFormData, string>> = {};

        if (!formData.name.trim()) {
            newErrors.name = 'Tên ca làm việc là bắt buộc';
        }

        if (!formData.start_time) {
            newErrors.start_time = 'Giờ bắt đầu là bắt buộc';
        }

        if (!formData.end_time) {
            newErrors.end_time = 'Giờ kết thúc là bắt buộc';
        }

        if (
            formData.start_time &&
            formData.end_time &&
            formData.start_time >= formData.end_time
        ) {
            newErrors.end_time = 'Giờ kết thúc phải sau giờ bắt đầu';
        }

        if (formData.work_days.length === 0) {
            newErrors.work_days = 'Phải chọn ít nhất một ngày làm việc';
        }

        if (formData.grace_period_minutes < 0) {
            newErrors.grace_period_minutes =
                'Thời gian cho phép đi muộn không được âm';
        }

        if (formData.early_departure_minutes < 0) {
            newErrors.early_departure_minutes =
                'Thời gian cho phép về sớm không được âm';
        }

        if (formData.break_duration_minutes < 0) {
            newErrors.break_duration_minutes =
                'Thời gian nghỉ giải lao không được âm';
        }

        if (formData.overtime_after_minutes < 0) {
            newErrors.overtime_after_minutes =
                'Thời gian tính làm thêm giờ không được âm';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        // TODO: Call API to create/update shift
        console.log('Submitting shift data:', formData);

        // Navigate back to list page
        navigate('/shifts');
    };

    const calculateWorkHours = () => {
        if (!formData.start_time || !formData.end_time) return 0;

        const [startHour, startMin] = formData.start_time
            .split(':')
            .map(Number);
        const [endHour, endMin] = formData.end_time.split(':').map(Number);

        const totalMinutes =
            endHour * 60 +
            endMin -
            (startHour * 60 + startMin) -
            formData.break_duration_minutes;
        return Math.max(0, totalMinutes / 60);
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

            <Card>
                <CardContent>
                    <Typography variant="h5" fontWeight="bold" mb={3}>
                        {id ? 'Chỉnh sửa ca làm việc' : 'Thêm ca làm việc mới'}
                    </Typography>

                    <Box component="form" onSubmit={handleSubmit}>
                        <Grid container spacing={3}>
                            {/* Basic Information */}
                            <Grid item xs={12}>
                                <Typography
                                    variant="h6"
                                    fontWeight="medium"
                                    mb={2}
                                >
                                    Thông tin cơ bản
                                </Typography>
                            </Grid>

                            <Grid item xs={12} md={8}>
                                <TextField
                                    fullWidth
                                    label="Tên ca làm việc"
                                    required
                                    value={formData.name}
                                    onChange={handleChange('name')}
                                    error={!!errors.name}
                                    helperText={errors.name}
                                    placeholder="Ví dụ: Ca sáng, Ca chiều, Ca hành chính"
                                />
                            </Grid>

                            <Grid item xs={12} md={4}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={formData.is_active}
                                            onChange={handleSwitchChange(
                                                'is_active'
                                            )}
                                            color="primary"
                                        />
                                    }
                                    label="Kích hoạt"
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    label="Mô tả"
                                    multiline
                                    rows={3}
                                    value={formData.description}
                                    onChange={handleChange('description')}
                                    placeholder="Mô tả chi tiết về ca làm việc này"
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <Divider />
                            </Grid>

                            {/* Work Time */}
                            <Grid item xs={12}>
                                <Typography
                                    variant="h6"
                                    fontWeight="medium"
                                    mb={2}
                                >
                                    Thời gian làm việc
                                </Typography>
                            </Grid>

                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Giờ bắt đầu"
                                    type="time"
                                    required
                                    InputLabelProps={{ shrink: true }}
                                    value={formData.start_time}
                                    onChange={handleChange('start_time')}
                                    error={!!errors.start_time}
                                    helperText={errors.start_time}
                                />
                            </Grid>

                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Giờ kết thúc"
                                    type="time"
                                    required
                                    InputLabelProps={{ shrink: true }}
                                    value={formData.end_time}
                                    onChange={handleChange('end_time')}
                                    error={!!errors.end_time}
                                    helperText={errors.end_time}
                                />
                            </Grid>

                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Thời gian nghỉ giải lao (phút)"
                                    type="number"
                                    value={formData.break_duration_minutes}
                                    onChange={handleNumberChange(
                                        'break_duration_minutes'
                                    )}
                                    error={!!errors.break_duration_minutes}
                                    helperText={
                                        errors.break_duration_minutes ||
                                        'Thời gian nghỉ trong ca'
                                    }
                                    inputProps={{ min: 0 }}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <Alert severity="info" sx={{ mt: 1 }}>
                                    <strong>Tổng thời gian làm việc:</strong>{' '}
                                    {calculateWorkHours().toFixed(1)} giờ (
                                    {formData.start_time} - {formData.end_time},
                                    nghỉ {formData.break_duration_minutes} phút)
                                </Alert>
                            </Grid>

                            <Grid item xs={12}>
                                <FormControl
                                    fullWidth
                                    error={!!errors.work_days}
                                >
                                    <InputLabel>
                                        Ngày làm việc trong tuần
                                    </InputLabel>
                                    <Select
                                        multiple
                                        value={formData.work_days}
                                        onChange={handleWorkDaysChange}
                                        input={
                                            <OutlinedInput label="Ngày làm việc trong tuần" />
                                        }
                                        renderValue={(selected) => (
                                            <Box
                                                sx={{
                                                    display: 'flex',
                                                    flexWrap: 'wrap',
                                                    gap: 0.5,
                                                }}
                                            >
                                                {selected.map((value) => (
                                                    <Chip
                                                        key={value}
                                                        label={
                                                            DAYS_OF_WEEK.find(
                                                                (d) =>
                                                                    d.value ===
                                                                    value
                                                            )?.label
                                                        }
                                                        size="small"
                                                    />
                                                ))}
                                            </Box>
                                        )}
                                    >
                                        {DAYS_OF_WEEK.map((day) => (
                                            <MenuItem
                                                key={day.value}
                                                value={day.value}
                                            >
                                                {day.label}
                                            </MenuItem>
                                        ))}
                                    </Select>
                                    {errors.work_days && (
                                        <Typography
                                            variant="caption"
                                            color="error"
                                            sx={{ mt: 0.5, ml: 1.5 }}
                                        >
                                            {errors.work_days}
                                        </Typography>
                                    )}
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <Divider />
                            </Grid>

                            {/* Attendance Policy */}
                            <Grid item xs={12}>
                                <Typography
                                    variant="h6"
                                    fontWeight="medium"
                                    mb={2}
                                >
                                    Chính sách chấm công
                                </Typography>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Thời gian cho phép đi muộn (phút)"
                                    type="number"
                                    value={formData.grace_period_minutes}
                                    onChange={handleNumberChange(
                                        'grace_period_minutes'
                                    )}
                                    error={!!errors.grace_period_minutes}
                                    helperText={
                                        errors.grace_period_minutes ||
                                        'Vượt quá sẽ bị tính muộn'
                                    }
                                    inputProps={{ min: 0 }}
                                />
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Thời gian cho phép về sớm (phút)"
                                    type="number"
                                    value={formData.early_departure_minutes}
                                    onChange={handleNumberChange(
                                        'early_departure_minutes'
                                    )}
                                    error={!!errors.early_departure_minutes}
                                    helperText={
                                        errors.early_departure_minutes ||
                                        'Vượt quá sẽ bị tính về sớm'
                                    }
                                    inputProps={{ min: 0 }}
                                />
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Tính làm thêm giờ sau (phút)"
                                    type="number"
                                    value={formData.overtime_after_minutes}
                                    onChange={handleNumberChange(
                                        'overtime_after_minutes'
                                    )}
                                    error={!!errors.overtime_after_minutes}
                                    helperText={
                                        errors.overtime_after_minutes ||
                                        `Mặc định: 480 phút (8 giờ). Tương đương: ${(
                                            formData.overtime_after_minutes / 60
                                        ).toFixed(1)} giờ`
                                    }
                                    inputProps={{ min: 0 }}
                                />
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={formData.is_flexible}
                                            onChange={handleSwitchChange(
                                                'is_flexible'
                                            )}
                                            color="primary"
                                        />
                                    }
                                    label="Ca linh hoạt"
                                />
                                <Typography
                                    variant="caption"
                                    color="text.secondary"
                                    display="block"
                                    sx={{ ml: 4 }}
                                >
                                    Cho phép nhân viên linh hoạt giờ vào/ra
                                    trong khung giờ quy định
                                </Typography>
                            </Grid>
                        </Grid>

                        <Box mt={4} display="flex" gap={2}>
                            <Button
                                type="submit"
                                variant="contained"
                                startIcon={<Save />}
                                size="large"
                            >
                                {id ? 'Cập nhật' : 'Tạo mới'}
                            </Button>
                            <Button
                                variant="outlined"
                                onClick={() => navigate('/shifts')}
                                size="large"
                            >
                                Hủy
                            </Button>
                        </Box>
                    </Box>
                </CardContent>
            </Card>
        </Box>
    );
};
