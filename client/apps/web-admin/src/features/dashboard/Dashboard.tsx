// src/features/dashboard/Dashboard.tsx

import { useEffect, useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/Card/Card';
import { Badge } from '@/components/Badge/Badge';
import { AttendanceChart } from '@/components/charts/AttendanceChart';
import { getDashboardStats, getRecentActivity, getAttendanceChart } from '@/services/mock/attendance';
import type { DashboardStats, RecentActivity, ChartData } from '@/types';
import styles from './Dashboard.module.scss';

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
  const [chartData, setChartData] = useState<ChartData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadDashboardData = async () => {
      try {
        const [statsResponse, activityResponse, chartResponse] = await Promise.all([
          getDashboardStats(),
          getRecentActivity(),
          getAttendanceChart()
        ]);
        
        if (statsResponse.data) {
          setStats(statsResponse.data);
        }
        
        if (activityResponse.data) {
          setRecentActivity(activityResponse.data);
        }
        
        if (chartResponse.data) {
          setChartData(chartResponse.data);
        }
      } catch (error) {
        console.error('Error loading dashboard data:', error);
      } finally {
        setLoading(false);
      }
    };

    loadDashboardData();
  }, []);

  const getActivityVariant = (type: string) => {
    switch (type) {
      case 'check_in':
        return 'success';
      case 'check_out':
        return 'info';
      case 'device_sync':
        return 'neutral';
      case 'employee_added':
        return 'warning';
      default:
        return 'neutral';
    }
  };

  const getActivityLabel = (type: string) => {
    switch (type) {
      case 'check_in':
        return 'Check-in';
      case 'check_out':
        return 'Check-out';
      case 'device_sync':
        return 'Đồng bộ';
      case 'employee_added':
        return 'Nhân viên mới';
      default:
        return type;
    }
  };

  if (loading) {
    return (
      <div className={styles.loading}>
        <div className={styles.spinner}></div>
        <span>Đang tải dữ liệu...</span>
      </div>
    );
  }

  return (
    <div className={styles.dashboard}>
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <div>
            <h1 className={styles.title}>Dashboard</h1>
            <p className={styles.subtitle}>Tổng quan hệ thống chấm công</p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.refreshButton}>
              <RefreshIcon />
              Làm mới
            </button>
          </div>
        </div>
      </div>

      {stats && (
        <div className={styles.statsGrid}>
          <Card className={`${styles.statCard} ${styles.statCardPrimary}`}>
            <div className={styles.statCardInner}>
              <div className={styles.statIcon}>
                <UsersIcon />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statLabel}>Tổng nhân viên</div>
                <div className={styles.statValue}>{stats.totalEmployees}</div>
                <div className={styles.statChange}>
                  <TrendUpIcon />
                  <span>+12% so với tháng trước</span>
                </div>
              </div>
            </div>
          </Card>

          <Card className={`${styles.statCard} ${styles.statCardSuccess}`}>
            <div className={styles.statCardInner}>
              <div className={styles.statIcon}>
                <CheckCircleIcon />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statLabel}>Check-in hôm nay</div>
                <div className={styles.statValue}>{stats.todayCheckIns}</div>
                <div className={styles.statChange}>
                  <ClockIcon />
                  <span>Cập nhật 5 phút trước</span>
                </div>
              </div>
            </div>
          </Card>

          <Card className={`${styles.statCard} ${styles.statCardWarning}`}>
            <div className={styles.statCardInner}>
              <div className={styles.statIcon}>
                <AlertIcon />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statLabel}>Đi trễ hôm nay</div>
                <div className={styles.statValue}>{stats.lateArrivals}</div>
                <div className={styles.statChange}>
                  <TrendDownIcon />
                  <span>-8% so với hôm qua</span>
                </div>
              </div>
            </div>
          </Card>

          <Card className={`${styles.statCard} ${styles.statCardInfo}`}>
            <div className={styles.statCardInner}>
              <div className={styles.statIcon}>
                <DevicesIcon />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statLabel}>Thiết bị online</div>
                <div className={styles.statValue}>{stats.devicesOnline}</div>
                <div className={styles.statChange}>
                  <span className={styles.statTotal}>/ {stats.devicesOnline + 2} thiết bị</span>
                </div>
              </div>
            </div>
          </Card>
        </div>
      )}

      <div className={styles.contentGrid}>
        <div className={styles.mainColumn}>
          <Card className={styles.chartCard}>
            <CardHeader>
              <CardTitle>Biểu đồ chấm công 7 ngày</CardTitle>
              <div className={styles.chartLegend}>
                <span className={styles.legendItem}>
                  <span className={styles.legendDot} style={{ backgroundColor: '#1976d2' }}></span>
                  Check-in
                </span>
                <span className={styles.legendItem}>
                  <span className={styles.legendDot} style={{ backgroundColor: '#4caf50' }}></span>
                  Check-out
                </span>
              </div>
            </CardHeader>
            <CardContent>
              <AttendanceChart data={chartData} height={240} />
            </CardContent>
          </Card>

          <Card className={styles.activityCard}>
            <CardHeader>
              <CardTitle>Hoạt động gần đây</CardTitle>
              <a href="/activities" className={styles.viewAllLink}>
                Xem tất cả
                <ArrowRightIcon />
              </a>
            </CardHeader>
            <CardContent>
              <div className={styles.activityList}>
                {recentActivity.slice(0, 3).map((activity, index) => (
                  <div key={index} className={styles.activityItemWrapper}>
                    <div className={styles.activityIcon}>
                      {activity.type === 'check_in' && <CheckInIcon />}
                      {activity.type === 'check_out' && <CheckOutIcon />}
                      {activity.type === 'device_sync' && <SyncIcon />}
                      {activity.type === 'employee_added' && <UserPlusIcon />}
                    </div>
                    <div className={styles.activityDetails}>
                      <div className={styles.activityMessage}>
                        {activity.employeeName && <strong>{activity.employeeName}</strong>}
                        {activity.deviceName && <strong>{activity.deviceName}</strong>}
                        {' '}{activity.message}
                      </div>
                      <div className={styles.activityTime}>
                        {formatTimeAgo(activity.timestamp)}
                      </div>
                    </div>
                    <Badge variant={getActivityVariant(activity.type)} size="sm">
                      {getActivityLabel(activity.type)}
                    </Badge>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        <div className={styles.sideColumn}>
          <Card className={styles.quickActionsCard}>
            <CardHeader>
              <CardTitle>Thao tác nhanh</CardTitle>
            </CardHeader>
            <CardContent>
              <div className={styles.quickActions}>
                <button className={styles.quickAction}>
                  <UserPlusIcon />
                  <span>Thêm nhân viên</span>
                </button>
                <button className={styles.quickAction}>
                  <DevicePlusIcon />
                  <span>Thêm thiết bị</span>
                </button>
                <button className={styles.quickAction}>
                  <DownloadIcon />
                  <span>Xuất báo cáo</span>
                </button>
                <button className={styles.quickAction}>
                  <SettingsIcon />
                  <span>Cài đặt hệ thống</span>
                </button>
              </div>
            </CardContent>
          </Card>

          <Card className={styles.upcomingCard}>
            <CardHeader>
              <CardTitle>Sắp tới</CardTitle>
            </CardHeader>
            <CardContent>
              <div className={styles.upcomingList}>
                <div className={styles.upcomingItem}>
                  <div className={styles.upcomingDate}>
                    <div className={styles.upcomingDay}>15</div>
                    <div className={styles.upcomingMonth}>Th10</div>
                  </div>
                  <div className={styles.upcomingDetails}>
                    <div className={styles.upcomingTitle}>Báo cáo tháng</div>
                    <div className={styles.upcomingDesc}>Hạn nộp báo cáo chấm công</div>
                  </div>
                </div>
                <div className={styles.upcomingItem}>
                  <div className={styles.upcomingDate}>
                    <div className={styles.upcomingDay}>20</div>
                    <div className={styles.upcomingMonth}>Th10</div>
                  </div>
                  <div className={styles.upcomingDetails}>
                    <div className={styles.upcomingTitle}>Bảo trì hệ thống</div>
                    <div className={styles.upcomingDesc}>Nâng cấp phần mềm định kỳ</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

// Helper function
function formatTimeAgo(timestamp: string): string {
  const now = new Date();
  const time = new Date(timestamp);
  const diff = Math.floor((now.getTime() - time.getTime()) / 1000); // seconds

  if (diff < 60) return 'Vừa xong';
  if (diff < 3600) return `${Math.floor(diff / 60)} phút trước`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} giờ trước`;
  return `${Math.floor(diff / 86400)} ngày trước`;
}

// Icons
const RefreshIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
  </svg>
);

const UsersIcon = () => (
  <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
    <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
  </svg>
);

const CheckCircleIcon = () => (
  <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
  </svg>
);

const AlertIcon = () => (
  <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
    <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
  </svg>
);

const DevicesIcon = () => (
  <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
    <path d="M4 6h18V4H4c-1.1 0-2 .9-2 2v11H0v3h14v-3H4V6zm19 2h-6c-.55 0-1 .45-1 1v10c0 .55.45 1 1 1h6c.55 0 1-.45 1-1V9c0-.55-.45-1-1-1zm-1 9h-4v-7h4v7z"/>
  </svg>
);

const TrendUpIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
    <path d="M16 6l2.29 2.29-4.88 4.88-4-4L2 16.59 3.41 18l6-6 4 4 6.3-6.29L22 12V6z"/>
  </svg>
);

const TrendDownIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
    <path d="M16 18l2.29-2.29-4.88-4.88-4 4L2 7.41 3.41 6l6 6 4-4 6.3 6.29L22 12v6z"/>
  </svg>
);

const ClockIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
    <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z"/>
  </svg>
);

const ArrowRightIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 4l-1.41 1.41L16.17 11H4v2h12.17l-5.58 5.59L12 20l8-8z"/>
  </svg>
);

const CheckInIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M9 11H7v2h2v-2zm4 0h-2v2h2v-2zm4 0h-2v2h2v-2zm2-7h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V9h14v11z"/>
  </svg>
);

const CheckOutIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-1.99.9-1.99 2L3 19c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11zM7 10h5v5H7z"/>
  </svg>
);

const SyncIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 4V1L8 5l4 4V6c3.31 0 6 2.69 6 6 0 1.01-.25 1.97-.7 2.8l1.46 1.46C19.54 15.03 20 13.57 20 12c0-4.42-3.58-8-8-8zm0 14c-3.31 0-6-2.69-6-6 0-1.01.25-1.97.7-2.8L5.24 7.74C4.46 8.97 4 10.43 4 12c0 4.42 3.58 8 8 8v3l4-4-4-4v3z"/>
  </svg>
);

const UserPlusIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M15 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm-9-2V7H4v3H1v2h3v3h2v-3h3v-2H6zm9 4c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
  </svg>
);

const DevicePlusIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M18 1.01L8 1c-1.1 0-2 .9-2 2v3h2V5h10v14H8v-1H6v3c0 1.1.9 2 2 2h10c1.1 0 2-.9 2-2V3c0-1.1-.9-1.99-2-1.99zM10 15h2V8H5v2h3.59L3 15.59 4.41 17 10 11.41z"/>
  </svg>
);

const DownloadIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2z"/>
  </svg>
);

const SettingsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>
  </svg>
);