/**
 * 에러 타입 정의
 * 철칙: 명시적 타입, 불변성, discriminated union 사용
 */

/**
 * 기본 에러 인터페이스
 */
export interface BaseError {
  readonly code: string;
  readonly message: string;
  readonly timestamp: string;
  readonly traceId?: string;
}

/**
 * 필드 에러
 */
export interface FieldError extends BaseError {
  readonly type: 'field';
  readonly field: string;
  readonly value?: unknown;
}

/**
 * 네트워크 에러
 */
export interface NetworkError extends BaseError {
  readonly type: 'network';
  readonly status?: number;
  readonly statusText?: string;
  readonly url?: string;
}

/**
 * 검증 에러
 */
export interface ValidationError extends BaseError {
  readonly type: 'validation';
  readonly errors: readonly FieldError[];
}

/**
 * 비즈니스 로직 에러
 */
export interface BusinessError extends BaseError {
  readonly type: 'business';
  readonly details?: Record<string, unknown>;
}

/**
 * 시스템 에러
 */
export interface SystemError extends BaseError {
  readonly type: 'system';
  readonly stack?: string;
}

/**
 * 통합 에러 타입
 */
export type AppError = 
  | FieldError
  | NetworkError
  | ValidationError
  | BusinessError
  | SystemError;

/**
 * 에러 코드 상수
 */
export const ErrorCodes = {
  // 네트워크 에러
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  
  // 검증 에러
  VALIDATION_FAILED: 'VALIDATION_FAILED',
  REQUIRED_FIELD: 'REQUIRED_FIELD',
  INVALID_FORMAT: 'INVALID_FORMAT',
  DUPLICATE_VALUE: 'DUPLICATE_VALUE',
  
  // 비즈니스 에러
  ENTITY_NOT_FOUND: 'ENTITY_NOT_FOUND',
  OPERATION_FAILED: 'OPERATION_FAILED',
  CONSTRAINT_VIOLATION: 'CONSTRAINT_VIOLATION',
  CONCURRENT_MODIFICATION: 'CONCURRENT_MODIFICATION',
  
  // 시스템 에러
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
} as const;

export type ErrorCode = typeof ErrorCodes[keyof typeof ErrorCodes];

/**
 * 에러 타입 가드
 * 철칙: 타입 안전성, 런타임 체크
 */
export const isFieldError = (error: AppError): error is FieldError => {
  return error.type === 'field';
};

export const isNetworkError = (error: AppError): error is NetworkError => {
  return error.type === 'network';
};

export const isValidationError = (error: AppError): error is ValidationError => {
  return error.type === 'validation';
};

export const isBusinessError = (error: AppError): error is BusinessError => {
  return error.type === 'business';
};

export const isSystemError = (error: AppError): error is SystemError => {
  return error.type === 'system';
};