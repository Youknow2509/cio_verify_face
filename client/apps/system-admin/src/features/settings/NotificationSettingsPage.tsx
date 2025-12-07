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
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
} from '@mui/material';
import {
    ArrowBack,
    Email,
    Notifications,
    Sms,
    Webhook,
    CheckCircle,
    Error as ErrorIcon,
} from '@mui/icons-material';

// Mock email templates
const emailTemplates = [
    { id: 'welcome', name: 'Welcome Email', subject: 'Chào mừng đến với Face Attendance', status: 'active' },
    { id: 'password_reset', name: 'Password Reset', subject: 'Đặt lại mật khẩu', status: 'active' },
    { id: 'subscription_expiry', name: 'Subscription Expiry', subject: 'Gói dịch vụ sắp hết hạn', status: 'active' },
    { id: 'device_offline', name: 'Device Offline Alert', subject: 'Cảnh báo thiết bị offline', status: 'inactive' },
];

export const NotificationSettingsPage: React.FC = () => {
    const navigate = useNavigate();
    const [settings, setSettings] = useState({
        smtpHost: 'smtp.sendgrid.net',
        smtpPort: 587,
        smtpUser: 'apikey',
        smtpPassword: '••••••••••••',
        senderEmail: 'noreply@faceattendance.vn',
        senderName: 'Face Attendance System',
        enableEmailNotifications: true,
        enableWebhooks: false,
        webhookUrl: '',
        slackWebhook: '',
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
                        Notification Settings
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Cấu hình email server và mẫu thông báo
                    </Typography>
                </Box>
                <Button variant="contained" color="primary">
                    Lưu thay đổi
                </Button>
            </Box>

            <Grid container spacing={3}>
                {/* SMTP Configuration */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Email color="primary" />
                                <Typography variant="h6" fontWeight="600">
                                    SMTP Configuration
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                                <TextField
                                    label="SMTP Host"
                                    size="small"
                                    fullWidth
                                    value={settings.smtpHost}
                                    onChange={(e) => handleChange('smtpHost', e.target.value)}
                                />

                                <Box sx={{ display: 'flex', gap: 2 }}>
                                    <TextField
                                        label="Port"
                                        size="small"
                                        type="number"
                                        value={settings.smtpPort}
                                        onChange={(e) => handleChange('smtpPort', parseInt(e.target.value))}
                                        sx={{ width: 120 }}
                                    />
                                    <TextField
                                        label="Username"
                                        size="small"
                                        fullWidth
                                        value={settings.smtpUser}
                                        onChange={(e) => handleChange('smtpUser', e.target.value)}
                                    />
                                </Box>

                                <TextField
                                    label="Password"
                                    size="small"
                                    type="password"
                                    fullWidth
                                    value={settings.smtpPassword}
                                    onChange={(e) => handleChange('smtpPassword', e.target.value)}
                                />

                                <Divider />

                                <TextField
                                    label="Sender Email"
                                    size="small"
                                    fullWidth
                                    value={settings.senderEmail}
                                    onChange={(e) => handleChange('senderEmail', e.target.value)}
                                />

                                <TextField
                                    label="Sender Name"
                                    size="small"
                                    fullWidth
                                    value={settings.senderName}
                                    onChange={(e) => handleChange('senderName', e.target.value)}
                                />

                                <Button variant="outlined" size="small">
                                    Gửi Email Test
                                </Button>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Webhook Integration */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Webhook color="secondary" />
                                <Typography variant="h6" fontWeight="600">
                                    Webhook Integration
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.enableWebhooks}
                                            onChange={(e) => handleChange('enableWebhooks', e.target.checked)}
                                        />
                                    }
                                    label="Enable Webhooks"
                                />

                                {settings.enableWebhooks && (
                                    <>
                                        <TextField
                                            label="Webhook URL"
                                            size="small"
                                            fullWidth
                                            placeholder="https://your-server.com/webhook"
                                            value={settings.webhookUrl}
                                            onChange={(e) => handleChange('webhookUrl', e.target.value)}
                                            helperText="Nhận events: company.created, subscription.updated, etc."
                                        />

                                        <TextField
                                            label="Slack Webhook"
                                            size="small"
                                            fullWidth
                                            placeholder="https://hooks.slack.com/services/..."
                                            value={settings.slackWebhook}
                                            onChange={(e) => handleChange('slackWebhook', e.target.value)}
                                            helperText="Gửi alert tới Slack channel"
                                        />
                                    </>
                                )}

                                <Alert severity="info" sx={{ mt: 2 }}>
                                    Webhooks sẽ gửi real-time events cho các hệ thống bên ngoài
                                </Alert>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Email Templates */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Notifications color="warning" />
                                <Typography variant="h6" fontWeight="600">
                                    Email Templates
                                </Typography>
                            </Box>

                            <TableContainer>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Template Name</TableCell>
                                            <TableCell>Subject</TableCell>
                                            <TableCell>Status</TableCell>
                                            <TableCell align="right">Actions</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {emailTemplates.map((template) => (
                                            <TableRow key={template.id}>
                                                <TableCell>
                                                    <Typography fontWeight="500">{template.name}</Typography>
                                                </TableCell>
                                                <TableCell>{template.subject}</TableCell>
                                                <TableCell>
                                                    <Chip
                                                        size="small"
                                                        icon={template.status === 'active' ? <CheckCircle /> : <ErrorIcon />}
                                                        label={template.status}
                                                        color={template.status === 'active' ? 'success' : 'default'}
                                                    />
                                                </TableCell>
                                                <TableCell align="right">
                                                    <Button size="small">Edit</Button>
                                                    <Button size="small">Preview</Button>
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
