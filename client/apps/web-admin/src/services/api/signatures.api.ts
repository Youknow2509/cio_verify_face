// src/services/api/signatures.api.ts
/**
 * Signature Upload Service APIs
 * Quản lý chữ ký nhân viên
 */

import { http } from '../http';
import type { ApiResponse } from '@/types';

const API_PREFIX = '/api/v1/signatures';

export interface SignatureData {
  id: string;
  userId: string;
  signatureUrl: string;
  fileName: string;
  createdAt: string;
}

/**
 * Upload user signature
 * @param file Signature image file
 * @returns Uploaded signature data
 */
export async function uploadSignatureAPI(file: File): Promise<ApiResponse<SignatureData>> {
  const formData = new FormData();
  formData.append('file', file);
  
  return http.post(API_PREFIX, formData);
}

/**
 * Get user's current signature
 * @param userId User ID
 * @returns Signature data
 */
export async function getUserSignatureAPI(userId: string): Promise<ApiResponse<SignatureData>> {
  return http.get(`${API_PREFIX}/${userId}`);
}

/**
 * Delete user signature
 * @param signatureId Signature ID
 * @returns Void
 */
export async function deleteSignatureAPI(signatureId: string): Promise<ApiResponse<void>> {
  return http.delete(`${API_PREFIX}/${signatureId}`);
}
