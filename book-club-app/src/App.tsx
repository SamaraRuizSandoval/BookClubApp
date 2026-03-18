import { IonApp, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import { Switch, Route, Redirect } from 'react-router-dom';

import { AdminLayout } from './components/layout/AdminLayout';
import { UserLayout } from './components/layout/UserLayout';
import { useAuth } from './context/AuthContext';
import { ToastProvider } from './context/ToastContext';
import { LandingPage } from './pages/LandingPage';
import { Login } from './pages/Login';
import { NotFound } from './pages/Not found';
import { Register } from './pages/Register';
import './global.css';

export default function App() {
  const { auth, initializing } = useAuth();

  if (initializing) {
    return <div>Loading...</div>;
  }

  function RedirectIfAuthenticated({ children }: { children: JSX.Element }) {
    const { auth } = useAuth();

    if (auth.isAuthenticated) {
      return auth.user?.role === 'admin' ? (
        <Redirect to="/admin" />
      ) : (
        <Redirect to="/app" />
      );
    }

    return children;
  }

  return (
    <IonApp>
      <IonReactRouter>
        <ToastProvider>
          <IonRouterOutlet>
            <Switch>
              {/* 🏠 LANDING PAGE */}
              <Route
                exact
                path="/"
                render={() => (
                  <RedirectIfAuthenticated>
                    <LandingPage />
                  </RedirectIfAuthenticated>
                )}
              />

              {/* 🔓 LOGIN */}
              <Route
                path="/login"
                render={() => (
                  <RedirectIfAuthenticated>
                    <Login />
                  </RedirectIfAuthenticated>
                )}
              />

              {/* 📝 REGISTER */}
              <Route
                path="/register"
                render={() => (
                  <RedirectIfAuthenticated>
                    <Register />
                  </RedirectIfAuthenticated>
                )}
              />

              {/* 🔐 ADMIN */}
              <Route path="/admin" component={AdminLayout} />

              {/* 🔐 AUTHENTICATED USER AREA */}
              <Route path="/app" component={UserLayout} />

              <Route component={NotFound} />
            </Switch>
          </IonRouterOutlet>
        </ToastProvider>
      </IonReactRouter>
    </IonApp>
  );
}
