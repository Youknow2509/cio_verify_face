import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    LinearProgress,
    Divider,
    Button,
    Alert,
    Stack,
} from '@mui/material';
import {
    CheckCircle,
    Error as ErrorIcon,
    Warning,
    Storage,
    Refresh,
    Dns,
    Cloud,
    Speed,
    Memory,
} from '@mui/icons-material';
import {
    AreaChart,
    Area,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';
import { useState, useEffect } from 'react';

// Standard Microservices from System Architecture (Section 3.1.4)
const MICROSERVICES = [
    { id: 'auth', name: 'Auth Service', port: 4000 },
    { id: 'identity', name: 'Identity & Org Ops', port: 4001 },
    { id: 'device', name: 'Device Management', port: 4002 },
    { id: 'workforce', name: 'Workforce Service', port: 4005 },
    { id: 'attendance', name: 'Attendance Core', port: 4003 },
    { id: 'analytic', name: 'Analytics & Reporting', port: 4006 },
    { id: 'signature', name: 'Signature Service', port: 4008 },
    { id: 'ws-delivery', name: 'WS Delivery', port: 8080 },
];

interface HealthStatus {
    serviceId: string;
    status: 'running' | 'down' | 'degraded';
    latency: number;
    uptime: number; // seconds
    lastCheck: string;
}

// Mock API Call simulating real backend health check
const checkServiceHealth = async (id: string): Promise<HealthStatus> => {
    // Simulate network delay
    await new Promise(r => setTimeout(r, 200 + Math.random() * 500));

    // Simulate simulation logic (95% running, 5% degraded, rare down)
    const rand = Math.random();
    let status: 'running' | 'down' | 'degraded' = 'running';
    if (id === 'attendance' && rand > 0.8) status = 'degraded'; // Simulate heavy load on attendance

    return {
        serviceId: id,
        status,
        latency: Math.floor(Math.random() * 50) + 5, // Low latency internal network
        uptime: 3600 * 48, // 2 days uptime
        lastCheck: new Date().toISOString(),
    };
};

const statusColors = {
    running: 'success',
    down: 'error',
    degraded: 'warning',
} as const;

const statusIcons = {
    running: <CheckCircle color="success" fontSize="small" />,
    down: <ErrorIcon color="error" fontSize="small" />,
    degraded: <Warning color="warning" fontSize="small" />,
};

export const MonitoringPage: React.FC = () => {
    const [statuses, setStatuses] = useState<Record<string, HealthStatus>>({});
    const [loading, setLoading] = useState(false);
    const [lastUpdated, setLastUpdated] = useState<Date>(new Date());

    const fetchHealth = async () => {
        setLoading(true);
        const newStatuses: Record<string, HealthStatus> = {};

        // Parallel checks simulating real dashboard behavior
        await Promise.all(MICROSERVICES.map(async (service) => {
            const result = await checkServiceHealth(service.id);
            newStatuses[service.id] = result;
        }));

        setStatuses(newStatuses);
        setLastUpdated(new Date());
        setLoading(false);
    };

    useEffect(() => {
        fetchHealth();
        const interval = setInterval(fetchHealth, 30000); // Auto refresh every 30s
        return () => clearInterval(interval);
    }, []);

    // Derived Simulation Data for infrastructure
    const totalServices = MICROSERVICES.length;
    const onlineServices = Object.values(statuses).filter(s => s.status === 'running').length;
    const systemHealth = Math.round((onlineServices / totalServices) * 100);

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Box>
                    <Typography variant="h4" fontWeight="700" mb={1}>
                        System Health Monitor
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Infrastructure & Services Overview (Architecture 3.1)
                    </Typography>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="caption" color="text.secondary">
                        Last updated: {lastUpdated.toLocaleTimeString()}
                    </Typography>
                    <Button
                        variant="outlined"
                        size="small"
                        startIcon={<Refresh />}
                        onClick={fetchHealth}
                        disabled={loading}
                    >
                        Refresh
                    </Button>
                </Box>
            </Box>

            {/* Critical Infrastructure Status (Data Storage Layer - 3.1.7) */}
            <Typography variant="h6" fontWeight="700" mb={2} sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Storage color="secondary" /> Data & Messaging Layer
            </Typography>
            <Grid container spacing={2} sx={{ mb: 4 }}>
                <Grid item xs={12} sm={6} md={3}>
                    <Alert icon={<Storage fontSize="inherit" />} severity="success" variant="outlined" sx={{ height: '100%' }}>
                        <Typography fontWeight="700">PostgreSQL Primary</Typography>
                        <Stack direction="row" spacing={1} alignItems="center" mt={1}>
                            <Speed fontSize="small" color="disabled" />
                            <Typography variant="caption">Latency: 2ms</Typography>
                        </Stack>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <Cloud fontSize="small" color="disabled" />
                            <Typography variant="caption">Conn: 85/100</Typography>
                        </Stack>
                    </Alert>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <Alert icon={<Storage fontSize="inherit" />} severity="success" variant="outlined" sx={{ height: '100%' }}>
                        <Typography fontWeight="700">ScyllaDB Cluster</Typography>
                        <Typography variant="caption" component="div" sx={{ mb: 0.5 }}>Time-Series Storage</Typography>
                        <Stack direction="row" spacing={1} alignItems="center" mt={1}>
                            <Speed fontSize="small" color="disabled" />
                            <Typography variant="caption">Write: 0.5ms</Typography>
                        </Stack>
                        <Typography variant="caption">Nodes: 3 (Healthy)</Typography>
                    </Alert>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <Alert icon={<Memory fontSize="inherit" />} severity="info" variant="outlined" sx={{ height: '100%', borderColor: 'info.main' }}>
                        <Typography fontWeight="700">Apache Kafka</Typography>
                        <Typography variant="caption" component="div" sx={{ mb: 0.5 }}>Event Bus</Typography>
                        <Stack direction="row" spacing={1} alignItems="center" mt={1}>
                            <Speed fontSize="small" color="disabled" />
                            <Typography variant="caption">Lag: 12ms</Typography>
                        </Stack>
                        <Typography variant="caption">Topics: 15 Active</Typography>
                    </Alert>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <Alert icon={<Storage fontSize="inherit" />} severity="success" variant="outlined" sx={{ height: '100%' }}>
                        <Typography fontWeight="700">Redis Cache</Typography>
                        <Stack direction="row" spacing={1} alignItems="center" mt={1}>
                            <Speed fontSize="small" color="disabled" />
                            <Typography variant="caption">Hit Rate: 99.2%</Typography>
                        </Stack>
                        <Typography variant="caption">Mem: 1.2GB/4GB</Typography>
                    </Alert>
                </Grid>
            </Grid>

            {/* Microservices Status Grid (Microservices Layer - 3.1.4) */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" fontWeight="700" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Dns color="primary" /> Microservices Cluster
                </Typography>
                <Typography variant="caption" sx={{ bgcolor: 'success.light', color: 'black', px: 1, py: 0.5, borderRadius: 1, fontWeight: 'bold' }}>
                    System Status: {systemHealth}% Operational
                </Typography>
            </Box>

            <Grid container spacing={2} sx={{ mb: 4 }}>
                {MICROSERVICES.map((service) => {
                    const status = statuses[service.id];
                    return (
                        <Grid item xs={12} sm={6} md={4} lg={3} key={service.id}>
                            <Card variant="outlined" sx={{ '&:hover': { borderColor: 'primary.main' } }}>
                                <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                                        <Typography variant="subtitle2" fontWeight="700">{service.name}</Typography>
                                        {status ? statusIcons[status.status] : <Refresh sx={{ animation: 'spin 1s infinite' }} />}
                                    </Box>

                                    <Divider sx={{ my: 1, opacity: 0.5 }} />

                                    {status ? (
                                        <Box>
                                            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                                                <Typography variant="caption" color="text.secondary">Port</Typography>
                                                <Typography variant="caption">{service.port}</Typography>
                                            </Box>
                                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                <Typography variant="caption" color="text.secondary">Response</Typography>
                                                <Typography
                                                    variant="caption"
                                                    fontWeight="600"
                                                    color={status.latency > 100 ? 'warning.main' : 'success.main'}
                                                >
                                                    {status.latency}ms
                                                </Typography>
                                            </Box>
                                        </Box>
                                    ) : (
                                        <Typography variant="caption">Checking...</Typography>
                                    )}
                                </CardContent>
                            </Card>
                        </Grid>
                    );
                })}
            </Grid>

            {/* Performance Metrics Chart */}
            <Card variant="outlined">
                <CardContent>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                        <Box>
                            <Typography variant="h6" fontWeight="700">
                                Gateway Traffic
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                Requests/sec via API Gateway
                            </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', gap: 2 }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: '#6366f1' }} />
                                <Typography variant="caption">Inbound</Typography>
                            </Box>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: '#ef4444' }} />
                                <Typography variant="caption">Errors</Typography>
                            </Box>
                        </Box>
                    </Box>
                    <Box sx={{ height: 250 }}>
                        <ResponsiveContainer width="100%" height="100%">
                            <AreaChart data={[
                                { time: '10:00', requests: 1200, errors: 5 },
                                { time: '10:05', requests: 1500, errors: 8 },
                                { time: '10:10', requests: 1100, errors: 4 },
                                { time: '10:15', requests: 1800, errors: 12 },
                                { time: '10:20', requests: 2200, errors: 15 },
                                { time: '10:25', requests: 1900, errors: 10 },
                                { time: '10:30', requests: 2500, errors: 8 },
                            ]}>
                                <CartesianGrid strokeDasharray="3 3" opacity={0.3} vertical={false} />
                                <XAxis dataKey="time" tick={{ fontSize: 12 }} axisLine={false} tickLine={false} />
                                <YAxis tick={{ fontSize: 12 }} axisLine={false} tickLine={false} />
                                <Tooltip />
                                <Area type="monotone" dataKey="requests" stroke="#6366f1" strokeWidth={2} fill="#6366f1" fillOpacity={0.1} />
                                <Area type="monotone" dataKey="errors" stroke="#ef4444" strokeWidth={2} fill="#ef4444" fillOpacity={0.1} />
                            </AreaChart>
                        </ResponsiveContainer>
                    </Box>
                </CardContent>
            </Card>
        </Box>
    );
};
