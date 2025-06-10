# Ontology Metadata Service (OMS) Frontend Development Specification

## 1. Executive Summary

본 문서는 OpenFoundry OMS의 프론트엔드 구현을 위한 상세 개발 명세서입니다. Palantir Blueprint.js를 활용하여 엔터프라이즈급 데이터 모델링 인터페이스를 구축하며, PRD에 정의된 모든 기능을 사용자 친화적으로 구현합니다.

## 2. Architecture Overview

### 2.1 Technology Stack

```yaml
Core Framework:
  - React 18.2+
  - TypeScript 5.0+
  - Blueprint.js 5.x (@blueprintjs/core, @blueprintjs/table, @blueprintjs/select)
  - @blueprintjs/icons (562+ icons)
  - @blueprintjs/popover2 (improved popovers)

State Management:
  - Zustand 4.x (경량 상태 관리)
  - React Query 4.x (서버 상태 관리)

API Communication:
  - Apollo Client 3.x (GraphQL)
  - Axios (REST fallback)

Build Tools:
  - Vite 4.x
  - pnpm (패키지 매니저)

Testing:
  - Jest + React Testing Library
  - Cypress (E2E)

Code Quality:
  - ESLint + Prettier
  - TypeScript strict mode
  - ts-prune (unused exports)
  - madge (circular dependencies)

Styling:
  - SCSS with Blueprint variables
  - CSS Modules for component isolation
  - PostCSS for optimization
```

### 2.2 Application Architecture

```
src/
├── app/                      # 애플리케이션 진입점
│   ├── App.tsx
│   ├── Router.tsx
│   └── Providers.tsx
├── features/                 # 기능별 모듈
│   ├── object-types/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── types/
│   ├── link-types/
│   └── search/
├── shared/                   # 공통 모듈
│   ├── components/
│   ├── hooks/
│   ├── services/
│   └── utils/
├── design-system/           # Blueprint 확장 (design-system.md 참조)
│   ├── components/
│   ├── themes/
│   └── tokens/
└── assets/
```

## 3. User Interface Design

### 3.1 Layout Structure

```tsx
// Main Application Layout
<div className="oms-app">
  <Navbar className="bp5-dark" fixedToTop>
    <NavbarGroup align="left">
      <NavbarHeading>OpenFoundry OMS</NavbarHeading>
      <NavbarDivider />
      <Button icon="cube" text="Object Types" />
      <Button icon="link" text="Relationships" />
    </NavbarGroup>
    <NavbarGroup align="right">
      <InputGroup 
        leftIcon="search" 
        placeholder="Search objects..."
      />
      <NavbarDivider />
      <Button icon="user" minimal />
    </NavbarGroup>
  </Navbar>
  
  <div className="main-content">
    <Sidebar />
    <ContentArea />
  </div>
</div>
```

### 3.2 Core Views

#### 3.2.1 Object Types List View

```tsx
interface ObjectTypesListViewProps {
  objectTypes: ObjectType[];
  onSelect: (objectType: ObjectType) => void;
  onDelete: (id: string) => void;
}

// Features:
// - Card-based layout with Blueprint Cards
// - Filtering by category, tags
// - Sorting by name, created date
// - Bulk actions toolbar
// - Virtual scrolling for performance
```

#### 3.2.2 Object Type Designer

```tsx
interface ObjectTypeDesignerProps {
  objectType?: ObjectType;
  onSave: (objectType: ObjectType) => void;
}

// Layout using Blueprint components:
<div className="object-type-designer">
  <Tabs id="designer-tabs" large>
    <Tab id="properties" title="Properties" panel={<PropertiesPanel />} />
    <Tab id="relationships" title="Relationships" panel={<RelationshipsPanel />} />
    <Tab id="metadata" title="Metadata" panel={<MetadataPanel />} />
    <Tab id="versions" title="Version History" panel={<VersionsPanel />} />
  </Tabs>
</div>
```

## 4. Component Specifications

### 4.1 Object Type Components

#### 4.1.1 ObjectTypeCard
```tsx
interface ObjectTypeCardProps {
  objectType: ObjectType;
  onEdit: () => void;
  onDelete: () => void;
  onDuplicate: () => void;
}

const ObjectTypeCard: React.FC<ObjectTypeCardProps> = ({ objectType, ...actions }) => {
  return (
    <Card interactive elevation={Elevation.TWO}>
      <div className="card-header">
        <H5>{objectType.displayName}</H5>
        <Tag intent="primary">{objectType.category}</Tag>
      </div>
      
      <p className="bp5-text-muted">{objectType.description}</p>
      
      <div className="card-stats">
        <Tag minimal icon="properties">
          {objectType.properties.length} properties
        </Tag>
        <Tag minimal icon="updated">
          v{objectType.version}
        </Tag>
      </div>
      
      <div className="card-actions">
        <ButtonGroup minimal>
          <Button icon="edit" onClick={actions.onEdit}>Edit</Button>
          <Button icon="duplicate" onClick={actions.onDuplicate}>Duplicate</Button>
          <Popover
            content={
              <Menu>
                <MenuItem 
                  icon="trash" 
                  text="Delete" 
                  intent="danger"
                  onClick={actions.onDelete}
                />
              </Menu>
            }
          >
            <Button icon="more" />
          </Popover>
        </ButtonGroup>
      </div>
    </Card>
  );
};
```

#### 4.1.2 PropertyEditor
```tsx
interface PropertyEditorProps {
  property: Property;
  onChange: (property: Property) => void;
  onDelete: () => void;
}

const PropertyEditor: React.FC<PropertyEditorProps> = ({ property, onChange, onDelete }) => {
  return (
    <div className="property-editor">
      <ControlGroup fill>
        <InputGroup 
          placeholder="Property name"
          value={property.name}
          onChange={(e) => onChange({ ...property, name: e.target.value })}
        />
        
        <HTMLSelect
          value={property.dataType}
          onChange={(e) => onChange({ ...property, dataType: e.target.value as DataType })}
        >
          <option value="STRING">String</option>
          <option value="NUMBER">Number</option>
          <option value="BOOLEAN">Boolean</option>
          <option value="DATE">Date</option>
          <option value="DATETIME">DateTime</option>
          <option value="ARRAY">Array</option>
          <option value="OBJECT">Object</option>
          <option value="REFERENCE">Reference</option>
        </HTMLSelect>
        
        <Switch 
          label="Required"
          checked={property.required}
          onChange={(e) => onChange({ ...property, required: e.target.checked })}
        />
        
        <Button icon="trash" intent="danger" onClick={onDelete} />
      </ControlGroup>
      
      <Collapse isOpen={showAdvanced}>
        <FormGroup label="Default Value">
          <InputGroup 
            value={property.defaultValue}
            onChange={(e) => onChange({ ...property, defaultValue: e.target.value })}
          />
        </FormGroup>
        
        <FormGroup label="Description">
          <TextArea 
            value={property.description}
            onChange={(e) => onChange({ ...property, description: e.target.value })}
          />
        </FormGroup>
      </Collapse>
    </div>
  );
};
```

