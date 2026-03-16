import {
  IonHeader,
  IonToolbar,
  IonButton,
  IonButtons,
  IonPopover,
  IonList,
  IonItem,
  useIonRouter,
  IonMenuButton,
} from '@ionic/react';
import { useState } from 'react';

import { useAuth } from '../../context/AuthContext';

export function UpperNavigation() {
  const { auth } = useAuth();
  const user = auth.user;
  const [showOptions, setShowOptions] = useState(false);
  const [event, setEvent] = useState<any>(null);
  const { logout } = useAuth();
  const router = useIonRouter();
  const initials = auth.user?.username
    ?.split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase();

  const handleLogout = () => {
    logout();
    router.push('/login', 'root', 'replace');
  };

  return (
    <>
      <IonHeader>
        <IonToolbar>
          <IonButtons slot="start">
            <IonMenuButton></IonMenuButton>
          </IonButtons>
          <IonButtons slot="end">
            <IonButton
              className="user-chip"
              fill="clear"
              onClick={(e) => {
                setEvent(e.nativeEvent);
                setShowOptions(true);
              }}
            >
              <div className="user-avatar" aria-hidden="true">
                {initials}
              </div>
              <span className="user-name">{user?.username}</span>
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
