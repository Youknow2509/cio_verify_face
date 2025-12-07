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
    LinearProgress,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
} from '@mui/material';
import {
    ArrowBack,
    Storage,
    CloudUpload,
    Delete,
    Folder,
    Image,
    VideoFile,
    Description,
} from '@mui/icons-material';

// Mock storage buckets
const storageBuckets = [
    { name: 'face-images', size: 125.4, objects: 45620, type: 'images', lastModified: '2024-12-08' },
    { name: 'attendance-photos', size: 89.2, objects: 32100, type: 'images', lastModified: '2024-12-08' },
    { name: 'reports', size: 12.5, objects: 1250, type: 'documents', lastModified: '2024-12-07' },
    { name: 'signatures', size: 3.8, objects: 890, type: 'images', lastModified: '2024-12-06' },
    { name: 'backups', size: 45.0, objects: 30, type: 'archives', lastModified: '2024-12-01' },
];

const totalStorage = 500; // GB
const usedStorage = storageBuckets.reduce((acc, b) => acc + b.size, 0);

export const StorageSettingsPage: React.FC = () => {
    const navigate = useNavigate();
    const [settings, setSettings] = useState({
        minioEndpoint: 'minio.faceattendance.vn:9000',
        accessKey: 'AKIAIOSFODNN7EXAMPLE',
        secretKey: '••••••••••••••••••••••••••',
        region: 'ap-southeast-1',
        useSSL: true,
        autoCleanup: true,
        cleanupAfterDays: 365,
        maxFileSize: 10, // MB
        allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'application/pdf'],
    });

    const handleChange = (field: string, value: any) => {
        setSettings(prev => ({ ...prev, [field]: value }));
    };

    const getTypeIcon = (type: string) => {
        switch (type) {
            case 'images': return <Image color="primary" />;
            case 'documents': return <Description color="info" />;
            default: return <Folder color="action" />;
        }
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
                <Button startIcon={<ArrowBack />} onClick={() => navigate('/settings')}>
                    Quay lại
                </Button>
                <Box sx={{ flex: 1 }}>
                    <Typography variant="h4" fontWeight="700">
                        Storage Configuration
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Kết nối MinIO/S3 và chính sách lưu trữ
                    </Typography>
                </Box>
                <Button variant="contained" color="primary">
                    Lưu thay đổi
                </Button>
            </Box>

            <Grid container spacing={3}>
                {/* Storage Overview */}
                <Grid item xs={12}>
                    <Card sx={{ background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)', color: 'white' }}>
                        <CardContent>
                            <Grid container spacing={3} alignItems="center">
                                <Grid item xs={12} md={4}>
                                    <Typography variant="h6" sx={{ opacity: 0.9 }}>Total Storage</Typography>
                                    <Typography variant="h3" fontWeight="700">
                                        {usedStorage.toFixed(1)} GB
                                    </Typography>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>
                                        of {totalStorage} GB used ({((usedStorage / totalStorage) * 100).toFixed(1)}%)
                                    </Typography>
                                    <LinearProgress
                                        variant="determinate"
                                        value={(usedStorage / totalStorage) * 100}
                                        sx={{ mt: 2, height: 8, borderRadius: 4, bgcolor: 'rgba(255,255,255,0.2)', '& .MuiLinearProgress-bar': { bgcolor: 'white' } }}
                                    />
                                </Grid>
                                <Grid item xs={6} md={2}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Buckets</Typography>
                                    <Typography variant="h4" fontWeight="700">{storageBuckets.length}</Typography>
                                </Grid>
                                <Grid item xs={6} md={2}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Objects</Typography>
                                    <Typography variant="h4" fontWeight="700">
                                        {(storageBuckets.reduce((acc, b) => acc + b.objects, 0) / 1000).toFixed(1)}K
                                    </Typography>
                                </Grid>
                                <Grid item xs={6} md={2}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Face Images</Typography>
                                    <Typography variant="h4" fontWeight="700">45.6K</Typography>
                                </Grid>
                                <Grid item xs={6} md={2}>
                                    <Typography variant="body2" sx={{ opacity: 0.8 }}>Avg Size</Typography>
                                    <Typography variant="h4" fontWeight="700">2.8 MB</Typography>
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>

                {/* MinIO Connection */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <CloudUpload color="primary" />
                                <Typography variant="h6" fontWeight="600">
                                    MinIO/S3 Connection
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                                <TextField
                                    label="Endpoint"
                                    size="small"
                                    fullWidth
                                    value={settings.minioEndpoint}
                                    onChange={(e) => handleChange('minioEndpoint', e.target.value)}
                                />

                                <TextField
                                    label="Access Key"
                                    size="small"
                                    fullWidth
                                    value={settings.accessKey}
                                    onChange={(e) => handleChange('accessKey', e.target.value)}
                                />

                                <TextField
                                    label="Secret Key"
                                    size="small"
                                    type="password"
                                    fullWidth
                                    value={settings.secretKey}
                                    onChange={(e) => handleChange('secretKey', e.target.value)}
                                />

                                <Box sx={{ display: 'flex', gap: 2 }}>
                                    <TextField
                                        label="Region"
                                        size="small"
                                        value={settings.region}
                                        onChange={(e) => handleChange('region', e.target.value)}
                                        sx={{ flex: 1 }}
                                    />
                                    <FormControlLabel
                                        control={
                                            <Switch
                                                checked={settings.useSSL}
                                                onChange={(e) => handleChange('useSSL', e.target.checked)}
                                            />
                                        }
                                        label="Use SSL"
                                    />
                                </Box>

                                <Button variant="outlined" size="small">
                                    Test Connection
                                </Button>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Retention Policy */}
                <Grid item xs={12} lg={6}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Delete color="error" />
                                <Typography variant="h6" fontWeight="600">
                                    Retention Policy
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings.autoCleanup}
                                            onChange={(e) => handleChange('autoCleanup', e.target.checked)}
                                        />
                                    }
                                    label={
                                        <Box>
                                            <Typography fontWeight="500">Auto Cleanup</Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                Tự động xóa files cũ theo policy
                                            </Typography>
                                        </Box>
                                    }
                                />

                                {settings.autoCleanup && (
                                    <TextField
                                        label="Delete after (days)"
                                        type="number"
                                        size="small"
                                        value={settings.cleanupAfterDays}
                                        onChange={(e) => handleChange('cleanupAfterDays', parseInt(e.target.value))}
                                        helperText="Xóa attendance photos sau số ngày này"
                                    />
                                )}

                                <Divider />

                                <TextField
                                    label="Max File Size (MB)"
                                    type="number"
                                    size="small"
                                    value={settings.maxFileSize}
                                    onChange={(e) => handleChange('maxFileSize', parseInt(e.target.value))}
                                />

                                <Box>
                                    <Typography variant="body2" fontWeight="500" mb={1}>Allowed File Types</Typography>
                                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                                        {settings.allowedTypes.map((type) => (
                                            <Chip key={type} label={type} size="small" onDelete={() => { }} />
                                        ))}
                                    </Box>
                                </Box>

                                <Alert severity="warning">
                                    Thay đổi retention policy sẽ áp dụng cho các files mới. Files cũ không bị ảnh hưởng.
                                </Alert>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Storage Buckets */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                                <Storage color="info" />
                                <Typography variant="h6" fontWeight="600">
                                    Storage Buckets
                                </Typography>
                            </Box>

                            <TableContainer>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Bucket Name</TableCell>
                                            <TableCell>Type</TableCell>
                                            <TableCell align="right">Size</TableCell>
                                            <TableCell align="right">Objects</TableCell>
                                            <TableCell>Last Modified</TableCell>
                                            <TableCell align="right">Actions</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {storageBuckets.map((bucket) => (
                                            <TableRow key={bucket.name}>
                                                <TableCell>
                                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                        {getTypeIcon(bucket.type)}
                                                        <Typography fontWeight="500">{bucket.name}</Typography>
                                                    </Box>
                                                </TableCell>
                                                <TableCell>
                                                    <Chip label={bucket.type} size="small" variant="outlined" />
                                                </TableCell>
                                                <TableCell align="right">{bucket.size.toFixed(1)} GB</TableCell>
                                                <TableCell align="right">{bucket.objects.toLocaleString()}</TableCell>
                                                <TableCell>{bucket.lastModified}</TableCell>
                                                <TableCell align="right">
                                                    <Button size="small">Browse</Button>
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
