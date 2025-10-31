import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn } from 'typeorm';
import { Employee } from './employee.entity';
import { Device } from './device.entity';

export enum RecordType {
  CHECK_IN = 0,
  CHECK_OUT = 1,
}

@Entity('attendance_records')
export class AttendanceRecord {
  @PrimaryGeneratedColumn('uuid')
  record_id: string;

  @Column({ type: 'uuid' })
  employee_id: string;

  @Column({ type: 'uuid' })
  device_id: string;

  @Column({ type: 'timestamp' })
  timestamp: Date;

  @Column({ type: 'int', enum: RecordType })
  record_type: RecordType;

  @Column({ type: 'text', nullable: true })
  metadate: string;

  @ManyToOne(() => Employee, employee => employee.attendanceRecords)
  @JoinColumn({ name: 'employee_id' })
  employee: Employee;

  @ManyToOne(() => Device, device => device.attendanceRecords)
  @JoinColumn({ name: 'device_id' })
  device: Device;
}
