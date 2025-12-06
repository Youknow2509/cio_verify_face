import { Router } from 'express';
import { userController } from '../controllers/userController';
import { faceDataController } from '../controllers/faceDataController';
import multer from 'multer';

const upload = multer({ dest: 'uploads/' });
/**
 * @swagger
 * /api/v1/users:
 *   get:
 *     summary: Get all users
 *     description: Retrieve list of users, optionally filtered by company_id
 *     tags:
 *       - Users
 *     parameters:
 *       - in: query
 *         name: company_id
 *         schema:
 *           type: string
 *           format: uuid
 *         description: Filter users by company ID
 *     responses:
 *       200:
 *         description: Users retrieved successfully
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/users:
 *   post:
 *     summary: Create new user/employee
 *     description: Create a new user in the system. If role is 2 (EMPLOYEE), also creates employee record
 *     tags:
 *       - Users
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - email
 *               - phone
 *               - password
 *               - full_name
 *               - role
 *             properties:
 *               email:
 *                 type: string
 *                 format: email
 *               phone:
 *                 type: string
 *               password:
 *                 type: string
 *                 format: password
 *               full_name:
 *                 type: string
 *               avatar_url:
 *                 type: string
 *                 format: uri
 *               role:
 *                 type: integer
 *                 enum: [0, 1, 2]
 *               company_id:
 *                 type: string
 *                 format: uuid
 *               employee_code:
 *                 type: string
 *               department:
 *                 type: string
 *               position:
 *                 type: string
 *               hire_date:
 *                 type: string
 *                 format: date
 *               salary:
 *                 type: number
 *     responses:
 *       201:
 *         description: User created successfully
 *       400:
 *         description: Validation error or user already exists
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/users/{user_id}:
 *   get:
 *     summary: Get user by ID
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     responses:
 *       200:
 *         description: User retrieved successfully
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}:
 *   put:
 *     summary: Update user
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
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
 *               email:
 *                 type: string
 *               phone:
 *                 type: string
 *               full_name:
 *                 type: string
 *               avatar_url:
 *                 type: string
 *               status:
 *                 type: integer
 *               department:
 *                 type: string
 *               position:
 *                 type: string
 *               salary:
 *                 type: number
 *     responses:
 *       200:
 *         description: User updated successfully
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}:
 *   delete:
 *     summary: Delete user
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     responses:
 *       200:
 *         description: User deleted successfully
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/face-data:
 *   post:
 *     summary: Upload face data
 *     description: Enroll a new face profile for a user via AI service and persist to face_profiles
 *     tags:
 *       - Face Data
 *     parameters:
 *       - in: path
 *         name: user_id
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
 *             required:
 *               - image_url
 *               - company_id
 *             properties:
 *               image_url:
 *                 type: string
 *                 format: uri
 *               company_id:
 *                 type: string
 *                 format: uuid
 *               make_primary:
 *                 type: boolean
 *               metadata:
 *                 type: object
 *                 additionalProperties:
 *                   type: string
 *     responses:
 *       201:
 *         description: Face data created successfully
 *       400:
 *         description: Validation error
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/face-data/upload:
 *   post:
 *     summary: Upload face data (binary)
 *     description: Enroll face profile by uploading image file directly (multipart/form-data)
 *     tags:
 *       - Face Data
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     requestBody:
 *       required: true
 *       content:
 *         multipart/form-data:
 *           schema:
 *             type: object
 *             required:
 *               - image
 *               - company_id
 *             properties:
 *               image:
 *                 type: string
 *                 format: binary
 *                 description: Image file (JPEG/PNG)
 *               company_id:
 *                 type: string
 *                 format: uuid
 *               make_primary:
 *                 type: boolean
 *     responses:
 *       201:
 *         description: Face profile enrolled successfully
 *       400:
 *         description: Validation error
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/face-data:
 *   get:
 *     summary: Get face data list
 *     description: Retrieve face profiles for a user (requires company_id)
 *     tags:
 *       - Face Data
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *       - in: query
 *         name: company_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *     responses:
 *       200:
 *         description: Face data list retrieved
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/face-data/{fid}:
 *   delete:
 *     summary: Delete face data
 *     description: Soft delete a specific face profile (use hard=true to hard delete)
 *     tags:
 *       - Face Data
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *       - in: path
 *         name: fid
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *       - in: query
 *         name: company_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *       - in: query
 *         name: hard
 *         required: false
 *         schema:
 *           type: boolean
 *     responses:
 *       200:
 *         description: Face data deleted successfully
 *       404:
 *         description: Face data not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/face-data/{fid}/primary:
 *   put:
 *     summary: Update primary face profile
 *     description: Set or unset a face profile as primary for the user
 *     tags:
 *       - Face Data
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *       - in: path
 *         name: fid
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
 *             required:
 *               - company_id
 *               - status
 *             properties:
 *               company_id:
 *                 type: string
 *                 format: uuid
 *               status:
 *                 type: boolean
 *     responses:
 *       200:
 *         description: Primary face profile updated successfully
 *       404:
 *         description: Face profile not found
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/name:
 *   put:
 *     summary: Update user name
 *     description: Update the full name of a user by their ID
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: User ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - full_name
 *             properties:
 *               full_name:
 *                 type: string
 *                 description: New full name for the user
 *                 example: "John Doe"
 *     responses:
 *       200:
 *         description: User name updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 user_id:
 *                   type: string
 *                   format: uuid
 *                 full_name:
 *                   type: string
 *                 email:
 *                   type: string
 *                 updated_at:
 *                   type: string
 *                   format: date-time
 *       400:
 *         description: Validation error (missing full_name)
 *       404:
 *         description: User not found
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/phone:
 *   put:
 *     summary: Update user phone
 *     description: Update the phone number of a user by their ID
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: User ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - phone
 *             properties:
 *               phone:
 *                 type: string
 *                 description: New phone number for the user
 *                 example: "+84912345678"
 *     responses:
 *       200:
 *         description: User phone updated successfully
 *       400:
 *         description: Validation error (missing phone)
 *       404:
 *         description: User not found
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/department:
 *   put:
 *     summary: Update user department
 *     description: Update the department of a user by their ID
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: User ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - department
 *             properties:
 *               department:
 *                 type: string
 *                 description: New department for the user
 *                 example: "Sales"
 *     responses:
 *       200:
 *         description: User department updated successfully
 *       400:
 *         description: Validation error (missing department)
 *       404:
 *         description: User not found
 *       500:
 *         description: Internal server error
 */

