import { Request, Response } from 'express';
import { faceDataService } from '../services/faceDataService';
import { userService } from '../services/userService';
import { sendSuccess, sendError } from '../utils/response';
import fetch from 'node-fetch';

async function fetchImageBuffer(image_url: string): Promise<Buffer> {
    const resp = await fetch(image_url);
    if (!resp.ok) throw new Error(`Download image failed: ${resp.status}`);
    const arrayBuffer = await resp.arrayBuffer();
    return Buffer.from(arrayBuffer);
}

export class FaceDataController {
    async getFaceDataByUserId(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { company_id } = req.query as { company_id?: string };

            if (!company_id) {
                return sendError(
                    res,
                    'company_id is required',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const profiles = await faceDataService.listProfiles(
                user_id,
                company_id
            );
            return sendSuccess(
                res,
                profiles,
                'Face profiles retrieved successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to retrieve face profiles'
            );
        }
    }

    async createFaceData(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { image_url, company_id, make_primary, metadata } =
                req.body as {
                    image_url?: string;
                    company_id?: string;
                    make_primary?: boolean;
                    metadata?: Record<string, string>;
                };

            if (!image_url) {
                return sendError(
                    res,
                    'image_url is required',
                    400,
                    'Validation Error'
                );
            }
            if (!company_id) {
                return sendError(
                    res,
                    'company_id is required',
                    400,
                    'Validation Error'
                );
            }
            if (!/^https?:\/\//i.test(image_url)) {
                return sendError(
                    res,
                    'image_url must start with http or https',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const imgBuffer = await fetchImageBuffer(image_url);

            const profile = await faceDataService.enrollProfile({
                user_id,
                company_id,
                imageBuffer: imgBuffer,
                make_primary: make_primary,
                metadata,
                enroll_image_path: image_url,
            });
            return sendSuccess(
                res,
                profile,
                'Face profile enrolled successfully',
                201
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to enroll face profile'
            );
        }
    }

    async deleteFaceData(req: Request, res: Response) {
        try {
            const { user_id, fid } = req.params; // fid becomes profile_id
            const { company_id, hard } = req.query as {
                company_id?: string;
                hard?: string;
            };

            if (!company_id) {
                return sendError(
                    res,
                    'company_id is required',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const profile = await faceDataService.getProfile(fid, company_id);
            if (!profile) {
                return sendError(
                    res,
                    'Face profile not found',
                    404,
                    'Not Found'
                );
            }

            await faceDataService.deleteProfile(
                fid,
                company_id,
                hard === 'true'
            );
            return sendSuccess(
                res,
                { profile_id: fid },
                'Face profile deleted successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to delete face profile'
            );
        }
    }
}

export const faceDataController = new FaceDataController();
