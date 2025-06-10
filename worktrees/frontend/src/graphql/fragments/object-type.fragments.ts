import { gql } from '@apollo/client';

/**
 * ObjectType 프래그먼트 정의
 * 철칙: 재사용 가능한 작은 단위로 분리
 */

// 기본 ObjectType 정보
export const OBJECT_TYPE_BASIC_FRAGMENT = gql`
  fragment ObjectTypeBasic on ObjectType {
    id
    rid
    apiName
    displayName
    pluralDisplayName
    description
    icon
    color
    status
    version
  }
`;

// ObjectType 메타데이터
export const OBJECT_TYPE_METADATA_FRAGMENT = gql`
  fragment ObjectTypeMetadata on ObjectTypeMetadata {
    tags
    category
    datasources
    visibility
    permissions {
      create
      read
      update
      delete
    }
    capabilities {
      searchable
      writeable
      deletable
      versionable
      auditable
      timeSeries
      attachments
    }
  }
`;

// Property 기본 정보
export const PROPERTY_BASIC_FRAGMENT = gql`
  fragment PropertyBasic on Property {
    id
    rid
    apiName
    displayName
    description
    dataType
    constraints {
      required
      unique
      multiValued
      searchable
      primaryKey
      minValue
      maxValue
      minLength
      maxLength
      pattern
      enumValues
    }
    metadata {
      tags
      category
      defaultValue
      formula
      visibility
      deprecated
      deprecationMessage
    }
    version
  }
`;

// 전체 ObjectType 정보 (모든 속성 포함)
export const OBJECT_TYPE_FULL_FRAGMENT = gql`
  ${OBJECT_TYPE_BASIC_FRAGMENT}
  ${OBJECT_TYPE_METADATA_FRAGMENT}
  ${PROPERTY_BASIC_FRAGMENT}
  
  fragment ObjectTypeFull on ObjectType {
    ...ObjectTypeBasic
    titleProperty
    subtitleProperty
    properties {
      ...PropertyBasic
    }
    metadata {
      ...ObjectTypeMetadata
    }
    createdAt
    updatedAt
    createdBy
    updatedBy
  }
`;

// ObjectType 리스트용 경량 프래그먼트
export const OBJECT_TYPE_LIST_ITEM_FRAGMENT = gql`
  fragment ObjectTypeListItem on ObjectType {
    id
    rid
    apiName
    displayName
    description
    icon
    color
    status
    properties {
      id
    }
    metadata {
      tags
      visibility
    }
    updatedAt
  }
`;