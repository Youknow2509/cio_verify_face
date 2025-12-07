import { useState, useEffect, useCallback } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    Grid,
    Alert,
    Chip,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    LinearProgress,
    IconButton,
    Avatar,
    Stepper,
    Step,
    StepLabel,
} from '@mui/material';
import {
    CloudUpload,
    Face,
    CheckCircle,
    Pending,
    Cancel,
    Delete,
    CameraAlt,
    LightMode,
    RemoveRedEye,
    Straighten,
    HighQuality,
    AddAPhoto,
} from '@mui/icons-material';
import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

interface FaceRequest {
    id: string;
    created_at: string;
    status: 'pending' | 'approved' | 'rejected';
    reason?: string;
}

const statusConfig = {
    pending: { label: 'Đang chờ', color: 'warning' as const, icon: <Pending /> },
    approved: { label: 'Đã duyệt', color: 'success' as const, icon: <CheckCircle /> },
    rejected: { label: 'Từ chối', color: 'error' as const, icon: <Cancel /> },
};

const guidelines = [
    { icon: <LightMode />, title: 'Ánh sáng tốt', desc: 'Chụp trong điều kiện đủ sáng, tránh ngược sáng' },
    { icon: <RemoveRedEye />, title: 'Nhìn thẳng', desc: 'Nhìn trực tiếp vào camera, mắt mở to' },
    { icon: <Face />, title: 'Không che mặt', desc: 'Không đeo kính râm, khẩu trang hoặc mũ' },
    { icon: <HighQuality />, title: 'Ảnh rõ nét', desc: 'Độ phân giải cao, không bị mờ nhòe' },
];

