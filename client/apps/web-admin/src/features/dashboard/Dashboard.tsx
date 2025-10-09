import { useEffect, useMemo, useState } from 'react';
import Alert from 'react-bootstrap/Alert';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Spinner from 'react-bootstrap/Spinner';
import Stack from 'react-bootstrap/Stack';
import { AttendanceChart } from '@/components/charts/AttendanceChart';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import {
  getAttendanceChart,
  getDashboardStats,
  getRecentActivity,
} from '@/services/mock/attendance';
import type { ChartData, DashboardStats, RecentActivity } from '@/types';

interface StatCardConfig {
  title: string;
  value: string;
  accent: string;
  icon: string;
  helper?: string;
}

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
  const [chartData, setChartData] = useState<ChartData[]>([]);
  const [loading, setLoading] = useState(true);
  const { showToast } = useUi();

  useEffect(() => {
    const loadData = async () => {
      try {
        const [statsResponse, activityResponse, chartResponse] = await Promise.all([
          getDashboardStats(),
          getRecentActivity(),
          getAttendanceChart(),
        ]);

        if (statsResponse.error) {
          throw new Error(statsResponse.error);
        }
        if (activityResponse.error) {
          throw new Error(activityResponse.error);
        }
        if (chartResponse.error) {
          throw new Error(chartResponse.error);
        }

        setStats(statsResponse.data);
        setRecentActivity(activityResponse.data ?? []);
        setChartData(chartResponse.data ?? []);
      } catch (error) {
        console.error(error);
        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể tải dữ liệu dashboard',
        });
      } finally {
        setLoading(false);
      }
    };

    void loadData();
  }, [showToast]);

  const statCards = useMemo<StatCardConfig[]>(() => {
    if (!stats) return [];
    return [
      {
        title: 'Tổng nhân viên',
        value: stats.totalEmployees.toLocaleString('vi-VN'),
        accent: 'primary',
        icon: 'bi-people',
        helper: '+12% so với tháng trước',
      },
      {
        title: 'Check-in hôm nay',
        value: stats.todayCheckIns.toLocaleString('vi-VN'),
        accent: 'success',
        icon: 'bi-check-circle',
        helper: 'Cập nhật 5 phút trước',
      },
      {
        title: 'Đi trễ hôm nay',
        value: stats.lateArrivals.toLocaleString('vi-VN'),
        accent: 'warning',
        icon: 'bi-alarm',
        helper: '-8% so với hôm qua',
      },
      {
        title: 'Thiết bị online',
        value: stats.devicesOnline.toLocaleString('vi-VN'),
        accent: 'info',
        icon: 'bi-hdd-network',
        helper: `/${stats.devicesOnline + 2} thiết bị tổng`,
      },
    ];
  }, [stats]);

  const activityLabel = (type: RecentActivity['type']) => {
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

  const activityVariant = (type: RecentActivity['type']) => {
    switch (type) {
      case 'check_in':
        return 'success';
      case 'check_out':
        return 'primary';
      case 'device_sync':
        return 'info';
      case 'employee_added':
        return 'warning';
      default:
        return 'secondary';
    }
  };

  const formatTimeAgo = (timestamp: string) => {
    const now = new Date();
    const time = new Date(timestamp);
    const diff = Math.floor((now.getTime() - time.getTime()) / 1000);

    if (diff < 60) return 'Vừa xong';
    if (diff < 3600) return `${Math.floor(diff / 60)} phút trước`;
    if (diff < 86400) return `${Math.floor(diff / 3600)} giờ trước`;
    return `${Math.floor(diff / 86400)} ngày trước`;
  };

  const actions = (
    <Button variant="outline-secondary">
      <i className="bi bi-arrow-clockwise me-2" aria-hidden />
      Làm mới
    </Button>
  );

  return (
    <Page
      title="Dashboard"
      subtitle="Tổng quan hệ thống chấm công"
      actions={actions}
      breadcrumb={[{ label: 'Trang chủ' }, { label: 'Dashboard' }]}
    >
      {loading && (
        <Card className="border-0 shadow-sm">
          <Card.Body className="d-flex align-items-center justify-content-center py-5 gap-3">
            <Spinner animation="border" role="status" aria-hidden />
            <span>Đang tải dữ liệu...</span>
          </Card.Body>
        </Card>
      )}

      {!loading && stats && (
        <Row className="g-3">
          {statCards.map((card) => (
            <Col key={card.title} sm={6} xl={3}>
              <Card className="border-0 shadow-sm h-100">
                <Card.Body>
                  <Stack direction="horizontal" gap={3}>
                    <div
                      className={`d-flex align-items-center justify-content-center rounded-circle bg-${card.accent}-subtle text-${card.accent}`}
                      style={{ width: 48, height: 48 }}
                    >
                      <i className={`bi ${card.icon} fs-4`} aria-hidden />
                    </div>
                    <div>
                      <p className="text-secondary text-uppercase small mb-1">{card.title}</p>
                      <h3 className="fw-semibold mb-0">{card.value}</h3>
                      {card.helper && <small className="text-secondary">{card.helper}</small>}
                    </div>
                  </Stack>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}

      {!loading && (
        <Row className="g-3 mt-1">
          <Col xl={8}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white border-0 pb-0">
                <Stack direction="horizontal" className="justify-content-between align-items-center">
                  <div>
                    <Card.Title className="mb-0">Biểu đồ chấm công 7 ngày</Card.Title>
                    <small className="text-secondary">Theo dõi check-in, check-out và đi trễ</small>
                  </div>
                </Stack>
              </Card.Header>
              <Card.Body style={{ minHeight: 320 }}>
                {chartData.length > 0 ? (
                  <AttendanceChart data={chartData} height={260} />
                ) : (
                  <div className="d-flex align-items-center justify-content-center text-secondary h-100">
                    <span>Chưa có dữ liệu để hiển thị</span>
                  </div>
                )}
              </Card.Body>
            </Card>
          </Col>
          <Col xl={4}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white border-0 pb-0">
                <Card.Title className="mb-0">Hoạt động gần đây</Card.Title>
              </Card.Header>
              <Card.Body>
                {recentActivity.length === 0 ? (
                  <Alert variant="light" className="mb-0">
                    Chưa có hoạt động nào gần đây
                  </Alert>
                ) : (
                  <Stack gap={3}>
                    {recentActivity.slice(0, 5).map((activity) => (
                      <Stack direction="horizontal" gap={3} key={activity.id}>
                        <div className="flex-shrink-0">
                          <div className="rounded-circle bg-primary-subtle text-primary d-flex align-items-center justify-content-center" style={{ width: 40, height: 40 }}>
                            <i className="bi bi-activity" aria-hidden />
                          </div>
                        </div>
                        <div className="flex-grow-1">
                          <div className="fw-semibold">
                            {activity.employeeName ?? activity.deviceName ?? 'Hệ thống'}
                          </div>
                          <small className="text-secondary">{activity.message}</small>
                          <div className="text-secondary small">{formatTimeAgo(activity.timestamp)}</div>
                        </div>
                        <Badge bg={activityVariant(activity.type)} pill>
                          {activityLabel(activity.type)}
                        </Badge>
                      </Stack>
                    ))}
                  </Stack>
                )}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}

      {!loading && (
        <Row className="g-3 mt-1">
          <Col xl={8}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white border-0 pb-0">
                <Card.Title className="mb-0">Thao tác nhanh</Card.Title>
              </Card.Header>
              <Card.Body>
                <Row className="g-3">
                  {[
                    { label: 'Thêm nhân viên', icon: 'bi-person-plus-fill' },
                    { label: 'Thêm thiết bị', icon: 'bi-cpu' },
                    { label: 'Xuất báo cáo', icon: 'bi-download' },
                    { label: 'Cài đặt hệ thống', icon: 'bi-gear' },
                  ].map((action) => (
                    <Col sm={6} key={action.label}>
                      <Button variant="outline-primary" className="w-100 d-flex align-items-center justify-content-between">
                        <span>{action.label}</span>
                        <i className={`bi ${action.icon}`} aria-hidden />
                      </Button>
                    </Col>
                  ))}
                </Row>
              </Card.Body>
            </Card>
          </Col>
          <Col xl={4}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Header className="bg-white border-0 pb-0">
                <Card.Title className="mb-0">Sắp tới</Card.Title>
              </Card.Header>
              <Card.Body>
                <Stack gap={3}>
                  {[
                    {
                      day: '15',
                      month: 'Th10',
                      title: 'Báo cáo tháng',
                      description: 'Hạn nộp báo cáo chấm công',
                    },
                    {
                      day: '20',
                      month: 'Th10',
                      title: 'Bảo trì hệ thống',
                      description: 'Nâng cấp phần mềm định kỳ',
                    },
                  ].map((event) => (
                    <Stack direction="horizontal" gap={3} key={event.title}>
                      <div className="rounded bg-primary text-white d-flex flex-column align-items-center justify-content-center" style={{ width: 56, height: 64 }}>
                        <span className="fw-semibold fs-4">{event.day}</span>
                        <small>{event.month}</small>
                      </div>
                      <div>
                        <div className="fw-semibold">{event.title}</div>
                        <small className="text-secondary">{event.description}</small>
                      </div>
                    </Stack>
                  ))}
                </Stack>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      )}
    </Page>
  );
}