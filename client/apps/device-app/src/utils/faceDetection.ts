/**
 * Face detection utilities using canvas-based analysis
 * This provides basic face detection for liveness checking before sending to server
 */

export interface FaceDetectionResult {
    hasFace: boolean;
    confidence: number;
    brightness: number;
    sharpness: number;
}

/**
 * Detect face in video frame using brightness and edge detection
 * @param video - Video element to analyze
 * @param canvas - Canvas element for processing
 * @returns Detection result with confidence score
 */
export const detectFaceInFrame = (
    video: HTMLVideoElement,
    canvas: HTMLCanvasElement
): FaceDetectionResult => {
    // Wait until video has valid dimensions to avoid 0x0 reads
    if (
        !video.videoWidth ||
        !video.videoHeight ||
        video.videoWidth < 10 ||
        video.videoHeight < 10
    ) {
        return {
            hasFace: false,
            confidence: 0,
            brightness: 0,
            sharpness: 0,
        };
    }

    const ctx = canvas.getContext('2d');
    if (!ctx) {
        return {
            hasFace: false,
            confidence: 0,
            brightness: 0,
            sharpness: 0,
        };
    }

    // Set canvas size to match video
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    ctx.drawImage(video, 0, 0);

    // Define face detection area (center of frame)
    const centerX = video.videoWidth / 2;
    const centerY = video.videoHeight / 2;
    const faceWidth = Math.min(280, video.videoWidth * 0.4);
    const faceHeight = Math.min(350, video.videoHeight * 0.5);

    // Get image data from center region
    const imageData = ctx.getImageData(
        centerX - faceWidth / 2,
        centerY - faceHeight / 2,
        faceWidth,
        faceHeight
    );

    const data = imageData.data;
    let totalBrightness = 0;
    let edgeCount = 0;
    let prevBrightness = 0;

    // Analyze pixels
    for (let i = 0; i < data.length; i += 4) {
        const r = data[i];
        const g = data[i + 1];
        const b = data[i + 2];

        // Calculate brightness (luminance)
        const brightness = 0.299 * r + 0.587 * g + 0.114 * b;
        totalBrightness += brightness;

        // Detect edges (sudden brightness changes)
        if (i > 0 && Math.abs(brightness - prevBrightness) > 12) {
            edgeCount++;
        }
        prevBrightness = brightness;
    }

    const pixelCount = data.length / 4;
    const avgBrightness = totalBrightness / pixelCount;
    const edgeRatio = edgeCount / pixelCount;

    // Calculate sharpness (more edges = sharper image = likely a real face)
    const sharpness = Math.min(edgeRatio * 10, 1);

    // Determine if face is present
    // Good lighting: brightness between 40-220 (more tolerant)
    // Good edges: edge ratio > 0.05 (more tolerant)
    const hasGoodLighting = avgBrightness >= 40 && avgBrightness <= 220;
    const hasGoodEdges = edgeRatio > 0.03; // (more tolerant, accommodates softer images)
    const hasFace = hasGoodLighting && hasGoodEdges;

    // Calculate confidence score (0-1)
    let confidence = 0;
    if (hasFace) {
        const lightingScore = Math.min(
            (avgBrightness - 40) / 120,
            (220 - avgBrightness) / 120
        );
        const edgeScore = Math.min(edgeRatio / 0.08, 1);
        confidence = lightingScore * 0.4 + edgeScore * 0.6;
    }

    return {
        hasFace,
        confidence: Math.max(0, Math.min(1, confidence)),
        brightness: avgBrightness,
        sharpness,
    };
};

/**
 * Check if image quality is good enough for face verification
 * @param detectionResult - Face detection result
 * @returns true if quality is acceptable
 */
export const hasGoodImageQuality = (
    detectionResult: FaceDetectionResult
): boolean => {
    return (
        detectionResult.hasFace &&
        detectionResult.confidence >= 0.25 &&
        detectionResult.brightness >= 45 &&
        detectionResult.brightness <= 210 &&
        detectionResult.sharpness >= 0.2
    );
};

/**
 * Get feedback message for user based on detection result
 * @param detectionResult - Face detection result
 * @returns User-friendly feedback message
 */
export const getDetectionFeedback = (
    detectionResult: FaceDetectionResult
): string => {
    if (!detectionResult.hasFace) {
        if (detectionResult.brightness < 40) {
            return 'Ánh sáng quá tối';
        }
        if (detectionResult.brightness > 230) {
            return 'Ánh sáng quá sáng';
        }
        return 'Không phát hiện khuôn mặt';
    }

    if (detectionResult.confidence < 0.4) {
        return 'Khuôn mặt không rõ';
    }

    if (detectionResult.sharpness < 0.3) {
        return 'Hình ảnh không sắc nét';
    }

    return 'Khuôn mặt phát hiện được';
};
