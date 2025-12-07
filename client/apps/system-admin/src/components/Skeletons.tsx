import { Box, Skeleton, Grid, Card, CardContent } from '@mui/material';

export const DashboardSkeleton: React.FC = () => (
    <Box>
        <Skeleton variant="text" width={200} height={40} sx={{ mb: 1 }} />
        <Skeleton variant="text" width={300} height={24} sx={{ mb: 4 }} />

        <Grid container spacing={3} sx={{ mb: 4 }}>
            {[1, 2, 3, 4].map((i) => (
                <Grid item xs={12} sm={6} lg={3} key={i}>
                    <Card>
                        <CardContent>
                            <Skeleton variant="text" width="60%" />
                            <Skeleton variant="text" width="40%" height={48} />
                            <Skeleton variant="rounded" width={60} height={22} sx={{ mt: 1 }} />
                        </CardContent>
                    </Card>
                </Grid>
            ))}
        </Grid>

        <Grid container spacing={3}>
            <Grid item xs={12} lg={8}>
                <Card>
                    <CardContent>
                        <Skeleton variant="text" width={150} height={32} sx={{ mb: 2 }} />
                        <Skeleton variant="rectangular" height={250} />
                    </CardContent>
                </Card>
            </Grid>
            <Grid item xs={12} lg={4}>
                <Card>
                    <CardContent>
                        <Skeleton variant="text" width={120} height={32} sx={{ mb: 2 }} />
                        <Skeleton variant="rectangular" height={200} />
                    </CardContent>
                </Card>
            </Grid>
        </Grid>
    </Box>
);

export const TableSkeleton: React.FC<{ rows?: number }> = ({ rows = 5 }) => (
    <Card>
        <Box sx={{ p: 2 }}>
            <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
                <Skeleton variant="rounded" width={200} height={40} />
                <Skeleton variant="rounded" width={150} height={40} />
            </Box>

            <Skeleton variant="rectangular" height={48} sx={{ mb: 1 }} />

            {Array.from({ length: rows }).map((_, i) => (
                <Skeleton key={i} variant="rectangular" height={52} sx={{ mb: 0.5 }} />
            ))}

            <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                <Skeleton variant="rounded" width={200} height={32} />
            </Box>
        </Box>
    </Card>
);

export const CardSkeleton: React.FC = () => (
    <Card>
        <CardContent>
            <Skeleton variant="text" width="50%" height={28} sx={{ mb: 2 }} />
            <Skeleton variant="rectangular" height={100} sx={{ mb: 2 }} />
            <Skeleton variant="text" width="80%" />
            <Skeleton variant="text" width="60%" />
        </CardContent>
    </Card>
);

export const PageHeaderSkeleton: React.FC = () => (
    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Box>
            <Skeleton variant="text" width={200} height={40} />
            <Skeleton variant="text" width={300} height={24} />
        </Box>
        <Skeleton variant="rounded" width={140} height={40} />
    </Box>
);
