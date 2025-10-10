import { Entity, PrimaryGeneratedColumn, Column, OneToOne, JoinColumn } from 'typeorm';
import { Employee } from './employee.entity';

@Entity('face_data')
export class FaceData {
  @PrimaryGeneratedColumn('uuid')
  face_id: string;

  @Column({ type: 'uuid' })
  employee_id: string;

  @Column({ type: 'bytea' })
  face_embedding: Buffer;

  @OneToOne(() => Employee, employee => employee.faceData)
  @JoinColumn({ name: 'employee_id' })
  employee: Employee;
}
