// src/services/mock/reports.ts

import type { ReportRow, FilterOptions, ApiResponse } from '../../types';

const MOCK_DELAY = 800;

const mockReportData: ReportRow[] = [
  {
    employeeId: '1',
    employeeName: 'Nguyễn Văn An',
    date: '2024-10-01',
    totalHours: 9.5,
    lateMinutes: 0,
    department: 'Kỹ thuật',
    checkIn: '08:00:00',
    checkOut: '17:30:00'
  },
  {
    employeeId: '2',
    employeeName: 'Trần Thị Bình',
    date: '2024-10-01',
    totalHours: 9.25,
    lateMinutes: 15,
    department: 'Nhân sự',
    checkIn: '08:15:00',
    checkOut: '17:30:00'
  },
  {
    employeeId: '3',
    employeeName: 'Lê Minh Cường',
    date: '2024-10-01',
    totalHours: 9,
    lateMinutes: 0,
    department: 'Marketing',
    checkIn: '09:00:00',
    checkOut: '18:00:00'
  },
  {
    employeeId: '5',
    employeeName: 'Hoàng Đức Em',
    date: '2024-10-01',
    totalHours: 9.5,
    lateMinutes: 5,
    department: 'Kỹ thuật',
    checkIn: '08:05:00',
    checkOut: '17:35:00'
  },
  {
    employeeId: '6',
    employeeName: 'Võ Thị Phương',
    date: '2024-10-01',
    totalHours: 9,
    lateMinutes: 30,
    department: 'Marketing',
    checkIn: '09:30:00',
    checkOut: '18:30:00'
  },
  // October 2nd
  {
    employeeId: '1',
    employeeName: 'Nguyễn Văn An',
    date: '2024-10-02',
    totalHours: 9.75,
    lateMinutes: 0,
    department: 'Kỹ thuật',
    checkIn: '07:55:00',
    checkOut: '17:40:00'
  },
  {
    employeeId: '2',
    employeeName: 'Trần Thị Bình',
    date: '2024-10-02',
    totalHours: 9.5,
    lateMinutes: 10,
    department: 'Nhân sự',
    checkIn: '08:10:00',
    checkOut: '17:40:00'
  },
  {
    employeeId: '3',
    employeeName: 'Lê Minh Cường',
    date: '2024-10-02',
    totalHours: 8.5,
    lateMinutes: 0,
    department: 'Marketing',
    checkIn: '09:00:00',
    checkOut: '17:30:00'
  },
  {
    employeeId: '5',
    employeeName: 'Hoàng Đức Em',
    date: '2024-10-02',
    totalHours: 9.5,
    lateMinutes: 0,
    department: 'Kỹ thuật',
    checkIn: '08:00:00',
    checkOut: '17:30:00'
  },
  // October 3rd
  {
    employeeId: '1',
    employeeName: 'Nguyễn Văn An',
    date: '2024-10-03',
    totalHours: 9.5,
    lateMinutes: 0,
    department: 'Kỹ thuật',
    checkIn: '08:00:00',
    checkOut: '17:30:00'
  },
  {
    employeeId: '2',
    employeeName: 'Trần Thị Bình',
    date: '2024-10-03',
    totalHours: 9.17,
    lateMinutes: 20,
    department: 'Nhân sự',
    checkIn: '08:20:00',
    checkOut: '17:30:00'
  },
  {
    employeeId: '6',
    employeeName: 'Võ Thị Phương',
    date: '2024-10-03',
    totalHours: 9,
    lateMinutes: 0,
    department: 'Marketing',
    checkIn: '09:00:00',
    checkOut: '18:00:00'
  }
];

export async function generateReport(filters: FilterOptions = {}): Promise<ApiResponse<ReportRow[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = [...mockReportData];

  // Apply filters
  if (filters.startDate && filters.endDate) {
    filtered = filtered.filter(row => 
      row.date >= filters.startDate! && row.date <= filters.endDate!
    );
  } else if (filters.startDate) {
    filtered = filtered.filter(row => row.date >= filters.startDate!);
  } else if (filters.endDate) {
    filtered = filtered.filter(row => row.date <= filters.endDate!);
  }

  if (filters.department) {
    filtered = filtered.filter(row => row.department === filters.department);
  }

  if (filters.search) {
    const searchTerm = filters.search.toLowerCase();
    filtered = filtered.filter(row => 
      row.employeeName.toLowerCase().includes(searchTerm) ||
      row.department?.toLowerCase().includes(searchTerm)
    );
  }

  // Apply sorting
  if (filters.sortBy) {
    filtered.sort((a, b) => {
      const aValue = a[filters.sortBy as keyof ReportRow] as any;
      const bValue = b[filters.sortBy as keyof ReportRow] as any;
      
      if (!aValue && !bValue) return 0;
      if (!aValue) return 1;
      if (!bValue) return -1;
      
      if (filters.sortOrder === 'desc') {
        return bValue > aValue ? 1 : -1;
      }
      return aValue > bValue ? 1 : -1;
    });
  }

  return {
    data: filtered
  };
}

export async function getReportSummary(filters: FilterOptions = {}): Promise<ApiResponse<{
  totalRecords: number;
  totalHours: number;
  totalLateMinutes: number;
  averageHours: number;
  latePercentage: number;
}>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const reportData = await generateReport(filters);
  const data = reportData.data;

  if (!data || data.length === 0) {
    return {
      data: {
        totalRecords: 0,
        totalHours: 0,
        totalLateMinutes: 0,
        averageHours: 0,
        latePercentage: 0
      }
    };
  }

  const totalRecords = data.length;
  const totalHours = data.reduce((sum, row) => sum + row.totalHours, 0);
  const totalLateMinutes = data.reduce((sum, row) => sum + row.lateMinutes, 0);
  const averageHours = totalHours / totalRecords;
  const lateRecords = data.filter(row => row.lateMinutes > 0).length;
  const latePercentage = (lateRecords / totalRecords) * 100;

  return {
    data: {
      totalRecords,
      totalHours: Math.round(totalHours * 100) / 100,
      totalLateMinutes,
      averageHours: Math.round(averageHours * 100) / 100,
      latePercentage: Math.round(latePercentage * 100) / 100
    }
  };
}

export async function getDepartmentStats(filters: FilterOptions = {}): Promise<ApiResponse<{
  department: string;
  employeeCount: number;
  totalHours: number;
  averageHours: number;
  latePercentage: number;
}[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const reportData = await generateReport(filters);
  const data = reportData.data;

  if (!data || data.length === 0) {
    return { data: [] };
  }

  const departmentStats = data.reduce((acc, row) => {
    if (!row.department) return acc;

    if (!acc[row.department]) {
      acc[row.department] = {
        department: row.department,
        employeeIds: new Set(),
        totalHours: 0,
        totalRecords: 0,
        lateRecords: 0
      };
    }

    acc[row.department].employeeIds.add(row.employeeId);
    acc[row.department].totalHours += row.totalHours;
    acc[row.department].totalRecords++;
    if (row.lateMinutes > 0) {
      acc[row.department].lateRecords++;
    }

    return acc;
  }, {} as Record<string, any>);

  const result = Object.values(departmentStats).map((dept: any) => ({
    department: dept.department,
    employeeCount: dept.employeeIds.size,
    totalHours: Math.round(dept.totalHours * 100) / 100,
    averageHours: Math.round((dept.totalHours / dept.totalRecords) * 100) / 100,
    latePercentage: Math.round((dept.lateRecords / dept.totalRecords) * 100 * 100) / 100
  }));

  return { data: result };
}