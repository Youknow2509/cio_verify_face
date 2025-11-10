import { apiClient } from './api';
import type {
    Shift,
    CreateShiftReq,
    EditShiftReq,
    ChangeStatusShiftReq,
    AddShiftEmployeeReq,
    AddShiftEmployeeListReq,
    ShiftEmployeeEffectiveDateReq,
    ShiftEmployeeEditEffectiveDateReq,
    EnableShiftForUserReq,
    DisableShiftForUserReq,
    ResponseData,
    EmployeeShift,
} from '@face-attendance/types';

// ==================== SHIFT MANAGEMENT APIs ====================

/**
 * Get list of shifts with pagination
 * @param page - Page number (optional)
 * @returns List of shifts
 */
export const getShifts = async (
    page?: number
): Promise<ResponseData<Shift[]>> => {
    const params = page ? { page } : {};
    const response = await apiClient.get('/api/v1/shift', { params });
    return response.data;
};

/**
 * Get shift detail by ID
 * @param shiftId - Shift ID
 * @returns Shift detail
 */
export const getShiftDetail = async (
    shiftId: string
): Promise<ResponseData<Shift>> => {
    const response = await apiClient.get(`/api/v1/shift/${shiftId}`);
    return response.data;
};

/**
 * Create a new shift
 * @param data - Shift creation data
 * @returns Created shift
 */
export const createShift = async (
    data: CreateShiftReq
): Promise<ResponseData<Shift>> => {
    const response = await apiClient.post('/api/v1/shift', data);
    return response.data;
};

/**
 * Edit an existing shift
 * @param data - Shift edit data
 * @returns Updated shift
 */
export const editShift = async (
    data: EditShiftReq
): Promise<ResponseData<Shift>> => {
    const response = await apiClient.post('/api/v1/shift/edit', data);
    return response.data;
};

/**
 * Delete a shift
 * @param shiftId - Shift ID to delete
 * @returns Success response
 */
export const deleteShift = async (shiftId: string): Promise<ResponseData> => {
    const response = await apiClient.delete(`/api/v1/shift/${shiftId}`);
    return response.data;
};

/**
 * Change shift status (active/inactive)
 * @param data - Status change data
 * @returns Updated shift
 */
export const changeShiftStatus = async (
    data: ChangeStatusShiftReq
): Promise<ResponseData<Shift>> => {
    const response = await apiClient.post('/api/v1/shift/status', data);
    return response.data;
};

// ==================== EMPLOYEE SHIFT ASSIGNMENT APIs ====================

/**
 * Get employee shifts by effective date range with pagination
 * @param data - Query parameters including date range and pagination
 * @returns List of employee shifts
 */
export const getEmployeeShifts = async (
    data: ShiftEmployeeEffectiveDateReq
): Promise<ResponseData<EmployeeShift[]>> => {
    const response = await apiClient.post('/api/v1/employee/shift', data);
    return response.data;
};

/**
 * Add a single employee to a shift
 * @param data - Employee shift assignment data
 * @returns Created employee shift
 */
export const addEmployeeToShift = async (
    data: AddShiftEmployeeReq
): Promise<ResponseData<EmployeeShift>> => {
    const response = await apiClient.post('/api/v1/employee/shift/add', data);
    return response.data;
};

/**
 * Add multiple employees to a shift (bulk assignment)
 * @param data - Bulk assignment data
 * @returns Created employee shifts
 */
export const addEmployeeListToShift = async (
    data: AddShiftEmployeeListReq
): Promise<ResponseData<EmployeeShift[]>> => {
    const response = await apiClient.post(
        '/api/v1/employee/shift/add/list',
        data
    );
    return response.data;
};

/**
 * Delete an employee shift assignment
 * @param employeeShiftId - Employee shift ID to delete
 * @returns Success response
 */
export const deleteEmployeeShift = async (
    employeeShiftId: string
): Promise<ResponseData> => {
    const response = await apiClient.delete(
        `/api/v1/employee/shift/${employeeShiftId}`
    );
    return response.data;
};

/**
 * Enable an employee shift assignment
 * @param data - Enable request data
 * @returns Updated employee shift
 */
export const enableEmployeeShift = async (
    data: EnableShiftForUserReq
): Promise<ResponseData<EmployeeShift>> => {
    const response = await apiClient.post(
        '/api/v1/employee/shift/enable',
        data
    );
    return response.data;
};

/**
 * Disable an employee shift assignment
 * @param data - Disable request data
 * @returns Updated employee shift
 */
export const disableEmployeeShift = async (
    data: DisableShiftForUserReq
): Promise<ResponseData<EmployeeShift>> => {
    const response = await apiClient.post(
        '/api/v1/employee/shift/disable',
        data
    );
    return response.data;
};

/**
 * Edit effective dates of an employee shift assignment
 * @param data - Edit effective date data
 * @returns Updated employee shift
 */
export const editEmployeeShiftEffectiveDate = async (
    data: ShiftEmployeeEditEffectiveDateReq
): Promise<ResponseData<EmployeeShift>> => {
    const response = await apiClient.post('/api/v1/user/edit/effective', data);
    return response.data;
};

// ==================== HELPER FUNCTIONS ====================

/**
 * Convert HH:MM time string to Unix timestamp for today
 * @param timeString - Time in HH:MM format
 * @returns Unix timestamp in seconds
 */
export const timeStringToTimestamp = (timeString: string): number => {
    const [hours, minutes] = timeString.split(':').map(Number);
    const date = new Date();
    date.setHours(hours, minutes, 0, 0);
    return Math.floor(date.getTime() / 1000);
};

/**
 * Convert Unix timestamp to HH:MM time string
 * @param timestamp - Unix timestamp in seconds
 * @returns Time in HH:MM format
 */
export const timestampToTimeString = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return `${String(date.getHours()).padStart(2, '0')}:${String(
        date.getMinutes()
    ).padStart(2, '0')}`;
};

/**
 * Convert date string (YYYY-MM-DD) to Unix timestamp
 * @param dateString - Date in YYYY-MM-DD format
 * @returns Unix timestamp in seconds
 */
export const dateStringToTimestamp = (dateString: string): number => {
    return Math.floor(new Date(dateString).getTime() / 1000);
};

/**
 * Convert Unix timestamp to date string (YYYY-MM-DD)
 * @param timestamp - Unix timestamp in seconds
 * @returns Date in YYYY-MM-DD format
 */
export const timestampToDateString = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return date.toISOString().split('T')[0];
};
