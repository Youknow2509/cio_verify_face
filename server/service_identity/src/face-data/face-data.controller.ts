import {
  Controller,
  Get,
  Post,
  Delete,
  Param,
  UseInterceptors,
  UploadedFiles,
  UseGuards,
  Request,
  BadRequestException,
} from '@nestjs/common';
import { FilesInterceptor } from '@nestjs/platform-express';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiConsumes } from '@nestjs/swagger';
import { FaceDataService } from './face-data.service';
import { FaceDataResponseDto } from '../dto/face-data.dto';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';

@ApiTags('Face Data')
@Controller('api/v1/users/:userId/face-data')
@UseGuards(JwtAuthGuard)
@ApiBearerAuth()
export class FaceDataController {
  constructor(private readonly faceDataService: FaceDataService) {}

  @Post()
  @UseInterceptors(FilesInterceptor('images', 5)) // Max 5 files
  @ApiOperation({ summary: 'Upload face data images' })
  @ApiConsumes('multipart/form-data')
  @ApiResponse({ status: 201, description: 'Face data uploaded successfully' })
  @ApiResponse({ status: 400, description: 'Invalid file or processing failed' })
  @ApiResponse({ status: 404, description: 'User not found' })
  async uploadFaceData(
    @Param('userId') userId: string,
    @UploadedFiles() files: Express.Multer.File[],
  ): Promise<FaceDataResponseDto[]> {
    if (!files || files.length === 0) {
      throw new BadRequestException('No files uploaded');
    }

    const faceData = await this.faceDataService.uploadFaceData(userId, files);
    
    return faceData.map(data => ({
      face_id: data.face_id,
      employee_id: data.employee_id,
      image_name: data.image_name,
      image_path: data.image_path,
      created_at: data.created_at,
    }));
  }

  @Get()
  @ApiOperation({ summary: 'Get user face data' })
  @ApiResponse({ status: 200, description: 'Face data retrieved successfully' })
  @ApiResponse({ status: 404, description: 'User not found' })
  async getFaceData(@Param('userId') userId: string) {
    return this.faceDataService.getFaceData(userId);
  }

  @Delete(':faceId')
  @ApiOperation({ summary: 'Delete specific face data' })
  @ApiResponse({ status: 200, description: 'Face data deleted successfully' })
  @ApiResponse({ status: 404, description: 'Face data not found' })
  async deleteFaceData(
    @Param('userId') userId: string,
    @Param('faceId') faceId: string,
  ) {
    await this.faceDataService.deleteFaceData(userId, faceId);
    return { message: 'Face data deleted successfully' };
  }
}
