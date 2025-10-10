import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, JoinColumn, CreateDateColumn } from 'typeorm';
import { Employee } from './employee.entity';

@Entity('face_data')
export class FaceData {
  @PrimaryGeneratedColumn('uuid')
  face_id: string;

  @Column({ type: 'uuid' })
  employee_id: string;

  @Column({ type: 'bytea' })
  face_embedding: Buffer;

  @Column({ type: 'varchar', length: 500, nullable: true })
  image_path: string;

  @Column({ type: 'varchar', length: 100, nullable: true })
  image_name: string;

  @Column({ type: 'int', nullable: true })
  image_size: number;

  @Column({ type: 'varchar', length: 50, nullable: true })
  image_type: string;

  @Column({ type: 'json', nullable: true })
  metadata: any;

  @CreateDateColumn()
  created_at: Date;

  @ManyToOne(() => Employee, employee => employee.faceData)
  @JoinColumn({ name: 'employee_id' })
  employee: Employee;
}
