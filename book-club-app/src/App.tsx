import { IonApp, IonSplitPane, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import React from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';

import { LeftMenu } from './components/LeftMenu';
import { Login } from './pages/Login';
import { Page } from './pages/Page';
import { Register } from './pages/Register';
import {
  ReadingSection,
  WishlistSection,
  CompletedSection,
} from './utils/sections';
import './global.css';

export default function App() {
  return (
    <IonApp>
      <IonReactRouter>
        <IonRouterOutlet>
          <Switch>
            {/* 🔓 AUTH ROUTES */}
            <Route path="/login" component={Login} exact />
            <Route path="/register" component={Register} exact />

            {/* 🔐 APP ROUTES */}
            <Route>
              <IonSplitPane when="md" contentId="main-content">
                <LeftMenu />
                <IonRouterOutlet id="main-content">
                  <Switch>
                    <Route
                      path="/home"
                      component={() => <Page title="Home" />}
                    />
                    <Route path="/reading" component={ReadingSection} />
                    <Route path="/wishlist" component={WishlistSection} />
                    <Route path="/completed" component={CompletedSection} />
                    <Route
                      path="/settings"
                      component={() => <Page title="Settings" />}
                    />
                    <Redirect exact from="/" to="/home" />
                    <Route component={() => <Page title="Not Found" />} />
                  </Switch>
                </IonRouterOutlet>
              </IonSplitPane>
            </Route>
          </Switch>
        </IonRouterOutlet>
      </IonReactRouter>
    </IonApp>
  );
}
