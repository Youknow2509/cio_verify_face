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
import { Users, Clock, AlertTriangle, UserCheck } from 'lucide-react';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { DateRangePicker } from '@/components/DateRangePicker/DateRangePicker';
import { FilterBar, FilterGroup, FilterSelect } from '@/components/FilterBar/FilterBar';
import { SummaryCard } from '@/components/SummaryCard/SummaryCard';

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
    start: new Date(),
    end: new Date(),
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
      message: `Đã bắt đầu xuất báo cáo (${format.toUpperCase()}) cho giai đoạn ${dateRange.start.toLocaleDateString('vi-VN')} - ${dateRange.end.toLocaleDateString('vi-VN')}.`,
    });
  };

  const reportTypeOptions = REPORT_TYPES.map(type => ({
    value: type.value,
    label: type.label
  }));

  const hasActiveFilters = reportType !== 'daily' || includeLate !== true || includeAbsent !== true;

  const handleClearFilters = () => {
    setReportType('daily');
    setDateRange({
      start: new Date(),
      end: new Date(),
    });
    setIncludeLate(true);
    setIncludeAbsent(true);
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
      <FilterBar 
        hasActiveFilters={hasActiveFilters}
        onClear={handleClearFilters}
      >
        <FilterGroup label="Loại báo cáo">
          <FilterSelect
            value={reportType}
            onChange={setReportType}
            options={reportTypeOptions}
            placeholder="Chọn loại báo cáo"
          />
        </FilterGroup>

        <FilterGroup label="Khoảng thời gian">
          <DateRangePicker
            startDate={dateRange.start}
            endDate={dateRange.end}
            onChange={(start, end) => setDateRange({ start: start || new Date(), end: end || new Date() })}
            startPlaceholder="Chọn ngày bắt đầu"
            endPlaceholder="Chọn ngày kết thúc"
          />
        </FilterGroup>

        <FilterGroup label="Tùy chọn">
          <Stack direction="horizontal" className="gap-3">
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
        </FilterGroup>
      </FilterBar>

      <Row className="g-3">
        <Col xl={6}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Header className="bg-transparent border-0 pb-0">
              <Stack direction="horizontal" className="justify-content-between align-items-center">
                <div>
                  <h2 className="fs-5 fw-semibold mb-0">Tổng quan chấm công</h2>
                  <p className="text-secondary small mb-0">Thời gian {dateRange.start.toLocaleDateString('vi-VN')} - {dateRange.end.toLocaleDateString('vi-VN')}</p>
                </div>
              </Stack>
            </Card.Header>
            <Card.Body>
              <Row className="g-3">
                <Col sm={6}>
                  <SummaryCard
                    title="Tổng nhân viên"
                    value={summary.totalEmployees}
                    icon={<Users size={24} />}
                    variant="primary"
                    subtitle="Tổng số trong hệ thống"
                  />
                </Col>
                <Col sm={6}>
                  <SummaryCard
                    title="Đi làm đúng giờ"
                    value={summary.onTime}
                    icon={<UserCheck size={24} />}
                    variant="success"
                    subtitle="Chấm công đúng giờ"
                    trend={{
                      type: 'up',
                      value: 5.2,
                      text: 'so với tuần trước'
                    }}
                  />
                </Col>
                <Col sm={6}>
                  <SummaryCard
                    title="Đi trễ"
                    value={summary.late}
                    icon={<Clock size={24} />}
                    variant="warning"
                    subtitle="Chấm công muộn"
                    trend={{
                      type: 'down',
                      value: 2.1,
                      text: 'so với tuần trước'
                    }}
                  />
                </Col>
                <Col sm={6}>
                  <SummaryCard
                    title="Vắng mặt"
                    value={summary.absent}
                    icon={<AlertTriangle size={24} />}
                    variant="danger"
                    subtitle="Không chấm công"
                    trend={{
                      type: 'neutral',
                      value: 0,
                      text: 'không thay đổi'
                    }}
                  />
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