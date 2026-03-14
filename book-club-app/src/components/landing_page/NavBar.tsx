import {
  IonToolbar,
  IonButtons,
  IonBackButton,
  IonTitle,
  IonButton,
  IonContent,
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
              Already a member? Sign in
            </IonButton>
          ) : isLoginPage ? (
            <IonButton
              className="nav-login "
              onClick={() => history.push('/register')}
            >
              New here? Create an account
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
