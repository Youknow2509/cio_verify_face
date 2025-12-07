import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    Typography,
    Button,
    Chip,
    Avatar,
    IconButton,
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
    Tooltip,
    LinearProgress,
} from '@mui/material';
import {
    DataGrid,
    GridColDef,
    GridRenderCellParams,
    GridValueGetterParams,
} from '@mui/x-data-grid';
import {
    Add,
    MoreVert,
    Visibility,
    Edit,
    Block,
    CheckCircle,
    Delete,
    Business,
} from '@mui/icons-material';
import { EditCompanyDialog } from './EditCompanyDialog';

interface Company {
    id: string;
    name: string;
    logo?: string;
    email: string;
    plan: 'starter' | 'professional' | 'enterprise';
    status: 'active' | 'suspended' | 'expired';
    employeeCount: number;
    employeeLimit: number;
    deviceCount: number;
    storageUsed: number; // GB
    storageLimit: number; // GB
    registeredAt: string;
    expiresAt: string;
}

const initialCompanies: Company[] = [
    { id: '1', name: 'Vinamilk', email: 'admin@vinamilk.com', plan: 'enterprise', status: 'active', employeeCount: 1250, employeeLimit: 2000, deviceCount: 12, storageUsed: 42.5, storageLimit: 100, registeredAt: '2023-01-15', expiresAt: '2024-12-15' },
    { id: '2', name: 'FPT Software', email: 'hr@fpt.com', plan: 'enterprise', status: 'active', employeeCount: 980, employeeLimit: 1500, deviceCount: 8, storageUsed: 35.2, storageLimit: 100, registeredAt: '2023-03-20', expiresAt: '2024-11-20' },
    { id: '3', name: 'Tech Corp', email: 'admin@techcorp.vn', plan: 'professional', status: 'active', employeeCount: 756, employeeLimit: 1000, deviceCount: 6, storageUsed: 15.8, storageLimit: 50, registeredAt: '2023-06-10', expiresAt: '2024-12-10' },
    { id: '4', name: 'Viettel', email: 'hr@viettel.com', plan: 'enterprise', status: 'active', employeeCount: 650, employeeLimit: 1000, deviceCount: 5, storageUsed: 22.1, storageLimit: 50, registeredAt: '2023-02-28', expiresAt: '2025-02-28' },
    { id: '5', name: 'XYZ Ltd', email: 'admin@xyz.com', plan: 'starter', status: 'suspended', employeeCount: 45, employeeLimit: 50, deviceCount: 1, storageUsed: 2.5, storageLimit: 10, registeredAt: '2023-09-01', expiresAt: '2024-09-01' },
    { id: '6', name: 'ABC Company', email: 'contact@abc.vn', plan: 'professional', status: 'expired', employeeCount: 320, employeeLimit: 500, deviceCount: 3, storageUsed: 12.0, storageLimit: 50, registeredAt: '2023-04-15', expiresAt: '2024-10-15' },
];

const planColors = {
    starter: 'default',
    professional: 'primary',
    enterprise: 'secondary',
} as const;

const statusColors = {
    active: 'success',
    suspended: 'warning',
    expired: 'error',
} as const;

