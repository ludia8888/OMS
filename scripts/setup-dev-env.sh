#!/bin/bash
# OMS Development Environment Setup Script
# This script sets up the complete development environment for OMS

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Check OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    PACKAGE_MANAGER="brew"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if command -v apt-get &> /dev/null; then
        PACKAGE_MANAGER="apt"
    elif command -v yum &> /dev/null; then
        PACKAGE_MANAGER="yum"
    fi
else
    print_error "Unsupported operating system"
    exit 1
fi

print_info "Starting OMS development environment setup..."

# 1. Install system dependencies
print_info "Installing system dependencies..."

if [[ "$PACKAGE_MANAGER" == "brew" ]]; then
    # macOS
    brew update
    brew install git docker docker-compose kubectl helm go@1.21 node@18 postgresql@15 redis
    brew install --cask docker
elif [[ "$PACKAGE_MANAGER" == "apt" ]]; then
    # Ubuntu/Debian
    sudo apt update
    sudo apt install -y git docker.io docker-compose kubectl golang-1.21 nodejs npm postgresql-15 redis-server
fi

print_success "System dependencies installed"

# 2. Install Node.js tools
print_info "Installing Node.js tools..."
npm install -g pnpm
print_success "Node.js tools installed"

# 3. Create project structure
print_info "Creating project structure..."
mkdir -p ~/workspace/oms/{backend,frontend,k8s,scripts,docs}
cd ~/workspace/oms

# 4. Initialize Git repository
if [ ! -d .git ]; then
    print_info "Initializing Git repository..."
    git init
    cat > .gitignore << 'EOF'
# Dependencies
node_modules/
vendor/
*.log

# Build outputs
dist/
build/
*.exe
*.dll
*.so
*.dylib

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# Environment
.env
.env.local
.env.*.local

# OS
.DS_Store
Thumbs.db

# Test coverage
coverage/
*.cover
*.coverage

# Docker
docker-compose.override.yml
EOF
    git add .
    git commit -m "Initial commit"
    print_success "Git repository initialized"
fi

# 5. Setup Backend
print_info "Setting up Backend..."
cd backend

# Initialize Go module
if [ ! -f go.mod ]; then
    go mod init github.com/openfoundry/oms
fi

# Create backend structure
mkdir -p cmd/server
mkdir -p internal/{config,domain,infrastructure,interfaces,pkg}
mkdir -p internal/domain/{entity,repository,service}
mkdir -p internal/infrastructure/{database,cache,messaging}
mkdir -p internal/interfaces/{graphql,rest,grpc}
mkdir -p internal/pkg/{errors,logger,validator}

# Create main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("OMS Backend Server Starting...")
    log.Fatal("Not implemented yet")
}
EOF

# Create Makefile
cat > Makefile << 'EOF'
.PHONY: help build run test lint clean

help:
	@echo "Available commands:"
	@echo "  make build   - Build the application"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make lint    - Run linter"
	@echo "  make clean   - Clean build artifacts"

build:
	go build -o bin/oms-server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
EOF

print_success "Backend setup complete"

# 6. Setup Frontend
print_info "Setting up Frontend..."
cd ../frontend

# Create package.json
cat > package.json << 'EOF'
{
  "name": "oms-frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "jest",
    "lint": "eslint . --ext .ts,.tsx",
    "type-check": "tsc --noEmit",
    "format": "prettier --write ."
  },
  "dependencies": {
    "@apollo/client": "^3.8.0",
    "@blueprintjs/core": "^5.0.0",
    "@blueprintjs/icons": "^5.0.0",
    "@blueprintjs/select": "^5.0.0",
    "@blueprintjs/table": "^5.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.0.0",
    "zustand": "^4.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "eslint": "^8.0.0",
    "jest": "^29.0.0",
    "prettier": "^3.0.0",
    "typescript": "^5.0.0",
    "vite": "^4.0.0"
  }
}
EOF

# Create tsconfig.json
cat > tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitOverride": true,
    "exactOptionalPropertyTypes": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
EOF

