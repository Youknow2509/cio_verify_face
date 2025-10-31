import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, OneToOne, OneToMany } from 'typeorm';
import { Company } from './company.entity';
import { User } from './user.entity';
import { FaceData } from './face-data.entity';

@Entity('employees')
export class Employee {
  @PrimaryGeneratedColumn('uuid')
  employee_id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @Column({ type: 'uuid', unique: true })
  user_id: string;

  @Column({ type: 'varchar', length: 50 })
  employee_code: string;

  @Column({ type: 'varchar', length: 100, nullable: true })
  department: string;

  @Column({ type: 'varchar', length: 100, nullable: true })
  position: string;

  @Column({ type: 'varchar', length: 20, nullable: true })
  phone: string;

  @Column({ type: 'date', nullable: true })
  hire_date: Date;

  @Column({ type: 'boolean', default: true })
  is_active: boolean;

  @Column({ type: 'json', nullable: true })
  permissions: any; // Employee-specific permissions

  @Column({ type: 'varchar', length: 100, nullable: true })
  manager_id: string; // Reference to another employee

  @ManyToOne(() => Company, company => company.employees)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @OneToOne(() => User, user => user.employee)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @OneToMany(() => FaceData, faceData => faceData.employee)
  faceData: FaceData[];
}
