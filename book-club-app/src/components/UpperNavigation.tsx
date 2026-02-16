import {
  IonHeader,
  IonToolbar,
  IonTitle,
  IonButton,
  IonButtons,
  IonAvatar,
  IonPopover,
  IonList,
  IonItem,
  useIonRouter,
} from '@ionic/react';
import { useState } from 'react';

import profilePic from '../assets/images/person-circle.svg';

export function UpperNavigation() {
  const [showOptions, setShowOptions] = useState(false);
  const [event, setEvent] = useState<any>(null);
  //const { logout } = useAuth();
  const router = useIonRouter();

  const handleLogout = () => {
    //logout();
    router.push('/login', 'root', 'replace');
  };

  return (
    <>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Hola</IonTitle>
          <IonButtons slot="end">
            <IonButton
              fill="clear"
              onClick={(e) => {
                setEvent(e.nativeEvent);
                setShowOptions(true);
              }}
            >
              <IonAvatar>
                <img src={profilePic} alt="Profile" />
              </IonAvatar>
            </IonButton>
          </IonButtons>
        </IonToolbar>
      </IonHeader>

      <IonPopover
        isOpen={showOptions}
        event={event}
        onDidDismiss={() => setShowOptions(false)}
      >
        <IonList>
          <IonItem button disabled>
            Profile
          </IonItem>
          <IonItem button disabled>
            Settings
          </IonItem>
          <IonItem button onClick={handleLogout}>
            Logout
          </IonItem>
        </IonList>
      </IonPopover>
    </>
  );
}
