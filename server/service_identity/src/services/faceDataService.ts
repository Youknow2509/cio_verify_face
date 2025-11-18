import { query } from '../config/database';
import { FaceProfile } from '../types';
import { v4 as uuidv4 } from 'uuid';
import { faceVerificationClient } from '../grpc/faceVerificationClient';

// Utility: ensure embedding array length is 512 and convert to pgvector literal
function prepareEmbedding(embedding?: number[]): string {
    if (!embedding || embedding.length === 0) {
        throw new Error('Embedding missing from AI response');
    }
    if (embedding.length !== 512) {
        throw new Error(`Embedding length ${embedding.length} != 512`);
    }
    // pgvector accepts ARRAY casting like '[e1,e2,...]' as text
    return (
        '[' +
        embedding
            .map((v) => (Number.isFinite(v) ? v.toFixed(6) : '0'))
            .join(',') +
        ']'
    );
}

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
        const aiResp = await faceVerificationClient.enrollFace({
            user_id: params.user_id,
            company_id: params.company_id,
            image_data: params.imageBuffer,
            make_primary: params.make_primary || false,
            metadata: params.metadata || {},
        });

        if (aiResp.status !== 'ok') {
            throw new Error(
                aiResp.message || `AI enroll failed: ${aiResp.status}`
            );
        }

        const profile_id = aiResp.profile_id || uuidv4();
        const now = new Date().toISOString();

        // If make_primary true, unset previous primary
        if (params.make_primary) {
            await query(
                `UPDATE face_profiles SET is_primary = false, updated_at = $1
         WHERE user_id = $2 AND company_id = $3 AND is_primary = true AND deleted_at IS NULL`,
                [now, params.user_id, params.company_id]
            );
        }

        const embeddingLiteral = prepareEmbedding(aiResp.embedding);
        const embedding_version = aiResp.embedding_version || 'v1';
        const quality_score = aiResp.quality_score || null;
        const meta_data = JSON.stringify(params.metadata || {});

        const insertResult = await query(
            `INSERT INTO face_profiles (
         profile_id, user_id, company_id, embedding, embedding_version,
         enroll_image_path, is_primary, quality_score, meta_data,
         created_at, updated_at, indexed, index_version
       ) VALUES (
         $1, $2, $3, $4::vector, $5, $6, $7, $8, $9::jsonb, $10, $11, $12, $13
       ) RETURNING profile_id, user_id, company_id, embedding_version, enroll_image_path, is_primary,
         quality_score, meta_data, created_at, updated_at, deleted_at, indexed, index_version`,
            [
                profile_id,
                params.user_id,
                params.company_id,
                embeddingLiteral,
                embedding_version,
                params.enroll_image_path || null,
                params.make_primary || false,
                quality_score,
                meta_data,
                now,
                now,
                false,
                0,
            ]
        );

        return insertResult.rows[0];
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