#### 4.1.3 ObjectTypeForm
```tsx
const ObjectTypeForm: React.FC<ObjectTypeFormProps> = ({ objectType, onSubmit }) => {
  const [formData, setFormData] = useState<ObjectTypeFormData>({
    name: '',
    displayName: '',
    description: '',
    category: '',
    tags: [],
    properties: [],
    ...objectType
  });

  return (
    <form onSubmit={handleSubmit}>
      <FormGroup 
        label="System Name" 
        labelInfo="(required)"
        labelFor="name-input"
        helperText="Use only letters, numbers, and underscores"
        intent={nameError ? "danger" : "none"}
      >
        <InputGroup
          id="name-input"
          large
          leftIcon="tag"
          placeholder="e.g., customer_profile"
          value={formData.name}
          onChange={handleNameChange}
          intent={nameError ? "danger" : "none"}
          rightElement={
            nameError && <Icon icon="error" intent="danger" />
          }
        />
      </FormGroup>

      <FormGroup label="Display Name" labelInfo="(required)">
        <InputGroup
          large
          placeholder="Customer Profile"
          value={formData.displayName}
          onChange={(e) => setFormData({ ...formData, displayName: e.target.value })}
        />
      </FormGroup>

      <FormGroup label="Category">
        <Suggest
          items={categories}
          itemRenderer={renderCategory}
          onItemSelect={(item) => setFormData({ ...formData, category: item })}
          inputValueRenderer={(item) => item}
          createNewItemFromQuery={(query) => query}
        />
      </FormGroup>

      <FormGroup label="Tags">
        <MultiSelect
          items={availableTags}
          selectedItems={formData.tags}
          onItemSelect={handleTagSelect}
          tagRenderer={renderTag}
          placeholder="Add tags..."
        />
      </FormGroup>

      <FormGroup label="Description">
        <TextArea
          large
          fill
          rows={3}
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
        />
      </FormGroup>

      <Divider />

      <PropertiesSection 
        properties={formData.properties}
        onChange={(properties) => setFormData({ ...formData, properties })}
      />

      <div className="oms-form-actions">
        <Button large intent="primary" type="submit" icon="tick">
          {objectType ? 'Update Object Type' : 'Create Object Type'}
        </Button>
        <Button large outlined onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
};
```

### 4.2 Relationship Components

#### 4.2.1 LinkTypeEditor
```tsx
interface LinkTypeEditorProps {
  linkType?: LinkType;
  objectTypes: ObjectType[];
  onSave: (linkType: LinkType) => void;
}

const LinkTypeEditor: React.FC<LinkTypeEditorProps> = ({ linkType, objectTypes, onSave }) => {
  return (
    <Dialog
      isOpen={isOpen}
      title={linkType ? "Edit Relationship" : "Create Relationship"}
      icon="link"
    >
      <div className={Classes.DIALOG_BODY}>
        <FormGroup label="Relationship Name" labelInfo="(required)">
          <InputGroup placeholder="has_orders" />
        </FormGroup>

        <ControlGroup fill>
          <FormGroup label="Source Object">
            <HTMLSelect fill>
              {objectTypes.map(ot => (
                <option key={ot.id} value={ot.id}>{ot.displayName}</option>
              ))}
            </HTMLSelect>
          </FormGroup>

          <FormGroup label="Cardinality">
            <HTMLSelect>
              <option value="ONE_TO_ONE">1:1</option>
              <option value="ONE_TO_MANY">1:N</option>
              <option value="MANY_TO_MANY">N:N</option>
            </HTMLSelect>
          </FormGroup>

          <FormGroup label="Target Object">
            <HTMLSelect fill>
              {objectTypes.map(ot => (
                <option key={ot.id} value={ot.id}>{ot.displayName}</option>
              ))}
            </HTMLSelect>
          </FormGroup>
        </ControlGroup>
      </div>

      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={onClose}>Cancel</Button>
          <Button intent="primary" onClick={handleSave}>Save Relationship</Button>
        </div>
      </div>
    </Dialog>
  );
};
```

#### 4.2.2 RelationshipGraph
```tsx
const RelationshipGraph: React.FC<RelationshipGraphProps> = ({ objectTypes, linkTypes }) => {
  // Using D3.js or vis.js for graph visualization
  // Wrapped in Blueprint Panel component
  
  return (
    <div className="relationship-graph">
      <div className="graph-toolbar">
        <ButtonGroup>
          <Button icon="zoom-fit" text="Fit to Screen" />
          <Button icon="fullscreen" text="Fullscreen" />
          <Button icon="export" text="Export" />
        </ButtonGroup>
      </div>
      
      <div className="graph-container" ref={graphRef}>
        {/* D3/vis.js rendered here */}
      </div>
      
      <div className="graph-legend">
        <Tag>Object Type</Tag>
        <Tag minimal>1:1 Relationship</Tag>
        <Tag minimal>1:N Relationship</Tag>
        <Tag minimal>N:N Relationship</Tag>
      </div>
    </div>
  );
};
```

### 4.3 Search Components

#### 4.3.1 GlobalSearch
```tsx
const GlobalSearch: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResults>();

  return (
    <Popover
      content={
        <div className="search-results">
          <Menu>
            <MenuDivider title="Object Types" />
            {results?.objectTypes.map(ot => (
              <MenuItem
                key={ot.id}
                text={ot.displayName}
                label={ot.category}
                onClick={() => navigateToObjectType(ot.id)}
              />
            ))}
            
            <MenuDivider title="Properties" />
            {results?.properties.map(prop => (
              <MenuItem
                key={prop.id}
                text={prop.name}
                label={prop.objectTypeName}
                onClick={() => navigateToProperty(prop)}
              />
            ))}
          </Menu>
        </div>
      }
      isOpen={query.length > 0 && results != null}
      position="bottom"
    >
      <InputGroup
        large
        leftIcon="search"
        placeholder="Search objects, properties..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        rightElement={
          query && <Button icon="cross" minimal onClick={() => setQuery('')} />
        }
      />
    </Popover>
  );
};
```

