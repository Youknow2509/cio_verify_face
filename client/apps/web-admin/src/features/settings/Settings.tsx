// src/features/settings/Settings.tsx

import { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/Card/Card';
import styles from './Settings.module.scss';

export default function Settings() {
  const [activeTab, setActiveTab] = useState('general');

  return (
    <div className={styles.settings}>
      <div className={styles.header}>
        <h1 className={styles.title}>Cài đặt hệ thống</h1>
        <p className={styles.subtitle}>Cấu hình thông số cho hệ thống chấm công</p>
      </div>

      <div className={styles.settingsLayout}>
        <div className={styles.sidebar}>
          <nav className={styles.nav}>
            <button
              className={`${styles.navItem} ${activeTab === 'general' ? styles.active : ''}`}
              onClick={() => setActiveTab('general')}
            >
              <SettingsIcon />
              Cài đặt chung
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'attendance' ? styles.active : ''}`}
              onClick={() => setActiveTab('attendance')}
            >
              <ClockIcon />
              Chấm công
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'notification' ? styles.active : ''}`}
              onClick={() => setActiveTab('notification')}
            >
              <BellIcon />
              Thông báo
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'security' ? styles.active : ''}`}
              onClick={() => setActiveTab('security')}
            >
              <LockIcon />
              Bảo mật
            </button>
          </nav>
        </div>

        <div className={styles.content}>
          {activeTab === 'general' && (
            <Card>
              <CardHeader>
                <CardTitle>Cài đặt chung</CardTitle>
              </CardHeader>
              <CardContent>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Tên công ty</label>
                  <input type="text" className={styles.input} defaultValue="Công ty ABC" />
                </div>
                
                <div className={styles.formGroup}>
                  <label className={styles.label}>Múi giờ</label>
                  <select className={styles.select}>
                    <option>GMT+7 (Hà Nội)</option>
                  </select>
                </div>
                
                <div className={styles.formGroup}>
                  <label className={styles.label}>Ngôn ngữ</label>
                  <select className={styles.select}>
                    <option>Tiếng Việt</option>
                    <option>English</option>
                  </select>
                </div>
                
                <button className={styles.saveButton}>Lưu thay đổi</button>
              </CardContent>
            </Card>
          )}

          {activeTab === 'attendance' && (
            <Card>
              <CardHeader>
                <CardTitle>Cài đặt chấm công</CardTitle>
              </CardHeader>
              <CardContent>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Giờ bắt đầu làm việc</label>
                  <input type="time" className={styles.input} defaultValue="08:00" />
                </div>
                
                <div className={styles.formGroup}>
                  <label className={styles.label}>Giờ kết thúc làm việc</label>
                  <input type="time" className={styles.input} defaultValue="17:00" />
                </div>
                
                <div className={styles.formGroup}>
                  <label className={styles.label}>Thời gian cho phép đi trễ (phút)</label>
                  <input type="number" className={styles.input} defaultValue="15" />
                </div>
                
                <div className={styles.checkboxGroup}>
                  <label className={styles.checkbox}>
                    <input type="checkbox" defaultChecked />
                    <span>Cho phép check-in/out từ thiết bị di động</span>
                  </label>
                  <label className={styles.checkbox}>
                    <input type="checkbox" defaultChecked />
                    <span>Yêu cầu chụp ảnh khi chấm công</span>
                  </label>
                </div>
                
                <button className={styles.saveButton}>Lưu thay đổi</button>
              </CardContent>
            </Card>
          )}

          {activeTab === 'notification' && (
            <Card>
              <CardHeader>
                <CardTitle>Cài đặt thông báo</CardTitle>
              </CardHeader>
              <CardContent>
                <div className={styles.checkboxGroup}>
                  <label className={styles.checkbox}>
                    <input type="checkbox" defaultChecked />
                    <span>Thông báo qua email</span>
                  </label>
                  <label className={styles.checkbox}>
                    <input type="checkbox" />
                    <span>Thông báo qua SMS</span>
                  </label>
                  <label className={styles.checkbox}>
                    <input type="checkbox" defaultChecked />
                    <span>Thông báo trong ứng dụng</span>
                  </label>
                </div>
                
                <button className={styles.saveButton}>Lưu thay đổi</button>
              </CardContent>
            </Card>
          )}

          {activeTab === 'security' && (
            <Card>
              <CardHeader>
                <CardTitle>Cài đặt bảo mật</CardTitle>
              </CardHeader>
              <CardContent>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Mật khẩu mới</label>
                  <input type="password" className={styles.input} />
                </div>
                
                <div className={styles.formGroup}>
                  <label className={styles.label}>Xác nhận mật khẩu</label>
                  <input type="password" className={styles.input} />
                </div>
                
                <div className={styles.checkboxGroup}>
                  <label className={styles.checkbox}>
                    <input type="checkbox" defaultChecked />
                    <span>Yêu cầu xác thực 2 yếu tố (2FA)</span>
                  </label>
                  <label className={styles.checkbox}>
                    <input type="checkbox" />
                    <span>Tự động đăng xuất sau 30 phút không hoạt động</span>
                  </label>
                </div>
                
                <button className={styles.saveButton}>Lưu thay đổi</button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}

const SettingsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>
  </svg>
);

const ClockIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z"/>
  </svg>
);

const BellIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.89 2 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.63 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
  </svg>
);

const LockIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
  </svg>
);