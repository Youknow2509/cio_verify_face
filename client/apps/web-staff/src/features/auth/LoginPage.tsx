import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    TextField,
    Button,
    Typography,
    Alert,
    CircularProgress,
    Paper,
    InputAdornment,
    IconButton,
} from '@mui/material';
import {
    Email,
    Lock,
    Visibility,
    VisibilityOff,
    ArrowForward,
} from '@mui/icons-material';
import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';
import { ThemeToggle } from '@/components/ThemeToggle';

export const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const { setAuth } = useAuthStore();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!email || !password) {
            setError('Vui lòng nhập đầy đủ thông tin');
            return;
        }

        setLoading(true);

        try {
            const response = await axios.post('/api/v1/auth/login', {
                email,
                password,
            });

            const { user, access_token, refresh_token } = response.data.data;

            setAuth(user, access_token, refresh_token);
            navigate('/dashboard');
        } catch (err: any) {
            console.error('Login error:', err);
            setError(
                err.response?.data?.message || 'Đăng nhập thất bại. Vui lòng thử lại.'
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box
            sx={{
                width: '100%',
                maxWidth: 440,
                px: 3,
                animation: 'slideUp 0.6s ease-out',
            }}
        >
            <Paper
                elevation={0}
                sx={{
                    p: { xs: 3, sm: 4 },
                    borderRadius: 4,
                    bgcolor: 'background.paper',
                    border: '1px solid',
                    borderColor: 'divider',
                    boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.15)',
                }}
            >
                {/* Header */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                    <Box>
                        <Typography variant="h4" fontWeight="700" color="text.primary" mb={0.5}>
                            Xin chào! 👋
                        </Typography>
                        <Typography variant="body1" color="text.secondary">
                            Đăng nhập để tiếp tục
                        </Typography>
                    </Box>
                    <ThemeToggle />
                </Box>

                {/* Form */}
                <Box component="form" onSubmit={handleSubmit}>
                    {error && (
                        <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
                            {error}
                        </Alert>
                    )}

                    <Box sx={{ mb: 2.5 }}>
                        <Typography variant="subtitle2" fontWeight="600" mb={1} color="text.primary">
                            Email
                        </Typography>
                        <TextField
                            fullWidth
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="email@company.com"
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <Email sx={{ color: 'text.secondary' }} />
                                    </InputAdornment>
                                ),
                            }}
                        />
                    </Box>

                    <Box sx={{ mb: 3 }}>
                        <Typography variant="subtitle2" fontWeight="600" mb={1} color="text.primary">
                            Mật khẩu
                        </Typography>
                        <TextField
                            fullWidth
                            type={showPassword ? 'text' : 'password'}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="••••••••"
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <Lock sx={{ color: 'text.secondary' }} />
                                    </InputAdornment>
                                ),
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            onClick={() => setShowPassword(!showPassword)}
                                            edge="end"
                                            size="small"
                                        >
                                            {showPassword ? <VisibilityOff /> : <Visibility />}
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
                            background: 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                            borderRadius: 2.5,
                            '&:hover:not(:disabled)': {
                                transform: 'translateY(-1px)',
                                boxShadow: '0 10px 30px rgba(37, 99, 235, 0.3)',
                            },
                        }}
                    >
                        {loading ? (
                            <CircularProgress size={24} sx={{ color: '#fff' }} />
                        ) : (
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                Đăng nhập
                                <ArrowForward />
                            </Box>
                        )}
                    </Button>
                </Box>
            </Paper>

            <Typography
                variant="caption"
                color="text.secondary"
                sx={{ display: 'block', textAlign: 'center', mt: 3 }}
            >
                © 2025 CIO Verify Face - Staff Portal
            </Typography>
        </Box>
    );
};
