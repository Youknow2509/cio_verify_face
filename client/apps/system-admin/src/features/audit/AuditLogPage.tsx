import { useState } from 'react';
import {
    Box,
    Card,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    Avatar,
    Chip,
    TextField,
    InputAdornment,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    IconButton,
    Tooltip,
    Button,
} from '@mui/material';
import {
    Search,
    FilterList,
    Visibility,
    AdminPanelSettings,
    Edit,
    Delete,
    Lock,
    Login,
    Security,
} from '@mui/icons-material';

interface AuditLog {
    id: string;
    timestamp: string;
    actor: {
        id: string;
        name: string;
        email: string;
        avatar?: string;
        role: string;
    };
    action: string;
    target: string; // Target resource (e.g., "Company: Vinamilk")
    details: string;
    ipAddress: string;
    severity: 'info' | 'warning' | 'critical';
}

const mockLogs: AuditLog[] = [
    {
        id: '1',
        timestamp: '2024-12-07 23:45:12',
        actor: { id: 'u1', name: 'Super Admin', email: 'admin@system.com', role: 'Super Admin' },
        action: 'UPDATE_COMPANY_LIMITS',
        target: 'Company: Vinamilk',
        details: 'Increased employee limit from 2000 to 2500',
        ipAddress: '192.168.1.1',
        severity: 'warning',
    },
    {
        id: '2',
        timestamp: '2024-12-07 22:15:00',
        actor: { id: 'u2', name: 'System Manager', email: 'manager@system.com', role: 'Admin' },
        action: 'SUSPEND_COMPANY',
        target: 'Company: XYZ Ltd',
        details: 'Reason: Non-payment for 3 months',
        ipAddress: '192.168.1.25',
        severity: 'critical',
    },
    {
        id: '3',
        timestamp: '2024-12-07 20:10:33',
        actor: { id: 'u1', name: 'Super Admin', email: 'admin@system.com', role: 'Super Admin' },
        action: 'CREATE_PLAN',
        target: 'ServicePlan: Enterprise Plus',
        details: 'Created new service plan with ID pl_ent_plus',
        ipAddress: '192.168.1.1',
        severity: 'info',
    },
    {
        id: '4',
        timestamp: '2024-12-07 19:30:45',
        actor: { id: 'u1', name: 'Super Admin', email: 'admin@system.com', role: 'Super Admin' },
        action: 'LOGIN_SUCCESS',
        target: 'System',
        details: 'Successful authentication via MFA',
        ipAddress: '192.168.1.1',
        severity: 'info',
    },
    {
        id: '5',
        timestamp: '2024-12-07 18:05:22',
        actor: { id: 'u3', name: 'Support Staff', email: 'support@system.com', role: 'Support' },
        action: 'VIEW_COMPANY_DETAILS',
        target: 'Company: FPT Software',
        details: 'Accessed company configuration details',
        ipAddress: '10.0.0.52',
        severity: 'info',
    },
];

const severityColors = {
    info: 'default',
    warning: 'warning',
    critical: 'error',
} as const;

export const AuditLogPage: React.FC = () => {
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [search, setSearch] = useState('');
    const [severityFilter, setSeverityFilter] = useState('all');

    const filteredLogs = mockLogs.filter((log) => {
        const matchesSearch = log.actor.name.toLowerCase().includes(search.toLowerCase()) ||
            log.action.toLowerCase().includes(search.toLowerCase()) ||
            log.target.toLowerCase().includes(search.toLowerCase());
        const matchesSeverity = severityFilter === 'all' || log.severity === severityFilter;
        return matchesSearch && matchesSeverity;
    });

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" fontWeight="700" mb={1}>
                    System Audit Log
                </Typography>
                <Typography variant="body1" color="text.secondary">
                    Nhật ký truy cập và tác động hệ thống toàn cục (Security Trail)
                </Typography>
            </Box>

            <Card sx={{ mb: 3 }}>
                <Box sx={{ p: 2, display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'center' }}>
                    <TextField
                        placeholder="Tìm kiếm Actor, Action, Target..."
                        size="small"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        sx={{ minWidth: 350, flex: 1 }}
                        InputProps={{
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Search color="action" />
                                </InputAdornment>
                            ),
                        }}
                    />
                    <FormControl size="small" sx={{ minWidth: 150 }}>
                        <InputLabel>Mức độ</InputLabel>
                        <Select
                            value={severityFilter}
                            label="Mức độ"
                            onChange={(e) => setSeverityFilter(e.target.value)}
                        >
                            <MenuItem value="all">Tất cả</MenuItem>
                            <MenuItem value="info">Info</MenuItem>
                            <MenuItem value="warning">Warning</MenuItem>
                            <MenuItem value="critical">Critical</MenuItem>
                        </Select>
                    </FormControl>
                    <Button variant="outlined" startIcon={<FilterList />}>
                        Bộ lọc nâng cao
                    </Button>
                </Box>
            </Card>

            <Card>
                <TableContainer>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>Thời gian</TableCell>
                                <TableCell>Người thực hiện</TableCell>
                                <TableCell>Hành động</TableCell>
                                <TableCell>Đối tượng</TableCell>
                                <TableCell>Chi tiết</TableCell>
                                <TableCell>IP Address</TableCell>
                                <TableCell align="right">Chi tiết</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {filteredLogs
                                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                .map((log) => (
                                    <TableRow key={log.id} hover>
                                        <TableCell sx={{ whiteSpace: 'nowrap', color: 'text.secondary', fontSize: '0.85rem' }}>
                                            {log.timestamp}
                                        </TableCell>
                                        <TableCell>
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                                <Avatar sx={{ width: 28, height: 28, fontSize: '0.75rem', bgcolor: 'primary.main' }}>
                                                    {log.actor.name.charAt(0)}
                                                </Avatar>
                                                <Box>
                                                    <Typography variant="body2" fontWeight="600">
                                                        {log.actor.name}
                                                    </Typography>
                                                    <Typography variant="caption" color="text.secondary">
                                                        {log.actor.role}
                                                    </Typography>
                                                </Box>
                                            </Box>
                                        </TableCell>
                                        <TableCell>
                                            <Chip
                                                label={log.action}
                                                size="small"
                                                color={severityColors[log.severity]}
                                                variant={log.severity === 'critical' ? 'filled' : 'outlined'}
                                                sx={{ fontWeight: 600, fontSize: '0.7rem' }}
                                            />
                                        </TableCell>
                                        <TableCell sx={{ fontWeight: 500 }}>
                                            {log.target}
                                        </TableCell>
                                        <TableCell sx={{ maxWidth: 300 }}>
                                            <Typography variant="body2" noWrap title={log.details}>
                                                {log.details}
                                            </Typography>
                                        </TableCell>
                                        <TableCell sx={{ fontFamily: 'monospace', fontSize: '0.85rem' }}>
                                            {log.ipAddress}
                                        </TableCell>
                                        <TableCell align="right">
                                            <Tooltip title="Xem chi tiết">
                                                <IconButton size="small">
                                                    <Visibility fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                        </TableCell>
                                    </TableRow>
                                ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                <TablePagination
                    component="div"
                    count={filteredLogs.length}
                    page={page}
                    onPageChange={(_, newPage) => setPage(newPage)}
                    rowsPerPage={rowsPerPage}
                    onRowsPerPageChange={(e) => setRowsPerPage(parseInt(e.target.value, 10))}
                    labelRowsPerPage="Số dòng:"
                />
            </Card>
        </Box>
    );
};
