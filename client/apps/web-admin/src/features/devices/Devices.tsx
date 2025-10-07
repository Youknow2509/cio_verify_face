// src/features/devices/Devices.tsx

import { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Card } from '@/components/Card/Card';
import { Table } from '@/components/Table/Table';
import { Badge } from '@/components/Badge/Badge';
import { 
  getDevices,
  syncDevice,
  deleteDevice
} from '@/services';
import type { Device, FilterOptions, TableColumn } from '@/types';
import styles from './Devices.module.scss';

export default function Devices() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<FilterOptions>({
    page: 1,
    limit: 20,
    search: '',
    status: '',
    sortBy: 'name',
    sortOrder: 'asc'
  });
  const [total, setTotal] = useState(0);
  const [syncingId, setSyncingId] = useState<string | null>(null);

  // Load devices data
  const loadDevices = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getDevices(filter);
      setDevices(response.data);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể tải danh sách thiết bị');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDevices();
  }, [filter]);

  // Calculate statistics
  const stats = useMemo(() => {
    const online = devices.filter(d => d.status === 'online').length;
    const offline = devices.filter(d => d.status === 'offline').length;
    
    return {
      total: devices.length,
      online,
      offline
    };
  }, [devices]);

  // Handle sync device
  const handleSync = async (deviceId: string) => {
    try {
      setSyncingId(deviceId);
      await syncDevice(deviceId);
      await loadDevices();
    } catch (err) {
      alert('Đồng bộ thất bại: ' + (err instanceof Error ? err.message : 'Unknown error'));
    } finally {
      setSyncingId(null);
    }
  };

  // Handle delete device
  const handleDelete = async (deviceId: string, deviceName: string) => {
    if (!confirm(`Bạn có chắc muốn xóa thiết bị "${deviceName}"?`)) {
      return;
    }

    try {
      await deleteDevice(deviceId);
      await loadDevices();
    } catch (err) {
      alert('Xóa thất bại: ' + (err instanceof Error ? err.message : 'Unknown error'));
    }
  };

  // Table columns definition
  const columns: TableColumn<Device>[] = [
    {
      key: 'name',
      header: 'Tên thiết bị',
      sortable: true,
      render: (_, device) => (
        <Link to={`/devices/${device.id}`} className={styles.deviceName}>
          {device.name}
        </Link>
      )
    },
    {
      key: 'location',
      header: 'Vị trí',
      sortable: true,
      render: (_, device) => (
        <div className={styles.deviceLocation}>
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
          </svg>
          {device.location}
        </div>
      )
    },
    {
      key: 'model',
      header: 'Model',
      sortable: true,
      render: (_, device) => (
        <span className={styles.deviceModel}>{device.model}</span>
      )
    },
    {
      key: 'ipAddress',
      header: 'Địa chỉ IP',
      render: (_, device) => (
        <span className={styles.deviceIp}>{device.ipAddress}</span>
      )
    },
    {
      key: 'status',
      header: 'Trạng thái',
      sortable: true,
      render: (_, device) => (
        <Badge 
          variant={device.status === 'online' ? 'success' : 'error'}
        >
          {device.status === 'online' ? '● Hoạt động' : '● Ngoại tuyến'}
        </Badge>
      )
    },
    {
      key: 'lastSyncAt',
      header: 'Đồng bộ lần cuối',
      sortable: true,
      render: (_, device) => (
        <span className={styles.lastSync}>
          {device.lastSyncAt ? new Date(device.lastSyncAt).toLocaleString('vi-VN') : 'Chưa đồng bộ'}
        </span>
      )
    },
    {
      key: 'actions',
      header: 'Thao tác',
      render: (_, device) => (
        <div className={styles.actionButtons}>
          <button
            className={`${styles.actionBtn} ${styles.sync}`}
            onClick={() => handleSync(device.id)}
            disabled={syncingId === device.id}
          >
            {syncingId === device.id ? '⟳ Đang đồng bộ...' : '⟳ Đồng bộ'}
          </button>
          <button
            className={`${styles.actionBtn} ${styles.delete}`}
            onClick={() => handleDelete(device.id, device.name)}
          >
            🗑 Xóa
          </button>
        </div>
      )
    }
  ];

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <h1>📱 Quản lý thiết bị</h1>
        <div className={styles.actions}>
          <button className={styles.refreshButton} onClick={loadDevices}>
            🔄 Làm mới
          </button>
          <button className={styles.addButton}>
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
            Thêm thiết bị
          </button>
        </div>
      </div>

      {/* Statistics Cards */}
      <div className={styles.statsCards}>
        <div className={styles.statCard}>
          <div className={`${styles.statIcon} ${styles.total}`}>
            📱
          </div>
          <div className={styles.statContent}>
            <h3 className={styles.statValue}>{stats.total}</h3>
            <p className={styles.statLabel}>Tổng số thiết bị</p>
          </div>
        </div>

        <div className={styles.statCard}>
          <div className={`${styles.statIcon} ${styles.online}`}>
            ✓
          </div>
          <div className={styles.statContent}>
            <h3 className={styles.statValue}>{stats.online}</h3>
            <p className={styles.statLabel}>Đang hoạt động</p>
          </div>
        </div>

        <div className={styles.statCard}>
          <div className={`${styles.statIcon} ${styles.offline}`}>
            ✕
          </div>
          <div className={styles.statContent}>
            <h3 className={styles.statValue}>{stats.offline}</h3>
            <p className={styles.statLabel}>Ngoại tuyến</p>
          </div>
        </div>


      </div>

      {/* Toolbar */}
      <div className={styles.toolbar}>
        <div className={styles.searchBox}>
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
          </svg>
          <input
            type="text"
            placeholder="Tìm kiếm theo tên, vị trí, IP..."
            value={filter.search}
            onChange={(e) => setFilter({ ...filter, search: e.target.value, page: 1 })}
          />
        </div>

        <div className={styles.filters}>
          <select
            className={styles.filterSelect}
            value={filter.status || ''}
            onChange={(e) => setFilter({ ...filter, status: e.target.value, page: 1 })}
          >
            <option value="">Tất cả trạng thái</option>
            <option value="online">Hoạt động</option>
            <option value="offline">Ngoại tuyến</option>
          </select>

          <select
            className={styles.filterSelect}
            value={filter.sortBy || 'name'}
            onChange={(e) => setFilter({ ...filter, sortBy: e.target.value })}
          >
            <option value="name">Sắp xếp: Tên</option>
            <option value="location">Sắp xếp: Vị trí</option>
            <option value="status">Sắp xếp: Trạng thái</option>
            <option value="lastSyncAt">Sắp xếp: Đồng bộ</option>
          </select>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className={styles.error}>
          ⚠️ {error}
        </div>
      )}

      {/* Devices Table */}
      <Card>
        <Table
          columns={columns}
          data={devices}
          loading={loading}
          empty={<div style={{ textAlign: 'center', padding: '40px', color: 'var(--md-sys-color-on-surface-variant)' }}>Không có thiết bị nào</div>}
        />
      </Card>

      {/* Pagination info */}
      {total > 0 && (
        <div style={{ marginTop: '16px', textAlign: 'center', color: 'var(--md-sys-color-on-surface-variant)' }}>
          Hiển thị {devices.length} / {total} thiết bị
        </div>
      )}
    </div>
  );
}