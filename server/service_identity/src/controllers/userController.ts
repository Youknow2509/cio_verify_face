import { Request, Response, Express } from 'express';
import { userService } from '../services/userService';
import { sendSuccess, sendError } from '../utils/response';
import { sendToKafka, getKafkaTopics } from '../config/kafka';

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

    async updateUserName(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { full_name } = req.body;

            if (!full_name) {
                return sendError(
                    res,
                    'Missing required field: full_name',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const oldName = user.full_name;
            const updatedUser = await userService.updateUserName(
                user_id,
                full_name
            );

            if (!updatedUser) {
                return sendError(
                    res,
                    'Failed to update user',
                    500,
                    'Update Error'
                );
            }

            const responseData = {
                user_id: updatedUser.user_id,
                old_name: oldName,
                new_name: updatedUser.full_name,
                email: updatedUser.email,
                updated_at: updatedUser.updated_at,
            };

            return sendSuccess(
                res,
                responseData,
                'User name updated successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to update user name'
            );
        }
    }

    async updateUserPhone(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { phone } = req.body;

            if (!phone) {
                return sendError(
                    res,
                    'Missing required field: phone',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const oldPhone = user.phone;
            const updatedUser = await userService.updateUserPhone(
                user_id,
                phone
            );

            if (!updatedUser) {
                return sendError(
                    res,
                    'Failed to update user',
                    500,
                    'Update Error'
                );
            }

            const responseData = {
                user_id: updatedUser.user_id,
                old_phone: oldPhone,
                new_phone: updatedUser.phone,
                email: updatedUser.email,
                updated_at: updatedUser.updated_at,
            };

            return sendSuccess(
                res,
                responseData,
                'User phone updated successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to update user phone'
            );
        }
    }

    async updateUserDepartment(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { department } = req.body;

            if (!department) {
                return sendError(
                    res,
                    'Missing required field: department',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const oldDepartment = user.department;
            const updatedUser = await userService.updateUserDepartment(
                user_id,
                department
            );

            if (!updatedUser) {
                return sendError(
                    res,
                    'Failed to update user',
                    500,
                    'Update Error'
                );
            }

            const responseData = {
                user_id: updatedUser.user_id,
                old_department: oldDepartment,
                new_department: department,
                full_name: updatedUser.full_name,
                email: updatedUser.email,
                updated_at: updatedUser.updated_at,
            };

            return sendSuccess(
                res,
                responseData,
                'User department updated successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to update user department'
            );
        }
    }

    async updateUserPosition(req: Request, res: Response) {
        try {
            const { user_id } = req.params;
            const { position } = req.body;

            if (!position) {
                return sendError(
                    res,
                    'Missing required field: position',
                    400,
                    'Validation Error'
                );
            }

            const user = await userService.getUserById(user_id);
            if (!user) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const oldPosition = user.position;
            const updatedUser = await userService.updateUserPosition(
                user_id,
                position
            );

            if (!updatedUser) {
                return sendError(
                    res,
                    'Failed to update user',
                    500,
                    'Update Error'
                );
            }

            const responseData = {
                user_id: updatedUser.user_id,
                old_position: oldPosition,
                new_position: position,
                full_name: updatedUser.full_name,
                email: updatedUser.email,
                updated_at: updatedUser.updated_at,
            };

            return sendSuccess(
                res,
                responseData,
                'User position updated successfully'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to update user position'
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
            return sendSuccess(res, result, 'Employees deleted successfully');
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

    async resetPassword(req: Request, res: Response) {
        try {
            const { user_id } = req.params;

            if (!user_id) {
                return sendError(
                    res,
                    'Missing required parameter: user_id',
                    400,
                    'Validation Error'
                );
            }

            const result = await userService.resetPassword(user_id);

            if (!result) {
                return sendError(res, 'User not found', 404, 'Not Found');
            }

            const { newPassword, user } = result;

            // Send to Kafka
            const kafkaTopics = getKafkaTopics();
            const kafkaMessage = {
                event: 'user.password_reset',
                user_id: user.user_id,
                email: user.email,
                full_name: user.full_name,
                new_password: newPassword,
                timestamp: new Date().toISOString(),
            };

            try {
                await sendToKafka(kafkaTopics.userEvents, [kafkaMessage]);
            } catch (kafkaError) {
                console.error('Failed to send to Kafka:', kafkaError);
                // Continue anyway - don't fail the request if Kafka is down
            }

            // Don't return password in response - only confirmation
            const responseData = {
                user_id: user.user_id,
                email: user.email,
                full_name: user.full_name,
                reset_at: user.updated_at,
            };

            return sendSuccess(
                res,
                responseData,
                'Password reset successfully. New password has been sent to user via Kafka notification.'
            );
        } catch (error: any) {
            return sendError(
                res,
                error.message,
                500,
                'Failed to reset password'
            );
        }
    }
}

export const userController = new UserController();
