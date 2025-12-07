import { useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Grid,
    Button,
    Chip,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Switch,
    FormControlLabel,
    InputAdornment,
} from '@mui/material';
import {
    Check,
    Add,
    Edit,
    Star,
    Business,
    Storage,
    Devices,
    AttachMoney,
} from '@mui/icons-material';

interface ServicePlan {
    id: string;
    name: string;
    price: number;
    billingCycle: 'monthly' | 'yearly';
    features: string[];
    limits: {
        employees: number;
        storage: number; // GB
        devices: number;
    };
    isPopular?: boolean;
    color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning';
}

const initialPlans: ServicePlan[] = [
    {
        id: 'starter',
        name: 'Starter',
        price: 49,
        billingCycle: 'monthly',
        features: [
            'Basic Attendance Tracking',
            'Standard Reports',
            'Email Support',
            'Data Retention: 3 months',
        ],
        limits: {
            employees: 50,
            storage: 10,
            devices: 2,
        },
        color: 'info',
    },
    {
        id: 'professional',
        name: 'Professional',
        price: 99,
        billingCycle: 'monthly',
        features: [
            'Advanced Analytics',
            'Custom Reports',
            'Priority Support',
            'Data Retention: 1 year',
            'API Access',
        ],
        limits: {
            employees: 200,
            storage: 50,
            devices: 10,
        },
        isPopular: true,
        color: 'primary',
    },
    {
        id: 'enterprise',
        name: 'Enterprise',
        price: 299,
        billingCycle: 'monthly',
        features: [
            'Unlimited History',
            'Dedicated Account Manager',
            'SLA 99.9%',
            'SSO Integration',
            'Custom Domain',
            'On-premise Deployment Option',
        ],
        limits: {
            employees: 1000,
            storage: 500,
            devices: 50,
        },
        color: 'secondary',
    },
];

