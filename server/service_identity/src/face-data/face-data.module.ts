import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { MulterModule } from '@nestjs/platform-express';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { FaceDataService } from './face-data.service';
import { FaceDataController } from './face-data.controller';
import { FaceData } from '../entities/face-data.entity';
import { Employee } from '../entities/employee.entity';
import { RedisModule } from '../redis/redis.module';
import { diskStorage } from 'multer';
import { extname } from 'path';

@Module({
  imports: [
    TypeOrmModule.forFeature([FaceData, Employee]),
    RedisModule,
    MulterModule.registerAsync({
      imports: [ConfigModule],
      useFactory: (configService: ConfigService) => ({
        storage: diskStorage({
          destination: configService.get('UPLOAD_PATH') || './uploads',
          filename: (req, file, callback) => {
            const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9);
            const ext = extname(file.originalname);
            callback(null, `face-${uniqueSuffix}${ext}`);
          },
        }),
        fileFilter: (req, file, callback) => {
          if (file.mimetype.match(/\/(jpg|jpeg|png|gif)$/)) {
            callback(null, true);
          } else {
            callback(new Error('Only image files are allowed!'), false);
          }
        },
        limits: {
          fileSize: parseInt(configService.get('MAX_FILE_SIZE')) || 5242880, // 5MB
        },
      }),
      inject: [ConfigService],
    }),
  ],
  controllers: [FaceDataController],
  providers: [FaceDataService],
  exports: [FaceDataService],
})
export class FaceDataModule {}
