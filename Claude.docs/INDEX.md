# OMS (Ontology Metadata Service) Master Index

> **Last Updated**: 2024-01-10  
> **Version**: 1.0.0  
> **Status**: Sprint Ready

## 📚 Document Overview

### Core Documents
1. **[Product Requirements Document (PRD.md)](./PRD.md)** - 제품 요구사항 명세
2. **[Backend Development Specification (Backend.md)](./Backend.md)** - 백엔드 개발 명세
3. **[Frontend Development Specification (Frontend.md)](./Frontend.md)** - 프론트엔드 개발 명세
4. **[Design System Guide (design-system.md)](../design-system.md)** - 디자인 시스템 가이드
5. **[Development Rules (Claude.rules.md)](./Claude.rules.md)** - 개발 규칙 및 가이드라인

## 🎯 Quick Navigation by Role

### 📋 Product Manager (PM)
- **Primary Documents**:
  - [PRD.md](./PRD.md) - 전체 제품 비전 및 요구사항
  - [PRD.md#2-사용자-정의](./PRD.md#2-사용자-정의) - 사용자 페르소나 및 시나리오
  - [PRD.md#8-마일스톤](./PRD.md#8-마일스톤) - 개발 일정 및 마일스톤
  - [PRD.md#9-성공-지표](./PRD.md#9-성공-지표) - KPI 및 성공 지표
  - [PRD.md#11-마이크로서비스-아키텍처](./PRD.md#11-마이크로서비스-아키텍처) - MSA 전략

### 💻 Backend Developer
- **Primary Documents**:
  - [Backend.md](./Backend.md) - 백엔드 구현 전체 가이드
  - [Backend.md#3-database-design](./Backend.md#3-database-design) - 데이터베이스 설계
  - [Backend.md#4-api-implementation](./Backend.md#4-api-implementation) - API 구현 가이드
  - [Backend.md#16-microservices-architecture](./Backend.md#16-microservices-architecture) - MSA 구현
  - [Claude.rules.md](./Claude.rules.md) - 코딩 규칙 및 품질 기준

- **Reference Sections**:
  - [PRD.md#4-api-명세](./PRD.md#4-api-명세) - API 요구사항
  - [PRD.md#7-데이터-모델](./PRD.md#7-데이터-모델) - 데이터 구조 정의

### 🎨 Frontend Developer
- **Primary Documents**:
  - [Frontend.md](./Frontend.md) - 프론트엔드 구현 전체 가이드
  - [Frontend.md#4-component-specifications](./Frontend.md#4-component-specifications) - 컴포넌트 명세
  - [Frontend.md#5-state-management](./Frontend.md#5-state-management) - 상태 관리
  - [Frontend.md#17-msa-integration-strategy](./Frontend.md#17-msa-integration-strategy) - MSA 통합
  - [design-system.md](../design-system.md) - UI/UX 가이드라인

- **Reference Sections**:
  - [PRD.md#3-기능-요구사항](./PRD.md#3-기능-요구사항) - 기능 요구사항
  - [Backend.md#graphql-api](./Backend.md#4-api-implementation) - API 인터페이스

### 🎨 Designer
- **Primary Documents**:
  - [design-system.md](../design-system.md) - 디자인 시스템 전체 가이드
  - [design-system.md#color-system](../design-system.md#2-color-system) - 색상 시스템
  - [design-system.md#typography](../design-system.md#3-typography) - 타이포그래피
  - [design-system.md#component-patterns](../design-system.md#5-component-patterns) - 컴포넌트 패턴

- **Reference Sections**:
  - [Frontend.md#3-user-interface-design](./Frontend.md#3-user-interface-design) - UI 구현 가이드
  - [PRD.md#2-사용자-정의](./PRD.md#2-사용자-정의) - 사용자 경험 요구사항

### 🔧 DevOps Engineer
- **Primary Documents**:
  - [Backend.md#12-deployment-considerations](./Backend.md#12-deployment-considerations) - 배포 전략
  - [Backend.md#16-7-container-deployment](./Backend.md#16-7-container-deployment) - 컨테이너화
  - [Backend.md#16-8-health-checks](./Backend.md#16-8-health-checks) - 헬스체크 구현
  - [Frontend.md#13-deployment](./Frontend.md#13-deployment) - 프론트엔드 배포

- **Reference Sections**:
  - [PRD.md#5-비기능-요구사항](./PRD.md#5-비기능-요구사항) - 성능/가용성 요구사항
  - [PRD.md#11-마이크로서비스-아키텍처](./PRD.md#11-마이크로서비스-아키텍처) - MSA 인프라

### 🚀 CI/CD Engineer
- **Primary Documents**:
  - [Claude.rules.md#pre-commit-hooks](./Claude.rules.md#pre-commit-hooks) - Git Hooks 설정
  - [Claude.rules.md#sonarqube-quality-gates](./Claude.rules.md#sonarqube-quality-gates) - 품질 게이트
  - [Backend.md#10-testing-strategy](./Backend.md#10-testing-strategy) - 백엔드 테스트
  - [Frontend.md#10-testing-strategy](./Frontend.md#10-testing-strategy) - 프론트엔드 테스트

## 📊 Sprint Planning Resources

### Week 1-2: Foundation Sprint
**Goal**: 기초 인프라 및 개발 환경 구축

#### Backend Team
- [ ] PostgreSQL 15+ 설치 및 스키마 생성 ([Backend.md#3-database-design](./Backend.md#3-database-design))
- [ ] Go 프로젝트 구조 설정 ([Backend.md#4-1-project-structure](./Backend.md#4-1-project-structure))
- [ ] 기본 CRUD API 구현 ([Backend.md#4-3-repository-interface](./Backend.md#4-3-repository-interface))
- [ ] 단위 테스트 환경 구축 ([Backend.md#10-testing-strategy](./Backend.md#10-testing-strategy))

#### Frontend Team
- [ ] React + TypeScript 프로젝트 초기화 ([Frontend.md#2-1-technology-stack](./Frontend.md#2-1-technology-stack))
- [ ] Blueprint.js 5.x 설정 ([Frontend.md#11-theme-and-styling](./Frontend.md#11-theme-and-styling))
- [ ] 디자인 시스템 토큰 구현 ([design-system.md#design-tokens](../design-system.md#design-tokens))
- [ ] 기본 레이아웃 구현 ([Frontend.md#3-1-layout-structure](./Frontend.md#3-1-layout-structure))

#### DevOps Team
- [ ] Docker 환경 구성 ([Backend.md#dockerfile](./Backend.md#dockerfile))
- [ ] Kubernetes 매니페스트 작성 ([Backend.md#kubernetes-deployment](./Backend.md#kubernetes-deployment))
- [ ] CI/CD 파이프라인 초기 설정
- [ ] 개발 환경 프로비저닝

### Week 3-4: Core Features Sprint
**Goal**: 핵심 기능 구현

#### Backend Team
- [ ] GraphQL 스키마 구현 ([Backend.md#4-5-graphql-implementation](./Backend.md#4-5-graphql-implementation))
- [ ] 서비스 레이어 구현 ([Backend.md#4-4-service-layer](./Backend.md#4-4-service-layer))
- [ ] Redis 캐싱 구현 ([Backend.md#6-caching-strategy](./Backend.md#6-caching-strategy))
- [ ] 통합 테스트 작성

#### Frontend Team
- [ ] Object Type 관리 UI ([Frontend.md#4-1-object-type-components](./Frontend.md#4-1-object-type-components))
- [ ] Property Editor 구현 ([Frontend.md#4-1-2-propertyeditor](./Frontend.md#4-1-2-propertyeditor))
- [ ] 상태 관리 구현 ([Frontend.md#5-state-management](./Frontend.md#5-state-management))
- [ ] API 통합 ([Frontend.md#6-api-integration](./Frontend.md#6-api-integration))

#### QA Team
- [ ] 테스트 시나리오 작성
- [ ] E2E 테스트 자동화 ([Frontend.md#10-3-e2e-tests](./Frontend.md#10-3-e2e-tests))
- [ ] 성능 테스트 계획 수립

### Week 5-6: Production Ready Sprint
**Goal**: 프로덕션 준비 및 MSA 통합

#### Backend Team
- [ ] gRPC 서비스 구현 ([Backend.md#16-2-grpc-service-definition](./Backend.md#16-2-grpc-service-definition))
- [ ] Kafka 이벤트 시스템 ([Backend.md#5-event-system](./Backend.md#5-event-system))
- [ ] 분산 추적 구현 ([Backend.md#16-6-distributed-tracing](./Backend.md#16-6-distributed-tracing))
- [ ] 성능 최적화 ([Backend.md#13-performance-optimization](./Backend.md#13-performance-optimization))

#### Frontend Team
- [ ] MSA 통합 구현 ([Frontend.md#17-msa-integration-strategy](./Frontend.md#17-msa-integration-strategy))
- [ ] 실시간 업데이트 ([Frontend.md#17-4-real-time-updates](./Frontend.md#17-4-real-time-updates))
- [ ] 성능 최적화 ([Frontend.md#7-performance-optimization](./Frontend.md#7-performance-optimization))
- [ ] 접근성 개선 ([Frontend.md#9-accessibility](./Frontend.md#9-accessibility))

#### DevOps Team
- [ ] 프로덕션 환경 구성
- [ ] 모니터링 시스템 구축 ([Backend.md#8-monitoring-and-observability](./Backend.md#8-monitoring-and-observability))
- [ ] 보안 감사 및 강화 ([Backend.md#9-security-implementation](./Backend.md#9-security-implementation))
- [ ] 배포 자동화 완성

## 🔄 Cross-Team Dependencies

### API Contract
- **정의**: [Backend.md#graphql-schema](./Backend.md#graphql-schema)
- **구현**: Backend Team (Week 3)
- **소비**: Frontend Team (Week 3-4)
- **문서**: GraphQL Playground 자동 생성

### Design System
- **정의**: [design-system.md](../design-system.md)
- **구현**: Frontend Team + Designer (Week 1-2)
- **사용**: 모든 UI 컴포넌트
- **검증**: Designer 리뷰 필수

### Database Schema
- **정의**: [Backend.md#3-database-design](./Backend.md#3-database-design)
- **구현**: Backend Team (Week 1)
- **영향**: API 설계, 데이터 모델
- **마이그레이션**: DevOps 협력

### MSA Integration Points
- **정의**: [PRD.md#11-마이크로서비스-아키텍처](./PRD.md#11-마이크로서비스-아키텍처)
- **구현**: 모든 팀 (Week 5-6)
- **조율**: CTO/아키텍트 주도

## 📐 Technical Standards

### Code Quality
- **Backend**: [Claude.rules.md#go-standards](./Claude.rules.md)
- **Frontend**: [Claude.rules.md#typescript-ultra-strict](./Claude.rules.md#typescript-ultra-strict)
- **공통**: SonarQube Quality Gates 준수

### Git Workflow
```bash
main
  └── develop
       ├── feature/oms-XXX-description
       ├── bugfix/oms-XXX-description
       └── hotfix/oms-XXX-description
```

### Commit Convention
```
type(scope): subject

body

footer
```
- Types: feat, fix, docs, style, refactor, test, chore
- Scope: backend, frontend, devops, design

### API Versioning
- GraphQL: Schema 버전 관리
- gRPC: Proto 파일 버전 관리
- REST: URL 버전 포함 (/api/v1)

## 🚨 Critical Path Items

1. **Database Schema** (Week 1) - 모든 개발의 기초
2. **API Contract** (Week 2) - Frontend 개발 차단 요소
3. **Design System** (Week 1) - UI 일관성 필수
4. **CI/CD Pipeline** (Week 2) - 지속적 배포 기반
5. **MSA Infrastructure** (Week 4) - 확장성 준비

## 📞 Communication Channels

### Daily Standup
- **시간**: 매일 오전 10:00
- **참석**: 모든 팀
- **형식**: 어제/오늘/블로커

### Sprint Planning
- **주기**: 2주
- **참석**: PM, Tech Leads
- **산출물**: Sprint Backlog

### Technical Sync
- **주기**: 주 2회 (화/목)
- **참석**: Backend/Frontend/DevOps Leads
- **목적**: 기술적 이슈 해결

### Design Review
- **주기**: 주 1회 (금)
- **참석**: Designer, Frontend Lead, PM
- **목적**: UI/UX 검토 및 승인

## 🔍 Document Version Control

| Document | Version | Last Updated | Owner |
|----------|---------|--------------|-------|
| PRD.md | 1.1.0 | 2024-01-10 | PM |
| Backend.md | 1.2.0 | 2024-01-10 | Backend Lead |
| Frontend.md | 1.2.0 | 2024-01-10 | Frontend Lead |
| design-system.md | 1.0.0 | 2024-01-10 | Design Lead |
| Claude.rules.md | 1.0.0 | 2024-01-10 | CTO |

## 🎯 Success Criteria

### Technical Metrics
- [ ] API 응답 시간 < 100ms (P95)
- [ ] 프론트엔드 로딩 시간 < 2초
- [ ] 테스트 커버리지 > 80%
- [ ] SonarQube 품질 게이트 통과
- [ ] 가용성 > 99.9%

### Business Metrics
- [ ] 5개 도메인 모델링 완료
- [ ] API 문서 100% 완성
- [ ] 개발자 만족도 > 4.0/5.0
- [ ] 제로 크리티컬 버그

---

> **Note**: 이 인덱스는 살아있는 문서입니다. 프로젝트 진행에 따라 지속적으로 업데이트하세요.