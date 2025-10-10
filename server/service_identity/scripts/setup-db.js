const { Client } = require('pg');
const fs = require('fs');
const path = require('path');

async function setupDatabase() {
  const client = new Client({
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 5432,
    user: process.env.DB_USERNAME || 'postgres',
    password: process.env.DB_PASSWORD || 'password',
    database: process.env.DB_DATABASE || 'identity_service',
  });

  try {
    await client.connect();
    console.log('Connected to PostgreSQL database');

    // Read and execute migration file
    const migrationPath = path.join(__dirname, '../src/migrations/001-initial-schema.sql');
    const migrationSQL = fs.readFileSync(migrationPath, 'utf8');
    
    await client.query(migrationSQL);
    console.log('Database schema created successfully');

    // Insert sample data
    await insertSampleData(client);
    console.log('Sample data inserted successfully');

  } catch (error) {
    console.error('Error setting up database:', error);
    process.exit(1);
  } finally {
    await client.end();
  }
}

async function insertSampleData(client) {
  // Insert sample company
  const companyResult = await client.query(`
    INSERT INTO companies (name, address, email, phone) 
    VALUES ($1, $2, $3, $4) 
    RETURNING company_id
  `, [
    'ABC Corporation',
    '123 Main St, City, Country',
    'contact@abccorp.com',
    '+1234567890'
  ]);

  const companyId = companyResult.rows[0].company_id;

  // Insert company secret
  const bcrypt = require('bcrypt');
  const salt = await bcrypt.genSalt(10);
  const secret = 'company-secret-key';
  const secretHash = await bcrypt.hash(secret, salt);

  await client.query(`
    INSERT INTO companies_secret (company_id, salt, secret_hash) 
    VALUES ($1, $2, $3)
  `, [companyId, salt, secretHash]);

  // Insert sample admin user
  const userSalt = await bcrypt.genSalt(10);
  const userPassword = await bcrypt.hash('admin123', userSalt);

  const userResult = await client.query(`
    INSERT INTO users (email, salt, password_hash, full_name, role) 
    VALUES ($1, $2, $3, $4, $5) 
    RETURNING user_id
  `, [
    'admin@abccorp.com',
    userSalt,
    userPassword,
    'System Administrator',
    0 // SYSTEM_ADMIN
  ]);

  const userId = userResult.rows[0].user_id;

  // Insert sample employee
  await client.query(`
    INSERT INTO employees (company_id, user_id, employee_code, department, position, phone, hire_date) 
    VALUES ($1, $2, $3, $4, $5, $6, $7)
  `, [
    companyId,
    userId,
    'ADM001',
    'IT Department',
    'System Administrator',
    '+1234567890',
    new Date()
  ]);

  // Insert sample work shift
  await client.query(`
    INSERT INTO work_shifts (company_id, name, start_time, end_time, description) 
    VALUES ($1, $2, $3, $4, $5)
  `, [
    companyId,
    'Morning Shift',
    '08:00:00',
    '17:00:00',
    'Regular morning shift'
  ]);

  console.log('Sample data inserted:');
  console.log('- Company: ABC Corporation');
  console.log('- Admin User: admin@abccorp.com (password: admin123)');
  console.log('- Work Shift: Morning Shift (08:00-17:00)');
}

// Run setup if called directly
if (require.main === module) {
  setupDatabase();
}

module.exports = { setupDatabase };
