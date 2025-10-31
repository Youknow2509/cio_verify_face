import { clearAuthToken, http, HttpError } from '@/services/http';
import type { AccountProfile, ChangePasswordPayload } from '@/types';

function removeStoredUser() {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    window.localStorage?.removeItem('user');
  } catch {
    // Ignore storage errors
  }

  try {
    window.sessionStorage?.removeItem('user');
  } catch {
    // Ignore storage errors
  }
}

export async function fetchMyAccount(): Promise<AccountProfile> {
  return http.get<AccountProfile>('/me');
}

export async function changePassword(payload: ChangePasswordPayload): Promise<void> {
  await http.post<void>('/auth/change-password', payload);
}

export async function logout(): Promise<void> {
  try {
    await http.post<void>('/auth/logout');
  } catch (error) {
    const isOptionalEndpointMissing = error instanceof HttpError && error.status === 404;
    if (!isOptionalEndpointMissing) {
      throw error;
    }
  } finally {
    clearAuthToken();
    removeStoredUser();
  }
}
