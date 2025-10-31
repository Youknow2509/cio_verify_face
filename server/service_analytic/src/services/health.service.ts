import { Injectable } from '@nestjs/common';
import { InjectDataSource } from '@nestjs/typeorm';
import { DataSource } from 'typeorm';

@Injectable()
export class HealthService {
  constructor(
    @InjectDataSource()
    private dataSource: DataSource,
  ) {}

  async checkHealth(): Promise<{ status: string; database: string; timestamp: string }> {
    let databaseStatus = 'disconnected';
    
    try {
      await this.dataSource.query('SELECT 1');
      databaseStatus = 'connected';
    } catch (error) {
      databaseStatus = 'error';
    }

    return {
      status: 'ok',
      database: databaseStatus,
      timestamp: new Date().toISOString(),
    };
  }
}
