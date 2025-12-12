import { IonApp, IonSplitPane, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import React from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';

import { LeftMenu } from './components/LeftMenu';
import { Page } from './pages/Page';
import {
  ReadingSection,
  WishlistSection,
  CompletedSection,
} from './utils/sections';

export default function App() {
  return (
    <IonApp>
      <IonReactRouter>
        <IonSplitPane when="md" contentId="main-content">
          <LeftMenu />
          <IonRouterOutlet id="main-content">
            <Switch>
              <Route path="/home" component={() => <Page title="Home" />} />
              <Route path="/reading" component={ReadingSection} />
              <Route path="/wishlist" component={WishlistSection} />
              <Route path="/completed" component={CompletedSection} />
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
    </IonApp>
  );
}
