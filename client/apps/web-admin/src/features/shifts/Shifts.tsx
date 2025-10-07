// src/features/shifts/Shifts.tsx

import { useState } from 'react';
import { Card, CardContent } from '@/components/Card/Card';
import { Badge } from '@/components/Badge/Badge';
import { Toolbar, ToolbarSection, SearchBox } from '@/components/Toolbar/Toolbar';
import styles from './Shifts.module.scss';

interface Shift {
  id: string;
  name: string;
  startTime: string;
  endTime: string;
  workHours: number;
  breakTime: number;
  isActive: boolean;
}

export default function Shifts() {
  const [searchQuery, setSearchQuery] = useState('');
  const [shifts] = useState<Shift[]>([
    {
      id: '1',
      name: 'Ca sáng',
      startTime: '08:00',
      endTime: '12:00',
      workHours: 4,
      breakTime: 0,
      isActive: true
    },
    {
      id: '2',
      name: 'Ca chiều',
      startTime: '13:00',
      endTime: '17:00',
      workHours: 4,
      breakTime: 0,
      isActive: true
    },
    {
      id: '3',
      name: 'Ca hành chính',
      startTime: '08:00',
      endTime: '17:00',
      workHours: 8,
      breakTime: 1,
      isActive: true
    },
    {
      id: '4',
      name: 'Ca tối',
      startTime: '18:00',
      endTime: '22:00',
      workHours: 4,
      breakTime: 0,
      isActive: false
    }
  ]);

  const filteredShifts = shifts.filter(shift =>
    shift.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className={styles.shifts}>
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>Quản lý ca làm việc</h1>
          <p className={styles.subtitle}>Thiết lập và quản lý các ca làm việc</p>
        </div>
      </div>

      <Toolbar>
        <ToolbarSection>
          <SearchBox
            value={searchQuery}
            onChange={setSearchQuery}
            placeholder="Tìm kiếm ca làm việc..."
          />
        </ToolbarSection>
        
        <ToolbarSection align="right">
          <button className={styles.addButton}>
            <AddIcon />
            Thêm ca làm việc
          </button>
        </ToolbarSection>
      </Toolbar>

      <div className={styles.shiftGrid}>
        {filteredShifts.map(shift => (
          <Card key={shift.id} className={styles.shiftCard}>
            <CardContent>
              <div className={styles.shiftHeader}>
                <h3 className={styles.shiftName}>{shift.name}</h3>
                <Badge variant={shift.isActive ? 'success' : 'neutral'}>
                  {shift.isActive ? 'Đang dùng' : 'Tạm ngưng'}
                </Badge>
              </div>
              
              <div className={styles.shiftDetails}>
                <div className={styles.shiftTime}>
                  <ClockIcon />
                  <span>{shift.startTime} - {shift.endTime}</span>
                </div>
                
                <div className={styles.shiftInfo}>
                  <div className={styles.infoItem}>
                    <span className={styles.infoLabel}>Giờ làm:</span>
                    <span className={styles.infoValue}>{shift.workHours}h</span>
                  </div>
                  <div className={styles.infoItem}>
                    <span className={styles.infoLabel}>Nghỉ giữa ca:</span>
                    <span className={styles.infoValue}>{shift.breakTime}h</span>
                  </div>
                </div>
              </div>
              
              <div className={styles.shiftActions}>
                <button className={styles.editButton}>Sửa</button>
                <button className={styles.deleteButton}>Xóa</button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}

const AddIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
  </svg>
);

const ClockIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
    <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z"/>
  </svg>
);