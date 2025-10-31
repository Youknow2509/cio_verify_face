import { useMemo } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import Card from 'react-bootstrap/Card';
import Nav from 'react-bootstrap/Nav';
import Offcanvas from 'react-bootstrap/Offcanvas';
import Stack from 'react-bootstrap/Stack';

interface NavItem {
  to: string;
  icon: string;
  label: string;
}

interface AppSidebarProps {
  variant: 'static' | 'offcanvas';
  show?: boolean;
  onHide?: () => void;
  onNavigate: () => void;
}

export function AppSidebar({ variant, show = false, onHide, onNavigate }: AppSidebarProps) {
  const location = useLocation();

  const items = useMemo<NavItem[]>(
    () => [
      { to: '/dashboard', icon: 'bi-speedometer', label: 'Dashboard' },
      { to: '/employees', icon: 'bi-people', label: 'Nhân viên' },
      { to: '/attendance', icon: 'bi-calendar-check', label: 'Chấm công' },
      { to: '/reports', icon: 'bi-file-earmark-bar-graph', label: 'Báo cáo' },
      { to: '/shifts', icon: 'bi-clock-history', label: 'Ca làm việc' },
      { to: '/devices', icon: 'bi-cpu', label: 'Thiết bị' },
      { to: '/settings', icon: 'bi-gear', label: 'Cài đặt' },
    ],
    []
  );

  const navList = (
    <Nav className="flex-column" aria-label="Main navigation">
      <Stack gap={2}>
        {items.map((item) => {
          const isActive = location.pathname.startsWith(item.to);
          return (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive: linkActive }) =>
                [
                  'btn',
                  'btn-light',
                  'd-flex',
                  'align-items-center',
                  'gap-2',
                  'justify-content-start',
                  linkActive || isActive ? 'active bg-primary text-white border-0' : 'text-secondary',
                ].join(' ')
              }
              onClick={onNavigate}
            >
              <i className={`bi ${item.icon}`} aria-hidden />
              <span>{item.label}</span>
            </NavLink>
          );
        })}
      </Stack>
    </Nav>
  );

  if (variant === 'offcanvas') {
    return (
      <Offcanvas show={show} onHide={onHide} placement="start" className="bg-body" aria-label="Mobile navigation">
        <Offcanvas.Header closeButton>
          <Offcanvas.Title className="fw-semibold text-primary">Điều hướng</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>{navList}</Offcanvas.Body>
      </Offcanvas>
    );
  }

  return (
    <Card className="shadow-sm bg-white border-0 h-100" aria-label="Sidebar navigation">
      <Card.Body>{navList}</Card.Body>
    </Card>
  );
}
