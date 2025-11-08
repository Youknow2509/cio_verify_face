import { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Grid,
  Typography,
  Switch,
  FormControlLabel,
  Tabs,
  Tab,
} from '@mui/material';
import { Save } from '@mui/icons-material';

export const SettingsPage: React.FC = () => {
  const [tab, setTab] = useState(0);
  const [companyInfo, setCompanyInfo] = useState({
    name: 'Công ty TNHH ABC',
    address: 'Hà Nội, Việt Nam',
    phone: '0123456789',
    email: 'contact@company.com',
  });

  const [settings, setSettings] = useState({
    valid_time_before_shift: 15,
    valid_time_after_shift: 30,
    allow_offline_attendance: false,
    daily_email_report: true,
    device_offline_alert: true,
  });

  return (
    <Box>
      <Typography variant="h4" fontWeight="bold" mb={3}>
        Cài đặt
      </Typography>
      <Card>
        <Tabs value={tab} onChange={(_, v) => setTab(v)}>
          <Tab label="Thông tin công ty" />
          <Tab label="Cài đặt chấm công" />
          <Tab label="Thông báo" />
        </Tabs>
        <CardContent sx={{ p: 3 }}>
          {tab === 0 && (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Tên công ty"
                  value={companyInfo.name}
                  onChange={(e) =>
                    setCompanyInfo({ ...companyInfo, name: e.target.value })
                  }
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Địa chỉ"
                  value={companyInfo.address}
                  onChange={(e) =>
                    setCompanyInfo({ ...companyInfo, address: e.target.value })
                  }
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Số điện thoại"
                  value={companyInfo.phone}
                  onChange={(e) =>
                    setCompanyInfo({ ...companyInfo, phone: e.target.value })
                  }
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Email"
                  value={companyInfo.email}
                  onChange={(e) =>
                    setCompanyInfo({ ...companyInfo, email: e.target.value })
                  }
                />
              </Grid>
            </Grid>
          )}
          {tab === 1 && (
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Thời gian cho phép chấm công trước ca (phút)"
                  type="number"
                  value={settings.valid_time_before_shift}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      valid_time_before_shift: parseInt(e.target.value),
                    })
                  }
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Thời gian cho phép chấm công sau ca (phút)"
                  type="number"
                  value={settings.valid_time_after_shift}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      valid_time_after_shift: parseInt(e.target.value),
                    })
                  }
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.allow_offline_attendance}
                      onChange={(e) =>
                        setSettings({
                          ...settings,
                          allow_offline_attendance: e.target.checked,
                        })
                      }
                    />
                  }
                  label="Cho phép chấm công offline"
                />
              </Grid>
            </Grid>
          )}
          {tab === 2 && (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.daily_email_report}
                      onChange={(e) =>
                        setSettings({
                          ...settings,
                          daily_email_report: e.target.checked,
                        })
                      }
                    />
                  }
                  label="Gửi báo cáo email hàng ngày"
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.device_offline_alert}
                      onChange={(e) =>
                        setSettings({
                          ...settings,
                          device_offline_alert: e.target.checked,
                        })
                      }
                    />
                  }
                  label="Cảnh báo khi thiết bị offline"
                />
              </Grid>
            </Grid>
          )}
          <Box mt={3}>
            <Button variant="contained" startIcon={<Save />}>
              Lưu thay đổi
            </Button>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};
