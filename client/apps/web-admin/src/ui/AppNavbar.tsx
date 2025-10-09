import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

interface AppNavbarProps {
  onToggleSidebar: () => void;
}

export function AppNavbar({ onToggleSidebar }: AppNavbarProps) {
  return (
    <Navbar bg="white" expand="lg" className="shadow-sm border-bottom" sticky="top">
      <Container fluid>
        <Button
          variant="outline-secondary"
          className="d-xl-none me-2"
          onClick={onToggleSidebar}
          aria-label="Toggle navigation"
        >
          <i className="bi bi-list" />
        </Button>
        <Navbar.Brand className="fw-semibold text-primary">Face Attendance</Navbar.Brand>
        <Nav className="ms-auto d-flex align-items-center gap-3">
          <Button variant="outline-secondary" size="sm" aria-label="Notifications">
            <i className="bi bi-bell" aria-hidden />
          </Button>
          <Button variant="primary" size="sm">
            <i className="bi bi-plus-lg me-2" aria-hidden />
            Tạo mới
          </Button>
        </Nav>
      </Container>
    </Navbar>
  );
}