### 4.4 Data Grid Components

#### 4.4.1 ObjectInstanceGrid
```tsx
const ObjectInstanceGrid: React.FC<ObjectInstanceGridProps> = ({ objectType }) => {
  // Using @blueprintjs/table for virtual scrolling
  
  return (
    <div className="object-instance-grid">
      <div className="grid-toolbar">
        <ButtonGroup>
          <Button icon="add" intent="primary" text="New Instance" />
          <Button icon="import" text="Import CSV" />
          <Button icon="export" text="Export" />
        </ButtonGroup>
        
        <InputGroup 
          leftIcon="filter" 
          placeholder="Filter instances..."
          className="grid-filter"
        />
      </div>

      <Table2
        numRows={data.length}
        enableGhostCells
        enableRowResizing={false}
        rowHeights={data.map(() => 40)}
        columnWidths={columns.map(col => col.width)}
      >
        {columns.map((col, index) => (
          <Column
            key={col.id}
            name={col.displayName}
            cellRenderer={(rowIndex) => (
              <EditableCell
                value={data[rowIndex][col.name]}
                onConfirm={(value) => handleCellEdit(rowIndex, col.name, value)}
              />
            )}
          />
        ))}
      </Table2>

      <div className="grid-footer">
        <span>{data.length} instances</span>
        <Button minimal icon="refresh" onClick={refresh} />
      </div>
    </div>
  );
};
```

### 4.5 Loading & Empty States

#### 4.5.1 Loading States
```tsx
// Full page loading
const LoadingState: React.FC = () => (
  <NonIdealState
    icon={<Spinner size={SpinnerSize.LARGE} />}
    title="Loading Object Types..."
    description="Please wait while we fetch your data"
  />
);

// Skeleton loading for cards
const ObjectTypeCardSkeleton: React.FC = () => (
  <Card className="oms-object-type-card">
    <div className="oms-card-header">
      <H4><Skeleton width="60%" /></H4>
      <Skeleton width="80px" height="20px" />
    </div>
    <p><Skeleton count={2} /></p>
    <div className="oms-card-meta">
      <Skeleton width="100px" height="16px" />
      <Skeleton width="60px" height="16px" />
    </div>
  </Card>
);
```

#### 4.5.2 Empty States
```tsx
const EmptyState: React.FC<{ onCreateNew: () => void }> = ({ onCreateNew }) => (
  <NonIdealState
    icon="folder-open"
    title="No Object Types"
    description="Create your first object type to start modeling your data"
    action={
      <Button intent="primary" icon="plus" large onClick={onCreateNew}>
        Create Object Type
      </Button>
    }
  />
);
```

#### 4.5.3 Error States
```tsx
const ErrorState: React.FC<{ error: Error; onRetry: () => void }> = ({ error, onRetry }) => (
  <Callout intent="danger" icon="error" className="oms-error-callout">
    <H4>Error Loading Data</H4>
    <p>{error.message || 'An unexpected error occurred'}</p>
    <Button intent="primary" outlined onClick={onRetry}>
      <Icon icon="refresh" /> Retry
    </Button>
  </Callout>
);
```

## 5. State Management

### 5.1 Zustand Stores

```typescript
// Object Types Store
interface ObjectTypesStore {
  objectTypes: ObjectType[];
  selectedObjectType: ObjectType | null;
  loading: boolean;
  error: Error | null;
  
  // Actions
  fetchObjectTypes: () => Promise<void>;
  createObjectType: (input: CreateObjectTypeInput) => Promise<ObjectType>;
  updateObjectType: (id: string, input: UpdateObjectTypeInput) => Promise<ObjectType>;
  deleteObjectType: (id: string) => Promise<void>;
  selectObjectType: (objectType: ObjectType | null) => void;
}

const useObjectTypesStore = create<ObjectTypesStore>((set, get) => ({
  objectTypes: [],
  selectedObjectType: null,
  loading: false,
  error: null,

  fetchObjectTypes: async () => {
    set({ loading: true, error: null });
    try {
      const data = await objectTypeService.fetchAll();
      set({ objectTypes: data, loading: false });
    } catch (error) {
      set({ error, loading: false });
    }
  },

  createObjectType: async (input) => {
    const objectType = await objectTypeService.create(input);
    set(state => ({
      objectTypes: [...state.objectTypes, objectType]
    }));
    return objectType;
  },

  // ... other actions
}));
```

### 5.2 React Query Integration

