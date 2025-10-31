import { query } from '../config/database';
import { User, CreateUserRequest, UpdateUserRequest, Employee } from '../types';
import { v4 as uuidv4 } from 'uuid';
import { hashPassword, generateSalt } from '../utils/crypto';

export class UserService {
  async getAllUsers(companyId?: string): Promise<User[]> {
    if (companyId) {
      const result = await query(
        `SELECT u.* FROM users u 
         INNER JOIN employees e ON u.user_id = e.employee_id 
         WHERE e.company_id = $1 
         ORDER BY u.created_at DESC`,
        [companyId]
      );
      return result.rows;
    }

    const result = await query('SELECT * FROM users ORDER BY created_at DESC');
    return result.rows;
  }

  async getUserById(userId: string): Promise<User | null> {
    const result = await query('SELECT * FROM users WHERE user_id = $1', [userId]);
    return result.rows[0] || null;
  }

  async getUserByEmail(email: string): Promise<User | null> {
    const result = await query('SELECT * FROM users WHERE email = $1', [email]);
    return result.rows[0] || null;
  }

  async createUser(data: CreateUserRequest): Promise<User> {
    const userId = uuidv4();
    const salt = generateSalt();
    const passwordHash = hashPassword(data.password, salt);
    const now = new Date().toISOString();

    const result = await query(
      `INSERT INTO users 
       (user_id, email, phone, salt, password_hash, full_name, avatar_url, role, status, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
       RETURNING *`,
      [
        userId,
        data.email,
        data.phone,
        salt,
        passwordHash,
        data.full_name,
        data.avatar_url || null,
        data.role,
        0, // status: ACTIVE
        now,
        now,
      ]
    );

    const user = result.rows[0];

    // If company_id is provided and role is EMPLOYEE, create employee record
    if (data.company_id && data.role === 2) {
      const employeeId = uuidv4();
      await query(
        `INSERT INTO employees 
         (employee_id, company_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
        [
          userId,
          data.company_id,
          data.employee_code || `EMP-${Date.now()}`,
          data.department || null,
          data.position || null,
          data.hire_date || null,
          data.salary || null,
          0, // status: active
          now,
          now,
        ]
      );
    }

    return user;
  }

  async updateUser(userId: string, data: UpdateUserRequest): Promise<User | null> {
    const now = new Date().toISOString();
    const updates: string[] = [];
    const values: any[] = [];
    let paramIndex = 1;

    if (data.email !== undefined) {
      updates.push(`email = $${paramIndex++}`);
      values.push(data.email);
    }
    if (data.phone !== undefined) {
      updates.push(`phone = $${paramIndex++}`);
      values.push(data.phone);
    }
    if (data.full_name !== undefined) {
      updates.push(`full_name = $${paramIndex++}`);
      values.push(data.full_name);
    }
    if (data.avatar_url !== undefined) {
      updates.push(`avatar_url = $${paramIndex++}`);
      values.push(data.avatar_url);
    }
    if (data.status !== undefined) {
      updates.push(`status = $${paramIndex++}`);
      values.push(data.status);
    }

    if (updates.length === 0) {
      return this.getUserById(userId);
    }

    updates.push(`updated_at = $${paramIndex++}`);
    values.push(now);
    values.push(userId);

    const result = await query(
      `UPDATE users SET ${updates.join(', ')} WHERE user_id = $${paramIndex} RETURNING *`,
      values
    );

    const user = result.rows[0] || null;

    // If department, position, or salary is updated, update employee record too
    if (user && (data.department !== undefined || data.position !== undefined || data.salary !== undefined)) {
      const empUpdates: string[] = [];
      const empValues: any[] = [];
      let empParamIndex = 1;

      if (data.department !== undefined) {
        empUpdates.push(`department = $${empParamIndex++}`);
        empValues.push(data.department);
      }
      if (data.position !== undefined) {
        empUpdates.push(`position = $${empParamIndex++}`);
        empValues.push(data.position);
      }
      if (data.salary !== undefined) {
        empUpdates.push(`salary = $${empParamIndex++}`);
        empValues.push(data.salary);
      }

      empUpdates.push(`updated_at = $${empParamIndex++}`);
      empValues.push(now);
      empValues.push(userId);

      await query(
        `UPDATE employees SET ${empUpdates.join(', ')} WHERE employee_id = $${empParamIndex}`,
        empValues
      );
    }

    return user;
  }

  async deleteUser(userId: string): Promise<boolean> {
    const result = await query('DELETE FROM users WHERE user_id = $1', [userId]);
    return result.rowCount > 0;
  }
}

export const userService = new UserService();
