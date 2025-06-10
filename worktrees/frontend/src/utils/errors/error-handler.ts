import type { AppError } from './error-types';
import { 
  isFieldError, 
  isNetworkError, 
  isValidationError, 
  isBusinessError, 
  isSystemError 
} from './error-types';
import { fromError } from './error-factory';
import { logger } from '../logger';
import { useUIStore } from '@/stores/ui/ui.store';

/**
 * 에러 메시지 포맷터
 * 철칙: 순수 함수, 사용자 친화적 메시지
 */
const formatErrorMessage = (error: AppError): string => {
  if (isValidationError(error)) {
    const fieldErrors = error.errors
      .map(e => `${e.field}: ${e.message}`)
      .join(', ');
    return `Validation failed: ${fieldErrors}`;
  }
  
  if (isNetworkError(error)) {
    if (error.status === 401) {
      return 'You are not authorized. Please log in.';
    }
    if (error.status === 403) {
      return 'You do not have permission to perform this action.';
    }
    if (error.status === 404) {
      return 'The requested resource was not found.';
    }
    if (error.status && error.status >= 500) {
      return 'A server error occurred. Please try again later.';
    }
  }
  
  return error.message || 'An unexpected error occurred.';
};

/**
 * 에러 심각도 결정
 * 철칙: 명시적 반환 타입, exhaustive check
 */
const getErrorSeverity = (error: AppError): 'info' | 'warning' | 'error' => {
  if (isFieldError(error) || isValidationError(error)) {
    return 'warning';
  }
  
  if (isNetworkError(error)) {
    if (error.status && error.status >= 500) {
      return 'error';
    }
    return 'warning';
  }
  
  if (isBusinessError(error)) {
    return 'warning';
  }
  
  if (isSystemError(error)) {
    return 'error';
  }
  
  return 'error';
};

/**
 * 에러 처리 옵션
 */
export interface ErrorHandlerOptions {
  readonly showToast?: boolean;
  readonly logError?: boolean;
  readonly throwError?: boolean;
  readonly customMessage?: string;
}

/**
 * 중앙 에러 처리 함수
 * 철칙: 단일 책임, 설정 가능한 동작
 */
export const handleError = (
  error: unknown,
  options: ErrorHandlerOptions = {}
): AppError => {
  const {
    showToast = true,
    logError = true,
    throwError = false,
    customMessage,
  } = options;
  
  // AppError로 변환
  const appError = fromError(error);
  
  // 로깅
  if (logError) {
    const severity = getErrorSeverity(appError);
    logger[severity]('Error occurred', {
      error: appError,
      traceId: appError.traceId,
    });
  }
  
  // 토스트 표시
  if (showToast) {
    const message = customMessage || formatErrorMessage(appError);
    const { showErrorToast, showWarningToast } = useUIStore.getState();
    
    if (getErrorSeverity(appError) === 'error') {
      showErrorToast(message);
    } else {
      showWarningToast(message);
    }
  }
  
  // 에러 재발생
  if (throwError) {
    throw appError;
  }
  
  return appError;
};

/**
 * Promise 에러 처리 래퍼
 * 철칙: 타입 안전성, 제네릭 사용
 */
export const withErrorHandling = async <T>(
  promise: Promise<T>,
  options?: ErrorHandlerOptions
): Promise<T | null> => {
  try {
    return await promise;
  } catch (error) {
    handleError(error, options);
    return null;
  }
};

/**
 * 에러 복구 시도
 * 철칙: 명시적 재시도 전략
 */
export const retryWithBackoff = async <T>(
  fn: () => Promise<T>,
  options: {
    readonly maxRetries?: number;
    readonly initialDelay?: number;
    readonly maxDelay?: number;
    readonly backoffFactor?: number;
    readonly shouldRetry?: (error: AppError, attempt: number) => boolean;
  } = {}
): Promise<T> => {
  const {
    maxRetries = 3,
    initialDelay = 1000,
    maxDelay = 10000,
    backoffFactor = 2,
    shouldRetry = (error) => isNetworkError(error) && error.status !== 401 && error.status !== 403,
  } = options;
  
  let lastError: AppError;
  
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = fromError(error);
      
      if (!shouldRetry(lastError, attempt + 1) || attempt === maxRetries - 1) {
        throw lastError;
      }
      
      const delay = Math.min(
        initialDelay * Math.pow(backoffFactor, attempt),
        maxDelay
      );
      
      logger.info(`Retrying after ${delay}ms (attempt ${attempt + 1}/${maxRetries})`);
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
  
  throw lastError!;
};