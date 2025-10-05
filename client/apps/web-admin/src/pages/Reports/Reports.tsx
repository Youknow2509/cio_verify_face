// src/pages/Reports/Reports.tsx

import { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../../components/Card/Card';
import { Toolbar, ToolbarSection } from '../../components/Toolbar/Toolbar';
import styles from './Reports.module.scss';

export default function Reports() {
  const [reportType, setReportType] = useState('daily');
  const [dateRange, setDateRange] = useState({ start: '', end: '' });

  const handleExport = (format: string) => {
    console.log(`Exporting ${reportType} report as ${format}`);
  };

  return (
    <div className={styles.reports}>
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>Báo cáo</h1>
          <p className={styles.subtitle}>Tổng hợp và xuất báo cáo chấm công</p>
        </div>
      </div>

      <Toolbar>
        <ToolbarSection>
          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>Loại báo cáo:</label>
            <select
              value={reportType}
              onChange={(e) => setReportType(e.target.value)}
              className={styles.reportSelect}
            >
              <option value="daily">Báo cáo ngày</option>
              <option value="weekly">Báo cáo tuần</option>
              <option value="monthly">Báo cáo tháng</option>
              <option value="custom">Tùy chỉnh</option>
            </select>
          </div>
          
          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>Từ ngày:</label>
            <input
              type="date"
              value={dateRange.start}
              onChange={(e) => setDateRange({ ...dateRange, start: e.target.value })}
              className={styles.dateInput}
            />
          </div>
          
          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>Đến ngày:</label>
            <input
              type="date"
              value={dateRange.end}
              onChange={(e) => setDateRange({ ...dateRange, end: e.target.value })}
              className={styles.dateInput}
            />
          </div>
        </ToolbarSection>
        
        <ToolbarSection align="right">
          <button className={styles.exportButton} onClick={() => handleExport('excel')}>
            <ExcelIcon />
            Xuất Excel
          </button>
          <button className={styles.exportButton} onClick={() => handleExport('pdf')}>
            <PdfIcon />
            Xuất PDF
          </button>
        </ToolbarSection>
      </Toolbar>

      <div className={styles.reportGrid}>
        <Card>
          <CardHeader>
            <CardTitle>Tổng quan chấm công</CardTitle>
          </CardHeader>
          <CardContent>
            <div className={styles.statGrid}>
              <div className={styles.statItem}>
                <div className={styles.statLabel}>Tổng nhân viên</div>
                <div className={styles.statValue}>248</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statLabel}>Đi làm đúng giờ</div>
                <div className={styles.statValue}>215</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statLabel}>Đi trễ</div>
                <div className={styles.statValue}>18</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statLabel}>Vắng mặt</div>
                <div className={styles.statValue}>15</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Thống kê theo phòng ban</CardTitle>
          </CardHeader>
          <CardContent>
            <div className={styles.departmentList}>
              {['IT', 'HR', 'Sales', 'Finance', 'Marketing'].map(dept => (
                <div key={dept} className={styles.departmentItem}>
                  <span className={styles.departmentName}>{dept}</span>
                  <span className={styles.departmentCount}>42 nhân viên</span>
                  <div className={styles.departmentBar}>
                    <div className={styles.departmentProgress} style={{ width: '85%' }}></div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

const ExcelIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M14 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V8l-6-6zM6 20V4h7v5h5v11H6z"/>
  </svg>
);

const PdfIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M20 2H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-8.5 7.5c0 .83-.67 1.5-1.5 1.5H9v2H7.5V7H10c.83 0 1.5.67 1.5 1.5v1zm5 2c0 .83-.67 1.5-1.5 1.5h-2.5V7H15c.83 0 1.5.67 1.5 1.5v3zm4-3H19v1h1.5V11H19v2h-1.5V7h3v1.5zM9 9.5h1v-1H9v1zM4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm10 5.5h1v-3h-1v3z"/>
  </svg>
);