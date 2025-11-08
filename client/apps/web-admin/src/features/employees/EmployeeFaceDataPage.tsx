import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Button,
  Typography,
  Grid,
  Paper,
  IconButton,
} from '@mui/material';
import { ArrowBack, CloudUpload, Delete } from '@mui/icons-material';

export const EmployeeFaceDataPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [faceImages, setFaceImages] = useState<string[]>([]);

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      // Handle file upload
      console.log('Files to upload:', files);
    }
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
            Quản lý dữ liệu khuôn mặt
          </Typography>
          <Paper
            sx={{
              p: 4,
              border: '2px dashed',
              borderColor: 'primary.main',
              textAlign: 'center',
              cursor: 'pointer',
              mb: 3,
            }}
            component="label"
          >
            <input
              type="file"
              accept="image/*"
              multiple
              hidden
              onChange={handleFileUpload}
            />
            <CloudUpload sx={{ fontSize: 48, color: 'primary.main', mb: 1 }} />
            <Typography variant="h6">
              Kéo thả ảnh vào đây hoặc click để chọn
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Chọn nhiều ảnh khuôn mặt từ các góc độ khác nhau
            </Typography>
          </Paper>
          <Typography variant="h6" mb={2}>
            Ảnh đã tải lên ({faceImages.length})
          </Typography>
          <Grid container spacing={2}>
            {faceImages.map((image, index) => (
              <Grid item xs={6} sm={4} md={3} key={index}>
                <Paper sx={{ position: 'relative', p: 1 }}>
                  <img
                    src={image}
                    alt={`Face ${index + 1}`}
                    style={{ width: '100%', borderRadius: 4 }}
                  />
                  <IconButton
                    size="small"
                    color="error"
                    sx={{ position: 'absolute', top: 4, right: 4 }}
                  >
                    <Delete />
                  </IconButton>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </CardContent>
      </Card>
    </Box>
  );
};
