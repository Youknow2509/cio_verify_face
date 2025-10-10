import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, UpdateDateColumn, OneToMany, OneToOne } from 'typeorm';
import { Employee } from './employee.entity';
import { CompanySecret } from './company-secret.entity';

@Entity('companies')
export class Company {
  @PrimaryGeneratedColumn('uuid')
  company_id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'text', nullable: true })
  address: string;

  @Column({ type: 'varchar', length: 100, nullable: true })
  email: string;

  @Column({ type: 'varchar', length: 20, nullable: true })
  phone: string;

  @Column({ type: 'varchar', length: 50, default: 'ACTIVE' })
  status: string;

  @Column({ type: 'varchar', length: 100, nullable: true })
  plan: string; // FREE, BASIC, PREMIUM, ENTERPRISE

  @Column({ type: 'int', default: 0 })
  max_employees: number;

  @Column({ type: 'json', nullable: true })
  settings: any; // Company-specific settings

  @CreateDateColumn()
  created_at: Date;

  @UpdateDateColumn()
  updated_at: Date;

  @OneToMany(() => Employee, employee => employee.company)
  employees: Employee[];

  @OneToOne(() => CompanySecret, companySecret => companySecret.company)
  companySecret: CompanySecret;
}
