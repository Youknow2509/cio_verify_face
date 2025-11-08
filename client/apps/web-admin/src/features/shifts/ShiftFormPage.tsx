import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Grid,
  Typography,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';

export const ShiftFormPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [formData, setFormData] = useState({
    name: '',
    start_time: '08:00',
    end_time: '17:00',
    break_start: '12:00',
    break_end: '13:00',
  });

  const handleChange = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [field]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    navigate('/shifts');
  };

  return (
    <Box>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate('/shifts')}
        sx={{ mb: 2 }}
      >
        Quay lại
      </Button>
      <Card>
        <CardContent>
          <Typography variant="h5" fontWeight="bold" mb={3}>
            {id ? 'Chỉnh sửa ca làm việc' : 'Thêm ca làm việc mới'}
          </Typography>
          <Box component="form" onSubmit={handleSubmit}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Tên ca làm việc"
                  required
                  value={formData.name}
                  onChange={handleChange('name')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Giờ bắt đầu"
                  type="time"
                  required
                  InputLabelProps={{ shrink: true }}
                  value={formData.start_time}
                  onChange={handleChange('start_time')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Giờ kết thúc"
                  type="time"
                  required
                  InputLabelProps={{ shrink: true }}
                  value={formData.end_time}
                  onChange={handleChange('end_time')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Giờ nghỉ trưa bắt đầu"
                  type="time"
                  InputLabelProps={{ shrink: true }}
                  value={formData.break_start}
                  onChange={handleChange('break_start')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Giờ nghỉ trưa kết thúc"
                  type="time"
                  InputLabelProps={{ shrink: true }}
                  value={formData.break_end}
                  onChange={handleChange('break_end')}
                />
              </Grid>
            </Grid>
            <Box mt={3} display="flex" gap={2}>
              <Button type="submit" variant="contained" startIcon={<Save />}>
                Lưu
              </Button>
              <Button variant="outlined" onClick={() => navigate('/shifts')}>
                Hủy
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};
