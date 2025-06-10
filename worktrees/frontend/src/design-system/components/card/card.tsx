import React from 'react';
import { Card as BpCard } from '@blueprintjs/core';
import type { CardProps as BpCardProps } from '@blueprintjs/core';
import { ComponentSpacing } from '../../tokens';
import './card.css';

/**
 * 확장된 카드 속성
 */
export interface CardProps extends Omit<BpCardProps, 'children'> {
  readonly children: React.ReactNode;
  readonly variant?: 'default' | 'outlined' | 'elevated';
  readonly padding?: 'none' | 'sm' | 'md' | 'lg';
  readonly spacing?: 'sm' | 'md' | 'lg';
  readonly testId?: string;
}

/**
 * 패딩 매핑
 * 철칙: 순수 함수, 명시적 매핑
 */
const getPaddingValue = (padding: CardProps['padding']): number => {
  switch (padding) {
    case 'none':
      return 0;
    case 'sm':
      return ComponentSpacing.card.padding / 2;
    case 'md':
      return ComponentSpacing.card.padding;
    case 'lg':
      return ComponentSpacing.card.padding * 1.5;
    default:
      return ComponentSpacing.card.padding;
  }
};

/**
 * CSS 클래스명 생성
 */
const getCardClassName = (variant: CardProps['variant'], className?: string): string => {
  const baseClass = 'oms-card';
  const variantClass = variant ? `${baseClass}--${variant}` : '';
  
  return [baseClass, variantClass, className]
    .filter(Boolean)
    .join(' ');
};

/**
 * OMS 카드 컴포넌트
 * 철칙: 타입 안전성, 단일 책임, 접근성
 */
export const Card: React.FC<CardProps> = ({
  children,
  variant = 'default',
  padding = 'md',
  className,
  style,
  testId,
  ...blueprintProps
}) => {
  const paddingValue = getPaddingValue(padding);
  const cardClassName = getCardClassName(variant, className);
  
  const cardStyle: React.CSSProperties = {
    ...style,
    padding: paddingValue,
  };

  return (
    <BpCard
      {...blueprintProps}
      className={cardClassName}
      style={cardStyle}
      data-testid={testId}
    >
      {children}
    </BpCard>
  );
};