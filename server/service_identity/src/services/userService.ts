import { query } from '../config/database';
import { User, CreateUserRequest, UpdateUserRequest, Employee } from '../types';
import { v4 as uuidv4 } from 'uuid';
import { hashPassword, generateSalt } from '../utils/crypto';
import { parse } from 'csv-parse';
import fs from 'fs';

export class UserService {
    async getAllUsers(companyId?: string): Promise<User[]> {
        if (companyId) {
            const result = await query(
                `SELECT 
                u.user_id,
                u.email,
                u.phone,
                u.full_name,
                u.avatar_url,
                u.role,
                u.status,
                u.created_at,
                u.updated_at,
                e.employee_code,
                e.department,
                e.hire_date,
                e.position,
                e.company_id,
                COALESCE(
                    (SELECT COUNT(*) 
                     FROM face_profiles fp 
                     WHERE fp.user_id = u.user_id 
                     AND fp.company_id = e.company_id 
                     AND fp.deleted_at IS NULL), 0
                ) as face_data_count
            FROM users u 
            INNER JOIN employees e ON u.user_id = e.employee_id 
            WHERE e.company_id = $1 
            ORDER BY u.created_at DESC`,
                [companyId]
            );
            return result.rows;
        }

        const result = await query(
            `SELECT 
            u.*,
            0 as face_data_count
        FROM users u 
        ORDER BY u.created_at DESC`
        );
        return result.rows;
    }

    async getUserById(userId: string): Promise<User | null> {
        const result = await query(
            `SELECT 
            u.user_id,
            u.role,
            u.avatar_url,
            u.email,
            u.phone,
            u.full_name,
            e.employee_code,
            e.department,
            e.hire_date,
            e.position,
            u.status 
            FROM users u
            INNER JOIN employees e ON u.user_id = e.employee_id
            WHERE u.user_id = $1`,
            [userId]
        );
        return result.rows[0] || null;
    }

    async getCountProfileFaceUseer(
        userId: string,
        company_id: string
    ): Promise<number> {
        const sql_raw = `
        SELECT COUNT(*) AS count
        FROM face_profiles
        WHERE user_id = $1 
        AND company_id = $2
        AND deleted_at IS NULL
    `;

        const result = await query(sql_raw, [userId, company_id]);
        return parseInt(result.rows[0].count, 10);
    }

    async getFaceProfilesByUser(
        userId: string,
        company_id: string
    ): Promise<any[]> {
        const sql_raw = `
        SELECT 
            profile_id,
            user_id,
            company_id,
            embedding_version,
            enroll_image_path,
            is_primary,
            quality_score,
            meta_data,
            created_at,
            updated_at,
            indexed,
            index_version
        FROM face_profiles
        WHERE user_id = $1 
        AND company_id = $2
        AND deleted_at IS NULL
        ORDER BY is_primary DESC, created_at DESC
    `;

        const result = await query(sql_raw, [userId, company_id]);
        return result.rows;
    }

    async getPrimaryFaceProfile(
        userId: string,
        company_id: string
    ): Promise<any | null> {
        const sql_raw = `
        SELECT 
            profile_id,
            user_id,
            company_id,
            embedding_version,
            enroll_image_path,
            is_primary,
            quality_score,
            meta_data,
            created_at,
            updated_at,
            indexed,
            index_version
        FROM face_profiles
        WHERE user_id = $1 
        AND company_id = $2
        AND is_primary = true
        AND deleted_at IS NULL
        LIMIT 1
    `;

        const result = await query(sql_raw, [userId, company_id]);
        return result.rows[0] || null;
    }

    async getFaceProfileById(
        profileId: string,
        company_id: string
    ): Promise<any | null> {
        const sql_raw = `
        SELECT 
            profile_id,
            user_id,
            company_id,
            embedding_version,
            enroll_image_path,
            is_primary,
            quality_score,
            meta_data,
            created_at,
            updated_at,
            indexed,
            index_version
        FROM face_profiles
        WHERE profile_id = $1 
        AND company_id = $2
        AND deleted_at IS NULL
    `;

        const result = await query(sql_raw, [profileId, company_id]);
        return result.rows[0] || null;
    }

