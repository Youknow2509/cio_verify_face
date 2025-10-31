// src/features/shifts/Shifts.tsx

import { useMemo, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';

interface ShiftItem {
  id: string;
  name: string;
  startTime: string;
  endTime: string;
  workHours: number;
  breakTime: number;
  isActive: boolean;
  description?: string;
}

const INITIAL_SHIFTS: ShiftItem[] = [
  {
    id: '1',
    name: 'Ca sáng',
    startTime: '08:00',
    endTime: '12:00',
    workHours: 4,
    breakTime: 0,
    isActive: true,
    description: 'Áp dụng cho bộ phận văn phòng, hỗ trợ linh hoạt 15 phút.',
  },
  {
    id: '2',
    name: 'Ca chiều',
    startTime: '13:00',
    endTime: '17:00',
    workHours: 4,
    breakTime: 0,
    isActive: true,
    description: 'Phù hợp cho bộ phận chăm sóc khách hàng.',
  },
  {
    id: '3',
    name: 'Ca hành chính',
    startTime: '08:00',
    endTime: '17:00',
    workHours: 8,
    breakTime: 1,
    isActive: true,
    description: 'Ca làm việc tiêu chuẩn, nghỉ trưa 60 phút.',
  },
  {
    id: '4',
    name: 'Ca tối',
    startTime: '18:00',
    endTime: '22:00',
    workHours: 4,
    breakTime: 0,
    isActive: false,
    description: 'Ca dự phòng cho sự kiện đặc biệt buổi tối.',
  },
];

export default function Shifts() {
  const { showToast, confirm } = useUi();
  const [searchQuery, setSearchQuery] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [shifts, setShifts] = useState<ShiftItem[]>(INITIAL_SHIFTS);
  const [formState, setFormState] = useState({
    name: '',
    startTime: '08:00',
    endTime: '17:00',
    workHours: 8,
    breakTime: 1,
    isActive: true,
    description: '',
  });

  const filteredShifts = useMemo(() => {
    return shifts.filter((shift) =>
      shift.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      shift.description?.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [shifts, searchQuery]);

  const activeStats = useMemo(() => {
    const active = shifts.filter((shift) => shift.isActive).length;
    return {
      total: shifts.length,
      active,
      inactive: shifts.length - active,
    };
  }, [shifts]);

  const handleCreateShift = () => {
    if (!formState.name.trim()) {
      showToast({ variant: 'warning', message: 'Vui lòng nhập tên ca làm việc.' });
      return;
    }

    const newShift: ShiftItem = {
      id: String(Date.now()),
      ...formState,
    };

    setShifts((prev) => [newShift, ...prev]);
    setShowCreateModal(false);
    setFormState({ name: '', startTime: '08:00', endTime: '17:00', workHours: 8, breakTime: 1, isActive: true, description: '' });
    showToast({ variant: 'success', message: 'Đã tạo ca làm việc mới.' });
  };

  const handleDeleteShift = async (shift: ShiftItem) => {
    const shouldDelete = await confirm({
      title: 'Xóa ca làm việc',
      message: `Bạn có chắc muốn xóa ca "${shift.name}"? Thao tác này không thể hoàn tác.`,
      confirmLabel: 'Xóa ca',
      cancelLabel: 'Hủy',
      confirmVariant: 'danger',
    });

    if (!shouldDelete) {
      return;
    }

    setShifts((prev) => prev.filter((item) => item.id !== shift.id));
    showToast({ variant: 'success', message: 'Đã xóa ca làm việc.' });
  };

  const handleToggleStatus = (shift: ShiftItem) => {
    setShifts((prev) =>
      prev.map((item) =>
        item.id === shift.id
          ? { ...item, isActive: !item.isActive }
          : item
      )
    );
    showToast({
      variant: 'info',
      message: `${shift.name} đã được ${shift.isActive ? 'tạm ngưng' : 'kích hoạt'} lại.`,
    });
  };

  return (
    <Page
      title="Quản lý ca làm việc"
      subtitle="Thiết lập ca trực rõ ràng giúp tối ưu lịch làm việc và tính công chính xác"
      breadcrumb={[{ label: 'Trang chủ', path: '/dashboard' }, { label: 'Ca làm việc' }]}
      actions={
        <Button variant="primary" onClick={() => setShowCreateModal(true)}>
          + Thêm ca làm việc
        </Button>
      }
    >
      <Row className="g-3">
        <Col md={4}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <p className="text-uppercase text-secondary small mb-2">Tổng quan</p>
              <h2 className="fs-2 fw-semibold mb-3">{activeStats.total} ca</h2>
              <div className="d-flex flex-column gap-2">
                <div className="d-flex justify-content-between align-items-center">
                  <span className="text-secondary">Đang áp dụng</span>
                  <Badge bg="success" pill>
                    {activeStats.active}
                  </Badge>
                </div>
                <div className="d-flex justify-content-between align-items-center">
                  <span className="text-secondary">Tạm ngưng</span>
                  <Badge bg="secondary" pill>
                    {activeStats.inactive}
                  </Badge>
                </div>
              </div>
            </Card.Body>
          </Card>
        </Col>
        <Col md={8}>
          <Card className="border-0 shadow-sm h-100">
            <Card.Body>
              <Form.Group controlId="shift-search">
                <Form.Label className="text-secondary text-uppercase small">Tìm kiếm</Form.Label>
                <Form.Control
                  type="search"
                  placeholder="Tên ca, mô tả..."
                  value={searchQuery}
                  onChange={(event) => setSearchQuery(event.target.value)}
                />
              </Form.Group>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row className="g-3">
        {filteredShifts.map((shift) => (
          <Col key={shift.id} md={6} xl={4}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Body className="d-flex flex-column gap-3">
                <Stack direction="horizontal" className="justify-content-between align-items-start gap-3">
                  <div>
                    <h3 className="fs-5 fw-semibold mb-1">{shift.name}</h3>
                    <p className="text-secondary small mb-0">
                      {shift.startTime} - {shift.endTime} · {shift.workHours}h làm việc
                    </p>
                  </div>
                  <Badge bg={shift.isActive ? 'success' : 'secondary'}>
                    {shift.isActive ? 'Đang dùng' : 'Tạm ngưng'}
                  </Badge>
                </Stack>

                {shift.description && <p className="text-secondary small mb-0">{shift.description}</p>}

                <div className="bg-light border rounded p-3">
                  <div className="d-flex justify-content-between small mb-2">
                    <span className="text-secondary">Thời gian làm</span>
                    <span className="fw-semibold">{shift.workHours} giờ</span>
                  </div>
                  <div className="d-flex justify-content-between small">
                    <span className="text-secondary">Nghỉ giữa ca</span>
                    <span className="fw-semibold">{shift.breakTime} giờ</span>
                  </div>
                </div>

                <Stack direction="horizontal" gap={2} className="mt-auto flex-wrap">
                  <Button variant="outline-secondary" size="sm" onClick={() => handleToggleStatus(shift)}>
                    {shift.isActive ? 'Tạm ngưng' : 'Kích hoạt'}
                  </Button>
                  <Button variant="outline-danger" size="sm" onClick={() => handleDeleteShift(shift)}>
                    Xóa ca
                  </Button>
                </Stack>
              </Card.Body>
            </Card>
          </Col>
        ))}
        {filteredShifts.length === 0 && (
          <Col>
            <Card className="border-0 shadow-sm">
              <Card.Body className="text-center text-secondary py-5">
                Không tìm thấy ca làm việc phù hợp với từ khóa "{searchQuery}".
              </Card.Body>
            </Card>
          </Col>
        )}
      </Row>

      <Modal show={showCreateModal} onHide={() => setShowCreateModal(false)} centered>
        <Modal.Header closeButton>
          <Modal.Title>Thêm ca làm việc</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Stack gap={3}>
            <Form.Group controlId="shift-name">
              <Form.Label>Tên ca *</Form.Label>
              <Form.Control
                value={formState.name}
                onChange={(event) => setFormState((prev) => ({ ...prev, name: event.target.value }))}
                placeholder="Ví dụ: Ca sáng linh hoạt"
              />
            </Form.Group>
            <Row className="g-3">
              <Col xs={6}>
                <Form.Group controlId="shift-start">
                  <Form.Label>Bắt đầu</Form.Label>
                  <Form.Control
                    type="time"
                    value={formState.startTime}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, startTime: event.target.value }))
                    }
                  />
                </Form.Group>
              </Col>
              <Col xs={6}>
                <Form.Group controlId="shift-end">
                  <Form.Label>Kết thúc</Form.Label>
                  <Form.Control
                    type="time"
                    value={formState.endTime}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, endTime: event.target.value }))
                    }
                  />
                </Form.Group>
              </Col>
            </Row>
            <Row className="g-3">
              <Col xs={6}>
                <Form.Group controlId="shift-workHours">
                  <Form.Label>Giờ làm</Form.Label>
                  <Form.Control
                    type="number"
                    min={1}
                    max={12}
                    value={formState.workHours}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, workHours: Number(event.target.value) }))
                    }
                  />
                </Form.Group>
              </Col>
              <Col xs={6}>
                <Form.Group controlId="shift-break">
                  <Form.Label>Nghỉ giữa ca</Form.Label>
                  <Form.Control
                    type="number"
                    min={0}
                    max={4}
                    value={formState.breakTime}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, breakTime: Number(event.target.value) }))
                    }
                  />
                </Form.Group>
              </Col>
            </Row>
            <Form.Group controlId="shift-description">
              <Form.Label>Mô tả</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                value={formState.description}
                onChange={(event) =>
                  setFormState((prev) => ({ ...prev, description: event.target.value }))
                }
              />
            </Form.Group>
            <Form.Check
              type="switch"
              id="shift-status"
              label="Kích hoạt ngay"
              checked={formState.isActive}
              onChange={(event) =>
                setFormState((prev) => ({ ...prev, isActive: event.target.checked }))
              }
            />
          </Stack>
        </Modal.Body>
        <Modal.Footer className="d-flex justify-content-between">
          <Button variant="link" onClick={() => setShowCreateModal(false)}>
            Hủy
          </Button>
          <Button variant="primary" onClick={handleCreateShift}>
            Lưu ca làm việc
          </Button>
        </Modal.Footer>
      </Modal>
    </Page>
  );
}