# OMS Team Quick Start Guides

> **Purpose**: 각 팀이 즉시 작업을 시작할 수 있는 실용적인 가이드

## 🚀 Backend Developer Quick Start

### Day 1: Environment Setup
```bash
# 1. Clone repository
git clone https://github.com/openfoundry/oms.git
cd oms/backend

# 2. Install Go 1.21+
brew install go@1.21

# 3. Install dependencies
go mod download

# 4. Setup PostgreSQL 15
docker run -d \
  --name oms-postgres \
  -e POSTGRES_PASSWORD=oms123 \
  -e POSTGRES_DB=oms \
  -p 5432:5432 \
  postgres:15-alpine

# 5. Run migrations
migrate -path migrations -database "postgresql://localhost:5432/oms?sslmode=disable" up

# 6. Start Redis
docker run -d \
  --name oms-redis \
  -p 6379:6379 \
  redis:7-alpine
```

### Day 2: First API
```go
// 1. Check Backend.md Section 4.2 for domain models
// 2. Implement your first service in internal/domain/service/
// 3. Add GraphQL resolver in internal/interfaces/graphql/
// 4. Write unit tests
// 5. Run: go test ./...
```

### Key Files to Review
- `/internal/domain/entity/object_type.go` - Core domain model
- `/internal/interfaces/graphql/schema.graphql` - API contract
- `/internal/infrastructure/database/migrations/` - DB schema

### Common Commands
```bash
# Run server
go run cmd/server/main.go

# Generate GraphQL code
go run github.com/99designs/gqlgen generate

# Run tests with coverage
go test -cover ./...

# Lint code
golangci-lint run
```

---

## 🎨 Frontend Developer Quick Start

### Day 1: Environment Setup
```bash
# 1. Clone repository
git clone https://github.com/openfoundry/oms.git
cd oms/frontend

# 2. Install Node.js 18+ and pnpm
brew install node@18
npm install -g pnpm

# 3. Install dependencies
pnpm install

# 4. Setup environment
cp .env.example .env.local
# Edit .env.local with your API endpoints

# 5. Start development server
pnpm dev
```

### Day 2: First Component
```tsx
// 1. Check design-system.md for UI patterns
// 2. Review Frontend.md Section 4.1 for component examples
// 3. Create component in src/features/object-types/components/
// 4. Use Blueprint.js components
// 5. Write tests in __tests__/
```

### Key Files to Review
- `/src/design-system/tokens/` - Design tokens
- `/src/features/object-types/components/ObjectTypeCard.tsx` - Example component
- `/src/shared/hooks/useObjectTypes.ts` - Data fetching pattern

### Common Commands
```bash
# Start dev server
pnpm dev

# Run tests
pnpm test

# Type check
pnpm type-check

# Lint and format
pnpm lint
pnpm format

# Build production
pnpm build
```

---

## 🎨 Designer Quick Start

### Day 1: Design System Setup
```bash
# 1. Install Figma/Sketch plugins
# - Blueprint.js UI Kit
# - Design Tokens plugin

# 2. Review design-system.md completely

# 3. Access color palette
# Primary: #2D72D2
# Success: #238551
# Warning: #C87619
# Danger: #CD4246

# 4. Setup local preview
cd oms/frontend
pnpm install
pnpm storybook
```

### Day 2: First Design
1. Use 8px grid system
2. Follow typography scale (11px to 36px)
3. Use Blueprint.js components as base
4. Check Frontend.md Section 3 for layouts
5. Export assets in @1x, @2x, @3x

### Key Resources
- Blueprint.js Docs: https://blueprintjs.com/docs/
- Color system: `design-system.md#2-color-system`
- Component patterns: `design-system.md#5-component-patterns`
- Existing components: `frontend/src/design-system/components/`

---

## 🔧 DevOps Engineer Quick Start

### Day 1: Infrastructure Setup
```bash
# 1. Install tools
brew install kubectl helm terraform
brew install --cask docker

# 2. Setup local Kubernetes
# Option A: Docker Desktop with K8s enabled
# Option B: minikube start

# 3. Create namespaces
kubectl create namespace openfoundry
kubectl create namespace monitoring

# 4. Install Istio (Service Mesh)
istioctl install --set profile=demo -y
kubectl label namespace openfoundry istio-injection=enabled
```

### Day 2: First Deployment
```bash
# 1. Review Backend.md Section 16.9 for K8s configs
# 2. Deploy PostgreSQL
helm install oms-postgres bitnami/postgresql \
  --namespace openfoundry \
  --values k8s/postgres-values.yaml

# 3. Deploy Redis
helm install oms-redis bitnami/redis \
  --namespace openfoundry \
  --values k8s/redis-values.yaml

# 4. Build and deploy OMS
docker build -t oms-backend:latest ./backend
kubectl apply -f k8s/oms-backend/
```

