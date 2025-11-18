import { useEffect, useState } from 'react';
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
} from '@mui/material';
import {
    ArrowBack,
    CloudUpload,
    Delete,
    Star,
    StarBorder,
} from '@mui/icons-material';
import {
    listFaceProfiles,
    uploadFaceProfiles,
    setPrimaryFaceProfile,
    deleteFaceProfile,
} from '@face-attendance/utils';
import { FaceProfile } from '@face-attendance/types';

export const EmployeeFaceDataPage: React.FC = () => {
    const navigate = useNavigate();
    const { id: userId } = useParams();
    const [profiles, setProfiles] = useState<FaceProfile[]>([]);
    const [uploading, setUploading] = useState(false);
    const [loading, setLoading] = useState(false);

    // Get company ID from JWT token
    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
        }
    };

    const refreshProfiles = async () => {
        if (!userId) return;
        setLoading(true);
        try {
            const accessToken = localStorage.getItem('access_token');
            if (accessToken) {
                const companyId = getCompanyIdFromToken(accessToken);
                if (companyId) {
                    const data = await listFaceProfiles(userId, companyId);
                    setProfiles(data);
                }
            } else {
                setLoading(false);
            }
        } catch (e) {
            // Could integrate snackbar later
            console.error('Failed to load face profiles', e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        refreshProfiles();
    }, [userId]);

    const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!userId) return;
        const files = e.target.files;
        if (!files || files.length === 0) return;
        setUploading(true);
        try {
            const uploaded = await uploadFaceProfiles(userId, files);
            setProfiles((prev) => [...uploaded, ...prev]);
        } catch (err) {
            console.error('Upload failed', err);
        } finally {
            setUploading(false);
            e.target.value = '';
        }
    };

    const handleSetPrimary = async (profileId: string) => {
        if (!userId) return;
        try {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) return;

            const companyId = getCompanyIdFromToken(accessToken);
            if (!companyId) return;

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
            }
        } catch (e) {
            console.error('Set primary failed', e);
        }
    };

    const handleSetUnPrimary = async (profileId: string) => {
        if (!userId) return;
        try {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) return;

            const companyId = getCompanyIdFromToken(accessToken);
            if (!companyId) return;

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
            }
        } catch (e) {
            console.error('Unset primary failed', e);
        }
    };

    const handleDelete = async (profileId: string) => {
        if (!userId) return;
        try {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) return;

            const companyId = getCompanyIdFromToken(accessToken);
            if (!companyId) return;

            const ok = await deleteFaceProfile(userId, profileId, companyId);
            if (ok) {
                setProfiles((prev) =>
                    prev.filter((p) => p.profile_id !== profileId)
                );
            }
        } catch (e) {
            console.error('Delete failed', e);
        }
    };

    return (
        <Box>
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/employees')}
                sx={{ mb: 2 }}
            >
                Quay lại
            </Button>

            <Card>
                <CardContent>
                    <Typography variant="h5" fontWeight="bold" mb={3}>
                        Quản lý dữ liệu khuôn mặt
                    </Typography>

                    <Paper
                        sx={{
                            p: 4,
                            border: '2px dashed',
                            borderColor: 'primary.main',
                            textAlign: 'center',
                            cursor: 'pointer',
                            mb: 3,
                        }}
                        component="label"
                    >
                        <input
                            type="file"
                            accept="image/*"
                            multiple
                            hidden
                            onChange={handleFileUpload}
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
                    </Paper>

                    <Typography variant="h6" mb={2}>
                        Dữ liệu khuôn mặt ({profiles.length})
                    </Typography>
                    {loading && (
                        <Typography
                            variant="body2"
                            color="text.secondary"
                            mb={2}
                        >
                            Đang tải dữ liệu...
                        </Typography>
                    )}
                    {uploading && (
                        <Typography variant="body2" color="primary" mb={2}>
                            Đang xử lý ảnh tải lên...
                        </Typography>
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
