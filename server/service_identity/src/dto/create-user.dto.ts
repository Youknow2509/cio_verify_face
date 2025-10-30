import { IsString, IsEmail, IsNotEmpty, IsOptional, MaxLength, IsEnum } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';
import { UserRole } from '../entities/user.entity';

export class CreateUserDto {
  @ApiProperty({ description: 'User email', example: 'john.doe@company.com' })
  @IsEmail()
  @IsNotEmpty()
  email: string;

  @ApiProperty({ description: 'User password', example: 'securePassword123' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(255)
  password: string;

  @ApiProperty({ description: 'Full name', example: 'John Doe' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(255)
  full_name: string;

  @ApiProperty({ description: 'User role', enum: UserRole, example: UserRole.EMPLOYEE })
  @IsEnum(UserRole)
  @IsOptional()
  role?: UserRole;

  @ApiProperty({ description: 'Employee code', example: 'EMP001' })
  @IsString()
  @IsOptional()
  @MaxLength(50)
  employee_code?: string;

  @ApiProperty({ description: 'Department', example: 'IT Department' })
  @IsString()
  @IsOptional()
  @MaxLength(100)
  department?: string;

  @ApiProperty({ description: 'Position', example: 'Software Developer' })
  @IsString()
  @IsOptional()
  @MaxLength(100)
  position?: string;

  @ApiProperty({ description: 'Phone number', example: '+1234567890' })
  @IsString()
  @IsOptional()
  @MaxLength(20)
  phone?: string;
}
