import { Request, Response } from 'express';
import { userService } from '../services/userService';
import { sendSuccess, sendError } from '../utils/response';

export class UserController {
  async getAllUsers(req: Request, res: Response) {
    try {
      const { company_id } = req.query;
      const users = await userService.getAllUsers(company_id as string);
      return sendSuccess(res, users, 'Users retrieved successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to retrieve users');
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
      return sendError(res, error.message, 500, 'Failed to retrieve user');
    }
  }

  async createUser(req: Request, res: Response) {
    try {
      const { email, phone, password, full_name, avatar_url, role, company_id, employee_code, department, position, hire_date, salary } = req.body;

      // Validation
      if (!email || !phone || !password || !full_name || role === undefined) {
        return sendError(res, 'Missing required fields: email, phone, password, full_name, role', 400, 'Validation Error');
      }

      // Check if user already exists
      const existingUser = await userService.getUserByEmail(email);
      if (existingUser) {
        return sendError(res, 'User with this email already exists', 400, 'Validation Error');
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
      const { email, phone, full_name, avatar_url, status, department, position, salary } = req.body;

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
}

export const userController = new UserController();
