# Ontology Metadata Service (OMS) Backend Development Specification

## 1. Executive Summary

본 문서는 OpenFoundry의 핵심 서비스인 Ontology Metadata Service (OMS)의 백엔드 구현을 위한 상세 개발 명세서입니다. OMS는 Palantir Foundry의 Ontology 철학을 계승하여, 비즈니스 객체와 관계를 정의하는 메타데이터 관리 시스템의 백엔드를 구현합니다.

## 2. System Architecture

### 2.1 Overall Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                           │
│                    (GraphQL Federation)                      │
└─────────────────┬───────────────────────┬───────────────────┘
                  │                       │
                  ▼                       ▼
        ┌─────────────────┐     ┌─────────────────┐
        │   GraphQL API   │     │    REST API     │
        │    (Primary)    │     │  (Compatibility) │
        └────────┬────────┘     └────────┬────────┘
                 │                       │
                 └───────────┬───────────┘
                             │
                  ┌──────────▼──────────┐
                  │   Service Layer     │
                  │  (Business Logic)   │
                  └──────────┬──────────┘
                             │
                  ┌──────────▼──────────┐
                  │  Repository Layer   │
                  │   (Data Access)     │
                  └──────────┬──────────┘
                             │
         ┌───────────────────┴───────────────────┐
         │                                       │
    ┌────▼────┐                            ┌────▼────┐
    │PostgreSQL│                            │  Redis  │
    │ (Primary)│                            │ (Cache) │
    └──────────┘                            └─────────┘
                             │
                    ┌────────▼────────┐
                    │   Kafka Event   │
                    │     Stream      │
                    └─────────────────┘
```

### 2.2 Technology Stack

```yaml
Language: Go 1.21+
Framework: 
  - Gin (HTTP router)
  - gqlgen (GraphQL)
  - gRPC (internal communication)

Database:
  - PostgreSQL 15+ (primary storage)
  - Redis 7+ (caching layer)
  
Message Queue:
  - Apache Kafka (event streaming)
  
Monitoring:
  - Prometheus (metrics)
  - Jaeger (distributed tracing)
  - Zap (structured logging)
```

## 3. Database Design

### 3.1 Core Tables

#### object_types
```sql
CREATE TABLE object_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(64),
    tags JSONB DEFAULT '[]'::jsonb,
    properties JSONB NOT NULL DEFAULT '[]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    
    CONSTRAINT object_type_name_format CHECK (name ~ '^[a-zA-Z][a-zA-Z0-9_]*$'),
    CONSTRAINT object_type_name_length CHECK (char_length(name) <= 64)
);

CREATE INDEX idx_object_types_name ON object_types(name) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_category ON object_types(category) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_tags ON object_types USING GIN (tags) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_created_at ON object_types(created_at DESC) WHERE is_deleted = FALSE;

-- Full-text search index
CREATE INDEX idx_object_types_search ON object_types 
USING GIN (to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, ''))) 
WHERE is_deleted = FALSE;
```

#### object_type_versions
```sql
CREATE TABLE object_type_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_type_id UUID NOT NULL REFERENCES object_types(id),
    version INTEGER NOT NULL,
    snapshot JSONB NOT NULL,
    change_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    
    UNIQUE(object_type_id, version)
);

CREATE INDEX idx_object_type_versions_object_type_id ON object_type_versions(object_type_id);
```

#### link_types
```sql
CREATE TABLE link_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    source_object_type_id UUID NOT NULL REFERENCES object_types(id),
    target_object_type_id UUID NOT NULL REFERENCES object_types(id),
    cardinality VARCHAR(32) NOT NULL CHECK (cardinality IN ('ONE_TO_ONE', 'ONE_TO_MANY', 'MANY_TO_MANY')),
    description TEXT,
    properties JSONB DEFAULT '[]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    
    CONSTRAINT link_type_name_format CHECK (name ~ '^[a-zA-Z][a-zA-Z0-9_]*$')
);

