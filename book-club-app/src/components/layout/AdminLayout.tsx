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
import { bookOutline } from 'ionicons/icons';
import { Route, Redirect, Switch } from 'react-router-dom';

import { AdminBooks } from '../../pages/admin/ManageBooksPage';
import { AddBookPage } from '../../pages/admin/books/AddBookPage';
import { SearchGoogleBooks } from '../../pages/admin/books/SearchGoogleBooks';

import { UpperNavigation } from './UpperNavigation';

export function AdminLayout() {
  return (
    <IonSplitPane when="lg" contentId="admin-content">
      <IonMenu contentId="admin-content" type="overlay" side="start">
        <IonHeader>
          <IonToolbar>
            <IonTitle>Admin Panel</IonTitle>
          </IonToolbar>
        </IonHeader>

        <IonContent fullscreen className="ion-padding">
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

        <IonRouterOutlet id="admin-content">
          <Switch>
            <Route path="/admin/books" component={AdminBooks} exact />
            <Route
              path="/admin/search-google-books"
              component={SearchGoogleBooks}
              exact
            />
            <Route path="/admin/books/add-book" component={AddBookPage} exact />
            <Redirect exact from="/admin" to="/admin/books" />
          </Switch>
        </IonRouterOutlet>
      </IonPage>
    </IonSplitPane>
  );
}
