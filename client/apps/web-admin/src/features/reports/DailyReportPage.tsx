import { useState, useEffect, useCallback, useRef } from 'react';
import {
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
    Chip,
    Grid,
    CircularProgress,
    Alert,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Snackbar,
    ToggleButton,
    ToggleButtonGroup,
} from '@mui/material';
import {
    Download,
    People,
    CheckCircle,
    Cancel,
    Schedule,
    ExitToApp,
    AccessTime,
    CalendarMonth,
    CalendarToday,
} from '@mui/icons-material';
import { apiClient } from '@face-attendance/utils';

// Types for Attendance Record from API
interface AttendanceRecord {
    CompanyID: string;
    YearMonth: string;
    RecordTime: string; // ISO format
    EmployeeID: string;
    DeviceID: string;
    RecordType?: number; // 0 = check-in, 1 = check-out
    VerificationMethod: string;
    VerificationScore: {
        Float: number;
    };
    FaceImageURL: string;
    LocationCoordinates: string;
    Metadata?: any;
    SyncStatus?: string;
    CreatedAt: string;
}

// Types for Daily Attendance Status API
interface DailyEmployee {
    employee_id: string;
    name: string;
    check_in: string | null; // HH:mm:ss format
    check_out: string | null; // HH:mm:ss format
    status: 'on_time' | 'late' | 'early_leave' | 'absent' | 'overtime';
    late_minutes: number;
    total_hours: number;
}

interface DailyStatistics {
    total_employees: number;
    checked_in: number;
    not_checked_in: number;
    on_time: number;
    late: number;
    early_leave: number;
    overtime: number;
}

interface DailyAttendanceStatusResponse {
    success: boolean;
    data: {
        date: string; // YYYY-MM-DD
        attendance_records: AttendanceRecord[];
        daily_summaries: any;
        total_employees: number;
        total_records: number;
    };
}

// Types for Monthly Summary API
interface MonthlyEmployeeSummary {
    employee_id: string;
    name: string;
    present_days: number;
    absent_days: number;
    late_days: number;
    total_hours: number;
    average_hours_per_day: number;
    overtime_hours: number;
}

interface MonthlyStatistics {
    average_attendance_rate: number;
    total_late_instances: number;
    total_overtime_hours: number;
}

interface MonthlySummaryResponse {
    success: boolean;
    data: {
        month: string; // YYYY-MM
        absent_days: number;
        daily_summaries: any;
        early_leave_days: number;
        employee_statistics: Record<string, any>;
        late_days: number;
        present_days: number;
        records_no_shift: AttendanceRecord[];
        total_daily_summaries: number;
        total_early_leave_minutes: number;
        total_late_minutes: number;
        total_records_no_shift: number;
        total_work_minutes: number;
    };
}

// Export response type
interface ExportResponse {
    success: boolean;
    data: {
        file_id: string;
        download_url: string;
        expires_at: string;
    };
}

type ReportMode = 'daily' | 'monthly';

