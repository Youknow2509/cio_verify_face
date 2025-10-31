import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, OneToMany } from 'typeorm';
import { Company } from './company.entity';
import { Device } from './device.entity';

@Entity('companies_secret')
export class CompanySecret {
  @PrimaryGeneratedColumn('uuid')
  company_secret_id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @Column({ type: 'varchar', length: 255 })
  salt: string;

  @Column({ type: 'varchar', length: 255 })
  secret_hash: string;

  @ManyToOne(() => Company, company => company.companySecrets)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @OneToMany(() => Device, device => device.companySecret)
  devices: Device[];
}
