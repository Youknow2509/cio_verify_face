import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, OneToMany } from 'typeorm';
import { Company } from './company.entity';
import { DailyAttendanceSummary } from './daily-attendance-summary.entity';

@Entity('work_shifts')
export class WorkShift {
  @PrimaryGeneratedColumn('uuid')
  shift_id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'time' })
  start_time: string;

  @Column({ type: 'time' })
  end_time: string;

  @ManyToOne(() => Company, company => company.workShifts)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @OneToMany(() => DailyAttendanceSummary, summary => summary.workShift)
  dailyAttendanceSummaries: DailyAttendanceSummary[];
}
