// src/features/attendance/Attendance.tsx

import { useEffect, useMemo, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import { Page } from '@/ui/Page';
import { DataTable, type DataTableColumn } from '@/ui/DataTable';
import { useUi } from '@/app/providers/UiProvider';

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

const MOCK_ATTENDANCE: AttendanceRecord[] = [
  {
    id: '1',
    employeeCode: 'EMP001',
    employeeName: 'Nguyễn Văn A',
    department: 'IT',
    checkInTime: '2024-10-05T08:00:00Z',
    checkOutTime: '2024-10-05T17:30:00Z',
    status: 'on_time',
    workHours: 9.5,
    date: '2024-10-05',
  },
  {
    id: '2',
    employeeCode: 'EMP002',
    employeeName: 'Trần Thị B',
    department: 'HR',
    checkInTime: '2024-10-05T08:15:00Z',
    checkOutTime: '2024-10-05T17:00:00Z',
    status: 'late',
    workHours: 8.75,
    date: '2024-10-05',
  },
  {
    id: '3',
    employeeCode: 'EMP003',
    employeeName: 'Lê Văn C',
    department: 'Sales',
    checkInTime: '2024-10-05T07:45:00Z',
    checkOutTime: '2024-10-05T16:30:00Z',
    status: 'early',
    workHours: 8.75,
    date: '2024-10-05',
  },
  {
    id: '4',
    employeeCode: 'EMP004',
    employeeName: 'Phạm Thị D',
    department: 'Finance',
    checkInTime: '2024-10-05T08:00:00Z',
    checkOutTime: null,
    status: 'on_time',
    workHours: 0,
    date: '2024-10-05',
  },
  {
    id: '5',
    employeeCode: 'EMP005',
    employeeName: 'Đỗ Minh E',
    department: 'Marketing',
    checkInTime: '2024-10-05T08:40:00Z',
    checkOutTime: '2024-10-05T17:40:00Z',
    status: 'late',
    workHours: 8.2,
    date: '2024-10-05',
  },
  {
    id: '6',
    employeeCode: 'EMP006',
    employeeName: 'Phan Quốc F',
    department: 'IT',
    checkInTime: '2024-10-05T07:55:00Z',
    checkOutTime: '2024-10-05T17:05:00Z',
    status: 'on_time',
    workHours: 9.1,
    date: '2024-10-05',
  },
];

const STATUS_CONFIG: Record<AttendanceRecord['status'], { label: string; variant: string }> = {
  on_time: { label: 'Đúng giờ', variant: 'success' },
  late: { label: 'Đi trễ', variant: 'warning' },
  early: { label: 'Về sớm', variant: 'info' },
  absent: { label: 'Vắng mặt', variant: 'danger' },
};

const PAGE_SIZE = 8;

export default function Attendance() {
  const { showToast } = useUi();
  const [records, setRecords] = useState<AttendanceRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [dateFilter, setDateFilter] = useState(new Date().toISOString().split('T')[0]);
  const [statusFilter, setStatusFilter] = useState<'all' | AttendanceRecord['status']>('all');
  const [page, setPage] = useState(1);

  useEffect(() => {
    const loadAttendanceRecords = async () => {
      setLoading(true);
      try {
        await new Promise((resolve) => setTimeout(resolve, 300));
        setRecords(MOCK_ATTENDANCE);
      } finally {
        setLoading(false);
      }
    };

    loadAttendanceRecords();
  }, []);

  const filteredRecords = useMemo(() => {
    return records.filter((record) => {
      const matchesSearch =
        record.employeeName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        record.employeeCode.toLowerCase().includes(searchQuery.toLowerCase()) ||
        record.department.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = statusFilter === 'all' || record.status === statusFilter;
      const matchesDate = record.date === dateFilter;

      return matchesSearch && matchesStatus && matchesDate;
    });
  }, [records, searchQuery, statusFilter, dateFilter]);

  const paginatedRecords = useMemo(() => {
    const start = (page - 1) * PAGE_SIZE;
    return filteredRecords.slice(start, start + PAGE_SIZE);
  }, [filteredRecords, page]);

  useEffect(() => {
    setPage(1);
  }, [searchQuery, statusFilter, dateFilter]);

  const statusSummary = useMemo(() => {
    return filteredRecords.reduce(
      (acc, record) => {
        acc.total += 1;
        acc[record.status] = (acc[record.status] ?? 0) + 1;
        return acc;
      },
      { total: 0, on_time: 0, late: 0, early: 0, absent: 0 } as Record<string, number>
    );
  }, [filteredRecords]);

  const handleExport = (format: 'excel' | 'csv') => {
    showToast({
      variant: 'info',
      message: `Đang chuẩn bị xuất dữ liệu (${format.toUpperCase()}) cho ngày ${new Date(
        dateFilter
      ).toLocaleDateString('vi-VN')}.`,
    });
  };

  const columns: DataTableColumn<AttendanceRecord>[] = [
    {
      header: 'Mã NV',
      render: (row) => <span className="fw-semibold text-uppercase">{row.employeeCode}</span>,
    },
    {
      header: 'Tên nhân viên',
      render: (row) => <span className="fw-semibold text-dark">{row.employeeName}</span>,
    },
    {
      header: 'Phòng ban',
      render: (row) => <span className="text-secondary">{row.department}</span>,
    },
    {
      header: 'Giờ vào',
      render: (row) => formatTime(row.checkInTime),
      className: 'text-nowrap',
    },
    {
      header: 'Giờ ra',
      render: (row) => formatTime(row.checkOutTime),
      className: 'text-nowrap',
    },
    {
      header: 'Giờ công',
      render: (row) => (
        <span className="fw-semibold text-primary">{row.workHours > 0 ? `${row.workHours.toFixed(1)}h` : '-'}</span>
      ),
      className: 'text-center',
    },
    {
      header: 'Trạng thái',
      render: (row) => {
        const config = STATUS_CONFIG[row.status];
        return <Badge bg={config.variant}>{config.label}</Badge>;
      },
      className: 'text-center',
    },
  ];

  return (
    <Page
      title="Chấm công"
      subtitle="Theo dõi giờ vào ra của nhân viên theo thời gian thực"
      breadcrumb={[{ label: 'Trang chủ', path: '/dashboard' }, { label: 'Chấm công' }]}
      actions={
        <Stack direction="horizontal" gap={2}>
          <Button variant="outline-primary" onClick={() => handleExport('csv')}>
            Xuất CSV
          </Button>
          <Button variant="primary" onClick={() => handleExport('excel')}>
            Xuất Excel
          </Button>
        </Stack>
      }
    >
      <Card className="border-0 shadow-sm">
        <Card.Body>
          <Row className="g-3">
            <Col md={4}>
              <Form.Group controlId="attendance-search">
                <Form.Label className="text-secondary small text-uppercase">Tìm kiếm</Form.Label>
                <Form.Control
                  type="search"
                  placeholder="Tên, mã NV, phòng ban..."
                  value={searchQuery}
                  onChange={(event) => setSearchQuery(event.target.value)}
                />
              </Form.Group>
            </Col>
            <Col md={4}>
              <Form.Group controlId="attendance-date">
                <Form.Label className="text-secondary small text-uppercase">Ngày</Form.Label>
                <Form.Control
                  type="date"
                  value={dateFilter}
                  onChange={(event) => setDateFilter(event.target.value)}
                />
              </Form.Group>
            </Col>
            <Col md={4}>
              <Form.Group controlId="attendance-status">
                <Form.Label className="text-secondary small text-uppercase">Trạng thái</Form.Label>
                <Form.Select
                  value={statusFilter}
                  onChange={(event) => setStatusFilter(event.target.value as typeof statusFilter)}
                >
                  <option value="all">Tất cả</option>
                  <option value="on_time">Đúng giờ</option>
                  <option value="late">Đi trễ</option>
                  <option value="early">Về sớm</option>
                  <option value="absent">Vắng mặt</option>
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      <Row className="g-3">
        <Col md={3}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <p className="text-secondary text-uppercase small mb-1">Tổng bản ghi</p>
              <h3 className="fw-semibold mb-0">{statusSummary.total}</h3>
              <p className="text-secondary small mb-0">Ngày {new Date(dateFilter).toLocaleDateString('vi-VN')}</p>
            </Card.Body>
          </Card>
        </Col>
        <Col md={3}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <p className="text-secondary text-uppercase small mb-1">Đúng giờ</p>
              <h3 className="fw-semibold text-success mb-0">{statusSummary.on_time}</h3>
              <p className="text-secondary small mb-0">{percentage(statusSummary.on_time, statusSummary.total)} tổng</p>
            </Card.Body>
          </Card>
        </Col>
        <Col md={3}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <p className="text-secondary text-uppercase small mb-1">Đi trễ</p>
              <h3 className="fw-semibold text-warning mb-0">{statusSummary.late}</h3>
              <p className="text-secondary small mb-0">{percentage(statusSummary.late, statusSummary.total)} tổng</p>
            </Card.Body>
          </Card>
        </Col>
        <Col md={3}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <p className="text-secondary text-uppercase small mb-1">Vắng mặt</p>
              <h3 className="fw-semibold text-danger mb-0">{statusSummary.absent}</h3>
              <p className="text-secondary small mb-0">{percentage(statusSummary.absent, statusSummary.total)} tổng</p>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Card className="border-0 shadow-sm">
        <Card.Header className="bg-transparent border-0 pb-0">
          <Stack
            direction="horizontal"
            className="justify-content-between align-items-center flex-wrap gap-2"
          >
            <div>
              <h2 className="fs-5 fw-semibold mb-0">Danh sách chấm công</h2>
              <p className="text-secondary small mb-0">
                {filteredRecords.length} bản ghi · Cập nhật {new Date().toLocaleTimeString('vi-VN')}
              </p>
            </div>
            <Button variant="outline-secondary" size="sm" onClick={() => setRecords(MOCK_ATTENDANCE)}>
              Làm mới dữ liệu
            </Button>
          </Stack>
        </Card.Header>
        <Card.Body className="pt-3">
          <DataTable
            columns={columns}
            data={paginatedRecords}
            loading={loading}
            page={page}
            pageSize={PAGE_SIZE}
            total={filteredRecords.length}
            onPageChange={setPage}
            keySelector={(row) => row.id}
            emptyMessage="Không có dữ liệu chấm công cho bộ lọc hiện tại"
          />
        </Card.Body>
      </Card>
    </Page>
  );
}

function formatTime(time: string | null) {
  if (!time) {
    return '-';
  }

  return new Date(time).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function percentage(value: number, total: number) {
  if (total === 0) {
    return '0%';
  }

  return `${Math.round((value / total) * 100)}%`;
}
