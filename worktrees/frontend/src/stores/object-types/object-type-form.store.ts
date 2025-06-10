import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import type { ObjectTypeFormData, Property, PropertyId } from '@/types/domain';
import type { FormState, StoreError } from '../types';

/**
 * Property 폼 데이터
 */
interface PropertyFormData {
  readonly id: PropertyId;
  readonly apiName: string;
  readonly displayName: string;
  readonly description: string;
  readonly dataType: string;
  readonly required: boolean;
  readonly unique: boolean;
  readonly multiValued: boolean;
  readonly searchable: boolean;
  readonly defaultValue: string;
}

/**
 * ObjectType 폼 스토어 상태
 */
interface ObjectTypeFormStoreState {
  // 기본 폼 상태
  readonly form: FormState<ObjectTypeFormData>;
  
  // 속성 관리
  readonly properties: readonly Property[];
  readonly editingProperty: PropertyFormData | null;
  readonly propertyErrors: Map<PropertyId, Map<keyof PropertyFormData, string>>;
  
  // 유효성 검사 상태
  readonly isValidating: boolean;
  readonly validationErrors: readonly StoreError[];
}

/**
 * ObjectType 폼 스토어 액션
 */
interface ObjectTypeFormStoreActions {
  // 폼 액션
  setField: <K extends keyof ObjectTypeFormData>(field: K, value: ObjectTypeFormData[K]) => void;
  setFields: (fields: Partial<ObjectTypeFormData>) => void;
  setError: (field: keyof ObjectTypeFormData, error: string | null) => void;
  touchField: (field: keyof ObjectTypeFormData) => void;
  resetForm: (data?: ObjectTypeFormData) => void;
  validateForm: () => boolean;
  
  // 속성 액션
  addProperty: (property: Property) => void;
  updateProperty: (id: PropertyId, updates: Partial<Property>) => void;
  removeProperty: (id: PropertyId) => void;
  reorderProperties: (propertyIds: readonly PropertyId[]) => void;
  setEditingProperty: (property: PropertyFormData | null) => void;
  setPropertyError: (propertyId: PropertyId, field: keyof PropertyFormData, error: string | null) => void;
  
  // 유효성 검사
  validateApiName: (apiName: string) => string | null;
  validateDisplayName: (displayName: string) => string | null;
  validateProperty: (property: PropertyFormData) => boolean;
  setValidationErrors: (errors: readonly StoreError[]) => void;
}

/**
 * 초기 폼 데이터
 */
const initialFormData: ObjectTypeFormData = {
  apiName: '',
  displayName: '',
  pluralDisplayName: '',
  description: '',
  icon: '',
  color: '#2D72D2',
  tags: [],
  category: '',
  visibility: 'PUBLIC',
};

/**
 * API 이름 유효성 검사
 * 철칙: 순수 함수, 부작용 없음
 */
const validateApiNameFormat = (apiName: string): string | null => {
  if (!apiName) {
    return 'API name is required';
  }
  if (apiName.length < 3) {
    return 'API name must be at least 3 characters';
  }
  if (apiName.length > 50) {
    return 'API name must be less than 50 characters';
  }
  if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(apiName)) {
    return 'API name must start with a letter and contain only letters, numbers, and underscores';
  }
  return null;
};

/**
 * 표시 이름 유효성 검사
 */
const validateDisplayNameFormat = (displayName: string): string | null => {
  if (!displayName) {
    return 'Display name is required';
  }
  if (displayName.length < 2) {
    return 'Display name must be at least 2 characters';
  }
  if (displayName.length > 100) {
    return 'Display name must be less than 100 characters';
  }
  return null;
};

/**
 * ObjectType 폼 스토어
 * 철칙: 불변성, 타입 안전성, 단일 책임
 */
