import { useEffect, useState, useCallback, useMemo } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Button,
    Typography,
    Grid,
    Paper,
    IconButton,
    Alert,
    LinearProgress,
    Snackbar,
    CircularProgress,
} from '@mui/material';
import {
    ArrowBack,
    CloudUpload,
    Delete,
    Star,
    StarBorder,
    CheckCircle,
    Error as ErrorIcon,
} from '@mui/icons-material';
import {
    listFaceProfiles,
    uploadFaceProfiles,
    setPrimaryFaceProfile,
    deleteFaceProfile,
} from '@face-attendance/utils';
import { FaceProfile } from '@face-attendance/types';

interface UploadStatus {
    fileName: string;
    status: 'pending' | 'uploading' | 'success' | 'error';
    error?: string;
    progress?: number;
}

interface FileWithPreview {
    file: File;
    preview: string;
}

// Custom hook for auth token management
const useAuthToken = () => {
    return useMemo(() => {
        const token = localStorage.getItem('access_token');
        if (!token) return { token: null, companyId: null };

        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return {
                token,
                companyId: payload.company_id || null,
            };
        } catch {
            return { token: null, companyId: null };
        }
    }, []);
};

// Constants
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
const BATCH_SIZE = 5; // Upload 5 files at a time

export const EmployeeFaceDataPage: React.FC = () => {
    const navigate = useNavigate();
    const { id: userId } = useParams();
    const { companyId } = useAuthToken();

    const [profiles, setProfiles] = useState<FaceProfile[]>([]);
    const [uploading, setUploading] = useState(false);
    const [loading, setLoading] = useState(false);
    const [uploadStatuses, setUploadStatuses] = useState<UploadStatus[]>([]);
    const [error, setError] = useState<string>('');
    const [successMessage, setSuccessMessage] = useState<string>('');
    const [warningMessage, setWarningMessage] = useState<string>('');
    const [errorKey, setErrorKey] = useState(0);
    const [successKey, setSuccessKey] = useState(0);
    const [warningKey, setWarningKey] = useState(0);
    const [filePreviews, setFilePreviews] = useState<FileWithPreview[]>([]);

    const refreshProfiles = useCallback(async () => {
        if (!userId || !companyId) {
            setError('Thiếu thông tin xác thực');
            return;
        }

        setLoading(true);
        setError('');
        try {
            const data = await listFaceProfiles(userId, companyId);
            setProfiles(data);
        } catch (e: any) {
            const errorMessage =
                e?.response?.data?.message ||
                e?.message ||
                'Không thể tải dữ liệu khuôn mặt';
            setError(errorMessage);
            console.error('Failed to load face profiles', e);
        } finally {
            setLoading(false);
        }
    }, [userId, companyId]);

    useEffect(() => {
        refreshProfiles();
    }, [refreshProfiles]);

    // Cleanup preview URLs on unmount
    useEffect(() => {
        return () => {
            filePreviews.forEach(({ preview }) => URL.revokeObjectURL(preview));
        };
    }, [filePreviews]);

    // Change keys when messages change to force Snackbar to reopen
    useEffect(() => {
        if (error) setErrorKey((k) => k + 1);
    }, [error]);
    useEffect(() => {
        if (successMessage) setSuccessKey((k) => k + 1);
    }, [successMessage]);
    useEffect(() => {
        if (warningMessage) setWarningKey((k) => k + 1);
    }, [warningMessage]);

    // Validate single file
    const validateFile = useCallback((file: File): string | null => {
        if (!ALLOWED_TYPES.includes(file.type)) {
            return `Định dạng không được hỗ trợ (chỉ hỗ trợ JPG, PNG, WEBP)`;
        }
        if (file.size > MAX_FILE_SIZE) {
            return `Kích thước vượt quá 10MB`;
        }
        return null;
    }, []);

    // Upload files in batches
    const uploadInBatches = useCallback(
        async (
            userId: string,
            companyId: string,
            files: File[]
        ): Promise<FaceProfile[]> => {
            const results: FaceProfile[] = [];

            for (let i = 0; i < files.length; i += BATCH_SIZE) {
                const batch = files.slice(i, i + BATCH_SIZE);
                const batchIndices = Array.from(
                    { length: batch.length },
                    (_, idx) => i + idx
                );

                // Mark batch as uploading
                setUploadStatuses((prev) =>
                    prev.map((status, idx) =>
                        batchIndices.includes(idx)
                            ? { ...status, status: 'uploading' }
                            : status
                    )
                );

                try {
                    const uploaded = await uploadFaceProfiles(
                        userId,
                        companyId,
                        batch
                    );
                    const uploadedArr = Array.isArray(uploaded)
                        ? uploaded
                        : uploaded
                        ? [uploaded as unknown as FaceProfile]
                        : [];
                    results.push(...uploadedArr);

                    // Mark batch as success
                    setUploadStatuses((prev) =>
                        prev.map((status, idx) =>
                            batchIndices.includes(idx)
                                ? { ...status, status: 'success' }
                                : status
                        )
                    );
                } catch (err: any) {
                    const errorMessage =
                        err?.response?.data?.message ||
                        err?.message ||
                        'Tải ảnh lên thất bại';

                    // Mark batch as error
                    setUploadStatuses((prev) =>
                        prev.map((status, idx) =>
                            batchIndices.includes(idx)
                                ? {
                                      ...status,
                                      status: 'error',
                                      error: errorMessage,
                                  }
                                : status
                        )
                    );

                    console.error(
                        `Batch ${i / BATCH_SIZE + 1} upload failed`,
                        err
                    );
                }

                // Small delay between batches to prevent overwhelming the server
                if (i + BATCH_SIZE < files.length) {
                    await new Promise((resolve) => setTimeout(resolve, 500));
                }
            }

            return results;
        },
        []
    );

    const handleFileUpload = useCallback(
        async (e: React.ChangeEvent<HTMLInputElement>) => {
            if (!userId || !companyId) {
                setError('Thiếu thông tin xác thực');
                return;
            }

            const files = e.target.files;
            if (!files || files.length === 0) return;

            // Validate files
            const validFilesWithPreviews: FileWithPreview[] = [];
            const errors: string[] = [];

            Array.from(files).forEach((file) => {
                const validationError = validateFile(file);
                if (validationError) {
                    errors.push(`${file.name}: ${validationError}`);
                } else {
                    validFilesWithPreviews.push({
                        file,
                        preview: URL.createObjectURL(file),
                    });
                }
            });

            if (errors.length > 0 && validFilesWithPreviews.length === 0) {
                // Only invalid files selected -> show error and exit
                setError(errors.join('; '));
                e.target.value = '';
                return;
            }

            if (errors.length > 0 && validFilesWithPreviews.length > 0) {
                // Mixed valid/invalid -> continue with upload, but show a warning summary
                const shown = errors.slice(0, 2).join('; ');
                const more =
                    errors.length > 2
                        ? `; ...và ${errors.length - 2} lỗi khác`
                        : '';
                setWarningMessage(
                    `Bỏ qua ${errors.length} ảnh không hợp lệ: ${shown}${more}`
                );
            }

            if (validFilesWithPreviews.length === 0) {
                e.target.value = '';
                return;
            }

            // Cleanup old previews
            filePreviews.forEach(({ preview }) => URL.revokeObjectURL(preview));
            setFilePreviews(validFilesWithPreviews);

            // Initialize upload statuses
            const statuses: UploadStatus[] = validFilesWithPreviews.map(
                ({ file }) => ({
                    fileName: file.name,
                    status: 'pending',
                    progress: 0,
                })
            );
            setUploadStatuses(statuses);

            setUploading(true);
            setError('');

            try {
                const validFiles = validFilesWithPreviews.map(
                    ({ file }) => file
                );
                const uploaded = await uploadInBatches(
                    userId,
                    companyId,
                    validFiles
                );

                if (uploaded.length > 0) {
                    setProfiles((prev) => [...uploaded, ...prev]);
                    setSuccessMessage(
                        `Đã tải lên thành công ${uploaded.length}/${validFiles.length} ảnh`
                    );
                    if (uploaded.length < validFiles.length) {
                        setWarningMessage(
                            `${
                                validFiles.length - uploaded.length
                            } ảnh tải lên thất bại`
                        );
                    }
                }

                // Clear previews and statuses after success
                setTimeout(() => {
                    validFilesWithPreviews.forEach(({ preview }) =>
                        URL.revokeObjectURL(preview)
                    );
                    setFilePreviews([]);
                    setUploadStatuses([]);
                }, 3000);
            } catch (err: any) {
                console.error('Upload failed', err);
            } finally {
                setUploading(false);
                e.target.value = '';
            }
        },
        [userId, companyId, validateFile, uploadInBatches, filePreviews]
    );

    const handleSetPrimary = useCallback(
        async (profileId: string) => {
            if (!userId || !companyId) {
                setError('Thiếu thông tin xác thực');
                return;
            }

            try {
                const updated = await setPrimaryFaceProfile(
                    userId,
                    profileId,
                    companyId,
                    true
                );
                if (updated) {
                    setProfiles((prev) =>
                        prev.map((p) => ({
                            ...p,
                            is_primary: p.profile_id === profileId,
                        }))
                    );
                    setSuccessMessage('Đã đặt ảnh làm ảnh chính');
                }
            } catch (e: any) {
                const errorMessage =
                    e?.response?.data?.message ||
                    e?.message ||
                    'Không thể đặt ảnh chính';
                setError(errorMessage);
                console.error('Set primary failed', e);
            }
        },
        [userId, companyId]
    );

    const handleSetUnPrimary = useCallback(
        async (profileId: string) => {
            if (!userId || !companyId) {
                setError('Thiếu thông tin xác thực');
                return;
            }

            try {
                const updated = await setPrimaryFaceProfile(
                    userId,
                    profileId,
                    companyId,
                    false
                );
                if (updated) {
                    setProfiles((prev) =>
                        prev.map((p) =>
                            p.profile_id === profileId
                                ? { ...p, is_primary: false }
                                : p
                        )
                    );
                    setSuccessMessage('Đã bỏ đặt ảnh chính');
                }
            } catch (e: any) {
                const errorMessage =
                    e?.response?.data?.message ||
                    e?.message ||
                    'Không thể bỏ đặt ảnh chính';
                setError(errorMessage);
                console.error('Unset primary failed', e);
            }
        },
        [userId, companyId]
    );

    const handleDelete = useCallback(
        async (profileId: string) => {
            if (!userId || !companyId) {
                setError('Thiếu thông tin xác thực');
                return;
            }

            if (!window.confirm('Bạn có chắc chắn muốn xóa ảnh này?')) {
                return;
            }

            try {
                const ok = await deleteFaceProfile(
                    userId,
                    profileId,
                    companyId
                );
                if (ok) {
                    setProfiles((prev) =>
                        prev.filter((p) => p.profile_id !== profileId)
                    );
                    setSuccessMessage('Đã xóa ảnh thành công');
                }
            } catch (e: any) {
                const errorMessage =
                    e?.response?.data?.message ||
                    e?.message ||
                    'Không thể xóa ảnh';
                setError(errorMessage);
                console.error('Delete failed', e);
            }
        },
        [userId, companyId]
    );

    // Calculate upload progress
    const uploadProgress = useMemo(() => {
        if (uploadStatuses.length === 0) return 0;
        const completed = uploadStatuses.filter(
            (s) => s.status === 'success' || s.status === 'error'
        ).length;
        return Math.round((completed / uploadStatuses.length) * 100);
    }, [uploadStatuses]);

    // Calculate auto-hide duration based on message length (200ms per 10 characters, min 4s, max 10s)
    const getAutoHideDuration = useCallback((message: string) => {
        const baseTime = 4000; // 4 seconds minimum
        const timePerChar = 50; // 50ms per character
        const calculatedTime = baseTime + message.length * timePerChar;
        return Math.min(Math.max(calculatedTime, baseTime), 10000); // Between 4-10 seconds
    }, []);

    return (
        <Box>
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/employees')}
                sx={{ mb: 2 }}
            >
                Quay lại
            </Button>

            {/* Error Snackbar */}
            <Snackbar
                key={errorKey}
                open={!!error}
                autoHideDuration={getAutoHideDuration(error)}
                onClose={() => setError('')}
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
            >
                <Alert
                    onClose={() => setError('')}
                    severity="error"
                    sx={{ width: '100%', maxWidth: '600px' }}
                >
                    {error}
                </Alert>
            </Snackbar>

            {/* Success Snackbar */}
            <Snackbar
                key={successKey}
                open={!!successMessage}
                autoHideDuration={getAutoHideDuration(successMessage)}
                onClose={() => setSuccessMessage('')}
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
            >
                <Alert
                    onClose={() => setSuccessMessage('')}
                    severity="success"
                    sx={{ width: '100%', maxWidth: '600px' }}
                >
                    {successMessage}
                </Alert>
            </Snackbar>

            {/* Warning Snackbar */}
            <Snackbar
                key={warningKey}
                open={!!warningMessage}
                autoHideDuration={getAutoHideDuration(warningMessage)}
                onClose={() => setWarningMessage('')}
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
            >
                <Alert
                    onClose={() => setWarningMessage('')}
                    severity="warning"
                    sx={{ width: '100%', maxWidth: '600px' }}
                >
                    {warningMessage}
                </Alert>
            </Snackbar>

            <Card>
                <CardContent>
                    <Typography variant="h5" fontWeight="bold" mb={3}>
                        Quản lý dữ liệu khuôn mặt
                    </Typography>

                    <Paper
                        sx={{
                            p: 4,
                            border: '2px dashed',
                            borderColor: uploading
                                ? 'grey.400'
                                : 'primary.main',
                            textAlign: 'center',
                            cursor: uploading ? 'not-allowed' : 'pointer',
                            mb: 3,
                            opacity: uploading ? 0.6 : 1,
                        }}
                        component="label"
                    >
                        <input
                            type="file"
                            accept="image/*"
                            multiple
                            hidden
                            onChange={handleFileUpload}
                            disabled={uploading}
                        />
                        <CloudUpload
                            sx={{ fontSize: 48, color: 'primary.main', mb: 1 }}
                        />
                        <Typography variant="h6">
                            Kéo thả ảnh vào đây hoặc click để chọn
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                            Chọn nhiều ảnh khuôn mặt từ các góc độ khác nhau
                        </Typography>
                        <Typography
                            variant="caption"
                            color="textSecondary"
                            display="block"
                            mt={1}
                        >
                            Hỗ trợ: JPG, PNG, WEBP • Tối đa: 10MB/ảnh
                        </Typography>
                    </Paper>

                    {/* Upload Progress */}
                    {uploading && (
                        <Box sx={{ mb: 3 }}>
                            <Box
                                sx={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'space-between',
                                    mb: 1,
                                }}
                            >
                                <Box
                                    sx={{
                                        display: 'flex',
                                        alignItems: 'center',
                                    }}
                                >
                                    <CircularProgress
                                        size={20}
                                        sx={{ mr: 1 }}
                                    />
                                    <Typography variant="body2" color="primary">
                                        Đang xử lý {uploadStatuses.length}{' '}
                                        ảnh...
                                    </Typography>
                                </Box>
                                <Typography
                                    variant="body2"
                                    color="primary"
                                    fontWeight="bold"
                                >
                                    {uploadProgress}%
                                </Typography>
                            </Box>
                            <LinearProgress
                                variant="determinate"
                                value={uploadProgress}
                            />
                        </Box>
                    )}

                    {/* Upload Status List */}
                    {uploadStatuses.length > 0 && (
                        <Box sx={{ mb: 3 }}>
                            {uploadStatuses.map((status, idx) => (
                                <Box
                                    key={idx}
                                    sx={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        gap: 1,
                                        mb: 0.5,
                                    }}
                                >
                                    {status.status === 'success' && (
                                        <CheckCircle
                                            color="success"
                                            fontSize="small"
                                        />
                                    )}
                                    {status.status === 'error' && (
                                        <ErrorIcon
                                            color="error"
                                            fontSize="small"
                                        />
                                    )}
                                    {status.status === 'uploading' && (
                                        <CircularProgress size={16} />
                                    )}
                                    <Typography
                                        variant="caption"
                                        color={
                                            status.status === 'error'
                                                ? 'error'
                                                : status.status === 'success'
                                                ? 'success.main'
                                                : 'text.secondary'
                                        }
                                    >
                                        {status.fileName}
                                        {status.error && ` - ${status.error}`}
                                    </Typography>
                                </Box>
                            ))}
                        </Box>
                    )}

                    {/* Preview Images */}
                    {filePreviews.length > 0 && (
                        <Box sx={{ mb: 3 }}>
                            <Typography variant="body2" color="primary" mb={1}>
                                Xem trước ({filePreviews.length} ảnh)
                            </Typography>
                            <Grid container spacing={1}>
                                {filePreviews.map(({ preview, file }, idx) => (
                                    <Grid item xs={3} sm={2} md={1.5} key={idx}>
                                        <Paper
                                            sx={{
                                                aspectRatio: '1',
                                                overflow: 'hidden',
                                                borderRadius: 1,
                                                position: 'relative',
                                            }}
                                        >
                                            <img
                                                src={preview}
                                                alt={file.name}
                                                style={{
                                                    width: '100%',
                                                    height: '100%',
                                                    objectFit: 'cover',
                                                }}
                                                loading="lazy"
                                            />
                                        </Paper>
                                    </Grid>
                                ))}
                            </Grid>
                        </Box>
                    )}

                    <Typography variant="h6" mb={2}>
                        Dữ liệu khuôn mặt ({profiles.length})
                    </Typography>
                    {loading && (
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                mb: 2,
                                gap: 1,
                            }}
                        >
                            <CircularProgress size={20} />
                            <Typography variant="body2" color="text.secondary">
                                Đang tải dữ liệu...
                            </Typography>
                        </Box>
                    )}

                    <Grid container spacing={2}>
                        {profiles.map((p) => {
                            const img =
                                p.image_url || p.enroll_image_path || '';
                            return (
                                <Grid
                                    item
                                    xs={6}
                                    sm={4}
                                    md={3}
                                    key={p.profile_id}
                                >
                                    <Paper
                                        sx={{
                                            position: 'relative',
                                            borderRadius: 2,
                                            overflow: 'hidden',
                                            display: 'flex',
                                            flexDirection: 'column',
                                        }}
                                    >
                                        {img && (
                                            <img
                                                src={img}
                                                alt={p.profile_id}
                                                style={{
                                                    width: '100%',
                                                    height: 160,
                                                    objectFit: 'cover',
                                                    display: 'block',
                                                }}
                                            />
                                        )}
                                        <Box sx={{ p: 1 }}>
                                            <Typography
                                                variant="caption"
                                                display="block"
                                            >
                                                Chất lượng:{' '}
                                                {p.quality_score ?? 'N/A'}
                                            </Typography>
                                            <Typography
                                                variant="caption"
                                                display="block"
                                            >
                                                Phiên bản embedding:{' '}
                                                {p.embedding_version}
                                            </Typography>
                                            <Typography
                                                variant="caption"
                                                display="block"
                                            >
                                                Trạng thái index:{' '}
                                                {p.indexed
                                                    ? 'Indexed'
                                                    : 'Pending'}{' '}
                                                v{p.index_version}
                                            </Typography>
                                        </Box>
                                        <Box
                                            sx={{
                                                position: 'absolute',
                                                top: 6,
                                                right: 6,
                                                display: 'flex',
                                                gap: 0.5,
                                            }}
                                        >
                                            <IconButton
                                                size="small"
                                                color={
                                                    p.is_primary
                                                        ? 'warning'
                                                        : 'default'
                                                }
                                                onClick={() =>
                                                    p.is_primary
                                                        ? handleSetUnPrimary(
                                                              p.profile_id
                                                          )
                                                        : handleSetPrimary(
                                                              p.profile_id
                                                          )
                                                }
                                                sx={{
                                                    backgroundColor: 'white',
                                                    ':hover': {
                                                        backgroundColor: '#ffe',
                                                    },
                                                }}
                                            >
                                                {p.is_primary ? (
                                                    <Star fontSize="small" />
                                                ) : (
                                                    <StarBorder fontSize="small" />
                                                )}
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                color="error"
                                                onClick={() =>
                                                    handleDelete(p.profile_id)
                                                }
                                                sx={{
                                                    backgroundColor: 'white',
                                                    ':hover': {
                                                        backgroundColor: '#fdd',
                                                    },
                                                }}
                                            >
                                                <Delete fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    </Paper>
                                </Grid>
                            );
                        })}
                    </Grid>
                </CardContent>
            </Card>
        </Box>
    );
};
