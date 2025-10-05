// src/utils/csv.ts

export interface CSVColumn {
  key: string;
  header: string;
  formatter?: (value: any) => string;
}

export function exportToCSV<T extends Record<string, any>>(
  data: T[],
  columns: CSVColumn[],
  filename: string
): void {
  if (data.length === 0) {
    console.warn('No data to export');
    return;
  }

  // Create CSV content
  const headers = columns.map(col => col.header);
  const csvRows = [
    headers.join(','),
    ...data.map(row => 
      columns.map(col => {
        let value = row[col.key];
        
        // Apply formatter if provided
        if (col.formatter) {
          value = col.formatter(value);
        }
        
        // Handle null/undefined values
        if (value == null) {
          value = '';
        }
        
        // Convert to string and escape commas/quotes
        const stringValue = String(value);
        if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
          return `"${stringValue.replace(/"/g, '""')}"`;
        }
        
        return stringValue;
      }).join(',')
    )
  ];

  const csvContent = csvRows.join('\n');
  
  // Create and download file
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  
  if (link.download !== undefined) {
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `${filename}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  } else {
    console.error('CSV export not supported in this browser');
  }
}

export function generateEmployeeReportCSV<T extends Record<string, any>>(
  data: T[],
  filename: string = 'employee-report'
): void {
  const columns: CSVColumn[] = [
    { key: 'employeeName', header: 'Tên nhân viên' },
    { key: 'date', header: 'Ngày' },
    { key: 'checkIn', header: 'Giờ vào' },
    { key: 'checkOut', header: 'Giờ ra' },
    { key: 'totalHours', header: 'Tổng giờ làm', formatter: (value) => `${value || 0}h` },
    { key: 'lateMinutes', header: 'Phút đi trễ', formatter: (value) => `${value || 0}` },
    { key: 'department', header: 'Phòng ban' }
  ];

  exportToCSV(data, columns, filename);
}

export function generateAttendanceReportCSV<T extends Record<string, any>>(
  data: T[],
  filename: string = 'attendance-report'
): void {
  const columns: CSVColumn[] = [
    { key: 'employeeName', header: 'Nhân viên' },
    { key: 'date', header: 'Ngày' },
    { key: 'checkIn', header: 'Check In' },
    { key: 'checkOut', header: 'Check Out' },
    { key: 'totalHours', header: 'Tổng giờ' },
    { key: 'isLate', header: 'Trễ giờ', formatter: (value) => value ? 'Có' : 'Không' },
    { key: 'deviceId', header: 'Thiết bị' }
  ];

  exportToCSV(data, columns, filename);
}

export function generateDeviceReportCSV<T extends Record<string, any>>(
  data: T[],
  filename: string = 'device-report'
): void {
  const columns: CSVColumn[] = [
    { key: 'name', header: 'Tên thiết bị' },
    { key: 'location', header: 'Vị trí' },
    { key: 'status', header: 'Trạng thái' },
    { key: 'lastSyncAt', header: 'Đồng bộ lần cuối' },
    { key: 'model', header: 'Model' },
    { key: 'ipAddress', header: 'IP Address' }
  ];

  exportToCSV(data, columns, filename);
}