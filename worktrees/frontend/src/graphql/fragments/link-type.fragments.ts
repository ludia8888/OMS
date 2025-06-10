import { gql } from '@apollo/client';

/**
 * LinkType 프래그먼트 정의
 * 철칙: 단일 책임 원칙 - 각 프래그먼트는 하나의 목적만 수행
 */

// 기본 LinkType 정보
export const LINK_TYPE_BASIC_FRAGMENT = gql`
  fragment LinkTypeBasic on LinkType {
    id
    rid
    apiName
    displayName
    reverseDisplayName
    description
    cardinality
    required
    cascadeDelete
    status
    version
  }
`;

// LinkType 메타데이터
export const LINK_TYPE_METADATA_FRAGMENT = gql`
  fragment LinkTypeMetadata on LinkTypeMetadata {
    tags
    category
    visibility
    permissions {
      create
      read
      delete
    }
  }
`;

// 전체 LinkType 정보
export const LINK_TYPE_FULL_FRAGMENT = gql`
  ${LINK_TYPE_BASIC_FRAGMENT}
  ${LINK_TYPE_METADATA_FRAGMENT}
  
  fragment LinkTypeFull on LinkType {
    ...LinkTypeBasic
    sourceObjectType
    targetObjectType
    metadata {
      ...LinkTypeMetadata
    }
    createdAt
    updatedAt
    createdBy
    updatedBy
  }
`;

// LinkType 리스트용 경량 프래그먼트
export const LINK_TYPE_LIST_ITEM_FRAGMENT = gql`
  fragment LinkTypeListItem on LinkType {
    id
    rid
    apiName
    displayName
    reverseDisplayName
    sourceObjectType
    targetObjectType
    cardinality
    status
    metadata {
      tags
      visibility
    }
    updatedAt
  }
`;