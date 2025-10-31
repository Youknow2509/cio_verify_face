import { useCallback, useEffect, useMemo, useState, type FormEvent } from 'react';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Spinner from 'react-bootstrap/Spinner';
import { useNavigate, useParams } from 'react-router-dom';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { getEmployee, updateEmployee } from '@/services';
import { clearAuthToken, HttpError } from '@/services/http';
import type { Employee, UpdateEmployeePayload } from '@/types';

interface FormState {
  code: string;
  name: string;
  email: string;
  department: string;
  position: string;
  active: boolean;
}

type FormErrors = Partial<Record<keyof FormState, string>>;

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function toFormState(employee: Employee): FormState {
  return {
    code: employee.code,
    name: employee.name,
    email: employee.email,
    department: employee.department ?? '',
    position: employee.position ?? '',
    active: employee.active,
  };
}

export default function EmployeeEditPage() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { showToast } = useUi();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [formState, setFormState] = useState<FormState>({
    code: '',
    name: '',
    email: '',
    department: '',
    position: '',
    active: true,
  });
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const [original, setOriginal] = useState<Employee | null>(null);

  useEffect(() => {
    if (!id) {
      navigate('/employees', { replace: true });
      return;
    }

    let isMounted = true;

    const loadEmployee = async () => {
      try {
        setLoading(true);
        const res = await getEmployee(id);
        if (!isMounted) return;
        if (res.error) {
          showToast({ variant: 'warning', title: 'Không tìm thấy', message: res.error });
          navigate('/employees', { replace: true });
          return;
        }
        const employee = res.data;
        setFormState(toFormState(employee));
        setOriginal(employee);
      } catch (error) {
        if (!isMounted) return;
        if (error instanceof HttpError && error.status === 401) {
          clearAuthToken();
          navigate('/login', { replace: true });
          return;
        }

        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể tải thông tin nhân viên',
        });
        navigate('/employees', { replace: true });
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    void loadEmployee();

    return () => {
      isMounted = false;
    };
  }, [id, navigate, showToast]);

  const isDirty = useMemo(() => {
    if (!original) return true;
    return (
      original.code !== formState.code ||
      original.name !== formState.name ||
      original.email !== formState.email ||
      (original.department ?? '') !== (formState.department ?? '') ||
      (original.position ?? '') !== (formState.position ?? '') ||
      !!original.active !== !!formState.active
    );
  }, [original, formState]);

  const validateForm = useCallback(() => {
    const errors: FormErrors = {};

    if (!/^[A-Z0-9_-]{2,}$/.test(formState.code.trim())) {
      errors.code = 'Mã nhân viên chỉ gồm chữ in hoa, số, "-" và "_"';
    }

    if (formState.name.trim().length < 2) {
      errors.name = 'Tên nhân viên tối thiểu 2 ký tự';
    }

    if (!emailRegex.test(formState.email.trim())) {
      errors.email = 'Email không hợp lệ';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  }, [formState]);

  const handleChange = <K extends keyof FormState>(key: K, value: FormState[K]) => {
    setFormState((prev) => ({
      ...prev,
      [key]: value,
    }));
  };

  const handleSubmit = useCallback(
    async (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      if (!id) {
        return;
      }

      if (!validateForm()) {
        return;
      }

      try {
        setSubmitting(true);
        const payload: UpdateEmployeePayload = {
          code: formState.code.trim().toUpperCase(),
          name: formState.name.trim(),
          email: formState.email.trim(),
          department: formState.department.trim() || undefined,
          position: formState.position.trim() || undefined,
          active: formState.active,
        };
        const res = await updateEmployee(id, payload);
        if (res.error) {
          showToast({ variant: 'danger', title: 'Lỗi', message: res.error });
          return;
        }
        showToast({
          variant: 'success',
          title: 'Thành công',
          message: 'Đã cập nhật thông tin nhân viên',
        });
        navigate(`/employees/${id}`);
      } catch (error) {
        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể cập nhật nhân viên',
        });
      } finally {
        setSubmitting(false);
      }
    },
    [formState, id, navigate, showToast, validateForm]
  );

  const pageBreadcrumb = useMemo(
    () => [
      { label: 'Trang chủ', path: '/dashboard' },
      { label: 'Nhân viên', path: '/employees' },
      id ? { label: `Chỉnh sửa #${id}` } : { label: 'Chỉnh sửa' },
    ],
    [id]
  );

  const handleCancel = useCallback(() => {
    if (id) {
      navigate(`/employees/${id}`);
    } else {
      navigate('/employees');
    }
  }, [id, navigate]);

  return (
    <Page
      title="Chỉnh sửa nhân viên"
      subtitle="Cập nhật thông tin cơ bản của nhân viên"
      breadcrumb={pageBreadcrumb}
    >
      <Card className="border-0 shadow-sm">
        <Card.Body>
          {loading ? (
            <div className="d-flex justify-content-center align-items-center py-5">
              <Spinner animation="border" role="status" aria-hidden />
              <span className="ms-2">Đang tải dữ liệu...</span>
            </div>
          ) : (
            <Form noValidate onSubmit={handleSubmit} className="gy-3">
              <Row className="g-3">
                <Col md={6}>
                  <Form.Group controlId="employee-code">
                    <Form.Label>Mã nhân viên</Form.Label>
                    <Form.Control
                      value={formState.code}
                      onChange={(event) => handleChange('code', event.target.value.toUpperCase())}
                      isInvalid={Boolean(formErrors.code)}
                      required
                    />
                    <Form.Control.Feedback type="invalid">
                      {formErrors.code}
                    </Form.Control.Feedback>
                  </Form.Group>
                </Col>
                <Col md={6}>
                  <Form.Group controlId="employee-name">
                    <Form.Label>Họ và tên</Form.Label>
                    <Form.Control
                      value={formState.name}
                      onChange={(event) => handleChange('name', event.target.value)}
                      isInvalid={Boolean(formErrors.name)}
                      required
                    />
                    <Form.Control.Feedback type="invalid">
                      {formErrors.name}
                    </Form.Control.Feedback>
                  </Form.Group>
                </Col>
                <Col md={6}>
                  <Form.Group controlId="employee-email">
                    <Form.Label>Email</Form.Label>
                    <Form.Control
                      type="email"
                      value={formState.email}
                      onChange={(event) => handleChange('email', event.target.value)}
                      isInvalid={Boolean(formErrors.email)}
                      required
                    />
                    <Form.Control.Feedback type="invalid">
                      {formErrors.email}
                    </Form.Control.Feedback>
                  </Form.Group>
                </Col>
                <Col md={6}>
                  <Form.Group controlId="employee-department">
                    <Form.Label>Phòng ban</Form.Label>
                    <Form.Control
                      value={formState.department}
                      onChange={(event) => handleChange('department', event.target.value)}
                      placeholder="Ví dụ: Kinh doanh"
                    />
                  </Form.Group>
                </Col>
                <Col md={6}>
                  <Form.Group controlId="employee-position">
                    <Form.Label>Chức vụ</Form.Label>
                    <Form.Control
                      value={formState.position}
                      onChange={(event) => handleChange('position', event.target.value)}
                      placeholder="Ví dụ: Trưởng phòng"
                    />
                  </Form.Group>
                </Col>
                <Col md={6} className="d-flex align-items-center">
                  <Form.Group controlId="employee-active" className="mb-0">
                    <Form.Check
                      type="switch"
                      label={formState.active ? 'Đang hoạt động' : 'Tạm dừng'}
                      checked={formState.active}
                      onChange={(event) => handleChange('active', event.target.checked)}
                    />
                  </Form.Group>
                </Col>
              </Row>

              <div className="d-flex justify-content-end gap-2 mt-4">
                <Button variant="outline-secondary" onClick={handleCancel} disabled={submitting}>
                  Hủy
                </Button>
                <Button type="submit" variant="primary" disabled={submitting || !isDirty}>
                  {submitting ? (
                    <span className="d-inline-flex align-items-center gap-2">
                      <Spinner animation="border" size="sm" role="status" aria-hidden />
                      Đang lưu...
                    </span>
                  ) : (
                    'Lưu thay đổi'
                  )}
                </Button>
              </div>
            </Form>
          )}
        </Card.Body>
      </Card>
    </Page>
  );
}
