# OMS Documentation Consistency Fixes

> **Created**: 2024-01-10  
> **Purpose**: 문서 간 불일치 해결 및 정합성 개선 가이드

## 📋 Action Items by Priority

### 🔴 Critical (Week 1)

#### 1. Technology Stack Alignment
**Issue**: PRD와 Backend 문서 간 기술 스택 불일치

**PRD.md Updates Required**:
```markdown
# Line 332-336 수정
기술 스택:
  - 언어: Go 1.21+
  - 데이터베이스: PostgreSQL 15+
  - 캐시: Redis 7+
  - API: GraphQL (gqlgen for Go)
  - Message Queue: Apache Kafka
```

**Rationale**: Backend.md에서 이미 Go로 구현이 확정되었으므로 PRD를 업데이트

#### 2. Data Model Synchronization
**Issue**: object_types 테이블 구조 불일치

**PRD.md Updates Required**:
```sql
-- Line 348-363 수정
CREATE TABLE object_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(64) UNIQUE NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(64),
  tags JSONB DEFAULT '[]'::jsonb,
  properties JSONB NOT NULL DEFAULT '[]'::jsonb,
  base_datasets JSONB DEFAULT '[]'::jsonb,  -- 추가
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  version INTEGER NOT NULL DEFAULT 1,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_by VARCHAR(255) NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_by VARCHAR(255) NOT NULL
);
```

#### 3. Performance Requirements Alignment
**Issue**: API 응답 시간 목표 불일치

**Action**: 팀 미팅을 통해 현실적인 목표 설정
- 단일 조회: < 100ms (P95) - Backend 구현 기준
- 목록 조회: < 150ms (P95) - 페이지네이션 고려
- 생성/수정: < 200ms (P95) - 유지

### 🟡 Important (Week 2)

#### 4. API Specification Updates
**Issue**: GraphQL 스키마 불일치

**Create Common Schema File**:
```bash
/Users/sihyun/Desktop/OMS/schema/
├── common.graphql      # 공통 타입 정의
├── oms.graphql        # OMS 서비스 스키마
└── federation.graphql # Federation 스키마
```

#### 5. Security Implementation
**Issue**: API 키 인증 누락

**Backend.md Addition** (After line 1310):
```go
// API Key Authentication
type APIKeyAuth struct {
    validKeys map[string]string
}

func NewAPIKeyAuth(keys map[string]string) *APIKeyAuth {
    return &APIKeyAuth{validKeys: keys}
}

func (a *APIKeyAuth) Authenticate(key string) (bool, string) {
    if clientID, exists := a.validKeys[key]; exists {
        return true, clientID
    }
    return false, ""
}
```

#### 6. Caching Strategy Documentation
**Issue**: 캐시 TTL 전략 불명확

**Create Cache Strategy Document**:
```markdown
# /Users/sihyun/Desktop/OMS/Claude.docs/CACHE_STRATEGY.md

## Cache TTL Guidelines
- Object Type Metadata: 30분 (자주 변경되지 않음)
- Object Type List: 5분 (새로운 타입 추가 고려)
- Search Results: 1분 (실시간성 중요)
- User Session: 24시간
```

### 🟢 Nice to Have (Week 3+)

#### 7. Testing Coverage Goals
**All Documents**: 테스트 커버리지 목표 통일
- Unit Tests: 80% 이상
- Integration Tests: 70% 이상
- E2E Tests: Critical Path 100%

#### 8. Frontend State Management
**Frontend.md Clarification**:
```markdown
## State Management Strategy
- **Zustand**: Local UI state (modals, forms, filters)
- **Apollo Cache**: Server state (API responses)
- **React Query**: Non-GraphQL API calls
```

## 🔗 Cross-Reference Matrix

| Feature | PRD | Backend | Frontend | Design System |
|---------|-----|---------|----------|---------------|
| Object Types | §3.1 | §4.2 | §4.1 | §5.1 |
| Properties | §3.1.2 | §4.2 | §4.1.2 | §5.2 |
| Relationships | §3.1.3 | §4.2 | §4.2 | §5.3 |
| Search | §3.4 | §4.6 | §4.3 | - |
| API | §4 | §4.5 | §6 | - |
| MSA | §11 | §16 | §17 | - |

## 📝 Documentation Standards

### Version Format
```
Major.Minor.Patch
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes, clarifications
```

### Update Process
1. Create branch: `docs/oms-XXX-description`
2. Update version in document header
3. Add entry to CHANGELOG.md
4. Get review from document owner
5. Merge to main

### Document Owners
- **PRD.md**: Product Manager
- **Backend.md**: Backend Tech Lead
- **Frontend.md**: Frontend Tech Lead
- **design-system.md**: Design Lead
- **Claude.rules.md**: CTO/Engineering Manager

## 🚀 Implementation Plan

### Week 1: Critical Fixes
- [ ] Update PRD technology stack
- [ ] Synchronize data models
- [ ] Align performance requirements
- [ ] Create master schema file

### Week 2: Important Updates
- [ ] Implement API key authentication
- [ ] Document caching strategy
- [ ] Update GraphQL schemas
- [ ] Add cross-references

### Week 3: Enhancements
- [ ] Unify testing goals
- [ ] Clarify state management
- [ ] Add glossary
- [ ] Complete cross-reference matrix

## 📊 Tracking

### Consistency Metrics
- [ ] Technology Stack: ⚠️ In Progress
- [ ] Data Models: ⚠️ In Progress
- [ ] API Contracts: ⚠️ In Progress
- [ ] Performance Goals: ⚠️ In Progress
- [ ] Security Implementation: ❌ Pending
- [ ] Testing Strategy: ❌ Pending
- [ ] Deployment Configuration: ✅ Aligned

### Review Schedule
- Weekly consistency review: Every Friday
- Full document audit: End of each sprint
- Stakeholder sign-off: Before major releases

---

> **Note**: This document should be reviewed and updated weekly until all consistency issues are resolved.