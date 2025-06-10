import { gql } from '@apollo/client';
import { OBJECT_TYPE_FULL_FRAGMENT } from '../fragments/object-type.fragments';

/**
 * ObjectType 뮤테이션 정의
 * 철칙: 각 뮤테이션은 단일 작업만 수행
 */

// ObjectType 생성
export const CREATE_OBJECT_TYPE = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation CreateObjectType($input: CreateObjectTypeInput!) {
    createObjectType(input: $input) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// ObjectType 수정
export const UPDATE_OBJECT_TYPE = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation UpdateObjectType($id: ID!, $input: UpdateObjectTypeInput!) {
    updateObjectType(id: $id, input: $input) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// ObjectType 삭제
export const DELETE_OBJECT_TYPE = gql`
  mutation DeleteObjectType($id: ID!, $force: Boolean = false) {
    deleteObjectType(id: $id, force: $force) {
      success
      message
      errors {
        field
        message
        code
      }
    }
  }
`;

// ObjectType 상태 변경
export const CHANGE_OBJECT_TYPE_STATUS = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation ChangeObjectTypeStatus($id: ID!, $status: ObjectTypeStatus!) {
    changeObjectTypeStatus(id: $id, status: $status) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// Property 추가
export const ADD_PROPERTY = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation AddProperty($objectTypeId: ID!, $input: CreatePropertyInput!) {
    addProperty(objectTypeId: $objectTypeId, input: $input) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// Property 수정
export const UPDATE_PROPERTY = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation UpdateProperty(
    $objectTypeId: ID!
    $propertyId: ID!
    $input: UpdatePropertyInput!
  ) {
    updateProperty(
      objectTypeId: $objectTypeId
      propertyId: $propertyId
      input: $input
    ) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// Property 삭제
export const REMOVE_PROPERTY = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation RemoveProperty($objectTypeId: ID!, $propertyId: ID!) {
    removeProperty(objectTypeId: $objectTypeId, propertyId: $propertyId) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;

// Property 순서 변경
export const REORDER_PROPERTIES = gql`
  ${OBJECT_TYPE_FULL_FRAGMENT}
  
  mutation ReorderProperties($objectTypeId: ID!, $propertyIds: [ID!]!) {
    reorderProperties(objectTypeId: $objectTypeId, propertyIds: $propertyIds) {
      objectType {
        ...ObjectTypeFull
      }
      errors {
        field
        message
        code
      }
    }
  }
`;