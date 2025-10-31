import { Router } from 'express';
import { companyController } from '../controllers/companyController';

/**
 * @swagger
 * /api/v1/companies:
 *   get:
 *     summary: Get all companies
 *     description: Retrieve a list of all companies
 *     tags:
 *       - Companies
 *     responses:
 *       200:
 *         description: List of companies retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *                 data:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Company'
 *                 statusCode:
 *                   type: integer
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/companies:
 *   post:
 *     summary: Create a new company
 *     description: Create a new company in the system
 *     tags:
 *       - Companies
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - name
 *             properties:
 *               name:
 *                 type: string
 *               address:
 *                 type: string
 *               phone:
 *                 type: string
 *               email:
 *                 type: string
 *               website:
 *                 type: string
 *               status:
 *                 type: integer
 *                 default: 1
 *               subscription_plan:
 *                 type: integer
 *                 default: 0
 *               max_employees:
 *                 type: integer
 *                 default: 100
 *               max_devices:
 *                 type: integer
 *                 default: 10
 *     responses:
 *       201:
 *         description: Company created successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *                 data:
 *                   $ref: '#/components/schemas/Company'
 *                 statusCode:
 *                   type: integer
 *       400:
 *         description: Validation error
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/companies/{company_id}:
 *   get:
 *     summary: Get company by ID
 *     description: Retrieve company details by company_id
 *     tags:
 *       - Companies
 *     parameters:
 *       - in: path
 *         name: company_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: Company UUID
 *     responses:
 *       200:
 *         description: Company retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *                 data:
 *                   $ref: '#/components/schemas/Company'
 *       404:
 *         description: Company not found
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/companies/{company_id}:
 *   put:
 *     summary: Update company
 *     description: Update company information
 *     tags:
 *       - Companies
 *     parameters:
 *       - in: path
 *         name: company_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               name:
 *                 type: string
 *               address:
 *                 type: string
 *               phone:
 *                 type: string
 *               email:
 *                 type: string
 *               website:
 *                 type: string
 *               status:
 *                 type: integer
 *               subscription_plan:
 *                 type: integer
 *               max_employees:
 *                 type: integer
 *               max_devices:
 *                 type: integer
 *     responses:
 *       200:
 *         description: Company updated successfully
 *       404:
 *         description: Company not found
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/companies/{company_id}:
 *   delete:
 *     summary: Delete company
 *     description: Delete a company from the system
 *     tags:
 *       - Companies
 *     parameters:
 *       - in: path
 *         name: company_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     responses:
 *       200:
 *         description: Company deleted successfully
 *       404:
 *         description: Company not found
 *       500:
 *         description: Internal server error
 */

const router = Router();

// Companies endpoints
router.get('/', (req, res) => companyController.getAllCompanies(req, res));
router.post('/', (req, res) => companyController.createCompany(req, res));
router.get('/:company_id', (req, res) => companyController.getCompanyById(req, res));
router.put('/:company_id', (req, res) => companyController.updateCompany(req, res));
router.delete('/:company_id', (req, res) => companyController.deleteCompany(req, res));

export default router;
