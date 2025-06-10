import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Card } from '../card';

/**
 * Card 컴포넌트 테스트
 * 철칙: AAA 패턴, 하나의 assertion per test, 90% 커버리지
 */

describe('Card', () => {
  it('renders children correctly', () => {
    // Arrange
    const testContent = 'Test Card Content';
    
    // Act
    render(<Card testId="test-card">{testContent}</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveTextContent(testContent);
  });

  it('applies default variant class', () => {
    // Arrange & Act
    render(<Card testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveClass('oms-card');
    expect(screen.getByTestId('test-card')).toHaveClass('oms-card--default');
  });

  it('applies outlined variant class when specified', () => {
    // Arrange & Act
    render(<Card variant="outlined" testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveClass('oms-card--outlined');
  });

  it('applies elevated variant class when specified', () => {
    // Arrange & Act
    render(<Card variant="elevated" testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveClass('oms-card--elevated');
  });

  it('applies custom className along with card classes', () => {
    // Arrange
    const customClass = 'custom-card-class';
    
    // Act
    render(<Card className={customClass} testId="test-card">Content</Card>);
    
    // Assert
    const cardElement = screen.getByTestId('test-card');
    expect(cardElement).toHaveClass('oms-card');
    expect(cardElement).toHaveClass(customClass);
  });

  it('applies padding style based on padding prop', () => {
    // Arrange & Act
    render(<Card padding="lg" testId="test-card">Content</Card>);
    
    // Assert
    const cardElement = screen.getByTestId('test-card');
    expect(cardElement).toHaveStyle({ padding: '36px' }); // 24 * 1.5
  });

  it('applies no padding when padding is none', () => {
    // Arrange & Act
    render(<Card padding="none" testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveStyle({ padding: '0px' });
  });

  it('merges custom style with internal styles', () => {
    // Arrange
    const customStyle = { backgroundColor: 'red' };
    
    // Act
    render(<Card style={customStyle} testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveStyle(customStyle);
  });

  it('passes through Blueprint Card props', () => {
    // Arrange & Act
    render(<Card interactive testId="test-card">Content</Card>);
    
    // Assert
    expect(screen.getByTestId('test-card')).toHaveClass('bp5-interactive');
  });
});