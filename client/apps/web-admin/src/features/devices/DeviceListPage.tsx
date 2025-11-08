import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Button,
    Card,
    Grid,
    Typography,
    Chip,
    IconButton,
    CardContent,
    CardActions,
    Tooltip,
    Snackbar,
    Alert,
    CircularProgress,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
} from '@mui/material';
import { apiClient } from '@face-attendance/utils';
import {
    Add,
    Edit,
    Settings as SettingsIcon,
    Refresh as RefreshIcon,
    ContentCopy as ContentCopyIcon,
    Delete as DeleteIcon,
    PowerSettingsNew as PowerIcon,
} from '@mui/icons-material';
import type { Device } from '@face-attendance/types';

export const DeviceListPage: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [devices, setDevices] = useState<Device[]>([]);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error' | 'info';
    }>({ open: false, message: '', severity: 'success' });
    const [refreshingId, setRefreshingId] = useState<string | null>(null);
    const [hoverTokenId, setHoverTokenId] = useState<string | null>(null);
    const [deleteDialog, setDeleteDialog] = useState<{
        open: boolean;
        deviceId: string | null;
        deviceName: string;
    }>({ open: false, deviceId: null, deviceName: '' });
    const [deletingId, setDeletingId] = useState<string | null>(null);
    const [togglingStatusId, setTogglingStatusId] = useState<string | null>(
        null
    );

    // Get devices from API
    useEffect(() => {
        const fetchDevices = async () => {
            try {
                const response = await apiClient.get('/api/v1/device');
                if (response.status === 200) {
                    console.log('Fetched devices:', response.data);
                }
                const settingConfig = {
                    allow_check_in: true,
                    allow_check_out: true,
                    timeout: 30,
                    recognition_threshold: 0.8,
                    sound_enabled: true,
                };
                const repsData = response.data.data.devices;
                if (!Array.isArray(repsData)) {
                    console.error('Invalid data format for devices:', repsData);
                    setDevices([]);
                    return;
                }
                const listDevices: Device[] = repsData.map((device: any) => ({
                    id: device.device_id,
                    company_id: device.company_id,
                    name: device.name,
                    serial_number: device.serial_number,
                    location: device.address,
                    status: device.status == 1 ? 'online' : 'offline',
                    settings: settingConfig,
                    created_at: device.created_at,
                    updated_at: device.update_at,
                    mac_address: device.mac_address,
                    token: device.token,
                }));
                setDevices(listDevices);
            } catch (error) {
                console.error('Failed to fetch stats:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchDevices();
    }, []);

    const handleCopy = (token: string | undefined) => {
        if (!token) {
            setSnackbar({
                open: true,
                message: 'Token không khả dụng',
                severity: 'error',
            });
            return;
        }
        navigator.clipboard
            .writeText(token)
            .then(() =>
                setSnackbar({
                    open: true,
                    message: 'Đã copy token',
                    severity: 'success',
                })
            )
            .catch(() =>
                setSnackbar({
                    open: true,
                    message: 'Copy thất bại',
                    severity: 'error',
                })
            );
    };

    const handleRefreshToken = async (deviceId: string) => {
        setRefreshingId(deviceId);
        try {
            // Assumption: API endpoint returns { data: { token: string } }
            const response = await apiClient.post(
                `/api/v1/device/token/refresh/${deviceId}`
            );
            const newToken = response?.data?.data?.device_token;
            if (newToken) {
                setDevices((prev) =>
                    prev.map((d) =>
                        d.id === deviceId ? { ...d, token: newToken } : d
                    )
                );
                setSnackbar({
                    open: true,
                    message: 'Token đã được làm mới',
                    severity: 'success',
                });
            } else {
                setSnackbar({
                    open: true,
                    message: 'Không nhận được token mới',
                    severity: 'error',
                });
            }
        } catch (e) {
            console.error('Refresh token error', e);
            setSnackbar({
                open: true,
                message: 'Làm mới token thất bại',
                severity: 'error',
            });
        } finally {
            setRefreshingId(null);
        }
    };

    const handleDeleteClick = (device: Device) => {
        setDeleteDialog({
            open: true,
            deviceId: device.id,
            deviceName: device.name,
        });
    };

    const handleDeleteConfirm = async () => {
        if (!deleteDialog.deviceId) return;
        setDeletingId(deleteDialog.deviceId);
        try {
            await apiClient.delete(`/api/v1/device/${deleteDialog.deviceId}`);
            setDevices((prev) =>
                prev.filter((d) => d.id !== deleteDialog.deviceId)
            );
            setSnackbar({
                open: true,
                message: 'Đã xóa thiết bị thành công',
                severity: 'success',
            });
        } catch (e) {
            console.error('Delete device error', e);
            setSnackbar({
                open: true,
                message: 'Xóa thiết bị thất bại',
                severity: 'error',
            });
        } finally {
            setDeletingId(null);
            setDeleteDialog({ open: false, deviceId: null, deviceName: '' });
        }
    };

    const handleDeleteCancel = () => {
        setDeleteDialog({ open: false, deviceId: null, deviceName: '' });
    };

    const handleToggleStatus = async (device: Device) => {
        setTogglingStatusId(device.id);
        const newStatus = device.status === 'online' ? 0 : 1; // 0=offline, 1=online
        try {
            await apiClient.post(`/api/v1/device/status`, {
                device_id: device.id,
                status: newStatus,
            });
            setDevices((prev) =>
                prev.map((d) =>
                    d.id === device.id
                        ? {
                              ...d,
                              status: newStatus === 1 ? 'online' : 'offline',
                          }
                        : d
                )
            );
            setSnackbar({
                open: true,
                message: `Đã chuyển trạng thái sang ${
                    newStatus === 1 ? 'Online' : 'Offline'
                }`,
                severity: 'success',
            });
        } catch (e) {
            console.error('Toggle status error', e);
            setSnackbar({
                open: true,
                message: 'Thay đổi trạng thái thất bại',
                severity: 'error',
            });
        } finally {
            setTogglingStatusId(null);
        }
    };

    const renderToken = (device: Device) => {
        const isHover = hoverTokenId === device.id;
        const token = device.token || '';
        const tokenLen = token.length;
        const display = isHover ? token : `*`.repeat(tokenLen);
        return (
            <Box
                mt={1}
                display="flex"
                alignItems="center"
                gap={1}
                onMouseEnter={() => setHoverTokenId(device.id)}
                onMouseLeave={() => setHoverTokenId(null)}
            >
                <Tooltip
                    title={
                        isHover ? 'Bấm để copy token' : 'Di chuột để hiện token'
                    }
                >
                    <Box
                        onClick={() => handleCopy(device.token)}
                        sx={{
                            cursor: 'pointer',
                            fontFamily: 'monospace',
                            backgroundColor: '#f5f5f5',
                            px: 1.5,
                            py: 0.5,
                            borderRadius: 1,
                            border: '1px solid #ddd',
                            maxWidth: '100%',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            fontSize: 12,
                        }}
                    >
                        {display}
                    </Box>
                </Tooltip>
                <Tooltip title="Làm mới token">
                    <span>
                        <IconButton
                            size="small"
                            onClick={() => handleRefreshToken(device.id)}
                            disabled={refreshingId === device.id}
                        >
                            {refreshingId === device.id ? (
                                <CircularProgress size={18} />
                            ) : (
                                <RefreshIcon fontSize="small" />
                            )}
                        </IconButton>
                    </span>
                </Tooltip>
                <Tooltip title="Copy token">
                    <IconButton
                        size="small"
                        onClick={() => handleCopy(device.token)}
                    >
                        <ContentCopyIcon fontSize="small" />
                    </IconButton>
                </Tooltip>
            </Box>
        );
    };

    if (loading) {
        return <Typography>Đang tải thiết bị...</Typography>;
    }

    return (
        <Box>
            <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                mb={3}
            >
                <Typography variant="h4" fontWeight="bold">
                    Quản lý Thiết bị
                </Typography>
                <Button
                    variant="contained"
                    startIcon={<Add />}
                    onClick={() => navigate('/devices/add')}
                >
                    Thêm thiết bị
                </Button>
            </Box>
            <Grid container spacing={3}>
                {devices.map((device) => (
                    <Grid item xs={12} sm={6} md={4} key={device.id}>
                        <Card>
                            <CardContent>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="start"
                                    mb={2}
                                >
                                    <Typography variant="h6" fontWeight="bold">
                                        {device.name}
                                    </Typography>
                                    <Box
                                        display="flex"
                                        gap={1}
                                        alignItems="center"
                                    >
                                        <Tooltip
                                            title={
                                                device.status === 'online'
                                                    ? 'Chuyển sang Offline'
                                                    : 'Chuyển sang Online'
                                            }
                                        >
                                            <span>
                                                <IconButton
                                                    size="small"
                                                    onClick={() =>
                                                        handleToggleStatus(
                                                            device
                                                        )
                                                    }
                                                    disabled={
                                                        togglingStatusId ===
                                                        device.id
                                                    }
                                                    color={
                                                        device.status ===
                                                        'online'
                                                            ? 'success'
                                                            : 'error'
                                                    }
                                                >
                                                    {togglingStatusId ===
                                                    device.id ? (
                                                        <CircularProgress
                                                            size={16}
                                                        />
                                                    ) : (
                                                        <PowerIcon fontSize="small" />
                                                    )}
                                                </IconButton>
                                            </span>
                                        </Tooltip>
                                        <Chip
                                            label={
                                                device.status === 'online'
                                                    ? 'Online'
                                                    : 'Offline'
                                            }
                                            color={
                                                device.status === 'online'
                                                    ? 'success'
                                                    : 'error'
                                            }
                                            size="small"
                                        />
                                    </Box>
                                </Box>
                                <Typography
                                    variant="body2"
                                    color="textSecondary"
                                    mb={1}
                                >
                                    Mã: {device.serial_number}
                                </Typography>
                                <Typography
                                    variant="body2"
                                    color="textSecondary"
                                    mb={1}
                                >
                                    Vị trí: {device.location}
                                </Typography>
                                <Typography
                                    variant="body2"
                                    color="textSecondary"
                                >
                                    Mac Address: {device.mac_address || 'N/A'}
                                </Typography>
                                <Typography
                                    mt={1}
                                    variant="caption"
                                    color="textSecondary"
                                    display="block"
                                >
                                    Token thiết bị:
                                </Typography>
                                {renderToken(device)}
                            </CardContent>
                            <CardActions>
                                <IconButton
                                    size="small"
                                    onClick={() =>
                                        navigate(`/devices/${device.id}/config`)
                                    }
                                >
                                    <SettingsIcon />
                                </IconButton>
                                <IconButton
                                    size="small"
                                    onClick={() =>
                                        navigate(`/devices/${device.id}/edit`, {
                                            state: { device },
                                        })
                                    }
                                >
                                    <Edit />
                                </IconButton>
                                <Tooltip title="Xóa thiết bị">
                                    <IconButton
                                        size="small"
                                        onClick={() =>
                                            handleDeleteClick(device)
                                        }
                                        color="error"
                                        disabled={deletingId === device.id}
                                    >
                                        {deletingId === device.id ? (
                                            <CircularProgress size={18} />
                                        ) : (
                                            <DeleteIcon />
                                        )}
                                    </IconButton>
                                </Tooltip>
                            </CardActions>
                        </Card>
                    </Grid>
                ))}
            </Grid>
            <Dialog
                open={deleteDialog.open}
                onClose={handleDeleteCancel}
                aria-labelledby="delete-dialog-title"
                aria-describedby="delete-dialog-description"
            >
                <DialogTitle id="delete-dialog-title">
                    Xác nhận xóa thiết bị
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="delete-dialog-description">
                        Bạn có chắc chắn muốn xóa thiết bị "
                        {deleteDialog.deviceName}"? Hành động này không thể hoàn
                        tác.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={handleDeleteCancel}
                        disabled={deletingId !== null}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleDeleteConfirm}
                        color="error"
                        variant="contained"
                        disabled={deletingId !== null}
                        autoFocus
                    >
                        {deletingId !== null ? 'Đang xóa...' : 'Xóa'}
                    </Button>
                </DialogActions>
            </Dialog>
            <Snackbar
                open={snackbar.open}
                autoHideDuration={3000}
                onClose={() => setSnackbar((s) => ({ ...s, open: false }))}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert
                    severity={snackbar.severity}
                    onClose={() => setSnackbar((s) => ({ ...s, open: false }))}
                    variant="filled"
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};
