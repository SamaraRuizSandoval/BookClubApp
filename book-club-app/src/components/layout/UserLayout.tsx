import {
  IonContent,
  IonMenu,
  IonMenuToggle,
  IonPage,
  IonRouterOutlet,
  IonTitle,
  IonToolbar,
  IonHeader,
  IonSplitPane,
} from '@ionic/react';
import { Switch, Route } from 'react-router-dom';
import { useLocation } from 'react-router-dom';

import { Page } from '../../pages/Page';
import {
  ReadingSection,
  WishlistSection,
  CompletedSection,
} from '../../utils/sections';
import { LeftMenu } from '../LeftMenu';

import { UpperNavigation } from './UpperNavigation';

const navItems = [
  { to: '/home', icon: '🔍', label: 'Browse Books' },
  { to: '/my-shelf', icon: '🏠', label: 'My Shelf' },
];

const bookCollections = [
  {
    to: '/reading',
    style: { background: 'var(--reading)' },
    label: 'Reading',
    id: 'count-reading',
  },
  {
    to: '/wishlist',
    style: { background: 'var(--want)' },
    label: 'Wishlist',
    id: 'count-want',
  },
  {
    to: '/completed',
    style: { background: 'var(--done)' },
    label: 'Completed',
    id: 'count-done',
  },
];

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
      <LeftMenu />

      <IonPage id="main-content">
        <IonHeader>
          <UpperNavigation />
        </IonHeader>

        <IonRouterOutlet>
          <Switch>
            <Route exact path="/home" render={() => <Page title="Home" />} />
            <Route path="/reading" component={ReadingSection} />
            <Route path="/wishlist" component={WishlistSection} />
            <Route path="/completed" component={CompletedSection} />
            <Route path="/settings" render={() => <Page title="Settings" />} />
            <Route render={() => <Page title="Not Found" />} />
          </Switch>
        </IonRouterOutlet>
      </IonPage>
    </IonSplitPane>
  );
}
