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
  Grid,
} from '@mui/material';
import { Download } from '@mui/icons-material';

export const SummaryReportPage: React.FC = () => {
  const [startDate, setStartDate] = useState(new Date().toISOString().split('T')[0]);
  const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0]);

  const summaryData = [
    {
      employee_name: 'Nguyễn Văn A',
      total_days: 20,
      total_hours: 160,
      late_count: 2,
      early_leave_count: 1,
      compliance_rate: 0.95,
    },
  ];

  return (
    <Box>
      <Typography variant="h4" fontWeight="bold" mb={3}>
        Báo cáo Tổng hợp
      </Typography>
      <Card sx={{ mb: 3, p: 2 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              label="Từ ngày"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              label="Đến ngày"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12} md={4}>
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
                <TableCell align="right">Tổng ngày</TableCell>
                <TableCell align="right">Tổng giờ</TableCell>
                <TableCell align="right">Số lần trễ</TableCell>
                <TableCell align="right">Về sớm</TableCell>
                <TableCell align="right">Tỷ lệ tuân thủ</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {summaryData.map((data, index) => (
                <TableRow key={index}>
                  <TableCell>{data.employee_name}</TableCell>
                  <TableCell align="right">{data.total_days}</TableCell>
                  <TableCell align="right">{data.total_hours}h</TableCell>
                  <TableCell align="right">{data.late_count}</TableCell>
                  <TableCell align="right">{data.early_leave_count}</TableCell>
                  <TableCell align="right">
                    {(data.compliance_rate * 100).toFixed(1)}%
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
