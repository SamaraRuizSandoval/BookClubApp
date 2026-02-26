import {
  IonContent,
  IonItem,
  IonLabel,
  IonList,
  IonMenu,
  IonMenuToggle,
  IonPage,
  IonRouterOutlet,
  IonTitle,
  IonToolbar,
  IonHeader,
  IonIcon,
  IonSplitPane,
} from '@ionic/react';
import { bookOutline, peopleOutline, gridOutline } from 'ionicons/icons';
import { Route, Redirect, Switch } from 'react-router-dom';

import { UpperNavigation } from '../../components/UpperNavigation';

import { AdminBooks } from './AdminBooks';

export function AdminLayout() {
  return (
    <>
      <IonSplitPane when="lg" contentId="admin-content">
        <IonMenu contentId="admin-content" type="overlay" side="start">
          <IonHeader>
            <IonToolbar>
              <IonTitle>Admin Panel</IonTitle>
            </IonToolbar>
          </IonHeader>
          <IonContent>
            <IonList>
              <IonMenuToggle autoHide={false}>
                <IonItem routerLink="/admin/books" routerDirection="none">
                  <IonIcon icon={bookOutline} slot="start" />
                  <IonLabel>Manage Books</IonLabel>
                </IonItem>
              </IonMenuToggle>
            </IonList>
          </IonContent>
        </IonMenu>

        <IonPage id="admin-content">
          <IonHeader>
            <UpperNavigation />
          </IonHeader>
          <IonRouterOutlet>
            <Switch>
              <Route path="/admin/books" component={AdminBooks} exact />
              <Redirect exact from="/admin" to="/admin/books" />
            </Switch>
            {/* <Route path="/admin/dashboard" component={AdminDashboard} exact /> */}
          </IonRouterOutlet>
        </IonPage>
      </IonSplitPane>
    </>
  );
}