```typescript
// GraphQL Queries
const OBJECT_TYPES_QUERY = gql`
  query GetObjectTypes($filter: ObjectTypeFilter, $pagination: PaginationInput) {
    objectTypes(filter: $filter, pagination: $pagination) {
      edges {
        node {
          id
          name
          displayName
          description
          category
          tags
          properties {
            id
            name
            dataType
            required
          }
          version
          metadata {
            createdAt
            createdBy
            updatedAt
            updatedBy
          }
        }
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
`;

// Custom Hook
export const useObjectTypes = (filter?: ObjectTypeFilter) => {
  return useQuery({
    queryKey: ['objectTypes', filter],
    queryFn: () => graphqlClient.request(OBJECT_TYPES_QUERY, { filter }),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
```

## 6. API Integration

### 6.1 GraphQL Client Setup

```typescript
import { GraphQLClient } from 'graphql-request';
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';

// Apollo Client Configuration
const httpLink = createHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_ENDPOINT,
  headers: {
    authorization: `Bearer ${getAuthToken()}`,
  },
});

export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache({
    typePolicies: {
      ObjectType: {
        keyFields: ['id'],
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
  },
});
```

### 6.2 Service Layer

```typescript
// services/objectTypeService.ts
export class ObjectTypeService {
  async create(input: CreateObjectTypeInput): Promise<ObjectType> {
    const { data } = await apolloClient.mutate({
      mutation: CREATE_OBJECT_TYPE_MUTATION,
      variables: { input },
      update: (cache, { data }) => {
        // Update cache
        const existing = cache.readQuery({ query: OBJECT_TYPES_QUERY });
        cache.writeQuery({
          query: OBJECT_TYPES_QUERY,
          data: {
            objectTypes: {
              ...existing.objectTypes,
              edges: [...existing.objectTypes.edges, { node: data.createObjectType }],
            },
          },
        });
      },
    });
    
    return data.createObjectType;
  }

  async update(id: string, input: UpdateObjectTypeInput): Promise<ObjectType> {
    const { data } = await apolloClient.mutate({
      mutation: UPDATE_OBJECT_TYPE_MUTATION,
      variables: { id, input },
    });
    
    return data.updateObjectType;
  }

  // ... other methods
}

export const objectTypeService = new ObjectTypeService();
```

## 7. Performance Optimization

### 7.1 Virtual Scrolling

```typescript
// For large lists using @blueprintjs/table
const VirtualizedObjectList: React.FC = () => {
  const rowRenderer = useCallback((index: number) => {
    const objectType = objectTypes[index];
    return <ObjectTypeCard key={objectType.id} objectType={objectType} />;
  }, [objectTypes]);

  return (
    <Table2
      numRows={objectTypes.length}
      rowRenderer={rowRenderer}
      defaultRowHeight={120}
      enableRowHeader={false}
    />
  );
};
```

### 7.2 Code Splitting

```typescript
// Lazy load heavy components
const RelationshipGraph = lazy(() => import('./components/RelationshipGraph'));
const ObjectInstanceGrid = lazy(() => import('./components/ObjectInstanceGrid'));

// Route-based code splitting
const routes = [
  {
    path: '/object-types',
    component: lazy(() => import('./features/object-types')),
  },
  {
    path: '/relationships',
    component: lazy(() => import('./features/link-types')),
  },
];
```

### 7.3 Memoization

```typescript
// Expensive computations
const useFilteredObjectTypes = (objectTypes: ObjectType[], filter: Filter) => {
  return useMemo(() => {
    return objectTypes.filter(ot => {
      if (filter.category && ot.category !== filter.category) return false;
      if (filter.tags?.length && !filter.tags.some(tag => ot.tags.includes(tag))) return false;
      if (filter.search) {
        const searchLower = filter.search.toLowerCase();
        return ot.name.toLowerCase().includes(searchLower) ||
               ot.displayName.toLowerCase().includes(searchLower) ||
               ot.description?.toLowerCase().includes(searchLower);
      }
      return true;
    });
  }, [objectTypes, filter]);
};
```

## 8. Error Handling & Toasts

### 8.1 Toast Notifications (Palantir Style)

```typescript
// Toast configuration
import { Position, Toaster, Intent } from "@blueprintjs/core";

export const AppToaster = Toaster.create({
  position: Position.TOP,
  maxToasts: 3,
});

// Toast utilities
export const toasts = {
  success: (message: string, action?: IActionProps) => {
    AppToaster.show({
      message,
      intent: Intent.SUCCESS,
      icon: "tick",
      timeout: 3000,
      action,
    });
  },
  
  error: (message: string, error?: Error, action?: IActionProps) => {
    AppToaster.show({
      message,
      intent: Intent.DANGER,
      icon: "error",
      timeout: 5000,
      action: action || {
        text: "Retry",
        onClick: () => window.location.reload(),
      },
    });
    
    // Log to error tracking
    console.error(message, error);
  },
  
  warning: (message: string) => {
    AppToaster.show({
      message,
      intent: Intent.WARNING,
      icon: "warning-sign",
      timeout: 4000,
    });
  },
  
  info: (message: string) => {
    AppToaster.show({
      message,
      intent: Intent.PRIMARY,
      icon: "info-sign",
      timeout: 3000,
    });
  },
};

// Usage examples
toasts.success("Object type created successfully");
toasts.error("Failed to save changes", error, {
  text: "Retry",
  onClick: handleRetry,
});
```

### 8.2 Global Error Boundary

```typescript
const ErrorBoundary: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <ErrorBoundaryComponent
      fallback={({ error, resetErrorBoundary }) => (
        <NonIdealState
          icon="error"
          title="Something went wrong"
          description={error.message}
          action={
            <Button intent="primary" onClick={resetErrorBoundary}>
              Try Again
            </Button>
          }
        />
      )}
    >
      {children}
    </ErrorBoundaryComponent>
  );
};
```

### 8.3 Alert Dialogs

```typescript
// Confirmation dialog for destructive actions
const DeleteConfirmation: React.FC<{
  isOpen: boolean;
  itemName: string;
  onConfirm: () => void;
  onCancel: () => void;
}> = ({ isOpen, itemName, onConfirm, onCancel }) => (
  <Alert
    isOpen={isOpen}
    intent={Intent.DANGER}
    icon="trash"
    confirmButtonText="Delete"
    cancelButtonText="Cancel"
    onConfirm={onConfirm}
    onCancel={onCancel}
  >
    <p>
      Are you sure you want to delete <strong>{itemName}</strong>?
    </p>
    <p className="bp5-text-muted">
      This action cannot be undone. All associated data will be permanently removed.
    </p>
  </Alert>
);
```

## 9. Accessibility

### 9.1 ARIA Implementation

```typescript
// Accessible form components
const AccessibleFormGroup: React.FC<FormGroupProps> = ({ 
  label, 
  labelInfo, 
  helperText,
  children,
  error,
  required 
}) => {
  const inputId = useId();
  const helperId = useId();
  const errorId = useId();

  return (
    <FormGroup
      label={label}
      labelFor={inputId}
      labelInfo={labelInfo}
      helperText={error || helperText}
      intent={error ? 'danger' : 'none'}
    >
      {React.cloneElement(children as React.ReactElement, {
        id: inputId,
        'aria-describedby': `${helperId} ${error ? errorId : ''}`,
        'aria-required': required,
        'aria-invalid': !!error,
      })}
    </FormGroup>
  );
};
```

### 9.2 Keyboard Navigation

```typescript
// Custom keyboard shortcuts with Blueprint Hotkeys
import { useHotkeys } from "@blueprintjs/core";

const useKeyboardShortcuts = () => {
  const hotkeys = useMemo(() => [
    {
      combo: "cmd+n",
      global: true,
      label: "Create new object type",
      onKeyDown: () => navigateToCreate(),
    },
    {
      combo: "cmd+f",
      global: true,
      label: "Focus search",
      onKeyDown: () => focusSearch(),
    },
    {
      combo: "cmd+s",
      global: false,
      label: "Save changes",
      onKeyDown: () => handleSave(),
      preventDefault: true,
    },
    {
      combo: "esc",
      global: false,
      label: "Close dialog",
      onKeyDown: () => closeActiveDialog(),
    },
  ], []);

  useHotkeys(hotkeys);
};

// Focus management
const useFocusTrap = (isActive: boolean) => {
  const containerRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (isActive && containerRef.current) {
      const focusableElements = containerRef.current.querySelectorAll(
        'a[href], button, textarea, input, select, [tabindex]:not([tabindex="-1"])'
      );
      const firstElement = focusableElements[0] as HTMLElement;
      firstElement?.focus();
    }
  }, [isActive]);
  
  return containerRef;
};
```

## 10. Testing Strategy

### 10.1 Unit Tests

```typescript
// Component tests
describe('ObjectTypeCard', () => {
  it('should display object type information', () => {
    const objectType = mockObjectType();
    
    render(<ObjectTypeCard objectType={objectType} />);
    
    expect(screen.getByText(objectType.displayName)).toBeInTheDocument();
    expect(screen.getByText(objectType.description)).toBeInTheDocument();
    expect(screen.getByText(`${objectType.properties.length} properties`)).toBeInTheDocument();
  });

  it('should call onEdit when edit button is clicked', () => {
    const onEdit = jest.fn();
    const objectType = mockObjectType();
    
    render(<ObjectTypeCard objectType={objectType} onEdit={onEdit} />);
    
    fireEvent.click(screen.getByText('Edit'));
    expect(onEdit).toHaveBeenCalled();
  });
});
```

### 10.2 Integration Tests

```typescript
// API integration tests
describe('Object Type Creation Flow', () => {
  it('should create a new object type', async () => {
    const { result } = renderHook(() => useObjectTypes());

    await act(async () => {
      await result.current.createObjectType({
        name: 'test_object',
        displayName: 'Test Object',
        properties: [
          {
            name: 'id',
            displayName: 'ID',
            dataType: 'STRING',
            required: true,
          },
        ],
      });
    });

    expect(result.current.objectTypes).toHaveLength(1);
    expect(result.current.objectTypes[0].name).toBe('test_object');
  });
});
```

### 10.3 E2E Tests

```typescript
// Cypress tests
describe('Object Type Management', () => {
  beforeEach(() => {
    cy.login();
    cy.visit('/object-types');
  });

  it('should create a new object type', () => {
    cy.findByRole('button', { name: 'Create Object Type' }).click();
    
    cy.findByLabelText('System Name').type('customer');
    cy.findByLabelText('Display Name').type('Customer');
    cy.findByLabelText('Description').type('Customer information');
    
    cy.findByRole('button', { name: 'Add Property' }).click();
    cy.findByLabelText('Property Name').type('email');
    cy.findByLabelText('Data Type').select('STRING');
    
    cy.findByRole('button', { name: 'Create Object Type' }).click();
    
    cy.findByText('Object type created successfully').should('be.visible');
    cy.findByText('Customer').should('be.visible');
  });
});
```

## 11. Theme and Styling

### 11.1 Blueprint Theme Extension

```scss
// styles/theme.scss
@import "~normalize.css/normalize.css";
@import "~@blueprintjs/core/lib/css/blueprint.css";
@import "~@blueprintjs/icons/lib/css/blueprint-icons.css";
@import "~@blueprintjs/table/lib/css/table.css";
@import "~@blueprintjs/select/lib/css/blueprint-select.css";
@import "~@blueprintjs/popover2/lib/css/blueprint-popover2.css";

// OMS Design System Variables (Palantir Foundry Style)
:root {
  // Intent Colors
  --oms-primary: #2D72D2;        // Blueprint Blue3
  --oms-success: #238551;        // Blueprint Green3
  --oms-warning: #C87619;        // Blueprint Orange3
  --oms-danger: #CD4246;         // Blueprint Red3
  
  // Backgrounds
  --oms-bg-primary: #FFFFFF;
  --oms-bg-secondary: #F6F7F9;
  --oms-bg-tertiary: #EDEFF2;
  
  // Borders
  --oms-border-default: #D3D8DE;
  --oms-border-dark: #ABB3BF;
  
  // Text
  --oms-text-primary: #1C2127;
  --oms-text-secondary: #5F6B7C;
  --oms-text-disabled: #ABB3BF;
}

// Dark theme
.bp5-dark {
  --oms-bg-primary: #1C2127;
  --oms-bg-secondary: #252A31;
  --oms-bg-tertiary: #2F343C;
  
  --oms-text-primary: #F6F7F9;
  --oms-text-secondary: #ABB3BF;
}

// OMS Component Styles
.oms-app {
  background-color: var(--oms-bg-secondary);
  min-height: 100vh;
  
  .bp5-navbar {
    background-color: var(--oms-bg-primary);
    box-shadow: 0 1px 0 var(--oms-border-default);
    height: 50px;
  }
}

// Card Styles
.oms-object-type-card {
  transition: all 200ms cubic-bezier(0.4, 1, 0.75, 0.9);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(16, 22, 26, 0.1), 0 2px 4px rgba(16, 22, 26, 0.04);
  }
  
  .oms-card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;
  }
  
  .oms-card-meta {
    display: flex;
    gap: 16px;
    margin-top: 12px;
    
    .oms-meta-item {
      display: flex;
      align-items: center;
      gap: 4px;
      color: var(--oms-text-secondary);
      font-size: 12px;
    }
  }
  
  .oms-card-actions {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--oms-border-default);
  }
}

// Property Editor
.oms-property-editor {
  padding: 16px;
  background: var(--oms-bg-tertiary);
  border: 1px solid var(--oms-border-default);
  border-radius: 3px;
  margin-bottom: 8px;
  
  &:hover {
    border-color: var(--oms-border-dark);
  }
}

// Explorer Layout
.oms-explorer {
  &__search {
    margin-bottom: 16px;
  }
  
  &__filters {
    display: flex;
    gap: 8px;
    margin-bottom: 24px;
    flex-wrap: wrap;
  }
  
  &__results {
    &.oms-grid {
      display: grid;
      gap: 16px;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    }
  }
}

// Form Actions
.oms-form-actions {
  display: flex;
  gap: 8px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--oms-border-default);
}
```

### 11.2 Responsive Design

```scss
// Responsive breakpoints (aligned with design-system.md)
$oms-breakpoint-mobile: 480px;
$oms-breakpoint-tablet: 768px;
$oms-breakpoint-desktop: 1024px;
$oms-breakpoint-wide: 1440px;

// Spacing system (8px grid)
$oms-spacing-xs: 4px;
$oms-spacing-sm: 8px;
$oms-spacing-md: 16px;
$oms-spacing-lg: 24px;
$oms-spacing-xl: 32px;
$oms-spacing-xxl: 48px;

// Responsive grid
.oms-grid {
  display: grid;
  gap: $oms-spacing-md;
  
  // Mobile-first approach
  grid-template-columns: 1fr;
  
  @media (min-width: $oms-breakpoint-tablet) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (min-width: $oms-breakpoint-desktop) {
    grid-template-columns: repeat(3, 1fr);
  }
  
  @media (min-width: $oms-breakpoint-wide) {
    grid-template-columns: repeat(4, 1fr);
  }
}

// Layout structure
.oms-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  height: 100vh;
  
  &__sidebar {
    background: var(--oms-bg-secondary);
    border-right: 1px solid var(--oms-border-default);
    padding: $oms-spacing-md;
    
    @media (max-width: $oms-breakpoint-tablet) {
      position: fixed;
      left: 0;
      top: 50px; // navbar height
      bottom: 0;
      transform: translateX(-100%);
      transition: transform 300ms cubic-bezier(0.4, 1, 0.75, 0.9);
      z-index: 20;
      
      &--open {
        transform: translateX(0);
        box-shadow: 0 0 0 1px rgba(16, 22, 26, 0.1), 
                    0 2px 8px rgba(16, 22, 26, 0.2), 
                    0 8px 24px rgba(16, 22, 26, 0.2);
      }
    }
  }
  
  &__main {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  &__header {
    padding: $oms-spacing-md;
    border-bottom: 1px solid var(--oms-border-default);
  }
  
  &__content {
    flex: 1;
    overflow-y: auto;
    padding: $oms-spacing-lg;
    
    @media (max-width: $oms-breakpoint-mobile) {
      padding: $oms-spacing-md;
    }
  }
  
  @media (max-width: $oms-breakpoint-tablet) {
    grid-template-columns: 1fr;
  }
}
```

## 12. Development Guidelines

### 12.1 Component Structure

```typescript
// Template for new components
import React, { useState, useCallback, useMemo } from 'react';
import { Button, Card, FormGroup } from '@blueprintjs/core';
import { useObjectTypes } from '@/hooks/useObjectTypes';
import type { ObjectType } from '@/types';

interface ComponentNameProps {
  // Props definition
}

export const ComponentName: React.FC<ComponentNameProps> = ({ 
  // Destructured props
}) => {
  // State
  const [state, setState] = useState();
  
  // Hooks
  const { data, loading, error } = useObjectTypes();
  
  // Callbacks
  const handleAction = useCallback(() => {
    // Implementation
  }, [dependencies]);
  
  // Memoized values
  const computedValue = useMemo(() => {
    // Expensive computation
  }, [dependencies]);
  
  // Early returns
  if (loading) return <Spinner />;
  if (error) return <ErrorState error={error} />;
  
  // Render
  return (
    <div className="component-name">
      {/* Component content */}
    </div>
  );
};
```

### 12.2 TypeScript Best Practices

```typescript
// Strict type definitions
interface ObjectTypeFormData {
  name: string;
  displayName: string;
  description?: string;
  category?: string;
  tags: string[];
  properties: PropertyFormData[];
}

// Type guards
const isObjectType = (value: unknown): value is ObjectType => {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'name' in value &&
    'properties' in value
  );
};

// Generic components
interface DataTableProps<T> {
  data: T[];
  columns: Column<T>[];
  onRowClick?: (item: T) => void;
}

const DataTable = <T extends { id: string }>({ 
  data, 
  columns, 
  onRowClick 
}: DataTableProps<T>) => {
  // Implementation
};
```

### 12.3 Design System Integration

```typescript
// Use design system tokens and components
import { omsTheme } from '@/design-system/theme';
import { ObjectTypeCard } from '@/design-system/components';

// Follow the design-system.md guidelines for:
// - Color usage (Intent colors, grays)
// - Typography (H1-H6, body styles)
// - Spacing (8px grid system)
// - Component patterns
// - Animation timings
// - Icon usage

// Example: Using design tokens
const StyledCard = styled(Card)`
  background: var(--oms-bg-primary);
  border: 1px solid var(--oms-border-default);
  padding: ${omsTheme.spacing.md}px;
  
  &:hover {
    border-color: var(--oms-border-dark);
    transition: ${omsTheme.transitions.default};
  }
`;
```

## 13. Deployment

### 13.1 Build Configuration

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'blueprint': [
            '@blueprintjs/core', 
            '@blueprintjs/table', 
            '@blueprintjs/select',
            '@blueprintjs/icons',
            '@blueprintjs/popover2'
          ],
          'vendor': ['react', 'react-dom', 'react-router-dom'],
          'graphql': ['@apollo/client', 'graphql'],
        },
      },
    },
  },
  server: {
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL,
        changeOrigin: true,
      },
    },
  },
});
```

### 13.2 Docker Configuration

```dockerfile
# Build stage
FROM node:18-alpine as builder

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# Production stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 14. Monitoring and Analytics