/**
 * @swagger
 * /api/v1/users/{user_id}/position:
 *   put:
 *     summary: Update user position
 *     description: Update the position of a user by their ID
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: User ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - position
 *             properties:
 *               position:
 *                 type: string
 *                 description: New position for the user
 *                 example: "Sales Manager"
 *     responses:
 *       200:
 *         description: User position updated successfully
 *       400:
 *         description: Validation error (missing position)
 *       404:
 *         description: User not found
 *       500:
 *         description: Internal server error
 */

/**
 * @route /api/v1/users/update-base
 * @summary Update user's base information
 * @method POST
 * @tags Users
 * @param {string} user_id.body.required - User ID
 * @param {string} user_fullname.body - Full name
 * @param {string} user_phone.body - Phone number
 * @param {string} user_email.body - Email address
 * @param {string} user_department.body - Department
 * @param {string} user_data_join_company.body - Date joined company
 * @param {string} user_position.body - Position
 * @return {object} 200 - Success response - application/json
 * @return {object} 400 - Validation error - application/json
 */

const router = Router();

// Users endpoints
router.get('/', (req, res) => userController.getAllUsers(req, res));
router.post('/', (req, res) => userController.createUser(req, res));
router.get('/:user_id', (req, res) => userController.getUserById(req, res));
router.put('/:user_id', (req, res) => userController.updateUser(req, res));
router.put('/:user_id/name', (req, res) =>
    userController.updateUserName(req, res)
);
router.put('/:user_id/phone', (req, res) =>
    userController.updateUserPhone(req, res)
);
router.put('/:user_id/department', (req, res) =>
    userController.updateUserDepartment(req, res)
);
router.put('/:user_id/position', (req, res) =>
    userController.updateUserPosition(req, res)
);
router.delete('/:user_id', (req, res) => userController.deleteUser(req, res));
router.post('/update-base', (req, res) =>
    userController.updateBaseInfo(req, res)
);
router.delete('/delete-multiple', (req, res) =>
    userController.deleteListEmployee(req, res)
);
router.post('/import-from-file', upload.single('file'), (req, res) =>
    userController.importUsersFromFile(req, res)
);

// Face data endpoints
router.post('/:user_id/face-data', (req, res) =>
    faceDataController.createFaceData(req, res)
);
router.post('/:user_id/face-data/upload', upload.single('image'), (req, res) =>
    faceDataController.createFaceDataBinary(req, res)
);
router.get('/:user_id/face-data', (req, res) =>
    faceDataController.getFaceDataByUserId(req, res)
);
router.delete('/:user_id/face-data/:fid', (req, res) =>
    faceDataController.deleteFaceData(req, res)
);

/**
 * @swagger
 * /api/v1/users/{user_id}/reset-password:
 *   post:
 *     summary: Reset user password
 *     description: Generate and reset a new random password for a user. The new password will be sent via Kafka to the user notification service (not returned in response for security).
 *     tags:
 *       - Users
 *     parameters:
 *       - in: path
 *         name: user_id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: User ID
 *     responses:
 *       200:
 *         description: Password reset successfully and sent to user via Kafka
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 data:
 *                   type: object
 *                   properties:
 *                     user_id:
 *                       type: string
 *                       format: uuid
 *                     email:
 *                       type: string
 *                     full_name:
 *                       type: string
 *                     reset_at:
 *                       type: string
 *                       format: date-time
 *                 message:
 *                   type: string
 *                   example: "Password reset successfully. New password has been sent to user via Kafka notification."
 *       404:
 *         description: User not found
 *       500:
 *         description: Internal server error
 */

router.post('/:user_id/reset-password', (req, res) =>
    userController.resetPassword(req, res)
);

export default router;
