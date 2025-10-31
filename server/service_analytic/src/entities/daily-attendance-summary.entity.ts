import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn } from 'typeorm';
import { Employee } from './employee.entity';
import { WorkShift } from './work-shift.entity';

export enum AttendanceStatus {
  PRESENT = 0,
  LATE = 1,
  EARLY_LEAVE = 2,
  ABSENT = 3,
}

@Entity('daily_attendance_summary')
export class DailyAttendanceSummary {
  @PrimaryGeneratedColumn('uuid')
  summary_id: string;

  @Column({ type: 'uuid' })
  employee_id: string;

  @Column({ type: 'uuid' })
  shift_id: string;

  @Column({ type: 'date' })
  work_date: Date;

  @Column({ type: 'timestamp', nullable: true })
  check_in_time: Date;

  @Column({ type: 'timestamp', nullable: true })
  check_out_time: Date;

  @Column({ type: 'int', default: 0 })
  total_hours: number;

  @Column({ type: 'int', default: 0 })
  overtime_hours: number;

  @Column({ type: 'int', enum: AttendanceStatus })
  status: AttendanceStatus;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  created_at: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  updated_at: Date;

  @ManyToOne(() => Employee, employee => employee.dailyAttendanceSummaries)
  @JoinColumn({ name: 'employee_id' })
  employee: Employee;

  @ManyToOne(() => WorkShift, workShift => workShift.dailyAttendanceSummaries)
  @JoinColumn({ name: 'shift_id' })
  workShift: WorkShift;
}
