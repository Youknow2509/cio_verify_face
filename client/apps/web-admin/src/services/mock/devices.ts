// src/services/mock/devices.ts

import type { Device, PaginatedResponse, FilterOptions, ApiResponse } from '@/types';

const MOCK_DELAY = 500;

const mockDevices: Device[] = [
  {
    id: '1',
    name: 'Thiết bị cửa chính',
    location: 'Tầng 1 - Lối vào chính',
    status: 'online',
    model: 'FaceReader-Pro-X1',
    ipAddress: '192.168.1.101',
    lastSyncAt: '2024-10-05T08:30:00Z',
    createdAt: '2024-01-15T10:00:00Z'
  },
  {
    id: '2',
    name: 'Thiết bị phòng IT',
    location: 'Tầng 2 - Phòng IT',
    status: 'online',
    model: 'FaceReader-Lite-A2',
    ipAddress: '192.168.1.102',
    lastSyncAt: '2024-10-05T08:25:00Z',
    createdAt: '2024-02-01T09:15:00Z'
  },
  {
    id: '3',
    name: 'Thiết bị phòng họp',
    location: 'Tầng 3 - Phòng họp A',
    status: 'offline',
    model: 'FaceReader-Standard-B1',
    ipAddress: '192.168.1.103',
    lastSyncAt: '2024-10-04T16:45:00Z',
    createdAt: '2024-01-20T14:30:00Z'
  },
  {
    id: '4',
    name: 'Thiết bị khu vực ăn uống',
    location: 'Tầng 1 - Cafeteria',
    status: 'online',
    model: 'FaceReader-Lite-A2',
    ipAddress: '192.168.1.104',
    lastSyncAt: '2024-10-05T07:50:00Z',
    createdAt: '2024-03-01T11:00:00Z'
  },
  {
    id: '5',
    name: 'Thiết bị khu vực parking',
    location: 'Tầng B1 - Parking',
    status: 'offline',
    model: 'FaceReader-Standard-B1',
    ipAddress: '192.168.1.105',
    lastSyncAt: '2024-10-03T18:20:00Z',
    createdAt: '2024-02-15T08:45:00Z'
  }
];

export async function getDevices(filters: FilterOptions = {}): Promise<PaginatedResponse<Device>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  let filtered = [...mockDevices];

  // Apply filters
  if (filters.search) {
    const searchTerm = filters.search.toLowerCase();
    filtered = filtered.filter(device => 
      device.name.toLowerCase().includes(searchTerm) ||
      (device.location && device.location.toLowerCase().includes(searchTerm)) ||
      (device.model && device.model.toLowerCase().includes(searchTerm)) ||
      (device.ipAddress && device.ipAddress.includes(searchTerm))
    );
  }

  if (filters.status) {
    filtered = filtered.filter(device => device.status === filters.status);
  }

  // Apply sorting
  if (filters.sortBy) {
    filtered.sort((a, b) => {
      const aValue = a[filters.sortBy as keyof Device] as any;
      const bValue = b[filters.sortBy as keyof Device] as any;
      
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

export async function getDevice(id: string): Promise<ApiResponse<Device>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const device = mockDevices.find(d => d.id === id);
  
  if (!device) {
    return {
      data: {} as Device,
      error: 'Không tìm thấy thiết bị'
    };
  }

  return {
    data: device
  };
}

export async function createDevice(data: Partial<Device>): Promise<ApiResponse<Device>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const newDevice: Device = {
    id: String(mockDevices.length + 1),
    name: data.name || '',
    location: data.location,
    status: data.status || 'offline',
    model: data.model,
    ipAddress: data.ipAddress,
    lastSyncAt: undefined,
    createdAt: new Date().toISOString()
  };

  mockDevices.push(newDevice);

  return {
    data: newDevice
  };
}

export async function updateDevice(id: string, data: Partial<Device>): Promise<ApiResponse<Device>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockDevices.findIndex(d => d.id === id);
  
  if (index === -1) {
    return {
      data: {} as Device,
      error: 'Không tìm thấy thiết bị'
    };
  }

  mockDevices[index] = {
    ...mockDevices[index],
    ...data
  };

  return {
    data: mockDevices[index]
  };
}

export async function deleteDevice(id: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const index = mockDevices.findIndex(d => d.id === id);
  
  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy thiết bị'
    };
  }

  mockDevices.splice(index, 1);

  return {
    data: undefined
  };
}

export async function syncDevice(id: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY + 800)); // Simulate sync delay

  const index = mockDevices.findIndex(d => d.id === id);
  
  if (index === -1) {
    return {
      data: undefined,
      error: 'Không tìm thấy thiết bị'
    };
  }

  // Update last sync time and potentially status
  mockDevices[index].lastSyncAt = new Date().toISOString();
  mockDevices[index].status = Math.random() > 0.2 ? 'online' : 'offline'; // 80% success rate

  return {
    data: undefined
  };
}

export async function testDevice(id: string): Promise<ApiResponse<{ success: boolean; message: string }>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY + 600)); // Simulate test delay

  const device = mockDevices.find(d => d.id === id);
  
  if (!device) {
    return {
      data: { success: false, message: 'Không tìm thấy thiết bị' },
      error: 'Không tìm thấy thiết bị'
    };
  }

  const success = Math.random() > 0.1; // 90% success rate
  
  return {
    data: {
      success,
      message: success 
        ? `Thiết bị ${device.name} hoạt động bình thường`
        : `Thiết bị ${device.name} không phản hồi. Vui lòng kiểm tra kết nối.`
    }
  };
}