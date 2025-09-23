import {
  IonApp,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonPage,
  IonMenu,
  IonList,
  IonItem,
  IonRouterOutlet,
  IonMenuToggle,
  IonButtons,
  IonMenuButton,
  IonIcon,
  IonSplitPane,
  IonLabel,
} from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import { homeOutline, bookmarkOutline, settingsOutline } from 'ionicons/icons';
import React from 'react';
import { createRoot } from 'react-dom/client';
import { Route, Switch, Redirect, useLocation } from 'react-router-dom';

function Page({ title }: { title: string }) {
  const location = useLocation();
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonButtons slot="start">
            <IonMenuButton />
          </IonButtons>
          <IonTitle>{title}</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent className="ion-padding">
        <h1 className="text-2xl font-semibold">{title}</h1>
        <p className="opacity-80">Current route: {location.pathname}</p>
      </IonContent>
    </IonPage>
  );
}

const menuItems: Array<{
  to: string;
  icon: string;
  label: string;
}> = [
  { to: '/home', icon: homeOutline, label: 'Home' },
  { to: '/saved', icon: bookmarkOutline, label: 'Saved' },
  { to: '/settings', icon: settingsOutline, label: 'Settings' },
];

function LeftMenu() {
  const location = useLocation();

  return (
    <IonMenu contentId="main-content" type="overlay">
      <IonHeader>
        <IonToolbar>
          <IonTitle>BookClub</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent>
        <IonList inset>
          {menuItems.map(({ to, icon, label }) => {
            const active = location.pathname.startsWith(to);
            return (
              <IonMenuToggle key={to} autoHide={false}>
                <IonItem
                  button
                  detail={false}
                  color={active ? 'primary' : undefined}
                  routerLink={to}
                  routerDirection="root"
                >
                  <IonIcon slot="start" icon={icon} />
                  <IonLabel>{label}</IonLabel>
                </IonItem>
              </IonMenuToggle>
            );
          })}
        </IonList>
      </IonContent>
    </IonMenu>
  );
}

function AppShell() {
  return (
    <IonReactRouter>
      <IonSplitPane when="md" contentId="main-content">
        <LeftMenu />
        <IonRouterOutlet id="main-content">
          <Switch>
            <Route path="/home" component={() => <Page title="Home" />} />
            <Route path="/saved" component={() => <Page title="Saved" />} />
            <Route
              path="/settings"
              component={() => <Page title="Settings" />}
            />
            <Route exact path="/" render={() => <Redirect to="/home" />} />
            <Route component={() => <Page title="Not Found" />} />
          </Switch>
        </IonRouterOutlet>
      </IonSplitPane>
    </IonReactRouter>
  );
}

export default function App() {
  return (
    <IonApp>
      <AppShell />
    </IonApp>
  );
}
