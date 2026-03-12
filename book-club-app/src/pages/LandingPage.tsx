import {
  IonPage,
  IonHeader,
  IonToolbar,
  IonButtons,
  IonBackButton,
  IonTitle,
  IonButton,
  IonContent,
} from '@ionic/react';
import { useHistory } from 'react-router-dom';

import { AppFeatures } from '../components/landing_page/AppFeatures';
import { Hero } from '../components/landing_page/Hero';
import './LandingPage.css';

export function LandingPage() {
  const history = useHistory();
  return (
    <IonPage>
      {/* Nav */}
      <IonHeader>
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
      </IonHeader>

      <IonContent>
        <Hero />

        <AppFeatures />
      </IonContent>
    </IonPage>
  );
}
