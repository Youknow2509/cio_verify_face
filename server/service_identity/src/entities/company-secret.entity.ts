import { Entity, PrimaryGeneratedColumn, Column, OneToOne, JoinColumn } from 'typeorm';
import { Company } from './company.entity';

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

  @OneToOne(() => Company, company => company.companySecret)
  @JoinColumn({ name: 'company_id' })
  company: Company;
}
