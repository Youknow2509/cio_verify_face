import { useEffect, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Image from 'react-bootstrap/Image';
import Row from 'react-bootstrap/Row';
import Spinner from 'react-bootstrap/Spinner';
import Stack from 'react-bootstrap/Stack';
import { useNavigate } from 'react-router-dom';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { fetchMyAccount } from '@/services/api/account';
import { clearAuthToken, HttpError } from '@/services/http';
import { formatDateTime } from '@/lib/format';
import type { AccountProfile, UserRole } from '@/types';

const roleLabels: Record<UserRole, string> = {
  CompanyAdmin: 'Quản trị doanh nghiệp',
  Manager: 'Quản lý',
  Staff: 'Nhân viên',
};

export default function AccountProfilePage() {
  const navigate = useNavigate();
  const { showToast } = useUi();
  const [profile, setProfile] = useState<AccountProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    const loadProfile = async () => {
      try {
        setLoading(true);
        const account = await fetchMyAccount();
        if (!isMounted) {
          return;
        }
        setProfile(account);
      } catch (error) {
        if (!isMounted) {
          return;
        }

        if (error instanceof HttpError && error.status === 401) {
          clearAuthToken();
          navigate('/login', { replace: true });
          return;
        }

        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: 'Không thể tải thông tin tài khoản',
        });
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    void loadProfile();

    return () => {
      isMounted = false;
    };
  }, [navigate, showToast]);

  const lastLogin = profile?.lastLoginAt ? formatDateTime(profile.lastLoginAt) : 'Chưa đăng nhập';

  const renderAvatar = () => {
    if (loading) {
      return (
        <div className="d-flex align-items-center justify-content-center" style={{ width: 64, height: 64 }}>
          <Spinner animation="border" role="status" aria-hidden />
        </div>
      );
    }

    if (profile?.avatarUrl) {
      return (
        <Image
          src={profile.avatarUrl}
          alt={`${profile.name} avatar`}
          width={64}
          height={64}
          roundedCircle
        />
      );
    }

    return (
      <div
        className="d-flex align-items-center justify-content-center rounded-circle bg-primary-subtle text-primary"
        style={{ width: 64, height: 64 }}
        role="img"
        aria-label="Không có ảnh đại diện"
      >
        <i className="bi bi-person" aria-hidden />
      </div>
    );
  };

  return (
    <Page
      title="Tài khoản của tôi"
      subtitle="Xem thông tin tài khoản và hoạt động gần đây"
      breadcrumb={[
        { label: 'Trang chủ', path: '/dashboard' },
        { label: 'Tài khoản' },
      ]}
    >
      <Card className="border-0 shadow-sm">
        <Card.Body>
          {loading ? (
            <div className="d-flex align-items-center gap-3 py-4">
              <Spinner animation="border" role="status" aria-hidden />
              <span>Đang tải thông tin tài khoản...</span>
            </div>
          ) : profile ? (
            <Row className="g-4 align-items-center">
              <Col xs="auto">{renderAvatar()}</Col>
              <Col>
                <Stack gap={1}>
                  <h2 className="fs-4 mb-0">{profile.name}</h2>
                  <div className="text-secondary">{profile.email}</div>
                  <div>
                    <Badge bg="primary" className="text-uppercase">
                      {roleLabels[profile.role] ?? profile.role}
                    </Badge>
                  </div>
                  <div className="text-secondary">
                    Lần đăng nhập gần nhất: <strong>{lastLogin}</strong>
                  </div>
                </Stack>
              </Col>
              <Col xs={12} md="auto" className="ms-md-auto">
                <Button variant="outline-primary" size="sm">
                  Edit basic info
                </Button>
              </Col>
            </Row>
          ) : (
            <div className="text-danger">Không có dữ liệu tài khoản.</div>
          )}
        </Card.Body>
      </Card>
    </Page>
  );
}
