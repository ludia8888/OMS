/**
 * OMS 도메인 타입 정의
 * Claude.rules.md 기준: any 타입 절대 금지, unknown 최소화, 엄격한 타입 정의
 */

// 기본 ID 타입
export type ObjectId = string & { readonly __brand: 'ObjectId' };
export type LinkTypeId = string & { readonly __brand: 'LinkTypeId' };
export type PropertyId = string & { readonly __brand: 'PropertyId' };

// 타입 헬퍼 함수
export const createObjectId = (id: string): ObjectId => id as ObjectId;
export const createLinkTypeId = (id: string): LinkTypeId => id as LinkTypeId;
export const createPropertyId = (id: string): PropertyId => id as PropertyId;

// 데이터 타입
export const DataTypes = {
  STRING: 'STRING',
  INTEGER: 'INTEGER',
  DOUBLE: 'DOUBLE',
  BOOLEAN: 'BOOLEAN',
  DATE: 'DATE',
  TIMESTAMP: 'TIMESTAMP',
  ARRAY: 'ARRAY',
  STRUCT: 'STRUCT',
  ATTACHMENT: 'ATTACHMENT',
  MEDIA_REFERENCE: 'MEDIA_REFERENCE',
} as const;

export type DataType = typeof DataTypes[keyof typeof DataTypes];

// 속성 제약 조건
export interface PropertyConstraints {
  readonly required: boolean;
  readonly unique: boolean;
  readonly multiValued: boolean;
  readonly searchable: boolean;
  readonly primaryKey: boolean;
  readonly minValue?: number;
  readonly maxValue?: number;
  readonly minLength?: number;
  readonly maxLength?: number;
  readonly pattern?: string;
  readonly enumValues?: readonly string[];
}

// 속성 정의
export interface Property {
  readonly id: PropertyId;
  readonly rid: string;
  readonly apiName: string;
  readonly displayName: string;
  readonly description?: string;
  readonly dataType: DataType;
  readonly constraints: PropertyConstraints;
  readonly metadata: PropertyMetadata;
  readonly createdAt: string;
  readonly updatedAt: string;
  readonly version: number;
}

// 속성 메타데이터
export interface PropertyMetadata {
  readonly tags: readonly string[];
  readonly category?: string;
  readonly defaultValue?: string | number | boolean;
  readonly formula?: string;
  readonly visibility: 'PUBLIC' | 'PRIVATE' | 'INTERNAL';
  readonly deprecated: boolean;
  readonly deprecationMessage?: string;
}

// 객체 타입 정의
export interface ObjectType {
  readonly id: ObjectId;
  readonly rid: string;
  readonly apiName: string;
  readonly displayName: string;
  readonly pluralDisplayName: string;
  readonly description?: string;
  readonly icon?: string;
  readonly color?: string;
  readonly titleProperty: PropertyId;
  readonly subtitleProperty?: PropertyId;
  readonly properties: readonly Property[];
  readonly metadata: ObjectTypeMetadata;
  readonly status: ObjectTypeStatus;
  readonly createdAt: string;
  readonly updatedAt: string;
  readonly createdBy: string;
  readonly updatedBy: string;
  readonly version: number;
}

// 객체 타입 메타데이터
export interface ObjectTypeMetadata {
  readonly tags: readonly string[];
  readonly category?: string;
  readonly datasources: readonly string[];
  readonly visibility: 'PUBLIC' | 'PRIVATE' | 'INTERNAL';
  readonly permissions: ObjectTypePermissions;
  readonly capabilities: ObjectTypeCapabilities;
}

// 객체 타입 권한
export interface ObjectTypePermissions {
  readonly create: readonly string[];
  readonly read: readonly string[];
  readonly update: readonly string[];
  readonly delete: readonly string[];
}

// 객체 타입 기능
export interface ObjectTypeCapabilities {
  readonly searchable: boolean;
  readonly writeable: boolean;
  readonly deletable: boolean;
  readonly versionable: boolean;
  readonly auditable: boolean;
  readonly timeSeries: boolean;
  readonly attachments: boolean;
}

// 객체 타입 상태
export const ObjectTypeStatuses = {
  DRAFT: 'DRAFT',
  ACTIVE: 'ACTIVE',
  DEPRECATED: 'DEPRECATED',
  ARCHIVED: 'ARCHIVED',
} as const;

export type ObjectTypeStatus = typeof ObjectTypeStatuses[keyof typeof ObjectTypeStatuses];

// 링크 타입 카디널리티
export const Cardinalities = {
  ONE_TO_ONE: 'ONE_TO_ONE',
  ONE_TO_MANY: 'ONE_TO_MANY',
  MANY_TO_ONE: 'MANY_TO_ONE',
  MANY_TO_MANY: 'MANY_TO_MANY',
} as const;

export type Cardinality = typeof Cardinalities[keyof typeof Cardinalities];

// 링크 타입 정의
export interface LinkType {
  readonly id: LinkTypeId;
  readonly rid: string;
  readonly apiName: string;
  readonly displayName: string;
  readonly reverseDisplayName: string;
  readonly description?: string;
  readonly sourceObjectType: ObjectId;
  readonly targetObjectType: ObjectId;
  readonly cardinality: Cardinality;
  readonly required: boolean;
  readonly cascadeDelete: boolean;
  readonly metadata: LinkTypeMetadata;
  readonly status: LinkTypeStatus;
  readonly createdAt: string;
  readonly updatedAt: string;
  readonly createdBy: string;
  readonly updatedBy: string;
  readonly version: number;
}

