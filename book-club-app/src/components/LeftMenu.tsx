import {
  IonMenu,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonList,
  IonItem,
  IonMenuToggle,
  IonIcon,
  IonLabel,
} from '@ionic/react';
import {
  homeOutline,
  bookOutline,
  heartOutline,
  checkmarkDoneOutline,
  settingsOutline,
} from 'ionicons/icons';
import { useLocation } from 'react-router-dom';

const navItems = [
  { to: '/home', icon: homeOutline, label: 'Home' },
  { to: '/reading', icon: bookOutline, label: 'Reading' },
  { to: '/wishlist', icon: heartOutline, label: 'Wishlist' },
  { to: '/completed', icon: checkmarkDoneOutline, label: 'Completed' },
  { to: '/settings', icon: settingsOutline, label: 'Settings' },
];

export function LeftMenu() {
  const location = useLocation();

  return (
    <IonMenu contentId="main-content" type="overlay" side="start">
      <IonHeader>
        <IonToolbar>
          <IonTitle>BookClub</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent>
        <IonList inset>
          {navItems.map(({ to, icon, label }) => {
            const active = location.pathname.startsWith(to);
            return (
              <IonMenuToggle key={to} autoHide={false}>
                <IonItem
                  button
                  detail={false}
                  color={active ? 'primary' : undefined}
                  routerLink={to}
                  routerDirection="root"
                >
                  <IonIcon slot="start" icon={icon} />
                  <IonLabel>{label}</IonLabel>
                </IonItem>
              </IonMenuToggle>
            );
          })}
        </IonList>
      </IonContent>
    </IonMenu>
  );
}