### 14.1 Error Tracking

```typescript
// Sentry integration
import * as Sentry from '@sentry/react';

Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
  environment: import.meta.env.MODE,
  integrations: [
    new Sentry.BrowserTracing(),
    new Sentry.Replay(),
  ],
  tracesSampleRate: 0.1,
  replaysSessionSampleRate: 0.1,
});
```

### 14.2 Performance Monitoring

```typescript
// Web Vitals tracking
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';

const sendToAnalytics = ({ name, delta, id }) => {
  // Send to analytics service
  gtag('event', name, {
    event_category: 'Web Vitals',
    event_label: id,
    value: Math.round(name === 'CLS' ? delta * 1000 : delta),
    non_interaction: true,
  });
};

getCLS(sendToAnalytics);
getFID(sendToAnalytics);
getFCP(sendToAnalytics);
getLCP(sendToAnalytics);
getTTFB(sendToAnalytics);
```

## 15. Future Enhancements

### 15.1 Phase 2 Features
- Real-time collaboration with WebSocket
- Advanced property validators UI
- Bulk operations interface
- Import/Export workflows
- Version comparison UI

### 15.2 Phase 3 Features
- Visual schema designer (drag-and-drop)
- Advanced search with Elasticsearch
- Multi-language support
- Plugin system for custom property types
- Mobile-responsive design