export const ServicePlansPage: React.FC = () => {
    const [plans, setPlans] = useState<ServicePlan[]>(initialPlans);
    const [openDialog, setOpenDialog] = useState(false);
    const [editingPlan, setEditingPlan] = useState<ServicePlan | null>(null);

    const handleEdit = (plan: ServicePlan) => {
        setEditingPlan(plan);
        setOpenDialog(true);
    };

    const handleSave = () => {
        if (editingPlan) {
            setPlans(plans.map((p) => (p.id === editingPlan.id ? editingPlan : p)));
            setOpenDialog(false);
            setEditingPlan(null);
        }
    };

    return (
        <Box sx={{ animation: 'fadeIn 0.5s ease-out' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Box>
                    <Typography variant="h4" fontWeight="700" mb={1}>
                        Gói Dịch vụ (Service Plans)
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Quản lý các gói cước và giới hạn tài nguyên cho khách hàng
                    </Typography>
                </Box>
                <Button variant="contained" startIcon={<Add />} color="primary">
                    Thêm Gói mới
                </Button>
            </Box>

            <Grid container spacing={3}>
                {plans.map((plan) => (
                    <Grid item xs={12} md={4} key={plan.id}>
                        <Card
                            sx={{
                                height: '100%',
                                position: 'relative',
                                border: plan.isPopular ? '2px solid' : '1px solid',
                                borderColor: plan.isPopular ? 'primary.main' : 'divider',
                                transform: plan.isPopular ? 'scale(1.02)' : 'none',
                                transition: 'transform 0.2s',
                            }}
                        >
                            {plan.isPopular && (
                                <Chip
                                    label="Most Popular"
                                    color="primary"
                                    size="small"
                                    sx={{
                                        position: 'absolute',
                                        top: 12,
                                        right: 12,
                                        fontWeight: 600,
                                    }}
                                />
                            )}
                            <CardContent sx={{ p: 4 }}>
                                <Box sx={{ textAlign: 'center', mb: 3 }}>
                                    <Typography variant="h5" fontWeight="700" gutterBottom>
                                        {plan.name}
                                    </Typography>
                                    <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'baseline' }}>
                                        <Typography variant="h3" fontWeight="800" color="text.primary">
                                            ${plan.price}
                                        </Typography>
                                        <Typography variant="subtitle1" color="text.secondary" ml={1}>
                                            /{plan.billingCycle === 'monthly' ? 'tháng' : 'năm'}
                                        </Typography>
                                    </Box>
                                </Box>

                                <Typography variant="subtitle2" fontWeight="700" sx={{ mb: 2, mt: 3, textTransform: 'uppercase', letterSpacing: 1, fontSize: '0.75rem', color: 'text.secondary' }}>
                                    Giới hạn (Limits)
                                </Typography>

                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5, mb: 3 }}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                        <Storage fontSize="small" color="action" />
                                        <Typography variant="body2">
                                            <b>{plan.limits.storage} GB</b> Storage
                                        </Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                        <Business fontSize="small" color="action" />
                                        <Typography variant="body2">
                                            <b>{plan.limits.employees}</b> Nhân viên
                                        </Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                        <Devices fontSize="small" color="action" />
                                        <Typography variant="body2">
                                            <b>{plan.limits.devices}</b> Thiết bị
                                        </Typography>
                                    </Box>
                                </Box>

                                <Typography variant="subtitle2" fontWeight="700" sx={{ mb: 2, textTransform: 'uppercase', letterSpacing: 1, fontSize: '0.75rem', color: 'text.secondary' }}>
                                    Tính năng (Features)
                                </Typography>

                                <List dense disablePadding sx={{ mb: 3 }}>
                                    {plan.features.map((feature, index) => (
                                        <ListItem key={index} disablePadding sx={{ mb: 1 }}>
                                            <ListItemIcon sx={{ minWidth: 32 }}>
                                                <Check color="success" fontSize="small" />
                                            </ListItemIcon>
                                            <ListItemText primary={feature} />
                                        </ListItem>
                                    ))}
                                </List>

                                <Button
                                    variant="outlined"
                                    fullWidth
                                    startIcon={<Edit />}
                                    onClick={() => handleEdit(plan)}
                                >
                                    Chỉnh sửa
                                </Button>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>

            {/* Edit Dialog */}
            <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
                <DialogTitle>Chỉnh sửa Gói {editingPlan?.name}</DialogTitle>
                <DialogContent>
                    {editingPlan && (
                        <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
                            <Box sx={{ display: 'flex', gap: 2 }}>
                                <TextField
                                    fullWidth
                                    label="Tên gói"
                                    value={editingPlan.name}
                                    onChange={(e) => setEditingPlan({ ...editingPlan, name: e.target.value })}
                                />
                                <TextField
                                    fullWidth
                                    label="Giá ($)"
                                    type="number"
                                    value={editingPlan.price}
                                    onChange={(e) => setEditingPlan({ ...editingPlan, price: Number(e.target.value) })}
                                    InputProps={{
                                        startAdornment: <InputAdornment position="start">$</InputAdornment>,
                                    }}
                                />
                            </Box>

                            <Typography variant="subtitle2" color="primary">Giới hạn Tài nguyên</Typography>

                            <Grid container spacing={2}>
                                <Grid item xs={4}>
                                    <TextField
                                        fullWidth
                                        label="Nhân viên"
                                        type="number"
                                        value={editingPlan.limits.employees}
                                        onChange={(e) => setEditingPlan({
                                            ...editingPlan,
                                            limits: { ...editingPlan.limits, employees: Number(e.target.value) }
                                        })}
                                    />
                                </Grid>
                                <Grid item xs={4}>
                                    <TextField
                                        fullWidth
                                        label="Storage (GB)"
                                        type="number"
                                        value={editingPlan.limits.storage}
                                        onChange={(e) => setEditingPlan({
                                            ...editingPlan,
                                            limits: { ...editingPlan.limits, storage: Number(e.target.value) }
                                        })}
                                    />
                                </Grid>
                                <Grid item xs={4}>
                                    <TextField
                                        fullWidth
                                        label="Thiết bị"
                                        type="number"
                                        value={editingPlan.limits.devices}
                                        onChange={(e) => setEditingPlan({
                                            ...editingPlan,
                                            limits: { ...editingPlan.limits, devices: Number(e.target.value) }
                                        })}
                                    />
                                </Grid>
                            </Grid>

                            <FormControlLabel
                                control={
                                    <Switch
                                        checked={editingPlan.isPopular}
                                        onChange={(e) => setEditingPlan({ ...editingPlan, isPopular: e.target.checked })}
                                    />
                                }
                                label="Đánh dấu là Phổ biến (Recommended)"
                            />
                        </Box>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenDialog(false)} color="inherit">Hủy</Button>
                    <Button onClick={handleSave} variant="contained" color="primary">Lưu thay đổi</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
