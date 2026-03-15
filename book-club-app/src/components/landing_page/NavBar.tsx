import {
  IonToolbar,
  IonButtons,
  IonBackButton,
  IonTitle,
  IonButton,
} from '@ionic/react';
import { useHistory, useLocation } from 'react-router-dom';

export function LandingNavBar() {
  const history = useHistory();
  const location = useLocation();
  const isRegisterPage = location.pathname === '/register';
  const isLoginPage = location.pathname === '/login';
  return (
    <>
      <IonToolbar>
        <IonButtons slot="start">
          <IonBackButton />
        </IonButtons>
        <IonTitle className="nav-logo" onClick={() => history.push('/')}>
          📚 BookClub
        </IonTitle>

        <IonButtons slot="end">
          {isRegisterPage ? (
            <IonButton
              fill="clear"
              className="nav-login"
              onClick={() => history.push('/login')}
            >
              <span className="desktop-text">Already a member? Sign in</span>
              <span className="mobile-text">Sign in</span>
            </IonButton>
          ) : isLoginPage ? (
            <IonButton
              className="nav-login "
              onClick={() => history.push('/register')}
            >
              <span className="desktop-text">New here? Create an account</span>
              <span className="mobile-text">Sign up</span>
            </IonButton>
          ) : (
            <IonButton
              className="btn-nav"
              onClick={() => history.push('/register')}
            >
              Sign Up
            </IonButton>
          )}
        </IonButtons>
      </IonToolbar>
    </>
  );
}
