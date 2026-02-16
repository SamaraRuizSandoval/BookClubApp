import { IonApp, IonSplitPane, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import React, { useEffect, useState } from 'react';
import { Switch, Route, Redirect, useHistory } from 'react-router-dom';

import { LeftMenu } from './components/LeftMenu';
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
  const [initializing, setInitializing] = useState(true);
  const [auth, setAuth] = useState<AuthState>({
    token: '',
    user: null,
    isAuthenticated: false,
  });

  useEffect(() => {
    const storedToken = localStorage.getItem('authToken');

    if (!storedToken) {
      setInitializing(false);
      return;
    }

    console.log('Found token in localStorage:', storedToken);
    fetch(`${import.meta.env.VITE_API_BASE_URL}/me`, {
      headers: {
        Authorization: `Bearer ${storedToken}`,
      },
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error('Failed to fetch user data');
        }
        return res.json();
      })
      .then((user) => {
        setAuth({
          token: storedToken,
          user: user,
          isAuthenticated: true,
        });
      })
      .catch(() => {
        localStorage.removeItem('token');
      })
      .finally(() => setInitializing(false));
  }, []);

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
                auth.isAuthenticated ? (
                  <Redirect to="/" />
                ) : (
                  <Login
                    onLoginSuccess={(token, user) =>
                      setAuth({
                        token,
                        user,
                        isAuthenticated: true,
                      })
                    }
                  />
                )
              }
            />

            <Route path="/register" component={Register} exact />

            {/* 🔐 ADMIN DASHBOARD */}
            <Route
              path="/dashboard"
              render={() =>
                auth.isAuthenticated && auth.user?.role === 'admin' ? (
                  <AdminDashboard />
                ) : (
                  <Redirect to="/login" />
                )
              }
            />

            {/* 🔐 USER LAYOUT (SplitPane with LeftMenu) */}
            <Route
              path="/home"
              render={() =>
                auth.isAuthenticated && auth.user?.role === 'user' ? (
                  <IonSplitPane when="md" contentId="main-content">
                    <LeftMenu />
                    <IonRouterOutlet id="main-content">
                      <Switch>
                        <Route
                          exact
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
                        <Route component={() => <Page title="Not Found" />} />
                      </Switch>
                    </IonRouterOutlet>
                  </IonSplitPane>
                ) : (
                  <Redirect to="/login" />
                )
              }
            />
          </Switch>
        </IonRouterOutlet>
      </IonReactRouter>
    </IonApp>
  );
}
