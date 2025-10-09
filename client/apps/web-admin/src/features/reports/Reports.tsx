// src/features/reports/Reports.tsx

import { useMemo, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import ProgressBar from 'react-bootstrap/ProgressBar';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import Table from 'react-bootstrap/Table';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';

const MOCK_DEPARTMENTS = [
  { name: 'IT', total: 42, onTime: 37 },
  { name: 'HR', total: 28, onTime: 22 },
  { name: 'Sales', total: 36, onTime: 24 },
  { name: 'Finance', total: 24, onTime: 21 },
  { name: 'Marketing', total: 31, onTime: 27 },
];

const REPORT_TYPES = [
  { value: 'daily', label: 'Báo cáo ngày' },
  { value: 'weekly', label: 'Báo cáo tuần' },
  { value: 'monthly', label: 'Báo cáo tháng' },
  { value: 'custom', label: 'Tùy chỉnh' },
];

export default function Reports() {
  const { showToast } = useUi();
  const [reportType, setReportType] = useState('daily');
  const [dateRange, setDateRange] = useState({
    start: new Date().toISOString().slice(0, 10),
    end: new Date().toISOString().slice(0, 10),
  });
  const [includeLate, setIncludeLate] = useState(true);
  const [includeAbsent, setIncludeAbsent] = useState(true);

  const summary = useMemo(() => {
    const totalEmployees = MOCK_DEPARTMENTS.reduce((acc, dept) => acc + dept.total, 0);
    const onTime = MOCK_DEPARTMENTS.reduce((acc, dept) => acc + dept.onTime, 0);
    const late = includeLate ? Math.max(totalEmployees - onTime - 15, 0) : 0;
    const absent = includeAbsent ? 15 : 0;

    return {
      totalEmployees,
      onTime,
      late,
      absent,
    };
  }, [includeAbsent, includeLate]);

  const handleExport = (format: 'excel' | 'pdf') => {
    showToast({
      variant: 'success',
      message: `Đã bắt đầu xuất báo cáo (${format.toUpperCase()}) cho giai đoạn ${new Date(
        dateRange.start
      ).toLocaleDateString('vi-VN')} - ${new Date(dateRange.end).toLocaleDateString('vi-VN')}.`,
    });
  };

  return (
    <Page
      title="Báo cáo"
      subtitle="Tổng hợp tình hình chấm công và xuất báo cáo linh hoạt theo nhu cầu"
      breadcrumb={[{ label: 'Trang chủ', path: '/dashboard' }, { label: 'Báo cáo' }]}
      actions={
        <Stack direction="horizontal" gap={2}>
          <Button variant="outline-primary" onClick={() => handleExport('excel')}>
            Xuất Excel
          </Button>
          <Button variant="primary" onClick={() => handleExport('pdf')}>
            Xuất PDF
          </Button>
        </Stack>
      }
    >
      <Card className="border-0 shadow-sm">
        <Card.Body>
          <Row className="g-3 align-items-end">
            <Col lg={3}>
              <Form.Group controlId="report-type">
                <Form.Label className="text-secondary text-uppercase small">Loại báo cáo</Form.Label>
                <Form.Select value={reportType} onChange={(event) => setReportType(event.target.value)}>
                  {REPORT_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col lg={3}>
              <Form.Group controlId="report-start">
                <Form.Label className="text-secondary text-uppercase small">Từ ngày</Form.Label>
                <Form.Control
                  type="date"
                  value={dateRange.start}
                  max={dateRange.end}
                  onChange={(event) => setDateRange((prev) => ({ ...prev, start: event.target.value }))}
                />
              </Form.Group>
            </Col>
            <Col lg={3}>
              <Form.Group controlId="report-end">
                <Form.Label className="text-secondary text-uppercase small">Đến ngày</Form.Label>
                <Form.Control
                  type="date"
                  value={dateRange.end}
                  min={dateRange.start}
                  onChange={(event) => setDateRange((prev) => ({ ...prev, end: event.target.value }))}
                />
              </Form.Group>
            </Col>
            <Col lg={3}>
              <Stack direction="horizontal" className="gap-2 justify-content-end">
                <Form.Check
                  id="include-late"
                  type="switch"
                  label="Bao gồm đi trễ"
                  checked={includeLate}
                  onChange={(event) => setIncludeLate(event.target.checked)}
                />
                <Form.Check
                  id="include-absent"
                  type="switch"
                  label="Bao gồm vắng"
                  checked={includeAbsent}
                  onChange={(event) => setIncludeAbsent(event.target.checked)}
                />
              </Stack>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      <Row className="g-3">
        <Col xl={6}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Header className="bg-transparent border-0 pb-0">
              <Stack direction="horizontal" className="justify-content-between align-items-center">
                <div>
                  <h2 className="fs-5 fw-semibold mb-0">Tổng quan chấm công</h2>
                  <p className="text-secondary small mb-0">Thời gian {new Date(dateRange.start).toLocaleDateString('vi-VN')} - {new Date(dateRange.end).toLocaleDateString('vi-VN')}</p>
                </div>
              </Stack>
            </Card.Header>
            <Card.Body>
              <Row className="g-3">
                <Col sm={6}>
                  <SummaryBadge label="Tổng nhân viên" value={summary.totalEmployees} variant="primary" />
                </Col>
                <Col sm={6}>
                  <SummaryBadge label="Đi làm đúng giờ" value={summary.onTime} variant="success" />
                </Col>
                <Col sm={6}>
                  <SummaryBadge label="Đi trễ" value={summary.late} variant="warning" />
                </Col>
                <Col sm={6}>
                  <SummaryBadge label="Vắng mặt" value={summary.absent} variant="danger" />
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
        <Col xl={6}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Header className="bg-transparent border-0 pb-0">
              <h2 className="fs-5 fw-semibold mb-0">Hiệu suất theo phòng ban</h2>
              <p className="text-secondary small mb-0">Tỉ lệ vào đúng giờ</p>
            </Card.Header>
            <Card.Body>
              <Stack gap={3}>
                {MOCK_DEPARTMENTS.map((dept) => {
                  const percent = Math.round((dept.onTime / dept.total) * 100);
                  return (
                    <div key={dept.name}>
                      <div className="d-flex justify-content-between">
                        <span className="fw-semibold">{dept.name}</span>
                        <span className="text-secondary small">{dept.onTime}/{dept.total} đúng giờ</span>
                      </div>
                      <ProgressBar now={percent} label={`${percent}%`} className="mt-2" visuallyHidden={false} />
                    </div>
                  );
                })}
              </Stack>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Card className="border-0 shadow-sm">
        <Card.Header className="bg-transparent border-0 pb-0">
          <Stack direction="horizontal" className="justify-content-between align-items-center">
            <h2 className="fs-5 fw-semibold mb-0">Bảng chi tiết dữ liệu</h2>
            <Button variant="outline-secondary" size="sm" onClick={() => handleExport('excel')}>
              Tải file mẫu
            </Button>
          </Stack>
        </Card.Header>
        <Card.Body className="pt-3">
          <div className="table-responsive">
            <Table striped hover className="align-middle mb-0">
              <thead>
                <tr>
                  <th>Phòng ban</th>
                  <th>Nhân viên</th>
                  <th>Đúng giờ</th>
                  {includeLate && <th>Đi trễ</th>}
                  {includeAbsent && <th>Vắng mặt</th>}
                </tr>
              </thead>
              <tbody>
                {MOCK_DEPARTMENTS.map((dept) => {
                  const late = Math.max(dept.total - dept.onTime - 3, 0);
                  const absent = Math.max(3, 0);
                  return (
                    <tr key={dept.name}>
                      <td className="fw-semibold">{dept.name}</td>
                      <td>{dept.total}</td>
                      <td className="text-success fw-semibold">{dept.onTime}</td>
                      {includeLate && <td className="text-warning">{late}</td>}
                      {includeAbsent && <td className="text-danger">{absent}</td>}
                    </tr>
                  );
                })}
              </tbody>
            </Table>
          </div>
        </Card.Body>
      </Card>
    </Page>
  );
}

interface SummaryBadgeProps {
  label: string;
  value: number;
  variant: 'primary' | 'success' | 'warning' | 'danger';
}

function SummaryBadge({ label, value, variant }: SummaryBadgeProps) {
  const variantMap: Record<SummaryBadgeProps['variant'], string> = {
    primary: 'bg-primary-subtle text-primary',
    success: 'bg-success-subtle text-success',
    warning: 'bg-warning-subtle text-warning',
    danger: 'bg-danger-subtle text-danger',
  };

  return (
    <div className="border rounded-3 p-3 h-100">
      <p className="text-secondary text-uppercase small mb-1">{label}</p>
      <h3 className="fw-semibold mb-0">{value}</h3>
      <span className={`badge ${variantMap[variant]} mt-2`}>Cập nhật mới nhất</span>
    </div>
  );
}