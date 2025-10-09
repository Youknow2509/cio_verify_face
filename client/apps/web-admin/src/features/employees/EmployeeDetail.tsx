import { useCallback, useEffect, useMemo, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Spinner from 'react-bootstrap/Spinner';
import Stack from 'react-bootstrap/Stack';
import { useNavigate, useParams } from 'react-router-dom';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { fetchEmployeeById } from '@/services/api/employees';
import { clearAuthToken, HttpError } from '@/services/http';
import { formatDateTime } from '@/lib/format';
import type { Employee } from '@/types';

export default function EmployeeDetail() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { showToast } = useUi();
  const [employee, setEmployee] = useState<Employee | null>(null);
  const [loading, setLoading] = useState(true);
  const [reloading, setReloading] = useState(false);

  const loadEmployee = useCallback(async (employeeId: string) => {
    try {
      setLoading(true);
      const result = await fetchEmployeeById(employeeId);
      setEmployee(result);
    } catch (error) {
      if (error instanceof HttpError) {
        if (error.status === 401) {
          clearAuthToken();
          navigate('/login', { replace: true });
          return;
        }

        if (error.status === 404) {
          showToast({
            variant: 'warning',
            title: 'Không tìm thấy',
            message: 'Nhân viên không tồn tại hoặc đã bị xóa',
          });
          navigate('/employees', { replace: true });
          return;
        }
      }

      showToast({
        variant: 'danger',
        title: 'Lỗi',
        message: 'Không thể tải chi tiết nhân viên',
      });
    } finally {
      setLoading(false);
      setReloading(false);
    }
  }, [navigate, showToast]);

  useEffect(() => {
    if (!id) {
      navigate('/employees', { replace: true });
      return;
    }

    void loadEmployee(id);
  }, [id, loadEmployee, navigate]);

  const breadcrumb = useMemo(() => (
    [
      { label: 'Trang chủ', path: '/dashboard' },
      { label: 'Nhân viên', path: '/employees' },
      { label: employee ? employee.name : 'Chi tiết' },
    ]
  ), [employee]);

  const handleRefresh = useCallback(() => {
    if (!id) return;
    setReloading(true);
    void loadEmployee(id);
  }, [id, loadEmployee]);

  const pageActions = useMemo(() => (
    <Stack direction="horizontal" gap={2}>
      <Button variant="outline-secondary" onClick={handleRefresh} disabled={loading || reloading}>
        {reloading ? (
          <span className="d-inline-flex align-items-center gap-2">
            <Spinner animation="border" size="sm" role="status" aria-hidden />
            Đang tải
          </span>
        ) : (
          <span>
            <i className="bi bi-arrow-clockwise me-2" aria-hidden />
            Làm mới
          </span>
        )}
      </Button>
      {id && (
        <Button variant="primary" onClick={() => navigate(`/employees/${id}/edit`)}>
          <i className="bi bi-pencil-square me-2" aria-hidden />
          Chỉnh sửa
        </Button>
      )}
    </Stack>
  ), [handleRefresh, id, loading, reloading]);

  const renderContent = () => {
    if (loading) {
      return (
        <div className="d-flex align-items-center justify-content-center py-5">
          <Spinner animation="border" role="status" aria-hidden />
          <span className="ms-2">Đang tải dữ liệu nhân viên...</span>
        </div>
      );
    }

    if (!employee) {
      return (
        <div className="py-5 text-center text-secondary">
          Không có dữ liệu nhân viên để hiển thị.
        </div>
      );
    }

    return (
      <>
        <Card className="border-0 shadow-sm mb-4">
          <Card.Body>
            <Row className="g-4 align-items-center">
              <Col xs="auto">
                <div
                  className="d-flex align-items-center justify-content-center rounded-circle bg-primary text-white"
                  style={{ width: 64, height: 64 }}
                  aria-hidden
                >
                  <i className="bi bi-person-lines-fill fs-3" />
                </div>
              </Col>
              <Col>
                <Stack gap={1}>
                  <h2 className="fs-4 mb-0">{employee.name}</h2>
                  <div className="text-secondary">Mã nhân viên: <strong>{employee.code}</strong></div>
                  <div className="text-secondary">Email: <strong>{employee.email}</strong></div>
                  <div>
                    <Badge bg={employee.active ? 'success' : 'secondary'}>
                      {employee.active ? 'Đang hoạt động' : 'Tạm dừng'}
                    </Badge>
                  </div>
                </Stack>
              </Col>
              <Col xs={12} md="auto">
                <Card className="border-0 bg-body-secondary-subtle">
                  <Card.Body className="py-3 px-4 text-center">
                    <div className="text-secondary text-uppercase small">Số ảnh khuôn mặt</div>
                    <div className="fs-3 fw-semibold text-primary">{employee.faceCount}</div>
                  </Card.Body>
                </Card>
              </Col>
            </Row>
          </Card.Body>
        </Card>

        <Row className="g-4">
          <Col xl={8}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white fw-semibold">Thông tin công việc</Card.Header>
              <Card.Body>
                <Row className="g-3">
                  <Col md={6}>
                    <InfoRow label="Phòng ban" value={employee.department ?? 'Chưa cập nhật'} />
                  </Col>
                  <Col md={6}>
                    <InfoRow label="Chức vụ" value={employee.position ?? 'Chưa cập nhật'} />
                  </Col>
                  <Col md={6}>
                    <InfoRow
                      label="Ngày tạo"
                      value={formatDateTime(employee.createdAt)}
                    />
                  </Col>
                  <Col md={6}>
                    <InfoRow
                      label="Cập nhật lần cuối"
                      value={formatDateTime(employee.updatedAt)}
                    />
                  </Col>
                </Row>
              </Card.Body>
            </Card>
          </Col>
          <Col xl={4}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white fw-semibold">Tác vụ nhanh</Card.Header>
              <Card.Body>
                <Stack gap={2}>
                  <Button
                    variant="outline-primary"
                    onClick={() => {
                      if (id) {
                        navigate(`/employees/${id}/edit`);
                      } else {
                        navigate('/employees');
                      }
                    }}
                  >
                    <i className="bi bi-pencil me-2" aria-hidden />
                    Chỉnh sửa thông tin
                  </Button>
                  <Button
                    variant="outline-secondary"
                    onClick={() => navigate('/employees')}
                  >
                    <i className="bi bi-arrow-left me-2" aria-hidden />
                    Quay lại danh sách
                  </Button>
                </Stack>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </>
    );
  };

  return (
    <Page
      title={employee ? employee.name : 'Chi tiết nhân viên'}
      subtitle="Theo dõi thông tin và trạng thái hoạt động của nhân viên"
      breadcrumb={breadcrumb}
      actions={pageActions}
    >
      {renderContent()}
    </Page>
  );
}

interface InfoRowProps {
  label: string;
  value: string;
}

function InfoRow({ label, value }: InfoRowProps) {
  return (
    <Stack gap={1} className="border rounded-3 p-3 h-100 bg-body-secondary-subtle">
      <div className="text-secondary text-uppercase small">{label}</div>
      <div className="fw-semibold text-dark">{value}</div>
    </Stack>
  );
}