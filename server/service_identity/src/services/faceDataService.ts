import { query } from '../config/database';
import { FaceProfile } from '../types';
import {
    faceVerificationClient,
    EnrollParams,
    EnrollResult,
} from '../grpc/faceVerificationClient';

export class FaceDataService {
    async listProfiles(
        userId: string,
        companyId: string
    ): Promise<FaceProfile[]> {
        const result = await query(
            `SELECT 
         profile_id, user_id, company_id, embedding_version, enroll_image_path, is_primary,
         quality_score, meta_data, created_at, updated_at, deleted_at, indexed, index_version
       FROM face_profiles
       WHERE user_id = $1 AND company_id = $2 AND deleted_at IS NULL
       ORDER BY is_primary DESC, created_at DESC`,
            [userId, companyId]
        );
        return result.rows;
    }

    async getProfile(
        profileId: string,
        companyId: string
    ): Promise<FaceProfile | null> {
        const result = await query(
            `SELECT 
         profile_id, user_id, company_id, embedding_version, enroll_image_path, is_primary,
         quality_score, meta_data, created_at, updated_at, deleted_at, indexed, index_version
       FROM face_profiles
       WHERE profile_id = $1 AND company_id = $2 AND deleted_at IS NULL`,
            [profileId, companyId]
        );
        return result.rows[0] || null;
    }

    async enrollProfile(params: {
        user_id: string;
        company_id: string;
        imageBuffer: Buffer;
        make_primary?: boolean;
        metadata?: Record<string, string>;
        enroll_image_path?: string; // persisted path or URL after upload
    }): Promise<FaceProfile> {
        // Call gRPC AI service to handle face enrollment
        const enrollParams: EnrollParams = {
            user_id: params.user_id,
            company_id: params.company_id,
            image_data: params.imageBuffer,
            make_primary: params.make_primary || false,
            metadata: params.metadata || {},
        };

        // Use enrollFace instead of enrollFaceStream if stream not implemented on server
        const resp: EnrollResult = await faceVerificationClient.enrollFace(
            enrollParams
        );

        // Check response
        if (
            resp.status !== 'ok' ||
            !resp.profile_id ||
            resp.profile_id === ''
        ) {
            // Return error with more details
            return Promise.reject(new Error(`${resp.message || resp.status}`));
        }

        // Call db get data profile id
        const res: FaceProfile = {
            profile_id: resp.profile_id,
            user_id: params.user_id,
            company_id: params.company_id,
            enroll_image_path: params.enroll_image_path || '',
            embedding_version: '',
            is_primary: params.make_primary || false,
            quality_score: resp.quality_score || 0,
            meta_data: params.metadata || {},
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            indexed: false,
            index_version: 0,
        };
        return Promise.resolve(res);
    }

    async deleteProfile(
        profileId: string,
        companyId: string,
        hardDelete = false
    ): Promise<boolean> {
        if (hardDelete) {
            const result = await query(
                `DELETE FROM face_profiles WHERE profile_id = $1 AND company_id = $2`,
                [profileId, companyId]
            );
            return result.rowCount > 0;
        }
        const now = new Date().toISOString();
        const result = await query(
            `UPDATE face_profiles SET deleted_at = $1, updated_at = $1 WHERE profile_id = $2 AND company_id = $3 AND deleted_at IS NULL`,
            [now, profileId, companyId]
        );
        return result.rowCount > 0;
    }

    async updatePrimaryProfile(
        profileId: string,
        userId: string,
        companyId: string,
        status: boolean
    ): Promise<boolean> {
        const now = new Date().toISOString();

        if (status) {
            // Set this profile as primary, unset others
            await query(
                `UPDATE face_profiles 
             SET is_primary = false, updated_at = $1
             WHERE user_id = $2 AND company_id = $3 AND is_primary = true AND deleted_at IS NULL`,
                [now, userId, companyId]
            );

            const result = await query(
                `UPDATE face_profiles
             SET is_primary = true, updated_at = $1
             WHERE profile_id = $2 AND user_id = $3 AND company_id = $4 AND deleted_at IS NULL`,
                [now, profileId, userId, companyId]
            );
            return result.rowCount > 0;
        } else {
            // Unset primary: check if this is the only primary profile BEFORE unsetting
            const primaryCheck = await query(
                `SELECT COUNT(*) as count 
             FROM face_profiles 
             WHERE user_id = $1 AND company_id = $2 AND is_primary = true AND deleted_at IS NULL`,
                [userId, companyId]
            );

            const primaryCount = parseInt(primaryCheck.rows[0].count, 10);

            if (primaryCount <= 1) {
                // Don't allow unsetting the last primary profile
                throw new Error(
                    'Cannot unset the last primary profile. Set another profile as primary first.'
                );
            }

            const result = await query(
                `UPDATE face_profiles
             SET is_primary = false, updated_at = $1
             WHERE profile_id = $2 AND user_id = $3 AND company_id = $4 AND deleted_at IS NULL`,
                [now, profileId, userId, companyId]
            );
            return result.rowCount > 0;
        }
    }
}

export const faceDataService = new FaceDataService();
