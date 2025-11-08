import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Avatar,
  Chip,
  Typography,
  CircularProgress,
} from '@mui/material';
import { Add, Edit, Face, Delete } from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';
import type { Employee } from '@face-attendance/types';

export const EmployeeListPage: React.FC = () => {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      const response = await apiClient.get('/api/v1/users');
      // Mock data for demo
      setEmployees([
        {
          id: '1',
          email: 'nguyen.van.a@company.com',
          full_name: 'Nguyễn Văn A',
          employee_code: 'NV001',
          phone: '0123456789',
          role: 'employee',
          status: 'active',
          department: 'IT',
          position: 'Developer',
          face_data_count: 3,
          created_at: '2024-01-01',
          updated_at: '2024-01-01',
        },
      ]);
    } catch (error) {
      console.error('Failed to fetch employees:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredEmployees = employees.filter((emp) =>
    emp.full_name.toLowerCase().includes(search.toLowerCase()) ||
    emp.employee_code?.toLowerCase().includes(search.toLowerCase())
  );

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="bold">
          Quản lý Nhân viên
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => navigate('/employees/add')}
        >
          Thêm nhân viên
        </Button>
      </Box>
      <Card>
        <Box p={2}>
          <TextField
            fullWidth
            placeholder="Tìm kiếm theo tên hoặc mã nhân viên..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </Box>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nhân viên</TableCell>
                <TableCell>Mã NV</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Phòng ban</TableCell>
                <TableCell>Trạng thái</TableCell>
                <TableCell>Khuôn mặt</TableCell>
                <TableCell align="right">Thao tác</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredEmployees.map((employee) => (
                <TableRow key={employee.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center" gap={2}>
                      <Avatar src={employee.avatar_url}>
                        {employee.full_name[0]}
                      </Avatar>
                      <Typography>{employee.full_name}</Typography>
                    </Box>
                  </TableCell>
                  <TableCell>{employee.employee_code}</TableCell>
                  <TableCell>{employee.email}</TableCell>
                  <TableCell>{employee.department || '-'}</TableCell>
                  <TableCell>
                    <Chip
                      label={employee.status === 'active' ? 'Hoạt động' : 'Không hoạt động'}
                      color={employee.status === 'active' ? 'success' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip label={employee.face_data_count} size="small" />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/employees/${employee.id}/face-data`)}
                    >
                      <Face />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/employees/${employee.id}/edit`)}
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
