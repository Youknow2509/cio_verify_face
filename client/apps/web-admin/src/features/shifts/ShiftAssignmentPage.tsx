import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Typography,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
    Alert,
    CircularProgress,
    Snackbar,
    Tab,
    Tabs,
    Checkbox,
    Paper,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    InputAdornment,
} from '@mui/material';
import {
    ArrowBack,
    Delete,
    Add,
    UploadFile,
    PersonAdd,
    PersonRemove,
    Search,
} from '@mui/icons-material';
import type { DeleteShiftEmployeeReq, Shift } from '@face-attendance/types';
import {
    addEmployeeListToShift,
    deleteEmployeeShift,
    dateStringToTimestamp,
    getShifts,
    getShiftDetail,
    getEmployeesInShift,
    getEmployeesNotInShift,
} from '@face-attendance/utils';
import { Pagination } from '@mui/material';

interface Employee {
    id: string;
    name: string;
    employee_code: string;
    current_shift?: string;
    is_current_shift_active?: boolean;
}

interface AssignedEmployee extends Employee {
    shift_id: string;
    effective_from: string;
    effective_to?: string;
    is_active: boolean;
}

interface ParsedRow {
    [key: string]: string;
    _matchedEmployeeId?: string | null;
}

export const ShiftAssignmentPage: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams();
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [allEmployees, setAllEmployees] = useState<Employee[]>([]); // Chua phan cong
    const [assignedEmployees, setAssignedEmployees] = useState<
        AssignedEmployee[]
    >([]); // Da phan cong

    // Tab management
    const [currentTab, setCurrentTab] = useState(0); // 0: Chưa phân công, 1: Đã phân công

    // Selection management
    const [selectedUnassigned, setSelectedUnassigned] = useState<string[]>([]);
    const [selectedAssigned, setSelectedAssigned] = useState<string[]>([]);

    // Search
    const [searchUnassigned, setSearchUnassigned] = useState('');
    const [searchAssigned, setSearchAssigned] = useState('');
    // Pagination state for both tabs
    const [pageAssigned, setPageAssigned] = useState(1);
    const [pageUnassigned, setPageUnassigned] = useState(1);
    // Fixed backend page size = 20 (or fewer on last page)
    const sizeAssigned = 20;
    const sizeUnassigned = 20;
    const [totalAssigned, setTotalAssigned] = useState(0);
    const [totalUnassigned, setTotalUnassigned] = useState(0);

    // Form data for adding assignments
    const [formData, setFormData] = useState({
        shift_id: id || '',
        effective_from: new Date().toISOString().split('T')[0],
        effective_to: (() => {
            const d = new Date();
            d.setMonth(d.getMonth() + 3);
            return d.toISOString().split('T')[0];
        })(),
    });

    // Dialog
    const [deleteDialog, setDeleteDialog] = useState<{
        open: boolean;
        employeeIds: string[];
        employeeNames: string[];
    }>({
        open: false,
        employeeIds: [],
        employeeNames: [],
    });

    const [loading, setLoading] = useState(false);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({
        open: false,
        message: '',
        severity: 'success',
    });

    // CSV import dialog state
    const [importDialogOpen, setImportDialogOpen] = useState(false);
    const [csvPreview, setCsvPreview] = useState<ParsedRow[]>([]);
    const [csvHeaders, setCsvHeaders] = useState<string[]>([]);
    const [csvFileName, setCsvFileName] = useState<string>('');
    const [importProcessing, setImportProcessing] = useState(false);

    // Computed: Unassigned employees
    const unassignedEmployees = allEmployees.filter(
        (emp) => !assignedEmployees.some((assigned) => assigned.id === emp.id)
    );

    // Filter employees based on search
    const filteredUnassigned = unassignedEmployees.filter(
        (emp) =>
            emp.name.toLowerCase().includes(searchUnassigned.toLowerCase()) ||
            emp.employee_code
                .toLowerCase()
                .includes(searchUnassigned.toLowerCase())
    );

    const filteredAssigned = assignedEmployees.filter(
        (emp) =>
            emp.name.toLowerCase().includes(searchAssigned.toLowerCase()) ||
            emp.employee_code
                .toLowerCase()
                .includes(searchAssigned.toLowerCase())
    );

    // Normalize various paginated response shapes from backend
    const extractPaginated = (payload: any) => {
        if (Array.isArray(payload)) {
            return {
                items: payload,
                total: payload.length,
                page: 1,
                size: payload.length,
                total_pages: 1,
            };
        }
        const d = payload ?? {};
        const items =
            d.employees || d.items || d.data || d.records || d.list || [];
        const total =
            d.total ??
            d.total_items ??
            d.totalCount ??
            d.total_records ??
            items.length;
        const page = d.page ?? d.current_page ?? d.page_index ?? 1;
        const size = d.size ?? d.per_page ?? d.page_size ?? items.length;
        const total_pages =
            d.total_pages ??
            d.totalPages ??
            (size ? Math.ceil(total / size) : 1);
        return { items, total, page, size, total_pages };
    };

    // Format date 2025-11-14T00:00:00Z -> dd/mm/yyyy
    const formatDate = (dateStr: string) => {
        if (!dateStr || dateStr === '' || '0001-01-01T00:00:00Z' === dateStr) {
            return '';
        }
        const date = new Date(dateStr);
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();
        return `${day}/${month}/${year}`;
    };

    // Fetch employees assigned to shift (in + not in) using new endpoints
    const fetchEmployeesForShift = async (shiftId: string) => {
        setLoading(true);
        try {
            const [inRes, notInRes] = await Promise.all([
                getEmployeesInShift({
                    shift_id: shiftId,
                    page: pageAssigned,
                    size: sizeAssigned,
                }),
                getEmployeesNotInShift({
                    shift_id: shiftId,
                    page: pageUnassigned,
                    size: sizeUnassigned,
                }),
            ]);
            // Assigned (in)
            if (inRes.code === 200 && inRes.data) {
                const p = extractPaginated(inRes.data);
                const assigned: AssignedEmployee[] = (p.items || []).map(
                    (emp: any) => ({
                        id: emp.employee_id ?? '',
                        name: emp.employee_name ?? '',
                        employee_code: emp.employee_code ?? '',
                        shift_id: shiftId,
                        effective_from: formatDate(
                            emp.shift_effective_from ?? undefined
                        ),
                        effective_to: formatDate(
                            emp.shift_effective_to ?? undefined
                        ),
                        is_active: emp.employee_shift_active ?? true,
                    })
                );
                setAssignedEmployees(assigned);
                setTotalAssigned(p.total ?? assigned.length);
            } else {
                setAssignedEmployees([]);
                setTotalAssigned(0);
            }

            // Unassigned (not_in)
            if (notInRes.code === 200 && notInRes.data) {
                const p = extractPaginated(notInRes.data);
                const unassigned: Employee[] = (p.items || []).map(
                    (emp: any) => ({
                        id: emp.employee_id ?? emp.user_id ?? emp.id,
                        name: emp.employee_name ?? emp.user_name ?? emp.name,
                        employee_code:
                            emp.employee_code ??
                            emp.number_employee ??
                            emp.code,
                        current_shift:
                            emp.employee_shift_name ??
                            emp.shift_active ??
                            emp.current_shift ??
                            undefined,
                        is_current_shift_active:
                            emp.employee_shift_active ?? false,
                    })
                );
                setAllEmployees(unassigned);
                setTotalUnassigned(p.total ?? unassigned.length);
            } else {
                setAllEmployees([]);
                setTotalUnassigned(0);
            }
        } catch (err) {
            console.error('Error fetching employees:', err);
            setError('Không thể tải danh sách nhân viên');
        } finally {
            setLoading(false);
        }
    };

    // Fetch list of shifts (when no shiftId param is provided)
    const fetchShifts = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await getShifts();
            if (response.code === 200 && response.data) {
                setShifts(response.data);
            } else {
                setError(
                    response.message || 'Không thể tải danh sách ca làm việc'
                );
            }
        } catch (err: any) {
            console.error('Error fetching shifts:', err);
            setError(
                err.response?.data?.error || 'Đã xảy ra lỗi khi tải dữ liệu'
            );
        } finally {
            setLoading(false);
        }
    };

    const fetchShiftInfo = async (shiftId: string) => {
        try {
            setLoading(true);
            setError(null);
            const response = await getShiftDetail(shiftId);
            if (response.code === 200 && response.data) {
                const shiftDetail: Shift = response.data;
                setShifts([shiftDetail]);
            } else {
                setError(
                    response.message || 'Không thể tải thông tin ca làm việc'
                );
            }
        } catch (err: any) {
            console.error('Error fetching shift detail:', err);
            setError(
                err.response?.data?.error || 'Đã xảy ra lỗi khi tải dữ liệu'
            );
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!id) {
            fetchShifts();
        } else {
            fetchEmployeesForShift(id);
            fetchShiftInfo(id);
        }
    }, [id]);

    // When user selects a shift from dropdown (route without id), fetch data
    useEffect(() => {
        if (!id && formData.shift_id) {
            // Reset pages when shift changes
            setPageAssigned(1);
            setPageUnassigned(1);
            fetchEmployeesForShift(formData.shift_id);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [formData.shift_id]);

    // Refetch when pagination changes
    useEffect(() => {
        const targetShiftId = id || formData.shift_id;
        if (targetShiftId) {
            fetchEmployeesForShift(targetShiftId);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [pageAssigned, pageUnassigned, sizeAssigned, sizeUnassigned]);

    // Clear selections when page changes
    useEffect(() => {
        setSelectedAssigned([]);
    }, [pageAssigned, sizeAssigned]);
    useEffect(() => {
        setSelectedUnassigned([]);
    }, [pageUnassigned, sizeUnassigned]);

    // Handle select/deselect unassigned employees
    const handleSelectUnassigned = (employeeId: string) => {
        setSelectedUnassigned((prev) =>
            prev.includes(employeeId)
                ? prev.filter((id) => id !== employeeId)
                : [...prev, employeeId]
        );
    };

    // Simple CSV parser (handles basic quoted values)
    const parseCSV = (
        text: string
    ): { headers: string[]; rows: string[][] } => {
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
                    if (cur.trim() !== '') {
                        lines.push(cur);
                    }
                    cur = '';
                    // skip if next is \n in CRLF
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
                    } else {
                        quoted = !quoted;
                    }
                } else if (ch === ',' && !quoted) {
                    cells.push(cell.trim());
                    cell = '';
                } else {
                    cell += ch;
                }
            }
            cells.push(cell.trim());
            return cells;
        });

        const headers =
            rows.length > 0
                ? rows[0].map((h) => h.replace(/^"|"$/g, '').trim())
                : [];
        const dataRows = rows
            .slice(1)
            .map((r) => r.map((c) => c.replace(/^"|"$/g, '').trim()));
        return { headers, rows: dataRows };
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
            const rows: ParsedRow[] = parsed.rows.map((r) => {
                const obj: ParsedRow = {} as ParsedRow;
                parsed.headers.forEach((h, idx) => {
                    obj[h] = r[idx] ?? '';
                });
                return obj;
            });

            // Try to match rows to existing employees by employee_code or email
            const enriched = rows.map((row) => {
                const code = (row['employee_code'] || row['employee'] || '')
                    .toString()
                    .trim();
                const email = (row['email'] || '').toString().trim();
                let matched: Employee | undefined = undefined;
                if (code)
                    matched = allEmployees.find(
                        (e) =>
                            (e.employee_code || '').toLowerCase() ===
                            code.toLowerCase()
                    );
                if (!matched && email)
                    matched = allEmployees.find(
                        (e) =>
                            (e.name || '').toLowerCase() === email.toLowerCase()
                    );
                return {
                    ...row,
                    _matchedEmployeeId: matched ? matched.id : null,
                } as ParsedRow;
            });

            setCsvPreview(enriched);
        };
        reader.readAsText(file, 'utf-8');
    };

    const handleConfirmImport = async () => {
        if (!formData.shift_id) {
            setSnackbar({
                open: true,
                message: 'Vui lòng chọn ca làm việc trước khi import',
                severity: 'error',
            });
            return;
        }
        const matchedIds = csvPreview
            .map((r) => r._matchedEmployeeId)
            .filter(Boolean) as string[];
        if (matchedIds.length === 0) {
            setSnackbar({
                open: true,
                message:
                    'Không có nhân viên hợp lệ để thêm. Hãy đảm bảo file chứa `employee_code` khớp.',
                severity: 'error',
            });
            return;
        }

        try {
            setImportProcessing(true);
            const companyId = id
                ? shifts.find((shift) => shift.shift_id === id)?.company_id
                : '';
            const requestData = {
                company_id: companyId,
                shift_id: formData.shift_id,
                employee_ids: matchedIds,
                effective_from: dateStringToTimestamp(formData.effective_from),
                effective_to: formData.effective_to
                    ? dateStringToTimestamp(formData.effective_to)
                    : dateStringToTimestamp(
                          new Date(
                              new Date().setFullYear(
                                  new Date().getFullYear() + 10
                              )
                          )
                              .toISOString()
                              .split('T')[0]
                      ),
            };
            const response = await addEmployeeListToShift(requestData as any);
            if (response.code !== 200) {
                throw new Error(response.message || 'Import thất bại');
            }

            setSnackbar({
                open: true,
                message: `Đã thêm ${matchedIds.length} nhân viên từ file`,
                severity: 'success',
            });
            setImportDialogOpen(false);
            setCsvPreview([]);
            setCsvHeaders([]);
            // Refresh data
            if (id) {
                await fetchEmployeesForShift(id);
            } else if (formData.shift_id) {
                await fetchEmployeesForShift(formData.shift_id);
            }
        } catch (err: any) {
            console.error('Import error:', err);
            setSnackbar({
                open: true,
                message: err.message || 'Lỗi khi import',
                severity: 'error',
            });
        } finally {
            setImportProcessing(false);
        }
    };

    const handleSelectAllUnassigned = () => {
        if (selectedUnassigned.length === filteredUnassigned.length) {
            setSelectedUnassigned([]);
        } else {
            setSelectedUnassigned(filteredUnassigned.map((emp) => emp.id));
        }
    };

    // Handle select/deselect assigned employees
    const handleSelectAssigned = (employeeId: string) => {
        setSelectedAssigned((prev) =>
            prev.includes(employeeId)
                ? prev.filter((id) => id !== employeeId)
                : [...prev, employeeId]
        );
    };

    const handleSelectAllAssigned = () => {
        if (selectedAssigned.length === filteredAssigned.length) {
            setSelectedAssigned([]);
        } else {
            setSelectedAssigned(filteredAssigned.map((emp) => emp.id));
        }
    };

    // Add selected employees to shift
    const handleAddEmployeesToShift = async () => {
        if (selectedUnassigned.length === 0 || !formData.shift_id) {
            setSnackbar({
                open: true,
                message: 'Vui lòng chọn ít nhất một nhân viên và ca làm việc',
                severity: 'error',
            });
            return;
        }

        try {
            setLoading(true);

            // Get company id from cur shift
            const companyId = id
                ? shifts.find((shift) => shift.shift_id === id)?.company_id
                : '';
            const requestData = {
                company_id: companyId,
                shift_id: formData.shift_id,
                employee_ids: selectedUnassigned,
                effective_from: dateStringToTimestamp(formData.effective_from),
                effective_to: formData.effective_to
                    ? dateStringToTimestamp(formData.effective_to)
                    : dateStringToTimestamp(
                          new Date(
                              new Date().setFullYear(
                                  new Date().getFullYear() + 10
                              )
                          )
                              .toISOString()
                              .split('T')[0]
                      ),
            };
            const response = await addEmployeeListToShift(requestData);

            if (response.code !== 200) {
                throw new Error(response.message || 'Phân công thất bại');
            }

            setSnackbar({
                open: true,
                message: `Đã thêm ${selectedUnassigned.length} nhân viên vào ca làm việc`,
                severity: 'success',
            });

            // Refresh data
            setSelectedUnassigned([]);
            if (id) {
                await fetchEmployeesForShift(id);
            } else if (formData.shift_id) {
                await fetchEmployeesForShift(formData.shift_id);
            }
        } catch (err: any) {
            console.error('Error adding employees:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.error ||
                    'Đã xảy ra lỗi khi phân công ca làm việc',
                severity: 'error',
            });
        } finally {
            setLoading(false);
        }
    };

    // Remove employees from shift
    const handleRemoveEmployeesFromShift = () => {
        const employeesToRemove = assignedEmployees.filter((emp) =>
            selectedAssigned.includes(emp.id)
        );
        setDeleteDialog({
            open: true,
            employeeIds: employeesToRemove.map((emp) => emp.id),
            employeeNames: employeesToRemove.map((emp) => emp.name),
        });
    };

    const confirmDelete = async () => {
        // Validate ids before sending
        const rawIds = deleteDialog.employeeIds || [];
        const ids = rawIds.filter(Boolean);
        if (ids.length === 0) {
            setSnackbar({
                open: true,
                message: 'Không có phân công hợp lệ để xóa',
                severity: 'error',
            });
            setDeleteDialog({
                open: false,
                employeeIds: [],
                employeeNames: [],
            });
            return;
        }
        const shiftId = id || formData.shift_id;
        if (!shiftId) {
            setSnackbar({
                open: true,
                message: 'Không tìm thấy ca làm việc để xóa phân công',
                severity: 'error',
            });
            setDeleteDialog({
                open: false,
                employeeIds: [],
                employeeNames: [],
            });
            return;
        }

        try {
            setLoading(true);
            const req: DeleteShiftEmployeeReq = {
                employee_ids: ids,
                shift_id: shiftId,
            };
            console.log('Delete request:', req);
            const response = await deleteEmployeeShift(req);
            if (response.code !== 200) {
                throw new Error(response.message || 'Xóa phân công thất bại');
            }

            setSnackbar({
                open: true,
                message: `Đã xóa ${ids.length} phân công`,
                severity: 'success',
            });

            // Refresh data
            setSelectedAssigned([]);
            if (id) {
                await fetchEmployeesForShift(id);
            } else if (formData.shift_id) {
                await fetchEmployeesForShift(formData.shift_id);
            }
        } catch (err: any) {
            console.error('Error deleting assignments:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.error ||
                    'Đã xảy ra lỗi khi xóa phân công',
                severity: 'error',
            });
        } finally {
            setLoading(false);
            setDeleteDialog({
                open: false,
                employeeIds: [],
                employeeNames: [],
            });
        }
    };

    const handleCloseSnackbar = () => {
        setSnackbar({ ...snackbar, open: false });
    };

    const selectedShift = id ? shifts.find((s) => s.shift_id === id) : null;

    return (
        <Box>
            <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                mb={3}
            >
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/shifts')}
                >
                    Quay lại
                </Button>
                <Typography variant="h5" fontWeight="bold">
                    Phân công ca làm việc
                    {selectedShift && ` - ${selectedShift.name}`}
                </Typography>
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            <Grid container spacing={3}>
                {/* Shift Selection and Date Range */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Grid container spacing={2}>
                                <Grid item xs={12} md={4}>
                                    <FormControl fullWidth required>
                                        <InputLabel>Ca làm việc</InputLabel>
                                        <Select
                                            value={formData.shift_id}
                                            onChange={(e) =>
                                                setFormData({
                                                    ...formData,
                                                    shift_id: e.target.value,
                                                })
                                            }
                                            label="Ca làm việc"
                                            disabled={!!id || loading}
                                        >
                                            {shifts.map((shift) => (
                                                <MenuItem
                                                    key={shift.shift_id}
                                                    value={shift.shift_id}
                                                >
                                                    <Box>
                                                        <Typography variant="body2">
                                                            {shift.name}
                                                        </Typography>
                                                        <Typography
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            {shift.start_time} -{' '}
                                                            {shift.end_time}
                                                        </Typography>
                                                    </Box>
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                </Grid>

                                <Grid item xs={12} md={4}>
                                    <TextField
                                        fullWidth
                                        label="Có hiệu lực từ"
                                        type="date"
                                        required
                                        InputLabelProps={{ shrink: true }}
                                        value={formData.effective_from}
                                        onChange={(e) =>
                                            setFormData({
                                                ...formData,
                                                effective_from: e.target.value,
                                            })
                                        }
                                        disabled={loading}
                                    />
                                </Grid>

                                <Grid item xs={12} md={4}>
                                    <TextField
                                        fullWidth
                                        label="Có hiệu lực đến"
                                        type="date"
                                        InputLabelProps={{ shrink: true }}
                                        value={formData.effective_to}
                                        onChange={(e) =>
                                            setFormData({
                                                ...formData,
                                                effective_to: e.target.value,
                                            })
                                        }
                                        disabled={loading}
                                    />
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Employee Management Tabs */}
                <Grid item xs={12}>
                    <Card>
                        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                            <Tabs
                                value={currentTab}
                                onChange={(_, newValue) =>
                                    setCurrentTab(newValue)
                                }
                            >
                                <Tab
                                    label={`Chưa phân công (${totalUnassigned})`}
                                    icon={<PersonAdd />}
                                    iconPosition="start"
                                />
                                <Tab
                                    label={`Đã phân công (${totalAssigned})`}
                                    icon={<PersonRemove />}
                                    iconPosition="start"
                                />
                            </Tabs>
                        </Box>

                        {/* Tab 0: Unassigned Employees */}
                        {currentTab === 0 && (
                            <CardContent>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <TextField
                                        size="small"
                                        placeholder="Tìm kiếm nhân viên..."
                                        value={searchUnassigned}
                                        onChange={(e) =>
                                            setSearchUnassigned(e.target.value)
                                        }
                                        InputProps={{
                                            startAdornment: (
                                                <InputAdornment position="start">
                                                    <Search />
                                                </InputAdornment>
                                            ),
                                        }}
                                        sx={{ width: 300 }}
                                    />
                                    <Box display="flex" gap={1}>
                                        <Chip
                                            label={`Mỗi trang: ${sizeUnassigned}`}
                                            color="default"
                                        />
                                        <Chip
                                            label={`Đã chọn: ${selectedUnassigned.length}`}
                                            color={
                                                selectedUnassigned.length > 0
                                                    ? 'primary'
                                                    : 'default'
                                            }
                                        />
                                        <Button
                                            variant="contained"
                                            startIcon={
                                                loading ? (
                                                    <CircularProgress
                                                        size={20}
                                                    />
                                                ) : (
                                                    <Add />
                                                )
                                            }
                                            onClick={handleAddEmployeesToShift}
                                            disabled={
                                                selectedUnassigned.length ===
                                                    0 ||
                                                !formData.shift_id ||
                                                loading
                                            }
                                        >
                                            Thêm vào ca
                                        </Button>
                                        <Button
                                            variant="outlined"
                                            startIcon={<UploadFile />}
                                            onClick={openImportDialog}
                                            disabled={
                                                !formData.shift_id || loading
                                            }
                                        >
                                            Thêm từ file
                                        </Button>
                                    </Box>
                                </Box>

                                {filteredUnassigned.length === 0 ? (
                                    <Alert severity="info">
                                        Không có nhân viên nào chưa được phân
                                        công
                                    </Alert>
                                ) : (
                                    <TableContainer component={Paper}>
                                        <Table>
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell padding="checkbox">
                                                        <Checkbox
                                                            checked={
                                                                selectedUnassigned.length ===
                                                                    filteredUnassigned.length &&
                                                                filteredUnassigned.length >
                                                                    0
                                                            }
                                                            indeterminate={
                                                                selectedUnassigned.length >
                                                                    0 &&
                                                                selectedUnassigned.length <
                                                                    filteredUnassigned.length
                                                            }
                                                            onChange={
                                                                handleSelectAllUnassigned
                                                            }
                                                        />
                                                    </TableCell>
                                                    <TableCell>
                                                        Mã nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Tên nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Ca hiện tại
                                                    </TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {filteredUnassigned.map(
                                                    (employee) => (
                                                        <TableRow
                                                            key={employee.id}
                                                            hover
                                                            onClick={() =>
                                                                handleSelectUnassigned(
                                                                    employee.id
                                                                )
                                                            }
                                                            sx={{
                                                                cursor: 'pointer',
                                                            }}
                                                        >
                                                            <TableCell padding="checkbox">
                                                                <Checkbox
                                                                    checked={selectedUnassigned.includes(
                                                                        employee.id
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                {
                                                                    employee.employee_code
                                                                }
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.name}
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.current_shift ? (
                                                                    <Chip
                                                                        label={
                                                                            employee.current_shift
                                                                        }
                                                                        size="small"
                                                                        color="default"
                                                                    />
                                                                ) : (
                                                                    <Typography
                                                                        variant="caption"
                                                                        color="text.secondary"
                                                                    >
                                                                        Chưa có
                                                                        ca
                                                                    </Typography>
                                                                )}
                                                            </TableCell>
                                                        </TableRow>
                                                    )
                                                )}
                                            </TableBody>
                                        </Table>
                                    </TableContainer>
                                )}
                                {/* Pagination for Unassigned */}
                                {totalUnassigned > sizeUnassigned && (
                                    <Box
                                        display="flex"
                                        justifyContent="flex-end"
                                        mt={2}
                                    >
                                        <Pagination
                                            count={Math.max(
                                                1,
                                                Math.ceil(
                                                    totalUnassigned /
                                                        sizeUnassigned
                                                )
                                            )}
                                            page={pageUnassigned}
                                            onChange={(_, value) =>
                                                setPageUnassigned(value)
                                            }
                                            color="primary"
                                        />
                                    </Box>
                                )}
                            </CardContent>
                        )}

                        {/* Tab 1: Assigned Employees */}
                        {currentTab === 1 && (
                            <CardContent>
                                <Box
                                    display="flex"
                                    justifyContent="space-between"
                                    alignItems="center"
                                    mb={2}
                                >
                                    <TextField
                                        size="small"
                                        placeholder="Tìm kiếm nhân viên..."
                                        value={searchAssigned}
                                        onChange={(e) =>
                                            setSearchAssigned(e.target.value)
                                        }
                                        InputProps={{
                                            startAdornment: (
                                                <InputAdornment position="start">
                                                    <Search />
                                                </InputAdornment>
                                            ),
                                        }}
                                        sx={{ width: 300 }}
                                    />
                                    <Box display="flex" gap={1}>
                                        <Chip
                                            label={`Mỗi trang: ${sizeAssigned}`}
                                            color="default"
                                        />
                                        <Chip
                                            label={`Đã chọn: ${selectedAssigned.length}`}
                                            color={
                                                selectedAssigned.length > 0
                                                    ? 'error'
                                                    : 'default'
                                            }
                                        />
                                        <Button
                                            variant="contained"
                                            color="error"
                                            startIcon={
                                                loading ? (
                                                    <CircularProgress
                                                        size={20}
                                                    />
                                                ) : (
                                                    <Delete />
                                                )
                                            }
                                            onClick={
                                                handleRemoveEmployeesFromShift
                                            }
                                            disabled={
                                                selectedAssigned.length === 0 ||
                                                loading
                                            }
                                        >
                                            Xóa khỏi ca
                                        </Button>
                                    </Box>
                                </Box>

                                {filteredAssigned.length === 0 ? (
                                    <Alert severity="info">
                                        Chưa có nhân viên nào được phân công vào
                                        ca này
                                    </Alert>
                                ) : (
                                    <TableContainer component={Paper}>
                                        <Table>
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell padding="checkbox">
                                                        <Checkbox
                                                            checked={
                                                                selectedAssigned.length ===
                                                                    filteredAssigned.length &&
                                                                filteredAssigned.length >
                                                                    0
                                                            }
                                                            indeterminate={
                                                                selectedAssigned.length >
                                                                    0 &&
                                                                selectedAssigned.length <
                                                                    filteredAssigned.length
                                                            }
                                                            onChange={
                                                                handleSelectAllAssigned
                                                            }
                                                        />
                                                    </TableCell>
                                                    <TableCell>
                                                        Mã nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Tên nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Có hiệu lực từ
                                                    </TableCell>
                                                    <TableCell>
                                                        Có hiệu lực đến
                                                    </TableCell>
                                                    <TableCell>
                                                        Trạng thái
                                                    </TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {filteredAssigned.map(
                                                    (employee) => (
                                                        <TableRow
                                                            key={employee.id}
                                                            hover
                                                            onClick={() =>
                                                                handleSelectAssigned(
                                                                    employee.id
                                                                )
                                                            }
                                                            sx={{
                                                                cursor: 'pointer',
                                                            }}
                                                        >
                                                            <TableCell padding="checkbox">
                                                                <Checkbox
                                                                    checked={selectedAssigned.includes(
                                                                        employee.id
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                {
                                                                    employee.employee_code
                                                                }
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.name}
                                                            </TableCell>
                                                            <TableCell>
                                                                {
                                                                    employee.effective_from
                                                                }
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.effective_to ||
                                                                    'Không giới hạn'}
                                                            </TableCell>
                                                            <TableCell>
                                                                <Chip
                                                                    label={
                                                                        employee.is_active
                                                                            ? 'Đang hoạt động'
                                                                            : 'Không hoạt động'
                                                                    }
                                                                    size="small"
                                                                    color={
                                                                        employee.is_active
                                                                            ? 'success'
                                                                            : 'default'
                                                                    }
                                                                />
                                                            </TableCell>
                                                        </TableRow>
                                                    )
                                                )}
                                            </TableBody>
                                        </Table>
                                    </TableContainer>
                                )}
                                {/* Pagination for Assigned */}
                                {totalAssigned > sizeAssigned && (
                                    <Box
                                        display="flex"
                                        justifyContent="flex-end"
                                        mt={2}
                                    >
                                        <Pagination
                                            count={Math.max(
                                                1,
                                                Math.ceil(
                                                    totalAssigned / sizeAssigned
                                                )
                                            )}
                                            page={pageAssigned}
                                            onChange={(_, value) =>
                                                setPageAssigned(value)
                                            }
                                            color="primary"
                                        />
                                    </Box>
                                )}
                            </CardContent>
                        )}
                    </Card>
                </Grid>
            </Grid>

            {/* Import from CSV Dialog */}
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
                            File mẫu: các cột (order không bắt buộc):
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
                                                {row._matchedEmployeeId ? (
                                                    <Chip
                                                        label="Tìm thấy"
                                                        size="small"
                                                        color="success"
                                                    />
                                                ) : (
                                                    <Chip
                                                        label="Không tìm thấy"
                                                        size="small"
                                                        color="default"
                                                    />
                                                )}
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

            {/* Delete Confirmation Dialog */}
            <Dialog
                open={deleteDialog.open}
                onClose={() =>
                    setDeleteDialog({
                        open: false,
                        employeeIds: [],
                        employeeNames: [],
                    })
                }
            >
                <DialogTitle>Xác nhận xóa phân công</DialogTitle>
                <DialogContent>
                    <Typography>
                        Bạn có chắc chắn muốn xóa phân công của{' '}
                        {deleteDialog.employeeNames.length} nhân viên sau không?
                    </Typography>
                    <Box mt={2}>
                        {deleteDialog.employeeNames.map((name, index) => (
                            <Chip
                                key={index}
                                label={name}
                                sx={{ mr: 1, mb: 1 }}
                            />
                        ))}
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={() =>
                            setDeleteDialog({
                                open: false,
                                employeeIds: [],
                                employeeNames: [],
                            })
                        }
                        disabled={loading}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={confirmDelete}
                        color="error"
                        variant="contained"
                        disabled={loading}
                        startIcon={
                            loading ? (
                                <CircularProgress size={20} />
                            ) : (
                                <Delete />
                            )
                        }
                    >
                        {loading ? 'Đang xóa...' : 'Xác nhận xóa'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar for notifications */}
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
        </Box>
    );
};
