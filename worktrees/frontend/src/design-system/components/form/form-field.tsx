import React from 'react';
import { FormGroup, type FormGroupProps } from '@blueprintjs/core';
import { ComponentSpacing } from '../../tokens';
import type { ValidationResult } from '@/utils/validators/field-validators';
import './form-field.css';

/**
 * 폼 필드 속성
 */
export interface FormFieldProps extends Omit<FormGroupProps, 'children'> {
  readonly children: React.ReactNode;
  readonly validation?: ValidationResult;
  readonly showValidation?: boolean;
  readonly testId?: string;
}

/**
 * 인텐트 결정
 * 철칙: 순수 함수, 명시적 타입 반환
 */
const getFieldIntent = (validation?: ValidationResult): FormGroupProps['intent'] => {
  if (!validation || validation.isValid) {
    return undefined;
  }
  return 'danger';
};

/**
 * 헬퍼 텍스트 생성
 */
const getHelperText = (
  originalHelperText?: React.ReactNode,
  validation?: ValidationResult,
  showValidation = true
): React.ReactNode => {
  if (validation && !validation.isValid && showValidation) {
    return validation.errors.join(', ');
  }
  return originalHelperText;
};

/**
 * OMS 폼 필드 컴포넌트
 * 철칙: 타입 안전성, 접근성, 단일 책임
 */
export const FormField: React.FC<FormFieldProps> = ({
  children,
  validation,
  showValidation = true,
  helperText,
  className,
  testId,
  ...formGroupProps
}) => {
  const intent = getFieldIntent(validation);
  const finalHelperText = getHelperText(helperText, validation, showValidation);
  const fieldClassName = ['oms-form-field', className].filter(Boolean).join(' ');

  return (
    <FormGroup
      {...formGroupProps}
      intent={intent}
      helperText={finalHelperText}
      className={fieldClassName}
      data-testid={testId}
    >
      {children}
    </FormGroup>
  );
};