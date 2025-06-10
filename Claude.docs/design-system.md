# OMS (Ontology Metadata Service) 디자인 시스템

## 개요

OMS는 Palantir Foundry의 Blueprint 디자인 시스템을 기반으로 구축됩니다. 이 문서는 OMS의 UI/UX 구현을 위한 디자인 가이드라인을 제공합니다.

## 1. 디자인 철학

### 핵심 원칙
- **데이터 중심 인터페이스**: 복잡한 메타데이터 관리를 위한 고밀도 정보 표시
- **개발자 친화적**: 기술 사용자를 위한 직관적이고 효율적인 UI
- **일관성**: OpenFoundry 플랫폼 전체와 통일된 사용자 경험
- **확장성**: 향후 기능 추가를 고려한 모듈형 디자인

### 디자인 목표
- 복잡한 객체 타입 관계를 명확하게 시각화
- 빠른 탐색과 검색을 위한 효율적인 정보 구조
- 실시간 협업과 버전 관리를 위한 직관적인 인터페이스

## 2. 색상 시스템

### 주요 색상 (Intent Colors)
```scss
// OMS 브랜드 색상 (Blueprint 기반)
$oms-primary: #2D72D2;        // Blueprint Blue3 - 주요 액션, 링크
$oms-success: #238551;        // Blueprint Green3 - 성공 상태
$oms-warning: #C87619;        // Blueprint Orange3 - 경고
$oms-danger: #CD4246;         // Blueprint Red3 - 오류, 삭제

// 의미론적 색상
$oms-info: #2D72D2;           // 정보성 메시지
$oms-muted: #5F6B7C;          // 비활성 텍스트
```

### 그레이스케일
```scss
// 배경색
$oms-bg-primary: #FFFFFF;      // 메인 배경
$oms-bg-secondary: #F6F7F9;    // 카드, 섹션 배경
$oms-bg-tertiary: #EDEFF2;     // 테이블 스트라이프

// 테두리
$oms-border-default: #D3D8DE;  // 기본 테두리
$oms-border-dark: #ABB3BF;     // 강조 테두리

// 텍스트
$oms-text-primary: #1C2127;    // 본문
$oms-text-secondary: #5F6B7C;  // 보조 텍스트
$oms-text-disabled: #ABB3BF;   // 비활성 텍스트
```

### 다크 모드
```scss
// 다크 모드 배경
$oms-dark-bg-primary: #1C2127;
$oms-dark-bg-secondary: #252A31;
$oms-dark-bg-tertiary: #2F343C;

// 다크 모드 텍스트
$oms-dark-text-primary: #F6F7F9;
$oms-dark-text-secondary: #ABB3BF;
```

## 3. Typography

### 폰트 스택
```css
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, 
             "Helvetica Neue", Arial, sans-serif;

/* 코드/기술적 콘텐츠용 */
font-family-monospace: "SF Mono", Monaco, Consolas, "Courier New", monospace;
```

### 텍스트 스타일
```scss
// 제목
.oms-h1 { font-size: 36px; line-height: 40px; font-weight: 700; }
.oms-h2 { font-size: 28px; line-height: 32px; font-weight: 700; }
.oms-h3 { font-size: 22px; line-height: 25px; font-weight: 600; }
.oms-h4 { font-size: 18px; line-height: 21px; font-weight: 600; }

// 본문
.oms-body { font-size: 14px; line-height: 20px; }
.oms-body-large { font-size: 16px; line-height: 22px; }
.oms-body-small { font-size: 12px; line-height: 16px; }

// 특수 용도
.oms-code { font-family: $font-family-monospace; font-size: 13px; }
.oms-label { font-size: 12px; font-weight: 600; text-transform: uppercase; }
```

## 4. 컴포넌트 스타일

### 버튼
```tsx
// Primary Action
<Button intent="primary" icon="plus">
  Create Object Type
</Button>

// Secondary Action
<Button intent="none" outlined>
  Cancel
</Button>

// Danger Action
<Button intent="danger" icon="trash">
  Delete
</Button>

// Icon Button
<Button minimal icon="more" />
```

### 폼 요소
```tsx
// Text Input
<FormGroup label="Object Type Name" labelFor="name-input" intent="primary">
  <InputGroup 
    id="name-input"
    placeholder="e.g., Customer"
    leftIcon="tag"
  />
</FormGroup>

// Select
<Select
  items={dataTypes}
  itemRenderer={renderDataType}
  onItemSelect={handleSelect}
  filterable={true}
>
  <Button text="Select Type" rightIcon="caret-down" />
</Select>

// Multi-select (Tags)
<MultiSelect
  items={availableTags}
  selectedItems={selectedTags}
  placeholder="Add tags..."
/>
```

