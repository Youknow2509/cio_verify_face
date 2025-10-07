// src/features/employees/Employees.tsx

import { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Card } from '@/components/Card/Card';
import { Table } from '@/components/Table/Table';
import { Badge } from '@/components/Badge/Badge';
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
            to={`/employees/${record.id}`}
            className={styles.actionButton}
            aria-label={`Edit employee ${record.name}`}
          >
            <EditIcon />
          </Link>
          <button
            type="button"
            onClick={() => handleDeleteEmployee(record.id)}
            className={`${styles.actionButton} ${styles.danger}`}
            aria-label={`Delete employee ${record.name}`}
          >
            <DeleteIcon />
          </button>
        </div>
      )
    }
  ];

  const FilterBar = () => (
    <div className={styles.filterBar}>
      <div className={styles.filterBarLeft}>
        <div className={styles.searchBox}>
          <SearchIcon />
          <input
            type="text"
            placeholder="Tìm kiếm theo tên hoặc mã nhân viên..."
            value={filter.search || ''}
            onChange={(e) => handleSearch(e.target.value)}
            className={styles.searchInput}
            aria-label="Search employees"
          />
        </div>
        
        <div className={styles.filters}>
          <select
            value={filter.department || ''}
            onChange={(e) => handleDepartmentFilter(e.target.value)}
            className={styles.select}
            aria-label="Filter by department"
          >
            <option value="">Tất cả phòng ban</option>
            {departments.map(dept => (
              <option key={dept} value={dept}>{dept}</option>
            ))}
          </select>

          <select
            value={filter.active === undefined ? '' : filter.active.toString()}
            onChange={(e) => {
              const value = e.target.value;
              handleActiveFilter(value === '' ? undefined : value === 'true');
            }}
            className={styles.select}
            aria-label="Filter by status"
          >
            <option value="">Tất cả trạng thái</option>
            <option value="true">Hoạt động</option>
            <option value="false">Tạm dừng</option>
          </select>
        </div>
      </div>

      <div className={styles.filterBarRight}>
        <button
          type="button"
          onClick={() => setShowAddModal(true)}
          className={styles.addButton}
          aria-label="Add new employee"
        >
          <PlusIcon />
          Thêm nhân viên
        </button>
      </div>
    </div>
  );

  const EmptyState = () => (
    <div className={styles.emptyState}>
      <div className={styles.emptyIcon}>
        <UsersIcon />
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
        <PlusIcon />
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
        <AlertIcon />
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
        <FilterBar />
        
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

// Icons
const SearchIcon = () => (
  <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
    <path fillRule="evenodd" d="M12.9 14.32a8 8 0 111.41-1.41l5.35 5.33-1.42 1.42-5.33-5.34zM8 14A6 6 0 108 2a6 6 0 000 12z" clipRule="evenodd" />
  </svg>
);

const PlusIcon = () => (
  <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
    <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
  </svg>
);

const EditIcon = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
    <path d="M11.013 1.427a1.75 1.75 0 012.474 0l1.086 1.086a1.75 1.75 0 010 2.474l-8.61 8.61c-.21.21-.47.364-.756.445l-3.251.93a.75.75 0 01-.927-.928l.929-3.25a1.75 1.75 0 01.445-.758l8.61-8.61zm1.414 1.06a.25.25 0 00-.354 0L10.811 3.75l1.439 1.44 1.263-1.263a.25.25 0 000-.354l-1.086-1.086zM11.189 6.25L9.75 4.81l-6.286 6.287a.25.25 0 00-.064.108l-.558 1.953 1.953-.558a.249.249 0 00.108-.064l6.286-6.286z" />
  </svg>
);

const DeleteIcon = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
    <path fillRule="evenodd" d="M6.5 1.75a.25.25 0 01.25-.25h2.5a.25.25 0 01.25.25V3h-3V1.75zm4.5 0V3h2.25a.75.75 0 010 1.5H2.75a.75.75 0 010-1.5H5V1.75C5 .784 5.784 0 6.75 0h2.5C10.216 0 11 .784 11 1.75zM4.496 6.675a.75.75 0 10-1.492.15l.66 6.6A1.75 1.75 0 005.405 15h5.19c.9 0 1.652-.681 1.741-1.576l.66-6.6a.75.75 0 00-1.492-.149l-.66 6.6a.25.25 0 01-.249.225h-5.19a.25.25 0 01-.249-.225l-.66-6.6z" clipRule="evenodd" />
  </svg>
);

const UsersIcon = () => (
  <svg width="48" height="48" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 4.5C13.6569 4.5 15 5.84315 15 7.5C15 9.15685 13.6569 10.5 12 10.5C10.3431 10.5 9 9.15685 9 7.5C9 5.84315 10.3431 4.5 12 4.5Z"/>
    <path d="M17.5 7.5C17.5 8.32843 16.8284 9 16 9C15.1716 9 14.5 8.32843 14.5 7.5C14.5 6.67157 15.1716 6 16 6C16.8284 6 17.5 6.67157 17.5 7.5Z"/>
    <path d="M9.5 7.5C9.5 8.32843 8.82843 9 8 9C7.17157 9 6.5 8.32843 6.5 7.5C6.5 6.67157 7.17157 6 8 6C8.82843 6 9.5 6.67157 9.5 7.5Z"/>
    <path d="M7 14C5.34315 14 4 15.3431 4 17V18.5C4 19.3284 4.67157 20 5.5 20H18.5C19.3284 20 20 19.3284 20 18.5V17C20 15.3431 18.6569 14 17 14H7Z"/>
  </svg>
);

const AlertIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
    <path fillRule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM12 6a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0112 6zm0 9a1 1 0 100-2 1 1 0 000 2z" clipRule="evenodd" />
  </svg>
);

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