import { gql } from '@apollo/client';
import {
  OBJECT_TYPE_LIST_ITEM_FRAGMENT,
  OBJECT_TYPE_FULL_FRAGMENT,
} from '../fragments/object-type.fragments';

/**
 * ObjectType 쿼리 정의
 * 철칙: 명시적 반환 타입, 작은 단위로 분리
 */

// ObjectType 목록 조회
export const GET_OBJECT_TYPES = gql`
  ${OBJECT_TYPE_LIST_ITEM_FRAGMENT}
  
  query GetObjectTypes(
    $filter: ObjectTypeFilter
    $sort: [SortCriteria!]
    $pagination: PaginationInput!
  ) {
    objectTypes(filter: $filter, sort: $sort, pagination: $pagination) {
      data {
        ...ObjectTypeListItem
      }
      pagination {
        page
        pageSize
        totalPages
        totalCount
        hasNext
        hasPrevious
      }
    }
  }
`;

// 단일 ObjectType 상세 조회
export const GET_OBJECT_TYPE = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  query GetObjectType($id: ID!) {
    objectType(id: $id) {
      ...ObjectTypeFull
    }
  }
`;

// ObjectType 이름으로 조회
export const GET_OBJECT_TYPE_BY_API_NAME = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  query GetObjectTypeByApiName($apiName: String!) {
    objectTypeByApiName(apiName: $apiName) {
      ...ObjectTypeFull
    }
  }
`;

// ObjectType 검색
export const SEARCH_OBJECT_TYPES = gql`
  ${OBJECT_TYPE_LIST_ITEM_FRAGMENT}
  
  query SearchObjectTypes($query: String!, $limit: Int = 10) {
    searchObjectTypes(query: $query, limit: $limit) {
      ...ObjectTypeListItem
    }
  }
`;

// ObjectType 통계
export const GET_OBJECT_TYPE_STATS = gql`
  query GetObjectTypeStats {
    objectTypeStats {
      totalCount
      activeCount
      draftCount
      deprecatedCount
      archivedCount
      byCategory {
        category
        count
      }
      byVisibility {
        visibility
        count
      }
    }
  }
`;

// ObjectType 의존성 확인
export const CHECK_OBJECT_TYPE_DEPENDENCIES = gql`
  query CheckObjectTypeDependencies($id: ID!) {
    objectTypeDependencies(id: $id) {
      hasIncomingLinks
      hasOutgoingLinks
      incomingLinkTypes {
        id
        apiName
        displayName
        sourceObjectType
      }
      outgoingLinkTypes {
        id
        apiName
        displayName
        targetObjectType
      }
      dependentObjectsCount
    }
  }
`;