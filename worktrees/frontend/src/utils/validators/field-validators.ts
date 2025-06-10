/**
 * 필드 유효성 검사 유틸리티
 * 철칙: 순수 함수, 타입 안전성, 재사용성
 */

/**
 * 유효성 검사 결과
 */
export interface ValidationResult {
  readonly isValid: boolean;
  readonly errors: readonly string[];
}

/**
 * 유효성 검사 함수 타입
 */
export type Validator<T> = (value: T) => ValidationResult;

/**
 * 성공적인 유효성 검사 결과
 */
const validResult: ValidationResult = {
  isValid: true,
  errors: [],
};

/**
 * 실패한 유효성 검사 결과 생성
 */
const invalidResult = (errors: readonly string[]): ValidationResult => ({
  isValid: false,
  errors,
});

/**
 * 필수 필드 검사
 */
export const required = <T>(message = 'This field is required'): Validator<T> => {
  return (value: T): ValidationResult => {
    if (value === null || value === undefined) {
      return invalidResult([message]);
    }
    
    if (typeof value === 'string' && value.trim() === '') {
      return invalidResult([message]);
    }
    
    if (Array.isArray(value) && value.length === 0) {
      return invalidResult([message]);
    }
    
    return validResult;
  };
};

/**
 * 문자열 길이 검사
 */
export const stringLength = (
  min?: number,
  max?: number,
  message?: string
): Validator<string> => {
  return (value: string): ValidationResult => {
    const length = value?.length || 0;
    const errors: string[] = [];
    
    if (min !== undefined && length < min) {
      errors.push(message || `Must be at least ${min} characters`);
    }
    
    if (max !== undefined && length > max) {
      errors.push(message || `Must be no more than ${max} characters`);
    }
    
    return errors.length > 0 ? invalidResult(errors) : validResult;
  };
};

/**
 * 정규식 패턴 검사
 */
export const pattern = (
  regex: RegExp,
  message = 'Invalid format'
): Validator<string> => {
  return (value: string): ValidationResult => {
    if (!value) return validResult; // 빈 값은 required에서 처리
    
    return regex.test(value) ? validResult : invalidResult([message]);
  };
};

/**
 * 이메일 형식 검사
 */
export const email = (message = 'Invalid email format'): Validator<string> => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return pattern(emailRegex, message);
};

/**
 * API 이름 형식 검사
 */
export const apiName = (message = 'Invalid API name format'): Validator<string> => {
  const apiNameRegex = /^[a-zA-Z][a-zA-Z0-9_]*$/;
  return pattern(apiNameRegex, message);
};

/**
 * 숫자 범위 검사
 */
export const numberRange = (
  min?: number,
  max?: number,
  message?: string
): Validator<number> => {
  return (value: number): ValidationResult => {
    if (value === null || value === undefined || isNaN(value)) {
      return validResult; // 빈 값은 required에서 처리
    }
    
    const errors: string[] = [];
    
    if (min !== undefined && value < min) {
      errors.push(message || `Must be at least ${min}`);
    }
    
    if (max !== undefined && value > max) {
      errors.push(message || `Must be no more than ${max}`);
    }
    
    return errors.length > 0 ? invalidResult(errors) : validResult;
  };
};

/**
 * 배열 길이 검사
 */
export const arrayLength = <T>(
  min?: number,
  max?: number,
  message?: string
): Validator<readonly T[]> => {
  return (value: readonly T[]): ValidationResult => {
    const length = value?.length || 0;
    const errors: string[] = [];
    
    if (min !== undefined && length < min) {
      errors.push(message || `Must have at least ${min} items`);
    }
    
    if (max !== undefined && length > max) {
      errors.push(message || `Must have no more than ${max} items`);
    }
    
    return errors.length > 0 ? invalidResult(errors) : validResult;
  };
};

/**
 * 유효성 검사 조합
 */
export const combine = <T>(...validators: readonly Validator<T>[]): Validator<T> => {
  return (value: T): ValidationResult => {
    const allErrors: string[] = [];
    
    for (const validator of validators) {
      const result = validator(value);
      if (!result.isValid) {
        allErrors.push(...result.errors);
      }
    }
    
    return allErrors.length > 0 ? invalidResult(allErrors) : validResult;
  };
};

/**
 * 조건부 유효성 검사
 */
export const when = <T>(
  condition: (value: T) => boolean,
  validator: Validator<T>
): Validator<T> => {
  return (value: T): ValidationResult => {
    return condition(value) ? validator(value) : validResult;
  };
};

/**
 * 객체 필드 유효성 검사
 */
export const validateObject = <T extends Record<string, unknown>>(
  obj: T,
  rules: Partial<Record<keyof T, Validator<T[keyof T]>>>
): Record<keyof T, ValidationResult> => {
  const results = {} as Record<keyof T, ValidationResult>;
  
  for (const [field, validator] of Object.entries(rules) as Array<[keyof T, Validator<T[keyof T]>]>) {
    results[field] = validator(obj[field]);
  }
  
  return results;
};

/**
 * 객체 전체 유효성 검사
 */
export const isObjectValid = <T extends Record<string, unknown>>(
  results: Record<keyof T, ValidationResult>
): boolean => {
  return Object.values(results).every((result) => result.isValid);
};