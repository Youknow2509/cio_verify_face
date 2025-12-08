import { useEffect, useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    CircularProgress,
    Chip,
    Pagination,
} from '@mui/material';
import { format } from 'date-fns';
import { shiftApi } from '@/services/api';

export const ShiftsPage: React.FC = () => {
    const [loading, setLoading] = useState(true);
    const [shifts, setShifts] = useState<any[]>([]);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const size = 12;

    useEffect(() => {
        loadShifts();
    }, [page]);

    const loadShifts = async () => {
        try {
            setLoading(true);
            const response: any = await shiftApi.getEmployeeShifts({
                page,
                size,
            });
            setShifts(response.data?.shifts || []);
            setTotal(response.data?.total || 0);
        } catch (error) {
            console.error('Failed to load shifts:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    return (
        <Box>
            <Typography variant="h4" fontWeight="700" mb={3}>
                Ca làm việc của tôi
            </Typography>

            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <Typography variant="h6" fontWeight="600" mb={1}>
                        Tổng quan
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Bạn có {total} ca làm việc được gán
                    </Typography>
                </CardContent>
            </Card>

            {shifts.length > 0 ? (
                <>
                    <Grid container spacing={3}>
                        {shifts.map((shift, index) => (
                            <Grid item xs={12} sm={6} md={4} lg={3} key={index}>
                                <Card
                                    sx={{
                                        height: '100%',
                                        background: shift.is_active
                                            ? 'linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%)'
                                            : 'rgba(71, 85, 105, 0.1)',
                                        border: shift.is_active
                                            ? '1px solid rgba(59, 130, 246, 0.3)'
                                            : '1px solid rgba(71, 85, 105, 0.3)',
                                        transition: 'all 0.3s ease',
                                        '&:hover': {
                                            transform: 'translateY(-4px)',
                                            boxShadow: '0 10px 30px rgba(0, 0, 0, 0.3)',
                                        },
                                    }}
                                >
                                    <CardContent>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                                            <Typography variant="h6" fontWeight="600">
                                                {shift.shift_name}
                                            </Typography>
                                            <Chip
                                                label={shift.is_active ? 'Hoạt động' : 'Không hoạt động'}
                                                color={shift.is_active ? 'success' : 'default'}
                                                size="small"
                                            />
                                        </Box>

                                        <Box sx={{ mb: 2 }}>
                                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                                Giờ làm việc
                                            </Typography>
                                            <Typography variant="body1" fontWeight="600">
                                                {shift.shift_start} - {shift.shift_end}
                                            </Typography>
                                        </Box>

                                        <Box sx={{ mb: 1 }}>
                                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                                Hiệu lực từ
                                            </Typography>
                                            <Typography variant="body2">
                                                {shift.effective_from !== '0001-01-01T00:00:00Z'
                                                    ? format(new Date(shift.effective_from), 'dd/MM/yyyy')
                                                    : 'Không giới hạn'}
                                            </Typography>
                                        </Box>

                                        {shift.effective_to && shift.effective_to !== '0001-01-01T00:00:00Z' && (
                                            <Box>
                                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                                    Hiệu lực đến
                                                </Typography>
                                                <Typography variant="body2">
                                                    {format(new Date(shift.effective_to), 'dd/MM/yyyy')}
                                                </Typography>
                                            </Box>
                                        )}
                                    </CardContent>
                                </Card>
                            </Grid>
                        ))}
                    </Grid>

                    {Math.ceil(total / size) > 1 && (
                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                            <Pagination
                                count={Math.ceil(total / size)}
                                page={page}
                                onChange={(e, value) => setPage(value)}
                                color="primary"
                            />
                        </Box>
                    )}
                </>
            ) : (
                <Card>
                    <CardContent>
                        <Typography variant="body1" color="text.secondary" textAlign="center" py={4}>
                            Bạn chưa được gán ca làm việc nào
                        </Typography>
                    </CardContent>
                </Card>
            )}
        </Box>
    );
};
