import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

// Wrapper client for FaceVerification service
const PROTO_PATH = path.join(
    __dirname,
    '..',
    'proto',
    'face_verification.proto'
);

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const faceProto: any =
    grpc.loadPackageDefinition(packageDefinition).face_verification;

// Enroll interfaces
export interface EnrollParams {
    user_id: string;
    company_id: string;
    image_data: Buffer;
    device_id?: string;
    make_primary?: boolean;
    metadata?: Record<string, string>;
}

export interface DuplicateProfile {
    user_id: string;
    similarity: number;
}

export interface EnrollResult {
    status: string;
    profile_id?: string;
    message?: string;
    duplicate_profiles?: DuplicateProfile[];
    quality_score?: number;
}

// Verify interfaces
export interface VerifyParams {
    image_data: Buffer;
    company_id: string;
    user_id?: string;
    device_id?: string;
    search_mode: '1:1' | '1:N';
    top_k?: number;
}

export interface Match {
    user_id: string;
    profile_id: string;
    similarity: number;
    confidence: number;
    is_primary: boolean;
}

export interface VerifyResult {
    status: string;
    verified: boolean;
    matches?: Match[];
    best_match?: Match;
    message?: string;
    liveness_score?: number;
}

// Multi-frame verify interfaces
export interface VerifyMultiFrameParams {
    frames: Buffer[];
    company_id: string;
    user_id?: string;
    device_id?: string;
    search_mode: '1:1' | '1:N';
    top_k?: number;
}

// Profile management interfaces
export interface GetProfilesParams {
    user_id: string;
    company_id: string;
}

export interface FaceProfile {
    profile_id: string;
    user_id: string;
    company_id: string;
    embedding_version: string;
    is_primary: boolean;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
    metadata?: Record<string, string>;
    quality_score?: number;
}

export interface GetProfilesResult {
    profiles: FaceProfile[];
}

export interface UpdateProfileParams {
    profile_id: string;
    company_id: string;
    image_data?: Buffer;
    make_primary?: boolean;
    metadata?: Record<string, string>;
}

export interface UpdateProfileResult {
    status: string;
    message: string;
}

export interface DeleteProfileParams {
    profile_id: string;
    company_id: string;
    hard_delete?: boolean;
}

export interface DeleteProfileResult {
    status: string;
    message: string;
}

class FaceVerificationClient {
    private client: grpc.Client;

    constructor(address: string) {
        this.client = new faceProto.FaceVerification(
            address,
            grpc.credentials.createInsecure()
        );
    }

