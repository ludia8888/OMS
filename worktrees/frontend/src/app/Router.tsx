import { Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from '../shared/components/Layout/MainLayout';
import { ObjectTypesPage } from '../features/object-types/pages/ObjectTypesPage';
import { LinkTypesPage } from '../features/link-types/pages/LinkTypesPage';
import { SearchPage } from '../features/search/pages/SearchPage';

export function Router() {
  return (
    <Routes>
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Navigate to="/object-types" replace />} />
        <Route path="object-types" element={<ObjectTypesPage />} />
        <Route path="link-types" element={<LinkTypesPage />} />
        <Route path="search" element={<SearchPage />} />
      </Route>
    </Routes>
  );
}