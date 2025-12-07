import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Box,
    Typography,
    Grid,
    Alert,
} from '@mui/material';
import { useState } from 'react';

interface EditLimitsDialogProps {
    open: boolean;
    onClose: () => void;
    currentLimits: {
        employees: number;
        devices: number;
        storage: number;
    };
    onSave: (newLimits: { employees: number; devices: number; storage: number }) => void;
}

export const EditLimitsDialog: React.FC<EditLimitsDialogProps> = ({ open, onClose, currentLimits, onSave }) => {
    const [limits, setLimits] = useState(currentLimits);

    const handleSave = () => {
        onSave(limits);
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>Tùy chỉnh Giới hạn (Quota)</DialogTitle>
            <DialogContent>
                <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
                    <Alert severity="info">
                        Thay đổi giới hạn thủ công sẽ ghi đè giới hạn mặc định của gói dịch vụ hiện tại.
                    </Alert>

                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
                                Giới hạn Nhân viên
                            </Typography>
                            <TextField
                                fullWidth
                                type="number"
                                value={limits.employees}
                                onChange={(e) => setLimits({ ...limits, employees: Number(e.target.value) })}
                                helperText="Số lượng nhân viên tối đa được phép tạo"
                            />
                        </Grid>
                        <Grid item xs={6}>
                            <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
                                Giới hạn Thiết bị
                            </Typography>
                            <TextField
                                fullWidth
                                type="number"
                                value={limits.devices}
                                onChange={(e) => setLimits({ ...limits, devices: Number(e.target.value) })}
                                helperText="Số thiết bị chấm công tối đa"
                            />
                        </Grid>
                        <Grid item xs={6}>
                            <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
                                Giới hạn Lưu trữ (GB)
                            </Typography>
                            <TextField
                                fullWidth
                                type="number"
                                value={limits.storage}
                                onChange={(e) => setLimits({ ...limits, storage: Number(e.target.value) })}
                                helperText="Dung lượng lưu trữ hình ảnh"
                            />
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
