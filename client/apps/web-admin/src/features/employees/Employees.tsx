// src/features/employees/Employees.tsx

import { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Plus, Edit, Trash2, Users, AlertTriangle } from 'lucide-react';
import { Card } from '@/components/Card/Card';
import { Table } from '@/components/Table/Table';
import { Badge } from '@/components/Badge/Badge';
import { FilterBar, SearchBox, FilterGroup, FilterSelect } from '@/components/FilterBar/FilterBar';
import { 
  getEmployees,
  createEmployee,
  deleteEmployee
} from '@/services';
import type { Employee, EmployeeFilter, TableColumn } from '@/types';
import styles from './Employees.module.scss';

export default function Employees() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<EmployeeFilter>({
    page: 1,
    limit: 20,
    search: '',
    department: '',
    active: undefined,
    sortBy: 'name',
    sortOrder: 'asc'
  });
  const [total, setTotal] = useState(0);
  const [showAddModal, setShowAddModal] = useState(false);

  // Departments for filter dropdown
  const departments = useMemo(() => {
    const depts = [...new Set(employees.map(emp => emp.department).filter(Boolean))];
    return depts.sort();
  }, [employees]);

  const loadEmployees = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getEmployees(filter);
      setEmployees(response.data);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load employees');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadEmployees();
  }, [filter]);

  const handleSearch = (value: string) => {
    setFilter(prev => ({ ...prev, search: value, page: 1 }));
  };

  const handleDepartmentFilter = (department: string) => {
    setFilter(prev => ({ ...prev, department, page: 1 }));
  };

  const handleActiveFilter = (active: boolean | undefined) => {
    setFilter(prev => ({ ...prev, active, page: 1 }));
  };

  const handleSort = (key: string, direction: 'asc' | 'desc') => {
    setFilter(prev => ({ ...prev, sortBy: key, sortOrder: direction }));
  };

  const handlePageChange = (page: number) => {
    setFilter(prev => ({ ...prev, page }));
  };

  const handleAddEmployee = async (data: any) => {
    try {
      await createEmployee(data);
      setShowAddModal(false);
      loadEmployees();
    } catch (err) {
      console.error('Failed to add employee:', err);
    }
  };

  const handleDeleteEmployee = async (id: string) => {
    if (!window.confirm('Bạn có chắc chắn muốn xóa nhân viên này?')) {
      return;
    }
    
    try {
      await deleteEmployee(id);
      loadEmployees();
    } catch (err) {
      console.error('Failed to delete employee:', err);
    }
  };

  const columns: TableColumn<Employee>[] = [
    {
      key: 'code',
      header: 'Mã NV',
      sortable: true,
      width: '120px',
      render: (value, record) => (
        <Link 
          to={`/employees/${record.id}`} 
          className={styles.employeeLink}
          aria-label={`View details for employee ${record.name}`}
        >
          {value}
        </Link>
      )
    },
    {
      key: 'name',
      header: 'Tên nhân viên',
      sortable: true,
      render: (value, record) => (
        <div className={styles.employeeInfo}>
          <div className={styles.avatar}>
            {value.charAt(0).toUpperCase()}
          </div>
          <div>
            <Link 
              to={`/employees/${record.id}`} 
              className={styles.employeeName}
              aria-label={`View details for employee ${record.name}`}
            >
              {value}
            </Link>
            {record.email && (
              <div className={styles.employeeEmail}>{record.email}</div>
            )}
          </div>
        </div>
      )
    },
    {
      key: 'department',
      header: 'Phòng ban',
      sortable: true,
      render: (value) => value || <span className={styles.noData}>—</span>
    },
    {
      key: 'position',
      header: 'Chức vụ',
      render: (value) => value || <span className={styles.noData}>—</span>
    },
    {
      key: 'faceCount',
      header: 'Dữ liệu khuôn mặt',
      align: 'center' as const,
      width: '140px',
      render: (value) => (
        <Badge 
          variant={value > 0 ? 'success' : 'neutral'}
          size="sm"
        >
          {value} ảnh
        </Badge>
      )
    },
    {
      key: 'active',
      header: 'Trạng thái',
      sortable: true,
      align: 'center' as const,
      width: '120px',
      render: (value) => (
        <Badge 
          variant={value ? 'success' : 'error'}
          size="sm"
        >
          {value ? 'Hoạt động' : 'Tạm dừng'}
        </Badge>
      )
    },
    {
      key: 'actions',
      header: 'Thao tác',
      align: 'center' as const,
      width: '120px',
      render: (_, record) => (
        <div className={styles.actions}>
          <Link 
            to={`/employees/${record.id}/edit`}
            className={styles.actionButton}
            aria-label={`Edit employee ${record.name}`}
          >
            <Edit size={16} />
          </Link>
          <button
            type="button"
            onClick={() => handleDeleteEmployee(record.id)}
            className={`${styles.actionButton} ${styles.danger}`}
            aria-label={`Delete employee ${record.name}`}
          >
            <Trash2 size={16} />
          </button>
        </div>
      )
    }
  ];

  const departmentOptions = departments
    .filter((dept): dept is string => dept !== undefined && dept !== null && dept !== '')
    .map(dept => ({
      value: dept,
      label: dept
    }));

  const statusOptions = [
    { value: 'true', label: 'Hoạt động' },
    { value: 'false', label: 'Tạm dừng' }
  ];

  const hasActiveFilters = filter.search !== '' || filter.department !== '' || filter.active !== undefined;

  const handleClearFilters = () => {
    setFilter(prev => ({
      ...prev,
      search: '',
      department: '',
      active: undefined,
      page: 1
    }));
  };

  const FilterBarComponent = () => (
    <FilterBar hasActiveFilters={hasActiveFilters} onClear={handleClearFilters}>
      <SearchBox
        value={filter.search || ''}
        onChange={handleSearch}
        placeholder="Tìm kiếm theo tên hoặc mã nhân viên..."
      />
      
      <FilterGroup label="Phòng ban">
        <FilterSelect
          value={filter.department || ''}
          onChange={handleDepartmentFilter}
          options={departmentOptions}
          placeholder="Tất cả phòng ban"
        />
      </FilterGroup>

      <FilterGroup label="Trạng thái">
        <FilterSelect
          value={filter.active === undefined ? '' : filter.active.toString()}
          onChange={(value) => {
            handleActiveFilter(value === '' ? undefined : value === 'true');
          }}
          options={statusOptions}
          placeholder="Tất cả trạng thái"
        />
      </FilterGroup>

      <FilterGroup>
        <button
          type="button"
          onClick={() => setShowAddModal(true)}
          className={styles.addButton}
          aria-label="Add new employee"
        >
          <Plus size={20} />
          Thêm nhân viên
        </button>
      </FilterGroup>
    </FilterBar>
  );

  const EmptyState = () => (
    <div className={styles.emptyState}>
      <div className={styles.emptyIcon}>
        <Users size={48} />
      </div>
      <h3 className={styles.emptyTitle}>Chưa có nhân viên nào</h3>
      <p className={styles.emptyDescription}>
        Bắt đầu bằng cách thêm nhân viên đầu tiên vào hệ thống
      </p>
      <button
        type="button"
        onClick={() => setShowAddModal(true)}
        className={styles.emptyAction}
      >
        <Plus size={20} />
        Thêm nhân viên đầu tiên
      </button>
    </div>
  );

  const LoadingState = () => (
    <div className={styles.loading} aria-live="polite">
      <div className={styles.spinner} aria-hidden="true"></div>
      <span>Đang tải danh sách nhân viên...</span>
    </div>
  );

  const ErrorState = () => (
    <div className={styles.error} role="alert">
      <div className={styles.errorIcon}>
        <AlertTriangle size={24} />
      </div>
      <h3 className={styles.errorTitle}>Có lỗi xảy ra</h3>
      <p className={styles.errorMessage}>{error}</p>
      <button
        type="button"
        onClick={loadEmployees}
        className={styles.retryButton}
      >
        Thử lại
      </button>
    </div>
  );

  return (
    <div className={styles.employees}>
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <h1 className={styles.title}>Quản lý nhân viên</h1>
          <p className={styles.subtitle}>
            Quản lý thông tin và dữ liệu khuôn mặt của nhân viên
          </p>
        </div>
      </div>

      <Card className={styles.content}>
        <FilterBarComponent />
        
        {loading && <LoadingState />}
        {error && <ErrorState />}
        {!loading && !error && employees.length === 0 && <EmptyState />}
        {!loading && !error && employees.length > 0 && (
          <>
            <div className={styles.tableWrapper}>
              <Table
                columns={columns}
                data={employees}
                sortKey={filter.sortBy}
                sortDirection={filter.sortOrder}
                onSort={handleSort}
                className={styles.table}
              />
            </div>
            
            {total > filter.limit! && (
              <div className={styles.pagination}>
                <Pagination
                  current={filter.page!}
                  total={total}
                  pageSize={filter.limit!}
                  onChange={handlePageChange}
                />
              </div>
            )}
          </>
        )}
      </Card>

      {showAddModal && (
        <AddEmployeeModal
          onClose={() => setShowAddModal(false)}
          onSave={handleAddEmployee}
        />
      )}
    </div>
  );
}

