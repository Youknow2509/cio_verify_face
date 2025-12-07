import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    TextField,
    Button,
    Typography,
    Alert,
    CircularProgress,
    Container,
    Paper,
    FormControlLabel,
    Checkbox,
} from '@mui/material';
import {
    Lock as LockIcon,
    MailOutline as EmailIcon,
    Visibility,
    VisibilityOff,
    ArrowForward,
    Security,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';
import { apiClient } from '@face-attendance/utils';
import type { LoginRequest, LoginResponse } from '@face-attendance/types';

export const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const { setAuth } = useAuthStore();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [rememberMe, setRememberMe] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            console.log('Submitting login form with', { email, password });
            const response = await apiClient.post<LoginResponse>(
                '/api/v1/auth/login/admin',
                {
                    username: email,
                    password,
                } as LoginRequest
            );
            const { access_token, refresh_token } = response.data.data;
            setAuth(null as any, access_token, refresh_token);
            navigate('/dashboard');
        } catch (err: any) {
            setError(
                err.response?.data?.message ||
                    'Đăng nhập thất bại. Vui lòng thử lại.'
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container maxWidth="lg" sx={{ zIndex: 1, py: 4 }}>
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' },
                    gap: { xs: 2, md: 6 },
                    alignItems: 'center',
                    minHeight: '100vh',
                }}
            >
                {/* Left side - Branding & Features */}
                <Box
                    sx={{
                        display: { xs: 'none', md: 'flex' },
                        flexDirection: 'column',
                        justifyContent: 'center',
                        animation: 'slideInLeft 0.8s ease-out',
                        '@keyframes slideInLeft': {
                            from: {
                                opacity: 0,
                                transform: 'translateX(-40px)',
                            },
                            to: {
                                opacity: 1,
                                transform: 'translateX(0)',
                            },
                        },
                    }}
                >
                    {/* Logo Section */}
                    <Box sx={{ mb: 4 }}>
                        <Box
                            sx={{
                                width: 80,
                                height: 80,
                                borderRadius: '16px',
                                background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                mb: 3,
                                boxShadow: '0 20px 40px rgba(59, 130, 246, 0.3)',
                            }}
                        >
                            <Security sx={{ fontSize: 48, color: 'white' }} />
                        </Box>

                        <Typography
                            variant="h3"
                            fontWeight="800"
                            mb={1}
                            sx={{
                                background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                backgroundClip: 'text',
                                WebkitBackgroundClip: 'text',
                                WebkitTextFillColor: 'transparent',
                                color: '#3b82f6',
                            }}
                        >
                            CIO Verify
                        </Typography>

                        <Typography
                            sx={{
                                color: '#cbd5e1',
                                fontSize: '1rem',
                                fontWeight: '500',
                            }}
                        >
                            Face Recognition System
                        </Typography>
                    </Box>

                    {/* Description */}
                    <Typography
                        variant="h6"
                        fontWeight="300"
                        mb={4}
                        sx={{ color: '#cbd5e1', fontSize: '1.1rem', lineHeight: 1.6 }}
                    >
                        Hệ thống quản lý chấm công và nhận dạng khuôn mặt tiên tiến. Đảm bảo bảo mật cao cấp, hiệu suất nhanh chóng và quản lý hiệu quả.
                    </Typography>

                    {/* Features List */}
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                        {[
                            { icon: '🔒', title: 'Bảo mật cao cấp', desc: 'Mã hóa end-to-end' },
                            { icon: '⚡', title: 'Hiệu suất nhanh', desc: 'Xử lý realtime' },
                            { icon: '📊', title: 'Thống kê chi tiết', desc: 'Dashboard toàn diện' },
                        ].map((feature, idx) => (
                            <Box
                                key={idx}
                                sx={{
                                    display: 'flex',
                                    gap: 3,
                                    p: 2,
                                    borderRadius: '12px',
                                    background: 'rgba(59, 130, 246, 0.05)',
                                    border: '1px solid rgba(59, 130, 246, 0.1)',
                                    transition: 'all 0.3s ease',
                                    '&:hover': {
                                        background: 'rgba(59, 130, 246, 0.1)',
                                        borderColor: 'rgba(59, 130, 246, 0.2)',
                                        transform: 'translateX(8px)',
                                    },
                                }}
                            >
                                <Typography sx={{ fontSize: '2rem' }}>
                                    {feature.icon}
                                </Typography>
                                <Box>
                                    <Typography
                                        sx={{
                                            color: '#e2e8f0',
                                            fontWeight: '600',
                                            mb: 0.5,
                                        }}
                                    >
                                        {feature.title}
                                    </Typography>
                                    <Typography
                                        sx={{
                                            color: '#94a3b8',
                                            fontSize: '0.875rem',
                                        }}
                                    >
                                        {feature.desc}
                                    </Typography>
                                </Box>
                            </Box>
                        ))}
                    </Box>
                </Box>

                {/* Right side - Login Form */}
                <Paper
                    elevation={0}
                    sx={{
                        p: { xs: 3, sm: 4, md: 5 },
                        borderRadius: '20px',
                        background: 'rgba(15, 23, 42, 0.5)',
                        backdropFilter: 'blur(20px)',
                        border: '1px solid rgba(71, 85, 105, 0.5)',
                        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
                        animation: 'slideInRight 0.8s ease-out',
                        '@keyframes slideInRight': {
                            from: {
                                opacity: 0,
                                transform: 'translateX(40px)',
                            },
                            to: {
                                opacity: 1,
                                transform: 'translateX(0)',
                            },
                        },
                    }}
                >
                    {/* Form Header */}
                    <Box sx={{ mb: 4 }}>
                        <Typography
                            variant="h4"
                            fontWeight="700"
                            mb={1}
                            sx={{ color: '#f1f5f9' }}
                        >
                            Đăng nhập
                        </Typography>
                        <Typography
                            variant="body2"
                            sx={{ color: '#cbd5e1', fontSize: '0.95rem' }}
                        >
                            Nhập thông tin tài khoản quản trị viên
                        </Typography>
                    </Box>

                    <Box component="form" onSubmit={handleSubmit}>
                        {/* Error Alert */}
                        {error && (
                            <Alert
                                severity="error"
                                sx={{
                                    mb: 3,
                                    borderRadius: '12px',
                                    background: 'rgba(239, 68, 68, 0.1)',
                                    border: '1px solid rgba(239, 68, 68, 0.3)',
                                    color: '#fca5a5',
                                    '& .MuiAlert-icon': {
                                        color: '#fca5a5',
                                    },
                                    animation: 'slideDown 0.3s ease-out',
                                    '@keyframes slideDown': {
                                        from: {
                                            opacity: 0,
                                            transform: 'translateY(-10px)',
                                        },
                                        to: {
                                            opacity: 1,
                                            transform: 'translateY(0)',
                                        },
                                    },
                                }}
                            >
                                {error}
                            </Alert>
                        )}

                        {/* Email Field */}
                        <Box sx={{ mb: 3 }}>
                            <Typography
                                variant="subtitle2"
                                fontWeight="600"
                                mb={1.5}
                                sx={{ color: '#e2e8f0', fontSize: '0.9rem' }}
                            >
                                Email
                            </Typography>
                            <TextField
                                fullWidth
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                placeholder="admin@example.com"
                                InputProps={{
                                    startAdornment: (
                                        <EmailIcon
                                            sx={{
                                                mr: 1.5,
                                                color: '#64748b',
                                                fontSize: '1.3rem',
                                            }}
                                        />
                                    ),
                                }}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: '12px',
                                        backgroundColor: 'rgba(30, 41, 59, 0.4)',
                                        transition: 'all 0.3s ease',
                                        border: '1.5px solid rgba(71, 85, 105, 0.3)',
                                        '&:hover': {
                                            backgroundColor: 'rgba(30, 41, 59, 0.6)',
                                            borderColor: 'rgba(100, 116, 139, 0.5)',
                                        },
                                        '&.Mui-focused': {
                                            backgroundColor: 'rgba(30, 41, 59, 0.8)',
                                            boxShadow: '0 0 0 3px rgba(59, 130, 246, 0.2)',
                                            '& fieldset': {
                                                borderColor: '#3b82f6 !important',
                                            },
                                        },
                                    },
                                    '& .MuiOutlinedInput-input': {
                                        p: '14px',
                                        fontSize: '0.95rem',
                                        color: '#e2e8f0',
                                        '&::placeholder': {
                                            color: '#64748b',
                                            opacity: 0.6,
                                        },
                                    },
                                }}
                            />
                        </Box>

                        {/* Password Field */}
                        <Box sx={{ mb: 3 }}>
                            <Typography
                                variant="subtitle2"
                                fontWeight="600"
                                mb={1.5}
                                sx={{ color: '#e2e8f0', fontSize: '0.9rem' }}
                            >
                                Mật khẩu
                            </Typography>
                            <TextField
                                fullWidth
                                type={showPassword ? 'text' : 'password'}
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                placeholder="••••••••"
                                InputProps={{
                                    startAdornment: (
                                        <LockIcon
                                            sx={{
                                                mr: 1.5,
                                                color: '#64748b',
                                                fontSize: '1.3rem',
                                            }}
                                        />
                                    ),
                                    endAdornment: (
                                        <Box
                                            onClick={() => setShowPassword(!showPassword)}
                                            sx={{
                                                cursor: 'pointer',
                                                display: 'flex',
                                                alignItems: 'center',
                                                color: '#64748b',
                                                transition: 'color 0.2s',
                                                '&:hover': {
                                                    color: '#94a3b8',
                                                },
                                            }}
                                        >
                                            {showPassword ? (
                                                <VisibilityOff fontSize="small" />
                                            ) : (
                                                <Visibility fontSize="small" />
                                            )}
                                        </Box>
                                    ),
                                }}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: '12px',
                                        backgroundColor: 'rgba(30, 41, 59, 0.4)',
                                        transition: 'all 0.3s ease',
                                        border: '1.5px solid rgba(71, 85, 105, 0.3)',
                                        '&:hover': {
                                            backgroundColor: 'rgba(30, 41, 59, 0.6)',
                                            borderColor: 'rgba(100, 116, 139, 0.5)',
                                        },
                                        '&.Mui-focused': {
                                            backgroundColor: 'rgba(30, 41, 59, 0.8)',
                                            boxShadow: '0 0 0 3px rgba(59, 130, 246, 0.2)',
                                            '& fieldset': {
                                                borderColor: '#3b82f6 !important',
                                            },
                                        },
                                    },
                                    '& .MuiOutlinedInput-input': {
                                        p: '14px',
                                        fontSize: '0.95rem',
                                        color: '#e2e8f0',
                                        '&::placeholder': {
                                            color: '#64748b',
                                            opacity: 0.6,
                                        },
                                    },
                                }}
                            />
                        </Box>

                        {/* Remember & Forgot Password */}
                        <Box
                            sx={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center',
                                mb: 4,
                            }}
                        >
                            <FormControlLabel
                                control={
                                    <Checkbox
                                        checked={rememberMe}
                                        onChange={(e) => setRememberMe(e.target.checked)}
                                        size="small"
                                        sx={{
                                            color: '#64748b',
                                            '&.Mui-checked': {
                                                color: '#3b82f6',
                                            },
                                        }}
                                    />
                                }
                                label={
                                    <Typography
                                        variant="body2"
                                        sx={{ color: '#cbd5e1', fontSize: '0.875rem' }}
                                    >
                                        Nhớ tài khoản
                                    </Typography>
                                }
                            />
                            <Typography
                                component="a"
                                href="#"
                                onClick={(e) => {
                                    e.preventDefault();
                                    // Handle forgot password
                                }}
                                sx={{
                                    color: '#60a5fa',
                                    textDecoration: 'none',
                                    fontSize: '0.875rem',
                                    fontWeight: '500',
                                    transition: 'all 0.2s',
                                    '&:hover': {
                                        color: '#3b82f6',
                                        textDecoration: 'underline',
                                    },
                                }}
                            >
                                Quên mật khẩu?
                            </Typography>
                        </Box>

                        {/* Login Button */}
                        <Button
                            fullWidth
                            type="submit"
                            variant="contained"
                            size="large"
                            disabled={loading}
                            sx={{
                                p: '14px 24px',
                                fontSize: '1rem',
                                fontWeight: '600',
                                textTransform: 'none',
                                borderRadius: '12px',
                                background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                transition: 'all 0.3s ease',
                                position: 'relative',
                                overflow: 'hidden',
                                mb: 2,
                                '&::before': {
                                    content: '""',
                                    position: 'absolute',
                                    top: 0,
                                    left: '-100%',
                                    width: '100%',
                                    height: '100%',
                                    background:
                                        'linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent)',
                                    transition: 'left 0.6s ease',
                                },
                                '&:hover:not(:disabled)': {
                                    '&::before': {
                                        left: '100%',
                                    },
                                    transform: 'translateY(-2px)',
                                    boxShadow: '0 20px 40px rgba(59, 130, 246, 0.4)',
                                },
                                '&:active:not(:disabled)': {
                                    transform: 'translateY(0)',
                                },
                                '&:disabled': {
                                    opacity: 0.6,
                                    cursor: 'not-allowed',
                                },
                            }}
                        >
                            {loading ? (
                                <CircularProgress size={24} sx={{ color: '#fff' }} />
                            ) : (
                                <Box
                                    sx={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        gap: 1,
                                    }}
                                >
                                    Đăng nhập
                                    <ArrowForward fontSize="small" />
                                </Box>
                            )}
                        </Button>

                        {/* Divider */}
                        <Box
                            sx={{
                                height: '1px',
                                background:
                                    'linear-gradient(90deg, transparent, rgba(71, 85, 105, 0.3), transparent)',
                                my: 3,
                            }}
                        />

                        {/* Footer Text */}
                        <Typography
                            variant="caption"
                            textAlign="center"
                            display="block"
                            sx={{
                                color: '#64748b',
                                fontSize: '0.8rem',
                            }}
                        >
                            © 2025 CIO Verify Face. Bảo mật dữ liệu của bạn.
                        </Typography>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
};
