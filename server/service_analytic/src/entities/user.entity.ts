import { Entity, PrimaryGeneratedColumn, Column, OneToOne, JoinColumn } from 'typeorm';
import { Employee } from './employee.entity';

export enum UserRole {
  SYSTEM_ADMIN = 0,
  COMPANY_ADMIN = 1,
  EMPLOYEE = 2,
}

@Entity('users')
export class User {
  @PrimaryGeneratedColumn('uuid')
  user_id: string;

  @Column({ type: 'varchar', length: 255, unique: true })
  email: string;

  @Column({ type: 'varchar', length: 255 })
  salt: string;

  @Column({ type: 'varchar', length: 255 })
  password_hash: string;

  @Column({ type: 'varchar', length: 255 })
  full_name: string;

  @Column({ type: 'int', enum: UserRole })
  role: UserRole;

  @OneToOne(() => Employee, employee => employee.user)
  employee: Employee;
}