### Key Files to Review
- `/k8s/` - Kubernetes manifests
- `/backend/Dockerfile` - Container definition
- `Backend.md#16-8` - Health check implementation
- `Backend.md#16-9` - Deployment configs

### Monitoring Setup
```bash
# Prometheus + Grafana
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring

# Jaeger for tracing
kubectl apply -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
```

---

## 👨‍💼 Product Manager Quick Start

### Day 1: Context Understanding
1. **Read PRD.md completely** - Your bible
2. **Review INDEX.md** - Understand all documents
3. **Check CONSISTENCY_FIXES.md** - Know current issues
4. **Setup access to**:
   - Jira/Linear for sprint management
   - Figma for design reviews
   - GitHub for code reviews
   - Slack/Discord for team communication

### Day 2: Sprint Planning
```markdown
## Sprint 1 Checklist
- [ ] Review PRD Section 8 (Milestones)
- [ ] Create Jira epics for Week 1-2 goals
- [ ] Schedule kick-off meeting
- [ ] Define success criteria (PRD Section 9)
- [ ] Setup daily standup (10 AM)
- [ ] Create team channels
```

### Key Metrics to Track
- Sprint velocity
- Bug count by severity
- API response times
- Test coverage percentage
- Team satisfaction score

### Stakeholder Communication
```markdown
## Weekly Status Template
### Completed This Week
- Feature X (link to demo)
- Bug fixes (#123, #124)

### In Progress
- Feature Y (60% complete)
- Performance optimization

### Blockers
- Waiting for design approval on Z

### Next Week Plan
- Complete Feature Y
- Start Feature Z
```

---

## 🏃‍♂️ First Sprint Checklist

### All Teams - Before Sprint 1
- [ ] Read your primary documents (see INDEX.md)
- [ ] Setup development environment
- [ ] Join team Slack channels
- [ ] Attend kick-off meeting
- [ ] Review sprint goals

### Backend Team - Week 1
- [ ] Database schema implementation
- [ ] Basic CRUD APIs
- [ ] Unit test setup
- [ ] GraphQL schema draft

### Frontend Team - Week 1
- [ ] Project setup with Blueprint.js
- [ ] Design system implementation
- [ ] Basic layout components
- [ ] Storybook setup

### DevOps Team - Week 1
- [ ] Local development environment
- [ ] CI/CD pipeline skeleton
- [ ] Container configurations
- [ ] Basic monitoring

### End of Week 1 Deliverables
- [ ] Working local environment for all
- [ ] Basic API endpoints (Backend)
- [ ] Component library started (Frontend)
- [ ] CI pipeline running (DevOps)
- [ ] Sprint 2 planning complete (PM)

---

## 📚 Learning Resources

### Backend
- [Effective Go](https://golang.org/doc/effective_go.html)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Go Microservices](https://www.oreilly.com/library/view/building-microservices-with/9781786468666/)

### Frontend
- [Blueprint.js Docs](https://blueprintjs.com/docs/)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Apollo Client Docs](https://www.apollographql.com/docs/react/)

### DevOps
- [Kubernetes Patterns](https://k8spatterns.io/)
- [Istio Documentation](https://istio.io/latest/docs/)
- [12 Factor App](https://12factor.net/)

### Design
- [Palantir Design Philosophy](https://medium.com/palantir/design-at-palantir-5t634tg)
- [Enterprise UX Patterns](https://www.uxpin.com/studio/blog/enterprise-ux-design/)

---

## 🆘 Getting Help

### Slack Channels
- `#oms-general` - General discussions
- `#oms-backend` - Backend specific
- `#oms-frontend` - Frontend specific
- `#oms-devops` - Infrastructure
- `#oms-design` - Design discussions
- `#oms-help` - Quick questions

### Key Contacts
- **CTO**: @cto-handle (Architecture decisions)
- **Backend Lead**: @backend-lead (API, Database)
- **Frontend Lead**: @frontend-lead (UI, State management)
- **DevOps Lead**: @devops-lead (Infrastructure, CI/CD)
- **Design Lead**: @design-lead (UI/UX decisions)
- **PM**: @pm-handle (Requirements, Sprint planning)

### Office Hours
- **Architecture Review**: Tue 2-3 PM
- **Code Review**: Thu 3-4 PM
- **Design Review**: Fri 2-3 PM

---

> **Remember**: When in doubt, check the documentation first, ask in Slack second, and schedule a meeting last.