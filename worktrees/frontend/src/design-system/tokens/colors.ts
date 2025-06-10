/**
 * 디자인 토큰 - 색상 시스템
 * 철칙: 불변성, 타입 안전성, const assertion 사용
 */

/**
 * 기본 팔레트
 */
export const BasePalette = {
  // Primary Colors (Blueprint 기반)
  BLUE_50: '#E3F2FD',
  BLUE_100: '#BBDEFB',
  BLUE_200: '#90CAF9',
  BLUE_300: '#64B5F6',
  BLUE_400: '#42A5F5',
  BLUE_500: '#2196F3', // Primary
  BLUE_600: '#1E88E5',
  BLUE_700: '#1976D2',
  BLUE_800: '#1565C0',
  BLUE_900: '#0D47A1',

  // Gray Scale
  GRAY_50: '#FAFAFA',
  GRAY_100: '#F5F5F5',
  GRAY_200: '#EEEEEE',
  GRAY_300: '#E0E0E0',
  GRAY_400: '#BDBDBD',
  GRAY_500: '#9E9E9E',
  GRAY_600: '#757575',
  GRAY_700: '#616161',
  GRAY_800: '#424242',
  GRAY_900: '#212121',

  // Semantic Colors
  SUCCESS_500: '#4CAF50',
  WARNING_500: '#FF9800',
  DANGER_500: '#F44336',
  INFO_500: '#2196F3',

  // White & Black
  WHITE: '#FFFFFF',
  BLACK: '#000000',
} as const;

/**
 * 테마 색상 정의
 */
export const LightTheme = {
  // Background
  background: {
    primary: BasePalette.WHITE,
    secondary: BasePalette.GRAY_50,
    tertiary: BasePalette.GRAY_100,
    overlay: 'rgba(0, 0, 0, 0.3)',
  },

  // Text
  text: {
    primary: BasePalette.GRAY_900,
    secondary: BasePalette.GRAY_700,
    tertiary: BasePalette.GRAY_500,
    inverse: BasePalette.WHITE,
    disabled: BasePalette.GRAY_400,
  },

  // Border
  border: {
    primary: BasePalette.GRAY_300,
    secondary: BasePalette.GRAY_200,
    focus: BasePalette.BLUE_500,
    error: BasePalette.DANGER_500,
  },

  // Interactive
  interactive: {
    primary: BasePalette.BLUE_500,
    primaryHover: BasePalette.BLUE_600,
    primaryActive: BasePalette.BLUE_700,
    secondary: BasePalette.GRAY_500,
    secondaryHover: BasePalette.GRAY_600,
    secondaryActive: BasePalette.GRAY_700,
  },

  // Status
  status: {
    success: BasePalette.SUCCESS_500,
    warning: BasePalette.WARNING_500,
    danger: BasePalette.DANGER_500,
    info: BasePalette.INFO_500,
  },
} as const;

export const DarkTheme = {
  // Background
  background: {
    primary: BasePalette.GRAY_900,
    secondary: BasePalette.GRAY_800,
    tertiary: BasePalette.GRAY_700,
    overlay: 'rgba(0, 0, 0, 0.5)',
  },

  // Text
  text: {
    primary: BasePalette.WHITE,
    secondary: BasePalette.GRAY_300,
    tertiary: BasePalette.GRAY_500,
    inverse: BasePalette.BLACK,
    disabled: BasePalette.GRAY_600,
  },

  // Border
  border: {
    primary: BasePalette.GRAY_600,
    secondary: BasePalette.GRAY_700,
    focus: BasePalette.BLUE_400,
    error: BasePalette.DANGER_500,
  },

  // Interactive
  interactive: {
    primary: BasePalette.BLUE_400,
    primaryHover: BasePalette.BLUE_300,
    primaryActive: BasePalette.BLUE_200,
    secondary: BasePalette.GRAY_400,
    secondaryHover: BasePalette.GRAY_300,
    secondaryActive: BasePalette.GRAY_200,
  },

  // Status
  status: {
    success: BasePalette.SUCCESS_500,
    warning: BasePalette.WARNING_500,
    danger: BasePalette.DANGER_500,
    info: BasePalette.INFO_500,
  },
} as const;

/**
 * 색상 타입 정의
 */
export type ColorTheme = typeof LightTheme;
export type ColorToken = keyof ColorTheme;
export type BackgroundColor = keyof ColorTheme['background'];
export type TextColor = keyof ColorTheme['text'];
export type BorderColor = keyof ColorTheme['border'];
export type InteractiveColor = keyof ColorTheme['interactive'];
export type StatusColor = keyof ColorTheme['status'];