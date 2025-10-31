import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, OneToOne } from 'typeorm';
import { Company } from './company.entity';
import { User } from './user.entity';
import { FaceData } from './face-data.entity';
import { AttendanceRecord } from './attendance-record.entity';
import { DailyAttendanceSummary } from './daily-attendance-summary.entity';

@Entity('employees')
export class Employee {
  @PrimaryGeneratedColumn('uuid')
  employee_id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @Column({ type: 'varchar', length: 50, unique: true })
  employee_code: string;

  @ManyToOne(() => Company, company => company.employees)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @OneToOne(() => User, user => user.employee)
  user: User;

  @OneToOne(() => FaceData, faceData => faceData.employee)
  faceData: FaceData;

  @OneToMany(() => AttendanceRecord, attendanceRecord => attendanceRecord.employee)
  attendanceRecords: AttendanceRecord[];

  @OneToMany(() => DailyAttendanceSummary, summary => summary.employee)
  dailyAttendanceSummaries: DailyAttendanceSummary[];
}
