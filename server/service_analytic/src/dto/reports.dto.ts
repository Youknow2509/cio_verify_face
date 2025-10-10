import { IsOptional, IsString, IsDateString, IsUUID, IsEnum } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export enum ReportFormat {
  EXCEL = 'excel',
  PDF = 'pdf',
  CSV = 'csv',
}

export class DailyReportQueryDto {
  @ApiProperty({ description: 'Ngày báo cáo (YYYY-MM-DD)', example: '2024-01-15' })
  @IsDateString()
  date: string;

  @ApiProperty({ description: 'ID công ty', required: false })
  @IsOptional()
  @IsUUID()
  company_id?: string;

  @ApiProperty({ description: 'ID địa điểm', required: false })
  @IsOptional()
  @IsUUID()
  location_id?: string;
}

export class SummaryReportQueryDto {
  @ApiProperty({ description: 'Tháng báo cáo (YYYY-MM)', example: '2024-01' })
  @IsString()
  month: string;

  @ApiProperty({ description: 'ID công ty', required: false })
  @IsOptional()
  @IsUUID()
  company_id?: string;
}

export class ExportReportDto {
  @ApiProperty({ description: 'Ngày bắt đầu (YYYY-MM-DD)', example: '2024-01-01' })
  @IsDateString()
  start_date: string;

  @ApiProperty({ description: 'Ngày kết thúc (YYYY-MM-DD)', example: '2024-01-31' })
  @IsDateString()
  end_date: string;

  @ApiProperty({ description: 'Định dạng xuất file', enum: ReportFormat })
  @IsEnum(ReportFormat)
  format: ReportFormat;

  @ApiProperty({ description: 'ID công ty', required: false })
  @IsOptional()
  @IsUUID()
  company_id?: string;

  @ApiProperty({ description: 'Email nhận file', required: false })
  @IsOptional()
  @IsString()
  email?: string;
}
