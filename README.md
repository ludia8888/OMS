# Ontology Metadata Service (OMS)

[![Go Report Card](https://goreportcard.com/badge/github.com/openfoundry/oms)](https://goreportcard.com/report/github.com/openfoundry/oms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📖 개요

Ontology Metadata Service (OMS)는 OpenFoundry 플랫폼의 핵심 메타데이터 관리 서비스입니다. 조직의 비즈니스 개체와 관계를 정의하는 스키마 레지스트리로, Palantir Foundry의 Ontology 철학을 계승하여 데이터를 현실 세계의 객체로 추상화하는 구조적 정의를 제공합니다.

### 주요 기능

- 🏗️ **스키마 레지스트리**: 조직 전체의 통일된 객체 정의
- 💼 **비즈니스 친화적**: 기술 용어가 아닌 비즈니스 용어로 모델링
- 🔄 **확장 가능한 설계**: 향후 서비스 통합을 위한 견고한 기반
- 📊 **메타데이터 관리**: 객체의 속성과 관계 정의
- 🔍 **검색 기능**: 강력한 검색 및 필터링 기능
- 📝 **버전 관리**: 모든 변경사항 추적 및 관리

## 🚀 시작하기

### 필수 요구사항

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

### 빠른 시작

1. 저장소 클론
```bash
git clone https://github.com/openfoundry/oms.git
cd oms
```

2. 백엔드 설정
```bash
cd backend
go mod download
docker-compose up -d  # PostgreSQL & Redis 시작
go run cmd/server/main.go
```

3. 프론트엔드 설정
```bash
cd frontend
pnpm install
pnpm dev
```

자세한 설정 가이드는 [TEAM_QUICKSTART.md](Claude.docs/TEAM_QUICKSTART.md)를 참조하세요.

## 📚 문서

- [제품 요구사항 문서 (PRD)](Claude.docs/PRD.md)
- [백엔드 개발 가이드](Claude.docs/Backend.md)
- [프론트엔드 개발 가이드](Claude.docs/Frontend.md)
- [디자인 시스템](Claude.docs/design-system.md)
- [팀 퀵스타트 가이드](Claude.docs/TEAM_QUICKSTART.md)

## 🏗️ 아키텍처

OMS는 다음과 같은 주요 컴포넌트로 구성됩니다:

- **Backend**: Go 기반의 GraphQL/REST API 서버
- **Frontend**: React + TypeScript 기반의 웹 인터페이스
- **Database**: PostgreSQL을 사용한 메타데이터 저장
- **Cache**: Redis를 사용한 성능 최적화

## 🤝 기여하기

1. 이슈 생성 또는 기존 이슈 확인
2. 브랜치 생성 (`feature/`, `fix/`, `docs/` 등)
3. 변경사항 커밋
4. Pull Request 생성

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 👥 팀

- Backend Team: [@backend-team](https://github.com/orgs/openfoundry/teams/backend-team)
- Frontend Team: [@frontend-team](https://github.com/orgs/openfoundry/teams/frontend-team)
- DevOps Team: [@devops-team](https://github.com/orgs/openfoundry/teams/devops-team)

## 📞 지원

- 이슈 트래커: [GitHub Issues](https://github.com/openfoundry/oms/issues)
- 이메일: support@openfoundry.io