## 16. Design System References

### 16.1 Key Design Documents
- **Primary Reference**: `design-system.md` - OMS 전용 디자인 시스템 가이드
- **Blueprint.js Docs**: https://blueprintjs.com/docs/
- **Palantir Foundry Style**: 엔터프라이즈 데이터 애플리케이션 패턴

### 16.2 Component Library
모든 UI 컴포넌트는 `design-system.md`에 정의된 스타일과 패턴을 따라야 합니다:
- Color system (Intent colors, Grayscale)
- Typography (Font sizes, weights, line heights)
- Spacing (8px grid system)
- Component patterns (Cards, Forms, Tables, etc.)
- Interaction patterns (Hover, Loading, Empty states)
- Accessibility guidelines

### 16.3 Implementation Checklist
- [ ] Blueprint.js 5.x 및 모든 필요 패키지 설치
- [ ] Design system 토큰 설정
- [ ] 커스텀 컴포넌트 라이브러리 구축
- [ ] 다크 모드 지원
- [ ] 반응형 레이아웃 구현
- [ ] 접근성 (WCAG 2.0 Level AA) 준수
- [ ] 성능 최적화 (가상 스크롤, 코드 분할)

## 17. MSA Integration Strategy

### 17.1 마이크로서비스 통합 개요

OMS 프론트엔드는 마이크로서비스 아키텍처 환경에서 여러 백엔드 서비스와 효율적으로 통신해야 합니다.

