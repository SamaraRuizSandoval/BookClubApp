import {
  IonPage,
  IonContent,
  IonCard,
  IonCardHeader,
  IonCardContent,
  IonCardTitle,
  IonItem,
  IonInput,
  IonButton,
} from '@ionic/react';

export function Login() {
  return (
    <IonPage>
      <IonContent>
        <div className="flex-center">
          <IonCard>
            <IonCardHeader>
              <IonCardTitle>Login</IonCardTitle>
            </IonCardHeader>

            <IonCardContent>
              <IonItem>
                <IonInput placeholder="Username" type="text" />
              </IonItem>
              <IonItem>
                <IonInput placeholder="Password" type="password" />
              </IonItem>
              <IonButton
                color="primary"
                shape="round"
                size="default"
                expand="full"
                className="primary-button"
              >
                Login with username
              </IonButton>
            </IonCardContent>
          </IonCard>
        </div>
      </IonContent>
    </IonPage>
  );
}
