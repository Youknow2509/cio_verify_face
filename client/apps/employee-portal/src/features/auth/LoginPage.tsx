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
    InputAdornment,
    IconButton,
} from '@mui/material';
import {
    Lock as LockIcon,
    MailOutline as EmailIcon,
    Visibility,
    VisibilityOff,
    PersonOutline,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';
import { authApi } from '@/services/api';

export const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const { setAuth } = useAuthStore();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const loginResponse: any = await authApi.login({
                username,
                password,
            });
            const loginPayload = loginResponse?.data ?? loginResponse; // interceptor already unwraps to data
            const tokenSource = loginPayload?.data ?? loginPayload;
            const accessToken = tokenSource?.access_token;
            const refreshToken = tokenSource?.refresh_token;

            // Some APIs already return user data alongside tokens
            const userFromLogin = tokenSource?.user ?? tokenSource?.data;

            if (!accessToken || !refreshToken) {
                throw new Error('Không nhận được token từ máy chủ.');
            }

            // Put tokens in storage so request interceptor can attach them for /me
            localStorage.setItem('access_token', accessToken);
            localStorage.setItem('refresh_token', refreshToken);

            const userInfo: any = userFromLogin ?? (await authApi.getMe());
            const userData = userInfo?.data ?? userInfo;

            setAuth(userData, accessToken, refreshToken);
            navigate('/dashboard');
        } catch (err: any) {
            setError(
                err.message ||
                    'Đăng nhập thất bại. Vui lòng kiểm tra lại thông tin.'
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
                position: 'relative',
                overflow: 'hidden',
                '&::before': {
                    content: '""',
                    position: 'absolute',
                    top: '-50%',
                    left: '-50%',
                    width: '200%',
                    height: '200%',
                    background:
                        'radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 70%)',
                    animation: 'pulse 15s ease-in-out infinite',
                },
                '@keyframes pulse': {
                    '0%, 100%': { opacity: 0.5 },
                    '50%': { opacity: 1 },
                },
            }}
        >
            <Container maxWidth="sm" sx={{ zIndex: 1 }}>
                <Paper
                    elevation={24}
                    sx={{
                        p: { xs: 3, sm: 5 },
                        borderRadius: 4,
                        background: 'rgba(15, 23, 42, 0.8)',
                        backdropFilter: 'blur(20px)',
                        border: '1px solid rgba(71, 85, 105, 0.5)',
                    }}
                >
                    {/* Logo and Title */}
                    <Box sx={{ textAlign: 'center', mb: 4 }}>
                        <Box
                            sx={{
                                width: 80,
                                height: 80,
                                borderRadius: '50%',
                                background:
                                    'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                margin: '0 auto 20px',
                                boxShadow:
                                    '0 10px 40px rgba(59, 130, 246, 0.4)',
                            }}
                        >
                            <PersonOutline
                                sx={{ fontSize: 48, color: 'white' }}
                            />
                        </Box>
                        <Typography variant="h4" fontWeight="700" mb={1}>
                            Cổng Nhân Viên
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Hệ thống chấm công nhận dạng khuôn mặt
                        </Typography>
                    </Box>

                    <Box component="form" onSubmit={handleSubmit}>
                        {error && (
                            <Alert
                                severity="error"
                                sx={{
                                    mb: 3,
                                    borderRadius: 2,
                                }}
                            >
                                {error}
                            </Alert>
                        )}

                        <Box sx={{ mb: 3 }}>
                            <Typography
                                variant="subtitle2"
                                fontWeight="600"
                                mb={1}
                            >
                                Tên đăng nhập
                            </Typography>
                            <TextField
                                fullWidth
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                required
                                placeholder="Nhập email hoặc mã nhân viên"
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <EmailIcon
                                                sx={{ color: 'text.secondary' }}
                                            />
                                        </InputAdornment>
                                    ),
                                }}
                            />
                        </Box>

                        <Box sx={{ mb: 4 }}>
                            <Typography
                                variant="subtitle2"
                                fontWeight="600"
                                mb={1}
                            >
                                Mật khẩu
                            </Typography>
                            <TextField
                                fullWidth
                                type={showPassword ? 'text' : 'password'}
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                placeholder="Nhập mật khẩu"
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <LockIcon
                                                sx={{ color: 'text.secondary' }}
                                            />
                                        </InputAdornment>
                                    ),
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                onClick={() =>
                                                    setShowPassword(
                                                        !showPassword
                                                    )
                                                }
                                                edge="end"
                                            >
                                                {showPassword ? (
                                                    <VisibilityOff />
                                                ) : (
                                                    <Visibility />
                                                )}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                        </Box>

                        <Button
                            fullWidth
                            type="submit"
                            variant="contained"
                            size="large"
                            disabled={loading}
                            sx={{
                                py: 1.5,
                                fontSize: '1rem',
                                fontWeight: '600',
                                background:
                                    'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                '&:hover': {
                                    background:
                                        'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                                    transform: 'translateY(-2px)',
                                    boxShadow:
                                        '0 10px 30px rgba(59, 130, 246, 0.4)',
                                },
                                transition: 'all 0.3s ease',
                            }}
                        >
                            {loading ? (
                                <CircularProgress
                                    size={24}
                                    sx={{ color: '#fff' }}
                                />
                            ) : (
                                'Đăng nhập'
                            )}
                        </Button>

                        <Typography
                            variant="caption"
                            display="block"
                            textAlign="center"
                            sx={{ mt: 3, color: 'text.secondary' }}
                        >
                            © 2025 CIO Verify Face - Cổng thông tin nhân viên
                        </Typography>
                    </Box>
                </Paper>
            </Container>
        </Box>
    );
};
