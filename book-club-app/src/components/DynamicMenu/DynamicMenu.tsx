import {
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonMenu,
  IonList,
  IonItem,
  IonRouterOutlet,
  IonSplitPane,
  IonLabel,
} from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import React, { useEffect, useState } from 'react';

// Example fetch function (replace with your actual API call)
async function fetchMenuItems() {
  // Replace with your backend endpoint
  const response = await fetch('/api/menu-items');
  if (!response.ok) throw new Error('Failed to fetch menu items');
  return response.json();
}

export function DynamicMenu() {
  const [items, setItems] = useState<Array<{ id: string; label: string }>>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMenuItems()
      .then(setItems)
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <IonMenu
      side="end"
      menuId="dynamic-menu"
      contentId="main-content"
      type="overlay"
    >
      <IonHeader>
        <IonToolbar>
          <IonTitle>Dynamic Menu</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent>
        {loading ? (
          <p>Loading...</p>
        ) : (
          <IonList>
            {items.map((item) => (
              <IonItem key={item.id} button>
                <IonLabel>{item.label}</IonLabel>
              </IonItem>
            ))}
          </IonList>
        )}
      </IonContent>
    </IonMenu>
  );
}