export const CompanyListPage: React.FC = () => {
    const navigate = useNavigate();
    const [companies, setCompanies] = useState<Company[]>(initialCompanies);
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [selectedCompanyId, setSelectedCompanyId] = useState<string | null>(null);
    const [editDialogOpen, setEditDialogOpen] = useState(false);

    const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, id: string) => {
        setAnchorEl(event.currentTarget);
        setSelectedCompanyId(id);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
        setSelectedCompanyId(null);
    };

    const handleEditSave = (updatedCompany: any) => {
        setCompanies(companies.map(c =>
            c.id === updatedCompany.id ? { ...c, ...updatedCompany } : c
        ));
    };

    const handleStatusChange = (status: 'active' | 'suspended') => {
        if (selectedCompanyId) {
            setCompanies(companies.map(c =>
                c.id === selectedCompanyId ? { ...c, status } : c
            ));
            handleMenuClose();
        }
    };

    const columns: GridColDef[] = [
        {
            field: 'name',
            headerName: 'Công ty',
            width: 250,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32, fontSize: '0.875rem' }}>
                        {params.value.charAt(0)}
                    </Avatar>
                    <Box>
                        <Typography variant="body2" fontWeight="600" lineHeight={1.2}>
                            {params.value}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                            {params.row.email}
                        </Typography>
                    </Box>
                </Box>
            ),
        },
        {
            field: 'plan',
            headerName: 'Gói dịch vụ',
            width: 140,
            renderCell: (params: GridRenderCellParams) => (
                <Chip
                    label={params.value.toUpperCase()}
                    size="small"
                    color={planColors[params.value as keyof typeof planColors]}
                    variant="outlined"
                    sx={{ fontWeight: 600, fontSize: '0.7rem' }}
                />
            ),
        },
        {
            field: 'usage',
            headerName: 'Sử dụng (NV/TB)',
            width: 160,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                    <Typography variant="caption">
                        NV: <b>{params.row.employeeCount}</b> / {params.row.employeeLimit}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                        TB: <b>{params.row.deviceCount}</b>
                    </Typography>
                </Box>
            ),
        },
        {
            field: 'storage', // NEW COLUMN
            headerName: 'Dung lượng',
            width: 140,
            renderCell: (params: GridRenderCellParams) => (
                <Box sx={{ width: '100%' }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                        <Typography variant="caption" fontWeight="600">
                            {params.row.storageUsed} GB
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                            / {params.row.storageLimit}
                        </Typography>
                    </Box>
                    <LinearProgress
                        variant="determinate"
                        value={(params.row.storageUsed / params.row.storageLimit) * 100}
                        color={(params.row.storageUsed / params.row.storageLimit) > 0.8 ? 'warning' : 'primary'}
                        sx={{ height: 4, borderRadius: 2 }}
                    />
                </Box>
            ),
        },
        {
            field: 'status',
            headerName: 'Trạng thái',
            width: 120,
            renderCell: (params: GridRenderCellParams) => (
                <Chip
                    label={params.value}
                    size="small"
                    color={statusColors[params.value as keyof typeof statusColors]}
                    sx={{ textTransform: 'capitalize' }}
                />
            ),
        },
        {
            field: 'expiresAt',
            headerName: 'Hết hạn',
            width: 150,
            valueFormatter: (params) => new Date(params.value as string).toLocaleDateString('vi-VN'),
        },
        {
            field: 'actions',
            headerName: '',
            width: 80,
            sortable: false,
            renderCell: (params: GridRenderCellParams) => (
                <IconButton
                    size="small"
                    onClick={(e) => handleMenuOpen(e, params.row.id)}
                >
                    <MoreVert fontSize="small" />
                </IconButton>
            ),
        },
    ];

    const selectedCompany = companies.find(c => c.id === selectedCompanyId);

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out', height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Box>
                    <Typography variant="h4" fontWeight="700" mb={1}>
                        Quản lý Công ty
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Quản lý khách hàng (Tenants)
                    </Typography>
                </Box>
                <Button
                    variant="contained"
                    startIcon={<Add />}
                    sx={{
                        background: 'linear-gradient(135deg, #6366f1, #ec4899)',
                        boxShadow: '0 4px 12px rgba(99, 102, 241, 0.3)',
                    }}
                >
                    Thêm Công ty
                </Button>
            </Box>

            <Card sx={{ flex: 1 }}>
                <DataGrid
                    rows={companies}
                    columns={columns}
                    initialState={{
                        pagination: {
                            paginationModel: { page: 0, pageSize: 10 },
                        },
                    }}
                    pageSizeOptions={[10, 20]}
                    checkboxSelection
                    disableRowSelectionOnClick
                    sx={{
                        border: 'none',
                        '& .MuiDataGrid-cell': {
                            borderBottom: '1px solid #f0f0f0',
                        },
                        '& .MuiDataGrid-columnHeaders': {
                            bgcolor: 'background.default',
                            fontWeight: 700,
                        },
                    }}
                />
            </Card>


            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
            >
                <MenuItem onClick={() => { navigate(`/companies/${selectedCompanyId}`); handleMenuClose(); }}>
                    <ListItemIcon><Visibility fontSize="small" /></ListItemIcon>
                    <ListItemText>Xem chi tiết</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => { setEditDialogOpen(true); handleMenuClose(); }}>
                    <ListItemIcon><Edit fontSize="small" /></ListItemIcon>
                    <ListItemText>Chỉnh sửa</ListItemText>
                </MenuItem>
                {selectedCompany?.status === 'active' ? (
                    <MenuItem onClick={() => handleStatusChange('suspended')}>
                        <ListItemIcon><Block fontSize="small" color="warning" /></ListItemIcon>
                        <ListItemText>Tạm khóa (Suspend)</ListItemText>
                    </MenuItem>
                ) : (
                    <MenuItem onClick={() => handleStatusChange('active')}>
                        <ListItemIcon><CheckCircle fontSize="small" color="success" /></ListItemIcon>
                        <ListItemText>Kích hoạt (Activate)</ListItemText>
                    </MenuItem>
                )}
                <MenuItem onClick={handleMenuClose} sx={{ color: 'error.main' }}>
                    <ListItemIcon><Delete fontSize="small" color="error" /></ListItemIcon>
                    <ListItemText>Xóa</ListItemText>
                </MenuItem>
            </Menu>

            <EditCompanyDialog
                open={editDialogOpen}
                onClose={() => setEditDialogOpen(false)}
                company={selectedCompany || null}
                onSave={handleEditSave}
            />
        </Box>
    );
};