CREATE INDEX idx_link_types_source ON link_types(source_object_type_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_link_types_target ON link_types(target_object_type_id) WHERE is_deleted = FALSE;
```

#### audit_logs
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    actor VARCHAR(255) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    old_value JSONB,
    new_value JSONB,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs(actor);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

### 3.2 Property Schema Structure

```json
{
  "id": "uuid",
  "name": "customerName",
  "displayName": "Customer Name",
  "dataType": "STRING",
  "required": true,
  "unique": false,
  "indexed": true,
  "defaultValue": null,
  "description": "The full name of the customer",
  "validators": [
    {
      "type": "maxLength",
      "value": 255
    },
    {
      "type": "pattern",
      "value": "^[a-zA-Z\\s]+$"
    }
  ],
  "metadata": {
    "uiHints": {
      "width": "medium",
      "placeholder": "Enter customer name"
    }
  }
}
```

## 4. API Implementation

### 4.1 Project Structure

```
oms-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── object_type.go
│   │   │   ├── link_type.go
│   │   │   └── property.go
│   │   ├── repository/
│   │   │   ├── object_type_repository.go
│   │   │   └── link_type_repository.go
│   │   └── service/
│   │       ├── object_type_service.go
│   │       └── link_type_service.go
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go
│   │   │   └── migrations/
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── messaging/
│   │   │   └── kafka.go
│   │   └── repository/
│   │       ├── postgres_object_type_repository.go
│   │       └── postgres_link_type_repository.go
│   ├── interfaces/
│   │   ├── graphql/
│   │   │   ├── schema.graphql
│   │   │   ├── resolver.go
│   │   │   └── generated.go
│   │   └── rest/
│   │       ├── handler/
│   │       └── middleware/
│   └── pkg/
│       ├── errors/
│       ├── logger/
│       └── validator/
├── pkg/
│   └── proto/
│       └── oms.proto
├── go.mod
└── go.sum
```

### 4.2 Domain Models

#### Object Type Entity
```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type ObjectType struct {
    ID           uuid.UUID              `json:"id"`
    Name         string                 `json:"name"`
    DisplayName  string                 `json:"displayName"`
    Description  *string                `json:"description,omitempty"`
    Category     *string                `json:"category,omitempty"`
    Tags         []string               `json:"tags"`
    Properties   []Property             `json:"properties"`
    BaseDatasets []DatasetReference     `json:"baseDatasets,omitempty"`
    Metadata     map[string]interface{} `json:"metadata"`
    Version      int                    `json:"version"`
    IsDeleted    bool                   `json:"-"`
    CreatedAt    time.Time              `json:"createdAt"`
    CreatedBy    string                 `json:"createdBy"`
    UpdatedAt    time.Time              `json:"updatedAt"`
    UpdatedBy    string                 `json:"updatedBy"`
}

type DatasetReference struct {
    DatasetRID string `json:"datasetRid"`
    Name       string `json:"name"`
}

type Property struct {
    ID           uuid.UUID              `json:"id"`
    Name         string                 `json:"name"`
    DisplayName  string                 `json:"displayName"`
    DataType     DataType               `json:"dataType"`
    Required     bool                   `json:"required"`
    Unique       bool                   `json:"unique"`
    Indexed      bool                   `json:"indexed"`
    DefaultValue interface{}            `json:"defaultValue,omitempty"`
    Description  *string                `json:"description,omitempty"`
    Validators   []Validator            `json:"validators,omitempty"`
    Metadata     map[string]interface{} `json:"metadata"`
}

type DataType string

const (
    DataTypeString    DataType = "STRING"
    DataTypeNumber    DataType = "NUMBER"
    DataTypeBoolean   DataType = "BOOLEAN"
    DataTypeDate      DataType = "DATE"
    DataTypeDateTime  DataType = "DATETIME"
    DataTypeArray     DataType = "ARRAY"
    DataTypeObject    DataType = "OBJECT"
    DataTypeReference DataType = "REFERENCE"
)

type Validator struct {
    Type  ValidatorType `json:"type"`
    Value interface{}   `json:"value"`
}

type ValidatorType string

const (
    ValidatorMinLength ValidatorType = "minLength"
    ValidatorMaxLength ValidatorType = "maxLength"
    ValidatorPattern   ValidatorType = "pattern"
    ValidatorMin       ValidatorType = "min"
    ValidatorMax       ValidatorType = "max"
)
```

### 4.3 Repository Interface

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/openfoundry/oms/internal/domain/entity"
)

type ObjectTypeRepository interface {
    // Basic CRUD
    Create(ctx context.Context, objectType *entity.ObjectType) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error)
    GetByName(ctx context.Context, name string) (*entity.ObjectType, error)
    Update(ctx context.Context, objectType *entity.ObjectType) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // Query operations
    List(ctx context.Context, filter ObjectTypeFilter) ([]*entity.ObjectType, error)
    Count(ctx context.Context, filter ObjectTypeFilter) (int64, error)
    Search(ctx context.Context, query string, limit int) ([]*entity.ObjectType, error)
    
    // Version management
    GetVersion(ctx context.Context, id uuid.UUID, version int) (*entity.ObjectType, error)
    ListVersions(ctx context.Context, id uuid.UUID) ([]*ObjectTypeVersion, error)
    CompareVersions(ctx context.Context, id uuid.UUID, v1, v2 int) (*VersionDiff, error)
}

type VersionDiff struct {
    ObjectTypeID uuid.UUID      `json:"objectTypeId"`
    Version1     int            `json:"version1"`
    Version2     int            `json:"version2"`
    Changes      []FieldChange  `json:"changes"`
}

type FieldChange struct {
    Field    string      `json:"field"`
    OldValue interface{} `json:"oldValue"`
    NewValue interface{} `json:"newValue"`
    Type     string      `json:"type"` // "added", "removed", "modified"
}

type ObjectTypeFilter struct {
    Category   *string
    Tags       []string
    PageSize   int
    PageCursor string // Cursor-based pagination: encoded timestamp + id
    SortBy     string
    SortOrder  string
}

type PageCursor struct {
    Timestamp time.Time
    ID        uuid.UUID
}

// EncodeCursor encodes pagination cursor
func EncodeCursor(timestamp time.Time, id uuid.UUID) string {
    data := fmt.Sprintf("%d:%s", timestamp.Unix(), id.String())
    return base64.StdEncoding.EncodeToString([]byte(data))
}

// DecodeCursor decodes pagination cursor
func DecodeCursor(cursor string) (*PageCursor, error) {
    data, err := base64.StdEncoding.DecodeString(cursor)
    if err != nil {
        return nil, err
    }
    
    parts := strings.Split(string(data), ":")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid cursor format")
    }
    
    timestamp, err := strconv.ParseInt(parts[0], 10, 64)
    if err != nil {
        return nil, err
    }
    
    id, err := uuid.Parse(parts[1])
    if err != nil {
        return nil, err
    }
    
    return &PageCursor{
        Timestamp: time.Unix(timestamp, 0),
        ID:        id,
    }, nil
}
```

### 4.4 Service Layer

```go
package service

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/openfoundry/oms/internal/domain/entity"
    "github.com/openfoundry/oms/internal/domain/repository"
    "github.com/openfoundry/oms/internal/infrastructure/messaging"
)

type ObjectTypeService struct {
    repo      repository.ObjectTypeRepository
    cache     CacheService
    publisher messaging.EventPublisher
}

func NewObjectTypeService(
    repo repository.ObjectTypeRepository,
    cache CacheService,
    publisher messaging.EventPublisher,
) *ObjectTypeService {
    return &ObjectTypeService{
        repo:      repo,
        cache:     cache,
        publisher: publisher,
    }
}

// CompareVersions compares two versions of an object type
func (s *ObjectTypeService) CompareVersions(ctx context.Context, id uuid.UUID, v1, v2 int) (*VersionDiff, error) {
    // Get both versions
    version1, err := s.repo.GetVersion(ctx, id, v1)
    if err != nil {
        return nil, fmt.Errorf("failed to get version %d: %w", v1, err)
    }
    
    version2, err := s.repo.GetVersion(ctx, id, v2)
    if err != nil {
        return nil, fmt.Errorf("failed to get version %d: %w", v2, err)
    }
    
    // Compare versions
    diff := &VersionDiff{
        ObjectTypeID: id,
        Version1:     v1,
        Version2:     v2,
        Changes:      []FieldChange{},
    }
    
    // Compare fields
    if version1.Name != version2.Name {
        diff.Changes = append(diff.Changes, FieldChange{
            Field:    "name",
            OldValue: version1.Name,
            NewValue: version2.Name,
            Type:     "modified",
        })
    }
    
    if version1.DisplayName != version2.DisplayName {
        diff.Changes = append(diff.Changes, FieldChange{
            Field:    "displayName",
            OldValue: version1.DisplayName,
            NewValue: version2.DisplayName,
            Type:     "modified",
        })
    }
    
    // Compare properties
    propChanges := s.compareProperties(version1.Properties, version2.Properties)
    diff.Changes = append(diff.Changes, propChanges...)
    
    return diff, nil
}

func (s *ObjectTypeService) CreateObjectType(
    ctx context.Context,
    input CreateObjectTypeInput,
) (*entity.ObjectType, error) {
    // Validate input
    if err := s.validateObjectType(input); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Check name uniqueness
    existing, _ := s.repo.GetByName(ctx, input.Name)
    if existing != nil {
        return nil, ErrObjectTypeNameExists
    }
    
    // Create entity
    objectType := &entity.ObjectType{
        ID:          uuid.New(),
        Name:        input.Name,
        DisplayName: input.DisplayName,
        Description: input.Description,
        Category:    input.Category,
        Tags:        input.Tags,
        Properties:  s.buildProperties(input.Properties),
        Metadata:    input.Metadata,
        Version:     1,
        CreatedBy:   getUserFromContext(ctx),
        UpdatedBy:   getUserFromContext(ctx),
    }
    
    // Save to database
    if err := s.repo.Create(ctx, objectType); err != nil {
        return nil, fmt.Errorf("failed to create object type: %w", err)
    }
    
    // Invalidate cache
    s.cache.Delete(ctx, "object_type:"+objectType.ID.String())
    s.cache.Delete(ctx, "object_types:list")
    
    // Publish event
    event := messaging.Event{
        Type:      "ObjectTypeCreated",
        EntityID:  objectType.ID.String(),
        Actor:     getUserFromContext(ctx),
        Timestamp: time.Now(),
        Data:      objectType,
    }
    
    if err := s.publisher.Publish(ctx, event); err != nil {
        // Log error but don't fail the operation
        log.Error("failed to publish event", "error", err)
    }
    
    return objectType, nil
}

func (s *ObjectTypeService) validateObjectType(input CreateObjectTypeInput) error {
    // Name validation
    if !isValidName(input.Name) {
        return fmt.Errorf("invalid name format: must start with letter and contain only alphanumeric and underscore")
    }
    
    if len(input.Name) > 64 {
        return fmt.Errorf("name too long: maximum 64 characters")
    }
    
    // Properties validation
    propertyNames := make(map[string]bool)
    for _, prop := range input.Properties {
        if propertyNames[prop.Name] {
            return fmt.Errorf("duplicate property name: %s", prop.Name)
        }
        propertyNames[prop.Name] = true
        
        if err := s.validateProperty(prop); err != nil {
            return fmt.Errorf("invalid property %s: %w", prop.Name, err)
        }
    }
    
    return nil
}

// validatePropertyValue validates a value against property validators
func (s *ObjectTypeService) validatePropertyValue(prop Property, value interface{}) error {
    if prop.Required && value == nil {
        return fmt.Errorf("required property %s is missing", prop.Name)
    }
    
    for _, validator := range prop.Validators {
        if err := applyValidator(validator, value, prop.DataType); err != nil {
            return fmt.Errorf("validation failed for %s: %w", prop.Name, err)
        }
    }
    
    return nil
}

func applyValidator(validator Validator, value interface{}, dataType DataType) error {
    switch validator.Type {
    case ValidatorMinLength:
        if dataType != DataTypeString {
            return fmt.Errorf("minLength validator only applies to string type")
        }
        str, ok := value.(string)
        if !ok {
            return fmt.Errorf("value is not a string")
        }
        minLen, ok := validator.Value.(float64)
        if !ok {
            return fmt.Errorf("invalid minLength value")
        }
        if len(str) < int(minLen) {
            return fmt.Errorf("string length %d is less than minimum %d", len(str), int(minLen))
        }
        
    case ValidatorMaxLength:
        if dataType != DataTypeString {
            return fmt.Errorf("maxLength validator only applies to string type")
        }
        str, ok := value.(string)
        if !ok {
            return fmt.Errorf("value is not a string")
        }
        maxLen, ok := validator.Value.(float64)
        if !ok {
            return fmt.Errorf("invalid maxLength value")
        }
        if len(str) > int(maxLen) {
            return fmt.Errorf("string length %d exceeds maximum %d", len(str), int(maxLen))
        }
        
    case ValidatorPattern:
        if dataType != DataTypeString {
            return fmt.Errorf("pattern validator only applies to string type")
        }
        str, ok := value.(string)
        if !ok {
            return fmt.Errorf("value is not a string")
        }
        pattern, ok := validator.Value.(string)
        if !ok {
            return fmt.Errorf("invalid pattern value")
        }
        matched, err := regexp.MatchString(pattern, str)
        if err != nil {
            return fmt.Errorf("invalid regex pattern: %w", err)
        }
        if !matched {
            return fmt.Errorf("value does not match pattern %s", pattern)
        }
        
    case ValidatorMin:
        if dataType != DataTypeNumber {
            return fmt.Errorf("min validator only applies to number type")
        }
        num, ok := value.(float64)
        if !ok {
            return fmt.Errorf("value is not a number")
        }
        min, ok := validator.Value.(float64)
        if !ok {
            return fmt.Errorf("invalid min value")
        }
        if num < min {
            return fmt.Errorf("value %f is less than minimum %f", num, min)
        }
        
    case ValidatorMax:
        if dataType != DataTypeNumber {
            return fmt.Errorf("max validator only applies to number type")
        }
        num, ok := value.(float64)
        if !ok {
            return fmt.Errorf("value is not a number")
        }
        max, ok := validator.Value.(float64)
        if !ok {
            return fmt.Errorf("invalid max value")
        }
        if num > max {
            return fmt.Errorf("value %f exceeds maximum %f", num, max)
        }
    }
    
    return nil
}
```

### 4.5 GraphQL Implementation

#### Schema
```graphql
scalar Time
scalar JSON

type ObjectType {
    id: ID!
    name: String!
    displayName: String!
    description: String
    category: String
    tags: [String!]
    properties: [Property!]!
    metadata: JSON!
    version: Int!
    createdAt: Time!
    createdBy: String!
    updatedAt: Time!
    updatedBy: String!
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
    validators: [Validator!]
    metadata: JSON!
}

enum DataType {
    STRING
    NUMBER
    BOOLEAN
    DATE
    DATETIME
    ARRAY
    OBJECT
    REFERENCE
}

type Query {
    objectType(id: ID!): ObjectType
    objectTypes(
        filter: ObjectTypeFilter
        pagination: PaginationInput
    ): ObjectTypeConnection!
    searchObjectTypes(query: String!, limit: Int = 10): [ObjectType!]!
    compareObjectTypeVersions(id: ID!, version1: Int!, version2: Int!): VersionComparison!
}

type VersionComparison {
    objectTypeId: ID!
    version1: ObjectTypeVersion!
    version2: ObjectTypeVersion!
    changes: [FieldChange!]!
}

type ObjectTypeVersion {
    version: Int!
    snapshot: ObjectType!
    changeDescription: String
    createdAt: Time!
    createdBy: String!
}

type FieldChange {
    field: String!
    oldValue: JSON
    newValue: JSON
    changeType: ChangeType!
}

enum ChangeType {
    ADDED
    REMOVED
    MODIFIED
}

type Mutation {
    createObjectType(input: CreateObjectTypeInput!): ObjectType!
    updateObjectType(id: ID!, input: UpdateObjectTypeInput!): ObjectType!
    deleteObjectType(id: ID!): Boolean!
}

input CreateObjectTypeInput {
    name: String!
    displayName: String!
    description: String
    category: String
    tags: [String!]
    properties: [PropertyInput!]!
    metadata: JSON
}
```

#### Resolver
```go
package graphql

import (
    "context"
    "github.com/openfoundry/oms/internal/domain/service"
)

type Resolver struct {
    objectTypeService *service.ObjectTypeService
    linkTypeService   *service.LinkTypeService
}

func (r *queryResolver) ObjectType(ctx context.Context, id string) (*model.ObjectType, error) {
    objectTypeID, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid ID format")
    }
    
    // Try cache first
    cacheKey := fmt.Sprintf("object_type:%s", id)
    if cached, found := r.cache.Get(ctx, cacheKey); found {
        return cached.(*model.ObjectType), nil
    }
    
    // Get from service
    objectType, err := r.objectTypeService.GetByID(ctx, objectTypeID)
    if err != nil {
        return nil, err
    }
    
    // Convert to GraphQL model
    result := r.toGraphQLObjectType(objectType)
    
    // Cache the result
    r.cache.Set(ctx, cacheKey, result, 5*time.Minute)
    
    return result, nil
}
```

### 4.6 Repository Implementation

#### PostgreSQL Repository with Full-text Search

```go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "github.com/lib/pq"
    "github.com/openfoundry/oms/internal/domain/entity"
)

type PostgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}

// Search implements full-text search using PostgreSQL's tsvector
func (r *PostgresRepository) Search(ctx context.Context, query string, limit int) ([]*entity.ObjectType, error) {
    sql := `
        SELECT id, name, display_name, description, category, tags, properties, 
               base_datasets, metadata, version, created_at, created_by, updated_at, updated_by
        FROM object_types 
        WHERE to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, '')) 
        @@ plainto_tsquery('english', $1)
        AND is_deleted = FALSE
        ORDER BY ts_rank(to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, '')), 
                        plainto_tsquery('english', $1)) DESC
        LIMIT $2
    `
    
    rows, err := r.db.QueryContext(ctx, sql, query, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to search object types: %w", err)
    }
    defer rows.Close()
    
    var results []*entity.ObjectType
    for rows.Next() {
        var ot entity.ObjectType
        var props, datasets, metadata []byte
        
        err := rows.Scan(
            &ot.ID, &ot.Name, &ot.DisplayName, &ot.Description,
            &ot.Category, pq.Array(&ot.Tags), &props, &datasets,
            &metadata, &ot.Version, &ot.CreatedAt, &ot.CreatedBy,
            &ot.UpdatedAt, &ot.UpdatedBy,
        )
        if err != nil {
            return nil, err
        }
        
        // Unmarshal JSON fields
        if err := json.Unmarshal(props, &ot.Properties); err != nil {
            return nil, err
        }
        if err := json.Unmarshal(datasets, &ot.BaseDatasets); err != nil {
            return nil, err
        }
        if err := json.Unmarshal(metadata, &ot.Metadata); err != nil {
            return nil, err
        }
        
        results = append(results, &ot)
    }
    
    return results, rows.Err()
}

// List with cursor-based pagination
func (r *PostgresRepository) List(ctx context.Context, filter ObjectTypeFilter) ([]*entity.ObjectType, error) {
    query := `
        SELECT id, name, display_name, description, category, tags, 
               properties, base_datasets, metadata, version, 
               created_at, created_by, updated_at, updated_by
        FROM object_types
        WHERE is_deleted = FALSE
    `
    
    var args []interface{}
    argCount := 0
    
    // Handle cursor-based pagination
    if filter.PageCursor != "" {
        cursor, err := DecodeCursor(filter.PageCursor)
        if err != nil {
            return nil, fmt.Errorf("invalid cursor: %w", err)
        }
        argCount++
        query += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)", argCount, argCount+1)
        args = append(args, cursor.Timestamp, cursor.ID)
        argCount++
    }
    
    // Apply filters
    if filter.Category != nil {
        argCount++
        query += fmt.Sprintf(" AND category = $%d", argCount)
        args = append(args, *filter.Category)
    }
    
    if len(filter.Tags) > 0 {
        argCount++
        query += fmt.Sprintf(" AND tags && $%d", argCount)
        args = append(args, pq.Array(filter.Tags))
    }
    
    // Order and limit
    query += " ORDER BY created_at DESC, id DESC"
    if filter.PageSize > 0 {
        argCount++
        query += fmt.Sprintf(" LIMIT $%d", argCount)
        args = append(args, filter.PageSize)
    }
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Process results (similar to Search method)
    var results []*entity.ObjectType
    // ... (same processing as Search method)
    
    return results, nil
}
```

