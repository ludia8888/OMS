import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { InputGroup } from '@blueprintjs/core';
import { FormField } from '../form-field';
import type { ValidationResult } from '@/utils/validators/field-validators';

/**
 * FormField 컴포넌트 테스트
 * 철칙: AAA 패턴, 타입 안전성, 에러 케이스 포함
 */

describe('FormField', () => {
  it('renders children correctly', () => {
    // Arrange
    const testLabel = 'Test Field';
    
    // Act
    render(
      <FormField label={testLabel} testId="test-form-field">
        <InputGroup placeholder="Test input" />
      </FormField>
    );
    
    // Assert
    expect(screen.getByTestId('test-form-field')).toBeInTheDocument();
    expect(screen.getByText(testLabel)).toBeInTheDocument();
  });

  it('shows no intent when validation is valid', () => {
    // Arrange
    const validValidation: ValidationResult = {
      isValid: true,
      errors: [],
    };
    
    // Act
    render(
      <FormField label="Valid Field" validation={validValidation} testId="test-form-field">
        <InputGroup />
      </FormField>
    );
    
    // Assert
    const formGroup = screen.getByTestId('test-form-field');
    expect(formGroup).not.toHaveClass('bp5-intent-danger');
  });

  it('shows danger intent when validation has errors', () => {
    // Arrange
    const invalidValidation: ValidationResult = {
      isValid: false,
      errors: ['This field is required'],
    };
    
    // Act
    render(
      <FormField label="Invalid Field" validation={invalidValidation} testId="test-form-field">
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.getByTestId('test-form-field')).toHaveClass('bp5-intent-danger');
  });

  it('displays validation errors as helper text', () => {
    // Arrange
    const errorMessage = 'This field is required';
    const invalidValidation: ValidationResult = {
      isValid: false,
      errors: [errorMessage],
    };
    
    // Act
    render(
      <FormField label="Error Field" validation={invalidValidation} testId="test-form-field">
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('displays multiple validation errors joined by comma', () => {
    // Arrange
    const errors = ['Error 1', 'Error 2'];
    const invalidValidation: ValidationResult = {
      isValid: false,
      errors,
    };
    
    // Act
    render(
      <FormField label="Multi Error Field" validation={invalidValidation}>
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.getByText('Error 1, Error 2')).toBeInTheDocument();
  });

  it('shows original helper text when validation is valid', () => {
    // Arrange
    const helperText = 'This is a helpful hint';
    const validValidation: ValidationResult = {
      isValid: true,
      errors: [],
    };
    
    // Act
    render(
      <FormField 
        label="Helper Field" 
        helperText={helperText}
        validation={validValidation}
      >
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.getByText(helperText)).toBeInTheDocument();
  });

  it('hides validation errors when showValidation is false', () => {
    // Arrange
    const errorMessage = 'This field is required';
    const invalidValidation: ValidationResult = {
      isValid: false,
      errors: [errorMessage],
    };
    
    // Act
    render(
      <FormField 
        label="Hidden Error Field" 
        validation={invalidValidation}
        showValidation={false}
      >
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.queryByText(errorMessage)).not.toBeInTheDocument();
  });

  it('applies custom className', () => {
    // Arrange
    const customClass = 'custom-form-field';
    
    // Act
    render(
      <FormField className={customClass} label="Custom Field" testId="test-form-field">
        <InputGroup />
      </FormField>
    );
    
    // Assert
    expect(screen.getByTestId('test-form-field')).toHaveClass('oms-form-field');
    expect(screen.getByTestId('test-form-field')).toHaveClass(customClass);
  });
});