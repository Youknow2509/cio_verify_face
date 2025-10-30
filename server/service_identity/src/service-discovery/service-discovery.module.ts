import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ServiceDiscoveryService } from './service-discovery.service';

@Module({
  imports: [ConfigModule],
  providers: [
    {
      provide: ServiceDiscoveryService,
      useFactory: (configService: ConfigService) => {
        return new ServiceDiscoveryService({
          serviceName: 'identity-service',
          port: configService.get('PORT') || 3000,
          host: configService.get('SERVICE_HOST') || 'localhost',
          version: '1.0.0',
          environment: configService.get('NODE_ENV') || 'development',
        });
      },
      inject: [ConfigService],
    },
  ],
  exports: [ServiceDiscoveryService],
})
export class ServiceDiscoveryModule {}
