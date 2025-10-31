const pool = require('../config/database');

class Signature {
  // Tạo chữ ký mới
  static async create(employeeId, companyId, filePath) {
    const connection = await pool.getConnection();
    try {
      const [result] = await connection.execute(
        'INSERT INTO signatures (employee_id, company_id, file_path) VALUES (?, ?, ?)',
        [employeeId, companyId, filePath]
      );
      
      const [rows] = await connection.execute(
        'SELECT * FROM signatures WHERE id = ?',
        [result.insertId]
      );
      return rows[0];
    } finally {
      connection.release();
    }
  }

  // Lấy danh sách chữ ký của nhân viên
  static async getByEmployeeId(employeeId) {
    const connection = await pool.getConnection();
    try {
      const [rows] = await connection.execute(
        'SELECT * FROM signatures WHERE employee_id = ? ORDER BY created_at DESC',
        [employeeId]
      );
      return rows;
    } finally {
      connection.release();
    }
  }

  // Lấy chữ ký theo ID
  static async getById(signatureId) {
    const connection = await pool.getConnection();
    try {
      const [rows] = await connection.execute(
        'SELECT * FROM signatures WHERE id = ?',
        [signatureId]
      );
      return rows[0];
    } finally {
      connection.release();
    }
  }

  // Xóa chữ ký
  static async delete(signatureId) {
    const connection = await pool.getConnection();
    try {
      const [result] = await connection.execute(
        'DELETE FROM signatures WHERE id = ?',
        [signatureId]
      );
      return result;
    } finally {
      connection.release();
    }
  }
}

module.exports = Signature;
