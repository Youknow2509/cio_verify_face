const pool = require('../config/database');

const migrate = async () => {
  let connection;
  try {
    console.log('Starting migration...');
    connection = await pool.getConnection();

    await connection.execute(`
      CREATE TABLE IF NOT EXISTS signatures (
        id INT AUTO_INCREMENT PRIMARY KEY,
        employee_id INT NOT NULL,
        company_id INT NOT NULL,
        file_path VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
      );
    `);

    console.log('Migration completed successfully!');
    process.exit(0);
  } catch (error) {
    console.error('Migration error:', error);
    process.exit(1);
  } finally {
    if (connection) connection.release();
  }
};

migrate();
