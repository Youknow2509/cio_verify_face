import { query } from '../config/database';
import { FaceData } from '../types';
import { v4 as uuidv4 } from 'uuid';

export class FaceDataService {
  async getFaceDataByUserId(userId: string): Promise<FaceData[]> {
    const result = await query(
      `SELECT * FROM face_data WHERE user_id = $1 ORDER BY created_at DESC`,
      [userId]
    );
    return result.rows;
  }

  async getFaceDataById(faceDataId: string, userId: string): Promise<FaceData | null> {
    const result = await query(
      `SELECT * FROM face_data WHERE fid = $1 AND user_id = $2`,
      [faceDataId, userId]
    );
    return result.rows[0] || null;
  }

  async createFaceData(
    userId: string,
    imageUrl: string,
    faceEncoding?: string,
    qualityScore?: number
  ): Promise<FaceData> {
    const faceDataId = uuidv4();
    const now = new Date().toISOString();

    const result = await query(
      `INSERT INTO face_data 
       (fid, user_id, image_url, face_encoding, quality_score, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6, $7)
       RETURNING *`,
      [faceDataId, userId, imageUrl, faceEncoding || null, qualityScore || null, now, now]
    );

    return result.rows[0];
  }

  async deleteFaceData(faceDataId: string, userId: string): Promise<boolean> {
    const result = await query(
      `DELETE FROM face_data WHERE fid = $1 AND user_id = $2`,
      [faceDataId, userId]
    );
    return result.rowCount > 0;
  }
}

export const faceDataService = new FaceDataService();