### 데이터 테이블
```tsx
<HTMLTable striped interactive>
  <thead>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Required</th>
      <th>Actions</th>
    </tr>
  </thead>
  <tbody>
    {properties.map(prop => (
      <tr key={prop.id}>
        <td><Code>{prop.name}</Code></td>
        <td><Tag>{prop.dataType}</Tag></td>
        <td>{prop.required && <Icon icon="tick" intent="success" />}</td>
        <td>
          <ButtonGroup minimal>
            <Button icon="edit" />
            <Button icon="trash" intent="danger" />
          </ButtonGroup>
        </td>
      </tr>
    ))}
  </tbody>
</HTMLTable>
```

### 카드와 패널
```tsx
// Object Type Card
<Card interactive elevation={1} className="oms-object-type-card">
  <div className="oms-card-header">
    <H4>{objectType.displayName}</H4>
    <Tag>{objectType.category}</Tag>
  </div>
  <p className="bp4-text-muted">{objectType.description}</p>
  <div className="oms-card-meta">
    <Icon icon="properties" />
    <span>{objectType.properties.length} properties</span>
  </div>
</Card>
```

## 5. 레이아웃 패턴

### 페이지 구조
```scss
.oms-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  height: 100vh;
  
  &__sidebar {
    background: $oms-bg-secondary;
    border-right: 1px solid $oms-border-default;
    padding: $pt-grid-size * 2;
  }
  
  &__main {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  &__header {
    padding: $pt-grid-size * 2;
    border-bottom: 1px solid $oms-border-default;
  }
  
  &__content {
    flex: 1;
    overflow-y: auto;
    padding: $pt-grid-size * 3;
  }
}
```

### 간격 시스템
```scss
// 기본 단위: 8px
$oms-spacing-xs: 4px;
$oms-spacing-sm: 8px;
$oms-spacing-md: 16px;
$oms-spacing-lg: 24px;
$oms-spacing-xl: 32px;
$oms-spacing-xxl: 48px;
```

## 6. 특수 컴포넌트

### Object Type Explorer
```tsx
<div className="oms-explorer">
  <div className="oms-explorer__search">
    <InputGroup
      large
      leftIcon="search"
      placeholder="Search object types..."
      rightElement={
        <Button minimal icon="filter-list">Filters</Button>
      }
    />
  </div>
  
  <div className="oms-explorer__filters">
    <Tag>Category: All</Tag>
    <Tag>Status: Active</Tag>
  </div>
  
  <div className="oms-explorer__results">
    {/* Object Type Cards Grid */}
  </div>
</div>
```

### Property Editor
```tsx
<div className="oms-property-editor">
  <ControlGroup fill>
    <InputGroup placeholder="Property name" />
    <HTMLSelect options={dataTypes} />
    <Switch label="Required" />
    <Button icon="tick" intent="primary" />
  </ControlGroup>
</div>
```

### Version Timeline
```tsx
<div className="oms-timeline">
  <Timeline>
    <TimelineEvent 
      icon="git-commit"
      intent="primary"
      title="Version 3"
      subtitle="Current"
    />
    <TimelineEvent 
      icon="git-commit"
      title="Version 2"
      subtitle="2 hours ago by John Doe"
    />
  </Timeline>
</div>
```

## 7. 상태 표시

### 로딩 상태
```tsx
// 전체 페이지 로딩
<NonIdealState
  icon={<Spinner size={SpinnerSize.LARGE} />}
  title="Loading Object Types..."
/>

// 인라인 로딩
<Button loading={isLoading}>Save Changes</Button>

// 스켈레톤 로딩
<Card>
  <H4><Skeleton /></H4>
  <p><Skeleton /></p>
  <Skeleton />
</Card>
```

### 빈 상태
```tsx
<NonIdealState
  icon="folder-open"
  title="No Object Types"
  description="Create your first object type to get started."
  action={<Button intent="primary" icon="plus">Create Object Type</Button>}
/>
```

### 오류 상태
```tsx
<Callout intent="danger" icon="error">
  <H4>Error Loading Data</H4>
  <p>Failed to fetch object types. Please try again.</p>
  <Button intent="primary" outlined>Retry</Button>
</Callout>
```

## 8. 인터랙션 패턴

