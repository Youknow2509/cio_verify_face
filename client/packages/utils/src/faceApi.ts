import { apiClient } from './api';
import {
    FaceProfile,
    FaceProfileUploadResponse,
    FaceProfileListResponse,
    FaceProfileActionResponse,
} from '@repo/types';

// List face profiles for a user
export const listFaceProfiles = async (
    userId: string,
    companyId: string
): Promise<FaceProfile[]> => {
    const { data } = await apiClient.get<FaceProfileListResponse>(
        `/api/v1/users/${userId}/face-data?company_id=${companyId}`
    );
    return data.data || [];
};

// Upload face images (multipart)
export const uploadFaceProfiles = async (
    userId: string,
    companyId: string,
    files: FileList | File[]
): Promise<FaceProfile[]> => {
    const formData = new FormData();
    formData.append('company_id', companyId);

    // Convert FileList to array with proper typing
    const fileArray: File[] = Array.isArray(files) ? files : Array.from(files);

    fileArray.forEach((file) => formData.append('image', file));

    const { data } = await apiClient.post<FaceProfileUploadResponse>(
        `/api/v1/users/${userId}/face-data/upload`,
        formData,
        {
            headers: { 'Content-Type': 'multipart/form-data' },
        }
    );
    const payload: any = (data as any)?.data;
    if (Array.isArray(payload)) return payload as FaceProfile[];
    if (payload) return [payload as FaceProfile];
    return [];
};

// Set a profile as primary
export const setPrimaryFaceProfile = async (
    userId: string,
    profileId: string,
    companyId: string,
    status: boolean = true
): Promise<FaceProfile | undefined> => {
    const { data } = await apiClient.put<FaceProfileActionResponse>(
        `/api/v1/users/${userId}/face-data/${profileId}/primary`,
        {
            company_id: companyId,
            status: status,
        }
    );
    return data.data;
};

// Delete (soft) face profile
export const deleteFaceProfile = async (
    userId: string,
    profileId: string,
    companyId: string
): Promise<boolean> => {
    const { data } = await apiClient.delete<FaceProfileActionResponse>(
        `/api/v1/users/${userId}/face-data/${profileId}`,
        {
            params: { company_id: companyId },
        }
    );
    return data.success;
};