# Create vite.config.ts
cat > vite.config.ts << 'EOF'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
EOF

# Create src structure
mkdir -p src/{app,features,shared,design-system,assets}
mkdir -p src/features/{object-types,link-types,search}
mkdir -p src/shared/{components,hooks,services,utils}

# Install dependencies
print_info "Installing frontend dependencies..."
pnpm install

print_success "Frontend setup complete"

# 7. Setup Docker Compose
print_info "Creating Docker Compose configuration..."
cd ..

cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: oms
      POSTGRES_PASSWORD: oms123
      POSTGRES_DB: oms
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

volumes:
  postgres_data:
  redis_data:
EOF

print_success "Docker Compose configuration created"

# 8. Create environment files
print_info "Creating environment configuration..."

cat > backend/.env.example << 'EOF'
# Server
SERVER_PORT=8080
SERVER_MODE=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=oms
DB_USER=oms
DB_PASSWORD=oms123
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=oms-events

# Security
JWT_SECRET=your-secret-key-here
API_KEY_HEADER=X-API-Key

# Monitoring
METRICS_PATH=/metrics
TRACE_ENDPOINT=http://localhost:14268/api/traces
EOF

cat > frontend/.env.example << 'EOF'
VITE_API_URL=http://localhost:8080
VITE_GRAPHQL_ENDPOINT=http://localhost:8080/graphql
VITE_WS_ENDPOINT=ws://localhost:8080/graphql
EOF

# Copy example files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

print_success "Environment configuration created"

# 9. Setup Git hooks
print_info "Setting up Git hooks..."

mkdir -p .git/hooks

cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for OMS

# Backend checks
echo "Running backend checks..."
cd backend
go fmt ./...
go vet ./...
golangci-lint run --fast

# Frontend checks
echo "Running frontend checks..."
cd ../frontend
pnpm lint
pnpm type-check

echo "Pre-commit checks passed!"
EOF

chmod +x .git/hooks/pre-commit

print_success "Git hooks configured"

# 10. Create helpful scripts
print_info "Creating helper scripts..."

cat > scripts/start-services.sh << 'EOF'
#!/bin/bash
# Start all required services

echo "Starting Docker services..."
docker-compose up -d

echo "Waiting for services to be ready..."
sleep 10

echo "Services status:"
docker-compose ps

echo "Services are ready!"
echo "PostgreSQL: localhost:5432"
echo "Redis: localhost:6379"
echo "Kafka: localhost:9092"
EOF

cat > scripts/stop-services.sh << 'EOF'
#!/bin/bash
# Stop all services

echo "Stopping Docker services..."
docker-compose down

echo "Services stopped."
EOF

chmod +x scripts/*.sh

print_success "Helper scripts created"

# 11. Final setup
print_info "Finalizing setup..."

# Create README
cat > README.md << 'EOF'
# OMS - Ontology Metadata Service

## Quick Start

1. Start services:
   ```bash
   ./scripts/start-services.sh
   ```

2. Run backend:
   ```bash
   cd backend
   make run
   ```

3. Run frontend:
   ```bash
   cd frontend
   pnpm dev
   ```

## Documentation

- [Master Index](Claude.docs/INDEX.md)
- [PRD](Claude.docs/PRD.md)
- [Backend Spec](Claude.docs/Backend.md)
- [Frontend Spec](Claude.docs/Frontend.md)
- [Design System](design-system.md)

## Development

See [Team Quick Start](Claude.docs/TEAM_QUICKSTART.md) for role-specific guides.
EOF

# Summary
echo ""
print_success "🎉 OMS Development Environment Setup Complete! 🎉"
echo ""
echo "Next steps:"
echo "1. cd ~/workspace/oms"
echo "2. ./scripts/start-services.sh  # Start Docker services"
echo "3. cd backend && make run       # Start backend"
echo "4. cd frontend && pnpm dev      # Start frontend"
echo ""
echo "Documentation: ~/workspace/oms/Claude.docs/"
echo ""
print_info "Happy coding! 🚀"