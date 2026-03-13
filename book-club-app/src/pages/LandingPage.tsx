import { IonPage, IonHeader, IonContent } from '@ionic/react';
import { useHistory } from 'react-router-dom';


import { AppFeatures } from '../components/landing_page/AppFeatures';
import { Hero } from '../components/landing_page/Hero';
import { LandingNavBar } from '../components/landing_page/NavBar';
import { useScrollReveal } from '../hooks/useScrollReveal';
import './LandingPage.css';

export function LandingPage() {
  useScrollReveal();
  const history = useHistory();

  return (
    <IonPage>
      {/* Nav */}
      <IonHeader>
        <LandingNavBar />
      </IonHeader>

      <IonContent>
        <Hero />

        <AppFeatures />
      </IonContent>
    </IonPage>
  );
}
