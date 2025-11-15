import { Request, Response, Express } from 'express';
import { userService } from '../services/userService';
import { sendSuccess, sendError } from '../utils/response';

export class UserController {
    async getAllUsers(req: Request, res: Response) {
        try {
            const { company_id } = req.query;
            const users = await userService.getAllUsers(company_id as string);
            return sendSuccess(res, users, 'Users retrieved successfully');
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to retrieve users'
            );
        }
    }

    async getUserById(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const user = await userService.getUserById(user_id);

            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            return sendSuccess(res, user, 'User retrieved successfully');
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to retrieve user'
            );
        }
    }

    async createUser(req: Request, res: Response) {
        try {
            const {
                email,
                phone,
                password,
                full_name,
                avatar_url,
                role,
                company_id,
                employee_code,
                department,
                position,
                hire_date,
                salary,
            } = req.body;

            // Validation
            if (
                !email ||
                !phone ||
                !password ||
                !full_name ||
                role === undefined
            ) {
                return sendError(
                    res,
                    'Missing required fields: email, phone, password, full_name, role',
                    400,
                    'Validation Error'
                );
            }

            // Check if user already exists
            const existingUser = await userService.getUserByEmail(email);
            if (existingUser) {
                return sendError(
                    res,
                    'User with this email already exists',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.createUser({
                email,
                phone,
                password,
                full_name,
                avatar_url,
                role,
                company_id,
                employee_code,
                department,
                position,
                hire_date,
                salary,
            });

            return sendSuccess(res, user, 'User created successfully', 201);
        } catch (error: any) {
            return sendError(res, error.message, 500, 'Failed to create user');
        }
    }

    async updateUser(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const {
                email,
                phone,
                full_name,
                avatar_url,
                status,
                department,
                position,
                salary,
            } = req.body;

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const updatedUser = await userService.updateUser(user_id, {
                email,
                phone,
                full_name,
                avatar_url,
                status,
                department,
                position,
                salary,
            });

            return sendSuccess(res, updatedUser, 'User updated successfully');
        } catch (error: any) {
            return sendError(res, error.message, 500, 'Failed to update user');
        }
    }

    async deleteUser(req: Request, res: Response) {
        try {
            const { user_id } = req.params;

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            await userService.deleteUser(user_id);
            return sendSuccess(res, { user_id }, 'User deleted successfully');
        } catch (error: any) {
            return sendError(res, error.message, 500, 'Failed to delete user');
        }
    }

    async updateBaseInfo(req: Request, res: Response) {
        try {
            const {
                user_id,
                user_fullname,
                user_phone,
                user_email,
                user_department,
                user_data_join_company,
                user_position,
            } = req.body;

            if (!user_id) {
                return sendError(
                    res,
                    'Missing required field: user_id',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const updatedUser = await userService.updateBaseInfo(user_id, {
                user_fullname,
                user_phone,
                user_email,
                user_department,
                user_data_join_company,
                user_position,
            });

            return sendSuccess(
                res,
                updatedUser,
                'User base info updated successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to update user base info'
            );
        }
    }

    async deleteListEmployee(req: Request, res: Response) {
        try {
            const { user_ids } = req.body;

            if (!Array.isArray(user_ids) || user_ids.length === 0) {
                return sendError(
                    res,
                    'user_ids must be a non-empty array',
                    400,
                    'Validation Error'
                );
            }

            const result = await userService.deleteListEmployee(user_ids);
            return sendSuccess(
                res,
                result,
                'Employees deleted successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to delete employees'
            );
        }
    }

    async importUsersFromFile(req: Request, res: Response) {
        try {
            const file = (req as Request & { file?: Express.Multer.File }).file;
            if (!file) {
                return sendError(
                    res,
                    'No file uploaded',
                    400,
                    'Validation Error'
                );
            }
            const result = await userService.importUsersFromFile(file);
            return sendSuccess(res, result, 'Users imported successfully');
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to import users from file'
            );
        }
    }
}

export const userController = new UserController();
