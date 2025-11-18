// Face profile types aligned with face_profiles table
export interface FaceProfile {
    profile_id: string;
    user_id: string;
    company_id?: string;
    embedding_version: string;
    enroll_image_path?: string; // raw storage path
    is_primary: boolean;
    quality_score?: number;
    meta_data: Record<string, any>;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
    indexed: boolean;
    index_version: number;
    // Convenience fields from API (optional)
    image_url?: string; // resolved public URL for enroll_image_path
}

export interface FaceProfileUploadResponse {
    success: boolean;
    data: FaceProfile[];
    message?: string;
}

export interface FaceProfileListResponse {
    success: boolean;
    data: FaceProfile[];
    message?: string;
}

export interface FaceProfileActionResponse {
    success: boolean;
    data?: FaceProfile;
    message?: string;
}
