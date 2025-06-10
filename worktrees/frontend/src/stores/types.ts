/**
 * Store 공통 타입 정의
 * 철칙: any 금지, 모든 타입 명시적 정의
 */

// 로딩 상태 타입
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

// 에러 타입
export interface StoreError {
  readonly code: string;
  readonly message: string;
  readonly field?: string;
  readonly timestamp: string;
}

// 페이지네이션 상태
export interface PaginationState {
  readonly page: number;
  readonly pageSize: number;
  readonly totalPages: number;
  readonly totalCount: number;
  readonly hasNext: boolean;
  readonly hasPrevious: boolean;
}

// 필터 상태
export interface FilterState<T> {
  readonly active: boolean;
  readonly filters: T;
}

// 정렬 상태
export interface SortState {
  readonly field: string;
  readonly direction: 'ASC' | 'DESC';
}

// 선택 상태
export interface SelectionState<T> {
  readonly selected: ReadonlySet<T>;
  readonly lastSelected: T | null;
}

// 기본 엔티티 스토어 인터페이스
export interface BaseEntityStore<T, ID> {
  // 상태
  readonly items: ReadonlyMap<ID, T>;
  readonly loadingState: LoadingState;
  readonly error: StoreError | null;
  
  // 액션
  readonly setItems: (items: readonly T[]) => void;
  readonly addItem: (item: T) => void;
  readonly updateItem: (id: ID, updates: Partial<T>) => void;
  readonly removeItem: (id: ID) => void;
  readonly clearItems: () => void;
  readonly setLoadingState: (state: LoadingState) => void;
  readonly setError: (error: StoreError | null) => void;
}

// 리스트 스토어 인터페이스
export interface ListStore<T, ID, F> extends BaseEntityStore<T, ID> {
  // 리스트 상태
  readonly pagination: PaginationState;
  readonly filter: FilterState<F>;
  readonly sort: SortState;
  readonly selection: SelectionState<ID>;
  
  // 리스트 액션
  readonly setPagination: (pagination: PaginationState) => void;
  readonly setPage: (page: number) => void;
  readonly setPageSize: (pageSize: number) => void;
  readonly setFilter: (filter: Partial<F>) => void;
  readonly clearFilter: () => void;
  readonly setSort: (field: string, direction: 'ASC' | 'DESC') => void;
  readonly select: (id: ID) => void;
  readonly deselect: (id: ID) => void;
  readonly selectAll: () => void;
  readonly clearSelection: () => void;
}

// 폼 상태
export interface FormState<T> {
  readonly data: T;
  readonly isDirty: boolean;
  readonly isValid: boolean;
  readonly errors: ReadonlyMap<keyof T, string>;
  readonly touched: ReadonlySet<keyof T>;
}

// 폼 스토어 인터페이스
export interface FormStore<T> {
  readonly form: FormState<T>;
  readonly setField: <K extends keyof T>(field: K, value: T[K]) => void;
  readonly setFields: (fields: Partial<T>) => void;
  readonly setError: (field: keyof T, error: string | null) => void;
  readonly touchField: (field: keyof T) => void;
  readonly resetForm: (data?: T) => void;
  readonly validateForm: () => boolean;
}