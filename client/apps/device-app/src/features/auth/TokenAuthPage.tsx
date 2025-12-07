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
    Key as KeyIcon,
    Visibility,
    VisibilityOff,
    ArrowForward,
    Devices as DevicesIcon,
} from '@mui/icons-material';
import { useDeviceStore } from '@/stores/deviceStore';
import axios from 'axios';

export const TokenAuthPage: React.FC = () => {
    const navigate = useNavigate();
    const { setDeviceToken, setDeviceInfo } = useDeviceStore();
    const [token, setToken] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [showToken, setShowToken] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!token.trim()) {
            setError('Vui lòng nhập token thiết bị');
            return;
        }

        setLoading(true);

        try {
            // Authenticate device using the provided token
            const response = await axios.post('/api/v1/auth/device', {
                token: token.trim(),
            });

            const { device_id, device_name, company_id, company_name, location } = response.data.data;

            // Store device info
            setDeviceToken(token.trim());
            setDeviceInfo({
                deviceId: device_id,
                deviceName: device_name,
                companyId: company_id,
                companyName: company_name,
                location: location,
            });

            // Navigate to attendance page
            navigate('/attendance');
        } catch (err: any) {
            console.error('Device authentication error:', err);
            setError(
                err.response?.data?.message ||
                'Xác thực thiết bị thất bại. Vui lòng kiểm tra token và thử lại.'
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box
            sx={{
                width: '100%',
                maxWidth: 480,
                px: 3,
                animation: 'slideUp 0.6s ease-out',
                '@keyframes slideUp': {
                    from: {
                        opacity: 0,
                        transform: 'translateY(30px)',
                    },
                    to: {
                        opacity: 1,
                        transform: 'translateY(0)',
                    },
                },
            }}
        >
            <Paper
                elevation={0}
                sx={{
                    p: { xs: 3, sm: 4, md: 5 },
                    borderRadius: 4,
                    background: 'rgba(255, 255, 255, 0.95)',
                    backdropFilter: 'blur(20px)',
                    border: '1px solid rgba(226, 232, 240, 0.8)',
                    boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.15)',
                }}
            >
                {/* Logo & Header */}
                <Box sx={{ textAlign: 'center', mb: 4 }}>
                    <Box
                        sx={{
                            width: 80,
                            height: 80,
                            borderRadius: 3,
                            background: 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            mx: 'auto',
                            mb: 3,
                            boxShadow: '0 20px 40px rgba(37, 99, 235, 0.3)',
                        }}
                    >
                        <DevicesIcon sx={{ fontSize: 40, color: 'white' }} />
                    </Box>

                    <Typography
                        variant="h4"
                        fontWeight="700"
                        sx={{
                            background: 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                            backgroundClip: 'text',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent',
                            mb: 1,
                        }}
                    >
                        CIO Verify Device
                    </Typography>

                    <Typography variant="body1" color="text.secondary">
                        Nhập token để xác thực thiết bị chấm công
                    </Typography>
                </Box>

                {/* Form */}
                <Box component="form" onSubmit={handleSubmit}>
                    {/* Error Alert */}
                    {error && (
                        <Alert
                            severity="error"
                            sx={{
                                mb: 3,
                                borderRadius: 2,
                                animation: 'fadeIn 0.3s ease-out',
                            }}
                        >
                            {error}
                        </Alert>
                    )}

                    {/* Token Input */}
                    <Box sx={{ mb: 3 }}>
                        <Typography
                            variant="subtitle2"
                            fontWeight="600"
                            mb={1.5}
                            color="text.primary"
                        >
                            Device Token
                        </Typography>
                        <TextField
                            fullWidth
                            type={showToken ? 'text' : 'password'}
                            value={token}
                            onChange={(e) => setToken(e.target.value)}
                            placeholder="Nhập token thiết bị..."
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <KeyIcon sx={{ color: 'text.secondary' }} />
                                    </InputAdornment>
                                ),
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            onClick={() => setShowToken(!showToken)}
                                            edge="end"
                                            size="small"
                                        >
                                            {showToken ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                            sx={{
                                '& .MuiOutlinedInput-root': {
                                    backgroundColor: '#f8fafc',
                                    '&:hover': {
                                        backgroundColor: '#f1f5f9',
                                    },
                                    '&.Mui-focused': {
                                        backgroundColor: '#fff',
                                        boxShadow: '0 0 0 3px rgba(37, 99, 235, 0.15)',
                                    },
                                },
                            }}
                        />
                    </Box>

                    {/* Submit Button */}
                    <Button
                        fullWidth
                        type="submit"
                        variant="contained"
                        size="large"
                        disabled={loading}
                        sx={{
                            py: 1.75,
                            fontSize: '1.1rem',
                            fontWeight: '600',
                            background: 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)',
                            borderRadius: 3,
                            transition: 'all 0.3s ease',
                            '&:hover:not(:disabled)': {
                                transform: 'translateY(-2px)',
                                boxShadow: '0 20px 40px rgba(37, 99, 235, 0.35)',
                            },
                            '&:active:not(:disabled)': {
                                transform: 'translateY(0)',
                            },
                            '&:disabled': {
                                opacity: 0.7,
                            },
                        }}
                    >
                        {loading ? (
                            <CircularProgress size={26} sx={{ color: '#fff' }} />
                        ) : (
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                Xác thực thiết bị
                                <ArrowForward />
                            </Box>
                        )}
                    </Button>

                    {/* Help Text */}
                    <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ display: 'block', textAlign: 'center', mt: 3 }}
                    >
                        Token được cung cấp bởi quản trị viên công ty
                    </Typography>
                </Box>
            </Paper>

            {/* Footer */}
            <Typography
                variant="caption"
                color="text.secondary"
                sx={{ display: 'block', textAlign: 'center', mt: 4 }}
            >
                © 2025 CIO Verify Face - Thiết bị chấm công
            </Typography>
        </Box>
    );
};