    async getUserByEmail(email: string): Promise<User | null> {
        const result = await query('SELECT * FROM users WHERE email = $1', [
            email,
        ]);
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
        console.log('Company id', data.company_id);
        console.log('Role', data.role);
        if (data.company_id) {
            const result = await query(
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

    async updateUser(
        userId: string,
        data: UpdateUserRequest
    ): Promise<User | null> {
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
            `UPDATE users SET ${updates.join(
                ', '
            )} WHERE user_id = $${paramIndex} RETURNING *`,
            values
        );

        const user = result.rows[0] || null;

        // If department, position, or salary is updated, update employee record too
        if (
            user &&
            (data.department !== undefined ||
                data.position !== undefined ||
                data.salary !== undefined)
        ) {
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
                `UPDATE employees SET ${empUpdates.join(
                    ', '
                )} WHERE employee_id = $${empParamIndex}`,
                empValues
            );
        }

        return user;
    }

    async deleteUser(userId: string): Promise<boolean> {
        const deleteEmployeeResult = await query(
            'DELETE FROM employees WHERE employee_id = $1',
            [userId]
        );
        const result = await query('DELETE FROM users WHERE user_id = $1', [
            userId,
        ]);
        return result.rowCount > 0;
    }

    async updateBaseInfo(
        userId: string,
        baseInfo: {
            user_fullname?: string;
            user_phone?: string;
            user_email?: string;
            user_department?: string;
            user_data_join_company?: string;
            user_position?: string;
        }
    ): Promise<User | null> {
        const queryUsers = `
            UPDATE users SET
                full_name = COALESCE($1, full_name),
                phone = COALESCE($2, phone),
                email = COALESCE($3, email)
            WHERE user_id = $4
            RETURNING *;
        `;
        const queryEmployees = `
            UPDATE employees SET
                department = COALESCE($1, department),
                hire_date = COALESCE($2, hire_date),
                position = COALESCE($3, position)
            WHERE employee_id = $4;
        `;

        const resultUser = await query(queryUsers, [
            baseInfo.user_fullname,
            baseInfo.user_phone,
            baseInfo.user_email,
            userId,
        ]);
        await query(queryEmployees, [
            baseInfo.user_department,
            baseInfo.user_data_join_company,
            baseInfo.user_position,
            userId,
        ]);

        return resultUser.rows[0] || null;
    }

    async deleteListEmployee(userIds: string[]): Promise<number> {
        const queryDeleteEmployees = `
            DELETE FROM employees
            WHERE employee_id = ANY($1::uuid[]);
        `;
        const queryDeleteUsers = `
            DELETE FROM users
            WHERE user_id = ANY($1::uuid[]);
        `;

        await query(queryDeleteEmployees, [userIds]);
        const result = await query(queryDeleteUsers, [userIds]);
        return result.rowCount;
    }

    async updateUserName(
        userId: string,
        fullName: string
    ): Promise<User | null> {
        if (!fullName || fullName.trim().length === 0) {
            throw new Error('Full name cannot be empty');
        }

        const now = new Date().toISOString();
        const result = await query(
            `UPDATE users SET full_name = $1, updated_at = $2 WHERE user_id = $3 RETURNING *`,
            [fullName.trim(), now, userId]
        );

        return result.rows[0] || null;
    }

    async updateUserPhone(userId: string, phone: string): Promise<User | null> {
        if (!phone || phone.trim().length === 0) {
            throw new Error('Phone cannot be empty');
        }

        const now = new Date().toISOString();
        const result = await query(
            `UPDATE users SET phone = $1, updated_at = $2 WHERE user_id = $3 RETURNING *`,
            [phone.trim(), now, userId]
        );

        return result.rows[0] || null;
    }

    async updateUserDepartment(
        userId: string,
        department: string
    ): Promise<User | null> {
        if (!department || department.trim().length === 0) {
            throw new Error('Department cannot be empty');
        }

        const now = new Date().toISOString();

        // Update employees table
        await query(
            `UPDATE employees SET department = $1, updated_at = $2 WHERE employee_id = $3`,
            [department.trim(), now, userId]
        );

        // Return updated user
        const result = await query(`SELECT * FROM users WHERE user_id = $1`, [
            userId,
        ]);

        return result.rows[0] || null;
    }

    async updateUserPosition(
        userId: string,
        position: string
    ): Promise<User | null> {
        if (!position || position.trim().length === 0) {
            throw new Error('Position cannot be empty');
        }

        const now = new Date().toISOString();

        // Update employees table
        await query(
            `UPDATE employees SET position = $1, updated_at = $2 WHERE employee_id = $3`,
            [position.trim(), now, userId]
        );

        // Return updated user
        const result = await query(`SELECT * FROM users WHERE user_id = $1`, [
            userId,
        ]);

        return result.rows[0] || null;
    }

    async resetPassword(
        userId: string
    ): Promise<{ newPassword: string; user: User } | null> {
        const {
            generateRandomPassword,
            generateSalt,
            hashPassword,
        } = require('../utils/crypto');

        const user = await this.getUserById(userId);
        if (!user) {
            return null;
        }

        const newPassword = generateRandomPassword();
        const salt = generateSalt();
        const passwordHash = hashPassword(newPassword, salt);
        const now = new Date().toISOString();

        const result = await query(
            `UPDATE users SET salt = $1, password_hash = $2, updated_at = $3 WHERE user_id = $4 RETURNING *`,
            [salt, passwordHash, now, userId]
        );

        return {
            newPassword,
            user: result.rows[0],
        };
    }

    async importUsersFromFile(file: Express.Multer.File): Promise<number> {
        const fileContent = fs.readFileSync(file.path, 'utf-8');

        const records: any[] = await new Promise((resolve, reject) => {
            parse(
                fileContent,
                { columns: true, skip_empty_lines: true },
                (err, records) => {
                    if (err) return reject(err);
                    resolve(records);
                }
            );
        });

        const users: CreateUserRequest[] = records.map((record) => ({
            email: record.email,
            phone: record.phone,
            password: record.password,
            full_name: record.full_name,
            avatar_url: record.avatar_url || undefined,
            company_id: record.company_id || undefined,
            employee_code: record.employee_code || undefined,
            department: record.department || undefined,
            position: record.position || undefined,
            hire_date: record.hire_date || undefined,
            salary: record.salary ? parseFloat(record.salary) : undefined,
            role: parseInt(record.role, 10),
        }));

        // insert
        for (const userData of users) {
            try {
                await this.createUser(userData);
            } catch (err) {
                console.error(`Failed to import ${userData.email}`, err);
            }
        }

        fs.unlinkSync(file.path);
        return users.length;
    }
}

export const userService = new UserService();
