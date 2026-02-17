import { IonContent } from '@ionic/react';
import React from 'react';

import { UpperNavigation } from '../components/UpperNavigation';

export function AdminDashboard() {
  return (
    <>
      <UpperNavigation />
      <IonContent>Admin Dashboard content</IonContent>
    </>
  );
}
