import { IonApp, IonSplitPane, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import React, { useEffect, useState } from 'react';
import { Switch, Route, Redirect, useHistory } from 'react-router-dom';

import { LeftMenu } from './components/LeftMenu';
import { useAuth } from './context/AuthContext';
import { AdminDashboard } from './pages/AdminDashboard';
import { Login } from './pages/Login';
import { Page } from './pages/Page';
import { Register } from './pages/Register';
import { AuthState } from './types/auth';
import {
  ReadingSection,
  WishlistSection,
  CompletedSection,
} from './utils/sections';
import './global.css';

export default function App() {
  const { auth, initializing } = useAuth();

  if (initializing) {
    return <div>Loading...</div>;
  }

  return (
    <IonApp>
      <IonReactRouter>
        <IonRouterOutlet>
          <Switch>
            {/* 🏠 LANDING */}
            <Route
              exact
              path="/"
              render={() =>
                auth.isAuthenticated ? (
                  auth.user?.role === 'admin' ? (
                    <Redirect to="/dashboard" />
                  ) : (
                    <Redirect to="/home" />
                  )
                ) : (
                  <Redirect to="/login" />
                )
              }
            />

            {/* 🔓 LOGIN */}
            <Route
              path="/login"
              render={() =>
                auth.isAuthenticated ? <Redirect to="/" /> : <Login />
              }
            />

            <Route path="/register" component={Register} exact />

            {/* 🔐 ADMIN DASHBOARD */}
            <Route path="/dashboard" component={AdminDashboard} exact />

            {/* 🔐 USER LAYOUT (SplitPane with LeftMenu) */}
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
