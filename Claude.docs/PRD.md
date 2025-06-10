# Ontology Metadata Service (OMS) Product Requirements Document

## Executive Summary

Ontology Metadata Service (OMS)는 OpenFoundry 플랫폼의 핵심 메타데이터 관리 서비스로, 조직의 비즈니스 개체와 관계를 정의하는 스키마 레지스트리입니다. Palantir Foundry의 Ontology 철학을 계승하여, 데이터를 현실 세계의 객체로 추상화하는 구조적 정의를 제공합니다.

## 1. 제품 개요

### 1.1 목적
OMS는 비즈니스 객체의 구조적 정의와 메타데이터를 중앙 관리하는 단일 진실 원천(Single Source of Truth)입니다.

### 1.2 핵심 가치
- **스키마 레지스트리**: 조직 전체의 통일된 객체 정의
- **비즈니스 친화적**: 기술 용어가 아닌 비즈니스 용어로 모델링
- **확장 가능한 설계**: 향후 서비스 통합을 위한 견고한 기반

### 1.3 MVP 범위
- 객체 타입 정의 및 관리
- 속성 스키마 정의
- 객체 간 관계(링크) 정의
- 메타데이터 관리
- API 제공 (GraphQL/REST)

## 2. 사용자 정의

### 2.1 주요 사용자

#### 데이터 모델러
- **목표**: 조직의 비즈니스 모델을 디지털로 표현
- **니즈**: 직관적인 모델링 도구, 유연한 스키마 정의
- **과제**: 복잡한 비즈니스 관계의 정확한 표현

#### API 개발자
- **목표**: 정의된 스키마를 프로그래밍적으로 활용
- **니즈**: 명확한 API 문서, 일관된 인터페이스
- **과제**: 동적 스키마 변경에 대한 대응

### 2.2 핵심 사용 시나리오

#### 시나리오: 전자상거래 도메인 모델링
```
1. 'Customer' 객체 타입 생성
   - 속성: customerId, name, email, registeredDate
   
2. 'Product' 객체 타입 생성
   - 속성: productId, name, price, category
   
3. 'Order' 객체 타입 생성
   - 속성: orderId, orderDate, totalAmount, status
   
4. 관계 정의
   - Customer -> Order (1:N)
   - Order -> Product (N:N, through OrderItem)
   
결과: 완전한 전자상거래 도메인 모델
```

## 3. 기능 요구사항

### 3.1 객체 타입 관리

#### 3.1.1 객체 타입 CRUD
**생성 (Create)**
- 고유 ID 자동 생성 (UUID v4)
- 필수 필드: name, displayName
- 선택 필드: description, category, tags
- 이름 중복 검증

**조회 (Read)**
- 단일 객체 타입 조회 (by ID)
- 목록 조회 (페이지네이션, 필터, 정렬)
- 전문 검색 (이름, 설명, 태그)

**수정 (Update)**
- 부분 업데이트 지원
- 버전 자동 증가
- 수정 이력 기록

**삭제 (Delete)**
- Soft delete (실제 삭제 대신 비활성화)
- 의존성 검증 (사용 중인 경우 삭제 방지)

#### 3.1.2 속성 정의

**지원 데이터 타입**
```typescript
enum DataType {
  STRING = 'STRING',
  NUMBER = 'NUMBER',
  BOOLEAN = 'BOOLEAN',
  DATE = 'DATE',
  DATETIME = 'DATETIME',
  ARRAY = 'ARRAY',
  OBJECT = 'OBJECT',
  REFERENCE = 'REFERENCE'  // 다른 객체 참조
}
```

**속성 설정**
```typescript
interface PropertyDefinition {
  name: string;           // 속성 이름
  displayName: string;    // 표시 이름
  dataType: DataType;     // 데이터 타입
  required: boolean;      // 필수 여부
  unique: boolean;        // 고유값 여부
  indexed: boolean;       // 인덱스 여부
  defaultValue?: any;     // 기본값
  description?: string;   // 설명
}
```

#### 3.1.3 관계(링크) 정의

**관계 유형**
- ONE_TO_ONE: 1:1 관계
- ONE_TO_MANY: 1:N 관계
- MANY_TO_MANY: N:N 관계

**관계 설정**
```typescript
interface LinkDefinition {
  name: string;                  // 관계 이름
  displayName: string;           // 표시 이름
  sourceObjectTypeId: string;    // 출발 객체
  targetObjectTypeId: string;    // 도착 객체
  cardinality: Cardinality;      // 관계 유형
  description?: string;          // 설명
}
```

### 3.2 메타데이터 관리

#### 3.2.1 시스템 메타데이터
- createdAt: 생성 시간
- createdBy: 생성자
- updatedAt: 수정 시간
- updatedBy: 수정자
- version: 버전 번호

#### 3.2.2 비즈니스 메타데이터
- tags: 태그 목록
- category: 분류
- businessOwner: 비즈니스 담당자
- dataSource: 데이터 출처

### 3.3 버전 관리

#### 3.3.1 변경 추적
- 모든 수정사항 자동 기록
- 변경 내용 상세 (이전값/이후값)
- 변경 사유 기록 옵션