## 5. Event System

### 5.1 Event Types

```go
package messaging

type EventType string

const (
    EventObjectTypeCreated EventType = "ObjectTypeCreated"
    EventObjectTypeUpdated EventType = "ObjectTypeUpdated"
    EventObjectTypeDeleted EventType = "ObjectTypeDeleted"
    EventLinkTypeCreated   EventType = "LinkTypeCreated"
    EventLinkTypeUpdated   EventType = "LinkTypeUpdated"
    EventLinkTypeDeleted   EventType = "LinkTypeDeleted"
)

type Event struct {
    ID           string                 `json:"id"`
    Type         EventType              `json:"type"`
    EntityID     string                 `json:"entityId"`
    Actor        string                 `json:"actor"`
    Timestamp    time.Time              `json:"timestamp"`
    Data         interface{}            `json:"data"`
    Metadata     map[string]interface{} `json:"metadata"`
    CorrelationID string                `json:"correlationId,omitempty"`
}
```

### 5.2 Kafka Publisher

```go
package messaging

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
)

type KafkaEventPublisher struct {
    writer *kafka.Writer
}

func NewKafkaEventPublisher(brokers []string, topic string) *KafkaEventPublisher {
    writer := &kafka.Writer{
        Addr:     kafka.TCP(brokers...),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
    
    return &KafkaEventPublisher{writer: writer}
}

func (p *KafkaEventPublisher) Publish(ctx context.Context, event Event) error {
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    message := kafka.Message{
        Key:   []byte(event.EntityID),
        Value: data,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte(event.Type)},
            {Key: "correlation-id", Value: []byte(event.CorrelationID)},
        },
    }
    
    return p.writer.WriteMessages(ctx, message)
}
```

