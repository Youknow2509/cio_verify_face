import { ApiProperty } from '@nestjs/swagger';

export class DailyReportResponseDto {
  @ApiProperty({ description: 'Ngày báo cáo' })
  date: string;

  @ApiProperty({ description: 'Tổng số nhân viên' })
  total_employees: number;

  @ApiProperty({ description: 'Số nhân viên có mặt' })
  present_employees: number;

  @ApiProperty({ description: 'Số nhân viên đi muộn' })
  late_employees: number;

  @ApiProperty({ description: 'Số nhân viên về sớm' })
  early_leave_employees: number;

  @ApiProperty({ description: 'Số nhân viên vắng mặt' })
  absent_employees: number;

  @ApiProperty({ description: 'Tỷ lệ có mặt (%)' })
  attendance_rate: number;

  @ApiProperty({ description: 'Chi tiết theo phòng ban' })
  departments: DepartmentReportDto[];

  @ApiProperty({ description: 'Chi tiết theo ca làm việc' })
  shifts: ShiftReportDto[];
}

export class DepartmentReportDto {
  @ApiProperty({ description: 'Tên phòng ban' })
  department_name: string;

  @ApiProperty({ description: 'Tổng số nhân viên' })
  total_employees: number;

  @ApiProperty({ description: 'Số nhân viên có mặt' })
  present_employees: number;

  @ApiProperty({ description: 'Tỷ lệ có mặt (%)' })
  attendance_rate: number;
}

export class ShiftReportDto {
  @ApiProperty({ description: 'Tên ca làm việc' })
  shift_name: string;

  @ApiProperty({ description: 'Giờ bắt đầu' })
  start_time: string;

  @ApiProperty({ description: 'Giờ kết thúc' })
  end_time: string;

  @ApiProperty({ description: 'Tổng số nhân viên' })
  total_employees: number;

  @ApiProperty({ description: 'Số nhân viên có mặt' })
  present_employees: number;

  @ApiProperty({ description: 'Tỷ lệ có mặt (%)' })
  attendance_rate: number;
}

export class SummaryReportResponseDto {
  @ApiProperty({ description: 'Tháng báo cáo' })
  month: string;

  @ApiProperty({ description: 'Tổng số ngày làm việc' })
  total_working_days: number;

  @ApiProperty({ description: 'Tổng số nhân viên' })
  total_employees: number;

  @ApiProperty({ description: 'Tỷ lệ có mặt trung bình (%)' })
  average_attendance_rate: number;

  @ApiProperty({ description: 'Tổng giờ làm việc' })
  total_working_hours: number;

  @ApiProperty({ description: 'Tổng giờ làm thêm' })
  total_overtime_hours: number;

  @ApiProperty({ description: 'Chi tiết theo tuần' })
  weekly_summary: WeeklySummaryDto[];

  @ApiProperty({ description: 'Top nhân viên có mặt nhiều nhất' })
  top_attendance_employees: EmployeeAttendanceDto[];

  @ApiProperty({ description: 'Nhân viên có tỷ lệ vắng mặt cao' })
  low_attendance_employees: EmployeeAttendanceDto[];
}

export class WeeklySummaryDto {
  @ApiProperty({ description: 'Tuần' })
  week: number;

  @ApiProperty({ description: 'Ngày bắt đầu tuần' })
  start_date: string;

  @ApiProperty({ description: 'Ngày kết thúc tuần' })
  end_date: string;

  @ApiProperty({ description: 'Tỷ lệ có mặt (%)' })
  attendance_rate: number;

  @ApiProperty({ description: 'Tổng giờ làm việc' })
  total_hours: number;
}

export class EmployeeAttendanceDto {
  @ApiProperty({ description: 'Mã nhân viên' })
  employee_code: string;

  @ApiProperty({ description: 'Tên nhân viên' })
  full_name: string;

  @ApiProperty({ description: 'Số ngày có mặt' })
  present_days: number;

  @ApiProperty({ description: 'Tỷ lệ có mặt (%)' })
  attendance_rate: number;

  @ApiProperty({ description: 'Tổng giờ làm việc' })
  total_hours: number;
}

export class ExportReportResponseDto {
  @ApiProperty({ description: 'ID job xuất báo cáo' })
  job_id: string;

  @ApiProperty({ description: 'Trạng thái job' })
  status: string;

  @ApiProperty({ description: 'Thông báo' })
  message: string;

  @ApiProperty({ description: 'Link download (nếu có)', required: false })
  download_url?: string;
}
