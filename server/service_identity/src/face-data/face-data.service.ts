import { Injectable, NotFoundException, BadRequestException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { FaceData } from '../entities/face-data.entity';
import { Employee } from '../entities/employee.entity';
import { RedisService } from '../redis/redis.service';
import { v4 as uuidv4 } from 'uuid';
import * as fs from 'fs';
import * as path from 'path';

@Injectable()
export class FaceDataService {
  constructor(
    @InjectRepository(FaceData)
    private faceDataRepository: Repository<FaceData>,
    @InjectRepository(Employee)
    private employeeRepository: Repository<Employee>,
    private redisService: RedisService,
  ) {}

  async uploadFaceData(
    userId: string,
    files: Express.Multer.File[],
  ): Promise<FaceData[]> {
    // Check if user exists and get employee info
    const employee = await this.employeeRepository.findOne({
      where: { user_id: userId },
      relations: ['user'],
    });

    if (!employee) {
      throw new NotFoundException('Employee not found');
    }

    // Check current face data count
    const currentCount = await this.faceDataRepository.count({
      where: { employee_id: employee.employee_id },
    });

    const maxImages = parseInt(process.env.MAX_FACE_IMAGES_PER_USER) || 5;
    if (currentCount + files.length > maxImages) {
      throw new BadRequestException(
        `Maximum ${maxImages} face images allowed per user. Current: ${currentCount}`,
      );
    }

    const uploadedFaceData: FaceData[] = [];

    for (const file of files) {
      try {
        // Simulate face processing (in real implementation, use face recognition library)
        const faceEmbedding = await this.processFaceImage(file.path);

        // Create face data record
        const faceData = this.faceDataRepository.create({
          employee_id: employee.employee_id,
          face_embedding: faceEmbedding,
          image_path: file.path,
          image_name: file.filename,
          image_size: file.size,
          image_type: file.mimetype,
          metadata: {
            original_name: file.originalname,
            uploaded_at: new Date(),
          },
        });

        const savedFaceData = await this.faceDataRepository.save(faceData);
        uploadedFaceData.push(savedFaceData);
      } catch (error) {
        // Clean up file if processing failed
        if (fs.existsSync(file.path)) {
          fs.unlinkSync(file.path);
        }
        throw new BadRequestException(`Failed to process image: ${error.message}`);
      }
    }

    // Update user face_registered status
    if (uploadedFaceData.length > 0) {
      await this.employeeRepository.update(
        { user_id: userId },
        { user: { face_registered: true } },
      );
    }

    // Clear cache
    await this.redisService.del(`user_face_data:${userId}`);
    await this.redisService.del(`user:${userId}`);

    return uploadedFaceData;
  }

  async getFaceData(userId: string): Promise<FaceData[]> {
    const cacheKey = `user_face_data:${userId}`;
    const cached = await this.redisService.get(cacheKey);
    
    if (cached) {
      return JSON.parse(cached);
    }

    const employee = await this.employeeRepository.findOne({
      where: { user_id: userId },
    });

    if (!employee) {
      throw new NotFoundException('Employee not found');
    }

    const faceData = await this.faceDataRepository.find({
      where: { employee_id: employee.employee_id },
      order: { created_at: 'DESC' },
    });

    // Cache for 30 minutes
    await this.redisService.set(cacheKey, JSON.stringify(faceData), 1800);

    return faceData;
  }

  async deleteFaceData(userId: string, faceId: string): Promise<void> {
    const employee = await this.employeeRepository.findOne({
      where: { user_id: userId },
    });

    if (!employee) {
      throw new NotFoundException('Employee not found');
    }

    const faceData = await this.faceDataRepository.findOne({
      where: { face_id: faceId, employee_id: employee.employee_id },
    });

    if (!faceData) {
      throw new NotFoundException('Face data not found');
    }

    // Delete file from filesystem
    if (faceData.image_path && fs.existsSync(faceData.image_path)) {
      fs.unlinkSync(faceData.image_path);
    }

    // Delete from database
    await this.faceDataRepository.remove(faceData);

    // Check if user still has face data
    const remainingCount = await this.faceDataRepository.count({
      where: { employee_id: employee.employee_id },
    });

    // Update user face_registered status if no face data left
    if (remainingCount === 0) {
      await this.employeeRepository.update(
        { user_id: userId },
        { user: { face_registered: false } },
      );
    }

    // Clear cache
    await this.redisService.del(`user_face_data:${userId}`);
    await this.redisService.del(`user:${userId}`);
  }

  private async processFaceImage(imagePath: string): Promise<Buffer> {
    // In a real implementation, this would:
    // 1. Load the image
    // 2. Detect faces using OpenCV or similar
    // 3. Extract face embeddings using a face recognition model
    // 4. Return the embedding vector as Buffer
    
    // For now, we'll simulate this with a random buffer
    // In production, use libraries like:
    // - @vladmandic/face-api
    // - opencv4nodejs
    // - face-recognition
    
    const mockEmbedding = Buffer.alloc(128); // Typical face embedding size
    for (let i = 0; i < 128; i++) {
      mockEmbedding[i] = Math.floor(Math.random() * 256);
    }
    
    return mockEmbedding;
  }
}
