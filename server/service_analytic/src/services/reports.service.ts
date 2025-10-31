import { Injectable, NotFoundException, BadRequestException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Between } from 'typeorm';
import { DailyAttendanceSummary } from '../entities/daily-attendance-summary.entity';
import { Employee } from '../entities/employee.entity';
import { WorkShift } from '../entities/work-shift.entity';
import { Company } from '../entities/company.entity';
import { DailyReportQueryDto, SummaryReportQueryDto, ExportReportDto } from '../dto/reports.dto';
import { 
  DailyReportResponseDto, 
  SummaryReportResponseDto, 
  ExportReportResponseDto 
} from '../dto/reports-response.dto';
import { AttendanceStatus } from '../entities/daily-attendance-summary.entity';

@Injectable()
export class ReportsService {
  constructor(
    @InjectRepository(DailyAttendanceSummary)
    private readonly attendanceSummaryRepository: Repository<DailyAttendanceSummary>,
    @InjectRepository(Employee)
    private readonly employeeRepository: Repository<Employee>,
    @InjectRepository(WorkShift)
    private readonly workShiftRepository: Repository<WorkShift>,
    @InjectRepository(Company)
    private readonly companyRepository: Repository<Company>,
  ) {}

  async getDailyReport(query: DailyReportQueryDto): Promise<DailyReportResponseDto> {
    const { date, company_id, location_id } = query;
    const reportDate = new Date(date);

    // Validate date
    if (isNaN(reportDate.getTime())) {
      throw new BadRequestException('Invalid date format');
    }

    // Build query conditions
    const whereConditions: any = {
      work_date: reportDate,
    };

    if (company_id) {
      whereConditions.employee = { company_id };
    }

    // Get attendance summaries for the date
    const attendanceSummaries = await this.attendanceSummaryRepository.find({
      where: whereConditions,
      relations: ['employee', 'workShift'],
    });

    // Calculate statistics
    const totalEmployees = attendanceSummaries.length;
    const presentEmployees = attendanceSummaries.filter(
      summary => summary.status === AttendanceStatus.PRESENT
    ).length;
    const lateEmployees = attendanceSummaries.filter(
      summary => summary.status === AttendanceStatus.LATE
    ).length;
    const earlyLeaveEmployees = attendanceSummaries.filter(
      summary => summary.status === AttendanceStatus.EARLY_LEAVE
    ).length;
    const absentEmployees = attendanceSummaries.filter(
      summary => summary.status === AttendanceStatus.ABSENT
    ).length;

    const attendanceRate = totalEmployees > 0 ? (presentEmployees / totalEmployees) * 100 : 0;

    // Group by departments (assuming department info is in employee or company)
    const departments = await this.getDepartmentReport(attendanceSummaries);
    
    // Group by shifts
    const shifts = await this.getShiftReport(attendanceSummaries);

    return {
      date,
      total_employees: totalEmployees,
      present_employees: presentEmployees,
      late_employees: lateEmployees,
      early_leave_employees: earlyLeaveEmployees,
      absent_employees: absentEmployees,
      attendance_rate: Math.round(attendanceRate * 100) / 100,
      departments,
      shifts,
    };
  }

