import { IsString, IsEmail, IsOptional, IsNotEmpty, MaxLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateCompanyDto {
  @ApiProperty({ description: 'Company name', example: 'ABC Corporation' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(255)
  name: string;

  @ApiProperty({ description: 'Company address', example: '123 Main St, City, Country' })
  @IsString()
  @IsOptional()
  address?: string;

  @ApiProperty({ description: 'Company email', example: 'contact@abccorp.com' })
  @IsEmail()
  @IsOptional()
  email?: string;

  @ApiProperty({ description: 'Company phone', example: '+1234567890' })
  @IsString()
  @IsOptional()
  @MaxLength(20)
  phone?: string;
}