// Simplified components for demo
const Pagination = ({ current, total, pageSize, onChange }: any) => (
  <div className={styles.paginationComponent}>
    <span>Showing {(current - 1) * pageSize + 1}-{Math.min(current * pageSize, total)} of {total}</span>
    <div className={styles.paginationButtons}>
      <button 
        disabled={current === 1}
        onClick={() => onChange(current - 1)}
      >
        Previous
      </button>
      <span>{current}</span>
      <button 
        disabled={current * pageSize >= total}
        onClick={() => onChange(current + 1)}
      >
        Next
      </button>
    </div>
  </div>
);

interface AddEmployeeModalProps {
  onClose: () => void;
  onSave: (data: any) => Promise<void>;
}

const AddEmployeeModal = ({ onClose, onSave }: AddEmployeeModalProps) => {
  const [formData, setFormData] = useState({
    code: '',
    name: '',
    email: '',
    phone: '',
    department: '',
    position: '',
    active: true
  });
  
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const newValue = type === 'checkbox' ? (e.target as HTMLInputElement).checked : value;
    
    setFormData(prev => ({
      ...prev,
      [name]: newValue
    }));
    
    // Clear error when user types
    if (errors[name]) {
      setErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.code.trim()) {
      newErrors.code = 'Mã nhân viên là bắt buộc';
    } else if (!/^[A-Z0-9]+$/.test(formData.code)) {
      newErrors.code = 'Mã nhân viên chỉ chứa chữ in hoa và số';
    }

    if (!formData.name.trim()) {
      newErrors.name = 'Tên nhân viên là bắt buộc';
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Tên phải có ít nhất 2 ký tự';
    }

    if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Email không hợp lệ';
    }

    if (formData.phone && !/^[0-9]{10,11}$/.test(formData.phone.replace(/\s/g, ''))) {
      newErrors.phone = 'Số điện thoại không hợp lệ (10-11 số)';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    try {
      setSubmitting(true);
      await onSave(formData);
    } catch (error) {
      console.error('Failed to save employee:', error);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className={styles.modal} onClick={onClose}>
      <div className={styles.modalContent} onClick={e => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h2 className={styles.modalTitle}>Thêm nhân viên mới</h2>
          <button
            type="button"
            onClick={onClose}
            className={styles.modalClose}
            aria-label="Close modal"
          >
            <CloseIcon />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGrid}>
            <div className={styles.formGroup}>
              <label htmlFor="code" className={styles.formLabel}>
                Mã nhân viên <span className={styles.required}>*</span>
              </label>
              <input
                type="text"
                id="code"
                name="code"
                value={formData.code}
                onChange={handleChange}
                className={`${styles.formInput} ${errors.code ? styles.formInputError : ''}`}
                placeholder="VD: EMP001"
                disabled={submitting}
              />
              {errors.code && <span className={styles.formError}>{errors.code}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="name" className={styles.formLabel}>
                Tên nhân viên <span className={styles.required}>*</span>
              </label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className={`${styles.formInput} ${errors.name ? styles.formInputError : ''}`}
                placeholder="Nguyễn Văn A"
                disabled={submitting}
              />
              {errors.name && <span className={styles.formError}>{errors.name}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="email" className={styles.formLabel}>
                Email
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className={`${styles.formInput} ${errors.email ? styles.formInputError : ''}`}
                placeholder="email@company.com"
                disabled={submitting}
              />
              {errors.email && <span className={styles.formError}>{errors.email}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="phone" className={styles.formLabel}>
                Số điện thoại
              </label>
              <input
                type="tel"
                id="phone"
                name="phone"
                value={formData.phone}
                onChange={handleChange}
                className={`${styles.formInput} ${errors.phone ? styles.formInputError : ''}`}
                placeholder="0901234567"
                disabled={submitting}
              />
              {errors.phone && <span className={styles.formError}>{errors.phone}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="department" className={styles.formLabel}>
                Phòng ban
              </label>
              <input
                type="text"
                id="department"
                name="department"
                value={formData.department}
                onChange={handleChange}
                className={styles.formInput}
                placeholder="Kỹ thuật"
                disabled={submitting}
              />
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="position" className={styles.formLabel}>
                Chức vụ
              </label>
              <input
                type="text"
                id="position"
                name="position"
                value={formData.position}
                onChange={handleChange}
                className={styles.formInput}
                placeholder="Developer"
                disabled={submitting}
              />
            </div>
          </div>

          <div className={styles.formGroup}>
            <label className={styles.formCheckbox}>
              <input
                type="checkbox"
                name="active"
                checked={formData.active}
                onChange={handleChange}
                disabled={submitting}
              />
              <span className={styles.formCheckboxLabel}>
                Kích hoạt tài khoản ngay
              </span>
            </label>
          </div>

          <div className={styles.modalActions}>
            <button 
              type="button"
              onClick={onClose}
              className={styles.modalButtonSecondary}
              disabled={submitting}
            >
              Hủy
            </button>
            <button 
              type="submit"
              className={styles.modalButtonPrimary}
              disabled={submitting}
            >
              {submitting ? (
                <>
                  <span className={styles.buttonSpinner}></span>
                  Đang lưu...
                </>
              ) : (
                <>
                  <SaveIcon />
                  Lưu nhân viên
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Additional Icons
const CloseIcon = () => (
  <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
    <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
  </svg>
);

const SaveIcon = () => (
  <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor">
    <path d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V6h5a2 2 0 012 2v7a2 2 0 01-2 2H4a2 2 0 01-2-2V8a2 2 0 012-2h5v5.586l-1.293-1.293zM9 4a1 1 0 012 0v2H9V4z" />
  </svg>
);