#### 3.3.2 버전 조회
- 특정 버전 조회
- 버전 간 차이 비교
- 변경 이력 타임라인

### 3.4 검색 기능

#### 3.4.1 기본 검색
- 이름으로 검색 (exact match)
- 와일드카드 검색 지원

#### 3.4.2 고급 검색
- 전문 검색 (이름, 설명, 태그)
- 필터 조합 (AND/OR)
- 정렬 옵션 (이름, 생성일, 수정일)

## 4. API 명세

### 4.1 GraphQL API

#### 4.1.1 스키마 정의
```graphql
type ObjectType {
  id: ID!
  name: String!
  displayName: String!
  description: String
  category: String
  tags: [String!]
  properties: [Property!]!
  version: Int!
  metadata: Metadata!
}

type Property {
  id: ID!
  name: String!
  displayName: String!
  dataType: DataType!
  required: Boolean!
  unique: Boolean!
  indexed: Boolean!
  defaultValue: JSON
  description: String
}

type LinkType {
  id: ID!
  name: String!
  displayName: String!
  sourceObjectType: ObjectType!
  targetObjectType: ObjectType!
  cardinality: Cardinality!
  description: String
}
```

#### 4.1.2 Query Operations
```graphql
type Query {
  # 객체 타입 조회
  objectType(id: ID!): ObjectType
  objectTypes(
    filter: ObjectTypeFilter
    pagination: PaginationInput
    sort: SortInput
  ): ObjectTypeConnection!
  
  # 링크 타입 조회
  linkType(id: ID!): LinkType
  linkTypes(
    filter: LinkTypeFilter
  ): [LinkType!]!
  
  # 검색
  searchObjectTypes(
    query: String!
    limit: Int = 10
  ): [ObjectType!]!
}
```

#### 4.1.3 Mutation Operations
```graphql
type Mutation {
  # 객체 타입 관리
  createObjectType(
    input: CreateObjectTypeInput!
  ): ObjectType!
  
  updateObjectType(
    id: ID!
    input: UpdateObjectTypeInput!
  ): ObjectType!
  
  deleteObjectType(id: ID!): Boolean!
  
  # 링크 타입 관리
  createLinkType(
    input: CreateLinkTypeInput!
  ): LinkType!
  
  updateLinkType(
    id: ID!
    input: UpdateLinkTypeInput!
  ): LinkType!
  
  deleteLinkType(id: ID!): Boolean!
}
```

### 4.2 REST API (호환성)

#### 4.2.1 엔드포인트
```
GET    /api/v1/object-types
POST   /api/v1/object-types
GET    /api/v1/object-types/{id}
PUT    /api/v1/object-types/{id}
DELETE /api/v1/object-types/{id}

GET    /api/v1/link-types
POST   /api/v1/link-types
GET    /api/v1/link-types/{id}
PUT    /api/v1/link-types/{id}
DELETE /api/v1/link-types/{id}
```

### 4.3 이벤트 발행

#### 4.3.1 이벤트 유형
- ObjectTypeCreated
- ObjectTypeUpdated
- ObjectTypeDeleted
- LinkTypeCreated
- LinkTypeUpdated
- LinkTypeDeleted

#### 4.3.2 이벤트 구조
```json
{
  "eventId": "uuid",
  "eventType": "ObjectTypeCreated",
  "timestamp": "2024-01-01T00:00:00Z",
  "actor": "user-id",
  "data": {
    "objectType": { ... }
  }
}
```

## 5. 비기능 요구사항

### 5.1 성능
- 단일 조회: < 50ms (P95)
- 목록 조회: < 100ms (P95)
- 생성/수정: < 200ms (P95)
- 동시 요청: 100 RPS

### 5.2 확장성
- 객체 타입: 최대 10,000개
- 속성: 객체당 최대 1,000개
- 수평 확장 가능 설계

### 5.3 가용성
- 목표: 99.9% uptime
- 자동 장애 복구
- 우아한 성능 저하

### 5.4 보안
- JWT 기반 인증
- API 키 지원
- TLS 1.3 암호화
- 감사 로그

## 6. 기술 제약사항

### 6.1 기술 스택
- 언어: Go 또는 Node.js/TypeScript
- 데이터베이스: PostgreSQL 13+
- 캐시: Redis
- API: GraphQL (Apollo Server)

### 6.2 인프라
- Container: Docker
- Orchestration: Kubernetes
- 메시지 큐: Apache Kafka

## 7. 데이터 모델

### 7.1 핵심 테이블

#### object_types
```sql
CREATE TABLE object_types (
  id UUID PRIMARY KEY,
  name VARCHAR(64) UNIQUE NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(64),
  tags JSONB,
  properties JSONB NOT NULL,
  metadata JSONB NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  updated_by VARCHAR(255) NOT NULL
);
```

