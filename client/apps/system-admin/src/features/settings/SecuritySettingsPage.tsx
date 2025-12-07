import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    Grid,
    Switch,
    FormControlLabel,
    TextField,
    Divider,
    Alert,
    Chip,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
} from '@mui/material';
import {
    ArrowBack,
    Security,
    Key,
    Shield,
    Lock,
    Delete,
    Add,
} from '@mui/icons-material';

// Mock data for active sessions
const activeSessions = [
    { id: '1', user: 'admin@system.vn', ip: '192.168.1.100', device: 'Chrome / Windows', lastActive: '2024-12-08 00:10', location: 'Ho Chi Minh City' },
    { id: '2', user: 'operator@system.vn', ip: '192.168.1.105', device: 'Firefox / macOS', lastActive: '2024-12-07 23:45', location: 'Hanoi' },
];

export const SecuritySettingsPage: React.FC = () => {
    const navigate = useNavigate();
    const [settings, setSettings] = useState({
        mfaEnabled: true,
        sessionTimeout: 30,
        maxLoginAttempts: 5,
        passwordMinLength: 12,
        passwordRequireSpecial: true,
        passwordRequireNumbers: true,
        passwordExpireDays: 90,
        ipWhitelistEnabled: false,
    });

    const handleChange = (field: string, value: any) => {
        setSettings(prev => ({ ...prev, [field]: value }));
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
                <Button startIcon={<ArrowBack />} onClick={() => navigate('/settings')}>
                    Quay lại
                </Button>
                <Box sx={{ flex: 1 }}>
                    <Typography variant="h4" fontWeight="700">
                        System Security
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Cấu hình bảo mật và chính sách truy cập
                    </Typography>
                </Box>
                <Button variant="contained" color="primary">
                    Lưu thay đổi
                </Button>
            </Box>

            <Grid container spacing={3}>
                {/* Authentication Settings */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Key color="primary" />
                                <Typography variant="h6" fontWeight="600">
                                    Authentication
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.mfaEnabled}
                                            onChange={(e) => handleChange('mfaEnabled', e.target.checked)}
                                        />
                                    }
                                    label={
                                        <Box>
                                            <Typography fontWeight="500">Bắt buộc MFA (2FA)</Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                Yêu cầu xác thực 2 yếu tố cho tất cả System Admin
                                            </Typography>
                                        </Box>
                                    }
                                />

                                <TextField
                                    label="Session Timeout (phút)"
                                    type="number"
                                    size="small"
                                    value={settings.sessionTimeout}
                                    onChange={(e) => handleChange('sessionTimeout', parseInt(e.target.value))}
                                    helperText="Thời gian tự động đăng xuất khi không hoạt động"
                                />

                                <TextField
                                    label="Max Login Attempts"
                                    type="number"
                                    size="small"
                                    value={settings.maxLoginAttempts}
                                    onChange={(e) => handleChange('maxLoginAttempts', parseInt(e.target.value))}
                                    helperText="Số lần đăng nhập sai tối đa trước khi khóa tài khoản"
                                />
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Password Policy */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Lock color="secondary" />
                                <Typography variant="h6" fontWeight="600">
                                    Password Policy
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                                <TextField
                                    label="Độ dài tối thiểu"
                                    type="number"
                                    size="small"
                                    value={settings.passwordMinLength}
                                    onChange={(e) => handleChange('passwordMinLength', parseInt(e.target.value))}
                                />

                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.passwordRequireSpecial}
                                            onChange={(e) => handleChange('passwordRequireSpecial', e.target.checked)}
                                        />
                                    }
                                    label="Yêu cầu ký tự đặc biệt (!@#$%)"
                                />

                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.passwordRequireNumbers}
                                            onChange={(e) => handleChange('passwordRequireNumbers', e.target.checked)}
                                        />
                                    }
                                    label="Yêu cầu chữ số (0-9)"
                                />

                                <TextField
                                    label="Password Expiry (ngày)"
                                    type="number"
                                    size="small"
                                    value={settings.passwordExpireDays}
                                    onChange={(e) => handleChange('passwordExpireDays', parseInt(e.target.value))}
                                    helperText="Số ngày trước khi yêu cầu đổi mật khẩu (0 = không hết hạn)"
                                />
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* IP Whitelist */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 3 }}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <Shield color="warning" />
                                    <Typography variant="h6" fontWeight="600">
                                        IP Whitelist
                                    </Typography>
                                </Box>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.ipWhitelistEnabled}
                                            onChange={(e) => handleChange('ipWhitelistEnabled', e.target.checked)}
                                        />
                                    }
                                    label="Bật IP Whitelist"
                                />
                            </Box>

                            {settings.ipWhitelistEnabled ? (
                                <Box>
                                    <Alert severity="warning" sx={{ mb: 2 }}>
                                        Chỉ các IP trong danh sách mới có thể truy cập System Admin Portal
                                    </Alert>
                                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
                                        <Chip label="192.168.1.0/24" onDelete={() => { }} />
                                        <Chip label="10.0.0.0/8" onDelete={() => { }} />
                                        <Chip label="103.45.67.89" onDelete={() => { }} />
                                    </Box>
                                    <Button startIcon={<Add />} size="small" variant="outlined">
                                        Thêm IP/CIDR
                                    </Button>
                                </Box>
                            ) : (
                                <Typography color="text.secondary">
                                    IP Whitelist đang tắt. Mọi IP đều có thể truy cập (sau khi xác thực).
                                </Typography>
                            )}
                        </CardContent>
                    </Card>
                </Grid>

                {/* Active Sessions */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Security color="info" />
                                <Typography variant="h6" fontWeight="600">
                                    Active Sessions
                                </Typography>
                            </Box>

                            <TableContainer>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>User</TableCell>
                                            <TableCell>IP Address</TableCell>
                                            <TableCell>Device</TableCell>
                                            <TableCell>Location</TableCell>
                                            <TableCell>Last Active</TableCell>
                                            <TableCell align="right">Action</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {activeSessions.map((session) => (
                                            <TableRow key={session.id}>
                                                <TableCell>{session.user}</TableCell>
                                                <TableCell><code>{session.ip}</code></TableCell>
                                                <TableCell>{session.device}</TableCell>
                                                <TableCell>{session.location}</TableCell>
                                                <TableCell>{session.lastActive}</TableCell>
                                                <TableCell align="right">
                                                    <IconButton size="small" color="error">
                                                        <Delete fontSize="small" />
                                                    </IconButton>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};
