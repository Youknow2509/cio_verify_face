import { useEffect, useState } from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  CircularProgress,
} from '@mui/material';
import {
  People,
  CheckCircle,
  Warning,
  Devices,
} from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';
import type { DashboardStats } from '@face-attendance/types';

export const DashboardPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<DashboardStats>({
    total_employees: 0,
    attendance_today: 0,
    late_rate_this_month: 0,
    active_devices: 0,
  });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await apiClient.get('/api/v1/reports/daily');
        if (response.data) {
          setStats(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch stats:', error);
        // Keep default zeros - no mock data
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" mb={3} fontWeight="bold">
        Dashboard
      </Typography>
      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography color="textSecondary" variant="body2">
                    Tổng nhân viên
                  </Typography>
                  <Typography variant="h4" fontWeight="bold">
                    {stats.total_employees}
                  </Typography>
                </Box>
                <People sx={{ fontSize: 48, color: 'primary.main' }} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography color="textSecondary" variant="body2">
                    Chấm công hôm nay
                  </Typography>
                  <Typography variant="h4" fontWeight="bold">
                    {stats.attendance_today}
                  </Typography>
                </Box>
                <CheckCircle sx={{ fontSize: 48, color: 'success.main' }} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography color="textSecondary" variant="body2">
                    Tỷ lệ đi trễ
                  </Typography>
                  <Typography variant="h4" fontWeight="bold">
                    {(stats.late_rate_this_month * 100).toFixed(1)}%
                  </Typography>
                </Box>
                <Warning sx={{ fontSize: 48, color: 'warning.main' }} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography color="textSecondary" variant="body2">
                    Thiết bị hoạt động
                  </Typography>
                  <Typography variant="h4" fontWeight="bold">
                    {stats.active_devices}
                  </Typography>
                </Box>
                <Devices sx={{ fontSize: 48, color: 'info.main' }} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};
