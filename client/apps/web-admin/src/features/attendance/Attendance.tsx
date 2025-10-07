// src/features/attendance/Attendance.tsx

import { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/Card/Card';
import { Table } from '@/components/Table/Table';
import { Badge } from '@/components/Badge/Badge';
import { Toolbar, ToolbarSection, SearchBox } from '@/components/Toolbar/Toolbar';
import styles from './Attendance.module.scss';

interface AttendanceRecord {
  id: string;
  employeeCode: string;
  employeeName: string;
  department: string;
  checkInTime: string;
  checkOutTime: string | null;
  status: 'on_time' | 'late' | 'early' | 'absent';
  workHours: number;
  date: string;
}

export default function Attendance() {
  const [records, setRecords] = useState<AttendanceRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [dateFilter, setDateFilter] = useState(new Date().toISOString().split('T')[0]);
  const [statusFilter, setStatusFilter] = useState<string>('all');

  useEffect(() => {
    loadAttendanceRecords();
  }, [dateFilter, statusFilter]);

  const loadAttendanceRecords = async () => {
    setLoading(true);
    try {
      // Mock data - replace with actual API call
      const mockData: AttendanceRecord[] = [
        {
          id: '1',
          employeeCode: 'EMP001',
          employeeName: 'Nguyễn Văn A',
          department: 'IT',
          checkInTime: '2024-10-05 08:00:00',
          checkOutTime: '2024-10-05 17:30:00',
          status: 'on_time',
          workHours: 9.5,
          date: '2024-10-05'
        },
        {
          id: '2',
          employeeCode: 'EMP002',
          employeeName: 'Trần Thị B',
          department: 'HR',
          checkInTime: '2024-10-05 08:15:00',
          checkOutTime: '2024-10-05 17:00:00',
          status: 'late',
          workHours: 8.75,
          date: '2024-10-05'
        },
        {
          id: '3',
          employeeCode: 'EMP003',
          employeeName: 'Lê Văn C',
          department: 'Sales',
          checkInTime: '2024-10-05 07:45:00',
          checkOutTime: '2024-10-05 16:30:00',
          status: 'early',
          workHours: 8.75,
          date: '2024-10-05'
        },
        {
          id: '4',
          employeeCode: 'EMP004',
          employeeName: 'Phạm Thị D',
          department: 'Finance',
          checkInTime: '2024-10-05 08:00:00',
          checkOutTime: null,
          status: 'on_time',
          workHours: 0,
          date: '2024-10-05'
        }
      ];
      
      setRecords(mockData);
    } catch (error) {
      console.error('Error loading attendance records:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredRecords = records.filter(record => {
    const matchesSearch = 
      record.employeeName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      record.employeeCode.toLowerCase().includes(searchQuery.toLowerCase()) ||
      record.department.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesStatus = statusFilter === 'all' || record.status === statusFilter;
    
    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { variant: 'success' | 'warning' | 'error' | 'info' | 'neutral', label: string }> = {
      on_time: { variant: 'success', label: 'Đúng giờ' },
      late: { variant: 'warning', label: 'Đi trễ' },
      early: { variant: 'info', label: 'Về sớm' },
      absent: { variant: 'error', label: 'Vắng mặt' }
    };
    
    const config = statusConfig[status] || { variant: 'neutral' as const, label: status };
    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  const formatTime = (timeString: string | null) => {
    if (!timeString) return '-';
    return new Date(timeString).toLocaleTimeString('vi-VN', { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const columns = [
    {
      key: 'employeeCode',
      header: 'Mã NV',
      render: (value: string) => <span className={styles.employeeCode}>{value}</span>
    },
    {
      key: 'employeeName',
      header: 'Tên nhân viên',
      render: (value: string) => <strong>{value}</strong>
    },
    {
      key: 'department',
      header: 'Phòng ban'
    },
    {
      key: 'checkInTime',
      header: 'Giờ vào',
      render: (value: string) => formatTime(value)
    },
    {
      key: 'checkOutTime',
      header: 'Giờ ra',
      render: (value: string | null) => formatTime(value)
    },
    {
      key: 'workHours',
      header: 'Giờ công',
      render: (value: number) => (
        <span className={styles.workHours}>
          {value > 0 ? `${value.toFixed(1)}h` : '-'}
        </span>
      )
    },
    {
      key: 'status',
      header: 'Trạng thái',
      render: (value: string) => getStatusBadge(value)
    }
  ];

  return (
    <div className={styles.attendance}>
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>Chấm công</h1>
          <p className={styles.subtitle}>Theo dõi giờ vào ra của nhân viên</p>
        </div>
      </div>

      <Toolbar>
        <ToolbarSection>
          <SearchBox
            value={searchQuery}
            onChange={setSearchQuery}
            placeholder="Tìm theo tên, mã NV, phòng ban..."
          />
          
          <div className={styles.filterGroup}>
            <label htmlFor="date-filter" className={styles.filterLabel}>Ngày:</label>
            <input
              id="date-filter"
              type="date"
              value={dateFilter}
              onChange={(e) => setDateFilter(e.target.value)}
              className={styles.dateInput}
            />
          </div>
          
          <div className={styles.filterGroup}>
            <label htmlFor="status-filter" className={styles.filterLabel}>Trạng thái:</label>
            <select
              id="status-filter"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className={styles.statusSelect}
            >
              <option value="all">Tất cả</option>
              <option value="on_time">Đúng giờ</option>
              <option value="late">Đi trễ</option>
              <option value="early">Về sớm</option>
              <option value="absent">Vắng mặt</option>
            </select>
          </div>
        </ToolbarSection>
        
        <ToolbarSection align="right">
          <button className={styles.exportButton}>
            <ExportIcon />
            Xuất Excel
          </button>
          <button className={styles.refreshButton} onClick={loadAttendanceRecords}>
            <RefreshIcon />
            Làm mới
          </button>
        </ToolbarSection>
      </Toolbar>

      <Card>
        <CardHeader>
          <CardTitle>
            Danh sách chấm công - {new Date(dateFilter).toLocaleDateString('vi-VN')}
            <span className={styles.recordCount}>({filteredRecords.length} bản ghi)</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table
            columns={columns}
            data={filteredRecords}
            loading={loading}
          />
        </CardContent>
      </Card>
    </div>
  );
}

// Icons
const ExportIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2z"/>
  </svg>
);

const RefreshIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
  </svg>
);
