import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, OneToMany } from 'typeorm';
import { Company } from './company.entity';
import { CompanySecret } from './company-secret.entity';
import { AttendanceRecord } from './attendance-record.entity';

@Entity('devices')
export class Device {
  @PrimaryGeneratedColumn('uuid')
  device_id: string;

  @Column({ type: 'uuid', nullable: true })
  location_id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @Column({ type: 'varchar', length: 255 })
  address: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'varchar', length: 50 })
  status: string;

  @Column({ type: 'varchar', length: 255 })
  token: string;

  @Column({ type: 'uuid' })
  company_secret_id: string;

  @ManyToOne(() => Company, company => company.devices)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @ManyToOne(() => CompanySecret, companySecret => companySecret.devices)
  @JoinColumn({ name: 'company_secret_id' })
  companySecret: CompanySecret;

  @OneToMany(() => AttendanceRecord, attendanceRecord => attendanceRecord.device)
  attendanceRecords: AttendanceRecord[];
}