```yaml
프론트엔드 통합 포인트:
  API Gateway:
    - GraphQL Federation Gateway
    - 단일 진입점 제공
    - 인증/인가 중앙화
    
  직접 연결 서비스:
    - OMS Service (메타데이터)
    - OSv2 Service (객체 인스턴스)
    - OSS Service (검색)
    - FOO Service (함수 실행)
```

### 17.2 GraphQL Federation Client

#### 17.2.1 Apollo Client Federation 설정
```typescript
import { ApolloClient, InMemoryCache, split } from '@apollo/client';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import { RetryLink } from '@apollo/client/link/retry';
import { onError } from '@apollo/client/link/error';

// Error handling for distributed services
const errorLink = onError(({ graphQLErrors, networkError, operation, forward }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path, extensions }) => {
      // Service-specific error handling
      const serviceName = extensions?.serviceName;
      
      switch (serviceName) {
        case 'oms-service':
          handleOMSError(message, extensions);
          break;
        case 'osv2-service':
          handleOSv2Error(message, extensions);
          break;
        default:
          toasts.error(`Error in ${serviceName}: ${message}`);
      }
    });
  }

  if (networkError) {
    // Circuit breaker pattern on client
    if (networkError.statusCode === 503) {
      toasts.warning('Service temporarily unavailable. Retrying...');
      return forward(operation);
    }
  }
});

// Retry logic for resilience
const retryLink = new RetryLink({
  delay: {
    initial: 300,
    max: Infinity,
    jitter: true
  },
  attempts: {
    max: 3,
    retryIf: (error, _operation) => {
      return !!error && error.statusCode === 503;
    }
  }
});

// WebSocket for real-time updates
const wsLink = new WebSocketLink({
  uri: import.meta.env.VITE_WS_ENDPOINT,
  options: {
    reconnect: true,
    connectionParams: () => ({
      authorization: `Bearer ${getAuthToken()}`
    })
  }
});

// Split between WebSocket and HTTP
const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink
);

// Apollo Client with federation support
export const apolloClient = new ApolloClient({
  link: ApolloLink.from([errorLink, retryLink, splitLink]),
  cache: new InMemoryCache({
    typePolicies: {
      ObjectType: {
        keyFields: ['id'],
        fields: {
          // Federation field merging
          instances: {
            keyArgs: ['filter'],
            merge(existing = { edges: [] }, incoming) {
              return {
                ...incoming,
                edges: [...existing.edges, ...incoming.edges]
              };
            }
          }
        }
      }
    }
  })
});
```