#### link_types
```sql
CREATE TABLE link_types (
  id UUID PRIMARY KEY,
  name VARCHAR(64) UNIQUE NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  source_object_type_id UUID NOT NULL,
  target_object_type_id UUID NOT NULL,
  cardinality VARCHAR(32) NOT NULL,
  description TEXT,
  metadata JSONB NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  
  FOREIGN KEY (source_object_type_id) 
    REFERENCES object_types(id),
  FOREIGN KEY (target_object_type_id) 
    REFERENCES object_types(id)
);
```

## 8. 마일스톤

### Week 1-2: 기초 구현
- 데이터베이스 스키마
- 기본 CRUD API
- 단위 테스트

### Week 3-4: API 완성
- GraphQL 스키마
- REST 래퍼
- 통합 테스트

### Week 5-6: 프로덕션 준비
- 성능 최적화
- 보안 강화
- 문서화
- 배포 준비

## 9. 성공 지표

### 기술적 지표
- API 응답 시간 목표 달성
- 99.9% 가용성
- 제로 크리티컬 버그
- 80% 테스트 커버리지

### 비즈니스 지표
- 5개 도메인 모델링 완료
- API 문서 100% 완성
- 개발자 만족도 4.0/5.0

## 10. 리스크

| 리스크 | 영향 | 완화 방안 |
|--------|------|-----------|
| 스키마 복잡도 | 높음 | 명확한 가이드라인, 검증 규칙 |
| 성능 저하 | 중간 | 캐싱 전략, 인덱스 최적화 |
| 보안 취약점 | 높음 | 정기 보안 감사, 자동 스캔 |
| MSA 복잡도 | 높음 | 단계적 전환, 서비스 메시 도입 |

## 11. 마이크로서비스 아키텍처

### 11.1 아키텍처 결정

OMS는 OpenFoundry 플랫폼의 핵심 마이크로서비스로 설계되어, 향후 추가될 서비스들과의 유연한 통합을 지원합니다.

#### MSA 채택 이유
- **독립적 확장성**: 메타데이터 서비스의 부하 패턴에 맞춘 최적화
- **기술 독립성**: 각 서비스별 최적의 기술 스택 선택 가능
- **팀 자율성**: 독립적인 개발 및 배포 주기
- **장애 격리**: 한 서비스의 장애가 전체 시스템에 미치는 영향 최소화

### 11.2 서비스 포지셔닝

```yaml
OMS의 역할:
  - 기반 서비스: 다른 모든 서비스가 의존하는 메타데이터 레지스트리
  - 스키마 권한: 비즈니스 객체의 구조와 관계 정의의 단일 진실 원천
  - 이벤트 소스: 메타데이터 변경 이벤트의 발행자
  
향후 통합될 서비스:
  - OSv2: 실제 객체 인스턴스 데이터 관리
  - OSS: 통합 검색 및 쿼리 서비스
  - FOO: 객체 기반 함수 실행 플랫폼
```

### 11.3 통신 인터페이스

#### 외부 API (클라이언트용)
- **GraphQL**: 유연한 쿼리와 구독을 위한 메인 API
- **REST**: 레거시 시스템 호환성을 위한 래퍼 API

#### 내부 API (서비스간 통신)
- **gRPC**: 고성능 서비스간 통신
- **Event Streaming**: Kafka를 통한 비동기 이벤트 전파

### 11.4 데이터 일관성 전략

```yaml
일관성 모델:
  - 강한 일관성: OMS 내부 데이터 (PostgreSQL 트랜잭션)
  - 최종 일관성: 서비스간 데이터 동기화 (이벤트 기반)
  
캐싱 전략:
  - L1 캐시: 애플리케이션 메모리 (5분)
  - L2 캐시: Redis (30분)
  - 캐시 무효화: 이벤트 기반 즉시 무효화
```

### 11.5 확장성 고려사항

#### 수평 확장
- 무상태 서비스 설계로 Pod 수평 확장 지원
- 읽기 전용 레플리카를 통한 조회 성능 향상
- 샤딩 준비된 데이터 모델

#### 멀티테넌시
- 테넌트별 네임스페이스 격리
- 리소스 쿼터 및 제한 설정
- 테넌트별 메타데이터 격리

### 11.6 운영 및 모니터링

#### 관찰성 (Observability)
- **분산 추적**: Jaeger를 통한 요청 흐름 추적
- **메트릭**: Prometheus/Grafana 대시보드
- **로깅**: 구조화된 로그 (JSON) + ELK 스택

#### SLA 목표
- 가용성: 99.9% (월간 43분 다운타임 허용)
- 응답 시간: P95 < 100ms (조회), P95 < 200ms (생성/수정)
- 오류율: < 0.1%

### 11.7 배포 전략

#### 컨테이너화
- Docker 이미지 크기 최적화 (< 50MB)
- 멀티스테이지 빌드
- 보안 스캔 자동화

#### 오케스트레이션
- Kubernetes 네이티브 배포
- Helm 차트를 통한 패키지화
- GitOps (ArgoCD) 기반 배포

#### 배포 패턴
- Blue-Green 배포: 프로덕션 무중단 배포
- Canary 배포: 단계적 롤아웃 (5% → 25% → 50% → 100%)
- Feature Flags: 기능별 점진적 활성화