## 6. Caching Strategy

### 6.1 Cache Keys

```
object_type:{id}                    - Single object type
object_type:name:{name}             - Object type by name
object_types:list:{hash}            - List query results
object_types:search:{query}:{limit} - Search results
link_types:source:{id}              - Links by source
link_types:target:{id}              - Links by target
```

### 6.2 Cache Implementation

```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
    return &RedisCache{
        client: client,
        ttl:    ttl,
    }
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return ErrCacheMiss
    }
    if err != nil {
        return fmt.Errorf("redis get error: %w", err)
    }
    
    return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("marshal error: %w", err)
    }
    
    return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisCache) InvalidatePattern(ctx context.Context, pattern string) error {
    var cursor uint64
    for {
        keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return err
        }
        
        if len(keys) > 0 {
            if err := c.client.Del(ctx, keys...).Err(); err != nil {
                return err
            }
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    
    return nil
}
```

## 7. Error Handling

### 7.1 Error Types

```go
package errors

import "errors"

var (
    // Domain errors
    ErrObjectTypeNotFound   = &UserError{Code: "OBJECT_TYPE_NOT_FOUND", Message: "요청하신 객체 타입을 찾을 수 없습니다"}
    ErrObjectTypeNameExists = &UserError{Code: "OBJECT_TYPE_NAME_EXISTS", Message: "이미 존재하는 객체 타입 이름입니다"}
    ErrInvalidObjectType    = &UserError{Code: "INVALID_OBJECT_TYPE", Message: "유효하지 않은 객체 타입입니다"}
    ErrCircularReference    = &UserError{Code: "CIRCULAR_REFERENCE", Message: "순환 참조가 감지되었습니다"}
    
    // Validation errors
    ErrInvalidName          = &UserError{Code: "INVALID_NAME", Message: "이름 형식이 올바르지 않습니다. 영문자로 시작하고 영숫자와 밑줄만 포함해야 합니다"}
    ErrInvalidDataType      = &UserError{Code: "INVALID_DATA_TYPE", Message: "유효하지 않은 데이터 타입입니다"}
    ErrRequiredField        = &UserError{Code: "REQUIRED_FIELD_MISSING", Message: "필수 필드가 누락되었습니다"}
    
    // System errors
    ErrDatabaseConnection   = errors.New("database connection error")
    ErrCacheConnection      = errors.New("cache connection error")
    ErrEventPublishing      = errors.New("event publishing error")
)

type UserError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *UserError) Error() string {
    return e.Message
}

type ValidationError struct {
    Field   string
    Message string
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
    // Implementation
}
```

