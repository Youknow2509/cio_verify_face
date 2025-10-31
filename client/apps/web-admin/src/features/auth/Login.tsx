import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import { Page } from '@/ui/Page';
import styles from './Login.module.scss';
import { login as mockLogin } from '@/services/mock/auth';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await mockLogin(email, password);
      if (res.error) {
        setError(res.error);
        return;
      }

      // on success redirect
      navigate('/dashboard', { replace: true });
    } catch (err) {
      setError('Lỗi kết nối');
    } finally {
      setLoading(false);
    }
  };

  
  return (
    <Page title="Đăng nhập">
      <div className={styles.loginWrap}>
        <div className={styles.loginCard}>
          <div className={styles.title}>Đăng nhập vào hệ thống</div>
          <Form onSubmit={handleSubmit}>
            <Form.Group className={styles.formGroup} controlId="email">
              <Form.Label>Email</Form.Label>
              <Form.Control
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@company.com"
                required
              />
            </Form.Group>

            <Form.Group className={styles.formGroup} controlId="password">
              <Form.Label>Mật khẩu</Form.Label>
              <Form.Control
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Mật khẩu"
                required
              />
            </Form.Group>

            {error && <div className={styles.error}>{error}</div>}

            <div className={styles.actions}>
              <Button variant="secondary" size="sm" onClick={() => (setEmail('admin@company.com'), setPassword('password'))}>Demo</Button>
              <Button type="submit" variant="primary" size="sm" disabled={loading}>{loading ? 'Đang...' : 'Đăng nhập'}</Button>
            </div>
          </Form>
        </div>
      </div>
    </Page>
  );
}
