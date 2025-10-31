import { 
  Controller, 
  Get, 
  Post, 
  Query, 
  Body, 
  UseGuards, 
  HttpStatus,
  HttpCode 
} from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
import { ReportsService } from '../services/reports.service';
import { 
  DailyReportQueryDto, 
  SummaryReportQueryDto, 
  ExportReportDto 
} from '../dto/reports.dto';
import { 
  DailyReportResponseDto, 
  SummaryReportResponseDto, 
  ExportReportResponseDto 
} from '../dto/reports-response.dto';
import { JwtAuthGuard } from '../guards/jwt-auth.guard';
import { RolesGuard } from '../guards/roles.guard';
import { Roles } from '../decorators/roles.decorator';
import { UserRole } from '../entities/user.entity';

@ApiTags('Reports')
@Controller('api/v1/reports')
@UseGuards(JwtAuthGuard, RolesGuard)
@ApiBearerAuth()
export class ReportsController {
  constructor(private readonly reportsService: ReportsService) {}

  @Get('daily')
  @ApiOperation({ 
    summary: 'Báo cáo chi tiết ngày',
    description: 'Lấy báo cáo chi tiết về chấm công trong một ngày cụ thể theo công ty/địa điểm'
  })
  @ApiResponse({ 
    status: HttpStatus.OK, 
    description: 'Báo cáo ngày được trả về thành công',
    type: DailyReportResponseDto 
  })
  @ApiResponse({ 
    status: HttpStatus.BAD_REQUEST, 
    description: 'Dữ liệu đầu vào không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.UNAUTHORIZED, 
    description: 'Token không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.FORBIDDEN, 
    description: 'Không có quyền truy cập báo cáo' 
  })
  @Roles(UserRole.SYSTEM_ADMIN, UserRole.COMPANY_ADMIN)
  async getDailyReport(@Query() query: DailyReportQueryDto): Promise<DailyReportResponseDto> {
    return await this.reportsService.getDailyReport(query);
  }

  @Get('summary')
  @ApiOperation({ 
    summary: 'Báo cáo tổng hợp tháng',
    description: 'Lấy báo cáo tổng hợp về chấm công trong một tháng'
  })
  @ApiResponse({ 
    status: HttpStatus.OK, 
    description: 'Báo cáo tổng hợp tháng được trả về thành công',
    type: SummaryReportResponseDto 
  })
  @ApiResponse({ 
    status: HttpStatus.BAD_REQUEST, 
    description: 'Dữ liệu đầu vào không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.UNAUTHORIZED, 
    description: 'Token không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.FORBIDDEN, 
    description: 'Không có quyền truy cập báo cáo' 
  })
  @Roles(UserRole.SYSTEM_ADMIN, UserRole.COMPANY_ADMIN)
  async getSummaryReport(@Query() query: SummaryReportQueryDto): Promise<SummaryReportResponseDto> {
    return await this.reportsService.getSummaryReport(query);
  }

  @Post('export')
  @HttpCode(HttpStatus.ACCEPTED)
  @ApiOperation({ 
    summary: 'Xuất báo cáo',
    description: 'Xuất báo cáo chấm công ra file Excel/PDF trong khoảng thời gian chỉ định'
  })
  @ApiResponse({ 
    status: HttpStatus.ACCEPTED, 
    description: 'Yêu cầu xuất báo cáo đã được chấp nhận',
    type: ExportReportResponseDto 
  })
  @ApiResponse({ 
    status: HttpStatus.BAD_REQUEST, 
    description: 'Dữ liệu đầu vào không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.UNAUTHORIZED, 
    description: 'Token không hợp lệ' 
  })
  @ApiResponse({ 
    status: HttpStatus.FORBIDDEN, 
    description: 'Không có quyền xuất báo cáo' 
  })
  @ApiResponse({ 
    status: HttpStatus.TOO_MANY_REQUESTS, 
    description: 'Quá nhiều yêu cầu xuất báo cáo' 
  })
  @Roles(UserRole.SYSTEM_ADMIN, UserRole.COMPANY_ADMIN)
  async exportReport(@Body() exportDto: ExportReportDto): Promise<ExportReportResponseDto> {
    return await this.reportsService.exportReport(exportDto);
  }
}
