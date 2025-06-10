import { H1, Card, Elevation, InputGroup } from '@blueprintjs/core';

export function SearchPage() {
  return (
    <div className="page-container">
      <H1>Search</H1>
      <Card elevation={Elevation.ONE} style={{ marginTop: 20 }}>
        <InputGroup
          large
          leftIcon="search"
          placeholder="Search for objects, links, or types..."
          style={{ marginBottom: 20 }}
        />
        <p>Advanced search functionality - 개발 예정</p>
      </Card>
    </div>
  );
}