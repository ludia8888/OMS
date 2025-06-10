import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import type { ObjectType, ObjectId, ObjectTypeStatus } from '@/types/domain';
import type { ListStore, StoreError, PaginationState, FilterState, SortState, SelectionState, LoadingState } from '../types';

/**
 * ObjectType 필터 타입
 */
export interface ObjectTypeFilter {
  readonly apiName?: string;
  readonly displayName?: string;
  readonly status?: readonly ObjectTypeStatus[];
  readonly visibility?: readonly ('PUBLIC' | 'PRIVATE' | 'INTERNAL')[];
  readonly tags?: readonly string[];
  readonly category?: string;
}

/**
 * ObjectType 스토어 상태
 */
interface ObjectTypeStoreState {
  // 엔티티 상태
  readonly items: Map<ObjectId, ObjectType>;
  readonly loadingState: LoadingState;
  readonly error: StoreError | null;
  
  // 리스트 상태
  readonly pagination: PaginationState;
  readonly filter: FilterState<ObjectTypeFilter>;
  readonly sort: SortState;
  readonly selection: SelectionState<ObjectId>;
  
  // 상세 뷰 상태
  readonly currentObjectType: ObjectType | null;
  readonly isDetailLoading: boolean;
}

/**
 * ObjectType 스토어 액션
 */
interface ObjectTypeStoreActions {
  // 엔티티 액션
  setItems: (items: readonly ObjectType[]) => void;
  addItem: (item: ObjectType) => void;
  updateItem: (id: ObjectId, updates: Partial<ObjectType>) => void;
  removeItem: (id: ObjectId) => void;
  clearItems: () => void;
  setLoadingState: (state: LoadingState) => void;
  setError: (error: StoreError | null) => void;
  
  // 리스트 액션
  setPagination: (pagination: PaginationState) => void;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  setFilter: (filter: Partial<ObjectTypeFilter>) => void;
  clearFilter: () => void;
  setSort: (field: string, direction: 'ASC' | 'DESC') => void;
  select: (id: ObjectId) => void;
  deselect: (id: ObjectId) => void;
  selectAll: () => void;
  clearSelection: () => void;
  
  // 상세 뷰 액션
  setCurrentObjectType: (objectType: ObjectType | null) => void;
  setDetailLoading: (loading: boolean) => void;
  
  // 유틸리티
  getItemById: (id: ObjectId) => ObjectType | undefined;
  getSelectedItems: () => readonly ObjectType[];
}

/**
 * ObjectType 스토어
 * 철칙: 불변성 보장, 단일 책임, 타입 안전성
 */
export const useObjectTypeStore = create<ObjectTypeStoreState & ObjectTypeStoreActions>()(
  immer((set, get) => ({
    // 초기 상태
    items: new Map(),
    loadingState: 'idle',
    error: null,
    
    pagination: {
      page: 1,
      pageSize: 20,
      totalPages: 0,
      totalCount: 0,
      hasNext: false,
      hasPrevious: false,
    },
    
    filter: {
      active: false,
      filters: {},
    },
    
    sort: {
      field: 'updatedAt',
      direction: 'DESC',
    },
    
    selection: {
      selected: new Set(),
      lastSelected: null,
    },
    
    currentObjectType: null,
    isDetailLoading: false,
    
    // 엔티티 액션 구현
    setItems: (items) => set((state) => {
      state.items.clear();
      items.forEach((item) => {
        state.items.set(item.id, item);
      });
    }),
    
    addItem: (item) => set((state) => {
      state.items.set(item.id, item);
    }),
    
    updateItem: (id, updates) => set((state) => {
      const item = state.items.get(id);
      if (item) {
        state.items.set(id, { ...item, ...updates });
      }
    }),
    
    removeItem: (id) => set((state) => {
      state.items.delete(id);
      state.selection.selected.delete(id);
    }),
    
    clearItems: () => set((state) => {
      state.items.clear();
      state.selection.selected.clear();
      state.selection.lastSelected = null;
    }),
    
    setLoadingState: (loadingState) => set((state) => {
      state.loadingState = loadingState;
    }),
    
    setError: (error) => set((state) => {
      state.error = error;
    }),
    
    // 리스트 액션 구현
    setPagination: (pagination) => set((state) => {
      state.pagination = pagination;
    }),
    
    setPage: (page) => set((state) => {
      state.pagination.page = page;
    }),
    
    setPageSize: (pageSize) => set((state) => {
      state.pagination.pageSize = pageSize;
      state.pagination.page = 1; // 페이지 크기 변경 시 첫 페이지로
    }),
    
    setFilter: (filter) => set((state) => {
      state.filter.filters = { ...state.filter.filters, ...filter };
      state.filter.active = true;
      state.pagination.page = 1; // 필터 변경 시 첫 페이지로
    }),
    
    clearFilter: () => set((state) => {
      state.filter.filters = {};
      state.filter.active = false;
      state.pagination.page = 1;
    }),
    
    setSort: (field, direction) => set((state) => {
      state.sort.field = field;
      state.sort.direction = direction;
      state.pagination.page = 1; // 정렬 변경 시 첫 페이지로
    }),
    
    select: (id) => set((state) => {
      state.selection.selected.add(id);
      state.selection.lastSelected = id;
    }),
    
    deselect: (id) => set((state) => {
      state.selection.selected.delete(id);
      if (state.selection.lastSelected === id) {
        state.selection.lastSelected = null;
      }
    }),
    
    selectAll: () => set((state) => {
      state.items.forEach((_, id) => {
        state.selection.selected.add(id);
      });
    }),
    
    clearSelection: () => set((state) => {
      state.selection.selected.clear();
      state.selection.lastSelected = null;
    }),
    
    // 상세 뷰 액션 구현
    setCurrentObjectType: (objectType) => set((state) => {
      state.currentObjectType = objectType;
    }),
    
    setDetailLoading: (loading) => set((state) => {
      state.isDetailLoading = loading;
    }),
    
    // 유틸리티 구현
    getItemById: (id) => get().items.get(id),
    
    getSelectedItems: () => {
      const state = get();
      const selected: ObjectType[] = [];
      state.selection.selected.forEach((id) => {
        const item = state.items.get(id);
        if (item) {
          selected.push(item);
        }
      });
      return selected;
    },
  }))
);