export const useObjectTypeFormStore = create<ObjectTypeFormStoreState & ObjectTypeFormStoreActions>()(
  immer((set, get) => ({
    // 초기 상태
    form: {
      data: initialFormData,
      isDirty: false,
      isValid: false,
      errors: new Map(),
      touched: new Set(),
    },
    
    properties: [],
    editingProperty: null,
    propertyErrors: new Map(),
    isValidating: false,
    validationErrors: [],
    
    // 폼 액션 구현
    setField: (field, value) => set((state) => {
      state.form.data[field] = value;
      state.form.isDirty = true;
      state.form.touched.add(field);
      
      // 실시간 유효성 검사
      const error = get()[field === 'apiName' ? 'validateApiName' : field === 'displayName' ? 'validateDisplayName' : '']?.(value as string);
      if (error) {
        state.form.errors.set(field, error);
      } else {
        state.form.errors.delete(field);
      }
    }),
    
    setFields: (fields) => set((state) => {
      Object.entries(fields).forEach(([key, value]) => {
        const field = key as keyof ObjectTypeFormData;
        if (value !== undefined) {
          state.form.data[field] = value as never;
          state.form.touched.add(field);
        }
      });
      state.form.isDirty = true;
    }),
    
    setError: (field, error) => set((state) => {
      if (error) {
        state.form.errors.set(field, error);
      } else {
        state.form.errors.delete(field);
      }
      state.form.isValid = state.form.errors.size === 0;
    }),
    
    touchField: (field) => set((state) => {
      state.form.touched.add(field);
    }),
    
    resetForm: (data) => set((state) => {
      state.form.data = data ?? initialFormData;
      state.form.isDirty = false;
      state.form.isValid = false;
      state.form.errors.clear();
      state.form.touched.clear();
      state.properties = [];
      state.editingProperty = null;
      state.propertyErrors.clear();
    }),
    
    validateForm: () => {
      const state = get();
      const errors = new Map<keyof ObjectTypeFormData, string>();
      
      // API 이름 검증
      const apiNameError = validateApiNameFormat(state.form.data.apiName);
      if (apiNameError) errors.set('apiName', apiNameError);
      
      // 표시 이름 검증
      const displayNameError = validateDisplayNameFormat(state.form.data.displayName);
      if (displayNameError) errors.set('displayName', displayNameError);
      
      // 복수형 표시 이름 검증
      if (!state.form.data.pluralDisplayName) {
        errors.set('pluralDisplayName', 'Plural display name is required');
      }
      
      set((state) => {
        state.form.errors = errors;
        state.form.isValid = errors.size === 0;
      });
      
      return errors.size === 0;
    },
    
    // 속성 액션 구현
    addProperty: (property) => set((state) => {
      state.properties = [...state.properties, property];
      state.form.isDirty = true;
    }),
    
    updateProperty: (id, updates) => set((state) => {
      const index = state.properties.findIndex(p => p.id === id);
      if (index !== -1) {
        state.properties = [
          ...state.properties.slice(0, index),
          { ...state.properties[index], ...updates },
          ...state.properties.slice(index + 1),
        ];
        state.form.isDirty = true;
      }
    }),
    
    removeProperty: (id) => set((state) => {
      state.properties = state.properties.filter(p => p.id !== id);
      state.propertyErrors.delete(id);
      state.form.isDirty = true;
    }),
    
    reorderProperties: (propertyIds) => set((state) => {
      const propertyMap = new Map(state.properties.map(p => [p.id, p]));
      state.properties = propertyIds
        .map(id => propertyMap.get(id))
        .filter((p): p is Property => p !== undefined);
      state.form.isDirty = true;
    }),
    
    setEditingProperty: (property) => set((state) => {
      state.editingProperty = property;
    }),
    
    setPropertyError: (propertyId, field, error) => set((state) => {
      if (!state.propertyErrors.has(propertyId)) {
        state.propertyErrors.set(propertyId, new Map());
      }
      const fieldErrors = state.propertyErrors.get(propertyId)!;
      if (error) {
        fieldErrors.set(field, error);
      } else {
        fieldErrors.delete(field);
      }
    }),
    
    // 유효성 검사 구현
    validateApiName: validateApiNameFormat,
    validateDisplayName: validateDisplayNameFormat,
    
    validateProperty: (property) => {
      const errors = new Map<keyof PropertyFormData, string>();
      
      if (!property.apiName) {
        errors.set('apiName', 'API name is required');
      } else if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(property.apiName)) {
        errors.set('apiName', 'Invalid API name format');
      }
      
      if (!property.displayName) {
        errors.set('displayName', 'Display name is required');
      }
      
      if (!property.dataType) {
        errors.set('dataType', 'Data type is required');
      }
      
      set((state) => {
        state.propertyErrors.set(property.id, errors);
      });
      
      return errors.size === 0;
    },
    
    setValidationErrors: (errors) => set((state) => {
      state.validationErrors = errors;
    }),
  }))
);