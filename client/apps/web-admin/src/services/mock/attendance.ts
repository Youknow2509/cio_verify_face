// src/services/mock/attendance.ts

import type { AttendanceRecord, PaginatedResponse, FilterOptions, ApiResponse, DashboardStats, ChartData, RecentActivity } from '@/types';

const MOCK_DELAY = 700;

const mockAttendanceRecords: AttendanceRecord[] = [
  {
    id: '1',
    employeeId: '1',
    employeeName: 'Nguyễn Văn An',
    date: '2024-10-05',
    checkIn: '08:00:00',
    checkOut: '17:30:00',
    totalHours: 9.5,
    isLate: false,
    shiftId: '1',
    deviceId: '1'
  },
  {
    id: '2',
    employeeId: '2',
    employeeName: 'Trần Thị Bình',
    date: '2024-10-05',
    checkIn: '08:15:00',
    checkOut: '17:45:00',
    totalHours: 9.5,
    isLate: true,
    shiftId: '1',
    deviceId: '1'
  },
  {
    id: '3',
    employeeId: '3',
    employeeName: 'Lê Minh Cường',
    date: '2024-10-05',
    checkIn: '09:00:00',
    checkOut: '18:00:00',
    totalHours: 9,
    isLate: false,
    shiftId: '2',
    deviceId: '2'
  },
  {
    id: '4',
    employeeId: '5',
    employeeName: 'Hoàng Đức Em',
    date: '2024-10-05',
    checkIn: '08:05:00',
    checkOut: undefined,
    totalHours: undefined,
    isLate: false,
    shiftId: '1',
    deviceId: '1'
  },
  {
    id: '5',
    employeeId: '1',
    employeeName: 'Nguyễn Văn An',
    date: '2024-10-04',
    checkIn: '07:55:00',
    checkOut: '17:25:00',
    totalHours: 9.5,
    isLate: false,
    shiftId: '1',
    deviceId: '1'
  },
  {
    id: '6',
    employeeId: '2',
    employeeName: 'Trần Thị Bình',
    date: '2024-10-04',
    checkIn: '08:20:00',
    checkOut: '17:30:00',
    totalHours: 9.17,
    isLate: true,
    shiftId: '1',
    deviceId: '1'
  },
  {
    id: '7',
    employeeId: '6',
    employeeName: 'Võ Thị Phương',
    date: '2024-10-04',
    checkIn: '09:30:00',
    checkOut: '18:30:00',
    totalHours: 9,
    isLate: false,
    shiftId: '2',
    deviceId: '4'
  }
];

export async function getAttendanceRecords(filters: FilterOptions = {}): Promise<PaginatedResponse<AttendanceRecord>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = [...mockAttendanceRecords];

  // Apply filters
  if (filters.search) {
    const searchTerm = filters.search.toLowerCase();
    filtered = filtered.filter(record => 
      record.employeeName.toLowerCase().includes(searchTerm) ||
      record.date.includes(searchTerm)
    );
  }

  if (filters.startDate && filters.endDate) {
    filtered = filtered.filter(record => 
      record.date >= filters.startDate! && record.date <= filters.endDate!
    );
  } else if (filters.startDate) {
    filtered = filtered.filter(record => record.date >= filters.startDate!);
  } else if (filters.endDate) {
    filtered = filtered.filter(record => record.date <= filters.endDate!);
  }

  if (filters.department) {
    // In a real app, you'd join with employee data
    // For mock, we'll filter by employee names that might belong to departments
    const techNames = ['Nguyễn Văn An', 'Hoàng Đức Em'];
    const hrNames = ['Trần Thị Bình'];
    const marketingNames = ['Lê Minh Cường', 'Võ Thị Phương'];
    
    switch (filters.department) {
      case 'Kỹ thuật':
        filtered = filtered.filter(record => techNames.includes(record.employeeName));
        break;
      case 'Nhân sự':
        filtered = filtered.filter(record => hrNames.includes(record.employeeName));
        break;
      case 'Marketing':
        filtered = filtered.filter(record => marketingNames.includes(record.employeeName));
        break;
    }
  }

  // Apply sorting
  if (filters.sortBy) {
    filtered.sort((a, b) => {
      const aValue = a[filters.sortBy as keyof AttendanceRecord] as any;
      const bValue = b[filters.sortBy as keyof AttendanceRecord] as any;
      
      if (!aValue && !bValue) return 0;
      if (!aValue) return 1;
      if (!bValue) return -1;
      
      if (filters.sortOrder === 'desc') {
        return bValue > aValue ? 1 : -1;
      }
      return aValue > bValue ? 1 : -1;
    });
  }

  // Apply pagination
  const page = filters.page || 1;
  const limit = filters.limit || 10;
  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;
  const paginatedData = filtered.slice(startIndex, endIndex);

  return {
    data: paginatedData,
    total: filtered.length,
    page,
    limit,
    totalPages: Math.ceil(filtered.length / limit)
  };
}

