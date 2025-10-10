import { IsString, IsOptional, IsArray, IsUUID } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class UploadFaceDataDto {
  @ApiProperty({ description: 'Face images', type: 'array', items: { type: 'string', format: 'binary' } })
  @IsArray()
  images: Express.Multer.File[];
}

export class FaceDataResponseDto {
  @ApiProperty({ description: 'Face data ID' })
  @IsUUID()
  face_id: string;

  @ApiProperty({ description: 'Employee ID' })
  @IsUUID()
  employee_id: string;

  @ApiProperty({ description: 'Image name' })
  @IsString()
  @IsOptional()
  image_name?: string;

  @ApiProperty({ description: 'Image path' })
  @IsString()
  @IsOptional()
  image_path?: string;

  @ApiProperty({ description: 'Created at' })
  created_at: Date;
}
