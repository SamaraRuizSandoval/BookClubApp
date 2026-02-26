import {
  IonContent,
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonText,
} from '@ionic/react';
import React from 'react';

export function AdminBooks() {
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Admin Books</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent className="ion-padding">
        <IonText>
          <h2>Admin Books Section</h2>
          <p>This is where admins can manage books.</p>
        </IonText>
      </IonContent>
    </IonPage>
  );
}