export async function getEmployeeAttendance(employeeId: string, filters: FilterOptions = {}): Promise<PaginatedResponse<AttendanceRecord>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = mockAttendanceRecords.filter(record => record.employeeId === employeeId);

  if (filters.startDate && filters.endDate) {
    filtered = filtered.filter(record => 
      record.date >= filters.startDate! && record.date <= filters.endDate!
    );
  }

  // Apply pagination
  const page = filters.page || 1;
  const limit = filters.limit || 10;
  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;
  const paginatedData = filtered.slice(startIndex, endIndex);

  return {
    data: paginatedData,
    total: filtered.length,
    page,
    limit,
    totalPages: Math.ceil(filtered.length / limit)
  };
}

export async function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const today = new Date().toISOString().split('T')[0];
  const todayRecords = mockAttendanceRecords.filter(record => record.date === today);

  const stats: DashboardStats = {
    totalEmployees: 6,
    todayCheckIns: todayRecords.length,
    lateArrivals: todayRecords.filter(record => record.isLate).length,
    devicesOnline: 3,
    attendanceRate: 85.5
  };

  return {
    data: stats
  };
}

export async function getAttendanceChart(): Promise<ApiResponse<ChartData[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  // Generate 7 days of mock chart data
  const data: ChartData[] = [];
  const today = new Date();
  
  for (let i = 6; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    const dateStr = date.toISOString().split('T')[0];
    
    const dayRecords = mockAttendanceRecords.filter(record => record.date === dateStr);
    
    data.push({
      date: dateStr,
      checkIns: dayRecords.length,
      checkOuts: dayRecords.filter(record => record.checkOut).length,
      lateArrivals: dayRecords.filter(record => record.isLate).length
    });
  }

  return {
    data
  };
}

export async function getRecentActivity(): Promise<ApiResponse<RecentActivity[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const activities: RecentActivity[] = [
    {
      id: '1',
      type: 'check_in',
      message: 'đã check-in',
      timestamp: '2024-10-05T08:00:00Z',
      employeeName: 'Nguyễn Văn An'
    },
    {
      id: '2',
      type: 'check_in',
      message: 'đã check-in (trễ 15 phút)',
      timestamp: '2024-10-05T08:15:00Z',
      employeeName: 'Trần Thị Bình'
    },
    {
      id: '3',
      type: 'device_sync',
      message: 'đã đồng bộ thành công',
      timestamp: '2024-10-05T08:30:00Z',
      deviceName: 'Thiết bị cửa chính'
    },
    {
      id: '4',
      type: 'check_in',
      message: 'đã check-in',
      timestamp: '2024-10-05T08:05:00Z',
      employeeName: 'Hoàng Đức Em'
    },
    {
      id: '5',
      type: 'employee_added',
      message: 'đã được thêm vào hệ thống',
      timestamp: '2024-10-04T16:30:00Z',
      employeeName: 'Võ Thị Phương'
    }
  ];

  return {
    data: activities
  };
}

export async function createAttendanceRecord(data: Partial<AttendanceRecord>): Promise<ApiResponse<AttendanceRecord>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const newRecord: AttendanceRecord = {
    id: String(mockAttendanceRecords.length + 1),
    employeeId: data.employeeId || '',
    employeeName: data.employeeName || '',
    date: data.date || new Date().toISOString().split('T')[0],
    checkIn: data.checkIn,
    checkOut: data.checkOut,
    totalHours: data.totalHours,
    isLate: data.isLate,
    shiftId: data.shiftId,
    deviceId: data.deviceId
  };

  mockAttendanceRecords.push(newRecord);

  return {
    data: newRecord
  };
}

export async function updateAttendanceRecord(id: string, data: Partial<AttendanceRecord>): Promise<ApiResponse<AttendanceRecord>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockAttendanceRecords.findIndex(record => record.id === id);
  
  if (index === -1) {
    return {
      data: {} as AttendanceRecord,
      error: 'Không tìm thấy bản ghi chấm công'
    };
  }

  mockAttendanceRecords[index] = {
    ...mockAttendanceRecords[index],
    ...data
  };

  return {
    data: mockAttendanceRecords[index]
  };
}