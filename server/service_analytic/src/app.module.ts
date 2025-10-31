import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { JwtModule } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { CacheModule } from '@nestjs/cache-manager';

// Entities
import { Company } from './entities/company.entity';
import { Employee } from './entities/employee.entity';
import { User } from './entities/user.entity';
import { Device } from './entities/device.entity';
import { CompanySecret } from './entities/company-secret.entity';
import { FaceData } from './entities/face-data.entity';
import { WorkShift } from './entities/work-shift.entity';
import { AttendanceRecord } from './entities/attendance-record.entity';
import { DailyAttendanceSummary } from './entities/daily-attendance-summary.entity';

// Services
import { ReportsService } from './services/reports.service';
import { HealthService } from './services/health.service';

// Controllers
import { ReportsController } from './controllers/reports.controller';
import { HealthController } from './controllers/health.controller';

// Guards
import { JwtAuthGuard } from './guards/jwt-auth.guard';
import { RolesGuard } from './guards/roles.guard';

// Config
import databaseConfig, { jwtConfig } from './config/database.config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [databaseConfig, jwtConfig],
    }),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: (configService: ConfigService) => ({
        type: 'postgres',
        host: configService.get('database.host'),
        port: configService.get('database.port'),
        username: configService.get('database.username'),
        password: configService.get('database.password'),
        database: configService.get('database.database'),
        entities: [
          Company,
          Employee,
          User,
          Device,
          CompanySecret,
          FaceData,
          WorkShift,
          AttendanceRecord,
          DailyAttendanceSummary,
        ],
        synchronize: configService.get('database.synchronize'),
        logging: configService.get('database.logging'),
      }),
      inject: [ConfigService],
    }),
    TypeOrmModule.forFeature([
      Company,
      Employee,
      User,
      Device,
      CompanySecret,
      FaceData,
      WorkShift,
      AttendanceRecord,
      DailyAttendanceSummary,
    ]),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      useFactory: (configService: ConfigService) => ({
        secret: configService.get('jwt.secret'),
        signOptions: { expiresIn: configService.get('jwt.expiresIn') },
      }),
      inject: [ConfigService],
    }),
    PassportModule,
    CacheModule.register({
      ttl: 300, // 5 minutes
      max: 100, // maximum number of items in cache
    }),
  ],
  controllers: [ReportsController, HealthController],
  providers: [
    ReportsService,
    HealthService,
    JwtAuthGuard,
    RolesGuard,
  ],
})
export class AppModule {}