  async getSummaryReport(query: SummaryReportQueryDto): Promise<SummaryReportResponseDto> {
    const { month, company_id } = query;
    
    // Parse month (YYYY-MM)
    const [year, monthNum] = month.split('-').map(Number);
    if (!year || !monthNum || monthNum < 1 || monthNum > 12) {
      throw new BadRequestException('Invalid month format. Use YYYY-MM');
    }

    const startDate = new Date(year, monthNum - 1, 1);
    const endDate = new Date(year, monthNum, 0); // Last day of month

    // Build query conditions
    const whereConditions: any = {
      work_date: Between(startDate, endDate),
    };

    if (company_id) {
      whereConditions.employee = { company_id };
    }

    // Get all attendance summaries for the month
    const attendanceSummaries = await this.attendanceSummaryRepository.find({
      where: whereConditions,
      relations: ['employee', 'workShift'],
    });

    // Calculate monthly statistics
    const totalWorkingDays = endDate.getDate();
    const totalEmployees = await this.getTotalEmployees(company_id);
    
    const totalPresentDays = attendanceSummaries.filter(
      summary => summary.status === AttendanceStatus.PRESENT
    ).length;
    
    const averageAttendanceRate = totalEmployees > 0 
      ? (totalPresentDays / (totalEmployees * totalWorkingDays)) * 100 
      : 0;

    const totalWorkingHours = attendanceSummaries.reduce(
      (sum, summary) => sum + summary.total_hours, 0
    );

    const totalOvertimeHours = attendanceSummaries.reduce(
      (sum, summary) => sum + summary.overtime_hours, 0
    );

    // Get weekly summaries
    const weeklySummary = await this.getWeeklySummary(startDate, endDate, company_id);
    
    // Get top attendance employees
    const topAttendanceEmployees = await this.getTopAttendanceEmployees(
      startDate, endDate, company_id, 10
    );
    
    // Get low attendance employees
    const lowAttendanceEmployees = await this.getLowAttendanceEmployees(
      startDate, endDate, company_id, 10
    );

    return {
      month,
      total_working_days: totalWorkingDays,
      total_employees: totalEmployees,
      average_attendance_rate: Math.round(averageAttendanceRate * 100) / 100,
      total_working_hours: Math.round(totalWorkingHours / 60), // Convert minutes to hours
      total_overtime_hours: Math.round(totalOvertimeHours / 60),
      weekly_summary: weeklySummary,
      top_attendance_employees: topAttendanceEmployees,
      low_attendance_employees: lowAttendanceEmployees,
    };
  }

  async exportReport(exportDto: ExportReportDto): Promise<ExportReportResponseDto> {
    const { start_date, end_date, format, company_id, email } = exportDto;
    
    const startDate = new Date(start_date);
    const endDate = new Date(end_date);

    // Validate date range
    if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
      throw new BadRequestException('Invalid date format');
    }

    if (startDate > endDate) {
      throw new BadRequestException('Start date must be before end date');
    }