## 8. Monitoring and Observability

### 8.1 Metrics

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "oms_request_duration_seconds",
            Help: "Duration of OMS requests in seconds",
        },
        []string{"method", "endpoint", "status"},
    )
    
    objectTypeCount = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "oms_object_types_total",
            Help: "Total number of object types",
        },
    )
    
    cacheHitRate = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "oms_cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache_key_type"},
    )
)
```

### 8.2 Logging

```go
package logger

import (
    "go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    return config.Build()
}

// Usage
logger.Info("object type created",
    zap.String("id", objectType.ID.String()),
    zap.String("name", objectType.Name),
    zap.String("actor", actor),
)
```

## 9. Security Implementation

### 9.1 Authentication Middleware

```go
package middleware

import (
    "context"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "authorization header missing"})
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }
        
        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["sub"])
        c.Set("user_roles", claims["roles"])
        
        c.Next()
    }
}
```

### 9.2 RBAC Implementation

```go
package auth

type Permission string

const (
    PermissionObjectTypeRead   Permission = "object_type:read"
    PermissionObjectTypeWrite  Permission = "object_type:write"
    PermissionObjectTypeDelete Permission = "object_type:delete"
)

func HasPermission(ctx context.Context, permission Permission) bool {
    roles := getRolesFromContext(ctx)
    
    // Check role-permission mapping
    for _, role := range roles {
        if roleHasPermission(role, permission) {
            return true
        }
    }
    
    return false
}
```

## 10. Testing Strategy

### 10.1 Unit Test Example

```go
package service_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestObjectTypeService_Create(t *testing.T) {
    // Setup
    mockRepo := new(MockObjectTypeRepository)
    mockCache := new(MockCacheService)
    mockPublisher := new(MockEventPublisher)
    
    service := NewObjectTypeService(mockRepo, mockCache, mockPublisher)
    
    // Test case
    input := CreateObjectTypeInput{
        Name:        "Customer",
        DisplayName: "Customer",
        Properties: []PropertyInput{
            {
                Name:     "name",
                DataType: "STRING",
                Required: true,
            },
        },
    }
    
    mockRepo.On("GetByName", mock.Anything, "Customer").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)
    mockPublisher.On("Publish", mock.Anything, mock.Anything).Return(nil)
    
    // Execute
    result, err := service.CreateObjectType(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Customer", result.Name)
    mockRepo.AssertExpectations(t)
}
```

## 11. Configuration

### 11.1 Environment Variables

```yaml
# Server
SERVER_PORT: 8080
SERVER_MODE: production # development, production

