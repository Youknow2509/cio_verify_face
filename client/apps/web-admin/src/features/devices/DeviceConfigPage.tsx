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
  Switch,
  FormControlLabel,
  Slider,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';

export const DeviceConfigPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [config, setConfig] = useState({
    allow_check_in: true,
    allow_check_out: true,
    timeout: 30,
    recognition_threshold: 0.8,
    sound_enabled: true,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    navigate('/devices');
  };

  return (
    <Box>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate('/devices')}
        sx={{ mb: 2 }}
      >
        Quay lại
      </Button>
      <Card>
        <CardContent>
          <Typography variant="h5" fontWeight="bold" mb={3}>
            Cấu hình thiết bị
          </Typography>
          <Box component="form" onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={config.allow_check_in}
                      onChange={(e) =>
                        setConfig({ ...config, allow_check_in: e.target.checked })
                      }
                    />
                  }
                  label="Cho phép check-in"
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={config.allow_check_out}
                      onChange={(e) =>
                        setConfig({ ...config, allow_check_out: e.target.checked })
                      }
                    />
                  }
                  label="Cho phép check-out"
                />
              </Grid>
              <Grid item xs={12}>
                <Typography gutterBottom>Timeout (giây)</Typography>
                <Slider
                  value={config.timeout}
                  onChange={(_, value) =>
                    setConfig({ ...config, timeout: value as number })
                  }
                  min={10}
                  max={60}
                  marks
                  valueLabelDisplay="on"
                />
              </Grid>
              <Grid item xs={12}>
                <Typography gutterBottom>Độ nhạy nhận diện</Typography>
                <Slider
                  value={config.recognition_threshold}
                  onChange={(_, value) =>
                    setConfig({ ...config, recognition_threshold: value as number })
                  }
                  min={0.5}
                  max={1}
                  step={0.05}
                  marks
                  valueLabelDisplay="on"
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={config.sound_enabled}
                      onChange={(e) =>
                        setConfig({ ...config, sound_enabled: e.target.checked })
                      }
                    />
                  }
                  label="Bật âm thanh"
                />
              </Grid>
            </Grid>
            <Box mt={3} display="flex" gap={2}>
              <Button type="submit" variant="contained" startIcon={<Save />}>
                Lưu cấu hình
              </Button>
              <Button variant="outlined" onClick={() => navigate('/devices')}>
                Hủy
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};
