/**
 * 디자인 토큰 - 간격 시스템
 * 철칙: 8px 그리드 시스템, 타입 안전성, const assertion
 */

/**
 * 기본 간격 단위 (8px 그리드)
 */
const BASE_UNIT = 8;

/**
 * 간격 스케일
 */
export const SpacingScale = {
  NONE: 0,
  XS: BASE_UNIT * 0.5, // 4px
  SM: BASE_UNIT * 1,   // 8px
  MD: BASE_UNIT * 2,   // 16px
  LG: BASE_UNIT * 3,   // 24px
  XL: BASE_UNIT * 4,   // 32px
  XXL: BASE_UNIT * 6,  // 48px
  XXXL: BASE_UNIT * 8, // 64px
} as const;

/**
 * 컴포넌트별 간격
 */
export const ComponentSpacing = {
  // 버튼
  button: {
    paddingX: SpacingScale.MD,
    paddingY: SpacingScale.SM,
    gap: SpacingScale.SM,
  },

  // 카드
  card: {
    padding: SpacingScale.LG,
    gap: SpacingScale.MD,
  },

  // 폼
  form: {
    fieldGap: SpacingScale.MD,
    labelGap: SpacingScale.XS,
    sectionGap: SpacingScale.XXL,
  },

  // 레이아웃
  layout: {
    sidebarWidth: BASE_UNIT * 32, // 256px
    headerHeight: BASE_UNIT * 8,  // 64px
    contentPadding: SpacingScale.LG,
    pageMargin: SpacingScale.XL,
  },

  // 목록
  list: {
    itemPadding: SpacingScale.MD,
    itemGap: SpacingScale.SM,
  },

  // 모달
  modal: {
    padding: SpacingScale.XXL,
    headerGap: SpacingScale.LG,
    footerGap: SpacingScale.LG,
  },
} as const;

/**
 * 반응형 간격
 */
export const ResponsiveSpacing = {
  mobile: {
    contentPadding: SpacingScale.MD,
    pageMargin: SpacingScale.MD,
  },
  tablet: {
    contentPadding: SpacingScale.LG,
    pageMargin: SpacingScale.LG,
  },
  desktop: {
    contentPadding: SpacingScale.XL,
    pageMargin: SpacingScale.XL,
  },
} as const;

/**
 * 간격 타입 정의
 */
export type SpacingValue = typeof SpacingScale[keyof typeof SpacingScale];
export type ComponentSpacingType = keyof typeof ComponentSpacing;

/**
 * 간격 유틸리티 함수
 * 철칙: 순수 함수, 타입 안전성
 */
export const getSpacing = (multiplier: number): number => {
  return BASE_UNIT * multiplier;
};

/**
 * CSS 간격 값 생성
 */
export const toCssValue = (spacing: SpacingValue): string => {
  return `${spacing}px`;
};