# Database
DB_HOST: localhost
DB_PORT: 5432
DB_NAME: oms
DB_USER: oms_user
DB_PASSWORD: ${DB_PASSWORD}
DB_SSL_MODE: require

# Redis
REDIS_HOST: localhost
REDIS_PORT: 6379
REDIS_PASSWORD: ${REDIS_PASSWORD}
REDIS_DB: 0

# Kafka
KAFKA_BROKERS: localhost:9092
KAFKA_TOPIC: oms-events

# Security
JWT_SECRET: ${JWT_SECRET}
API_KEY_HEADER: X-API-Key

# Monitoring
METRICS_PATH: /metrics
TRACE_ENDPOINT: http://jaeger:14268/api/traces
```

## 12. Deployment Considerations

### 12.1 Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o oms-service cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/oms-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./oms-service"]
```

### 12.2 Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: oms-service
  namespace: openfoundry
spec:
  replicas: 3
  selector:
    matchLabels:
      app: oms-service
  template:
    metadata:
      labels:
        app: oms-service
    spec:
      containers:
      - name: oms-service
        image: openfoundry/oms-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: oms-secrets
              key: db-password
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 13. Performance Optimization

### 13.1 Database Optimizations

1. **Connection Pooling**
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

2. **Query Optimization**
   - Use prepared statements
   - Implement pagination with cursor-based approach
   - Add appropriate indexes

3. **Batch Operations**
   ```go
   func (r *PostgresRepository) BatchCreate(ctx context.Context, items []*entity.ObjectType) error {
       tx, err := r.db.BeginTx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()
       
       stmt, err := tx.Prepare(insertQuery)
       if err != nil {
           return err
       }
       defer stmt.Close()
       
       for _, item := range items {
           if _, err := stmt.Exec(item.ID, item.Name, ...); err != nil {
               return err
           }
       }
       
       return tx.Commit()
   }
   ```

## 14. Migration Strategy

### 14.1 Database Migration Tool

Using golang-migrate:

```bash
migrate create -ext sql -dir migrations create_object_types_table
migrate -path migrations -database "postgresql://..." up
```

### 14.2 Zero-Downtime Migration

1. Add new columns as nullable
2. Deploy new code that writes to both old and new columns
3. Backfill data
4. Deploy code that reads from new columns
5. Remove old columns

## 15. Next Steps

1. **Phase 2 Integration Points**
   - Prepare interfaces for OSv2 integration
   - Design event contracts for Object Data Funnel
   - Plan API extensions for Functions on Objects

2. **Scalability Preparations**
   - Implement database sharding strategy
   - Design multi-tenant architecture
   - Plan for horizontal scaling

3. **Production Readiness**
   - Complete security audit
   - Performance testing and optimization
   - Disaster recovery planning

## 16. Microservices Architecture

### 16.1 Service Overview

OMS는 OpenFoundry 플랫폼의 기반 마이크로서비스로, 다른 서비스들이 의존하는 메타데이터 레지스트리 역할을 합니다.

```yaml
Service Identity:
  Name: oms-service
  Domain: metadata.openfoundry.io
  Port: 
    HTTP/GraphQL: 8080
    gRPC: 9090
    Metrics: 9091
  
Dependencies:
  - PostgreSQL (데이터 저장)
  - Redis (캐싱)
  - Kafka (이벤트 스트리밍)
  
Consumers:
  - OSv2 (객체 인스턴스 서비스)
  - OSS (검색 서비스)
  - FOO (함수 실행 서비스)
```

### 16.2 gRPC Service Definition

