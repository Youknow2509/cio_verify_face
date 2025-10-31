// src/features/settings/Settings.tsx

import { useState, type ReactNode } from 'react';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Nav from 'react-bootstrap/Nav';
import Row from 'react-bootstrap/Row';
import Stack from 'react-bootstrap/Stack';
import Tab from 'react-bootstrap/Tab';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';

export default function Settings() {
  const { showToast } = useUi();
  const [general, setGeneral] = useState({
    companyName: 'Công ty ABC',
    timezone: 'Asia/Bangkok',
    language: 'vi',
  });
  const [attendance, setAttendance] = useState({
    startTime: '08:00',
    endTime: '17:00',
    allowLateMinutes: 15,
    mobileCheckIn: true,
    requirePhoto: true,
  });
  const [notifications, setNotifications] = useState({
    email: true,
    sms: false,
    inApp: true,
  });
  const [security, setSecurity] = useState({
    twoFactor: true,
    autoLogout: true,
  });

  const handleSave = (section: 'general' | 'attendance' | 'notifications' | 'security') => {
    showToast({
      variant: 'success',
      message: `Đã lưu cài đặt ${translateSection(section)} thành công.`,
    });
  };

  return (
    <Page
      title="Cài đặt hệ thống"
      subtitle="Điều chỉnh các tham số vận hành, thông báo và bảo mật cho hệ thống chấm công"
      breadcrumb={[{ label: 'Trang chủ', path: '/dashboard' }, { label: 'Cài đặt' }]}
    >
      <Tab.Container defaultActiveKey="general">
        <Row className="g-3">
          <Col xl={3} lg={4}>
            <Card className="border-0 shadow-sm h-100">
              <Card.Body className="p-0">
                <Nav variant="pills" className="flex-column">
                  <Nav.Item>
                    <Nav.Link eventKey="general" className="py-3 px-4">
                      Cài đặt chung
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="attendance" className="py-3 px-4">
                      Chấm công
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="notifications" className="py-3 px-4">
                      Thông báo
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="security" className="py-3 px-4">
                      Bảo mật
                    </Nav.Link>
                  </Nav.Item>
                </Nav>
              </Card.Body>
            </Card>
          </Col>
          <Col xl={9} lg={8}>
            <Tab.Content>
              <Tab.Pane eventKey="general">
                <SettingsSection
                  title="Cài đặt chung"
                  description="Thông tin tổ chức và ngôn ngữ hiển thị."
                  onSave={() => handleSave('general')}
                >
                  <Row className="g-3">
                    <Col md={6}>
                      <Form.Group controlId="company-name">
                        <Form.Label>Tên công ty</Form.Label>
                        <Form.Control
                          value={general.companyName}
                          onChange={(event) =>
                            setGeneral((prev) => ({ ...prev, companyName: event.target.value }))
                          }
                        />
                      </Form.Group>
                    </Col>
                    <Col md={6}>
                      <Form.Group controlId="company-timezone">
                        <Form.Label>Múi giờ</Form.Label>
                        <Form.Select
                          value={general.timezone}
                          onChange={(event) =>
                            setGeneral((prev) => ({ ...prev, timezone: event.target.value }))
                          }
                        >
                          <option value="Asia/Bangkok">GMT+7 (Hà Nội)</option>
                          <option value="Asia/Singapore">GMT+8 (Singapore)</option>
                          <option value="Asia/Tokyo">GMT+9 (Tokyo)</option>
                        </Form.Select>
                      </Form.Group>
                    </Col>
                    <Col md={6}>
                      <Form.Group controlId="company-language">
                        <Form.Label>Ngôn ngữ</Form.Label>
                        <Form.Select
                          value={general.language}
                          onChange={(event) =>
                            setGeneral((prev) => ({ ...prev, language: event.target.value }))
                          }
                        >
                          <option value="vi">Tiếng Việt</option>
                          <option value="en">English</option>
                        </Form.Select>
                      </Form.Group>
                    </Col>
                  </Row>
                </SettingsSection>
              </Tab.Pane>

              <Tab.Pane eventKey="attendance">
                <SettingsSection
                  title="Cài đặt chấm công"
                  description="Thiết lập khung giờ làm việc và chính sách chấm công."
                  onSave={() => handleSave('attendance')}
                >
                  <Row className="g-3">
                    <Col md={4}>
                      <Form.Group controlId="attendance-start">
                        <Form.Label>Giờ bắt đầu</Form.Label>
                        <Form.Control
                          type="time"
                          value={attendance.startTime}
                          onChange={(event) =>
                            setAttendance((prev) => ({ ...prev, startTime: event.target.value }))
                          }
                        />
                      </Form.Group>
                    </Col>
                    <Col md={4}>
                      <Form.Group controlId="attendance-end">
                        <Form.Label>Giờ kết thúc</Form.Label>
                        <Form.Control
                          type="time"
                          value={attendance.endTime}
                          onChange={(event) =>
                            setAttendance((prev) => ({ ...prev, endTime: event.target.value }))
                          }
                        />
                      </Form.Group>
                    </Col>
                    <Col md={4}>
                      <Form.Group controlId="attendance-late">
                        <Form.Label>Đi trễ cho phép (phút)</Form.Label>
                        <Form.Control
                          type="number"
                          min={0}
                          max={120}
                          value={attendance.allowLateMinutes}
                          onChange={(event) =>
                            setAttendance((prev) => ({ ...prev, allowLateMinutes: Number(event.target.value) }))
                          }
                        />
                      </Form.Group>
                    </Col>
                  </Row>
                  <Stack gap={2} className="mt-3">
                    <Form.Check
                      type="switch"
                      id="mobile-checkin"
                      label="Cho phép check-in/out từ thiết bị di động"
                      checked={attendance.mobileCheckIn}
                      onChange={(event) =>
                        setAttendance((prev) => ({ ...prev, mobileCheckIn: event.target.checked }))
                      }
                    />
                    <Form.Check
                      type="switch"
                      id="require-photo"
                      label="Yêu cầu chụp ảnh khi chấm công"
                      checked={attendance.requirePhoto}
                      onChange={(event) =>
                        setAttendance((prev) => ({ ...prev, requirePhoto: event.target.checked }))
                      }
                    />
                  </Stack>
                </SettingsSection>
              </Tab.Pane>

              <Tab.Pane eventKey="notifications">
                <SettingsSection
                  title="Cài đặt thông báo"
                  description="Lựa chọn kênh thông báo cho sự kiện quan trọng."
                  onSave={() => handleSave('notifications')}
                >
                  <Stack gap={3}>
                    <Form.Check
                      type="switch"
                      id="notif-email"
                      label="Thông báo qua email"
                      checked={notifications.email}
                      onChange={(event) =>
                        setNotifications((prev) => ({ ...prev, email: event.target.checked }))
                      }
                    />
                    <Form.Check
                      type="switch"
                      id="notif-sms"
                      label="Thông báo qua SMS"
                      checked={notifications.sms}
                      onChange={(event) =>
                        setNotifications((prev) => ({ ...prev, sms: event.target.checked }))
                      }
                    />
                    <Form.Check
                      type="switch"
                      id="notif-inapp"
                      label="Thông báo trong ứng dụng"
                      checked={notifications.inApp}
                      onChange={(event) =>
                        setNotifications((prev) => ({ ...prev, inApp: event.target.checked }))
                      }
                    />
                  </Stack>
                </SettingsSection>
              </Tab.Pane>

              <Tab.Pane eventKey="security">
                <SettingsSection
                  title="Cài đặt bảo mật"
                  description="Tăng cường bảo mật tài khoản và phiên đăng nhập."
                  onSave={() => handleSave('security')}
                >
                  <Stack gap={3}>
                    <Form.Group controlId="security-password">
                      <Form.Label>Mật khẩu mới</Form.Label>
                      <Form.Control type="password" placeholder="Nhập mật khẩu mới" />
                    </Form.Group>
                    <Form.Group controlId="security-password-confirm">
                      <Form.Label>Xác nhận mật khẩu</Form.Label>
                      <Form.Control type="password" placeholder="Nhập lại mật khẩu" />
                    </Form.Group>
                    <Form.Check
                      type="switch"
                      id="security-2fa"
                      label="Yêu cầu xác thực 2 yếu tố (2FA)"
                      checked={security.twoFactor}
                      onChange={(event) =>
                        setSecurity((prev) => ({ ...prev, twoFactor: event.target.checked }))
                      }
                    />
                    <Form.Check
                      type="switch"
                      id="security-autologout"
                      label="Tự động đăng xuất sau 30 phút không hoạt động"
                      checked={security.autoLogout}
                      onChange={(event) =>
                        setSecurity((prev) => ({ ...prev, autoLogout: event.target.checked }))
                      }
                    />
                  </Stack>
                </SettingsSection>
              </Tab.Pane>
            </Tab.Content>
          </Col>
        </Row>
      </Tab.Container>
    </Page>
  );
}

interface SettingsSectionProps {
  title: string;
  description: string;
  children: ReactNode;
  onSave: () => void;
}

function SettingsSection({ title, description, children, onSave }: SettingsSectionProps) {
  return (
    <Card className="border-0 shadow-sm">
      <Card.Header className="bg-transparent border-0 pb-0">
        <h2 className="fs-5 fw-semibold mb-1">{title}</h2>
        <p className="text-secondary small mb-0">{description}</p>
      </Card.Header>
      <Card.Body className="pt-3">
        <Stack gap={3}>{children}</Stack>
      </Card.Body>
      <Card.Footer className="bg-transparent border-0 d-flex justify-content-end gap-2">
        <Button variant="outline-secondary" onClick={onSave}>
          Hủy
        </Button>
        <Button variant="primary" onClick={onSave}>
          Lưu thay đổi
        </Button>
      </Card.Footer>
    </Card>
  );
}

function translateSection(section: 'general' | 'attendance' | 'notifications' | 'security') {
  switch (section) {
    case 'general':
      return 'chung';
    case 'attendance':
      return 'chấm công';
    case 'notifications':
      return 'thông báo';
    case 'security':
      return 'bảo mật';
  }
}