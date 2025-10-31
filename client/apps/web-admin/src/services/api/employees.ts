import { http } from '@/services/http';
import type { Employee, UpdateEmployeePayload } from '@/types';

export async function fetchEmployeeById(id: string): Promise<Employee> {
  return http.get<Employee>(`/employees/${id}`);
}

export async function updateEmployeeById(
  id: string,
  payload: UpdateEmployeePayload
): Promise<Employee> {
  return http.patch<Employee>(`/employees/${id}`, payload);
}