```protobuf
syntax = "proto3";

package openfoundry.oms.v1;

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

// ObjectType Service - 다른 마이크로서비스를 위한 내부 API
service ObjectTypeService {
  // 단일 객체 타입 조회
  rpc GetObjectType(GetObjectTypeRequest) returns (ObjectType);
  
  // 객체 타입 목록 조회
  rpc ListObjectTypes(ListObjectTypesRequest) returns (ListObjectTypesResponse);
  
  // 객체 타입 검증
  rpc ValidateObjectInstance(ValidateObjectInstanceRequest) returns (ValidateObjectInstanceResponse);
  
  // 스트리밍 API - 실시간 변경사항 구독
  rpc WatchObjectTypes(WatchObjectTypesRequest) returns (stream ObjectTypeEvent);
}

message ObjectType {
  string id = 1;
  string name = 2;
  string display_name = 3;
  string description = 4;
  repeated Property properties = 5;
  map<string, string> metadata = 6;
  int32 version = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}

message Property {
  string id = 1;
  string name = 2;
  DataType data_type = 3;
  bool required = 4;
  bool unique = 5;
  string default_value = 6;
  repeated Validator validators = 7;
}

enum DataType {
  DATA_TYPE_UNSPECIFIED = 0;
  DATA_TYPE_STRING = 1;
  DATA_TYPE_NUMBER = 2;
  DATA_TYPE_BOOLEAN = 3;
  DATA_TYPE_DATE = 4;
  DATA_TYPE_DATETIME = 5;
  DATA_TYPE_ARRAY = 6;
  DATA_TYPE_OBJECT = 7;
  DATA_TYPE_REFERENCE = 8;
}

message Validator {
  ValidatorType type = 1;
  string value = 2;
}

enum ValidatorType {
  VALIDATOR_TYPE_UNSPECIFIED = 0;
  VALIDATOR_TYPE_MIN_LENGTH = 1;
  VALIDATOR_TYPE_MAX_LENGTH = 2;
  VALIDATOR_TYPE_PATTERN = 3;
  VALIDATOR_TYPE_MIN = 4;
  VALIDATOR_TYPE_MAX = 5;
}
```

### 16.3 API Versioning Strategy

```go
// API 버전 관리를 위한 미들웨어
func APIVersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("X-API-Version")
        if version == "" {
            version = "v1" // 기본 버전
        }
        
        c.Set("api_version", version)
        c.Next()
    }
}

// 버전별 응답 변환
type VersionedSerializer interface {
    SerializeV1() interface{}
    SerializeV2() interface{}
}

func (ot *ObjectType) SerializeV1() interface{} {
    return struct {
        ID          string     `json:"id"`
        Name        string     `json:"name"`
        DisplayName string     `json:"displayName"`
        Properties  []Property `json:"properties"`
    }{
        ID:          ot.ID.String(),
        Name:        ot.Name,
        DisplayName: ot.DisplayName,
        Properties:  ot.Properties,
    }
}

func (ot *ObjectType) SerializeV2() interface{} {
    // v2에는 모든 필드 포함
    return ot
}
```

### 16.4 Service Discovery & Circuit Breaker

```go
// Service discovery using Consul
type ServiceDiscovery struct {
    client *consul.Client
}

func NewServiceDiscovery(consulAddr string) (*ServiceDiscovery, error) {
    config := consul.DefaultConfig()
    config.Address = consulAddr
    
    client, err := consul.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ServiceDiscovery{client: client}, nil
}

func (sd *ServiceDiscovery) RegisterService(name, address string, port int) error {
    return sd.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s-%d", name, address, port),
        Name:    name,
        Address: address,
        Port:    port,
        Check: &consul.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
            Interval: "10s",
            Timeout:  "3s",
        },
    })
}

// Circuit breaker pattern
type CircuitBreaker struct {
    maxFailures  int
    timeout      time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
    mu           sync.RWMutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.RLock()
    state := cb.state
    cb.mu.RUnlock()
    
    if state == "open" {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.mu.Lock()
            cb.state = "half-open"
            cb.failures = 0
            cb.mu.Unlock()
        } else {
            return ErrCircuitBreakerOpen
        }
    }
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    cb.state = "closed"
    cb.failures = 0
    return nil
}
```

### 16.5 Event-Driven Architecture

```go
// Domain events
type DomainEvent interface {
    EventType() string
    AggregateID() string
    EventVersion() int
    OccurredAt() time.Time
}

type ObjectTypeCreatedEvent struct {
    ID          string    `json:"id"`
    ObjectType  ObjectType `json:"objectType"`
    CreatedBy   string    `json:"createdBy"`
    Timestamp   time.Time `json:"timestamp"`
}

func (e ObjectTypeCreatedEvent) EventType() string { return "ObjectTypeCreated" }
func (e ObjectTypeCreatedEvent) AggregateID() string { return e.ID }
func (e ObjectTypeCreatedEvent) EventVersion() int { return 1 }
func (e ObjectTypeCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

// Event publisher with retry and dead letter queue
type EventPublisher struct {
    producer *kafka.Producer
    dlq      *kafka.Producer
}

func (ep *EventPublisher) PublishWithRetry(ctx context.Context, event DomainEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    msg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &[]string{"oms.events"}[0],
            Partition: kafka.PartitionAny,
        },
        Key:   []byte(event.AggregateID()),
        Value: data,
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte(event.EventType())},
            {Key: "event_version", Value: []byte(fmt.Sprintf("%d", event.EventVersion()))},
        },
    }
    
    // Retry logic
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err = ep.producer.Produce(msg, nil)
        if err == nil {
            return nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
        }
    }
    
    // Send to DLQ if all retries failed
    dlqMsg := *msg
    dlqMsg.TopicPartition.Topic = &[]string{"oms.events.dlq"}[0]
    dlqMsg.Headers = append(dlqMsg.Headers, kafka.Header{
        Key:   "original_error",
        Value: []byte(err.Error()),
    })
    
    return ep.dlq.Produce(&dlqMsg, nil)
}
```

### 16.6 Distributed Tracing

