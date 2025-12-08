import { useEffect, useRef, useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    TextField,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    CircularProgress,
    Alert,
    Chip,
    Grid,
    Avatar,
    InputAdornment,
    IconButton,
} from '@mui/material';
import {
    Person as PersonIcon,
    Edit as EditIcon,
    PhotoCamera as PhotoCameraIcon,
    CheckCircle as CheckCircleIcon,
    Pending as PendingIcon,
    Cancel as CancelIcon,
    ContentCopy as ContentCopyIcon,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';
import { profileUpdateApi } from '@/services/api';

export const ProfilePage: React.FC = () => {
    const { user } = useAuthStore();
    const [loading, setLoading] = useState(false);
    const [requestStatus, setRequestStatus] = useState<any>(null);
    const [requestHistory, setRequestHistory] = useState<any[]>([]);
    const [openRequestDialog, setOpenRequestDialog] = useState(false);
    const [openUploadDialog, setOpenUploadDialog] = useState(false);
    const [reason, setReason] = useState('');
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [updateToken, setUpdateToken] = useState('');
    const [capturing, setCapturing] = useState(false);
    const [capturedBlob, setCapturedBlob] = useState<Blob | null>(null);
    const [capturePreviewUrl, setCapturePreviewUrl] = useState('');
    const [cameraStream, setCameraStream] = useState<MediaStream | null>(null);
    const videoRef = useRef<HTMLVideoElement | null>(null);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [showToken, setShowToken] = useState(false);

    useEffect(() => {
        loadRequestStatus();
    }, []);

    const loadRequestStatus = async () => {
        try {
            const response: any = await profileUpdateApi.getMyRequest();
            const parsed = response?.data ?? response; // interceptor already unwraps
            const data = parsed?.data ?? parsed ?? null;

            const requestsRaw = Array.isArray(data)
                ? data
                : Array.isArray(data?.requests)
                ? data.requests
                : data
                ? [data]
                : [];

            const requests = requestsRaw.map((req) => {
                const tokenFromLink = extractToken(
                    req?.update_link || req?.updateLink
                );
                return {
                    ...req,
                    token: req?.token || tokenFromLink || '',
                };
            });

            setRequestHistory(requests);
            const latest = requests[0] ?? null;
            setRequestStatus(latest);

            const tokenFromLatest = extractToken(
                latest?.update_link || latest?.updateLink
            );
            if (tokenFromLatest) setUpdateToken(tokenFromLatest);
        } catch (error: any) {
            console.error('Failed to load request status:', error);
            setRequestStatus(null);
            setRequestHistory([]);
        }
    };

    const handleCreateRequest = async () => {
        if (!reason.trim()) {
            setError('Vui lòng nhập lý do yêu cầu');
            return;
        }

        try {
            setLoading(true);
            setError('');
            const res: any = await profileUpdateApi.createRequest({ reason });
            const parsed = res?.data ?? res;
            setSuccess(parsed?.message || 'Tạo yêu cầu thành công!');
            setOpenRequestDialog(false);
            setReason('');
            await loadRequestStatus();
        } catch (err: any) {
            setError(err.message || 'Không thể tạo yêu cầu. Vui lòng thử lại.');
        } finally {
            setLoading(false);
        }
    };

    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            if (file.size > 10 * 1024 * 1024) {
                setError('Kích thước file phải nhỏ hơn 10MB');
                return;
            }
            setSelectedFile(file);
            setCapturedBlob(null);
            if (capturePreviewUrl) URL.revokeObjectURL(capturePreviewUrl);
            setCapturePreviewUrl('');
        }
    };

    const handleUploadFace = async () => {
        const fileToSend =
            selectedFile ||
            (capturedBlob
                ? new File([capturedBlob], 'capture.jpg', {
                      type: capturedBlob.type || 'image/jpeg',
                  })
                : null);

        if (!fileToSend || !updateToken.trim()) {
            setError('Vui lòng chọn ảnh (hoặc chụp từ camera) và nhập token');
            return;
        }

        if (fileToSend.size > 10 * 1024 * 1024) {
            setError('Kích thước file phải nhỏ hơn 10MB');
            return;
        }

        try {
            setLoading(true);
            setError('');

            const formData = new FormData();
            formData.append('token', updateToken);
            formData.append('image', fileToSend);

            const res: any = await profileUpdateApi.updateFace(formData);
            const parsed = res?.data ?? res;
            setSuccess(parsed?.message || 'Cập nhật khuôn mặt thành công!');
            setOpenUploadDialog(false);
            setSelectedFile(null);
            setUpdateToken('');
            setCapturedBlob(null);
            if (capturePreviewUrl) URL.revokeObjectURL(capturePreviewUrl);
            setCapturePreviewUrl('');
            stopCamera();
            await loadRequestStatus();
        } catch (err: any) {
            setError(
                err.message || 'Không thể cập nhật khuôn mặt. Vui lòng thử lại.'
            );
        } finally {
            setLoading(false);
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'pending':
                return {
                    label: 'Đang chờ',
                    color: 'warning' as const,
                    icon: <PendingIcon />,
                };
            case 'approved':
                return {
                    label: 'Đã duyệt',
                    color: 'success' as const,
                    icon: <CheckCircleIcon />,
                };
            case 'rejected':
                return {
                    label: 'Từ chối',
                    color: 'error' as const,
                    icon: <CancelIcon />,
                };
            default:
                return {
                    label: status || 'Không rõ',
                    color: 'default' as const,
                    icon: <PendingIcon />,
                };
        }
    };

    const extractToken = (link?: string) => {
        if (!link) return '';
        try {
            const url = new URL(link);
            return url.searchParams.get('token') || '';
        } catch (_err) {
            return '';
        }
    };

    const stopCamera = () => {
        cameraStream?.getTracks().forEach((track) => track.stop());
        if (videoRef.current) {
            videoRef.current.srcObject = null;
        }
        setCameraStream(null);
        setCapturing(false);
    };

    const startCamera = async () => {
        try {
            setError('');
            const stream = await navigator.mediaDevices.getUserMedia({
                video: { facingMode: 'user' },
            });
            if (videoRef.current) {
                videoRef.current.srcObject = stream;
                try {
                    await videoRef.current.play(); // ensure autoplay starts
                } catch (_err) {
                    /* ignore */
                }
            }
            setCameraStream(stream);
            setCapturing(true);
        } catch (err: any) {
            setError('Không thể bật camera. Vui lòng kiểm tra quyền truy cập.');
        }
    };

    const capturePhoto = () => {
        if (!videoRef.current) return;
        const video = videoRef.current;
        if (!video.videoWidth || !video.videoHeight) {
            setError('Camera chưa sẵn sàng. Vui lòng thử lại.');
            return;
        }
        const canvas = document.createElement('canvas');
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        const ctx = canvas.getContext('2d');
        if (!ctx) return;
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
        canvas.toBlob(
            (blob) => {
                if (blob) {
                    setCapturedBlob(blob);
                    setSelectedFile(null);
                    if (capturePreviewUrl)
                        URL.revokeObjectURL(capturePreviewUrl);
                    setCapturePreviewUrl(URL.createObjectURL(blob));
                }
            },
            'image/jpeg',
            0.9
        );
    };

    const handleCopyToken = async (value: string) => {
        if (!value) return;
        try {
            await navigator.clipboard.writeText(value);
            setSuccess('Đã copy token');
        } catch (_err) {
            setError('Không thể copy token');
        }
    };

    useEffect(() => {
        return () => {
            stopCamera();
            if (capturePreviewUrl) URL.revokeObjectURL(capturePreviewUrl);
        };
    }, [capturePreviewUrl]);

    useEffect(() => {
        if (cameraStream && videoRef.current) {
            videoRef.current.srcObject = cameraStream;
            videoRef.current.play().catch(() => undefined);
        }
    }, [cameraStream]);

    return (
        <Box>
            <Typography variant="h4" fontWeight="700" mb={3}>
                Hồ sơ cá nhân
            </Typography>

            {error && (
                <Alert
                    severity="error"
                    sx={{ mb: 3 }}
                    onClose={() => setError('')}
                >
                    {error}
                </Alert>
            )}

            {success && (
                <Alert
                    severity="success"
                    sx={{ mb: 3 }}
                    onClose={() => setSuccess('')}
                >
                    {success}
                </Alert>
            )}

            <Grid container spacing={3}>
                {/* User Info Card */}
                <Grid item xs={12} md={6}>
                    <Card>
                        <CardContent>
                            <Box
                                sx={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 2,
                                    mb: 3,
                                }}
                            >
                                <Avatar
                                    sx={{
                                        width: 80,
                                        height: 80,
                                        background:
                                            'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                                        fontSize: '2rem',
                                    }}
                                >
                                    {user?.name?.charAt(0) ||
                                        user?.email?.charAt(0) ||
                                        'U'}
                                </Avatar>
                                <Box>
                                    <Typography variant="h5" fontWeight="600">
                                        {user?.name || 'Người dùng'}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        color="text.secondary"
                                    >
                                        {user?.email}
                                    </Typography>
                                </Box>
                            </Box>

                            <Box sx={{ mb: 2 }}>
                                <Typography
                                    variant="body2"
                                    color="text.secondary"
                                    gutterBottom
                                >
                                    Mã nhân viên
                                </Typography>
                                <Typography variant="body1">
                                    {user?.employee_code || 'Chưa có'}
                                </Typography>
                            </Box>

                            <Box>
                                <Typography
                                    variant="body2"
                                    color="text.secondary"
                                    gutterBottom
                                >
                                    Vai trò
                                </Typography>
                                <Typography variant="body1">
                                    Nhân viên
                                </Typography>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Face Update Request Card */}
                <Grid item xs={12} md={6}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" fontWeight="600" mb={2}>
                                Cập nhật khuôn mặt
                            </Typography>

                            {requestStatus ? (
                                <Box>
                                    <Box sx={{ mb: 2 }}>
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                            gutterBottom
                                        >
                                            Trạng thái yêu cầu
                                        </Typography>
                                        <Chip
                                            label={
                                                getStatusLabel(
                                                    requestStatus.status
                                                ).label
                                            }
                                            color={
                                                getStatusLabel(
                                                    requestStatus.status
                                                ).color
                                            }
                                            icon={
                                                getStatusLabel(
                                                    requestStatus.status
                                                ).icon
                                            }
                                        />
                                    </Box>

                                    <Box sx={{ mb: 2 }}>
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                            gutterBottom
                                        >
                                            Lý do
                                        </Typography>
                                        <Typography variant="body2">
                                            {requestStatus.reason}
                                        </Typography>
                                    </Box>

                                    {requestStatus.status === 'approved' &&
                                        requestStatus.token && (
                                            <Box
                                                sx={{ mb: 2 }}
                                                onMouseEnter={() =>
                                                    setShowToken(true)
                                                }
                                                onMouseLeave={() =>
                                                    setShowToken(false)
                                                }
                                            >
                                                <Typography
                                                    variant="body2"
                                                    color="text.secondary"
                                                    gutterBottom
                                                >
                                                    Token xác thực (di chuột để
                                                    xem)
                                                </Typography>
                                                <TextField
                                                    fullWidth
                                                    value={
                                                        showToken
                                                            ? requestStatus.token
                                                            : '••••••••••••••••'
                                                    }
                                                    InputProps={{
                                                        readOnly: true,
                                                        endAdornment: (
                                                            <InputAdornment position="end">
                                                                <IconButton
                                                                    size="small"
                                                                    onClick={() =>
                                                                        handleCopyToken(
                                                                            requestStatus.token
                                                                        )
                                                                    }
                                                                >
                                                                    <ContentCopyIcon fontSize="small" />
                                                                </IconButton>
                                                            </InputAdornment>
                                                        ),
                                                    }}
                                                    size="small"
                                                />
                                            </Box>
                                        )}

                                    {requestStatus.status === 'approved' &&
                                        requestStatus.token && (
                                            <Button
                                                variant="contained"
                                                startIcon={<PhotoCameraIcon />}
                                                onClick={() =>
                                                    setOpenUploadDialog(true)
                                                }
                                                fullWidth
                                            >
                                                Upload khuôn mặt mới
                                            </Button>
                                        )}

                                    {requestStatus.status !== 'pending' && (
                                        <Button
                                            variant="outlined"
                                            startIcon={<EditIcon />}
                                            onClick={() =>
                                                setOpenRequestDialog(true)
                                            }
                                            fullWidth
                                            sx={{ mt: 2 }}
                                        >
                                            Tạo yêu cầu khác
                                        </Button>
                                    )}

                                    <Box sx={{ mt: 3 }}>
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                            gutterBottom
                                        >
                                            Lịch sử yêu cầu
                                        </Typography>
                                        {requestHistory.length === 0 ? (
                                            <Typography
                                                variant="body2"
                                                color="text.secondary"
                                            >
                                                Chưa có lịch sử yêu cầu.
                                            </Typography>
                                        ) : (
                                            <Box
                                                sx={{
                                                    display: 'flex',
                                                    flexDirection: 'column',
                                                    gap: 1,
                                                    maxHeight: 240,
                                                    overflowY: 'auto',
                                                }}
                                            >
                                                {requestHistory.map(
                                                    (req, idx) => {
                                                        const statusMeta =
                                                            getStatusLabel(
                                                                req.status
                                                            );
                                                        return (
                                                            <Card
                                                                key={
                                                                    req.id ||
                                                                    idx
                                                                }
                                                                variant="outlined"
                                                                sx={{
                                                                    borderRadius: 2,
                                                                }}
                                                            >
                                                                <CardContent
                                                                    sx={{
                                                                        py: 1.5,
                                                                    }}
                                                                >
                                                                    <Box
                                                                        sx={{
                                                                            display:
                                                                                'flex',
                                                                            justifyContent:
                                                                                'space-between',
                                                                            alignItems:
                                                                                'center',
                                                                            gap: 2,
                                                                        }}
                                                                    >
                                                                        <Box>
                                                                            <Typography
                                                                                variant="body2"
                                                                                fontWeight={
                                                                                    600
                                                                                }
                                                                            >
                                                                                {req.reason ||
                                                                                    'Không có lý do'}
                                                                            </Typography>
                                                                            <Typography
                                                                                variant="caption"
                                                                                color="text.secondary"
                                                                            >
                                                                                {req.created_at ||
                                                                                    req.createdAt ||
                                                                                    'Chưa rõ thời gian'}
                                                                            </Typography>
                                                                        </Box>
                                                                        <Chip
                                                                            label={
                                                                                statusMeta.label
                                                                            }
                                                                            color={
                                                                                statusMeta.color
                                                                            }
                                                                            icon={
                                                                                statusMeta.icon
                                                                            }
                                                                            size="small"
                                                                        />
                                                                    </Box>
                                                                </CardContent>
                                                            </Card>
                                                        );
                                                    }
                                                )}
                                            </Box>
                                        )}
                                    </Box>
                                </Box>
                            ) : (
                                <Box>
                                    <Typography
                                        variant="body2"
                                        color="text.secondary"
                                        mb={2}
                                    >
                                        Bạn chưa có yêu cầu cập nhật khuôn mặt
                                        nào
                                    </Typography>
                                    <Button
                                        variant="contained"
                                        startIcon={<EditIcon />}
                                        onClick={() =>
                                            setOpenRequestDialog(true)
                                        }
                                        fullWidth
                                    >
                                        Tạo yêu cầu mới
                                    </Button>
                                </Box>
                            )}
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>

            {/* Create Request Dialog */}
            <Dialog
                open={openRequestDialog}
                onClose={() => setOpenRequestDialog(false)}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>Tạo yêu cầu cập nhật khuôn mặt</DialogTitle>
                <DialogContent>
                    <TextField
                        fullWidth
                        multiline
                        rows={4}
                        label="Lý do yêu cầu"
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        placeholder="Ví dụ: Tôi muốn đổi khuôn mặt vì đã thay đổi diện mạo"
                        sx={{ mt: 2 }}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenRequestDialog(false)}>
                        Hủy
                    </Button>
                    <Button
                        onClick={handleCreateRequest}
                        variant="contained"
                        disabled={loading}
                    >
                        {loading ? (
                            <CircularProgress size={24} />
                        ) : (
                            'Gửi yêu cầu'
                        )}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Upload Face Dialog */}
            <Dialog
                open={openUploadDialog}
                onClose={() => setOpenUploadDialog(false)}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>Upload khuôn mặt mới</DialogTitle>
                <DialogContent>
                    <TextField
                        fullWidth
                        label="Token xác thực"
                        value={updateToken}
                        onChange={(e) => setUpdateToken(e.target.value)}
                        placeholder="Nhập token từ email"
                        sx={{ mt: 2, mb: 2 }}
                    />
                    <Button
                        variant="outlined"
                        fullWidth
                        sx={{ mb: 2 }}
                        onClick={capturing ? stopCamera : startCamera}
                    >
                        {capturing ? 'Tắt camera' : 'Bật camera (live)'}
                    </Button>
                    {capturing && (
                        <Box sx={{ mb: 2 }}>
                            <video
                                ref={videoRef}
                                autoPlay
                                playsInline
                                muted
                                style={{ width: '100%', borderRadius: 8 }}
                            />
                            <Button
                                variant="contained"
                                fullWidth
                                sx={{ mt: 1 }}
                                onClick={capturePhoto}
                            >
                                Chụp ảnh từ camera
                            </Button>
                        </Box>
                    )}
                    {capturePreviewUrl && !selectedFile && (
                        <Box sx={{ mb: 2 }}>
                            <Typography
                                variant="body2"
                                color="text.secondary"
                                gutterBottom
                            >
                                Ảnh đã chụp
                            </Typography>
                            <Box
                                component="img"
                                src={capturePreviewUrl}
                                alt="Xem trước ảnh chụp"
                                sx={{
                                    width: '100%',
                                    borderRadius: 2,
                                    border: '1px solid rgba(255,255,255,0.1)',
                                }}
                            />
                        </Box>
                    )}
                    <Button
                        variant="outlined"
                        component="label"
                        fullWidth
                        startIcon={<PhotoCameraIcon />}
                    >
                        {selectedFile
                            ? selectedFile.name
                            : 'Chọn ảnh khuôn mặt'}
                        <input
                            type="file"
                            hidden
                            accept="image/*"
                            onChange={handleFileSelect}
                        />
                    </Button>
                    {selectedFile && (
                        <Typography
                            variant="caption"
                            color="text.secondary"
                            display="block"
                            mt={1}
                        >
                            Kích thước: {(selectedFile.size / 1024).toFixed(2)}{' '}
                            KB
                        </Typography>
                    )}
                    {capturedBlob && !selectedFile && (
                        <Typography
                            variant="caption"
                            color="text.secondary"
                            display="block"
                            mt={1}
                        >
                            Đã chụp ảnh từ camera
                        </Typography>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenUploadDialog(false)}>
                        Hủy
                    </Button>
                    <Button
                        onClick={handleUploadFace}
                        variant="contained"
                        disabled={loading}
                    >
                        {loading ? <CircularProgress size={24} /> : 'Upload'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
