// src/pages/NotFound.tsx
import { IonPage, IonContent, IonButton } from '@ionic/react';
import { useHistory } from 'react-router-dom';

export function NotFound() {
  const history = useHistory();

  return (
    <IonContent className="body-bg ion-padding ion-text-center">
      <div className="padding-top padding-side">
        <h1>404</h1>
        <h2>Page not found</h2>
        <p>The page you're looking for doesn't exist.</p>
      </div>

      <IonButton onClick={() => history.push('/app')}>Go back to app</IonButton>
    </IonContent>
  );
}
