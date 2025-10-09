// src/features/devices/Devices.tsx

import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import Tooltip from 'react-bootstrap/Tooltip';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { DataTable, type DataTableColumn } from '@/ui/DataTable';
import { getDevices, syncDevice, deleteDevice } from '@/services';
import type { Device, FilterOptions } from '@/types';

const DEFAULT_FILTER: FilterOptions = {
  page: 1,
  limit: 10,
  search: '',
  status: '',
  sortBy: 'name',
  sortOrder: 'asc',
};

export default function Devices() {
  const { showToast, confirm } = useUi();
  const [filter, setFilter] = useState<FilterOptions>(DEFAULT_FILTER);
  const [devices, setDevices] = useState<Device[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [syncingId, setSyncingId] = useState<string | null>(null);

  const page = filter.page ?? 1;
  const pageSize = filter.limit ?? 10;

  useEffect(() => {
    const loadDevices = async () => {
      setLoading(true);
      try {
        const response = await getDevices(filter);
        setDevices(response.data);
        setTotal(response.total);
      } catch (error) {
        showToast({
          variant: 'danger',
          message:
            error instanceof Error ? error.message : 'Không thể tải danh sách thiết bị. Vui lòng thử lại.',
        });
      } finally {
        setLoading(false);
      }
    };

    loadDevices();
  }, [filter, showToast]);

  const stats = useMemo(() => {
    const online = devices.filter((device) => device.status === 'online').length;
    const offline = devices.filter((device) => device.status === 'offline').length;

    return {
      total: total,
      online,
      offline,
    };
  }, [devices, total]);

  const handleSync = async (device: Device) => {
    try {
      setSyncingId(device.id);
      await syncDevice(device.id);
      setFilter((prev) => ({ ...prev }));
      showToast({ variant: 'success', message: `Đã đồng bộ "${device.name}".` });
    } catch (error) {
      showToast({
        variant: 'danger',
        message: error instanceof Error ? error.message : 'Không thể đồng bộ thiết bị.',
      });
    } finally {
      setSyncingId(null);
    }
  };

  const handleDelete = async (device: Device) => {
    const shouldDelete = await confirm({
      title: 'Xóa thiết bị',
      message: `Thiết bị "${device.name}" sẽ bị xóa khỏi hệ thống. Bạn có chắc chắn?`,
      confirmLabel: 'Xóa thiết bị',
      cancelLabel: 'Hủy',
      confirmVariant: 'danger',
    });

    if (!shouldDelete) {
      return;
    }

    try {
      await deleteDevice(device.id);
      setFilter((prev) => ({ ...prev }));
      showToast({ variant: 'success', message: 'Đã xóa thiết bị thành công.' });
    } catch (error) {
      showToast({
        variant: 'danger',
        message: error instanceof Error ? error.message : 'Không thể xóa thiết bị.',
      });
    }
  };

  const columns: DataTableColumn<Device>[] = [
    {
      header: 'Thiết bị',
      render: (row) => (
        <div>
          <Link to={`/devices/${row.id}`} className="fw-semibold text-decoration-none">
            {row.name}
          </Link>
          {row.location && <p className="text-secondary small mb-0">{row.location}</p>}
        </div>
      ),
      className: 'text-nowrap',
    },
    {
      header: 'Model',
      render: (row) => <span className="text-secondary small">{row.model ?? '—'}</span>,
    },
    {
      header: 'Địa chỉ IP',
      render: (row) => (
        <span className="badge bg-light text-dark fw-normal">{row.ipAddress ?? '—'}</span>
      ),
      className: 'text-nowrap',
    },
    {
      header: 'Trạng thái',
      render: (row) => (
        <Badge bg={row.status === 'online' ? 'success' : 'secondary'}>
          {row.status === 'online' ? '● Hoạt động' : '● Ngoại tuyến'}
        </Badge>
      ),
      className: 'text-center',
    },
    {
      header: 'Đồng bộ lần cuối',
      render: (row) => (
        <span className="text-secondary small">
          {row.lastSyncAt
            ? new Date(row.lastSyncAt).toLocaleString('vi-VN', { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' })
            : 'Chưa đồng bộ'}
        </span>
      ),
      className: 'text-nowrap',
    },
    {
      header: 'Thao tác',
      render: (row) => (
        <Stack direction="horizontal" gap={2} className="justify-content-end">
          <Button
            variant="outline-primary"
            size="sm"
            disabled={syncingId === row.id}
            onClick={() => handleSync(row)}
          >
            {syncingId === row.id ? 'Đang đồng bộ...' : 'Đồng bộ'}
          </Button>
          <Button variant="outline-danger" size="sm" onClick={() => handleDelete(row)}>
            Xóa
          </Button>
        </Stack>
      ),
      className: 'text-end text-nowrap',
    },
  ];

  return (
    <Page
      title="Quản lý thiết bị"
      subtitle="Theo dõi tình trạng kết nối và đồng bộ của toàn bộ thiết bị điểm danh"
      breadcrumb={[{ label: 'Trang chủ', path: '/dashboard' }, { label: 'Thiết bị' }]}
      actions={
        <Stack direction="horizontal" gap={2}>
          <Button variant="outline-secondary" onClick={() => setFilter((prev) => ({ ...prev }))}>
            Làm mới
          </Button>
          <Button variant="primary">+ Thêm thiết bị</Button>
        </Stack>
      }
    >
      <Row className="g-3">
        <Col md={4}>
          <StatCard
            title="Tổng thiết bị"
            value={stats.total}
            description={`Đang xem ${devices.length} trên mỗi trang`}
            variant="primary"
          />
        </Col>
        <Col md={4}>
          <StatCard
            title="Đang hoạt động"
            value={stats.online}
            description={`${percentage(stats.online, stats.total)} trực tuyến`}
            variant="success"
          />
        </Col>
        <Col md={4}>
          <StatCard
            title="Ngoại tuyến"
            value={stats.offline}
            description={`${percentage(stats.offline, stats.total)} cần kiểm tra`}
            variant="warning"
          />
        </Col>
      </Row>

      <Card className="border-0 shadow-sm">
        <Card.Body>
          <Row className="g-3 align-items-end">
            <Col lg={4}>
              <Form.Group controlId="device-search">
                <Form.Label className="text-secondary small text-uppercase">Tìm kiếm</Form.Label>
                <Form.Control
                  type="search"
                  placeholder="Tên, vị trí, model..."
                  value={filter.search ?? ''}
                  onChange={(event) => setFilter({ ...filter, search: event.target.value, page: 1 })}
                />
              </Form.Group>
            </Col>
            <Col lg={3}>
              <Form.Group controlId="device-status">
                <Form.Label className="text-secondary small text-uppercase">Trạng thái</Form.Label>
                <Form.Select
                  value={filter.status ?? ''}
                  onChange={(event) => setFilter({ ...filter, status: event.target.value || undefined, page: 1 })}
                >
                  <option value="">Tất cả</option>
                  <option value="online">Hoạt động</option>
                  <option value="offline">Ngoại tuyến</option>
                </Form.Select>
              </Form.Group>
            </Col>
            <Col lg={3}>
              <Form.Group controlId="device-sort">
                <Form.Label className="text-secondary small text-uppercase">Sắp xếp theo</Form.Label>
                <Form.Select
                  value={filter.sortBy ?? 'name'}
                  onChange={(event) => setFilter({ ...filter, sortBy: event.target.value, page: 1 })}
                >
                  <option value="name">Tên thiết bị</option>
                  <option value="location">Vị trí</option>
                  <option value="status">Trạng thái</option>
                  <option value="lastSyncAt">Đồng bộ</option>
                </Form.Select>
              </Form.Group>
            </Col>
            <Col lg={2}>
              <Form.Group controlId="device-limit">
                <Form.Label className="text-secondary small text-uppercase">Hiển thị</Form.Label>
                <Form.Select
                  value={pageSize}
                  onChange={(event) =>
                    setFilter({ ...filter, limit: Number(event.target.value), page: 1 })
                  }
                >
                  {[5, 10, 20, 50].map((size) => (
                    <option key={size} value={size}>
                      {size} thiết bị
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      <Card className="border-0 shadow-sm">
        <Card.Header className="bg-transparent border-0 pb-0">
          <Stack direction="horizontal" className="justify-content-between align-items-center flex-wrap gap-2">
            <div>
              <h2 className="fs-5 fw-semibold mb-0">Danh sách thiết bị</h2>
              <p className="text-secondary small mb-0">
                {total} thiết bị · Cập nhật {new Date().toLocaleTimeString('vi-VN')}
              </p>
            </div>
            <OverlayTrigger
              placement="left"
              overlay={<Tooltip id="refresh-tooltip">Tải lại dữ liệu</Tooltip>}
            >
              <Button variant="outline-secondary" size="sm" onClick={() => setFilter((prev) => ({ ...prev }))}>
                Làm mới
              </Button>
            </OverlayTrigger>
          </Stack>
        </Card.Header>
        <Card.Body className="pt-3">
          <DataTable
            columns={columns}
            data={devices}
            loading={loading}
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={(nextPage) => setFilter({ ...filter, page: nextPage })}
            keySelector={(row) => row.id}
            emptyMessage="Không có thiết bị nào phù hợp với bộ lọc hiện tại"
          />
        </Card.Body>
      </Card>
    </Page>
  );
}

interface StatCardProps {
  title: string;
  value: number;
  description: string;
  variant: 'primary' | 'success' | 'warning';
}

function StatCard({ title, value, description, variant }: StatCardProps) {
  const variantMap: Record<StatCardProps['variant'], string> = {
    primary: 'bg-primary-subtle text-primary',
    success: 'bg-success-subtle text-success',
    warning: 'bg-warning-subtle text-warning',
  };

  return (
    <Card className="border-0 shadow-sm h-100">
      <Card.Body className="d-flex flex-column gap-2">
        <span className="text-secondary text-uppercase small">{title}</span>
        <h2 className="fs-2 fw-semibold mb-0">{value}</h2>
        <span className={`badge ${variantMap[variant]} align-self-start`}>{description}</span>
      </Card.Body>
    </Card>
  );
}

function percentage(value: number, total: number) {
  if (total === 0) {
    return '0%';
  }

  return `${Math.round((value / total) * 100)}%`;
}