const express = require('express');
const router = express.Router();
const signatureController = require('../controllers/signatureController');
const upload = require('../middleware/upload');

/**
 * @swagger
 * /api/v1/signatures:
 *   post:
 *     summary: Upload digital signature
 *     description: Upload a signature image for an employee
 *     tags:
 *       - Signatures
 *     requestBody:
 *       required: true
 *       content:
 *         multipart/form-data:
 *           schema:
 *             type: object
 *             properties:
 *               file:
 *                 type: string
 *                 format: binary
 *                 description: Signature image file (PNG, JPEG, GIF)
 *               employee_id:
 *                 type: integer
 *                 description: Employee ID
 *               company_id:
 *                 type: integer
 *                 description: Company ID
 *             required:
 *               - file
 *               - employee_id
 *               - company_id
 *     responses:
 *       201:
 *         description: Signature uploaded successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                 data:
 *                   type: object
 *                   properties:
 *                     id:
 *                       type: integer
 *                     employee_id:
 *                       type: integer
 *                     company_id:
 *                       type: integer
 *                     file_path:
 *                       type: string
 *                     created_at:
 *                       type: string
 *       400:
 *         description: Missing required fields or invalid file
 *       500:
 *         description: Server error
 */
router.post('/', upload.single('file'), signatureController.uploadSignature);

/**
 * @swagger
 * /api/v1/signatures/{user_id}:
 *   get:
 *     summary: Get employee signatures
 *     description: Retrieve all signatures for a specific employee
 *     tags:
 *       - Signatures
 *     parameters:
 *       - in: path
 *         name: user_id
 *         schema:
 *           type: integer
 *         required: true
 *         description: Employee ID
 *     responses:
 *       200:
 *         description: List of signatures retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                 data:
 *                   type: array
 *                   items:
 *                     type: object
 *                     properties:
 *                       id:
 *                         type: integer
 *                       employee_id:
 *                         type: integer
 *                       company_id:
 *                         type: integer
 *                       file_path:
 *                         type: string
 *                       created_at:
 *                         type: string
 *       400:
 *         description: Missing user_id parameter
 *       500:
 *         description: Server error
 */
router.get('/:user_id', signatureController.getSignatures);

/**
 * @swagger
 * /api/v1/signatures/{signature_id}:
 *   delete:
 *     summary: Delete a signature
 *     description: Delete a signature and its associated file
 *     tags:
 *       - Signatures
 *     parameters:
 *       - in: path
 *         name: signature_id
 *         schema:
 *           type: integer
 *         required: true
 *         description: Signature ID
 *     responses:
 *       200:
 *         description: Signature deleted successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *       404:
 *         description: Signature not found
 *       500:
 *         description: Server error
 */
router.delete('/:signature_id', signatureController.deleteSignature);

module.exports = router;
