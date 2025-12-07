import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Box,
    Grid,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
} from '@mui/material';
import { useState, useEffect } from 'react';

// Reusing the interface from list page or defining a shared one is better, 
// but for now we define local props matching the data structure
interface CompanyData {
    id: string;
    name: string;
    email: string;
    plan: 'starter' | 'professional' | 'enterprise';
    status: 'active' | 'suspended' | 'expired';
}

interface EditCompanyDialogProps {
    open: boolean;
    onClose: () => void;
    company: CompanyData | null;
    onSave: (updatedCompany: CompanyData) => void;
}

export const EditCompanyDialog: React.FC<EditCompanyDialogProps> = ({ open, onClose, company, onSave }) => {
    const [formData, setFormData] = useState<CompanyData | null>(null);

    useEffect(() => {
        if (company) {
            setFormData({ ...company });
        }
    }, [company]);

    const handleChange = (field: keyof CompanyData, value: string) => {
        if (formData) {
            setFormData({ ...formData, [field]: value });
        }
    };

    const handleSave = () => {
        if (formData) {
            onSave(formData);
            onClose();
        }
    };

    if (!formData) return null;

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>Chỉnh sửa thông tin Công ty</DialogTitle>
            <DialogContent>
                <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                label="Tên Công ty"
                                value={formData.name}
                                onChange={(e) => handleChange('name', e.target.value)}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                label="Email liên hệ"
                                value={formData.email}
                                onChange={(e) => handleChange('email', e.target.value)}
                            />
                        </Grid>
                        <Grid item xs={6}>
                            <FormControl fullWidth>
                                <InputLabel>Gói dịch vụ</InputLabel>
                                <Select
                                    value={formData.plan}
                                    label="Gói dịch vụ"
                                    onChange={(e) => handleChange('plan', e.target.value)}
                                >
                                    <MenuItem value="starter">Starter</MenuItem>
                                    <MenuItem value="professional">Professional</MenuItem>
                                    <MenuItem value="enterprise">Enterprise</MenuItem>
                                </Select>
                            </FormControl>
                        </Grid>
                        <Grid item xs={6}>
                            <FormControl fullWidth>
                                <InputLabel>Trạng thái</InputLabel>
                                <Select
                                    value={formData.status}
                                    label="Trạng thái"
                                    onChange={(e) => handleChange('status', e.target.value)}
                                >
                                    <MenuItem value="active">Active</MenuItem>
                                    <MenuItem value="suspended">Suspended</MenuItem>
                                    <MenuItem value="expired">Expired</MenuItem>
                                </Select>
                            </FormControl>
                        </Grid>
                    </Grid>
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="inherit">Hủy</Button>
                <Button onClick={handleSave} variant="contained" color="primary">Lưu thay đổi</Button>
            </DialogActions>
        </Dialog>
    );
};
