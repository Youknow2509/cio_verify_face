import { Injectable } from '@nestjs/common';
import { HealthCheckService, TypeOrmHealthIndicator, MemoryHealthIndicator, DiskHealthIndicator } from '@nestjs/terminus';
import { RedisService } from '../redis/redis.service';

@Injectable()
export class HealthService {
  constructor(
    private health: HealthCheckService,
    private db: TypeOrmHealthIndicator,
    private memory: MemoryHealthIndicator,
    private disk: DiskHealthIndicator,
    private redisService: RedisService,
  ) {}

  async checkDatabase() {
    return this.health.check([
      () => this.db.pingCheck('database'),
    ]);
  }

  async checkRedis() {
    return this.health.check([
      async () => {
        const isHealthy = await this.redisService.exists('health-check');
        if (isHealthy) {
          return { redis: { status: 'up' } };
        }
        throw new Error('Redis is not available');
      },
    ]);
  }

  async checkMemory() {
    return this.health.check([
      () => this.memory.checkHeap('memory_heap', 150 * 1024 * 1024),
      () => this.memory.checkRSS('memory_rss', 150 * 1024 * 1024),
    ]);
  }

  async checkDisk() {
    return this.health.check([
      () => this.disk.checkStorage('storage', { path: '/', thresholdPercent: 0.5 }),
    ]);
  }

  async checkAll() {
    return this.health.check([
      () => this.db.pingCheck('database'),
      async () => {
        const isHealthy = await this.redisService.exists('health-check');
        if (isHealthy) {
          return { redis: { status: 'up' } };
        }
        throw new Error('Redis is not available');
      },
      () => this.memory.checkHeap('memory_heap', 150 * 1024 * 1024),
      () => this.memory.checkRSS('memory_rss', 150 * 1024 * 1024),
      () => this.disk.checkStorage('storage', { path: '/', thresholdPercent: 0.5 }),
    ]);
  }
}
