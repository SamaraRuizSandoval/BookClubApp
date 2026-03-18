import {
  IonMenuToggle,
  IonPage,
  IonRouterOutlet,
  IonHeader,
  IonSplitPane,
} from '@ionic/react';
import { Switch, Route } from 'react-router-dom';
import { useLocation } from 'react-router-dom';

import { UserStatsProvider } from '../../context/UserStatsContext';
import { NotFound } from '../../pages/Not found';
import { Page } from '../../pages/Page';
import { DiscoverBooks } from '../../pages/user/DiscoverBooks';
import { MyShelf } from '../../pages/user/MyShelf';
import {
  ReadingSection,
  WishlistSection,
  CompletedSection,
} from '../../utils/sections';
import { LeftMenu } from '../LeftMenu';
import { UpperNavigation } from '../UpperNavigation';

export function UserLayout() {
  const location = useLocation();

  const renderNavItem = (to: string, icon: any, label: string) => {
    const active = location.pathname.startsWith(to);
    return (
      <IonMenuToggle key={to} autoHide={false}>
        <a className={`nav-item ${active ? 'active' : ''}`} href={to}>
          <span className="nav-icon">{icon}</span>
          {label}
        </a>
      </IonMenuToggle>
    );
  };
  return (
    <IonSplitPane when="lg" contentId="main-content">
      <UserStatsProvider>
        <LeftMenu />

        <IonPage id="main-content">
          <IonHeader>
            <UpperNavigation />
          </IonHeader>

          <IonRouterOutlet>
            <Switch>
              <Route exact path="/app" component={DiscoverBooks} />
              <Route path="/app/library" component={MyShelf} />
              <Route path="/app/library/reading" component={ReadingSection} />
              <Route path="/app/library/wishlist" component={WishlistSection} />
              <Route
                path="/app/library/completed"
                component={CompletedSection}
              />
              <Route
                path="/app/settings"
                render={() => <Page title="Settings" />}
              />
              <Route component={NotFound} />
            </Switch>
          </IonRouterOutlet>
        </IonPage>
      </UserStatsProvider>
    </IonSplitPane>
  );
}