// 링크 타입 메타데이터
export interface LinkTypeMetadata {
  readonly tags: readonly string[];
  readonly category?: string;
  readonly visibility: 'PUBLIC' | 'PRIVATE' | 'INTERNAL';
  readonly permissions: LinkTypePermissions;
}

// 링크 타입 권한
export interface LinkTypePermissions {
  readonly create: readonly string[];
  readonly read: readonly string[];
  readonly delete: readonly string[];
}

// 링크 타입 상태
export const LinkTypeStatuses = {
  ACTIVE: 'ACTIVE',
  DEPRECATED: 'DEPRECATED',
  ARCHIVED: 'ARCHIVED',
} as const;

export type LinkTypeStatus = typeof LinkTypeStatuses[keyof typeof LinkTypeStatuses];

// 검색 필터
export interface SearchFilter {
  readonly field: string;
  readonly operator: FilterOperator;
  readonly value: string | number | boolean | readonly (string | number | boolean)[];
}

// 필터 연산자
export const FilterOperators = {
  EQUALS: 'EQUALS',
  NOT_EQUALS: 'NOT_EQUALS',
  CONTAINS: 'CONTAINS',
  NOT_CONTAINS: 'NOT_CONTAINS',
  STARTS_WITH: 'STARTS_WITH',
  ENDS_WITH: 'ENDS_WITH',
  IN: 'IN',
  NOT_IN: 'NOT_IN',
  GREATER_THAN: 'GREATER_THAN',
  GREATER_THAN_OR_EQUALS: 'GREATER_THAN_OR_EQUALS',
  LESS_THAN: 'LESS_THAN',
  LESS_THAN_OR_EQUALS: 'LESS_THAN_OR_EQUALS',
  BETWEEN: 'BETWEEN',
  IS_NULL: 'IS_NULL',
  IS_NOT_NULL: 'IS_NOT_NULL',
} as const;

export type FilterOperator = typeof FilterOperators[keyof typeof FilterOperators];

// 정렬 방향
export const SortDirections = {
  ASC: 'ASC',
  DESC: 'DESC',
} as const;

export type SortDirection = typeof SortDirections[keyof typeof SortDirections];

// 정렬 기준
export interface SortCriteria {
  readonly field: string;
  readonly direction: SortDirection;
}

// 페이지네이션
export interface Pagination {
  readonly page: number;
  readonly pageSize: number;
}

// 검색 요청
export interface SearchRequest {
  readonly filters: readonly SearchFilter[];
  readonly sort?: readonly SortCriteria[];
  readonly pagination: Pagination;
  readonly includeCount: boolean;
}

// 검색 응답
export interface SearchResponse<T> {
  readonly data: readonly T[];
  readonly pagination: PaginationInfo;
}

// 페이지네이션 정보
export interface PaginationInfo {
  readonly page: number;
  readonly pageSize: number;
  readonly totalPages: number;
  readonly totalCount: number;
  readonly hasNext: boolean;
  readonly hasPrevious: boolean;
}

// 감사 로그
export interface AuditLog {
  readonly id: string;
  readonly entityType: 'OBJECT_TYPE' | 'LINK_TYPE' | 'PROPERTY';
  readonly entityId: string;
  readonly action: AuditAction;
  readonly userId: string;
  readonly timestamp: string;
  readonly changes: readonly AuditChange[];
  readonly metadata: Record<string, string | number | boolean>;
}

// 감사 액션
export const AuditActions = {
  CREATE: 'CREATE',
  UPDATE: 'UPDATE',
  DELETE: 'DELETE',
  RESTORE: 'RESTORE',
  ARCHIVE: 'ARCHIVE',
  PUBLISH: 'PUBLISH',
  DEPRECATE: 'DEPRECATE',
} as const;

export type AuditAction = typeof AuditActions[keyof typeof AuditActions];

// 감사 변경 사항
export interface AuditChange {
  readonly field: string;
  readonly oldValue: string | number | boolean | null;
  readonly newValue: string | number | boolean | null;
}

// API 에러 응답
export interface ApiError {
  readonly code: string;
  readonly message: string;
  readonly details?: Record<string, string | number | boolean>;
  readonly timestamp: string;
  readonly traceId: string;
}

// 폼 데이터 타입
export interface ObjectTypeFormData {
  readonly apiName: string;
  readonly displayName: string;
  readonly pluralDisplayName: string;
  readonly description: string;
  readonly icon: string;
  readonly color: string;
  readonly tags: readonly string[];
  readonly category: string;
  readonly visibility: 'PUBLIC' | 'PRIVATE' | 'INTERNAL';
}

export interface PropertyFormData {
  readonly apiName: string;
  readonly displayName: string;
  readonly description: string;
  readonly dataType: DataType;
  readonly required: boolean;
  readonly unique: boolean;
  readonly multiValued: boolean;
  readonly searchable: boolean;
  readonly defaultValue: string;
}

export interface LinkTypeFormData {
  readonly apiName: string;
  readonly displayName: string;
  readonly reverseDisplayName: string;
  readonly description: string;
  readonly sourceObjectType: ObjectId;
  readonly targetObjectType: ObjectId;
  readonly cardinality: Cardinality;
  readonly required: boolean;
  readonly cascadeDelete: boolean;
}