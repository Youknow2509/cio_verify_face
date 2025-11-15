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
    MenuItem,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils/src/api';

export const EmployeeFormPage: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams();
    const isEdit = !!id;

    const [formData, setFormData] = useState({
        full_name: '',
        employee_code: '',
        email: '',
        phone: '',
        department: '',
        position: '',
        hire_date: '',
        password: '',
        avatar_url: '',
        salary: '',
        role: '',
        company_id: '',
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Get data and populate form if editing
    useEffect(() => {
        setLoading(true);
        setError(null);
        if (isEdit && id) {
            fetchEmployee(id);
        } else {
            // when adding, prefill company_id from token if available
            const accessToken = localStorage.getItem('access_token');
            if (accessToken) {
                try {
                    const payload = JSON.parse(atob(accessToken.split('.')[1]));
                    if (payload?.company_id) {
                        setFormData((prev) => ({
                            ...prev,
                            company_id: prev.company_id || payload.company_id,
                        }));
                    }
                } catch {
                    // ignore
                }
            }
        }
        setLoading(false);
    }, [isEdit, id]);

    // fetch employee data
    const fetchEmployee = async (employeeId: string) => {
        const response = await apiClient.get(`/api/v1/users/${employeeId}`);
        if (response.status === 200) {
            const data = response.data.data;
            console.log('Fetched employee data:', data);
            setFormData({
                full_name: data.full_name || '',
                employee_code: data.employee_code || '',
                email: data.email || '',
                phone: data.phone || '',
                department: data.department || '',
                position: data.position || '',
                hire_date: formatDate(data.hire_date) || '',
                password: '',
                avatar_url: data.avatar_url || '',
                salary: data.salary ?? '',
                role: data.role ?? '',
                company_id: data.company_id || '',
            });
        }
    };

    // convert "2025-11-11T07:14:51.930Z" to yyyy-mm-dd
    const formatDate = (isoDate: string) => {
        const date = new Date(isoDate);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    };

    const handleChange =
        (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
            setFormData({ ...formData, [field]: e.target.value });
        };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        // basic validation
        if (
            !formData.full_name ||
            !formData.email ||
            (!isEdit && !formData.password)
        ) {
            setError('Vui lòng nhập Họ tên, Email và Mật khẩu (khi thêm mới).');
            return;
        }

        setLoading(true);
        try {
            if (isEdit && id) {
                // existing behavior: update-base
                const mapData = {
                    user_id: id,
                    user_fullname: formData.full_name,
                    user_phone: formData.phone,
                    user_email: formData.email,
                    user_department: formData.department,
                    user_data_join_company: formData.hire_date,
                    user_position: formData.position,
                };
                const response = await apiClient.post(
                    '/api/v1/users/update-base',
                    mapData
                );
                if (response.status === 201 || response.status === 200) {
                    navigate('/employees');
                    return;
                }
                throw new Error('Lưu nhân viên thất bại');
            } else {
                // create new user
                const payload: any = {
                    email: formData.email,
                    password: formData.password,
                    full_name: formData.full_name,
                    phone: formData.phone,
                    avatar_url: formData.avatar_url,
                    company_id: formData.company_id,
                    employee_code: formData.employee_code,
                    department: formData.department,
                    position: formData.position,
                    hire_date: formData.hire_date,
                    salary: formData.salary,
                    role: formData.role,
                };
                const response = await apiClient.post('/api/v1/users', payload);
                if (response.status === 201 || response.status === 200) {
                    navigate('/employees');
                    return;
                }
                throw new Error('Tạo nhân viên thất bại');
            }
        } catch (err: any) {
            console.error('Failed to save/create employee:', err);
            setError(
                err.response?.data?.error ||
                    err.message ||
                    'Lỗi khi lưu nhân viên'
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box>
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/employees')}
                sx={{ mb: 2 }}
            >
                Quay lại
            </Button>
            {error && (
                <Typography color="error" mb={2}>
                    {error}
                </Typography>
            )}
            {loading && (
                <Typography color="textSecondary" mb={2}>
                    Đang tải...
                </Typography>
            )}
            <Card>
                <CardContent>
                    <Typography variant="h5" fontWeight="bold" mb={3}>
                        {isEdit ? 'Chỉnh sửa nhân viên' : 'Thêm nhân viên mới'}
                    </Typography>
                    <Box component="form" onSubmit={handleSubmit}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Họ và tên"
                                    required
                                    value={formData.full_name}
                                    onChange={handleChange('full_name')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Mã nhân viên"
                                    required
                                    value={formData.employee_code}
                                    InputProps={{ disabled: true }}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Email"
                                    type="email"
                                    required
                                    value={formData.email}
                                    onChange={handleChange('email')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Số điện thoại"
                                    value={formData.phone}
                                    onChange={handleChange('phone')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Phòng ban"
                                    value={formData.department}
                                    onChange={handleChange('department')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Chức vụ"
                                    value={formData.position}
                                    onChange={handleChange('position')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Ngày vào làm"
                                    type="date"
                                    InputLabelProps={{ shrink: true }}
                                    value={formData.hire_date || new Date().toISOString().split('T')[0]}
                                    onChange={handleChange('hire_date')}
                                />
                            </Grid>
                            {!isEdit && (
                                <Grid item xs={12} md={6}>
                                    <TextField
                                        fullWidth
                                        label="Mật khẩu"
                                        type="password"
                                        required
                                        value={formData.password}
                                        onChange={handleChange('password')}
                                    />
                                </Grid>
                            )}

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Avatar URL"
                                    value={formData.avatar_url}
                                    onChange={handleChange('avatar_url')}
                                />
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Salary"
                                    value={formData.salary}
                                    onChange={handleChange('salary')}
                                />
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    select
                                    fullWidth
                                    label="Role"
                                    value={formData.role}
                                    onChange={handleChange('role')}
                                >
                                    <MenuItem value={''}>Chọn vai trò</MenuItem>
                                    <MenuItem value={1}>Manager</MenuItem>
                                    <MenuItem value={2}>Employee</MenuItem>
                                </TextField>
                            </Grid>

                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Company ID"
                                    value={formData.company_id}
                                    InputProps={{ disabled: true }}
                                />
                            </Grid>
                        </Grid>
                        <Box mt={3} display="flex" gap={2}>
                            <Button
                                type="submit"
                                variant="contained"
                                startIcon={<Save />}
                                onClick={handleSubmit}
                            >
                                Lưu
                            </Button>
                            <Button
                                variant="outlined"
                                onClick={() => navigate('/employees')}
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
