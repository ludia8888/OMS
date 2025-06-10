/**
 * 디자인 토큰 - 타이포그래피 시스템
 * 철칙: 타입 안전성, 웹 폰트 최적화, 접근성 고려
 */

/**
 * 폰트 패밀리
 */
export const FontFamily = {
  primary: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, sans-serif',
  mono: '"SF Mono", Monaco, Inconsolata, "Roboto Mono", "Source Code Pro", Consolas, monospace',
  serif: 'Georgia, "Times New Roman", Times, serif',
} as const;

/**
 * 폰트 무게
 */
export const FontWeight = {
  light: 300,
  regular: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
} as const;

/**
 * 줄 높이
 */
export const LineHeight = {
  tight: 1.2,
  normal: 1.4,
  relaxed: 1.6,
  loose: 1.8,
} as const;

/**
 * 폰트 크기 스케일
 */
export const FontSize = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 20,
  xxl: 24,
  xxxl: 32,
  xxxxl: 48,
} as const;

/**
 * 타이포그래피 스타일 정의
 */
export const TypographyStyles = {
  // 헤딩
  h1: {
    fontSize: FontSize.xxxxl,
    fontWeight: FontWeight.bold,
    lineHeight: LineHeight.tight,
    fontFamily: FontFamily.primary,
  },
  h2: {
    fontSize: FontSize.xxxl,
    fontWeight: FontWeight.bold,
    lineHeight: LineHeight.tight,
    fontFamily: FontFamily.primary,
  },
  h3: {
    fontSize: FontSize.xxl,
    fontWeight: FontWeight.semibold,
    lineHeight: LineHeight.tight,
    fontFamily: FontFamily.primary,
  },
  h4: {
    fontSize: FontSize.xl,
    fontWeight: FontWeight.semibold,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },
  h5: {
    fontSize: FontSize.lg,
    fontWeight: FontWeight.medium,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },
  h6: {
    fontSize: FontSize.md,
    fontWeight: FontWeight.medium,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },

  // 본문
  body1: {
    fontSize: FontSize.md,
    fontWeight: FontWeight.regular,
    lineHeight: LineHeight.relaxed,
    fontFamily: FontFamily.primary,
  },
  body2: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.regular,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },

  // 캡션
  caption: {
    fontSize: FontSize.xs,
    fontWeight: FontWeight.regular,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },

  // 라벨
  label: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.medium,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.primary,
  },

  // 코드
  code: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.regular,
    lineHeight: LineHeight.normal,
    fontFamily: FontFamily.mono,
  },

  // 버튼
  button: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.medium,
    lineHeight: LineHeight.tight,
    fontFamily: FontFamily.primary,
  },
} as const;

/**
 * 타이포그래피 타입 정의
 */
export type TypographyVariant = keyof typeof TypographyStyles;
export type FontFamilyValue = typeof FontFamily[keyof typeof FontFamily];
export type FontWeightValue = typeof FontWeight[keyof typeof FontWeight];
export type FontSizeValue = typeof FontSize[keyof typeof FontSize];
export type LineHeightValue = typeof LineHeight[keyof typeof LineHeight];

/**
 * CSS 타이포그래피 스타일 생성
 * 철칙: 순수 함수, 타입 안전성
 */
export const getTypographyStyle = (variant: TypographyVariant): React.CSSProperties => {
  const style = TypographyStyles[variant];
  return {
    fontSize: `${style.fontSize}px`,
    fontWeight: style.fontWeight,
    lineHeight: style.lineHeight,
    fontFamily: style.fontFamily,
  };
};