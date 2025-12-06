import * as crypto from 'crypto';

export function hashPassword(password: string, salt: string): string {
    return crypto.createHmac('sha256', salt).update(password).digest('hex');
}

export function generateSalt(): string {
    return crypto.randomBytes(16).toString('hex');
}

export function generateRandomPassword(length: number = 12): string {
    const chars =
        'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*';
    let password = '';
    for (let i = 0; i < length; i++) {
        password += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return password;
}

export function verifyPassword(
    password: string,
    salt: string,
    hash: string
): boolean {
    const hashAttempt = hashPassword(password, salt);
    return hashAttempt === hash;
}
