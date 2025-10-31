import { query } from '../config/database';
import { Company, CreateCompanyRequest, UpdateCompanyRequest } from '../types';
import { v4 as uuidv4 } from 'uuid';

export class CompanyService {
  async getAllCompanies(): Promise<Company[]> {
    const result = await query('SELECT * FROM companies ORDER BY created_at DESC');
    return result.rows;
  }

  async getCompanyById(companyId: string): Promise<Company | null> {
    const result = await query('SELECT * FROM companies WHERE company_id = $1', [companyId]);
    return result.rows[0] || null;
  }

  async createCompany(data: CreateCompanyRequest): Promise<Company> {
    const companyId = uuidv4();
    const now = new Date().toISOString();

    const result = await query(
      `INSERT INTO companies 
       (company_id, name, address, phone, email, website, status, subscription_plan, 
        subscription_start_date, subscription_end_date, max_employees, max_devices, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
       RETURNING *`,
      [
        companyId,
        data.name,
        data.address || null,
        data.phone || null,
        data.email || null,
        data.website || null,
        data.status ?? 1,
        data.subscription_plan ?? 0,
        data.subscription_start_date || null,
        data.subscription_end_date || null,
        data.max_employees ?? 100,
        data.max_devices ?? 10,
        now,
        now,
      ]
    );

    return result.rows[0];
  }

  async updateCompany(companyId: string, data: UpdateCompanyRequest): Promise<Company | null> {
    const now = new Date().toISOString();
    const updates: string[] = [];
    const values: any[] = [];
    let paramIndex = 1;

    if (data.name !== undefined) {
      updates.push(`name = $${paramIndex++}`);
      values.push(data.name);
    }
    if (data.address !== undefined) {
      updates.push(`address = $${paramIndex++}`);
      values.push(data.address);
    }
    if (data.phone !== undefined) {
      updates.push(`phone = $${paramIndex++}`);
      values.push(data.phone);
    }
    if (data.email !== undefined) {
      updates.push(`email = $${paramIndex++}`);
      values.push(data.email);
    }
    if (data.website !== undefined) {
      updates.push(`website = $${paramIndex++}`);
      values.push(data.website);
    }
    if (data.status !== undefined) {
      updates.push(`status = $${paramIndex++}`);
      values.push(data.status);
    }
    if (data.subscription_plan !== undefined) {
      updates.push(`subscription_plan = $${paramIndex++}`);
      values.push(data.subscription_plan);
    }
    if (data.subscription_start_date !== undefined) {
      updates.push(`subscription_start_date = $${paramIndex++}`);
      values.push(data.subscription_start_date);
    }
    if (data.subscription_end_date !== undefined) {
      updates.push(`subscription_end_date = $${paramIndex++}`);
      values.push(data.subscription_end_date);
    }
    if (data.max_employees !== undefined) {
      updates.push(`max_employees = $${paramIndex++}`);
      values.push(data.max_employees);
    }
    if (data.max_devices !== undefined) {
      updates.push(`max_devices = $${paramIndex++}`);
      values.push(data.max_devices);
    }

    if (updates.length === 0) {
      return this.getCompanyById(companyId);
    }

    updates.push(`updated_at = $${paramIndex++}`);
    values.push(now);
    values.push(companyId);

    const result = await query(
      `UPDATE companies SET ${updates.join(', ')} WHERE company_id = $${paramIndex} RETURNING *`,
      values
    );

    return result.rows[0] || null;
  }

  async deleteCompany(companyId: string): Promise<boolean> {
    const result = await query('DELETE FROM companies WHERE company_id = $1', [companyId]);
    return result.rowCount > 0;
  }
}

export const companyService = new CompanyService();
