import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { CompaniesModule } from './companies/companies.module';
import { UsersModule } from './users/users.module';
import { FaceDataModule } from './face-data/face-data.module';
import { AuthModule } from './auth/auth.module';
import { DatabaseModule } from './database/database.module';
import { RedisModule } from './redis/redis.module';
import { HealthModule } from './health/health.module';
import { ServiceDiscoveryModule } from './service-discovery/service-discovery.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    DatabaseModule,
    RedisModule,
    HealthModule,
    ServiceDiscoveryModule,
    AuthModule,
    CompaniesModule,
    UsersModule,
    FaceDataModule,
  ],
})
export class AppModule {}
