import { IonApp, IonSplitPane, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import { Switch, Route, Redirect } from 'react-router-dom';

import { LeftMenu } from './components/LeftMenu';
import { AdminLayout } from './components/layout/AdminLayout';
import { useAuth } from './context/AuthContext';
import { ToastProvider } from './context/ToastContext';
import { LandingPage } from './pages/LandingPage';
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
  const { auth, initializing } = useAuth();

  if (initializing) {
    return <div>Loading...</div>;
  }

  return (
    <IonApp>
      <IonReactRouter>
        <IonRouterOutlet>
          <ToastProvider>
            <Switch>
              {/* 🏠 LANDING PAGE */}
              <Route
                exact
                path="/"
                render={() =>
                  auth.isAuthenticated ? (
                    auth.user?.role === 'admin' ? (
                      <Redirect to="/admin" />
                    ) : (
                      <Redirect to="/home" />
                    )
                  ) : (
                    <LandingPage />
                  )
                }
              />

              {/* 🔓 LOGIN */}
              <Route
                path="/login"
                render={() =>
                  auth.isAuthenticated ? <Redirect to="/home" /> : <Login />
                }
              />

              {/* 📝 REGISTER */}
              <Route
                path="/register"
                render={() =>
                  auth.isAuthenticated ? <Redirect to="/home" /> : <Register />
                }
              />

              {/* 🔐 ADMIN */}
              <Route path="/admin" component={AdminLayout} />

              {/* 🔐 AUTHENTICATED USER AREA */}
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

                      {/* fallback */}
                      <Route component={() => <Page title="Not Found" />} />
                    </Switch>
                  </IonRouterOutlet>
                </IonSplitPane>
              </Route>
            </Switch>
          </ToastProvider>
        </IonRouterOutlet>
      </IonReactRouter>
    </IonApp>
  );
}
