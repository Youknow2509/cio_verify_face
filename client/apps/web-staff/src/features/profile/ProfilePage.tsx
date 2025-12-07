import { useState, useEffect } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    TextField,
    Button,
    Avatar,
    Divider,
    Alert,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    IconButton,
    LinearProgress,
    Chip,
    Tooltip,
} from '@mui/material';
import {
    Person,
    Email,
    Phone,
    Business,
    Badge,
    Lock,
    Edit,
    Save,
    Close,
    Face,
    Cake,
    Home,
    Info,
} from '@mui/icons-material';
import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

export const ProfilePage: React.FC = () => {
    const { user, setUser, accessToken } = useAuthStore();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [editing, setEditing] = useState(false);
    const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const [formData, setFormData] = useState({
        phone: '',
        address: '',
    });

    const [passwordData, setPasswordData] = useState({
        current_password: '',
        new_password: '',
        confirm_password: '',
    });

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await axios.get('/api/v1/auth/me', {
                    headers: { Authorization: `Bearer ${accessToken}` },
                });
                setUser(response.data.data);
            } catch (err) {
                console.error('Failed to fetch profile:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchProfile();
    }, [accessToken, setUser]);

    useEffect(() => {
        if (user) {
            setFormData({
                phone: user.phone || '',
                address: user.address || '',
            });
        }
    }, [user]);

    const handleSave = async () => {
        setSaving(true);
        setError('');
        try {
            await axios.put(
                `/api/v1/users/${user?.id}`,
                formData,
                { headers: { Authorization: `Bearer ${accessToken}` } }
            );
            setSuccess('Cập nhật thông tin thành công!');
            setEditing(false);
            // Refresh user data
            const response = await axios.get('/api/v1/auth/me', {
                headers: { Authorization: `Bearer ${accessToken}` },
            });
            setUser(response.data.data);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Cập nhật thất bại');
        } finally {
            setSaving(false);
        }
    };

    const handlePasswordChange = async () => {
        if (passwordData.new_password !== passwordData.confirm_password) {
            setError('Mật khẩu xác nhận không khớp');
            return;
        }

        setSaving(true);
        setError('');
        try {
            await axios.post(
                '/api/v1/password/reset',
                {
                    current_password: passwordData.current_password,
                    new_password: passwordData.new_password,
                },
                { headers: { Authorization: `Bearer ${accessToken}` } }
            );
            setSuccess('Đổi mật khẩu thành công! Vui lòng đăng nhập lại.');
            setPasswordDialogOpen(false);
            setPasswordData({ current_password: '', new_password: '', confirm_password: '' });
        } catch (err: any) {
            setError(err.response?.data?.message || 'Đổi mật khẩu thất bại');
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return <LinearProgress />;
    }

    const workInfo = [
        { icon: <Badge />, label: 'Mã nhân viên', value: user?.employee_id || 'N/A' },
        { icon: <Email />, label: 'Email công ty', value: user?.email || 'N/A' },
        { icon: <Business />, label: 'Công ty', value: user?.company_name || 'N/A' },
        { icon: <Person />, label: 'Phòng ban', value: user?.department || 'N/A' },
        { icon: <Badge />, label: 'Chức vụ', value: user?.position || 'N/A' },
    ];

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Typography variant="h4" fontWeight="700" mb={3}>
                Hồ sơ Nhân viên
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess('')}>
                    {success}
                </Alert>
            )}

            <Grid container spacing={3}>
                {/* Left Column: Avatar & Actions */}
                <Grid item xs={12} md={4}>
                    <Card sx={{ textAlign: 'center', p: 3, height: '100%' }}>
                        <Box sx={{ position: 'relative', display: 'inline-block' }}>
                            <Avatar
                                src={user?.avatar}
                                sx={{
                                    width: 140,
                                    height: 140,
                                    mx: 'auto',
                                    mb: 2,
                                    fontSize: '3.5rem',
                                    bgcolor: 'primary.main',
                                    border: '4px solid',
                                    borderColor: 'primary.light',
                                    boxShadow: '0 8px 16px rgba(0,0,0,0.1)',
                                }}
                            >
                                {user?.full_name?.charAt(0) || 'U'}
                            </Avatar>
                            <Tooltip title="Cập nhật ảnh đại diện">
                                <IconButton
                                    sx={{
                                        position: 'absolute',
                                        bottom: 16,
                                        right: 0,
                                        bgcolor: 'background.paper',
                                        boxShadow: 2,
                                        '&:hover': { bgcolor: 'action.hover' },
                                    }}
                                    onClick={() => (window.location.href = '/face-update')}
                                >
                                    <Face color="primary" />
                                </IconButton>
                            </Tooltip>
                        </Box>

                        <Typography variant="h5" fontWeight="700" mb={0.5}>
                            {user?.full_name || 'Nhân viên'}
                        </Typography>
                        <Typography variant="body1" color="text.secondary" mb={2}>
                            {user?.position || 'Chức vụ'}
                        </Typography>

                        <Box sx={{ display: 'flex', justifyContent: 'center', gap: 1, flexWrap: 'wrap', mb: 3 }}>
                            <Chip
                                label={user?.department || 'Phòng ban'}
                                color="primary"
                                variant="outlined"
                                size="small"
                            />
                            <Chip
                                label={user?.status === 'active' ? 'Đang hoạt động' : 'Tạm khóa'}
                                color={user?.status === 'active' ? 'success' : 'error'}
                                size="small"
                            />
                        </Box>

                        <Divider sx={{ my: 3 }} />

                        <Button
                            variant="outlined"
                            startIcon={<Lock />}
                            fullWidth
                            onClick={() => setPasswordDialogOpen(true)}
                            color="inherit"
                        >
                            Đổi mật khẩu
                        </Button>
                    </Card>
                </Grid>

                {/* Right Column: Info Details */}
                <Grid item xs={12} md={8}>
                    {/* Personal Information */}
                    <Card sx={{ mb: 3 }}>
                        <CardContent>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <Person color="primary" />
                                    <Typography variant="h6" fontWeight="600">
                                        Thông tin cá nhân
                                    </Typography>
                                </Box>
                                {editing ? (
                                    <Box sx={{ display: 'flex', gap: 1 }}>
                                        <Button
                                            variant="outlined"
                                            startIcon={<Close />}
                                            onClick={() => setEditing(false)}
                                            size="small"
                                        >
                                            Hủy
                                        </Button>
                                        <Button
                                            variant="contained"
                                            startIcon={<Save />}
                                            onClick={handleSave}
                                            disabled={saving}
                                            size="small"
                                        >
                                            Lưu thay đổi
                                        </Button>
                                    </Box>
                                ) : (
                                    <Button
                                        variant="text"
                                        startIcon={<Edit />}
                                        onClick={() => setEditing(true)}
                                        size="small"
                                    >
                                        Chỉnh sửa
                                    </Button>
                                )}
                            </Box>

                            <Grid container spacing={3}>
                                {/* Read-only Personal Fields */}
                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        label="Họ và tên"
                                        value={user?.full_name || ''}
                                        disabled
                                        InputProps={{
                                            startAdornment: <Person sx={{ mr: 1, color: 'text.secondary' }} />,
                                            endAdornment: (
                                                <Tooltip title="Liên hệ Admin để thay đổi">
                                                    <Info fontSize="small" color="disabled" sx={{ cursor: 'help' }} />
                                                </Tooltip>
                                            ),
                                        }}
                                    />
                                </Grid>
                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        label="Ngày sinh"
                                        value={user?.date_of_birth ? new Date(user.date_of_birth).toLocaleDateString('vi-VN') : 'Chưa cập nhật'}
                                        disabled
                                        InputProps={{
                                            startAdornment: <Cake sx={{ mr: 1, color: 'text.secondary' }} />,
                                        }}
                                    />
                                </Grid>

                                {/* Editable Fields */}
                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        label="Số điện thoại"
                                        value={editing ? formData.phone : (user?.phone || 'Chưa cập nhật')}
                                        onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                        disabled={!editing}
                                        InputProps={{
                                            startAdornment: <Phone sx={{ mr: 1, color: editing ? 'primary.main' : 'text.secondary' }} />,
                                        }}
                                        sx={{
                                            '& .MuiInputBase-root': {
                                                bgcolor: editing ? 'background.paper' : 'action.hover',
                                            }
                                        }}
                                    />
                                </Grid>
                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        label="Địa chỉ liên hệ"
                                        value={editing ? formData.address : (user?.address || 'Chưa cập nhật')}
                                        onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                                        disabled={!editing}
                                        InputProps={{
                                            startAdornment: <Home sx={{ mr: 1, color: editing ? 'primary.main' : 'text.secondary' }} />,
                                        }}
                                        sx={{
                                            '& .MuiInputBase-root': {
                                                bgcolor: editing ? 'background.paper' : 'action.hover',
                                            }
                                        }}
                                    />
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>

                    {/* Work Information (Read-only) */}
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Business color="primary" />
                                <Typography variant="h6" fontWeight="600">
                                    Thông tin công việc
                                </Typography>
                            </Box>
                            <Grid container spacing={3}>
                                {workInfo.map((item, index) => (
                                    <Grid item xs={12} sm={6} key={index}>
                                        <TextField
                                            fullWidth
                                            label={item.label}
                                            value={item.value}
                                            disabled
                                            InputProps={{
                                                startAdornment: <Box sx={{ mr: 1, color: 'text.secondary' }}>{item.icon}</Box>,
                                                readOnly: true,
                                            }}
                                            variant="outlined"
                                            sx={{
                                                '& .MuiInputBase-root': {
                                                    bgcolor: 'action.hover',
                                                }
                                            }}
                                        />
                                    </Grid>
                                ))}
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>

            {/* Password Dialog - Kept same logic, just minor UI tweaks */}
            <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)} maxWidth="xs" fullWidth>
                <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    Đổi mật khẩu
                    <IconButton size="small" onClick={() => setPasswordDialogOpen(false)}>
                        <Close />
                    </IconButton>
                </DialogTitle>
                <DialogContent>
                    <Box sx={{ pt: 1, display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                        <TextField
                            fullWidth
                            type="password"
                            label="Mật khẩu hiện tại"
                            placeholder="Nhập mật khẩu hiện tại"
                            value={passwordData.current_password}
                            onChange={(e) => setPasswordData({ ...passwordData, current_password: e.target.value })}
                        />
                        <TextField
                            fullWidth
                            type="password"
                            label="Mật khẩu mới"
                            placeholder="Ít nhất 6 ký tự"
                            value={passwordData.new_password}
                            onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })}
                        />
                        <TextField
                            fullWidth
                            type="password"
                            label="Xác nhận mật khẩu"
                            placeholder="Nhập lại mật khẩu mới"
                            value={passwordData.confirm_password}
                            onChange={(e) => setPasswordData({ ...passwordData, confirm_password: e.target.value })}
                        />
                    </Box>
                </DialogContent>
                <DialogActions sx={{ p: 2.5 }}>
                    <Button onClick={() => setPasswordDialogOpen(false)} color="inherit">Hủy</Button>
                    <Button variant="contained" onClick={handlePasswordChange} disabled={saving}>
                        Đổi mật khẩu
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
