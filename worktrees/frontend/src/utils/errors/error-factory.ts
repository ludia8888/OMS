import type { 
  AppError, 
  FieldError, 
  NetworkError, 
  ValidationError, 
  BusinessError, 
  SystemError,
  ErrorCode 
} from './error-types';
import { ErrorCodes } from './error-types';

/**
 * 타임스탬프 생성
 * 철칙: 순수 함수, 부작용 없음
 */
const createTimestamp = (): string => {
  return new Date().toISOString();
};

/**
 * 트레이스 ID 생성
 */
const createTraceId = (): string => {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

/**
 * 필드 에러 생성
 */
export const createFieldError = (
  field: string,
  message: string,
  code: ErrorCode = ErrorCodes.VALIDATION_FAILED,
  value?: unknown
): FieldError => ({
  type: 'field',
  field,
  message,
  code,
  value,
  timestamp: createTimestamp(),
  traceId: createTraceId(),
});

/**
 * 네트워크 에러 생성
 */
export const createNetworkError = (
  message: string,
  options?: {
    status?: number;
    statusText?: string;
    url?: string;
    code?: ErrorCode;
  }
): NetworkError => ({
  type: 'network',
  message,
  code: options?.code ?? ErrorCodes.NETWORK_ERROR,
  status: options?.status,
  statusText: options?.statusText,
  url: options?.url,
  timestamp: createTimestamp(),
  traceId: createTraceId(),
});

/**
 * 검증 에러 생성
 */
export const createValidationError = (
  message: string,
  errors: readonly FieldError[]
): ValidationError => ({
  type: 'validation',
  message,
  code: ErrorCodes.VALIDATION_FAILED,
  errors,
  timestamp: createTimestamp(),
  traceId: createTraceId(),
});

/**
 * 비즈니스 에러 생성
 */
export const createBusinessError = (
  message: string,
  code: ErrorCode,
  details?: Record<string, unknown>
): BusinessError => ({
  type: 'business',
  message,
  code,
  details,
  timestamp: createTimestamp(),
  traceId: createTraceId(),
});

/**
 * 시스템 에러 생성
 */
export const createSystemError = (
  message: string,
  stack?: string
): SystemError => ({
  type: 'system',
  message,
  code: ErrorCodes.INTERNAL_ERROR,
  stack,
  timestamp: createTimestamp(),
  traceId: createTraceId(),
});

/**
 * Error 객체에서 AppError 변환
 * 철칙: 방어적 프로그래밍, 타입 안전성
 */
export const fromError = (error: unknown): AppError => {
  // 이미 AppError인 경우
  if (isAppError(error)) {
    return error;
  }
  
  // 표준 Error 객체인 경우
  if (error instanceof Error) {
    return createSystemError(error.message, error.stack);
  }
  
  // 문자열인 경우
  if (typeof error === 'string') {
    return createSystemError(error);
  }
  
  // 그 외의 경우
  return createSystemError('An unknown error occurred');
};

/**
 * AppError 타입 가드
 */
const isAppError = (error: unknown): error is AppError => {
  return (
    typeof error === 'object' &&
    error !== null &&
    'type' in error &&
    'code' in error &&
    'message' in error &&
    'timestamp' in error
  );
};

/**
 * HTTP 상태 코드에서 에러 코드 결정
 * 철칙: 순수 함수, exhaustive check
 */
export const getErrorCodeFromStatus = (status: number): ErrorCode => {
  switch (status) {
    case 401:
      return ErrorCodes.UNAUTHORIZED;
    case 403:
      return ErrorCodes.FORBIDDEN;
    case 404:
      return ErrorCodes.NOT_FOUND;
    case 408:
    case 504:
      return ErrorCodes.TIMEOUT;
    case 503:
      return ErrorCodes.SERVICE_UNAVAILABLE;
    default:
      if (status >= 400 && status < 500) {
        return ErrorCodes.VALIDATION_FAILED;
      }
      if (status >= 500) {
        return ErrorCodes.INTERNAL_ERROR;
      }
      return ErrorCodes.NETWORK_ERROR;
  }
};