#### 17.2.2 Federated Queries
```typescript
// Cross-service query example
const OBJECT_WITH_INSTANCES_QUERY = gql`
  query GetObjectTypeWithInstances($objectTypeId: ID!, $instanceFilter: InstanceFilter) {
    # From OMS service
    objectType(id: $objectTypeId) {
      id
      name
      displayName
      properties {
        id
        name
        dataType
      }
      
      # From OSv2 service (federated)
      instances(filter: $instanceFilter) @connection(key: "instances") {
        edges {
          node {
            id
            data
            createdAt
            updatedAt
          }
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
      
      # From OSS service (federated)
      searchMetrics @client {
        totalCount
        lastSearched
      }
    }
  }
`;
```

### 17.3 Service Health Monitoring

#### 17.3.1 Service Status Component
```typescript
interface ServiceHealth {
  name: string;
  status: 'healthy' | 'degraded' | 'down';
  latency: number;
  lastChecked: Date;
}

const ServiceHealthIndicator: React.FC = () => {
  const [services, setServices] = useState<ServiceHealth[]>([]);
  
  useEffect(() => {
    const checkHealth = async () => {
      const healthChecks = await Promise.allSettled([
        checkServiceHealth('oms-service'),
        checkServiceHealth('osv2-service'),
        checkServiceHealth('oss-service'),
        checkServiceHealth('foo-service')
      ]);
      
      const serviceStatuses = healthChecks.map((result, index) => {
        const serviceName = ['oms', 'osv2', 'oss', 'foo'][index];
        
        if (result.status === 'fulfilled') {
          return result.value;
        } else {
          return {
            name: serviceName,
            status: 'down' as const,
            latency: -1,
            lastChecked: new Date()
          };
        }
      });
      
      setServices(serviceStatuses);
    };
    
    // Initial check
    checkHealth();
    
    // Periodic health checks
    const interval = setInterval(checkHealth, 30000); // 30 seconds
    
    return () => clearInterval(interval);
  }, []);
  
  const overallHealth = services.every(s => s.status === 'healthy') 
    ? 'healthy' 
    : services.some(s => s.status === 'down') 
      ? 'degraded' 
      : 'warning';
  
  return (
    <Popover
      content={
        <Menu>
          <MenuDivider title="Service Health" />
          {services.map(service => (
            <MenuItem
              key={service.name}
              icon={getStatusIcon(service.status)}
              text={service.name.toUpperCase()}
              label={`${service.latency}ms`}
              intent={getStatusIntent(service.status)}
            />
          ))}
        </Menu>
      }
    >
      <Button
        minimal
        icon={getStatusIcon(overallHealth)}
        intent={getStatusIntent(overallHealth)}
      />
    </Popover>
  );
};
```

### 17.4 Real-time Updates

#### 17.4.1 WebSocket Subscriptions
```typescript
// Real-time object type updates
const OBJECT_TYPE_SUBSCRIPTION = gql`
  subscription OnObjectTypeChange($objectTypeIds: [ID!]) {
    objectTypeChanged(objectTypeIds: $objectTypeIds) {
      id
      type # CREATED | UPDATED | DELETED
      objectType {
        id
        name
        displayName
        version
        updatedAt
        updatedBy
      }
    }
  }
`;

// Hook for real-time updates
export const useObjectTypeSubscription = (objectTypeIds: string[]) => {
  const queryClient = useQueryClient();
  
  useSubscription(OBJECT_TYPE_SUBSCRIPTION, {
    variables: { objectTypeIds },
    onSubscriptionData: ({ subscriptionData }) => {
      const { type, objectType } = subscriptionData.data.objectTypeChanged;
      
      switch (type) {
        case 'UPDATED':
          // Update cache
          queryClient.setQueryData(
            ['objectType', objectType.id],
            objectType
          );
          
          // Show notification
          toasts.info(
            `${objectType.displayName} was updated by ${objectType.updatedBy}`,
            {
              action: {
                text: 'Refresh',
                onClick: () => window.location.reload()
              }
            }
          );
          break;
          
        case 'DELETED':
          // Remove from cache
          queryClient.removeQueries(['objectType', objectType.id]);
          
          // Redirect if viewing deleted object
          if (window.location.pathname.includes(objectType.id)) {
            navigate('/object-types');
            toasts.warning(`${objectType.displayName} was deleted`);
          }
          break;
      }
    }
  });
};
```

### 17.5 Distributed Error Handling

#### 17.5.1 Service-Specific Error Handlers
```typescript
// Error context for distributed systems
interface ServiceError {
  service: string;
  code: string;
  message: string;
  details?: Record<string, any>;
  traceId?: string;
}

const ServiceErrorBoundary: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <ErrorBoundary
      fallback={({ error, resetErrorBoundary }) => {
        const serviceError = parseServiceError(error);
        
        return (
          <NonIdealState
            icon="error"
            title="Service Error"
            description={
              <>
                <p>{serviceError.message}</p>
                {serviceError.traceId && (
                  <p className="bp5-text-muted">
                    Trace ID: <code>{serviceError.traceId}</code>
                  </p>
                )}
              </>
            }
            action={
              <ButtonGroup>
                <Button intent="primary" onClick={resetErrorBoundary}>
                  Retry
                </Button>
                <Button onClick={() => copyTraceId(serviceError.traceId)}>
                  Copy Trace ID
                </Button>
              </ButtonGroup>
            }
          />
        );
      }}
    >
      {children}
    </ErrorBoundary>
  );
};

// Graceful degradation for service failures
const useServiceWithFallback = <T,>(
  serviceName: string,
  query: DocumentNode,
  variables?: any
) => {
  const [fallbackMode, setFallbackMode] = useState(false);
  
  const { data, loading, error } = useQuery(query, {
    variables,
    errorPolicy: 'all',
    onError: (error) => {
      if (isServiceUnavailable(error)) {
        setFallbackMode(true);
        // Use cached data or local storage
        return getCachedData(serviceName, variables);
      }
    }
  });
  
  return {
    data: fallbackMode ? getCachedData(serviceName, variables) : data,
    loading,
    error,
    fallbackMode
  };
};
```

### 17.6 Performance Optimization for MSA

#### 17.6.1 Request Batching
```typescript
// Batch multiple service requests
import { BatchHttpLink } from '@apollo/client/link/batch-http';

const batchLink = new BatchHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_ENDPOINT,
  batchMax: 5, // Max 5 queries per batch
  batchInterval: 20 // 20ms window
});
```

#### 17.6.2 Response Caching Strategy
```typescript
// Service-specific cache policies
const cacheConfig = {
  typePolicies: {
    // OMS data - longer cache
    ObjectType: {
      fields: {
        properties: {
          merge: false // Replace on update
        }
      }
    },
    // OSv2 data - shorter cache
    ObjectInstance: {
      fields: {
        data: {
          read(existing, { args, toReference }) {
            // Check freshness
            if (existing && Date.now() - existing.timestamp < 60000) {
              return existing.value;
            }
            // Trigger refetch
            return undefined;
          }
        }
      }
    }
  }
};
```

### 17.7 Development Tools for MSA

#### 17.7.1 Service Mocking
```typescript
// Mock service responses for development
const mocks = [
  {
    request: {
      query: OBJECT_TYPES_QUERY,
      variables: {}
    },
    result: {
      data: {
        objectTypes: {
          edges: mockObjectTypes,
          pageInfo: { hasNextPage: false, endCursor: null }
        }
      }
    },
    delay: 100 // Simulate network latency
  }
];

// Apollo MockedProvider for tests
<MockedProvider mocks={mocks} addTypename={false}>
  <App />
</MockedProvider>
```

#### 17.7.2 Distributed Tracing UI
```typescript
const TraceViewer: React.FC<{ traceId: string }> = ({ traceId }) => {
  const { data, loading } = useQuery(GET_TRACE_QUERY, {
    variables: { traceId }
  });
  
  if (loading) return <Spinner />;
  
  return (
    <div className="trace-viewer">
      <H4>Request Trace: {traceId}</H4>
      <Tree
        contents={buildTraceTree(data.trace)}
        onNodeClick={handleNodeClick}
      />
      <div className="trace-timeline">
        {data.trace.spans.map(span => (
          <div 
            key={span.id}
            className="span-bar"
            style={{
              left: `${span.startTime}%`,
              width: `${span.duration}%`,
              backgroundColor: getServiceColor(span.service)
            }}
          >
            <Tooltip content={`${span.service}: ${span.duration}ms`}>
              <div>{span.operation}</div>
            </Tooltip>
          </div>
        ))}
      </div>
    </div>
  );
};
```

## 18. Conclusion

이 프론트엔드 명세서는 OMS PRD의 모든 요구사항을 Blueprint.js와 Palantir Foundry 스타일의 디자인 시스템을 활용하여 구현하는 방법을 상세히 정의합니다. 특히 MSA 환경에서의 효율적인 통합 전략을 포함하여, `design-system.md`와 함께 참조하여 일관되고 전문적인 엔터프라이즈 애플리케이션을 구축할 수 있습니다.
