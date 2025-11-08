import { useState } from 'react';
import {
  Box,
  Card,
  TextField,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Chip,
  Grid,
} from '@mui/material';
import { Download } from '@mui/icons-material';
import type { AttendanceRecord } from '@face-attendance/types';

export const DailyReportPage: React.FC = () => {
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [records, setRecords] = useState<AttendanceRecord[]>([
    {
      id: '1',
      employee_id: '1',
      employee_name: 'Nguyễn Văn A',
      employee_code: 'NV001',
      device_id: '1',
      device_name: 'Thiết bị tầng 1',
      shift_name: 'Ca sáng',
      check_in_time: '2024-01-15T08:05:00Z',
      check_out_time: '2024-01-15T17:02:00Z',
      check_in_status: 'on_time',
      check_out_status: 'on_time',
      total_hours: 8.95,
      date: '2024-01-15',
      created_at: '2024-01-15',
      updated_at: '2024-01-15',
    },
  ]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'on_time':
        return 'success';
      case 'late':
        return 'error';
      case 'early':
        return 'warning';
      default:
        return 'default';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'on_time':
        return 'Đúng giờ';
      case 'late':
        return 'Trễ';
      case 'early':
        return 'Sớm';
      default:
        return '-';
    }
  };

  return (
    <Box>
      <Typography variant="h4" fontWeight="bold" mb={3}>
        Báo cáo Chấm công Hàng ngày
      </Typography>
      <Card sx={{ mb: 3, p: 2 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              label="Ngày"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12} md={8}>
            <Button variant="contained" startIcon={<Download />}>
              Xuất Excel
            </Button>
          </Grid>
        </Grid>
      </Card>
      <Card>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nhân viên</TableCell>
                <TableCell>Ca làm việc</TableCell>
                <TableCell>Giờ vào</TableCell>
                <TableCell>Giờ ra</TableCell>
                <TableCell>Tổng giờ</TableCell>
                <TableCell>Trạng thái vào</TableCell>
                <TableCell>Trạng thái ra</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {records.map((record) => (
                <TableRow key={record.id}>
                  <TableCell>
                    <Box>
                      <Typography fontWeight="bold">{record.employee_name}</Typography>
                      <Typography variant="caption" color="textSecondary">
                        {record.employee_code}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>{record.shift_name || '-'}</TableCell>
                  <TableCell>
                    {record.check_in_time
                      ? new Date(record.check_in_time).toLocaleTimeString('vi-VN')
                      : '-'}
                  </TableCell>
                  <TableCell>
                    {record.check_out_time
                      ? new Date(record.check_out_time).toLocaleTimeString('vi-VN')
                      : '-'}
                  </TableCell>
                  <TableCell>
                    {record.total_hours ? `${record.total_hours.toFixed(1)}h` : '-'}
                  </TableCell>
                  <TableCell>
                    {record.check_in_status && (
                      <Chip
                        label={getStatusText(record.check_in_status)}
                        color={getStatusColor(record.check_in_status)}
                        size="small"
                      />
                    )}
                  </TableCell>
                  <TableCell>
                    {record.check_out_status && (
                      <Chip
                        label={getStatusText(record.check_out_status)}
                        color={getStatusColor(record.check_out_status)}
                        size="small"
                      />
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Card>
    </Box>
  );
};
