// src/services/mock/reports.ts

import type { ReportRow, FilterOptions, ApiResponse } from '@/types';

const MOCK_DELAY = 800;

const mockReportData: ReportRow[] = [];

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