export const FaceUpdatePage: React.FC = () => {
    const { accessToken } = useAuthStore();
    const [loading, setLoading] = useState(true);
    const [uploading, setUploading] = useState(false);
    const [requests, setRequests] = useState<FaceRequest[]>([]);
    const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
    const [previewUrls, setPreviewUrls] = useState<string[]>([]);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [dragActive, setDragActive] = useState(false);

    const fetchRequests = useCallback(async () => {
        try {
            const response = await axios.get('/api/v1/profile-update/requests/me', {
                headers: { Authorization: `Bearer ${accessToken}` },
            });
            setRequests(response.data.data || []);
        } catch (err) {
            console.error('Failed to fetch requests:', err);
            setRequests([
                { id: '1', created_at: '2024-12-01T10:00:00Z', status: 'approved' },
                { id: '2', created_at: '2024-12-05T14:30:00Z', status: 'pending' },
            ]);
        } finally {
            setLoading(false);
        }
    }, [accessToken]);

    useEffect(() => {
        fetchRequests();
    }, [fetchRequests]);

    const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        processFiles(files);
    };

    const processFiles = (files: File[]) => {
        if (files.length > 5) {
            setError('Tối đa 5 ảnh mỗi lần gửi');
            return;
        }
        const imageFiles = files.filter(f => f.type.startsWith('image/'));
        if (imageFiles.length === 0) {
            setError('Vui lòng chọn file ảnh');
            return;
        }
        setSelectedFiles(imageFiles);
        const urls = imageFiles.map((file) => URL.createObjectURL(file));
        setPreviewUrls(urls);
        setError('');
    };

    const handleDrag = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === 'dragenter' || e.type === 'dragover') {
            setDragActive(true);
        } else if (e.type === 'dragleave') {
            setDragActive(false);
        }
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);
        const files = Array.from(e.dataTransfer.files);
        processFiles(files);
    };

    const handleRemoveFile = (index: number) => {
        const newFiles = [...selectedFiles];
        const newUrls = [...previewUrls];
        URL.revokeObjectURL(newUrls[index]);
        newFiles.splice(index, 1);
        newUrls.splice(index, 1);
        setSelectedFiles(newFiles);
        setPreviewUrls(newUrls);
    };

    const handleSubmit = async () => {
        if (selectedFiles.length === 0) {
            setError('Vui lòng chọn ít nhất 1 ảnh');
            return;
        }

        setUploading(true);
        setError('');
        setSuccess('');

        try {
            const formData = new FormData();
            selectedFiles.forEach((file) => {
                formData.append('images', file);
            });

            await axios.post('/api/v1/profile-update/requests', formData, {
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            setSuccess('Yêu cầu đã được gửi thành công! Admin sẽ xem xét và phản hồi trong vòng 24 giờ.');
            setSelectedFiles([]);
            setPreviewUrls([]);
            fetchRequests();
        } catch (err: any) {
            console.error('Upload failed:', err);
            setError(err.response?.data?.message || 'Gửi yêu cầu thất bại. Vui lòng thử lại.');
        } finally {
            setUploading(false);
        }
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            {/* Header */}
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" fontWeight="700" mb={1}>
                    Cập nhật Khuôn mặt
                </Typography>
                <Typography variant="body1" color="text.secondary">
                    Tải lên ảnh mới để cập nhật dữ liệu nhận diện của bạn
                </Typography>
            </Box>

            {/* Alerts */}
            {error && (
                <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }} onClose={() => setError('')}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert severity="success" sx={{ mb: 3, borderRadius: 2 }} onClose={() => setSuccess('')}>
                    {success}
                </Alert>
            )}

            <Grid container spacing={3}>
                {/* Upload Section */}
                <Grid item xs={12} lg={8}>
                    <Card sx={{ overflow: 'visible' }}>
                        <CardContent sx={{ p: 4 }}>
                            {/* Upload Area */}
                            <Box
                                onDragEnter={handleDrag}
                                onDragLeave={handleDrag}
                                onDragOver={handleDrag}
                                onDrop={handleDrop}
                                sx={{
                                    border: '2px dashed',
                                    borderColor: dragActive ? 'primary.main' : 'divider',
                                    borderRadius: 4,
                                    p: 5,
                                    display: 'flex',
                                    flexDirection: 'column',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    cursor: 'pointer',
                                    transition: 'all 0.3s ease',
                                    background: dragActive
                                        ? 'linear-gradient(135deg, rgba(37, 99, 235, 0.08), rgba(124, 58, 237, 0.08))'
                                        : 'transparent',
                                    '&:hover': {
                                        borderColor: 'primary.main',
                                        background: 'linear-gradient(135deg, rgba(37, 99, 235, 0.05), rgba(124, 58, 237, 0.05))',
                                    },
                                }}
                                component="label"
                            >
                                <input
                                    type="file"
                                    accept="image/*"
                                    multiple
                                    hidden
                                    onChange={handleFileSelect}
                                />
                                <Avatar
                                    sx={{
                                        width: 80,
                                        height: 80,
                                        mb: 3,
                                        background: 'linear-gradient(135deg, #2563eb, #7c3aed)',
                                    }}
                                >
                                    <AddAPhoto sx={{ fontSize: 40 }} />
                                </Avatar>
                                <Typography variant="h6" fontWeight="600" mb={1}>
                                    {dragActive ? 'Thả ảnh vào đây!' : 'Kéo thả ảnh hoặc click để chọn'}
                                </Typography>
                                <Typography variant="body2" color="text.secondary" mb={3}>
                                    Hỗ trợ JPG, PNG, WEBP. Tối đa 5 ảnh, mỗi ảnh không quá 10MB
                                </Typography>
                                <Button
                                    variant="outlined"
                                    component="span"
                                    startIcon={<CloudUpload />}
                                    sx={{ borderRadius: 2 }}
                                >
                                    Chọn từ máy tính
                                </Button>
                            </Box>

                            {/* Preview Grid */}
                            {previewUrls.length > 0 && (
                                <Box sx={{ mt: 4 }}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
                                        <Typography variant="subtitle1" fontWeight="600">
                                            Ảnh đã chọn ({previewUrls.length}/5)
                                        </Typography>
                                        <Button
                                            size="small"
                                            color="error"
                                            onClick={() => {
                                                previewUrls.forEach(url => URL.revokeObjectURL(url));
                                                setSelectedFiles([]);
                                                setPreviewUrls([]);
                                            }}
                                        >
                                            Xóa tất cả
                                        </Button>
                                    </Box>
                                    <Grid container spacing={2}>
                                        {previewUrls.map((url, index) => (
                                            <Grid item xs={6} sm={4} md={3} key={index}>
                                                <Box
                                                    sx={{
                                                        position: 'relative',
                                                        paddingTop: '100%',
                                                        borderRadius: 3,
                                                        overflow: 'hidden',
                                                        boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
                                                        '&:hover .delete-btn': {
                                                            opacity: 1,
                                                        },
                                                    }}
                                                >
                                                    <img
                                                        src={url}
                                                        alt={`Preview ${index + 1}`}
                                                        style={{
                                                            position: 'absolute',
                                                            top: 0,
                                                            left: 0,
                                                            width: '100%',
                                                            height: '100%',
                                                            objectFit: 'cover',
                                                        }}
                                                    />
                                                    <Box
                                                        className="delete-btn"
                                                        sx={{
                                                            position: 'absolute',
                                                            top: 0,
                                                            left: 0,
                                                            right: 0,
                                                            bottom: 0,
                                                            bgcolor: 'rgba(0,0,0,0.5)',
                                                            display: 'flex',
                                                            alignItems: 'center',
                                                            justifyContent: 'center',
                                                            opacity: 0,
                                                            transition: 'opacity 0.2s',
                                                        }}
                                                    >
                                                        <IconButton
                                                            onClick={() => handleRemoveFile(index)}
                                                            sx={{
                                                                bgcolor: 'error.main',
                                                                color: 'white',
                                                                '&:hover': { bgcolor: 'error.dark' },
                                                            }}
                                                        >
                                                            <Delete />
                                                        </IconButton>
                                                    </Box>
                                                    <Chip
                                                        label={index + 1}
                                                        size="small"
                                                        sx={{
                                                            position: 'absolute',
                                                            bottom: 8,
                                                            left: 8,
                                                            bgcolor: 'rgba(0,0,0,0.6)',
                                                            color: 'white',
                                                        }}
                                                    />
                                                </Box>
                                            </Grid>
                                        ))}
                                    </Grid>
                                </Box>
                            )}

                            {/* Submit Button */}
                            <Button
                                variant="contained"
                                fullWidth
                                size="large"
                                onClick={handleSubmit}
                                disabled={uploading || selectedFiles.length === 0}
                                sx={{
                                    mt: 4,
                                    py: 1.5,
                                    borderRadius: 3,
                                    fontSize: '1rem',
                                    fontWeight: 600,
                                    background: selectedFiles.length > 0
                                        ? 'linear-gradient(135deg, #2563eb 0%, #7c3aed 100%)'
                                        : undefined,
                                    '&:hover': {
                                        background: selectedFiles.length > 0
                                            ? 'linear-gradient(135deg, #1e40af 0%, #6d28d9 100%)'
                                            : undefined,
                                    },
                                }}
                            >
                                {uploading ? (
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                        <LinearProgress sx={{ width: 100 }} />
                                        Đang tải lên...
                                    </Box>
                                ) : (
                                    `Gửi yêu cầu cập nhật (${selectedFiles.length} ảnh)`
                                )}
                            </Button>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Guidelines */}
                <Grid item xs={12} lg={4}>
                    <Card
                        sx={{
                            background: (theme) =>
                                theme.palette.mode === 'dark'
                                    ? 'linear-gradient(135deg, rgba(37, 99, 235, 0.1), rgba(124, 58, 237, 0.1))'
                                    : 'linear-gradient(135deg, rgba(37, 99, 235, 0.05), rgba(124, 58, 237, 0.05))',
                            border: '1px solid',
                            borderColor: (theme) =>
                                theme.palette.mode === 'dark'
                                    ? 'rgba(37, 99, 235, 0.3)'
                                    : 'rgba(37, 99, 235, 0.2)',
                        }}
                    >
                        <CardContent sx={{ p: 3 }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 3 }}>
                                <Avatar sx={{ bgcolor: 'primary.main', width: 36, height: 36 }}>
                                    <CameraAlt fontSize="small" />
                                </Avatar>
                                <Typography variant="h6" fontWeight="700">
                                    Hướng dẫn chụp ảnh
                                </Typography>
                            </Box>

                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
                                {guidelines.map((item, index) => (
                                    <Box
                                        key={index}
                                        sx={{
                                            display: 'flex',
                                            gap: 2,
                                            p: 2,
                                            borderRadius: 2,
                                            bgcolor: 'background.paper',
                                            transition: 'transform 0.2s',
                                            '&:hover': {
                                                transform: 'translateX(4px)',
                                            },
                                        }}
                                    >
                                        <Avatar
                                            sx={{
                                                bgcolor: 'primary.main',
                                                width: 40,
                                                height: 40,
                                            }}
                                        >
                                            {item.icon}
                                        </Avatar>
                                        <Box>
                                            <Typography variant="subtitle2" fontWeight="600">
                                                {item.title}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {item.desc}
                                            </Typography>
                                        </Box>
                                    </Box>
                                ))}
                            </Box>

                            {/* Sample poses */}
                            <Box sx={{ mt: 3, p: 2, bgcolor: 'background.paper', borderRadius: 2 }}>
                                <Typography variant="subtitle2" fontWeight="600" mb={2} textAlign="center">
                                    Các góc chụp khuyến nghị
                                </Typography>
                                <Box sx={{ display: 'flex', justifyContent: 'center', gap: 2 }}>
                                    {['Trái', 'Thẳng', 'Phải'].map((label, i) => (
                                        <Box key={i} sx={{ textAlign: 'center' }}>
                                            <Avatar
                                                sx={{
                                                    width: 56,
                                                    height: 56,
                                                    bgcolor: i === 1 ? 'primary.main' : 'action.selected',
                                                    mb: 0.5,
                                                }}
                                            >
                                                <Face sx={{ transform: i === 0 ? 'rotateY(30deg)' : i === 2 ? 'rotateY(-30deg)' : 'none' }} />
                                            </Avatar>
                                            <Typography variant="caption" color="text.secondary">
                                                {label}
                                            </Typography>
                                        </Box>
                                    ))}
                                </Box>
                            </Box>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>

            {/* Request History */}
            <Card sx={{ mt: 4 }}>
                <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3}>
                        Lịch sử yêu cầu
                    </Typography>
                    {loading ? (
                        <LinearProgress />
                    ) : requests.length === 0 ? (
                        <Box sx={{ textAlign: 'center', py: 6 }}>
                            <Avatar sx={{ width: 64, height: 64, mx: 'auto', mb: 2, bgcolor: 'action.selected' }}>
                                <Face sx={{ fontSize: 32 }} />
                            </Avatar>
                            <Typography color="text.secondary">
                                Bạn chưa có yêu cầu cập nhật nào
                            </Typography>
                        </Box>
                    ) : (
                        <TableContainer component={Paper} elevation={0} sx={{ borderRadius: 2 }}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell sx={{ fontWeight: 600 }}>Ngày gửi</TableCell>
                                        <TableCell sx={{ fontWeight: 600 }}>Trạng thái</TableCell>
                                        <TableCell sx={{ fontWeight: 600 }}>Ghi chú</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {requests.map((req) => (
                                        <TableRow key={req.id} hover>
                                            <TableCell>
                                                {new Date(req.created_at).toLocaleDateString('vi-VN', {
                                                    day: '2-digit',
                                                    month: '2-digit',
                                                    year: 'numeric',
                                                })}
                                            </TableCell>
                                            <TableCell>
                                                <Chip
                                                    icon={statusConfig[req.status].icon}
                                                    label={statusConfig[req.status].label}
                                                    color={statusConfig[req.status].color}
                                                    size="small"
                                                    sx={{ fontWeight: 500 }}
                                                />
                                            </TableCell>
                                            <TableCell>
                                                <Typography variant="body2" color="text.secondary">
                                                    {req.reason || '—'}
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    )}
                </CardContent>
            </Card>
        </Box>
    );
};
