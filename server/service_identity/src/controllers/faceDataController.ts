import { Request, Response } from 'express';
import { faceDataService } from '../services/faceDataService';
import { userService } from '../services/userService';
import { sendSuccess, sendError } from '../utils/response';

export class FaceDataController {
  async getFaceDataByUserId(req: Request, res: Response) {
    try {
      const { user_id } = req.params;

      const user = await userService.getUserById(user_id);
      if (!user) {
        return sendError(res, 'User not found', 404, 'Not Found');
      }

      const faceDataList = await faceDataService.getFaceDataByUserId(user_id);
      return sendSuccess(res, faceDataList, 'Face data retrieved successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to retrieve face data');
    }
  }

  async createFaceData(req: Request, res: Response) {
    try {
      const { user_id } = req.params;
      const { image_url, face_encoding, quality_score } = req.body;

      if (!image_url) {
        return sendError(res, 'image_url is required', 400, 'Validation Error');
      }

      const user = await userService.getUserById(user_id);
      if (!user) {
        return sendError(res, 'User not found', 404, 'Not Found');
      }

      const faceData = await faceDataService.createFaceData(user_id, image_url, face_encoding, quality_score);
      return sendSuccess(res, faceData, 'Face data created successfully', 201);
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to create face data');
    }
  }

  async deleteFaceData(req: Request, res: Response) {
    try {
      const { user_id, fid } = req.params;

      const user = await userService.getUserById(user_id);
      if (!user) {
        return sendError(res, 'User not found', 404, 'Not Found');
      }

      const faceData = await faceDataService.getFaceDataById(fid, user_id);
      if (!faceData) {
        return sendError(res, 'Face data not found', 404, 'Not Found');
      }

      await faceDataService.deleteFaceData(fid, user_id);
      return sendSuccess(res, { fid }, 'Face data deleted successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to delete face data');
    }
  }
}

export const faceDataController = new FaceDataController();
