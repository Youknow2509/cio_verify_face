import { useCallback, useEffect, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Dropdown from 'react-bootstrap/Dropdown';
import Image from 'react-bootstrap/Image';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Row from 'react-bootstrap/Row';
import Spinner from 'react-bootstrap/Spinner';
import { Link, Outlet, useNavigate } from 'react-router-dom';
import { UiProvider, useUi } from '@/app/providers/UiProvider';
import { AppSidebar } from '@/ui/AppSidebar';
import logoAsset from '@/assets/logo.svg';
import { fetchMyAccount, logout as apiLogout } from '@/services/api/account';
import { clearAuthToken, HttpError } from '@/services/http';
import type { AccountProfile } from '@/types';

export function Layout() {
  return (
    <UiProvider>
      <LayoutShell />
    </UiProvider>
  );
}

function LayoutShell() {
  const navigate = useNavigate();
  const { showToast } = useUi();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [profile, setProfile] = useState<AccountProfile | null>(null);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [signingOut, setSigningOut] = useState(false);

  useEffect(() => {
    let isMounted = true;

    const loadProfile = async () => {
      try {
        setLoadingProfile(true);
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
          setLoadingProfile(false);
        }
      }
    };

    void loadProfile();

    return () => {
      isMounted = false;
    };
  }, [navigate, showToast]);

  const handleLogout = useCallback(async () => {
    try {
      setSigningOut(true);
      await apiLogout();
      showToast({
        variant: 'success',
        title: 'Đăng xuất',
        message: 'Bạn đã đăng xuất thành công',
      });
    } catch (error) {
      showToast({
        variant: 'warning',
        title: 'Đăng xuất',
        message: 'Không thể kết nối máy chủ, phiên làm việc đã được kết thúc',
      });
    } finally {
      setSigningOut(false);
      navigate('/login', { replace: true });
    }
  }, [navigate, showToast]);

  const renderAvatar = () => {
    if (loadingProfile) {
      return <Spinner animation="border" size="sm" role="status" aria-hidden="true" />;
    }

    if (profile?.avatarUrl) {
      return (
        <Image
          src={profile.avatarUrl}
          alt={`${profile.name} avatar`}
          width={36}
          height={36}
          roundedCircle
        />
      );
    }

    return <i className="bi bi-person-circle fs-4 text-secondary" aria-hidden />;
  };

  return (
    <div className="d-flex flex-column min-vh-100 bg-body-tertiary">
      <Navbar bg="white" expand="lg" className="shadow-sm border-bottom" sticky="top">
        <Container fluid>
          <Button
            variant="outline-secondary"
            className="d-xl-none me-2"
            onClick={() => setSidebarOpen(true)}
            aria-label="Mở menu điều hướng"
          >
            <i className="bi bi-list" aria-hidden />
          </Button>
          <Navbar.Brand as={Link} to="/dashboard" className="d-flex align-items-center gap-2">
            <Image src={logoAsset} alt="Face Attendance logo" width={32} height={32} roundedCircle />
            <span className="fw-semibold text-primary">Face Attendance</span>
          </Navbar.Brand>
          <Nav className="ms-auto d-flex align-items-center gap-3">
            <Button variant="outline-secondary" size="sm" aria-label="Thông báo">
              <i className="bi bi-bell" aria-hidden />
            </Button>
            <Dropdown align="end">
              <Dropdown.Toggle
                id="navbar-user-menu"
                variant="link"
                className="p-0 border-0 shadow-none d-flex align-items-center"
                aria-label="Mở menu tài khoản"
              >
                {renderAvatar()}
              </Dropdown.Toggle>
              <Dropdown.Menu className="shadow-sm">
                <Dropdown.Header>
                  <div className="fw-semibold">{profile?.name ?? 'Người dùng'}</div>
                  <small className="text-muted">{profile?.email ?? ''}</small>
                </Dropdown.Header>
                <Dropdown.Divider />
                <Dropdown.Item as={Link} to="/account">
                  <i className="bi bi-person me-2" aria-hidden />
                  My Account
                </Dropdown.Item>
                <Dropdown.Item as={Link} to="/account/password">
                  <i className="bi bi-shield-lock me-2" aria-hidden />
                  Change Password
                </Dropdown.Item>
                <Dropdown.Divider />
                <Dropdown.Item onClick={handleLogout} disabled={signingOut} className="text-danger">
                  {signingOut ? (
                    <span className="d-inline-flex align-items-center gap-2">
                      <Spinner animation="border" size="sm" role="status" aria-hidden />
                      Đang đăng xuất...
                    </span>
                  ) : (
                    <span>
                      <i className="bi bi-box-arrow-right me-2" aria-hidden />
                      Sign out
                    </span>
                  )}
                </Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          </Nav>
        </Container>
      </Navbar>
      <Container fluid className="flex-grow-1 py-3">
        <Row className="flex-nowrap g-3">
          <Col xl="auto" className="d-none d-xl-flex">
            <AppSidebar variant="static" onNavigate={() => undefined} />
          </Col>
          <Col xs={12} xl className="px-2 px-md-3">
            <Outlet />
          </Col>
        </Row>
      </Container>
      <AppSidebar
        variant="offcanvas"
        show={sidebarOpen}
        onHide={() => setSidebarOpen(false)}
        onNavigate={() => setSidebarOpen(false)}
      />
    </div>
  );
}