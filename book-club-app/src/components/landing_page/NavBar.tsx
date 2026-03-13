import {
  IonToolbar,
  IonButtons,
  IonBackButton,
  IonTitle,
  IonButton,
  IonContent,
} from '@ionic/react';
import { useHistory } from 'react-router-dom';

export function LandingNavBar() {
  const history = useHistory();
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
          <IonButton
            className="btn-nav"
            onClick={() => history.push('/register')}
          >
            Sign Up
          </IonButton>
        </IonButtons>
      </IonToolbar>
    </>
  );
}
