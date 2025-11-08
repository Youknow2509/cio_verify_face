import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Grid,
  Typography,
  MenuItem,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';

export const EmployeeFormPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEdit = !!id;

  const [formData, setFormData] = useState({
    full_name: '',
    employee_code: '',
    email: '',
    phone: '',
    department: '',
    position: '',
    hire_date: '',
  });

  const handleChange = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [field]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // API call here
    navigate('/employees');
  };

  return (
    <Box>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate('/employees')}
        sx={{ mb: 2 }}
      >
        Quay lại
      </Button>
      <Card>
        <CardContent>
          <Typography variant="h5" fontWeight="bold" mb={3}>
            {isEdit ? 'Chỉnh sửa nhân viên' : 'Thêm nhân viên mới'}
          </Typography>
          <Box component="form" onSubmit={handleSubmit}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Họ và tên"
                  required
                  value={formData.full_name}
                  onChange={handleChange('full_name')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Mã nhân viên"
                  required
                  value={formData.employee_code}
                  onChange={handleChange('employee_code')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleChange('email')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Số điện thoại"
                  value={formData.phone}
                  onChange={handleChange('phone')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Phòng ban"
                  value={formData.department}
                  onChange={handleChange('department')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Chức vụ"
                  value={formData.position}
                  onChange={handleChange('position')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Ngày vào làm"
                  type="date"
                  InputLabelProps={{ shrink: true }}
                  value={formData.hire_date}
                  onChange={handleChange('hire_date')}
                />
              </Grid>
            </Grid>
            <Box mt={3} display="flex" gap={2}>
              <Button
                type="submit"
                variant="contained"
                startIcon={<Save />}
              >
                Lưu
              </Button>
              <Button
                variant="outlined"
                onClick={() => navigate('/employees')}
              >
                Hủy
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};
