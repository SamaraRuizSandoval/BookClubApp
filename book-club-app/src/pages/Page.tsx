import {
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonButtons,
  IonMenuButton,
} from '@ionic/react';
import { useLocation } from 'react-router-dom';

import { UpperNavigation } from '../components/UpperNavigation';

export function Page({ title }: { title: string }) {
  const location = useLocation();
  return (
    <IonPage>
      <IonHeader>
        <UpperNavigation />
      </IonHeader>
      <IonContent className="ion-padding">
        <h1 className="text-2xl font-semibold">{title}</h1>
        <p className="opacity-80">Current route: {location.pathname}</p>
      </IonContent>
    </IonPage>
  );
}
