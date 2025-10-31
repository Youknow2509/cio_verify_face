import { useCallback, useEffect, useMemo, useState, type FormEvent } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import FloatingLabel from 'react-bootstrap/FloatingLabel';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import { useNavigate } from 'react-router-dom';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { DataTable, type DataTableColumn } from '@/ui/DataTable';
import { FilterBar, SearchBox, FilterGroup, FilterSelect } from '@/components/FilterBar/FilterBar';
import { createEmployee, deleteEmployee, getEmployees } from '@/services';
import type { Employee, EmployeeFilter, FilterOptions } from '@/types';

interface EmployeeFormState {
  code: string;
  name: string;
  email: string;
  department: string;
  position: string;
  active: boolean;
}

type FormErrors = Partial<Record<keyof EmployeeFormState, string>>;

const defaultFilter: EmployeeFilter = {
  page: 1,
  limit: 10,
  search: '',
  department: '',
  active: undefined,
  sortBy: 'name',
  sortOrder: 'asc',
};

const defaultFormState: EmployeeFormState = {
  code: '',
  name: '',
  email: '',
  department: '',
  position: '',
  active: true,
};

export default function EmployeesPage() {
  const navigate = useNavigate();
  const [filter, setFilter] = useState<EmployeeFilter>(defaultFilter);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [departments, setDepartments] = useState<string[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [formState, setFormState] = useState<EmployeeFormState>(defaultFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { showToast, confirm } = useUi();

  // Helper functions for filter
  const departmentOptions = useMemo(() => 
    departments.map(dept => ({ value: dept, label: dept })), 
    [departments]
  );

  const statusOptions = [
    { value: 'true', label: 'Hoạt động' },
    { value: 'false', label: 'Tạm dừng' }
  ];

  const hasActiveFilters = useMemo(() => 
    filter.search !== '' || filter.department !== '' || filter.active !== undefined,
    [filter.search, filter.department, filter.active]
  );

  const handleClearFilters = useCallback(() => {
    setFilter(prev => ({
      ...prev,
      search: '',
      department: '',
      active: undefined,
      page: 1
    }));
  }, []);

  const normalizedFilters = useMemo<FilterOptions>(() => {
    const status =
      typeof filter.active === 'boolean' ? (filter.active ? 'active' : 'inactive') : undefined;

    return {
      search: filter.search,
      department: filter.department,
      status,
      page: filter.page,
      limit: filter.limit,
      sortBy: filter.sortBy,
      sortOrder: filter.sortOrder,
    };
  }, [filter]);

  const loadEmployees = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getEmployees(normalizedFilters);
      setEmployees(response.data);
      setTotal(response.total);
      const uniqueDepartments = Array.from(
        new Set(response.data.map((item) => item.department).filter(Boolean) as string[])
      ).sort((a, b) => a.localeCompare(b));
      setDepartments(uniqueDepartments);
    } catch (error) {
      console.error(error);
      showToast({
        variant: 'danger',
        title: 'Lỗi',
        message: 'Không thể tải danh sách nhân viên',
      });
    } finally {
      setLoading(false);
    }
  }, [normalizedFilters, showToast]);

  useEffect(() => {
    void loadEmployees();
  }, [loadEmployees]);

  const handleChangeFilter = <K extends keyof EmployeeFilter>(key: K, value: EmployeeFilter[K]) => {
    setFilter((prev) => {
      const next: EmployeeFilter = {
        ...prev,
        [key]: value,
      };
      if (key !== 'page') {
        next.page = 1;
      }
      return next;
    });
  };

  const validateForm = useCallback(() => {
    const errors: FormErrors = {};

    if (!/^[A-Z0-9]+$/.test(formState.code.trim())) {
      errors.code = 'Mã nhân viên chỉ gồm chữ in hoa và số';
    }

    if (formState.name.trim().length < 2) {
      errors.name = 'Tên nhân viên tối thiểu 2 ký tự';
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formState.email.trim())) {
      errors.email = 'Email không hợp lệ';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  }, [formState]);

  const handleSubmit = useCallback(
    async (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      if (!validateForm()) return;

      try {
        const payload = {
          code: formState.code,
          name: formState.name,
          email: formState.email,
          department: formState.department,
          position: formState.position,
          active: formState.active,
        };
        const response = await createEmployee(payload);
        if (response.error) {
          throw new Error(response.error);
        }
        showToast({
          variant: 'success',
          title: 'Thành công',
          message: 'Đã thêm nhân viên mới',
        });
        setShowModal(false);
        setFormState(defaultFormState);
        await loadEmployees();
      } catch (error) {
        console.error(error);
        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể thêm nhân viên',
        });
      }
    },
    [formState, loadEmployees, showToast, validateForm]
  );

  const handleDelete = useCallback(
    async (employee: Employee) => {
      const confirmed = await confirm({
        title: 'Xóa nhân viên',
        message: `Bạn có chắc chắn muốn xóa nhân viên "${employee.name}"?`,
        confirmLabel: 'Xóa',
        cancelLabel: 'Hủy',
        confirmVariant: 'danger',
      });

      if (!confirmed) return;

      try {
        const response = await deleteEmployee(employee.id);
        if (response.error) {
          throw new Error(response.error);
        }
        showToast({
          variant: 'success',
          title: 'Thành công',
          message: 'Đã xóa nhân viên',
        });
        await loadEmployees();
      } catch (error) {
        console.error(error);
        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể xóa nhân viên',
        });
      }
    },
    [confirm, loadEmployees, showToast]
  );

  const columns = useMemo<DataTableColumn<Employee>[]>(
    () => [
      {
        header: 'Mã NV',
        accessor: 'code',
        className: 'fw-semibold text-nowrap',
      },
      {
        header: 'Tên nhân viên',
        render: (employee) => (
          <Stack direction="horizontal" gap={3}>
            <i className="bi bi-person-circle fs-4 text-primary" aria-hidden />
            <div>
              <div className="fw-semibold">{employee.name}</div>
              <small className="text-secondary">{employee.email}</small>
            </div>
          </Stack>
        ),
      },
      {
        header: 'Phòng ban',
        accessor: 'department',
      },
      {
        header: 'Chức vụ',
        accessor: 'position',
      },
      {
        header: 'Khuôn mặt',
        render: (employee) => (
          <Badge bg={employee.faceCount > 0 ? 'success' : 'secondary'}>{employee.faceCount} ảnh</Badge>
        ),
      },
      {
        header: 'Trạng thái',
        render: (employee) => (
          <Badge bg={employee.active ? 'success' : 'danger'}>
            {employee.active ? 'Hoạt động' : 'Tạm dừng'}
          </Badge>
        ),
      },
      {
        header: 'Thao tác',
        className: 'text-end text-nowrap',
        render: (employee) => (
          <Stack direction="horizontal" gap={2} className="justify-content-end">
            <Button
              size="sm"
              variant="outline-primary"
              onClick={() => navigate(`/employees/${employee.id}`)}
              aria-label={`Chỉnh sửa ${employee.name}`}
            >
              <i className="bi bi-pencil" aria-hidden />
            </Button>
            <Button
              size="sm"
              variant="outline-danger"
              onClick={() => {
                void handleDelete(employee);
              }}
              aria-label={`Xóa ${employee.name}`}
            >
              <i className="bi bi-trash" aria-hidden />
            </Button>
          </Stack>
        ),
      },
    ],
    [handleDelete]
  );

  const actions = (
    <>
      <Button variant="outline-secondary">
        <i className="bi bi-download me-2" aria-hidden />
        Xuất Excel
      </Button>
      <Button variant="primary" onClick={() => setShowModal(true)}>
        <i className="bi bi-plus-lg me-2" aria-hidden />
        Thêm nhân viên
      </Button>
    </>
  );

  return (
    <Page
      title="Quản lý nhân viên"
      subtitle="Quản lý thông tin và dữ liệu khuôn mặt của nhân viên"
      breadcrumb={[
        { label: 'Trang chủ', path: '/dashboard' },
        { label: 'Nhân viên' },
      ]}
      actions={actions}
    >
      <FilterBar 
        hasActiveFilters={hasActiveFilters}
        onClear={handleClearFilters}
      >
        <SearchBox
          value={filter.search ?? ''}
          onChange={(value) => handleChangeFilter('search', value)}
          placeholder="Tìm theo tên, mã nhân viên..."
        />
        
        <FilterGroup label="Phòng ban">
          <FilterSelect
            value={filter.department ?? ''}
            onChange={(value) => handleChangeFilter('department', value)}
            options={departmentOptions}
            placeholder="Tất cả phòng ban"
          />
        </FilterGroup>

        <FilterGroup label="Trạng thái">
          <FilterSelect
            value={typeof filter.active === 'boolean' ? String(filter.active) : ''}
            onChange={(value) => {
              handleChangeFilter('active', value === '' ? undefined : value === 'true');
            }}
            options={statusOptions}
            placeholder="Tất cả trạng thái"
          />
        </FilterGroup>
      </FilterBar>

      <div className="mt-4">
        <DataTable
          columns={columns}
          data={employees}
          loading={loading}
          page={filter.page ?? 1}
          pageSize={filter.limit ?? 10}
          total={total}
          onPageChange={(pageNumber) => handleChangeFilter('page', pageNumber)}
          keySelector={(employee) => employee.id}
        />
      </div>

      <Modal show={showModal} onHide={() => setShowModal(false)} backdrop="static" size="lg">
        <Form onSubmit={handleSubmit} noValidate>
          <Modal.Header closeButton>
            <Modal.Title className="fw-semibold">Thêm nhân viên</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Row className="g-3">
              <Col md={6}>
                <FloatingLabel label="Mã nhân viên">
                  <Form.Control
                    value={formState.code}
                    onChange={(event) =>
                      setFormState((prev) => ({
                        ...prev,
                        code: event.target.value.toUpperCase(),
                      }))
                    }
                    isInvalid={Boolean(formErrors.code)}
                    required
                  />
                  <Form.Control.Feedback type="invalid">
                    {formErrors.code}
                  </Form.Control.Feedback>
                </FloatingLabel>
              </Col>
              <Col md={6}>
                <FloatingLabel label="Họ và tên">
                  <Form.Control
                    value={formState.name}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, name: event.target.value }))
                    }
                    isInvalid={Boolean(formErrors.name)}
                    required
                  />
                  <Form.Control.Feedback type="invalid">
                    {formErrors.name}
                  </Form.Control.Feedback>
                </FloatingLabel>
              </Col>
              <Col md={6}>
                <FloatingLabel label="Email">
                  <Form.Control
                    type="email"
                    value={formState.email}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, email: event.target.value }))
                    }
                    isInvalid={Boolean(formErrors.email)}
                    required
                  />
                  <Form.Control.Feedback type="invalid">
                    {formErrors.email}
                  </Form.Control.Feedback>
                </FloatingLabel>
              </Col>
              <Col md={6}>
                <FloatingLabel label="Phòng ban">
                  <Form.Control
                    value={formState.department}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, department: event.target.value }))
                    }
                  />
                </FloatingLabel>
              </Col>
              <Col md={6}>
                <FloatingLabel label="Chức vụ">
                  <Form.Control
                    value={formState.position}
                    onChange={(event) =>
                      setFormState((prev) => ({ ...prev, position: event.target.value }))
                    }
                  />
                </FloatingLabel>
              </Col>
              <Col md={6} className="d-flex align-items-center">
                <Form.Check
                  type="switch"
                  id="employee-active"
                  label="Hoạt động"
                  checked={formState.active}
                  onChange={(event) =>
                    setFormState((prev) => ({ ...prev, active: event.target.checked }))
                  }
                />
              </Col>
            </Row>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="outline-secondary" onClick={() => setShowModal(false)}>
              Hủy
            </Button>
            <Button variant="primary" type="submit">
              Lưu nhân viên
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Page>
  );
}
