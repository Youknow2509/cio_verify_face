import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    Alert,
    MenuItem,
    InputAdornment,
    IconButton,
} from '@mui/material';
import { Save, ArrowBack, Shuffle } from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';
import { DeviceType } from '@face-attendance/types';

interface DeviceAddFormData {
    address: string;
    device_name: string;
    device_type: DeviceType; // 0 default type
    mac_address: string;
    serial_number: string;
}

export const DeviceAddPage: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [formData, setFormData] = useState<DeviceAddFormData>({
        address: '',
        device_name: '',
        device_type: DeviceType.FACE_TERMINAL,
        mac_address: '',
        serial_number: '',
    });

    const handleChange =
        (field: keyof DeviceAddFormData) =>
        (e: React.ChangeEvent<HTMLInputElement>) => {
            const raw = e.target.value;
            const value =
                field === 'device_type' ? (Number(raw) as DeviceType) : raw;
            setFormData({ ...formData, [field]: value as any });
        };

    const randomHexPair = () =>
        Math.floor(Math.random() * 256)
            .toString(16)
            .padStart(2, '0')
            .toUpperCase();
    const generateRandomMac = () =>
        `${randomHexPair()}:${randomHexPair()}:${randomHexPair()}:${randomHexPair()}:${randomHexPair()}:${randomHexPair()}`;
    const generateRandomSerial = () =>
        `SN-${Math.random().toString(36).slice(2, 8).toUpperCase()}${Date.now()
            .toString(36)
            .slice(-4)
            .toUpperCase()}`;

    const handleRandomMac = () =>
        setFormData((s) => ({ ...s, mac_address: generateRandomMac() }));
    const handleRandomSerial = () =>
        setFormData((s) => ({ ...s, serial_number: generateRandomSerial() }));

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            // POST new device
            await apiClient.post('/api/v1/device', formData);
            navigate('/devices');
        } catch (err: any) {
            console.error('Create device error', err);
            setError(
                err.response?.data?.message || 'Không thể tạo thiết bị mới'
            );
        } finally {
            setLoading(false);
        }
    };

    const isValid =
        formData.address.trim() !== '' &&
        formData.device_name.trim() !== '' &&
        formData.serial_number.trim() !== '' &&
        formData.mac_address.trim() !== '';

    return (
        <Box>
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/devices')}
                sx={{ mb: 2 }}
            >
                Quay lại
            </Button>
            <Card>
                <CardContent>
                    <Typography variant="h5" fontWeight="bold" mb={3}>
                        Thêm thiết bị mới
                    </Typography>
                    {error && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {error}
                        </Alert>
                    )}
                    <Box component="form" onSubmit={handleSubmit}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Tên thiết bị"
                                    required
                                    value={formData.device_name}
                                    onChange={handleChange('device_name')}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Serial Number"
                                    required
                                    value={formData.serial_number}
                                    onChange={handleChange('serial_number')}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton
                                                    aria-label="random serial"
                                                    onClick={handleRandomSerial}
                                                    edge="end"
                                                >
                                                    <Shuffle />
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Địa chỉ MAC"
                                    required
                                    value={formData.mac_address}
                                    onChange={handleChange('mac_address')}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton
                                                    aria-label="random mac"
                                                    onClick={handleRandomMac}
                                                    edge="end"
                                                >
                                                    <Shuffle />
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    select
                                    fullWidth
                                    label="Loại thiết bị"
                                    required
                                    value={formData.device_type}
                                    onChange={handleChange('device_type')}
                                >
                                    <MenuItem value={0}>FACE_TERMINAL</MenuItem>
                                    <MenuItem value={1}>MOBILE_APP</MenuItem>
                                    <MenuItem value={2}>WEB_CAMERA</MenuItem>
                                    <MenuItem value={3}>IOT_SENSOR</MenuItem>
                                </TextField>
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    label="Địa chỉ lắp đặt"
                                    required
                                    value={formData.address}
                                    onChange={handleChange('address')}
                                />
                            </Grid>
                        </Grid>
                        <Box mt={3} display="flex" gap={2}>
                            <Button
                                type="submit"
                                variant="contained"
                                startIcon={<Save />}
                                disabled={!isValid || loading}
                            >
                                {loading ? 'Đang lưu...' : 'Lưu'}
                            </Button>
                            <Button
                                variant="outlined"
                                onClick={() => navigate('/devices')}
                                disabled={loading}
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
