// src/services/mock/employees.ts

import type { Employee, FaceData, PaginatedResponse, FilterOptions, ApiResponse } from '@/types';

const MOCK_DELAY = 600;

const mockEmployees: Employee[] = [
  {
    id: '1',
    code: 'EMP001',
    name: 'Nguyễn Văn An',
    email: 'an.nguyen@company.com',
    department: 'Kỹ thuật',
    position: 'Developer',
    active: true,
    faceCount: 3,
    createdAt: '2024-01-15T08:30:00Z',
    updatedAt: '2024-10-01T10:15:00Z'
  },
  {
    id: '2',
    code: 'EMP002',
    name: 'Trần Thị Bình',
    email: 'binh.tran@company.com',
    department: 'Nhân sự',
    position: 'HR Manager',
    active: true,
    faceCount: 2,
    createdAt: '2024-01-10T09:00:00Z',
    updatedAt: '2024-09-25T14:20:00Z'
  },
  {
    id: '3',
    code: 'EMP003',
    name: 'Lê Minh Cường',
    email: 'cuong.le@company.com',
    department: 'Marketing',
    position: 'Marketing Specialist',
    active: true,
    faceCount: 1,
    createdAt: '2024-02-01T07:45:00Z',
    updatedAt: '2024-10-02T16:30:00Z'
  },
  {
    id: '4',
    code: 'EMP004',
    name: 'Phạm Thu Dung',
    email: 'dung.pham@company.com',
    department: 'Kế toán',
    position: 'Accountant',
    active: false,
    faceCount: 0,
    createdAt: '2024-01-20T11:15:00Z',
    updatedAt: '2024-09-15T13:45:00Z'
  },
  {
    id: '5',
    code: 'EMP005',
    name: 'Hoàng Đức Em',
    email: 'em.hoang@company.com',
    department: 'Kỹ thuật',
    position: 'QA Engineer',
    active: true,
    faceCount: 2,
    createdAt: '2024-03-01T08:00:00Z',
    updatedAt: '2024-10-01T09:30:00Z'
  },
  {
    id: '6',
    code: 'EMP006',
    name: 'Võ Thị Phương',
    email: 'phuong.vo@company.com',
    department: 'Marketing',
    position: 'Content Creator',
    active: true,
    faceCount: 4,
    createdAt: '2024-02-15T10:30:00Z',
    updatedAt: '2024-09-30T15:00:00Z'
  }
];

const mockFaceData: FaceData[] = [
  {
    id: 'face-1',
    employeeId: '1',
    imageUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
    fileName: 'face_1_1.jpg',
    createdAt: '2024-01-15T08:35:00Z'
  },
  {
    id: 'face-2',
    employeeId: '1',
    imageUrl: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
    fileName: 'face_1_2.jpg',
    createdAt: '2024-02-01T09:15:00Z'
  },
  {
    id: 'face-3',
    employeeId: '2',
    imageUrl: 'https://images.unsplash.com/photo-1494790108755-2616b612b65c?w=150&h=150&fit=crop&crop=face',
    fileName: 'face_2_1.jpg',
    createdAt: '2024-01-10T09:05:00Z'
  }
];

export async function getEmployees(filters: FilterOptions = {}): Promise<PaginatedResponse<Employee>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = [...mockEmployees];

  // Apply filters
  if (filters.search) {
    const searchTerm = filters.search.toLowerCase();
    filtered = filtered.filter(emp => 
      emp.name.toLowerCase().includes(searchTerm) ||
      emp.code.toLowerCase().includes(searchTerm) ||
      emp.email.toLowerCase().includes(searchTerm)
    );
  }

  if (filters.department) {
    filtered = filtered.filter(emp => emp.department === filters.department);
  }

  if (filters.status) {
    const isActive = filters.status === 'active';
    filtered = filtered.filter(emp => emp.active === isActive);
  }

  // Apply sorting
  if (filters.sortBy) {
    filtered.sort((a, b) => {
      const aValue = a[filters.sortBy as keyof Employee] as any;
      const bValue = b[filters.sortBy as keyof Employee] as any;
      
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

export async function getEmployee(id: string): Promise<ApiResponse<Employee>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const employee = mockEmployees.find(emp => emp.id === id);
  
  if (!employee) {
    return {
      data: {} as Employee,
      error: 'Không tìm thấy nhân viên'
    };
  }

  return {
    data: employee
  };
}

export async function createEmployee(data: Partial<Employee>): Promise<ApiResponse<Employee>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const newEmployee: Employee = {
    id: String(mockEmployees.length + 1),
    code: data.code || `EMP${String(mockEmployees.length + 1).padStart(3, '0')}`,
    name: data.name || '',
    email: data.email || '',
    department: data.department,
    position: data.position,
    active: data.active ?? true,
    faceCount: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  };

  mockEmployees.push(newEmployee);

  return {
    data: newEmployee
  };
}

export async function updateEmployee(id: string, data: Partial<Employee>): Promise<ApiResponse<Employee>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockEmployees.findIndex(emp => emp.id === id);
  
  if (index === -1) {
    return {
      data: {} as Employee,
      error: 'Không tìm thấy nhân viên'
    };
  }

  mockEmployees[index] = {
    ...mockEmployees[index],
    ...data,
    updatedAt: new Date().toISOString()
  };

  return {
    data: mockEmployees[index]
  };
}

export async function deleteEmployee(id: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockEmployees.findIndex(emp => emp.id === id);
  
  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy nhân viên'
    };
  }

  mockEmployees.splice(index, 1);

  return {
    data: undefined
  };
}

export async function getEmployeeFaceData(employeeId: string): Promise<ApiResponse<FaceData[]>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const faceData = mockFaceData.filter(face => face.employeeId === employeeId);

  return {
    data: faceData
  };
}

export async function uploadFaceData(employeeId: string, file: File): Promise<ApiResponse<FaceData>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY + 1000)); // Simulate upload delay

  const newFaceData: FaceData = {
    id: `face-${Date.now()}`,
    employeeId,
    imageUrl: URL.createObjectURL(file),
    fileName: file.name,
    createdAt: new Date().toISOString()
  };

  mockFaceData.push(newFaceData);

  // Update employee face count
  const employee = mockEmployees.find(emp => emp.id === employeeId);
  if (employee) {
    employee.faceCount++;
    employee.updatedAt = new Date().toISOString();
  }

  return {
    data: newFaceData
  };
}

export async function deleteFaceData(faceId: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockFaceData.findIndex(face => face.id === faceId);
  
  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy dữ liệu khuôn mặt'
    };
  }

  const faceData = mockFaceData[index];
  mockFaceData.splice(index, 1);

  // Update employee face count
  const employee = mockEmployees.find(emp => emp.id === faceData.employeeId);
  if (employee) {
    employee.faceCount = Math.max(0, employee.faceCount - 1);
    employee.updatedAt = new Date().toISOString();
  }

  return {
    data: undefined
  };
}