import { useCallback, useState, type FormEvent } from 'react';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Spinner from 'react-bootstrap/Spinner';
import Stack from 'react-bootstrap/Stack';
import { useNavigate } from 'react-router-dom';
import { useUi } from '@/app/providers/UiProvider';
import { Page } from '@/ui/Page';
import { changePassword } from '@/services/api/account';
import { clearAuthToken, HttpError } from '@/services/http';
import type { ChangePasswordPayload } from '@/types';

interface FormState {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
}

type FormErrors = Partial<Record<keyof FormState, string>>;

type VisibilityState = Record<keyof FormState, boolean>;

const defaultFormState: FormState = {
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
};

const defaultErrors: FormErrors = {};

export default function ChangePasswordPage() {
  const navigate = useNavigate();
  const { showToast } = useUi();
  const [formState, setFormState] = useState<FormState>(defaultFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>(defaultErrors);
  const [visible, setVisible] = useState<VisibilityState>({
    oldPassword: false,
    newPassword: false,
    confirmPassword: false,
  });
  const [submitting, setSubmitting] = useState(false);

  const toggleVisibility = (field: keyof FormState) => {
    setVisible((prev) => ({
      ...prev,
      [field]: !prev[field],
    }));
  };

  const handleChange = <K extends keyof FormState>(key: K, value: FormState[K]) => {
    setFormState((prev) => ({
      ...prev,
      [key]: value,
    }));
  };

  const validate = useCallback(() => {
    const errors: FormErrors = {};

    if (formState.oldPassword.length === 0) {
      errors.oldPassword = 'Vui lòng nhập mật khẩu hiện tại';
    }

    if (formState.newPassword.length < 8) {
      errors.newPassword = 'Mật khẩu mới tối thiểu 8 ký tự';
    }

    if (formState.newPassword === formState.oldPassword && formState.newPassword.length > 0) {
      errors.newPassword = 'Mật khẩu mới phải khác mật khẩu hiện tại';
    }

    if (formState.confirmPassword !== formState.newPassword) {
      errors.confirmPassword = 'Xác nhận mật khẩu không khớp';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  }, [formState]);

  const handleSubmit = useCallback(
    async (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      if (!validate()) {
        return;
      }

      try {
        setSubmitting(true);
        const payload: ChangePasswordPayload = {
          oldPassword: formState.oldPassword,
          newPassword: formState.newPassword,
        };
        await changePassword(payload);
        showToast({
          variant: 'success',
          title: 'Thành công',
          message: 'Đã đổi mật khẩu thành công',
        });
        setFormState(defaultFormState);
        setFormErrors(defaultErrors);
      } catch (error) {
        if (error instanceof HttpError && error.status === 401) {
          clearAuthToken();
          navigate('/login', { replace: true });
          return;
        }

        showToast({
          variant: 'danger',
          title: 'Lỗi',
          message: error instanceof Error ? error.message : 'Không thể đổi mật khẩu',
        });
      } finally {
        setSubmitting(false);
      }
    },
    [formState, navigate, showToast, validate]
  );

  const renderPasswordField = (
    field: keyof FormState,
    label: string,
    autoComplete: string
  ) => {
    const isVisible = visible[field];
    const inputType = isVisible ? 'text' : 'password';
    const error = formErrors[field];
    return (
      <Form.Group controlId={`change-password-${field}`}>
        <Form.Label>{label}</Form.Label>
  <InputGroup hasValidation>
          <Form.Control
            type={inputType}
            value={formState[field]}
            onChange={(event) => handleChange(field, event.target.value)}
            isInvalid={Boolean(error)}
            autoComplete={autoComplete}
            required
          />
          <Button
            variant="outline-secondary"
            onClick={() => toggleVisibility(field)}
            type="button"
            aria-label={isVisible ? 'Ẩn mật khẩu' : 'Hiển thị mật khẩu'}
          >
            <i className={`bi ${isVisible ? 'bi-eye-slash' : 'bi-eye'}`} aria-hidden />
          </Button>
          <Form.Control.Feedback type="invalid">{error}</Form.Control.Feedback>
        </InputGroup>
        {field === 'newPassword' && (
          <Form.Text className="text-secondary">
            Mật khẩu nên có tối thiểu 8 ký tự, bao gồm chữ hoa, chữ thường và số.
          </Form.Text>
        )}
      </Form.Group>
    );
  };

  return (
    <Page
      title="Đổi mật khẩu"
      subtitle="Cập nhật mật khẩu đăng nhập của bạn"
      breadcrumb={[
        { label: 'Trang chủ', path: '/dashboard' },
        { label: 'Tài khoản', path: '/account' },
        { label: 'Đổi mật khẩu' },
      ]}
    >
      <Card className="border-0 shadow-sm">
        <Card.Body>
          <Form noValidate onSubmit={handleSubmit} className="gy-4">
            <Stack gap={3}>
              {renderPasswordField('oldPassword', 'Mật khẩu hiện tại', 'current-password')}
              {renderPasswordField('newPassword', 'Mật khẩu mới', 'new-password')}
              {renderPasswordField('confirmPassword', 'Xác nhận mật khẩu mới', 'new-password')}
            </Stack>
            <div className="d-flex justify-content-end gap-2 mt-4">
              <Button
                variant="outline-secondary"
                type="button"
                onClick={() => {
                  setFormState(defaultFormState);
                  setFormErrors(defaultErrors);
                }}
                disabled={submitting}
              >
                Làm lại
              </Button>
              <Button type="submit" variant="primary" disabled={submitting}>
                {submitting ? (
                  <span className="d-inline-flex align-items-center gap-2">
                    <Spinner animation="border" size="sm" role="status" aria-hidden />
                    Đang cập nhật...
                  </span>
                ) : (
                  'Đổi mật khẩu'
                )}
              </Button>
            </div>
          </Form>
        </Card.Body>
      </Card>
    </Page>
  );
}
