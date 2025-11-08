import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Typography,
  Chip,
} from '@mui/material';
import { Add, Edit, Delete } from '@mui/icons-material';
import type { Shift } from '@face-attendance/types';

export const ShiftListPage: React.FC = () => {
  const navigate = useNavigate();
  const [shifts, setShifts] = useState<Shift[]>([
    {
      id: '1',
      company_id: '1',
      name: 'Ca sáng',
      start_time: '08:00',
      end_time: '17:00',
      break_start: '12:00',
      break_end: '13:00',
      valid_check_in_before: 15,
      valid_check_in_after: 30,
      valid_check_out_before: 30,
      valid_check_out_after: 15,
      days_of_week: [1, 2, 3, 4, 5],
      employee_count: 50,
      created_at: '2024-01-01',
      updated_at: '2024-01-01',
    },
  ]);

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="bold">
          Quản lý Ca làm việc
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => navigate('/shifts/add')}
        >
          Thêm ca làm việc
        </Button>
      </Box>
      <Card>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Tên ca</TableCell>
                <TableCell>Giờ làm việc</TableCell>
                <TableCell>Giờ nghỉ</TableCell>
                <TableCell>Số nhân viên</TableCell>
                <TableCell align="right">Thao tác</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {shifts.map((shift) => (
                <TableRow key={shift.id}>
                  <TableCell>
                    <Typography fontWeight="bold">{shift.name}</Typography>
                  </TableCell>
                  <TableCell>
                    {shift.start_time} - {shift.end_time}
                  </TableCell>
                  <TableCell>
                    {shift.break_start && shift.break_end
                      ? `${shift.break_start} - ${shift.break_end}`
                      : '-'}
                  </TableCell>
                  <TableCell>
                    <Chip label={shift.employee_count} size="small" />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/shifts/${shift.id}/edit`)}
                    >
                      <Edit />
                    </IconButton>
                    <IconButton size="small" color="error">
                      <Delete />
                    </IconButton>
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