export const DailyReportPage: React.FC = () => {
    const [mode, setMode] = useState<ReportMode>('daily');
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
    const [month, setMonth] = useState(
        new Date().toISOString().slice(0, 7) // YYYY-MM format
    );
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Daily report states
    const [dailyData, setDailyData] = useState<{
        date: string;
        statistics: DailyStatistics;
        employees: DailyEmployee[];
        rawRecords: AttendanceRecord[];
    } | null>(null);

    // Monthly report states
    const [monthlyData, setMonthlyData] = useState<{
        month: string;
        total_working_days: number;
        employees_summary: MonthlyEmployeeSummary[];
        statistics: MonthlyStatistics;
        rawRecords: AttendanceRecord[];
    } | null>(null);

    // Employee names cache
    const [employeeNames, setEmployeeNames] = useState<Record<string, string>>(
        {}
    );
    const employeeNamesRef = useRef<Record<string, string>>({});

    // Export states
    const [exportDialogOpen, setExportDialogOpen] = useState(false);
    const [exportFormat, setExportFormat] = useState<'excel' | 'pdf' | 'csv'>(
        'excel'
    );
    const [exportEmail, setExportEmail] = useState('');
    const [exportLoading, setExportLoading] = useState(false);
    const [snackbar, setSnackbar] = useState<{
        open: boolean;
        message: string;
        severity: 'success' | 'error';
    }>({ open: false, message: '', severity: 'success' });

    // Get company ID from JWT token
    const getCompanyIdFromToken = (token: string): string | null => {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return payload.company_id || null;
        } catch {
            return null;
        }
    };

    // Fetch employee names
    const fetchEmployeeNames = useCallback(
        async (companyId: string, employeeIds: string[]) => {
            try {
                const response = await apiClient.get(
                    `/api/v1/users?company_id=${companyId}`
                );
                if (response.data?.data) {
                    const namesMap: Record<string, string> = {};
                    response.data.data.forEach((emp: any) => {
                        if (employeeIds.includes(emp.user_id)) {
                            namesMap[emp.user_id] =
                                emp.full_name || emp.name || 'N/A';
                        }
                    });
                    setEmployeeNames((prev) => {
                        const next = { ...prev, ...namesMap };
                        employeeNamesRef.current = next;
                        return next;
                    });
                    return namesMap;
                }
            } catch (err) {
                console.error('Failed to fetch employee names:', err);
            }
            return {};
        },
        []
    );

    // Process daily attendance records
    const processDailyRecords = (
        records: AttendanceRecord[],
        employeeNamesMap: Record<string, string>
    ): {
        statistics: DailyStatistics;
        employees: DailyEmployee[];
    } => {
        // Group records by employee
        const employeeRecords: Record<
            string,
            { checkIn?: AttendanceRecord; checkOut?: AttendanceRecord }
        > = {};

        records.forEach((record) => {
            const employeeId = record.EmployeeID;
            if (!employeeRecords[employeeId]) {
                employeeRecords[employeeId] = {};
            }

            const recordTime = new Date(record.RecordTime);
            const recordType = record.RecordType ?? 0; // 0 = check-in, 1 = check-out

            if (recordType === 0) {
                // Check-in - keep the earliest one
                if (
                    !employeeRecords[employeeId].checkIn ||
                    new Date(employeeRecords[employeeId].checkIn!.RecordTime) >
                        recordTime
                ) {
                    employeeRecords[employeeId].checkIn = record;
                }
            } else {
                // Check-out - keep the latest one
                if (
                    !employeeRecords[employeeId].checkOut ||
                    new Date(employeeRecords[employeeId].checkOut!.RecordTime) <
                        recordTime
                ) {
                    employeeRecords[employeeId].checkOut = record;
                }
            }
        });

        // Convert to display format
        const employees: DailyEmployee[] = Object.entries(employeeRecords).map(
            ([employeeId, records]) => {
                const checkInTime = records.checkIn
                    ? new Date(records.checkIn.RecordTime)
                    : null;
                const checkOutTime = records.checkOut
                    ? new Date(records.checkOut.RecordTime)
                    : null;

                const checkInStr = checkInTime
                    ? checkInTime.toLocaleTimeString('vi-VN', {
                          hour: '2-digit',
                          minute: '2-digit',
                          second: '2-digit',
                      })
                    : null;
                const checkOutStr = checkOutTime
                    ? checkOutTime.toLocaleTimeString('vi-VN', {
                          hour: '2-digit',
                          minute: '2-digit',
                          second: '2-digit',
                      })
                    : null;

                // Calculate total hours
                let totalHours = 0;
                if (checkInTime && checkOutTime) {
                    const diffMs =
                        checkOutTime.getTime() - checkInTime.getTime();
                    totalHours = diffMs / (1000 * 60 * 60);
                }

                // Determine status (simplified - would need shift data for accurate status)
                let status: DailyEmployee['status'] = 'on_time';
                let lateMinutes = 0;

                if (checkInTime) {
                    // Assume work starts at 8:00 AM
                    const workStart = new Date(checkInTime);
                    workStart.setHours(8, 0, 0, 0);
                    if (checkInTime > workStart) {
                        const diffMs =
                            checkInTime.getTime() - workStart.getTime();
                        lateMinutes = Math.floor(diffMs / (1000 * 60));
                        if (lateMinutes > 0) {
                            status = 'late';
                        }
                    }
                } else {
                    status = 'absent';
                }

                if (checkOutTime && checkInTime) {
                    // Assume work ends at 17:00 (5 PM)
                    const workEnd = new Date(checkOutTime);
                    workEnd.setHours(17, 0, 0, 0);
                    if (
                        checkOutTime < workEnd &&
                        checkOutTime.getHours() < 17
                    ) {
                        status = 'early_leave';
                    }
                }

                return {
                    employee_id: employeeId,
                    name: employeeNamesMap[employeeId] || 'N/A',
                    check_in: checkInStr,
                    check_out: checkOutStr,
                    status,
                    late_minutes: lateMinutes,
                    total_hours: totalHours,
                };
            }
        );

        // Calculate statistics
        const statistics: DailyStatistics = {
            total_employees: employees.length,
            checked_in: employees.filter((e) => e.check_in !== null).length,
            not_checked_in: employees.filter((e) => e.check_in === null).length,
            on_time: employees.filter((e) => e.status === 'on_time').length,
            late: employees.filter((e) => e.status === 'late').length,
            early_leave: employees.filter((e) => e.status === 'early_leave')
                .length,
            overtime: employees.filter((e) => e.status === 'overtime').length,
        };

        return { statistics, employees };
    };

    // Process monthly records
    const processMonthlyRecords = (
        data: MonthlySummaryResponse['data'],
        employeeNamesMap: Record<string, string>
    ): {
        month: string;
        total_working_days: number;
        employees_summary: MonthlyEmployeeSummary[];
        statistics: MonthlyStatistics;
    } => {
        // Group records by employee
        const employeeStats: Record<
            string,
            {
                records: AttendanceRecord[];
                presentDays: Set<string>;
                lateCount: number;
            }
        > = {};

        data.records_no_shift.forEach((record) => {
            const employeeId = record.EmployeeID;
            if (!employeeStats[employeeId]) {
                employeeStats[employeeId] = {
                    records: [],
                    presentDays: new Set(),
                    lateCount: 0,
                };
            }

            employeeStats[employeeId].records.push(record);
            const recordDate = new Date(record.RecordTime)
                .toISOString()
                .split('T')[0];
            employeeStats[employeeId].presentDays.add(recordDate);
        });

        // Calculate total working days (estimate from month)
        const [year, month] = data.month.split('-').map(Number);
        const daysInMonth = new Date(year, month, 0).getDate();
        // Estimate working days (excluding weekends, ~22 days per month)
        const totalWorkingDays = Math.floor((daysInMonth / 7) * 5);

        // Convert to display format
        const employees_summary: MonthlyEmployeeSummary[] = Object.entries(
            employeeStats
        ).map(([employeeId, stats]) => {
            const presentDays = stats.presentDays.size;
            // Estimate absent days
            const absentDays = Math.max(0, totalWorkingDays - presentDays);

            // Calculate total hours (simplified)
            let totalHours = 0;
            const recordsByDate: Record<string, AttendanceRecord[]> = {};
            stats.records.forEach((record) => {
                const date = new Date(record.RecordTime)
                    .toISOString()
                    .split('T')[0];
                if (!recordsByDate[date]) {
                    recordsByDate[date] = [];
                }
                recordsByDate[date].push(record);
            });

            Object.values(recordsByDate).forEach((dayRecords) => {
                // Sort records by time
                dayRecords.sort(
                    (a, b) =>
                        new Date(a.RecordTime).getTime() -
                        new Date(b.RecordTime).getTime()
                );

                // If RecordType is available, use it
                if (dayRecords[0].RecordType !== undefined) {
                    const checkIns = dayRecords.filter(
                        (r) => r.RecordType === 0
                    );
                    const checkOuts = dayRecords.filter(
                        (r) => r.RecordType === 1
                    );
                    if (checkIns.length > 0 && checkOuts.length > 0) {
                        const checkIn = new Date(
                            checkIns[0].RecordTime
                        ).getTime();
                        const checkOut = new Date(
                            checkOuts[checkOuts.length - 1].RecordTime
                        ).getTime();
                        const hours = (checkOut - checkIn) / (1000 * 60 * 60);
                        totalHours += hours;
                    } else if (checkIns.length > 0) {
                        // Assume 8 hours if only check-in
                        totalHours += 8;
                    }
                } else {
                    // If no RecordType, assume first record is check-in, last is check-out
                    if (dayRecords.length >= 2) {
                        const checkIn = new Date(
                            dayRecords[0].RecordTime
                        ).getTime();
                        const checkOut = new Date(
                            dayRecords[dayRecords.length - 1].RecordTime
                        ).getTime();
                        const hours = (checkOut - checkIn) / (1000 * 60 * 60);
                        // Only count if it's a reasonable work day (1-12 hours)
                        if (hours >= 1 && hours <= 12) {
                            totalHours += hours;
                        } else {
                            // Assume 8 hours if time seems unreasonable
                            totalHours += 8;
                        }
                    } else if (dayRecords.length === 1) {
                        // Only one record, assume 8 hours
                        totalHours += 8;
                    }
                }
            });

            const averageHoursPerDay =
                presentDays > 0 ? totalHours / presentDays : 0;

            return {
                employee_id: employeeId,
                name: employeeNamesMap[employeeId] || 'N/A',
                present_days: presentDays,
                absent_days: absentDays,
                late_days: stats.lateCount,
                total_hours: totalHours,
                average_hours_per_day: averageHoursPerDay,
                overtime_hours: 0, // Would need shift data to calculate
            };
        });

        // Calculate statistics
        const totalEmployees = employees_summary.length;
        const totalPresentDays = employees_summary.reduce(
            (sum, emp) => sum + emp.present_days,
            0
        );
        const averageAttendanceRate =
            totalEmployees > 0
                ? (totalPresentDays / (totalEmployees * totalWorkingDays)) * 100
                : 0;

        const statistics: MonthlyStatistics = {
            average_attendance_rate: averageAttendanceRate,
            total_late_instances: data.late_days,
            total_overtime_hours: data.total_work_minutes / 60, // Convert minutes to hours
        };

        return {
            month: data.month,
            total_working_days: totalWorkingDays,
            employees_summary,
            statistics,
        };
    };

    // Fetch daily attendance status
    const fetchDailyReport = useCallback(
        async (selectedDate: string) => {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) {
                setError('Không tìm thấy token xác thực');
                return;
            }

            const companyId = getCompanyIdFromToken(accessToken);
            if (!companyId) {
                setError('Không tìm thấy company_id trong token');
                return;
            }

            setLoading(true);
            setError(null);
            try {
                const response =
                    await apiClient.get<DailyAttendanceStatusResponse>(
                        '/api/v1/company/daily-attendance-status',
                        {
                            params: {
                                company_id: companyId,
                                date: selectedDate,
                            },
                        }
                    );

                if (response.data?.success && response.data?.data) {
                    const apiData = response.data.data;
                    const records = apiData.attendance_records || [];

                    // Extract unique employee IDs
                    const employeeIds = [
                        ...new Set(records.map((r) => r.EmployeeID)),
                    ];

                    // Fetch employee names
                    let namesMap = employeeNamesRef.current;
                    if (employeeIds.length > 0) {
                        const fetchedNames = await fetchEmployeeNames(
                            companyId,
                            employeeIds
                        );
                        namesMap = { ...namesMap, ...fetchedNames };
                    }

                    // Process records
                    const processed = processDailyRecords(records, namesMap);

                    setDailyData({
                        date: apiData.date,
                        rawRecords: records,
                        ...processed,
                    });
                } else {
                    setError('Không thể tải dữ liệu báo cáo');
                    setDailyData(null);
                }
            } catch (err: any) {
                console.error('Failed to fetch daily report:', err);
                setError(
                    err.response?.data?.message ||
                        err.response?.data?.error ||
                        'Không thể tải báo cáo. Vui lòng thử lại.'
                );
                setDailyData(null);
            } finally {
                setLoading(false);
            }
        },
        [fetchEmployeeNames]
    );

    // Fetch monthly summary
    const fetchMonthlyReport = useCallback(
        async (selectedMonth: string) => {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) {
                setError('Không tìm thấy token xác thực');
                return;
            }

            const companyId = getCompanyIdFromToken(accessToken);
            if (!companyId) {
                setError('Không tìm thấy company_id trong token');
                return;
            }

            setLoading(true);
            setError(null);
            try {
                const response = await apiClient.get<MonthlySummaryResponse>(
                    '/api/v1/company/monthly-summary',
                    {
                        params: {
                            company_id: companyId,
                            month: selectedMonth,
                        },
                    }
                );

                if (response.data?.success && response.data?.data) {
                    const apiData = response.data.data;
                    const records = apiData.records_no_shift || [];

                    // Extract unique employee IDs
                    const employeeIds = [
                        ...new Set(records.map((r) => r.EmployeeID)),
                    ];

                    // Fetch employee names
                    let namesMap = employeeNamesRef.current;
                    if (employeeIds.length > 0) {
                        const fetchedNames = await fetchEmployeeNames(
                            companyId,
                            employeeIds
                        );
                        namesMap = { ...namesMap, ...fetchedNames };
                    }

                    // Process records
                    const processed = processMonthlyRecords(apiData, namesMap);

                    setMonthlyData({ ...processed, rawRecords: records });
                } else {
                    setError('Không thể tải dữ liệu báo cáo');
                    setMonthlyData(null);
                }
            } catch (err: any) {
                console.error('Failed to fetch monthly report:', err);
                setError(
                    err.response?.data?.message ||
                        err.response?.data?.error ||
                        'Không thể tải báo cáo. Vui lòng thử lại.'
                );
                setMonthlyData(null);
            } finally {
                setLoading(false);
            }
        },
        [fetchEmployeeNames]
    );

    useEffect(() => {
        if (mode === 'daily') {
            fetchDailyReport(date);
        } else {
            fetchMonthlyReport(month);
        }
    }, [mode, date, month, fetchDailyReport, fetchMonthlyReport]);

    const handleModeChange = (
        _event: React.MouseEvent<HTMLElement>,
        newMode: ReportMode | null
    ) => {
        if (newMode !== null) {
            setMode(newMode);
        }
    };

    const handleDateChange = (newDate: string) => {
        setDate(newDate);
    };

    const handleMonthChange = (newMonth: string) => {
        setMonth(newMonth);
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'on_time':
                return 'success';
            case 'late':
                return 'warning';
            case 'early_leave':
                return 'info';
            case 'absent':
                return 'error';
            case 'overtime':
                return 'secondary';
            default:
                return 'default';
        }
    };

    const getStatusText = (status: string) => {
        switch (status) {
            case 'on_time':
                return 'Đúng giờ';
            case 'late':
                return 'Trễ';
            case 'early_leave':
                return 'Về sớm';
            case 'absent':
                return 'Vắng mặt';
            case 'overtime':
                return 'Tăng ca';
            default:
                return '-';
        }
    };

    const formatHours = (hours: number) => {
        if (!hours) return '-';
        const h = Math.floor(hours);
        const m = Math.round((hours - h) * 60);
        if (h > 0 && m > 0) {
            return `${h}h${m}m`;
        } else if (h > 0) {
            return `${h}h`;
        } else {
            return `${m}m`;
        }
    };

    const formatDateTime = (value: string) => {
        if (!value) return '-';
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return '-';
        return date.toLocaleString('vi-VN', {
            hour12: false,
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
        });
    };

    // Handle export report
    const handleExportReport = async () => {
        const accessToken = localStorage.getItem('access_token');
        if (!accessToken) {
            setSnackbar({
                open: true,
                message: 'Không tìm thấy token xác thực',
                severity: 'error',
            });
            return;
        }

        const companyId = getCompanyIdFromToken(accessToken);
        if (!companyId) {
            setSnackbar({
                open: true,
                message: 'Không tìm thấy company_id trong token',
                severity: 'error',
            });
            return;
        }

        setExportLoading(true);
        try {
            const endpoint =
                mode === 'daily'
                    ? '/api/v1/company/export-daily-status'
                    : '/api/v1/company/export-monthly-summary';

            const requestBody: {
                company_id: string;
                date?: string;
                month?: string;
                format: string;
                email?: string;
            } = {
                company_id: companyId,
                format: exportFormat,
            };

            if (mode === 'daily') {
                requestBody.date = date;
            } else {
                requestBody.month = month;
            }

            if (exportEmail.trim()) {
                requestBody.email = exportEmail.trim();
            }

            const response = await apiClient.post<ExportResponse>(
                endpoint,
                requestBody
            );

            if (response.data?.success && response.data?.data) {
                const { download_url } = response.data.data;

                if (exportEmail.trim()) {
                    // Email sent
                    setSnackbar({
                        open: true,
                        message: `Báo cáo đã được gửi đến email ${exportEmail}`,
                        severity: 'success',
                    });
                } else {
                    // Download file
                    try {
                        const downloadResponse = await apiClient.get(
                            download_url,
                            {
                                responseType: 'blob',
                            }
                        );

                        const blob = new Blob([downloadResponse.data], {
                            type:
                                downloadResponse.headers['content-type'] ||
                                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
                        });

                        const url = window.URL.createObjectURL(blob);
                        const link = document.createElement('a');
                        link.href = url;
                        const extension =
                            exportFormat === 'excel' ? 'xlsx' : exportFormat;
                        const fileName = `${mode}-report-${
                            mode === 'daily' ? date : month
                        }.${extension}`;
                        link.setAttribute('download', fileName);
                        document.body.appendChild(link);
                        link.click();
                        link.remove();
                        window.URL.revokeObjectURL(url);

                        setSnackbar({
                            open: true,
                            message: 'Xuất báo cáo thành công',
                            severity: 'success',
                        });
                    } catch (downloadErr: any) {
                        console.error('Failed to download file:', downloadErr);
                        setSnackbar({
                            open: true,
                            message:
                                'Xuất báo cáo thành công nhưng không thể tải xuống. Vui lòng thử lại.',
                            severity: 'error',
                        });
                    }
                }

                setExportDialogOpen(false);
                setExportEmail('');
            } else {
                setSnackbar({
                    open: true,
                    message: 'Không thể xuất báo cáo. Vui lòng thử lại.',
                    severity: 'error',
                });
            }
        } catch (err: any) {
            console.error('Failed to export report:', err);
            setSnackbar({
                open: true,
                message:
                    err.response?.data?.message ||
                    err.response?.data?.error ||
                    'Không thể xuất báo cáo. Vui lòng thử lại.',
                severity: 'error',
            });
        } finally {
            setExportLoading(false);
        }
    };

    const handleOpenExportDialog = () => {
        setExportDialogOpen(true);
        setExportFormat('excel');
        setExportEmail('');
    };

    const handleCloseExportDialog = () => {
        if (!exportLoading) {
            setExportDialogOpen(false);
            setExportEmail('');
        }
    };

    return (
        <Box>
            <Typography variant="h4" fontWeight="bold" mb={3}>
                Báo cáo Chấm công
            </Typography>

            {/* Mode Toggle */}
            <Card sx={{ mb: 3, p: 2 }}>
                <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} md={4}>
                        <FormControl fullWidth>
                            <ToggleButtonGroup
                                value={mode}
                                exclusive
                                onChange={handleModeChange}
                                aria-label="report mode"
                                fullWidth
                            >
                                <ToggleButton value="daily" aria-label="daily">
                                    <CalendarToday sx={{ mr: 1 }} />
                                    Theo ngày
                                </ToggleButton>
                                <ToggleButton
                                    value="monthly"
                                    aria-label="monthly"
                                >
                                    <CalendarMonth sx={{ mr: 1 }} />
                                    Theo tháng
                                </ToggleButton>
                            </ToggleButtonGroup>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12} md={4}>
                        {mode === 'daily' ? (
                            <TextField
                                fullWidth
                                label="Ngày"
                                type="date"
                                value={date}
                                onChange={(e) =>
                                    handleDateChange(e.target.value)
                                }
                                InputLabelProps={{ shrink: true }}
                            />
                        ) : (
                            <TextField
                                fullWidth
                                label="Tháng"
                                type="month"
                                value={month}
                                onChange={(e) =>
                                    handleMonthChange(e.target.value)
                                }
                                InputLabelProps={{ shrink: true }}
                            />
                        )}
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <Button
                            variant="contained"
                            startIcon={<Download />}
                            onClick={handleOpenExportDialog}
                            disabled={loading}
                            fullWidth
                        >
                            Xuất báo cáo
                        </Button>
                    </Grid>
                </Grid>
            </Card>

            {error && (
                <Alert
                    severity="error"
                    sx={{ mb: 2 }}
                    onClose={() => setError(null)}
                >
                    {error}
                </Alert>
            )}

            {loading ? (
                <Box
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    minHeight="200px"
                >
                    <CircularProgress />
                </Box>
            ) : (
                <>
                    {/* Daily Report View */}
                    {mode === 'daily' && dailyData && dailyData.statistics && (
                        <>
                            {/* Statistics Cards */}
                            <Grid container spacing={2} sx={{ mb: 3 }}>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Tổng nhân viên
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                    >
                                                        {dailyData.statistics
                                                            ?.total_employees ??
                                                            0}
                                                    </Typography>
                                                </Box>
                                                <People
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'primary.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đã chấm công
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="success.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.checked_in ?? 0}
                                                    </Typography>
                                                </Box>
                                                <CheckCircle
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'success.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Chưa chấm công
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="error.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.not_checked_in ??
                                                            0}
                                                    </Typography>
                                                </Box>
                                                <Cancel
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'error.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đúng giờ
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="success.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.on_time ?? 0}
                                                    </Typography>
                                                </Box>
                                                <CheckCircle
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'success.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Đi trễ
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="warning.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.late ?? 0}
                                                    </Typography>
                                                </Box>
                                                <Schedule
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'warning.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Về sớm
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="info.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.early_leave ?? 0}
                                                    </Typography>
                                                </Box>
                                                <ExitToApp
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'info.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={12} sm={6} md={2.4}>
                                    <Card>
                                        <CardContent>
                                            <Box
                                                display="flex"
                                                justifyContent="space-between"
                                                alignItems="center"
                                            >
                                                <Box>
                                                    <Typography
                                                        color="textSecondary"
                                                        variant="body2"
                                                    >
                                                        Tăng ca
                                                    </Typography>
                                                    <Typography
                                                        variant="h5"
                                                        fontWeight="bold"
                                                        color="secondary.main"
                                                    >
                                                        {dailyData.statistics
                                                            ?.overtime ?? 0}
                                                    </Typography>
                                                </Box>
                                                <AccessTime
                                                    sx={{
                                                        fontSize: 40,
                                                        color: 'secondary.main',
                                                    }}
                                                />
                                            </Box>
                                        </CardContent>
                                    </Card>
                                </Grid>
                            </Grid>

                            {/* Employees Table */}
                            <Card>
                                <TableContainer>
                                    <Table>
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Nhân viên</TableCell>
                                                <TableCell>Giờ vào</TableCell>
                                                <TableCell>Giờ ra</TableCell>
                                                <TableCell>
                                                    Tổng giờ làm
                                                </TableCell>
                                                <TableCell>
                                                    Trễ (phút)
                                                </TableCell>
                                                <TableCell>
                                                    Trạng thái
                                                </TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {!dailyData.employees ||
                                            dailyData.employees.length === 0 ? (
                                                <TableRow>
                                                    <TableCell
                                                        colSpan={6}
                                                        align="center"
                                                    >
                                                        <Typography
                                                            color="textSecondary"
                                                            py={2}
                                                        >
                                                            Không có dữ liệu cho
                                                            ngày đã chọn
                                                        </Typography>
                                                    </TableCell>
                                                </TableRow>
                                            ) : (
                                                dailyData.employees.map(
                                                    (employee) => (
                                                        <TableRow
                                                            key={
                                                                employee.employee_id
                                                            }
                                                        >
                                                            <TableCell>
                                                                <Typography fontWeight="bold">
                                                                    {
                                                                        employee.name
                                                                    }
                                                                </Typography>
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.check_in ||
                                                                    '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.check_out ||
                                                                    '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                {formatHours(
                                                                    employee.total_hours
                                                                )}
                                                            </TableCell>
                                                            <TableCell>
                                                                {employee.late_minutes >
                                                                0
                                                                    ? `${employee.late_minutes}`
                                                                    : '-'}
                                                            </TableCell>
                                                            <TableCell>
                                                                <Chip
                                                                    label={getStatusText(
                                                                        employee.status
                                                                    )}
                                                                    color={getStatusColor(
                                                                        employee.status
                                                                    )}
                                                                    size="small"
                                                                />
                                                            </TableCell>
                                                        </TableRow>
                                                    )
                                                )
                                            )}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </Card>

                            {/* Raw attendance records for the selected day */}
                            <Card sx={{ mt: 3 }}>
                                <CardContent>
                                    <Typography variant="h6" mb={2}>
                                        Bản ghi chấm công
                                    </Typography>
                                    <TableContainer>
                                        <Table size="small">
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell>
                                                        Nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Thời gian
                                                    </TableCell>
                                                    <TableCell>Loại</TableCell>
                                                    <TableCell>
                                                        Thiết bị
                                                    </TableCell>
                                                    <TableCell>
                                                        Phương thức
                                                    </TableCell>
                                                    <TableCell align="right">
                                                        Điểm
                                                    </TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {!dailyData.rawRecords ||
                                                dailyData.rawRecords.length ===
                                                    0 ? (
                                                    <TableRow>
                                                        <TableCell
                                                            colSpan={6}
                                                            align="center"
                                                        >
                                                            <Typography
                                                                color="textSecondary"
                                                                py={2}
                                                            >
                                                                Không có bản ghi
                                                                cho ngày đã chọn
                                                            </Typography>
                                                        </TableCell>
                                                    </TableRow>
                                                ) : (
                                                    dailyData.rawRecords.map(
                                                        (record, idx) => {
                                                            const name =
                                                                employeeNames[
                                                                    record
                                                                        .EmployeeID
                                                                ] || 'N/A';
                                                            const typeLabel =
                                                                record.RecordType ===
                                                                0
                                                                    ? 'Giờ vào'
                                                                    : record.RecordType ===
                                                                      1
                                                                    ? 'Giờ ra'
                                                                    : '-';

                                                            return (
                                                                <TableRow
                                                                    key={`${record.EmployeeID}-${record.RecordTime}-${idx}`}
                                                                >
                                                                    <TableCell>
                                                                        {name}
                                                                    </TableCell>
                                                                    <TableCell>
                                                                        {formatDateTime(
                                                                            record.RecordTime
                                                                        )}
                                                                    </TableCell>
                                                                    <TableCell>
                                                                        {
                                                                            typeLabel
                                                                        }
                                                                    </TableCell>
                                                                    <TableCell>
                                                                        {record.DeviceID ||
                                                                            '-'}
                                                                    </TableCell>
                                                                    <TableCell>
                                                                        {record.VerificationMethod ||
                                                                            '-'}
                                                                    </TableCell>
                                                                    <TableCell align="right">
                                                                        {record.VerificationScore?.Float?.toFixed(
                                                                            2
                                                                        ) ??
                                                                            '-'}
                                                                    </TableCell>
                                                                </TableRow>
                                                            );
                                                        }
                                                    )
                                                )}
                                            </TableBody>
                                        </Table>
                                    </TableContainer>
                                </CardContent>
                            </Card>
                        </>
                    )}

                    {/* Monthly Report View */}
                    {mode === 'monthly' &&
                        monthlyData &&
                        monthlyData.statistics && (
                            <>
                                {/* Statistics Cards */}
                                <Grid container spacing={2} sx={{ mb: 3 }}>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng ngày làm việc
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                        >
                                                            {monthlyData.total_working_days ??
                                                                0}
                                                        </Typography>
                                                    </Box>
                                                    <CalendarMonth
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'primary.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tỷ lệ chấm công TB
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="success.main"
                                                        >
                                                            {(
                                                                monthlyData
                                                                    .statistics
                                                                    ?.average_attendance_rate ??
                                                                0
                                                            ).toFixed(1)}
                                                            %
                                                        </Typography>
                                                    </Box>
                                                    <CheckCircle
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'success.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng lần đi trễ
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="warning.main"
                                                        >
                                                            {monthlyData
                                                                .statistics
                                                                ?.total_late_instances ??
                                                                0}
                                                        </Typography>
                                                    </Box>
                                                    <Schedule
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'warning.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Card>
                                            <CardContent>
                                                <Box
                                                    display="flex"
                                                    justifyContent="space-between"
                                                    alignItems="center"
                                                >
                                                    <Box>
                                                        <Typography
                                                            color="textSecondary"
                                                            variant="body2"
                                                        >
                                                            Tổng giờ tăng ca
                                                        </Typography>
                                                        <Typography
                                                            variant="h5"
                                                            fontWeight="bold"
                                                            color="secondary.main"
                                                        >
                                                            {formatHours(
                                                                monthlyData
                                                                    .statistics
                                                                    ?.total_overtime_hours ??
                                                                    0
                                                            )}
                                                        </Typography>
                                                    </Box>
                                                    <AccessTime
                                                        sx={{
                                                            fontSize: 40,
                                                            color: 'secondary.main',
                                                        }}
                                                    />
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                </Grid>

                                {/* Employees Summary Table */}
                                <Card>
                                    <TableContainer>
                                        <Table>
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell>
                                                        Nhân viên
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày có mặt
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày vắng mặt
                                                    </TableCell>
                                                    <TableCell>
                                                        Ngày đi trễ
                                                    </TableCell>
                                                    <TableCell>
                                                        Tổng giờ làm
                                                    </TableCell>
                                                    <TableCell>
                                                        TB giờ/ngày
                                                    </TableCell>
                                                    <TableCell>
                                                        Giờ tăng ca
                                                    </TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {!monthlyData.employees_summary ||
                                                monthlyData.employees_summary
                                                    .length === 0 ? (
                                                    <TableRow>
                                                        <TableCell
                                                            colSpan={7}
                                                            align="center"
                                                        >
                                                            <Typography
                                                                color="textSecondary"
                                                                py={2}
                                                            >
                                                                Không có dữ liệu
                                                                cho tháng đã
                                                                chọn
                                                            </Typography>
                                                        </TableCell>
                                                    </TableRow>
                                                ) : (
                                                    monthlyData.employees_summary.map(
                                                        (employee) => (
                                                            <TableRow
                                                                key={
                                                                    employee.employee_id
                                                                }
                                                            >
                                                                <TableCell>
                                                                    <Typography fontWeight="bold">
                                                                        {
                                                                            employee.name
                                                                        }
                                                                    </Typography>
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.present_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.absent_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {
                                                                        employee.late_days
                                                                    }
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.total_hours
                                                                    )}
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.average_hours_per_day
                                                                    )}
                                                                </TableCell>
                                                                <TableCell>
                                                                    {formatHours(
                                                                        employee.overtime_hours
                                                                    )}
                                                                </TableCell>
                                                            </TableRow>
                                                        )
                                                    )
                                                )}
                                            </TableBody>
                                        </Table>
                                    </TableContainer>
                                </Card>

                                {/* Raw records without shift for the selected month */}
                                <Card sx={{ mt: 3 }}>
                                    <CardContent>
                                        <Typography variant="h6" mb={2}>
                                            Bản ghi không thuộc ca
                                        </Typography>
                                        <TableContainer>
                                            <Table size="small">
                                                <TableHead>
                                                    <TableRow>
                                                        <TableCell>
                                                            Nhân viên
                                                        </TableCell>
                                                        <TableCell>
                                                            Thời gian
                                                        </TableCell>
                                                        <TableCell>
                                                            Thiết bị
                                                        </TableCell>
                                                        <TableCell>
                                                            Phương thức
                                                        </TableCell>
                                                        <TableCell align="right">
                                                            Điểm
                                                        </TableCell>
                                                    </TableRow>
                                                </TableHead>
                                                <TableBody>
                                                    {!monthlyData.rawRecords ||
                                                    monthlyData.rawRecords
                                                        .length === 0 ? (
                                                        <TableRow>
                                                            <TableCell
                                                                colSpan={5}
                                                                align="center"
                                                            >
                                                                <Typography
                                                                    color="textSecondary"
                                                                    py={2}
                                                                >
                                                                    Không có bản
                                                                    ghi cho
                                                                    tháng đã
                                                                    chọn
                                                                </Typography>
                                                            </TableCell>
                                                        </TableRow>
                                                    ) : (
                                                        monthlyData.rawRecords.map(
                                                            (record, idx) => {
                                                                const name =
                                                                    employeeNames[
                                                                        record
                                                                            .EmployeeID
                                                                    ] || 'N/A';

                                                                return (
                                                                    <TableRow
                                                                        key={`${record.EmployeeID}-${record.RecordTime}-${idx}`}
                                                                    >
                                                                        <TableCell>
                                                                            {
                                                                                name
                                                                            }
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            {formatDateTime(
                                                                                record.RecordTime
                                                                            )}
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            {record.DeviceID ||
                                                                                '-'}
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            {record.VerificationMethod ||
                                                                                '-'}
                                                                        </TableCell>
                                                                        <TableCell align="right">
                                                                            {record.VerificationScore?.Float?.toFixed(
                                                                                2
                                                                            ) ??
                                                                                '-'}
                                                                        </TableCell>
                                                                    </TableRow>
                                                                );
                                                            }
                                                        )
                                                    )}
                                                </TableBody>
                                            </Table>
                                        </TableContainer>
                                    </CardContent>
                                </Card>
                            </>
                        )}

                    {/* No data message */}
                    {((mode === 'daily' && !dailyData) ||
                        (mode === 'monthly' && !monthlyData)) && (
                        <Card>
                            <Box
                                display="flex"
                                justifyContent="center"
                                alignItems="center"
                                minHeight="200px"
                            >
                                <Typography color="textSecondary">
                                    Không có dữ liệu để hiển thị
                                </Typography>
                            </Box>
                        </Card>
                    )}
                </>
            )}

            {/* Export Dialog */}
            <Dialog
                open={exportDialogOpen}
                onClose={handleCloseExportDialog}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>Xuất báo cáo chấm công</DialogTitle>
                <DialogContent>
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'column',
                            gap: 2,
                            pt: 2,
                        }}
                    >
                        <FormControl fullWidth>
                            <InputLabel>Định dạng</InputLabel>
                            <Select
                                value={exportFormat}
                                label="Định dạng"
                                onChange={(e) =>
                                    setExportFormat(
                                        e.target.value as
                                            | 'excel'
                                            | 'pdf'
                                            | 'csv'
                                    )
                                }
                            >
                                <MenuItem value="excel">Excel (.xlsx)</MenuItem>
                                <MenuItem value="pdf">PDF (.pdf)</MenuItem>
                                <MenuItem value="csv">CSV (.csv)</MenuItem>
                            </Select>
                        </FormControl>
                        <TextField
                            fullWidth
                            label="Email (tùy chọn)"
                            type="email"
                            value={exportEmail}
                            onChange={(e) => setExportEmail(e.target.value)}
                            placeholder="Nhập email để nhận báo cáo qua email"
                            helperText="Để trống nếu muốn tải xuống trực tiếp"
                        />
                        <Typography variant="body2" color="textSecondary">
                            {mode === 'daily'
                                ? `Ngày: ${date}`
                                : `Tháng: ${month}`}
                        </Typography>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={handleCloseExportDialog}
                        disabled={exportLoading}
                    >
                        Hủy
                    </Button>
                    <Button
                        onClick={handleExportReport}
                        variant="contained"
                        disabled={exportLoading}
                        startIcon={
                            exportLoading ? (
                                <CircularProgress size={20} />
                            ) : (
                                <Download />
                            )
                        }
                    >
                        {exportLoading ? 'Đang xử lý...' : 'Xuất báo cáo'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Snackbar for notifications */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={() => setSnackbar({ ...snackbar, open: false })}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={() => setSnackbar({ ...snackbar, open: false })}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};
