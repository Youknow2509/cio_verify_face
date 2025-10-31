import { Entity, PrimaryGeneratedColumn, Column, OneToMany } from 'typeorm';
import { Employee } from './employee.entity';
import { WorkShift } from './work-shift.entity';
import { Device } from './device.entity';
import { CompanySecret } from './company-secret.entity';

@Entity('companies')
export class Company {
  @PrimaryGeneratedColumn('uuid')
  company_id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'text', nullable: true })
  address: string;

  @OneToMany(() => Employee, employee => employee.company)
  employees: Employee[];

  @OneToMany(() => WorkShift, workShift => workShift.company)
  workShifts: WorkShift[];

  @OneToMany(() => Device, device => device.company)
  devices: Device[];

  @OneToMany(() => CompanySecret, companySecret => companySecret.company)
  companySecrets: CompanySecret[];
}