```go
// OpenTelemetry integration
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func InitTracing(serviceName string) (func(), error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    
    return func() {
        tp.Shutdown(context.Background())
    }, nil
}

// Trace middleware for Gin
func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tracer := otel.Tracer("oms-service")
        
        ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
        defer span.End()
        
        span.SetAttributes(
            attribute.String("http.method", c.Request.Method),
            attribute.String("http.url", c.Request.URL.String()),
            attribute.String("http.user_agent", c.Request.UserAgent()),
        )
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
        
        span.SetAttributes(
            attribute.Int("http.status_code", c.Writer.Status()),
        )
    }
}

// Propagate trace context in gRPC
func NewGRPCServer() *grpc.Server {
    return grpc.NewServer(
        grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
        grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
    )
}
```

### 16.7 Service Mesh Integration

```yaml
# Istio VirtualService configuration
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: oms-service
  namespace: openfoundry
spec:
  hosts:
  - oms-service
  http:
  - match:
    - headers:
        x-api-version:
          exact: v2
    route:
    - destination:
        host: oms-service
        subset: v2
      weight: 100
  - route:
    - destination:
        host: oms-service
        subset: v1
      weight: 100

---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: oms-service
  namespace: openfoundry
spec:
  host: oms-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 50
        http2MaxRequests: 100
    loadBalancer:
      simple: LEAST_REQUEST
    outlierDetection:
      consecutiveErrors: 5
      interval: 30s
      baseEjectionTime: 30s
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
```

### 16.8 Health Checks & Readiness Probes

```go
// Comprehensive health check implementation
type HealthChecker struct {
    db       *sql.DB
    redis    *redis.Client
    kafka    *kafka.Producer
}

type HealthStatus struct {
    Status string                 `json:"status"`
    Checks map[string]CheckResult `json:"checks"`
}

type CheckResult struct {
    Status  string                 `json:"status"`
    Message string                 `json:"message,omitempty"`
    Details map[string]interface{} `json:"details,omitempty"`
}

func (hc *HealthChecker) CheckHealth(ctx context.Context) HealthStatus {
    checks := make(map[string]CheckResult)
    overallHealthy := true
    
    // Database check
    dbCheck := hc.checkDatabase(ctx)
    checks["database"] = dbCheck
    if dbCheck.Status != "healthy" {
        overallHealthy = false
    }
    
    // Redis check
    redisCheck := hc.checkRedis(ctx)
    checks["redis"] = redisCheck
    if redisCheck.Status != "healthy" {
        overallHealthy = false
    }
    
    // Kafka check
    kafkaCheck := hc.checkKafka(ctx)
    checks["kafka"] = kafkaCheck
    if kafkaCheck.Status != "healthy" && kafkaCheck.Status != "degraded" {
        overallHealthy = false
    }
    
    status := "healthy"
    if !overallHealthy {
        status = "unhealthy"
    }
    
    return HealthStatus{
        Status: status,
        Checks: checks,
    }
}

// Kubernetes probes
func RegisterHealthEndpoints(router *gin.Engine, hc *HealthChecker) {
    // Liveness probe - basic check
    router.GET("/health/live", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "alive"})
    })
    
    // Readiness probe - comprehensive check
    router.GET("/health/ready", func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
        defer cancel()
        
        health := hc.CheckHealth(ctx)
        
        statusCode := 200
        if health.Status != "healthy" {
            statusCode = 503
        }
        
        c.JSON(statusCode, health)
    })
}
```

### 16.9 Container & Deployment Configuration

```dockerfile
# Multi-stage Dockerfile for OMS service
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o oms-service cmd/server/main.go

# Final stage
FROM scratch

# Copy timezone data and CA certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/oms-service /oms-service

# Copy migration files
COPY --from=builder /app/migrations /migrations

# Expose ports
EXPOSE 8080 9090 9091

# Run the service
ENTRYPOINT ["/oms-service"]
```

```yaml
# Kubernetes Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: oms-service
  namespace: openfoundry
  labels:
    app: oms-service
    version: v1
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: oms-service
      version: v1
  template:
    metadata:
      labels:
        app: oms-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9091"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: oms-service
      containers:
      - name: oms-service
        image: openfoundry/oms-service:latest
        ports:
        - name: http
          containerPort: 8080
        - name: grpc
          containerPort: 9090
        - name: metrics
          containerPort: 9091
        env:
        - name: SERVICE_NAME
          value: "oms-service"
        - name: LOG_LEVEL
          value: "info"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: oms-db-secret
              key: host
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
          initialDelaySeconds: 15
          periodSeconds: 5
        volumeMounts:
        - name: config
          mountPath: /etc/oms
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: oms-config
```

## 17. MSA Integration Roadmap

### 17.1 Phase 1 - Foundation (Current)
- Complete OMS service implementation
- Establish gRPC service contracts
- Implement comprehensive event publishing
- Set up distributed tracing

### 17.2 Phase 2 - Service Integration (Month 2-3)
- Integrate with service mesh (Istio)
- Implement circuit breakers for all external calls
- Establish service discovery patterns
- Create client SDKs for other services

### 17.3 Phase 3 - Advanced Patterns (Month 4-6)
- Implement CQRS for read-heavy operations
- Add GraphQL federation support
- Implement distributed caching strategies
- Advanced monitoring and alerting

### 17.4 Integration Points for Future Services

```yaml
OSv2 Integration:
  - Subscribe to: ObjectTypeCreated, ObjectTypeUpdated events
  - Call: ValidateObjectInstance gRPC method
  - Provide: Object instance data for OMS metadata enrichment

OSS Integration:
  - Subscribe to: All OMS events for indexing
  - Call: GetObjectType, ListObjectTypes for schema information
  - Provide: Search results with OMS metadata

FOO Integration:
  - Subscribe to: ObjectTypeUpdated for function schema updates
  - Call: GetObjectType for function parameter validation
  - Provide: Function execution results for OMS analytics
```