    /**
     * Enroll a face with single image
     */
    enrollFace(params: EnrollParams): Promise<EnrollResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).enrollFace(
                {
                    user_id: params.user_id,
                    company_id: params.company_id,
                    image_data: params.image_data,
                    device_id: params.device_id,
                    make_primary: params.make_primary || false,
                    metadata: params.metadata || {},
                },
                (err: grpc.ServiceError, response: EnrollResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Enroll a face with streaming image (optimized for bandwidth)
     */
    enrollFaceStream(params: EnrollParams): Promise<EnrollResult> {
        return new Promise((resolve, reject) => {
            const call = (this.client as any).enrollFaceStream(
                (err: grpc.ServiceError, response: EnrollResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );

            // Send metadata first
            call.write({
                metadata: {
                    user_id: params.user_id,
                    company_id: params.company_id,
                    device_id: params.device_id,
                    make_primary: params.make_primary || false,
                    metadata: params.metadata || {},
                    total_size: params.image_data.length,
                    image_format: 'JPEG',
                },
            });

            // Stream image in chunks (e.g., 64KB chunks)
            const CHUNK_SIZE = 64 * 1024;
            for (let i = 0; i < params.image_data.length; i += CHUNK_SIZE) {
                const chunk = params.image_data.slice(i, i + CHUNK_SIZE);
                call.write({ image_chunk: chunk });
            }

            call.end();
        });
    }

    /**
     * Verify face with single image (1:1 or 1:N)
     */
    verifyFace(params: VerifyParams): Promise<VerifyResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).verifyFace(
                {
                    image_data: params.image_data,
                    company_id: params.company_id,
                    user_id: params.user_id,
                    device_id: params.device_id,
                    search_mode: params.search_mode,
                    top_k: params.top_k || 10,
                },
                (err: grpc.ServiceError, response: VerifyResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Verify face with streaming image (optimized for bandwidth)
     */
    verifyFaceStream(params: VerifyParams): Promise<VerifyResult> {
        return new Promise((resolve, reject) => {
            const call = (this.client as any).verifyFaceStream(
                (err: grpc.ServiceError, response: VerifyResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );

            // Send metadata first
            call.write({
                metadata: {
                    company_id: params.company_id,
                    user_id: params.user_id,
                    device_id: params.device_id,
                    search_mode: params.search_mode,
                    top_k: params.top_k || 10,
                    total_size: params.image_data.length,
                    image_format: 'JPEG',
                },
            });

            // Stream image in chunks
            const CHUNK_SIZE = 64 * 1024;
            for (let i = 0; i < params.image_data.length; i += CHUNK_SIZE) {
                const chunk = params.image_data.slice(i, i + CHUNK_SIZE);
                call.write({ image_chunk: chunk });
            }

            call.end();
        });
    }

    /**
     * Verify face with multiple frames (more robust)
     */
    verifyFaceMultiFrame(
        params: VerifyMultiFrameParams
    ): Promise<VerifyResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).verifyFaceMultiFrame(
                {
                    frames: params.frames,
                    company_id: params.company_id,
                    user_id: params.user_id,
                    device_id: params.device_id,
                    search_mode: params.search_mode,
                    top_k: params.top_k || 10,
                },
                (err: grpc.ServiceError, response: VerifyResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Verify face with multiple frames using streaming (optimized for bandwidth)
     */
    verifyFaceMultiFrameStream(
        params: VerifyMultiFrameParams
    ): Promise<VerifyResult> {
        return new Promise((resolve, reject) => {
            const call = (this.client as any).verifyFaceMultiFrameStream(
                (err: grpc.ServiceError, response: VerifyResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );

            // Send metadata first
            call.write({
                metadata: {
                    company_id: params.company_id,
                    user_id: params.user_id,
                    device_id: params.device_id,
                    search_mode: params.search_mode,
                    top_k: params.top_k || 10,
                    frame_count: params.frames.length,
                    image_format: 'JPEG',
                },
            });

            // Stream each frame
            const CHUNK_SIZE = 64 * 1024;
            params.frames.forEach((frame, frameIndex) => {
                // Send frame delimiter
                call.write({
                    frame_delimiter: {
                        frame_number: frameIndex,
                        frame_size: frame.length,
                    },
                });

                // Stream frame in chunks
                for (let i = 0; i < frame.length; i += CHUNK_SIZE) {
                    const chunk = frame.slice(i, i + CHUNK_SIZE);
                    call.write({ frame_chunk: chunk });
                }
            });

            call.end();
        });
    }

    /**
     * Get user face profiles
     */
    getUserProfiles(params: GetProfilesParams): Promise<GetProfilesResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).getUserProfiles(
                {
                    user_id: params.user_id,
                    company_id: params.company_id,
                },
                (err: grpc.ServiceError, response: GetProfilesResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Update face profile
     */
    updateProfile(params: UpdateProfileParams): Promise<UpdateProfileResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).updateProfile(
                {
                    profile_id: params.profile_id,
                    company_id: params.company_id,
                    image_data: params.image_data,
                    make_primary: params.make_primary,
                    metadata: params.metadata || {},
                },
                (err: grpc.ServiceError, response: UpdateProfileResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Update face profile with streaming image (optimized for bandwidth)
     */
    updateProfileStream(
        params: UpdateProfileParams
    ): Promise<UpdateProfileResult> {
        return new Promise((resolve, reject) => {
            const call = (this.client as any).updateProfileStream(
                (err: grpc.ServiceError, response: UpdateProfileResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );

            const hasImage = !!params.image_data;

            // Send metadata first
            call.write({
                metadata: {
                    profile_id: params.profile_id,
                    company_id: params.company_id,
                    make_primary: params.make_primary,
                    metadata: params.metadata || {},
                    has_image: hasImage,
                    total_size: hasImage ? params.image_data!.length : 0,
                    image_format: hasImage ? 'JPEG' : '',
                },
            });

            // Stream image in chunks if provided
            if (hasImage && params.image_data) {
                const CHUNK_SIZE = 64 * 1024;
                for (let i = 0; i < params.image_data.length; i += CHUNK_SIZE) {
                    const chunk = params.image_data.slice(i, i + CHUNK_SIZE);
                    call.write({ image_chunk: chunk });
                }
            }

            call.end();
        });
    }

    /**
     * Delete face profile
     */
    deleteProfile(params: DeleteProfileParams): Promise<DeleteProfileResult> {
        return new Promise((resolve, reject) => {
            (this.client as any).deleteProfile(
                {
                    profile_id: params.profile_id,
                    company_id: params.company_id,
                    hard_delete: params.hard_delete || false,
                },
                (err: grpc.ServiceError, response: DeleteProfileResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }

    /**
     * Close the gRPC connection
     */
    close(): void {
        this.client.close();
    }
}

// Singleton instance; address should be configured via env FACE_AI_GRPC_ADDR (e.g. "localhost:50051")
const grpcAddress = process.env.FACE_AI_GRPC_ADDR || 'localhost:50051';
export const faceVerificationClient = new FaceVerificationClient(grpcAddress);
