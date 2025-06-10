import { Outlet, Link, useLocation } from 'react-router-dom';
import {
  Navbar,
  NavbarGroup,
  NavbarHeading,
  NavbarDivider,
  Button,
  Alignment,
  Classes,
} from '@blueprintjs/core';

export function MainLayout() {
  const location = useLocation();

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <div className="app-layout">
      <Navbar>
        <NavbarGroup align={Alignment.LEFT}>
          <NavbarHeading>
            <strong>OMS</strong>
          </NavbarHeading>
          <NavbarDivider />
          <Link to="/object-types" className={Classes.MINIMAL}>
            <Button
              className={Classes.MINIMAL}
              icon="cube"
              text="Object Types"
              active={isActive('/object-types')}
            />
          </Link>
          <Link to="/link-types" className={Classes.MINIMAL}>
            <Button
              className={Classes.MINIMAL}
              icon="link"
              text="Link Types"
              active={isActive('/link-types')}
            />
          </Link>
          <Link to="/search" className={Classes.MINIMAL}>
            <Button
              className={Classes.MINIMAL}
              icon="search"
              text="Search"
              active={isActive('/search')}
            />
          </Link>
        </NavbarGroup>
        <NavbarGroup align={Alignment.RIGHT}>
          <Button className={Classes.MINIMAL} icon="cog" />
          <NavbarDivider />
          <Button className={Classes.MINIMAL} icon="user" />
        </NavbarGroup>
      </Navbar>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}