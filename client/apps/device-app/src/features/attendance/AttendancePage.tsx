import { useState, useRef, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Typography,
    Paper,
    IconButton,
    Chip,
    CircularProgress,
} from '@mui/material';
import {
    CameraAlt,
    Logout,
    CheckCircle,
    Error as ErrorIcon,
    AccessTime,
    LocationOn,
    Refresh,
    Login as LoginIcon,
    ExitToApp,
} from '@mui/icons-material';
import { useDeviceStore } from '@/stores/deviceStore';
import { deviceFaceApi } from '@/services/api';
import { useTokenRefresh } from '@/hooks/useTokenRefresh';
import {
    detectFaceInFrame,
    hasGoodImageQuality,
    getDetectionFeedback,
} from '@/utils/faceDetection';

interface VerifyResult {
    success: boolean;
    employee?: {
        name: string;
        employeeId: string;
        avatar?: string;
    };
    attendanceType?: 'check-in' | 'check-out';
    time?: string;
    message?: string;
}

export const AttendancePage: React.FC = () => {
    const navigate = useNavigate();
    const { token, deviceInfo, clearDevice } = useDeviceStore();
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const streamRef = useRef<MediaStream | null>(null);
    const detectionIntervalRef = useRef<ReturnType<typeof setInterval> | null>(
        null
    );
    const lastAttendanceRef = useRef<number>(0);
    const lastDetectionLogRef = useRef<number>(0);

    // Auto-refresh token
    useTokenRefresh();

    const [currentTime, setCurrentTime] = useState(new Date());
    const [cameraReady, setCameraReady] = useState(false);
    const [cameraError, setCameraError] = useState<string | null>(null);
    const [isProcessing, setIsProcessing] = useState(false);
    const [result, setResult] = useState<VerifyResult | null>(null);
    const [faceDetected, setFaceDetected] = useState(false);
    const [faceQuality, setFaceQuality] = useState<number>(0);
    const [status, setStatus] = useState('Đang chờ...');

    // Cooldown period between attendances (5 seconds)
    const COOLDOWN_MS = 5000;
    // Interval for auto-detection (500ms)
    const DETECTION_INTERVAL = 500;

    // Determine attendance type based on time of day
    const getAttendanceType = (): 'check-in' | 'check-out' => {
        const hour = new Date().getHours();
        // Before 12:00 = check-in, after = check-out (simple logic)
        return hour < 12 ? 'check-in' : 'check-out';
    };

    // Update clock every second
    useEffect(() => {
        const timer = setInterval(() => setCurrentTime(new Date()), 1000);
        return () => clearInterval(timer);
    }, []);

    // Auto-hide result after 4 seconds
    useEffect(() => {
        if (result) {
            const timer = setTimeout(() => setResult(null), 4000);
            return () => clearTimeout(timer);
        }
    }, [result]);

    // Initialize camera
    const initCamera = useCallback(async () => {
        try {
            setCameraError(null);
            setCameraReady(false);
            const stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 1280 },
                    height: { ideal: 720 },
                    facingMode: 'user',
                },
                audio: false,
            });

            if (videoRef.current) {
                videoRef.current.srcObject = stream;
                streamRef.current = stream;

                const markReady = () => setCameraReady(true);

                if (videoRef.current.readyState >= 1) {
                    markReady();
                } else {
                    videoRef.current.onloadedmetadata = markReady;
                }
            }
        } catch (err: any) {
            console.error('Camera error:', err);
            setCameraError(
                err.name === 'NotAllowedError'
                    ? 'Vui lòng cho phép truy cập camera để sử dụng chấm công'
                    : 'Không thể kết nối camera. Vui lòng kiểm tra thiết bị.'
            );
        }
    }, []);

    useEffect(() => {
        initCamera();
        return () => {
            if (streamRef.current) {
                streamRef.current.getTracks().forEach((track) => track.stop());
            }
            if (detectionIntervalRef.current) {
                clearInterval(detectionIntervalRef.current);
            }
        };
    }, [initCamera]);

    // Capture image from camera
    const captureImage = (): Blob | null => {
        if (!videoRef.current || !canvasRef.current) return null;

        const video = videoRef.current;
        // Guard against early calls before video has dimensions
        if (
            !video.videoWidth ||
            !video.videoHeight ||
            video.videoWidth < 10 ||
            video.videoHeight < 10
        ) {
            return null;
        }

        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');

        if (!context) return null;

        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        context.drawImage(video, 0, 0);

        // Synchronous blob creation
        const dataUrl = canvas.toDataURL('image/jpeg', 0.9);
        const byteString = atob(dataUrl.split(',')[1]);
        const mimeString = dataUrl.split(',')[0].split(':')[1].split(';')[0];
        const ab = new ArrayBuffer(byteString.length);
        const ia = new Uint8Array(ab);
        for (let i = 0; i < byteString.length; i++) {
            ia[i] = byteString.charCodeAt(i);
        }
        return new Blob([ab], { type: mimeString });
    };

    // Handle automatic attendance
    const handleAutoAttendance = useCallback(async () => {
        // Check cooldown
        const now = Date.now();
        if (now - lastAttendanceRef.current < COOLDOWN_MS) {
            return;
        }

        if (!cameraReady || isProcessing) return;

        setIsProcessing(true);
        setStatus('Đang nhận diện...');

        try {
            const imageBlob = captureImage();
            if (!imageBlob) {
                console.warn(
                    '[Face] captureImage returned null (camera not ready or canvas error)'
                );
                throw new Error('Không thể chụp ảnh từ camera');
            }

            // Call device face verification API
            const response: any = await deviceFaceApi.verify(imageBlob);

            if (response.code !== 200) {
                throw new Error(response.message || 'Xác thực thất bại');
            }

            const verifyData = response.data;

            // Check if face matches and liveness is good
            if (
                verifyData.verified &&
                verifyData.status === 'match' &&
                verifyData.liveness_score >= 0.8
            ) {
                const match = verifyData.matches?.[0];

                if (match) {
                    const attendanceType = getAttendanceType();

                    // Set cooldown
                    lastAttendanceRef.current = Date.now();

                    setResult({
                        success: true,
                        employee: {
                            name: match.user_name || 'Nhân viên',
                            employeeId: match.user_id || '',
                        },
                        attendanceType,
                        time: new Date().toLocaleTimeString('vi-VN'),
                        message:
                            attendanceType === 'check-in'
                                ? 'Chấm công VÀO thành công!'
                                : 'Chấm công RA thành công!',
                    });

                    setStatus('Chấm công thành công!');
                }
            } else {
                // Face not verified
                const message = verifyData.message || 'Khuôn mặt không khớp';
                setStatus(message);

                // Show message briefly then reset
                setTimeout(() => setStatus('Đang chờ...'), 2000);
            }
        } catch (err: any) {
            console.error('Attendance error:', err);
            setStatus('Lỗi kết nối');
            setTimeout(() => setStatus('Đang chờ...'), 2000);
        } finally {
            setIsProcessing(false);
        }
    }, [cameraReady, isProcessing]);

    // Auto-detection loop
    useEffect(() => {
        if (!cameraReady) return;

        detectionIntervalRef.current = setInterval(() => {
            if (videoRef.current && canvasRef.current && !isProcessing) {
                // Skip detection until video has valid dimensions
                if (
                    !videoRef.current.videoWidth ||
                    !videoRef.current.videoHeight ||
                    videoRef.current.videoWidth < 10
                ) {
                    return;
                }

                const detection = detectFaceInFrame(
                    videoRef.current,
                    canvasRef.current
                );
                setFaceDetected(detection.hasFace);
                setFaceQuality(detection.confidence);

                // Log detection metrics at most every 3 seconds to avoid spam
                const now = Date.now();
                if (now - lastDetectionLogRef.current > 3000) {
                    lastDetectionLogRef.current = now;
                    console.log('[Face] detection', {
                        hasFace: detection.hasFace,
                        confidence: detection.confidence,
                        brightness: detection.brightness,
                        sharpness: detection.sharpness,
                        videoWidth: videoRef.current.videoWidth,
                        videoHeight: videoRef.current.videoHeight,
                    });
                }

                // Update status based on detection
                if (!detection.hasFace) {
                    const feedback = getDetectionFeedback(detection);
                    setStatus(feedback);
                } else if (detection.hasFace && !isProcessing) {
                    const hasGoodQuality = hasGoodImageQuality(detection);

                    if (hasGoodQuality) {
                        setStatus('Khuôn mặt được phát hiện');

                        // Check cooldown before triggering verification
                        if (now - lastAttendanceRef.current >= COOLDOWN_MS) {
                            handleAutoAttendance();
                        }
                    } else {
                        setStatus(getDetectionFeedback(detection));
                    }
                }
            }
        }, DETECTION_INTERVAL);

        return () => {
            if (detectionIntervalRef.current) {
                clearInterval(detectionIntervalRef.current);
            }
        };
    }, [cameraReady, isProcessing, handleAutoAttendance]);

    const handleLogout = () => {
        clearDevice();
        navigate('/token-auth');
    };

    const attendanceType = getAttendanceType();

    return (
        <Box
            sx={{
                height: '100vh',
                width: '100%',
                display: 'flex',
                flexDirection: 'column',
                background: 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%)',
                overflow: 'hidden',
            }}
        >
            {/* Header */}
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    p: 2,
                    px: 3,
                    background: 'rgba(255, 255, 255, 0.95)',
                    backdropFilter: 'blur(10px)',
                    borderBottom: '1px solid rgba(226, 232, 240, 0.8)',
                    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05)',
                }}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Box
                        sx={{
                            width: 44,
                            height: 44,
                            borderRadius: 2,
                            background:
                                'linear-gradient(135deg, #2563eb, #7c3aed)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                        }}
                    >
                        <CameraAlt sx={{ color: 'white', fontSize: 24 }} />
                    </Box>
                    <Box>
                        <Typography
                            variant="h6"
                            fontWeight="700"
                            color="text.primary"
                        >
                            {deviceInfo?.companyName || 'CIO Verify'}
                        </Typography>
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: 0.5,
                            }}
                        >
                            <LocationOn
                                sx={{ fontSize: 14, color: 'text.secondary' }}
                            />
                            <Typography
                                variant="caption"
                                color="text.secondary"
                            >
                                {deviceInfo?.location ||
                                    deviceInfo?.deviceName ||
                                    'Device'}
                            </Typography>
                        </Box>
                    </Box>
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                    <Box sx={{ textAlign: 'right' }}>
                        <Typography
                            variant="h4"
                            fontWeight="700"
                            color="primary.main"
                        >
                            {currentTime.toLocaleTimeString('vi-VN', {
                                hour: '2-digit',
                                minute: '2-digit',
                            })}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                            {currentTime.toLocaleDateString('vi-VN', {
                                weekday: 'long',
                                day: 'numeric',
                                month: 'long',
                            })}
                        </Typography>
                    </Box>
                    <Chip
                        icon={
                            attendanceType === 'check-in' ? (
                                <LoginIcon />
                            ) : (
                                <ExitToApp />
                            )
                        }
                        label={
                            attendanceType === 'check-in'
                                ? 'Chế độ VÀO'
                                : 'Chế độ RA'
                        }
                        color={
                            attendanceType === 'check-in'
                                ? 'success'
                                : 'warning'
                        }
                        sx={{ fontWeight: 600 }}
                    />
                    <Chip
                        label="Online"
                        size="small"
                        color="success"
                        sx={{ fontWeight: 600 }}
                    />
                    <IconButton
                        onClick={handleLogout}
                        size="small"
                        color="default"
                    >
                        <Logout />
                    </IconButton>
                </Box>
            </Box>

            {/* Main Content */}
            <Box
                sx={{
                    flex: 1,
                    display: 'flex',
                    gap: 3,
                    p: 3,
                    overflow: 'hidden',
                }}
            >
                {/* Camera Section */}
                <Paper
                    elevation={0}
                    sx={{
                        flex: 2,
                        borderRadius: 4,
                        overflow: 'hidden',
                        position: 'relative',
                        background: '#000',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    {cameraError ? (
                        <Box sx={{ textAlign: 'center', color: 'white', p: 4 }}>
                            <ErrorIcon
                                sx={{
                                    fontSize: 64,
                                    mb: 2,
                                    color: 'error.main',
                                }}
                            />
                            <Typography variant="h6" mb={2}>
                                {cameraError}
                            </Typography>
                            <IconButton
                                onClick={initCamera}
                                sx={{ color: 'white' }}
                            >
                                <Refresh sx={{ fontSize: 32 }} />
                            </IconButton>
                        </Box>
                    ) : (
                        <>
                            <video
                                ref={videoRef}
                                autoPlay
                                playsInline
                                muted
                                style={{
                                    width: '100%',
                                    height: '100%',
                                    objectFit: 'cover',
                                }}
                            />
                            {/* Face detection overlay */}
                            <Box
                                sx={{
                                    position: 'absolute',
                                    top: '50%',
                                    left: '50%',
                                    transform: 'translate(-50%, -50%)',
                                    width: 280,
                                    height: 350,
                                    border: '3px solid',
                                    borderColor: isProcessing
                                        ? 'warning.main'
                                        : faceDetected
                                        ? 'success.main'
                                        : 'primary.main',
                                    borderRadius: 4,
                                    boxShadow: faceDetected
                                        ? '0 0 40px rgba(16, 185, 129, 0.5)'
                                        : '0 0 40px rgba(37, 99, 235, 0.3)',
                                    animation: isProcessing
                                        ? 'pulse 1s infinite'
                                        : 'none',
                                    transition: 'all 0.3s ease',
                                }}
                            />
                            {/* Processing indicator */}
                            {isProcessing && (
                                <Box
                                    sx={{
                                        position: 'absolute',
                                        top: '50%',
                                        left: '50%',
                                        transform: 'translate(-50%, -50%)',
                                    }}
                                >
                                    <CircularProgress
                                        size={80}
                                        sx={{ color: 'white' }}
                                    />
                                </Box>
                            )}
                            {/* Status */}
                            <Box
                                sx={{
                                    position: 'absolute',
                                    bottom: 24,
                                    left: '50%',
                                    transform: 'translateX(-50%)',
                                    background: 'rgba(0, 0, 0, 0.7)',
                                    backdropFilter: 'blur(10px)',
                                    borderRadius: 2,
                                    px: 3,
                                    py: 1.5,
                                    minWidth: 200,
                                    textAlign: 'center',
                                }}
                            >
                                <Typography
                                    component="div"
                                    color="white"
                                    fontWeight="500"
                                    sx={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        gap: 1,
                                    }}
                                >
                                    {faceDetected && !isProcessing && (
                                        <Box
                                            sx={{
                                                width: 8,
                                                height: 8,
                                                borderRadius: '50%',
                                                bgcolor: 'success.main',
                                                animation: 'pulse 1s infinite',
                                            }}
                                        />
                                    )}
                                    {status}
                                </Typography>
                                {faceDetected && faceQuality > 0 && (
                                    <Box sx={{ mt: 1 }}>
                                        <Typography
                                            variant="caption"
                                            color="rgba(255,255,255,0.7)"
                                        >
                                            Chất lượng:{' '}
                                            {Math.round(faceQuality * 100)}%
                                        </Typography>
                                    </Box>
                                )}
                            </Box>
                        </>
                    )}
                    <canvas ref={canvasRef} style={{ display: 'none' }} />
                </Paper>

                {/* Info Panel */}
                <Box
                    sx={{
                        flex: 1,
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 2,
                    }}
                >
                    {/* Result Card */}
                    {result && (
                        <Paper
                            elevation={0}
                            sx={{
                                p: 3,
                                borderRadius: 3,
                                background: result.success
                                    ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(16, 185, 129, 0.05))'
                                    : 'linear-gradient(135deg, rgba(239, 68, 68, 0.1), rgba(239, 68, 68, 0.05))',
                                border: '1px solid',
                                borderColor: result.success
                                    ? 'success.light'
                                    : 'error.light',
                                animation: 'slideUp 0.3s ease-out',
                            }}
                        >
                            <Box
                                sx={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 2,
                                    mb: 2,
                                }}
                            >
                                {result.success ? (
                                    <CheckCircle
                                        sx={{
                                            fontSize: 48,
                                            color: 'success.main',
                                        }}
                                    />
                                ) : (
                                    <ErrorIcon
                                        sx={{
                                            fontSize: 48,
                                            color: 'error.main',
                                        }}
                                    />
                                )}
                                <Box>
                                    <Typography
                                        variant="h6"
                                        fontWeight="700"
                                        color={
                                            result.success
                                                ? 'success.main'
                                                : 'error.main'
                                        }
                                    >
                                        {result.success
                                            ? 'Thành công!'
                                            : 'Thất bại'}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        color="text.secondary"
                                    >
                                        {result.message}
                                    </Typography>
                                </Box>
                            </Box>
                            {result.success && result.employee && (
                                <Box
                                    sx={{
                                        mt: 2,
                                        p: 2,
                                        background: 'rgba(255,255,255,0.8)',
                                        borderRadius: 2,
                                    }}
                                >
                                    <Typography
                                        fontWeight="600"
                                        fontSize="1.2rem"
                                    >
                                        {result.employee.name}
                                    </Typography>
                                    <Typography
                                        variant="caption"
                                        color="text.secondary"
                                    >
                                        ID: {result.employee.employeeId}
                                    </Typography>
                                    <Box
                                        sx={{
                                            display: 'flex',
                                            alignItems: 'center',
                                            gap: 0.5,
                                            mt: 1,
                                        }}
                                    >
                                        <AccessTime
                                            sx={{
                                                fontSize: 16,
                                                color: 'text.secondary',
                                            }}
                                        />
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                        >
                                            {result.time}
                                        </Typography>
                                        <Chip
                                            size="small"
                                            label={
                                                result.attendanceType ===
                                                'check-in'
                                                    ? 'VÀO'
                                                    : 'RA'
                                            }
                                            color={
                                                result.attendanceType ===
                                                'check-in'
                                                    ? 'success'
                                                    : 'warning'
                                            }
                                            sx={{ ml: 1 }}
                                        />
                                    </Box>
                                </Box>
                            )}
                        </Paper>
                    )}

                    {/* Auto Mode Info */}
                    <Paper
                        elevation={0}
                        sx={{ p: 3, borderRadius: 3, bgcolor: 'primary.50' }}
                    >
                        <Typography
                            variant="h6"
                            fontWeight="700"
                            color="primary.main"
                            mb={1}
                        >
                            🤖 Chế độ Tự động
                        </Typography>
                        <Typography
                            variant="body2"
                            color="text.secondary"
                            mb={2}
                        >
                            Hệ thống sẽ tự động nhận diện và chấm công khi phát
                            hiện khuôn mặt trong khung.
                        </Typography>
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 1,
                            }}
                        >
                            <Typography variant="body2">
                                • <strong>Buổi sáng (trước 12h):</strong> Chấm
                                công VÀO
                            </Typography>
                            <Typography variant="body2">
                                • <strong>Buổi chiều (sau 12h):</strong> Chấm
                                công RA
                            </Typography>
                        </Box>
                    </Paper>

                    {/* Status Info */}
                    <Paper
                        elevation={0}
                        sx={{ p: 2, borderRadius: 3, mt: 'auto' }}
                    >
                        <Typography
                            variant="subtitle2"
                            color="text.secondary"
                            mb={1}
                        >
                            Thông tin thiết bị
                        </Typography>
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 0.5,
                            }}
                        >
                            <Typography variant="body2">
                                <strong>Thiết bị:</strong>{' '}
                                {deviceInfo?.deviceName || 'Test Device'}
                            </Typography>
                            <Typography variant="body2">
                                <strong>Địa điểm:</strong>{' '}
                                {deviceInfo?.location || 'Test Location'}
                            </Typography>
                            <Typography variant="body2">
                                <strong>Trạng thái:</strong>{' '}
                                {faceDetected
                                    ? '🟢 Phát hiện khuôn mặt'
                                    : '⚪ Đang chờ...'}
                            </Typography>
                        </Box>
                    </Paper>
                </Box>
            </Box>
        </Box>
    );
};
