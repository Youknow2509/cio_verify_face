import { Request, Response } from 'express';
import { companyService } from '../services/companyService';
import { sendSuccess, sendError } from '../utils/response';

export class CompanyController {
  async getAllCompanies(req: Request, res: Response) {
    try {
      const companies = await companyService.getAllCompanies();
      return sendSuccess(res, companies, 'Companies retrieved successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to retrieve companies');
    }
  }

  async getCompanyById(req: Request, res: Response) {
    try {
      const { company_id } = req.params;
      const company = await companyService.getCompanyById(company_id);

      if (!company) {
        return sendError(res, 'Company not found', 404, 'Not Found');
      }

      return sendSuccess(res, company, 'Company retrieved successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to retrieve company');
    }
  }

  async createCompany(req: Request, res: Response) {
    try {
      const { name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices } = req.body;

      if (!name) {
        return sendError(res, 'Company name is required', 400, 'Validation Error');
      }

      const company = await companyService.createCompany({
        name,
        address,
        phone,
        email,
        website,
        status,
        subscription_plan,
        subscription_start_date,
        subscription_end_date,
        max_employees,
        max_devices,
      });

      return sendSuccess(res, company, 'Company created successfully', 201);
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to create company');
    }
  }

  async updateCompany(req: Request, res: Response) {
    try {
      const { company_id } = req.params;
      const { name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices } = req.body;

      const company = await companyService.getCompanyById(company_id);
      if (!company) {
        return sendError(res, 'Company not found', 404, 'Not Found');
      }

      const updatedCompany = await companyService.updateCompany(company_id, {
        name,
        address,
        phone,
        email,
        website,
        status,
        subscription_plan,
        subscription_start_date,
        subscription_end_date,
        max_employees,
        max_devices,
      });

      return sendSuccess(res, updatedCompany, 'Company updated successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to update company');
    }
  }

  async deleteCompany(req: Request, res: Response) {
    try {
      const { company_id } = req.params;

      const company = await companyService.getCompanyById(company_id);
      if (!company) {
        return sendError(res, 'Company not found', 404, 'Not Found');
      }

      await companyService.deleteCompany(company_id);
      return sendSuccess(res, { company_id }, 'Company deleted successfully');
    } catch (error: any) {
      return sendError(res, error.message, 500, 'Failed to delete company');
    }
  }
}

export const companyController = new CompanyController();
