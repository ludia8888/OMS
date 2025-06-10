import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

/**
 * 토스트 타입
 */
export type ToastIntent = 'none' | 'primary' | 'success' | 'warning' | 'danger';

/**
 * 토스트 메시지
 */
export interface Toast {
  readonly id: string;
  readonly message: string;
  readonly intent: ToastIntent;
  readonly timeout?: number;
  readonly action?: {
    readonly text: string;
    readonly onClick: () => void;
  };
}

/**
 * 모달 타입
 */
export type ModalType = 
  | 'createObjectType'
  | 'editObjectType'
  | 'deleteObjectType'
  | 'createLinkType'
  | 'editLinkType'
  | 'deleteLinkType'
  | 'propertyEditor'
  | 'confirm';

/**
 * 모달 상태
 */
export interface ModalState {
  readonly type: ModalType;
  readonly isOpen: boolean;
  readonly data?: unknown;
}

/**
 * 사이드바 상태
 */
export interface SidebarState {
  readonly isOpen: boolean;
  readonly activePanel: 'objectTypes' | 'linkTypes' | 'search' | null;
}

/**
 * UI 스토어 상태
 */
interface UIStoreState {
  // 토스트
  readonly toasts: readonly Toast[];
  
  // 모달
  readonly modal: ModalState;
  
  // 사이드바
  readonly sidebar: SidebarState;
  
  // 테마
  readonly theme: 'light' | 'dark';
  
  // 레이아웃
  readonly isCompactMode: boolean;
  readonly showHelperText: boolean;
  
  // 글로벌 로딩
  readonly isGlobalLoading: boolean;
  readonly globalLoadingMessage?: string;
}

/**
 * UI 스토어 액션
 */
interface UIStoreActions {
  // 토스트 액션
  showToast: (toast: Omit<Toast, 'id'>) => void;
  removeToast: (id: string) => void;
  clearToasts: () => void;
  
  // 모달 액션
  openModal: (type: ModalType, data?: unknown) => void;
  closeModal: () => void;
  
  // 사이드바 액션
  toggleSidebar: () => void;
  setSidebarOpen: (isOpen: boolean) => void;
  setActivePanel: (panel: 'objectTypes' | 'linkTypes' | 'search' | null) => void;
  
  // 테마 액션
  setTheme: (theme: 'light' | 'dark') => void;
  toggleTheme: () => void;
  
  // 레이아웃 액션
  setCompactMode: (isCompact: boolean) => void;
  toggleCompactMode: () => void;
  setShowHelperText: (show: boolean) => void;
  toggleHelperText: () => void;
  
  // 글로벌 로딩 액션
  setGlobalLoading: (isLoading: boolean, message?: string) => void;
  
  // 유틸리티
  showSuccessToast: (message: string) => void;
  showErrorToast: (message: string) => void;
  showWarningToast: (message: string) => void;
}

/**
 * 토스트 ID 생성
 * 철칙: 순수 함수, 부작용 없음
 */
const generateToastId = (): string => {
  return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

/**
 * 기본 토스트 타임아웃 (밀리초)
 */
const DEFAULT_TOAST_TIMEOUT = 5000;

/**
 * UI 스토어
 * 철칙: 불변성, 타입 안전성, 단일 책임
 */
export const useUIStore = create<UIStoreState & UIStoreActions>()(
  immer((set, get) => ({
    // 초기 상태
    toasts: [],
    
    modal: {
      type: 'confirm',
      isOpen: false,
      data: undefined,
    },
    
    sidebar: {
      isOpen: true,
      activePanel: 'objectTypes',
    },
    
    theme: 'light',
    isCompactMode: false,
    showHelperText: true,
    isGlobalLoading: false,
    globalLoadingMessage: undefined,
    
    // 토스트 액션 구현
    showToast: (toast) => {
      const id = generateToastId();
      const newToast: Toast = {
        ...toast,
        id,
        timeout: toast.timeout ?? DEFAULT_TOAST_TIMEOUT,
      };
      
      set((state) => {
        state.toasts = [...state.toasts, newToast];
      });
      
      // 자동 제거 설정
      if (newToast.timeout > 0) {
        setTimeout(() => {
          get().removeToast(id);
        }, newToast.timeout);
      }
    },
    
    removeToast: (id) => set((state) => {
      state.toasts = state.toasts.filter(t => t.id !== id);
    }),
    
    clearToasts: () => set((state) => {
      state.toasts = [];
    }),
    
    // 모달 액션 구현
    openModal: (type, data) => set((state) => {
      state.modal.type = type;
      state.modal.isOpen = true;
      state.modal.data = data;
    }),
    
    closeModal: () => set((state) => {
      state.modal.isOpen = false;
      state.modal.data = undefined;
    }),
    
    // 사이드바 액션 구현
    toggleSidebar: () => set((state) => {
      state.sidebar.isOpen = !state.sidebar.isOpen;
    }),
    
    setSidebarOpen: (isOpen) => set((state) => {
      state.sidebar.isOpen = isOpen;
    }),
    
    setActivePanel: (panel) => set((state) => {
      state.sidebar.activePanel = panel;
    }),
    
    // 테마 액션 구현
    setTheme: (theme) => set((state) => {
      state.theme = theme;
      // DOM 클래스 업데이트
      document.documentElement.classList.remove('bp5-light', 'bp5-dark');
      document.documentElement.classList.add(theme === 'dark' ? 'bp5-dark' : 'bp5-light');
    }),
    
    toggleTheme: () => {
      const currentTheme = get().theme;
      get().setTheme(currentTheme === 'light' ? 'dark' : 'light');
    },
    
    // 레이아웃 액션 구현
    setCompactMode: (isCompact) => set((state) => {
      state.isCompactMode = isCompact;
    }),
    
    toggleCompactMode: () => set((state) => {
      state.isCompactMode = !state.isCompactMode;
    }),
    
    setShowHelperText: (show) => set((state) => {
      state.showHelperText = show;
    }),
    
    toggleHelperText: () => set((state) => {
      state.showHelperText = !state.showHelperText;
    }),
    
    // 글로벌 로딩 액션 구현
    setGlobalLoading: (isLoading, message) => set((state) => {
      state.isGlobalLoading = isLoading;
      state.globalLoadingMessage = message;
    }),
    
    // 유틸리티 구현
    showSuccessToast: (message) => {
      get().showToast({
        message,
        intent: 'success',
      });
    },
    
    showErrorToast: (message) => {
      get().showToast({
        message,
        intent: 'danger',
        timeout: 10000, // 에러는 더 오래 표시
      });
    },
    
    showWarningToast: (message) => {
      get().showToast({
        message,
        intent: 'warning',
      });
    },
  }))
);