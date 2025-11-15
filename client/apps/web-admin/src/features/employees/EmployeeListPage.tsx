import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Button,
    Card,
    TextField,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
    Avatar,
    Chip,
    Typography,
    CircularProgress,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Paper,
    Snackbar,
    Alert,
} from '@mui/material';
import { Add, Edit, Face, Delete, UploadFile } from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';
import type { Employee } from '@face-attendance/types';

export const EmployeeListPage: React.FC = () => {
    const navigate = useNavigate();
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');

    useEffect(() => {
        const accessToken = localStorage.getItem('access_token');
        if (accessToken) {
            const companyId = getCompanyIdFromToken(accessToken);
            if (companyId) {
                fetchEmployees(companyId);
            }
        } else {
            setLoading(false);
        }
    }, []);

    // Get company ID from JWT token
    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
        }
    };

    const fetchEmployees = async (companyId: string) => {
        try {
            const response = await apiClient.get(
                `/api/v1/users?company_id=${companyId}`
            );
            if (response.status === 200) {
                console.log('Fetched employees:', response.data);
                const mapData = response.data?.data.map((item: any) => ({
                    id: item.user_id,
                    email: item.email,
                    full_name: item.full_name,
                    employee_code: item.employee_code,
                    department: item.department,
                    phone: item.phone,
                    avatar_url: item.avatar_url,
                    role: item.role,
                    status: item.status === 0 ? 'active' : 'inactive',
                    position: item.position,
                    face_data_count: item.face_data_count || 0,
                    created_at: item.created_at,
                    updated_at: item.updated_at,
                }));
                setEmployees(mapData);
            }
        } catch (error) {
            console.error('Failed to fetch employees:', error);
        } finally {
            setLoading(false);
        }
    };

    // CSV import UI state and helpers
    const [importDialogOpen, setImportDialogOpen] = useState(false);
    const [csvPreview, setCsvPreview] = useState<any[]>([]);
    const [csvHeaders, setCsvHeaders] = useState<string[]>([]);
    const [csvFileName, setCsvFileName] = useState<string>('');
    const [importProcessing, setImportProcessing] = useState(false);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({ open: false, message: '', severity: 'success' });

    const parseCSV = (text: string) => {
        const lines: string[] = [];
        let cur = '';
        let inQuotes = false;
        for (let i = 0; i < text.length; i++) {
            const ch = text[i];
            if (ch === '"') {
                inQuotes = !inQuotes;
                cur += ch;
            } else if (ch === '\n' || ch === '\r') {
                if (!inQuotes) {
                    if (cur.trim() !== '') lines.push(cur);
                    cur = '';
                    if (ch === '\r' && text[i + 1] === '\n') i++;
                } else {
                    cur += ch;
                }
            } else {
                cur += ch;
            }
        }
        if (cur.trim() !== '') lines.push(cur);

        const rows = lines.map((line) => {
            const cells: string[] = [];
            let cell = '';
            let quoted = false;
            for (let i = 0; i < line.length; i++) {
                const ch = line[i];
                if (ch === '"') {
                    if (quoted && line[i + 1] === '"') {
                        cell += '"';
                        i++;
                    } else quoted = !quoted;
                } else if (ch === ',' && !quoted) {
                    cells.push(cell.trim());
                    cell = '';
                } else {
                    cell += ch;
                }
            }
            cells.push(cell.trim());
            return cells.map((c) => c.replace(/^"|"$/g, '').trim());
        });
        const headers = rows.length > 0 ? rows[0] : [];
        const data = rows.slice(1);
        return { headers, data };
    };

    const openImportDialog = () => {
        setCsvPreview([]);
        setCsvHeaders([]);
        setCsvFileName('');
        setImportDialogOpen(true);
    };

    const downloadSampleCsv = () => {
        const sample =
            `email,phone,password,full_name,avatar_url,company_id,employee_code,department,position,hire_date,salary,role\n` +
            `nva@gmail.com,0123124325,Passw0rd!,Nguyen Van A,https://example.com/avatars/alice.jpg,1267b1ef-52f2-425d-81db-8004c8a06316,EMP-001,HR,HR Manager,2020-06-15,1500.50,2\n` +
            `nvb@gmail.com,0937125623,Secret123,Nguyen Van Tran,,1267b1ef-52f2-425d-81db-8004c8a06316,EMP-002,Engineering,Software Engineer,2021-03-01,2000,2`;
        const blob = new Blob([sample], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'sample_employees.csv';
        a.click();
        URL.revokeObjectURL(url);
    };

    const handleFileChange = (file: File | null) => {
        if (!file) return;
        setCsvFileName(file.name);
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = String(e.target?.result || '');
            const parsed = parseCSV(text);
            setCsvHeaders(parsed.headers);
            const rows = parsed.data.map((r: string[]) => {
                const obj: any = {};
                parsed.headers.forEach((h: string, i: number) => {
                    obj[h] = r[i] ?? '';
                });
                return obj;
            });
            setCsvPreview(rows);
        };
        reader.readAsText(file, 'utf-8');
    };

    const handleConfirmImport = async () => {
        if (csvPreview.length === 0) {
            setSnackbar({
                open: true,
                message: 'Không có dữ liệu để nhập',
                severity: 'error',
            });
            return;
        }
        try {
            setImportProcessing(true);
            const promises = csvPreview.map((row) => {
                const payload: any = {
                    email: row['email'] || row['Email'] || '',
                    password: row['password'] || row['Password'] || '',
                    full_name:
                        row['full_name'] ||
                        row['fullName'] ||
                        row['full name'] ||
                        '',
                    phone: row['phone'] || row['Phone'] || '',
                    avatar_url: row['avatar_url'] || row['avatarUrl'] || '',
                    company_id: row['company_id'] || row['companyId'] || '',
                    employee_code:
                        row['employee_code'] ||
                        row['employee'] ||
                        row['employee_code'] ||
                        '',
                    department: row['department'] || '',
                    position: row['position'] || '',
                    hire_date: row['hire_date'] || row['hireDate'] || '',
                    salary: row['salary'] || '',
                    role: row['role'] || '',
                };
                return apiClient.post('/api/v1/users', payload);
            });
            const results = await Promise.allSettled(promises);
            const successCount = results.filter(
                (r) => r.status === 'fulfilled'
            ).length;
            const failCount = results.length - successCount;
            setSnackbar({
                open: true,
                message: `Hoàn tất: ${successCount} thành công, ${failCount} thất bại`,
                severity: failCount === 0 ? 'success' : 'error',
            });
            setImportDialogOpen(false);
            const accessToken = localStorage.getItem('access_token');
            if (accessToken) {
                const companyId = getCompanyIdFromToken(accessToken);
                if (companyId) await fetchEmployees(companyId);
            }
        } catch (err) {
            console.error('Import error', err);
            setSnackbar({
                open: true,
                message: 'Lỗi khi nhập từ file',
                severity: 'error',
            });
        } finally {
            setImportProcessing(false);
        }
    };

    const handleCloseSnackbar = () => setSnackbar({ ...snackbar, open: false });

    const filteredEmployees = employees.filter(
        (emp) =>
            emp.full_name.toLowerCase().includes(search.toLowerCase()) ||
            emp.employee_code?.toLowerCase().includes(search.toLowerCase()) ||
            emp.email.toLowerCase().includes(search.toLowerCase()) ||
            emp.phone?.toLowerCase().includes(search.toLowerCase())
    );

    // role to str (not used)

    // handle delete employee
    const handleDelete = async (employeeId: string) => {
        if (!window.confirm('Bạn có chắc chắn muốn xóa nhân viên này không?')) {
            return;
        }
        try {
            const response = await apiClient.delete(
                `/api/v1/users/${employeeId}`
            );
            if (response.status === 200) {
                setEmployees(employees.filter((emp) => emp.id !== employeeId));
            } else {
                console.error('Failed to delete employee:', response);
            }
        } catch (error) {
            console.error('Error deleting employee:', error);
        }
    };

    if (loading) {
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="400px"
            >
                <CircularProgress />
            </Box>
        );
    }

    return (
        <>
            <Box>
                <Box
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                    mb={3}
                >
                    <Typography variant="h4" fontWeight="bold">
                        Quản lý Nhân viên
                    </Typography>
                    <Box display="flex" gap={1}>
                        <Button
                            variant="contained"
                            startIcon={<Add />}
                            onClick={() => navigate('/employees/add')}
                        >
                            Thêm nhân viên
                        </Button>
                        <Button
                            variant="outlined"
                            startIcon={<UploadFile />}
                            onClick={openImportDialog}
                        >
                            Thêm từ file
                        </Button>
                    </Box>
                </Box>
                <Card>
                    <Box p={2}>
                        <TextField
                            fullWidth
                            placeholder="Tìm kiếm theo tên hoặc mã nhân viên..."
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                        />
                    </Box>
                    <TableContainer>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Nhân viên</TableCell>
                                    <TableCell>Mã NV</TableCell>
                                    <TableCell>Email</TableCell>
                                    <TableCell>Số điện thoại</TableCell>
                                    <TableCell>Phòng ban</TableCell>
                                    <TableCell>Chức vụ</TableCell>
                                    <TableCell>Trạng thái</TableCell>
                                    <TableCell>Khuôn mặt</TableCell>
                                    <TableCell align="right">
                                        Thao tác
                                    </TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {filteredEmployees.map((employee) => (
                                    <TableRow key={employee.id}>
                                        <TableCell>
                                            <Box
                                                display="flex"
                                                alignItems="center"
                                                gap={2}
                                            >
                                                <Avatar
                                                    src={employee.avatar_url}
                                                >
                                                    {employee.full_name[0]}
                                                </Avatar>
                                                <Typography>
                                                    {employee.full_name}
                                                </Typography>
                                            </Box>
                                        </TableCell>
                                        <TableCell>
                                            {employee.employee_code}
                                        </TableCell>
                                        <TableCell>{employee.email}</TableCell>
                                        <TableCell>{employee.phone}</TableCell>
                                        <TableCell>
                                            {employee.department}
                                        </TableCell>
                                        <TableCell>
                                            {employee.position}
                                        </TableCell>
                                        <TableCell>
                                            <Chip
                                                label={
                                                    employee.status === 'active'
                                                        ? 'Hoạt động'
                                                        : 'Không hoạt động'
                                                }
                                                color={
                                                    employee.status === 'active'
                                                        ? 'success'
                                                        : 'default'
                                                }
                                                size="small"
                                            />
                                        </TableCell>
                                        <TableCell>
                                            <Chip
                                                label={employee.face_data_count}
                                                size="small"
                                            />
                                        </TableCell>
                                        <TableCell align="right">
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    navigate(
                                                        `/employees/${employee.id}/face-data`
                                                    )
                                                }
                                            >
                                                <Face />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={() =>
                                                    navigate(
                                                        `/employees/${employee.id}/edit`
                                                    )
                                                }
                                            >
                                                <Edit />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                color="error"
                                            >
                                                <Delete
                                                    onClick={() =>
                                                        handleDelete(
                                                            employee.id
                                                        )
                                                    }
                                                />
                                            </IconButton>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Card>
            </Box>
            {/* Import Dialog */}
            <Dialog
                open={importDialogOpen}
                onClose={() => setImportDialogOpen(false)}
                maxWidth="lg"
                fullWidth
            >
                <DialogTitle>Nhập nhân viên từ file CSV</DialogTitle>
                <DialogContent>
                    <Box mb={2}>
                        <Typography variant="body2" color="text.secondary">
                            File mẫu: cột (không bắt buộc thứ tự):
                            <Box component="span" ml={1} fontFamily="monospace">
                                email,phone,password,full_name,avatar_url,company_id,employee_code,department,position,hire_date,salary,role
                            </Box>
                        </Typography>
                        <Box mt={1} display="flex" gap={1}>
                            <Button size="small" onClick={downloadSampleCsv}>
                                Tải file mẫu
                            </Button>
                            <Button size="small" component="label">
                                Chọn file
                                <input
                                    type="file"
                                    accept=".csv,text/csv"
                                    hidden
                                    onChange={(e) =>
                                        handleFileChange(
                                            e.target.files
                                                ? e.target.files[0]
                                                : null
                                        )
                                    }
                                />
                            </Button>
                            <Typography
                                variant="caption"
                                color="text.secondary"
                                sx={{ alignSelf: 'center' }}
                            >
                                {csvFileName || 'Chưa chọn file'}
                            </Typography>
                        </Box>
                    </Box>

                    {csvPreview.length === 0 ? (
                        <Alert severity="info">
                            Chưa có dữ liệu để xem trước
                        </Alert>
                    ) : (
                        <TableContainer
                            component={Paper}
                            sx={{ maxHeight: 360 }}
                        >
                            <Table stickyHeader size="small">
                                <TableHead>
                                    <TableRow>
                                        {csvHeaders.map((h) => (
                                            <TableCell key={h}>{h}</TableCell>
                                        ))}
                                        <TableCell>Trạng thái</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {csvPreview.map((row, idx) => (
                                        <TableRow key={idx} hover>
                                            {csvHeaders.map((h) => (
                                                <TableCell key={h}>
                                                    {row[h]}
                                                </TableCell>
                                            ))}
                                            <TableCell>
                                                {/* If backend needs matching to existing employees, implement here. For create, just show Ready */}
                                                <Chip
                                                    label="Sẵn sàng"
                                                    size="small"
                                                    color="default"
                                                />
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={() => setImportDialogOpen(false)}
                        disabled={importProcessing}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleConfirmImport}
                        variant="contained"
                        disabled={importProcessing || csvPreview.length === 0}
                    >
                        {importProcessing ? 'Đang xử lý...' : 'Nhập từ file'}
                    </Button>
                </DialogActions>
            </Dialog>

            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={handleCloseSnackbar}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={handleCloseSnackbar}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </>
    );
};
