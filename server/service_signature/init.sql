-- Init script for MySQL
USE mysql;
CREATE USER IF NOT EXISTS 'signature_user'@'%' IDENTIFIED BY 'signature_password';
GRANT ALL PRIVILEGES ON signature_db.* TO 'signature_user'@'%';
FLUSH PRIVILEGES;
