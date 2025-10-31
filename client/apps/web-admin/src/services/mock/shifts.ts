// src/services/mock/shifts.ts

import type { Shift, PaginatedResponse, FilterOptions, ApiResponse } from '@/types';

const MOCK_DELAY = 500;

const mockShifts: Shift[] = [];

// Mock employee-shift assignments
const mockEmployeeShifts: { employeeId: string; shiftId: string; assignedAt: string }[] = [];

export async function getShifts(filters: FilterOptions = {}): Promise<PaginatedResponse<Shift>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = [...mockShifts];

  // Apply filters
  if (filters.search) {
    const searchTerm = filters.search.toLowerCase();
    filtered = filtered.filter(shift => 
      shift.name.toLowerCase().includes(searchTerm) ||
      (shift.description && shift.description.toLowerCase().includes(searchTerm))
    );
  }

  if (filters.status) {
    const isActive = filters.status === 'active';
    filtered = filtered.filter(shift => shift.active === isActive);
  }

  // Apply sorting
  if (filters.sortBy) {
    filtered.sort((a, b) => {
      const aValue = a[filters.sortBy as keyof Shift] as any;
      const bValue = b[filters.sortBy as keyof Shift] as any;
      
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

export async function getShift(id: string): Promise<ApiResponse<Shift>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const shift = mockShifts.find(s => s.id === id);
  
  if (!shift) {
    return {
      data: {} as Shift,
      error: 'Không tìm thấy ca làm việc'
    };
  }

  return {
    data: shift
  };
}

export async function createShift(data: Partial<Shift>): Promise<ApiResponse<Shift>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const newShift: Shift = {
    id: String(mockShifts.length + 1),
    name: data.name || '',
    start: data.start || '08:00',
    end: data.end || '17:00',
    description: data.description,
    active: data.active ?? true,
    createdAt: new Date().toISOString()
  };

  mockShifts.push(newShift);

  return {
    data: newShift
  };
}

export async function updateShift(id: string, data: Partial<Shift>): Promise<ApiResponse<Shift>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockShifts.findIndex(s => s.id === id);
  
  if (index === -1) {
    return {
      data: {} as Shift,
      error: 'Không tìm thấy ca làm việc'
    };
  }

  mockShifts[index] = {
    ...mockShifts[index],
    ...data
  };

  return {
    data: mockShifts[index]
  };
}

export async function deleteShift(id: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockShifts.findIndex(s => s.id === id);
  
  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy ca làm việc'
    };
  }

  // Remove shift assignments
  for (let i = mockEmployeeShifts.length - 1; i >= 0; i--) {
    if (mockEmployeeShifts[i].shiftId === id) {
      mockEmployeeShifts.splice(i, 1);
    }
  }

  mockShifts.splice(index, 1);

  return {
    data: undefined
  };
}

export async function getEmployeeShifts(employeeId: string): Promise<ApiResponse<{ shiftId: string; shiftName: string; assignedAt: string }[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const employeeShiftAssignments = mockEmployeeShifts.filter(assignment => 
    assignment.employeeId === employeeId
  );

  const result = employeeShiftAssignments.map(assignment => {
    const shift = mockShifts.find(s => s.id === assignment.shiftId);
    return {
      shiftId: assignment.shiftId,
      shiftName: shift?.name || 'Unknown Shift',
      assignedAt: assignment.assignedAt
    };
  });

  return {
    data: result
  };
}

export async function assignShiftToEmployee(employeeId: string, shiftId: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  // Check if shift exists
  const shift = mockShifts.find(s => s.id === shiftId);
  if (!shift) {
    return {
      data: undefined,
      error: 'Không tìm thấy ca làm việc'
    };
  }

  // Check if assignment already exists
  const existingAssignment = mockEmployeeShifts.find(assignment => 
    assignment.employeeId === employeeId && assignment.shiftId === shiftId
  );

  if (existingAssignment) {
    return {
      data: undefined,
      error: 'Nhân viên đã được gán ca làm việc này'
    };
  }

  // Add new assignment
  mockEmployeeShifts.push({
    employeeId,
    shiftId,
    assignedAt: new Date().toISOString()
  });

  return {
    data: undefined
  };
}

export async function removeShiftFromEmployee(employeeId: string, shiftId: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockEmployeeShifts.findIndex(assignment => 
    assignment.employeeId === employeeId && assignment.shiftId === shiftId
  );

  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy phân công ca làm việc'
    };
  }

  mockEmployeeShifts.splice(index, 1);

  return {
    data: undefined
  };
}

export async function getShiftEmployees(shiftId: string): Promise<ApiResponse<{ employeeId: string; employeeName: string; assignedAt: string }[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const shiftAssignments = mockEmployeeShifts.filter(assignment => 
    assignment.shiftId === shiftId
  );

  // Mock employee names - in real app, you'd join with employee data
  const employeeNames: Record<string, string> = {
    '1': 'Nguyễn Văn An',
    '2': 'Trần Thị Bình',
    '3': 'Lê Minh Cường',
    '5': 'Hoàng Đức Em',
    '6': 'Võ Thị Phương'
  };

  const result = shiftAssignments.map(assignment => ({
    employeeId: assignment.employeeId,
    employeeName: employeeNames[assignment.employeeId] || 'Unknown Employee',
    assignedAt: assignment.assignedAt
  }));

  return {
    data: result
  };
}