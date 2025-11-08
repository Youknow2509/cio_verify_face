import { useState, useEffect } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    CircularProgress,
    Alert,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';
import { createApiClient } from '@face-attendance/utils';
import {
    DeviceEditNameForm,
    DeviceEditLocationForm,
} from '@face-attendance/types';

const apiClient = createApiClient();

interface DeviceEditShow {
    device_id: string;
    company_id: string;
    location_id: string;
    name: string;
    address: string;
    serial_number: string;
    mac_address?: string;
}

// Page for editing existing device (mapped at /devices/:id/edit)
export const DeviceFormPage: React.FC = () => {
    const navigate = useNavigate();
    const { state } = useLocation();
    const { id } = useParams();
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [formData, setFormData] = useState<DeviceEditShow>({
        device_id: state?.device?.device_id || id || '',
        company_id: state?.device?.company_id || '',
        location_id: state?.device?.location_id || '',
        name: state?.device?.name || '',
        address: state?.device?.location || '',
        serial_number: state?.device?.serial_number || '',
        mac_address: state?.device?.mac_address || '',
    });
    // Two editable sub-forms (name + location)
    const [nameForm, setNameForm] = useState<DeviceEditNameForm>({
        device_id: state?.device?.device_id || (id as string) || '',
        device_name: state?.device?.name || '',
    });
    const [locationForm, setLocationForm] = useState<DeviceEditLocationForm>({
        device_id: state?.device?.device_id || (id as string) || '',
        location_id: state?.device?.location_id,
        address: state?.device?.location || '',
    });

    // Keep initial snapshots to detect changes on submit
    const [initialNameForm, setInitialNameForm] = useState<DeviceEditNameForm>({
        device_id: state?.device?.device_id || (id as string) || '',
        device_name: state?.device?.name || '',
    });
    const [initialLocationForm, setInitialLocationForm] =
        useState<DeviceEditLocationForm>({
            device_id: state?.device?.device_id || (id as string) || '',
            location_id: state?.device?.location_id,
            address: state?.device?.location || '',
        });

    // Fetch device data from API if editing and no state data
    useEffect(() => {
        const fetchDeviceData = async () => {
            if (id && !state) {
                setLoading(true);
                setError(null);
                try {
                    const response = await apiClient.get(
                        `/api/v1/devices/${id}`
                    );
                    const deviceData = response.data?.device as any;

                    const mapped: DeviceEditShow = {
                        device_id: deviceData.device_id || deviceData.id || '',
                        company_id: deviceData.company_id || '',
                        location_id: deviceData.location_id || '',
                        name: deviceData.name || '',
                        address:
                            deviceData.address || deviceData.location || '',
                        serial_number: deviceData.serial_number || '',
                        mac_address: deviceData.mac_address || '',
                    };
                    setFormData(mapped);

                    const nextName: DeviceEditNameForm = {
                        device_id: mapped.device_id,
                        device_name: mapped.name,
                    };
                    const nextLoc: DeviceEditLocationForm = {
                        device_id: mapped.device_id,
                        location_id: mapped.location_id || undefined,
                        address: mapped.address,
                    };
                    setNameForm(nextName);
                    setLocationForm(nextLoc);
                    setInitialNameForm(nextName);
                    setInitialLocationForm(nextLoc);
                } catch (err: any) {
                    console.error('Error fetching device:', err);
                    setError(
                        err.response?.data?.message ||
                            'Không thể tải dữ liệu thiết bị'
                    );
                } finally {
                    setLoading(false);
                }
            }
        };

        fetchDeviceData();
    }, [id, state]);

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNameForm({ ...nameForm, device_name: e.target.value });
    };
    const handleLocationChange =
        (field: 'address' | 'location_id') =>
        (e: React.ChangeEvent<HTMLInputElement>) => {
            setLocationForm({ ...locationForm, [field]: e.target.value });
        };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const promises: Promise<any>[] = [];

            const nameChanged =
                nameForm.device_name !== initialNameForm.device_name;
            const locationChanged =
                locationForm.address !== initialLocationForm.address ||
                (locationForm.location_id || '') !==
                    (initialLocationForm.location_id || '');

            if (!formData.device_id) {
                setError('Thiếu device_id, không thể cập nhật.');
                return;
            }

            if (nameChanged) {
                promises.push(
                    apiClient.post(`/api/v1/device/name`, {
                        device_id: formData.device_id,
                        device_name: nameForm.device_name,
                    })
                );
            }
            if (locationChanged) {
                promises.push(
                    apiClient.post(`/api/v1/device/location`, {
                        device_id: formData.device_id,
                        location_id: locationForm.location_id,
                        address: locationForm.address,
                    })
                );
            }

            if (promises.length === 0) {
                navigate('/devices');
                return;
            }

            await Promise.all(promises);
            navigate('/devices');
        } catch (err: any) {
            console.error('Error saving device:', err);
            setError(err.response?.data?.message || 'Không thể lưu thiết bị');
        } finally {
            setLoading(false);
        }
    };

    if (loading && !formData.name) {
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="400px"
            >
                <CircularProgress />
            </Box>
        );
    }

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
                        Chỉnh sửa thiết bị
                    </Typography>

                    {error && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {error}
                        </Alert>
                    )}

                    <Box component="form" onSubmit={handleSubmit}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Device ID"
                                    value={formData.device_id}
                                    disabled
                                />
                            </Grid>
                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Company ID"
                                    value={formData.company_id}
                                    disabled
                                />
                            </Grid>
                            <Grid item xs={12} md={4}>
                                <TextField
                                    fullWidth
                                    label="Serial number"
                                    value={formData.serial_number}
                                    disabled
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Địa chỉ MAC"
                                    value={formData.mac_address || ''}
                                    disabled
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    label="Tên thiết bị"
                                    required
                                    value={nameForm.device_name}
                                    onChange={handleNameChange}
                                />
                            </Grid>
                            <Grid item xs={12} md={8}>
                                <TextField
                                    fullWidth
                                    label="Địa chỉ"
                                    required
                                    value={locationForm.address}
                                    onChange={handleLocationChange('address')}
                                />
                            </Grid>
                        </Grid>
                        <Box mt={3} display="flex" gap={2}>
                            <Button
                                type="submit"
                                variant="contained"
                                startIcon={<Save />}
                                disabled={
                                    loading ||
                                    !formData.device_id ||
                                    (nameForm.device_name ===
                                        initialNameForm.device_name &&
                                        locationForm.address ===
                                            initialLocationForm.address &&
                                        (locationForm.location_id || '') ===
                                            (initialLocationForm.location_id ||
                                                ''))
                                }
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
