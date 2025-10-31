import { Response } from 'express';
import { ApiResponse } from '../types';

export function sendSuccess<T>(res: Response, data: T, message: string = 'Success', statusCode: number = 200): Response {
  const response: ApiResponse<T> = {
    success: true,
    message,
    data,
    statusCode,
  };
  return res.status(statusCode).json(response);
}

export function sendError(res: Response, error: string, statusCode: number = 400, message: string = 'Error'): Response {
  const response: ApiResponse = {
    success: false,
    message,
    error,
    statusCode,
  };
  return res.status(statusCode).json(response);
}