### 토스트 알림
```tsx
// 성공
AppToaster.show({
  message: "Object type created successfully",
  intent: "success",
  icon: "tick"
});

// 오류
AppToaster.show({
  message: "Failed to save changes",
  intent: "danger",
  icon: "error",
  action: {
    text: "Retry",
    onClick: handleRetry
  }
});
```

### 확인 다이얼로그
```tsx
<Alert
  isOpen={isDeleteOpen}
  intent="danger"
  icon="trash"
  confirmButtonText="Delete"
  cancelButtonText="Cancel"
  onConfirm={handleDelete}
>
  <p>Are you sure you want to delete this object type?</p>
  <p>This action cannot be undone.</p>
</Alert>
```

### 드래그 앤 드롭
```scss
.oms-draggable {
  cursor: move;
  
  &--dragging {
    opacity: 0.5;
  }
  
  &--over {
    background-color: $oms-primary-light;
    border: 2px dashed $oms-primary;
  }
}
```

## 9. 반응형 디자인

### 브레이크포인트
```scss
$oms-breakpoint-mobile: 480px;
$oms-breakpoint-tablet: 768px;
$oms-breakpoint-desktop: 1024px;
$oms-breakpoint-wide: 1440px;
```

### 반응형 그리드
```scss
.oms-grid {
  display: grid;
  gap: $oms-spacing-md;
  
  // 모바일
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
```

## 10. 접근성

### 키보드 네비게이션
- Tab 키로 모든 인터랙티브 요소 접근 가능
- Enter/Space로 버튼 활성화
- Escape로 모달/팝오버 닫기
- 화살표 키로 메뉴 네비게이션

### ARIA 레이블
```tsx
<Button
  icon="filter"
  minimal
  aria-label="Filter object types"
/>

<InputGroup
  aria-label="Search object types"
  placeholder="Search..."
/>
```

### 색상 대비
- 텍스트: 최소 4.5:1 대비율
- 대형 텍스트: 최소 3:1 대비율
- 인터랙티브 요소: 최소 3:1 대비율

## 11. 애니메이션

### 트랜지션
```scss
// 기본 트랜지션
$oms-transition-default: all 200ms cubic-bezier(0.4, 1, 0.75, 0.9);

// 빠른 트랜지션 (호버, 포커스)
$oms-transition-fast: all 100ms cubic-bezier(0.4, 1, 0.75, 0.9);

// 느린 트랜지션 (페이지 전환)
$oms-transition-slow: all 300ms cubic-bezier(0.4, 1, 0.75, 0.9);
```

### 애니메이션 패턴
```scss
@keyframes oms-fade-in {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes oms-slide-in {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}

.oms-animate-fade {
  animation: oms-fade-in 200ms ease-out;
}
```

## 12. 아이콘 사용 가이드

### 주요 아이콘 매핑
```tsx
// Object Types
"cube": 객체 타입
"tag": 속성
"link": 관계/링크
"folder-open": 카테고리

// Actions
"plus": 추가
"edit": 편집
"trash": 삭제
"duplicate": 복제
"export": 내보내기
"import": 가져오기

// States
"tick": 성공/완료
"cross": 오류/닫기
"warning-sign": 경고
"info-sign": 정보

// Navigation
"chevron-right": 확장
"chevron-down": 축소
"arrow-left": 뒤로
"home": 홈
```

## 13. 구현 예제

### 설치
```bash
npm install @blueprintjs/core @blueprintjs/icons @blueprintjs/select
npm install @blueprintjs/popover2 @blueprintjs/table
```

### 테마 설정
```tsx
// theme.ts
export const omsTheme = {
  colors: {
    primary: '#2D72D2',
    success: '#238551',
    warning: '#C87619',
    danger: '#CD4246',
    // ...
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
  },
  typography: {
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI"...',
    // ...
  }
};
```

### 스타일 import
```scss
// styles/main.scss
@import "~normalize.css/normalize.css";
@import "~@blueprintjs/core/lib/css/blueprint.css";
@import "~@blueprintjs/icons/lib/css/blueprint-icons.css";
@import "~@blueprintjs/select/lib/css/blueprint-select.css";
@import "~@blueprintjs/table/lib/css/table.css";

// OMS 커스텀 스타일
@import "variables";
@import "components";
@import "layouts";
@import "utilities";
```

이 디자인 시스템은 OMS가 Palantir Foundry의 검증된 UX 패턴을 따르면서도 메타데이터 관리라는 특수한 목적에 최적화된 인터페이스를 제공할 수 있도록 합니다.