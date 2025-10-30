import { Controller, Get } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { HealthService } from './health.service';

@ApiTags('Health')
@Controller('health')
export class HealthController {
  constructor(private readonly healthService: HealthService) {}

  @Get()
  @ApiOperation({ summary: 'Health check endpoint' })
  @ApiResponse({ status: 200, description: 'Service is healthy' })
  @ApiResponse({ status: 503, description: 'Service is unhealthy' })
  check() {
    return this.healthService.checkAll();
  }

  @Get('database')
  @ApiOperation({ summary: 'Database health check' })
  @ApiResponse({ status: 200, description: 'Database is healthy' })
  checkDatabase() {
    return this.healthService.checkDatabase();
  }

  @Get('redis')
  @ApiOperation({ summary: 'Redis health check' })
  @ApiResponse({ status: 200, description: 'Redis is healthy' })
  checkRedis() {
    return this.healthService.checkRedis();
  }

  @Get('memory')
  @ApiOperation({ summary: 'Memory health check' })
  @ApiResponse({ status: 200, description: 'Memory usage is healthy' })
  checkMemory() {
    return this.healthService.checkMemory();
  }

  @Get('disk')
  @ApiOperation({ summary: 'Disk health check' })
  @ApiResponse({ status: 200, description: 'Disk usage is healthy' })
  checkDisk() {
    return this.healthService.checkDisk();
  }
}