    // Generate job ID
    const jobId = `export_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

    // TODO: Implement actual export logic with queue system
    // For now, return a mock response
    return {
      job_id: jobId,
      status: 'queued',
      message: 'Export job has been queued. You will receive an email when ready.',
      download_url: undefined,
    };
  }

  private async getDepartmentReport(attendanceSummaries: DailyAttendanceSummary[]) {
    // Group by department (assuming department info is available)
    const departmentMap = new Map();
    
    attendanceSummaries.forEach(summary => {
      const departmentName = summary.employee?.company?.name || 'Unknown';
      
      if (!departmentMap.has(departmentName)) {
        departmentMap.set(departmentName, {
          department_name: departmentName,
          total_employees: 0,
          present_employees: 0,
        });
      }
      
      const dept = departmentMap.get(departmentName);
      dept.total_employees++;
      if (summary.status === AttendanceStatus.PRESENT) {
        dept.present_employees++;
      }
    });

    return Array.from(departmentMap.values()).map(dept => ({
      ...dept,
      attendance_rate: dept.total_employees > 0 
        ? Math.round((dept.present_employees / dept.total_employees) * 10000) / 100
        : 0,
    }));
  }

  private async getShiftReport(attendanceSummaries: DailyAttendanceSummary[]) {
    const shiftMap = new Map();
    
    attendanceSummaries.forEach(summary => {
      const shiftName = summary.workShift?.name || 'Unknown';
      
      if (!shiftMap.has(shiftName)) {
        shiftMap.set(shiftName, {
          shift_name: shiftName,
          start_time: summary.workShift?.start_time || '00:00',
          end_time: summary.workShift?.end_time || '00:00',
          total_employees: 0,
          present_employees: 0,
        });
      }
      
      const shift = shiftMap.get(shiftName);
      shift.total_employees++;
      if (summary.status === AttendanceStatus.PRESENT) {
        shift.present_employees++;
      }
    });

    return Array.from(shiftMap.values()).map(shift => ({
      ...shift,
      attendance_rate: shift.total_employees > 0 
        ? Math.round((shift.present_employees / shift.total_employees) * 10000) / 100
        : 0,
    }));
  }

  private async getTotalEmployees(companyId?: string): Promise<number> {
    const whereConditions: any = {};
    if (companyId) {
      whereConditions.company_id = companyId;
    }
    
    return await this.employeeRepository.count({ where: whereConditions });
  }

  private async getWeeklySummary(
    startDate: Date, 
    endDate: Date, 
    companyId?: string
  ) {
    const weeklySummaries = [];
    const currentDate = new Date(startDate);
    
    let weekNumber = 1;
    while (currentDate <= endDate) {
      const weekStart = new Date(currentDate);
      const weekEnd = new Date(currentDate);
      weekEnd.setDate(weekEnd.getDate() + 6);
      
      if (weekEnd > endDate) {
        weekEnd.setTime(endDate.getTime());
      }

      // Get attendance data for this week
      const whereConditions: any = {
        work_date: Between(weekStart, weekEnd),
      };

      if (companyId) {
        whereConditions.employee = { company_id: companyId };
      }

      const weekAttendance = await this.attendanceSummaryRepository.find({
        where: whereConditions,
      });

      const totalPresentDays = weekAttendance.filter(
        summary => summary.status === AttendanceStatus.PRESENT
      ).length;

      const totalEmployees = await this.getTotalEmployees(companyId);
      const attendanceRate = totalEmployees > 0 
        ? (totalPresentDays / (totalEmployees * 7)) * 100 
        : 0;

      const totalHours = weekAttendance.reduce(
        (sum, summary) => sum + summary.total_hours, 0
      );

      weeklySummaries.push({
        week: weekNumber,
        start_date: weekStart.toISOString().split('T')[0],
        end_date: weekEnd.toISOString().split('T')[0],
        attendance_rate: Math.round(attendanceRate * 100) / 100,
        total_hours: Math.round(totalHours / 60),
      });

      currentDate.setDate(currentDate.getDate() + 7);
      weekNumber++;
    }

    return weeklySummaries;
  }

  private async getTopAttendanceEmployees(
    startDate: Date, 
    endDate: Date, 
    companyId?: string, 
    limit: number = 10
  ) {
    const whereConditions: any = {
      work_date: Between(startDate, endDate),
    };

    if (companyId) {
      whereConditions.employee = { company_id: companyId };
    }

    const attendanceSummaries = await this.attendanceSummaryRepository.find({
      where: whereConditions,
      relations: ['employee', 'employee.user'],
    });

    // Group by employee and calculate statistics
    const employeeMap = new Map();
    
    attendanceSummaries.forEach(summary => {
      const employeeId = summary.employee_id;
      
      if (!employeeMap.has(employeeId)) {
        employeeMap.set(employeeId, {
          employee_code: summary.employee.employee_code,
          full_name: summary.employee.user?.full_name || 'Unknown',
          present_days: 0,
          total_hours: 0,
        });
      }
      
      const emp = employeeMap.get(employeeId);
      if (summary.status === AttendanceStatus.PRESENT) {
        emp.present_days++;
      }
      emp.total_hours += summary.total_hours;
    });

    return Array.from(employeeMap.values())
      .sort((a, b) => b.present_days - a.present_days)
      .slice(0, limit)
      .map(emp => ({
        ...emp,
        attendance_rate: Math.round((emp.present_days / 30) * 10000) / 100, // Assuming 30 days
        total_hours: Math.round(emp.total_hours / 60),
      }));
  }

  private async getLowAttendanceEmployees(
    startDate: Date, 
    endDate: Date, 
    companyId?: string, 
    limit: number = 10
  ) {
    // Similar to getTopAttendanceEmployees but sorted by lowest attendance
    const topEmployees = await this.getTopAttendanceEmployees(startDate, endDate, companyId, 100);
    return topEmployees
      .sort((a, b) => a.attendance_rate - b.attendance_rate)
      .slice(0, limit);
  }
}
