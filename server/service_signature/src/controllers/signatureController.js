const Signature = require('../models/signature');
const fs = require('fs');
const path = require('path');

// POST - Upload chữ ký
exports.uploadSignature = async (req, res) => {
  try {
    if (!req.file) {
      return res.status(400).json({ message: 'No file uploaded' });
    }

    const { employee_id, company_id } = req.body;

    if (!employee_id || !company_id) {
      return res.status(400).json({ message: 'employee_id and company_id are required' });
    }

    const filePath = `/uploads/signatures/${req.file.filename}`;
    const signature = await Signature.create(employee_id, company_id, filePath);

    res.status(201).json({
      message: 'Signature uploaded successfully',
      data: signature,
    });
  } catch (error) {
    console.error('Upload error:', error);
    res.status(500).json({ message: 'Error uploading signature', error: error.message });
  }
};

// GET - Lấy danh sách chữ ký của nhân viên
exports.getSignatures = async (req, res) => {
  try {
    const { user_id } = req.params;

    if (!user_id) {
      return res.status(400).json({ message: 'user_id is required' });
    }

    const signatures = await Signature.getByEmployeeId(user_id);

    res.status(200).json({
      message: 'Signatures retrieved successfully',
      data: signatures,
    });
  } catch (error) {
    console.error('Get signatures error:', error);
    res.status(500).json({ message: 'Error retrieving signatures', error: error.message });
  }
};

// DELETE - Xóa chữ ký
exports.deleteSignature = async (req, res) => {
  try {
    const { signature_id } = req.params;

    if (!signature_id) {
      return res.status(400).json({ message: 'signature_id is required' });
    }

    const signature = await Signature.getById(signature_id);
    if (!signature) {
      return res.status(404).json({ message: 'Signature not found' });
    }

    // Xóa file
    const filePath = path.join(process.cwd(), signature.file_path);
    if (fs.existsSync(filePath)) {
      fs.unlinkSync(filePath);
    }

    // Xóa record từ database
    await Signature.delete(signature_id);

    res.status(200).json({
      message: 'Signature deleted successfully',
    });
  } catch (error) {
    console.error('Delete error:', error);
    res.status(500).json({ message: 'Error deleting signature', error: error.message });
  }
};
